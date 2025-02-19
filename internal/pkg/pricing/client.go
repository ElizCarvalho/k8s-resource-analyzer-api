package pricing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client é o cliente para obter preços e taxas de câmbio
type Client struct {
	httpClient *http.Client
	config     *Config
}

// Config contém as configurações do cliente
type Config struct {
	ExchangeURL string
	Timeout     time.Duration
}

// ResourcePricing contém os preços dos recursos
type ResourcePricing struct {
	CPU    float64 // Preço por vCPU/hora
	Memory float64 // Preço por GB/hora
}

// ExchangeRate contém a taxa de câmbio
type ExchangeRate struct {
	Rate      float64
	Timestamp time.Time
}

// NewClient cria um novo cliente
func NewClient(cfg *Config) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		config: cfg,
	}
}

// GetResourcePricing obtém os preços dos recursos do GCP
func (c *Client) GetResourcePricing(ctx context.Context) (*ResourcePricing, error) {
	// Valores baseados no GCP (região us-central1, E2 standard)
	return &ResourcePricing{
		CPU:    0.021811, // $0.021811 por vCPU/hora
		Memory: 0.002923, // $0.002923 por GB/hora
	}, nil
}

// GetExchangeRate obtém a taxa de câmbio USD/BRL
func (c *Client) GetExchangeRate(ctx context.Context) (*ExchangeRate, error) {
	url := fmt.Sprintf("%s/latest?base=USD&symbols=BRL", c.config.ExchangeURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter taxa de câmbio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao obter taxa de câmbio: status %d", resp.StatusCode)
	}

	var result struct {
		Rates struct {
			BRL float64 `json:"BRL"`
		} `json:"rates"`
		Timestamp int64 `json:"timestamp"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &ExchangeRate{
		Rate:      result.Rates.BRL,
		Timestamp: time.Unix(result.Timestamp, 0),
	}, nil
}
