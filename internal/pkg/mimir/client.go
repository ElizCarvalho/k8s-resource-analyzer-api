package mimir

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
)

// Client é o cliente para interagir com o Mimir
type Client struct {
	baseURL    string
	httpClient *http.Client
	config     *ClientConfig
}

// ClientConfig contém as configurações para o cliente Mimir
type ClientConfig struct {
	BaseURL     string
	Timeout     time.Duration
	ServiceName string
	Namespace   string
	LocalPort   string
	ServicePort string
	OrgID       string
}

// QueryResponse representa a resposta de uma query do Mimir
type QueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
			Values [][]interface{}   `json:"values,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// Constantes para validação
var (
	// Lista branca de serviços permitidos
	allowedServices = map[string]bool{
		"lgtm-mimir-query-frontend": true,
		"mimir-query-frontend":      true,
		"mimir":                     true,
	}

	// Lista branca de namespaces permitidos
	allowedNamespaces = map[string]bool{
		"monitoring":    true,
		"observability": true,
	}

	// Portas permitidas
	minPort = 1024
	maxPort = 65535
)

// NewClient cria uma nova instância do cliente Mimir
func NewClient(cfg *ClientConfig) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	return &Client{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		config: cfg,
	}
}

func parseValue(value string) (float64, error) {
	if value == "NaN" || value == "Inf" || value == "-Inf" || value == "" {
		return 0, nil
	}
	var floatValue float64
	_, err := fmt.Sscanf(value, "%f", &floatValue)
	if err != nil {
		return 0, err
	}
	if floatValue != floatValue || floatValue < 0 || floatValue > 1e6 { // Verifica NaN e valores absurdos
		return 0, nil
	}
	return floatValue, nil
}

// Query executa uma query instantânea no Mimir
func (c *Client) Query(ctx context.Context, query string) (*types.QueryResult, error) {
	fmt.Printf("DEBUG: Mimir baseURL = %s\n", c.baseURL)

	u, err := url.Parse(c.baseURL + "/prometheus/api/v1/query")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("query", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Scope-OrgID", c.config.OrgID)
	fmt.Printf("DEBUG: Enviando requisição para %s com X-Scope-OrgID: %s\n", req.URL.String(), c.config.OrgID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("DEBUG: Erro ao executar requisição: %v\n", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("DEBUG: Erro ao ler resposta: %v\n", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("DEBUG: Status code inesperado: %d, body: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var queryResp QueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		fmt.Printf("DEBUG: Erro ao fazer unmarshal da resposta: %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if queryResp.Status != "success" {
		fmt.Printf("DEBUG: Query falhou: %s\n", string(body))
		return nil, fmt.Errorf("query failed: %s", string(body))
	}

	if len(queryResp.Data.Result) == 0 {
		fmt.Printf("DEBUG: Nenhum resultado encontrado\n")
		return &types.QueryResult{
			Value:     0,
			Timestamp: time.Now(),
		}, nil
	}

	// Extrai o valor e timestamp
	result := queryResp.Data.Result[0]
	if len(result.Value) != 2 {
		return nil, fmt.Errorf("unexpected value format: %v", result.Value)
	}

	timestamp, ok := result.Value[0].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid timestamp format: %v", result.Value[0])
	}

	value, ok := result.Value[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid value format: %v", result.Value[1])
	}

	floatValue, err := parseValue(value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse value: %w", err)
	}

	return &types.QueryResult{
		Value:     floatValue,
		Timestamp: time.Unix(int64(timestamp), 0),
	}, nil
}

// QueryRange executa uma query de intervalo no Mimir
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error) {
	u, err := url.Parse(c.baseURL + "/prometheus/api/v1/query_range")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("start", fmt.Sprintf("%d", start.Unix()))
	q.Set("end", fmt.Sprintf("%d", end.Unix()))
	q.Set("step", fmt.Sprintf("%ds", int(step.Seconds())))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Scope-OrgID", c.config.OrgID)
	fmt.Printf("DEBUG: Enviando requisição para %s com X-Scope-OrgID: %s\n", req.URL.String(), c.config.OrgID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("DEBUG: Erro ao executar requisição: %v\n", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("DEBUG: Erro ao ler resposta: %v\n", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("DEBUG: Status code inesperado: %d, body: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var queryResp QueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		fmt.Printf("DEBUG: Erro ao fazer unmarshal da resposta: %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if queryResp.Status != "success" {
		fmt.Printf("DEBUG: Query falhou: %s\n", string(body))
		return nil, fmt.Errorf("query failed: %s", string(body))
	}

	if len(queryResp.Data.Result) == 0 {
		fmt.Printf("DEBUG: Nenhum resultado encontrado\n")
		return &types.QueryRangeResult{
			Values:    []types.QueryResult{},
			StartTime: start,
			EndTime:   end,
		}, nil
	}

	fmt.Printf("DEBUG: Resposta recebida com sucesso. Número de resultados: %d\n", len(queryResp.Data.Result))

	// Processa os valores
	result := queryResp.Data.Result[0]
	values := make([]types.QueryResult, 0, len(result.Values))

	for _, v := range result.Values {
		if len(v) != 2 {
			continue
		}

		timestamp, ok := v[0].(float64)
		if !ok {
			continue
		}

		value, ok := v[1].(string)
		if !ok {
			continue
		}

		floatValue, err := parseValue(value)
		if err != nil {
			continue
		}

		values = append(values, types.QueryResult{
			Value:     floatValue,
			Timestamp: time.Unix(int64(timestamp), 0),
		})
	}

	return &types.QueryRangeResult{
		Values:    values,
		StartTime: start,
		EndTime:   end,
	}, nil
}

// CheckConnection verifica a conexão com o Mimir
func (c *Client) CheckConnection(ctx context.Context) error {
	healthURL := fmt.Sprintf("%s/ready", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar request de verificação: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao Mimir: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Mimir retornou status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) setupPortForward(ctx context.Context) error {
	// Validação dos inputs
	if err := validatePortForwardInputs(c.config); err != nil {
		return fmt.Errorf("inputs inválidos para port-forward: %w", err)
	}

	// Tenta estabelecer port-forward com valores validados e seguros
	// #nosec G204 -- Inputs são validados contra uma lista branca de valores permitidos
	cmd := exec.CommandContext(ctx, "kubectl", "port-forward",
		fmt.Sprintf("svc/%s", c.config.ServiceName),
		"-n", c.config.Namespace,
		fmt.Sprintf("%s:%s", c.config.LocalPort, c.config.ServicePort))

	// Executa o comando
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar port-forward: %w", err)
	}

	// Aguarda um pouco para o port-forward estabelecer
	time.Sleep(2 * time.Second)

	// Verifica se a conexão foi estabelecida
	healthURL := fmt.Sprintf("http://localhost:%s/ready", c.config.LocalPort)
	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar request de verificação: %w", err)
	}

	// Verifica novamente a conexão
	resp, err := c.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro: não foi possível estabelecer conexão com o Mimir")
	}

	return nil
}

// validatePortForwardInputs valida os inputs para o port-forward
func validatePortForwardInputs(cfg *ClientConfig) error {
	// Valida ServiceName usando lista branca
	if !allowedServices[cfg.ServiceName] {
		return fmt.Errorf("serviço não permitido: %s", cfg.ServiceName)
	}

	// Valida Namespace usando lista branca
	if !allowedNamespaces[cfg.Namespace] {
		return fmt.Errorf("namespace não permitido: %s", cfg.Namespace)
	}

	// Valida LocalPort
	localPort, err := strconv.Atoi(cfg.LocalPort)
	if err != nil || localPort < minPort || localPort > maxPort {
		return fmt.Errorf("porta local inválida: deve estar entre %d e %d", minPort, maxPort)
	}

	// Valida ServicePort
	servicePort, err := strconv.Atoi(cfg.ServicePort)
	if err != nil || servicePort < minPort || servicePort > maxPort {
		return fmt.Errorf("porta do serviço inválida: deve estar entre %d e %d", minPort, maxPort)
	}

	return nil
}
