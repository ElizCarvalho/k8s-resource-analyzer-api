package mimir

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"sync"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
)

// CircuitBreaker implementa o padrão circuit breaker
type CircuitBreaker struct {
	config        CircuitBreakerConfig
	failures      int
	lastFailure   time.Time
	state         string // closed, open, half-open
	halfOpenCalls int
	mu            sync.RWMutex
}

func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  "closed",
	}
}

func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	switch cb.state {
	case "open":
		if now.Sub(cb.lastFailure) > cb.config.ResetTimeout {
			cb.state = "half-open"
			cb.halfOpenCalls = 0
			return true
		}
		return false
	case "half-open":
		if cb.halfOpenCalls < cb.config.HalfOpenMaxCalls {
			cb.halfOpenCalls++
			return true
		}
		return false
	default: // closed
		return true
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	if cb.state == "half-open" {
		cb.state = "closed"
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.failures >= cb.config.MaxFailures {
		cb.state = "open"
	}
}

// State retorna o estado atual do circuit breaker
func (cb *CircuitBreaker) State() string {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Client é o cliente para interagir com o Mimir
type Client struct {
	config         *Config
	http           *http.Client
	circuitBreaker *CircuitBreaker
}

// NewClient cria um novo cliente Mimir
func NewClient(cfg *Config) *Client {
	return &Client{
		config: cfg,
		http: &http.Client{
			Timeout: cfg.Timeouts.Connect,
		},
		circuitBreaker: NewCircuitBreaker(cfg.CircuitBreaker),
	}
}

// addOrgHeader adiciona o header X-Scope-OrgID à requisição
func (c *Client) addOrgHeader(req *http.Request) {
	req.Header.Add("X-Scope-OrgID", c.config.OrgID)
}

// CheckConnection verifica se a conexão com o Mimir está ativa
// Se não estiver, tenta estabelecer um port-forward
func (c *Client) CheckConnection(ctx context.Context) error {
	log := logger.Logger

	// Tenta fazer uma requisição simples para verificar a conexão
	url := fmt.Sprintf("%s/prometheus/api/v1/query?query=up", c.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar request: %w", err)
	}
	c.addOrgHeader(req)

	resp, err := c.http.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Info().Msg("Conexão com Mimir não está ativa. Tentando estabelecer port-forward...")

		// Tenta estabelecer port-forward
		cmd := exec.CommandContext(ctx, "kubectl", "port-forward",
			fmt.Sprintf("svc/%s", c.config.ServiceName),
			"-n", c.config.Namespace,
			fmt.Sprintf("%s:%s", c.config.LocalPort, c.config.ServicePort))

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("erro ao iniciar port-forward: %w", err)
		}

		// Aguarda um pouco para o port-forward estabelecer
		time.Sleep(5 * time.Second)

		// Verifica novamente a conexão
		resp, err = c.http.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return fmt.Errorf("erro: não foi possível estabelecer conexão com o Mimir")
		}

		log.Info().Msg("Conexão com Mimir estabelecida com sucesso!")
	}

	return nil
}

// QueryResult representa o resultado de uma consulta ao Mimir
type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string   `json:"metric"`
			Value  []json.RawMessage   `json:"value"`
			Values [][]json.RawMessage `json:"values,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// doWithRetry executa uma requisição HTTP com retentativas e backoff exponencial
func (c *Client) doWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	log := logger.Logger.With().
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Logger()

	if !c.circuitBreaker.AllowRequest() {
		log.Warn().Msg("Circuit breaker aberto, requisição bloqueada")
		return nil, fmt.Errorf("circuit breaker aberto")
	}

	var lastErr error
	backoff := c.config.Retry.InitialBackoff

	for i := 0; i < c.config.Retry.MaxRetries; i++ {
		log.Debug().
			Int("attempt", i+1).
			Str("backoff", backoff.String()).
			Msg("Tentando executar requisição")

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = err
			log.Error().
				Err(err).
				Int("attempt", i+1).
				Msg("Erro na requisição")

			c.circuitBreaker.RecordFailure()
			if c.circuitBreaker.State() == "open" {
				return nil, fmt.Errorf("circuit breaker aberto")
			}

			time.Sleep(backoff)
			backoff *= 2
			if backoff > c.config.Retry.MaxBackoff {
				backoff = c.config.Retry.MaxBackoff
			}
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			log.Warn().
				Str("retry_after", retryAfter).
				Msg("Rate limit atingido")

			if retryAfter != "" {
				if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
					time.Sleep(seconds)
					continue
				}
			}
			time.Sleep(backoff)
			backoff *= 2
			if backoff > c.config.Retry.MaxBackoff {
				backoff = c.config.Retry.MaxBackoff
			}
			continue
		}

		if resp.StatusCode >= 500 {
			c.circuitBreaker.RecordFailure()
			if c.circuitBreaker.State() == "open" {
				return nil, fmt.Errorf("circuit breaker aberto")
			}

			if i < c.config.Retry.MaxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2
				if backoff > c.config.Retry.MaxBackoff {
					backoff = c.config.Retry.MaxBackoff
				}
				continue
			}
		}

		c.circuitBreaker.RecordSuccess()
		return resp, nil
	}

	c.circuitBreaker.RecordFailure()
	if c.circuitBreaker.State() == "open" {
		return nil, fmt.Errorf("circuit breaker aberto")
	}

	if lastErr != nil {
		return nil, fmt.Errorf("máximo de tentativas excedido: %w", lastErr)
	}
	return nil, fmt.Errorf("máximo de tentativas excedido")
}

// Query executa uma consulta PromQL no Mimir
func (c *Client) Query(ctx context.Context, query string) (*QueryResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeouts.Query)
	defer cancel()

	log := logger.Logger.With().
		Str("operation", "query").
		Str("query", query).
		Logger()

	log.Debug().Msg("Iniciando consulta PromQL")

	url := fmt.Sprintf("%s/prometheus/api/v1/query?query=%s", c.config.BaseURL, url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Erro ao criar request")
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}
	c.addOrgHeader(req)

	start := time.Now()
	resp, err := c.doWithRetry(ctx, req)
	duration := time.Since(start)

	log.Debug().
		Dur("duration", duration).
		Msg("Consulta PromQL finalizada")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.circuitBreaker.RecordFailure()
		if c.circuitBreaker.State() == "open" {
			return nil, fmt.Errorf("circuit breaker aberto")
		}
		return nil, fmt.Errorf("erro na resposta do Mimir: status %d", resp.StatusCode)
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &result, nil
}

// QueryRange executa uma consulta PromQL com range de tempo no Mimir
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*QueryResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeouts.QueryRange)
	defer cancel()

	log := logger.Logger.With().
		Str("operation", "query_range").
		Str("query", query).
		Str("start", start.Format(time.RFC3339)).
		Str("end", end.Format(time.RFC3339)).
		Str("step", step.String()).
		Logger()

	log.Debug().Msg("Iniciando consulta PromQL com range")

	url := fmt.Sprintf("%s/prometheus/api/v1/query_range?query=%s&start=%d&end=%d&step=%d",
		c.config.BaseURL,
		url.QueryEscape(query),
		start.Unix(),
		end.Unix(),
		int(step.Seconds()))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Erro ao criar request")
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}
	c.addOrgHeader(req)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na resposta do Mimir: status %d", resp.StatusCode)
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &result, nil
}
