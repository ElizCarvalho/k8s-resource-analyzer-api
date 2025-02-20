package pricing

import (
	"context"
	"fmt"
	"time"
)

// Client é o cliente para obter informações de preços
type Client struct {
	config *Config
}

// Config contém as configurações do cliente
type Config struct {
	ExchangeURL string
	Timeout     time.Duration
}

// ResourcePrices representa preços dos recursos
type ResourcePrices struct {
	CPU struct {
		PerCore float64
	}
	Memory struct {
		PerGB float64
	}
}

// ExchangeRate representa taxa de câmbio
type ExchangeRate struct {
	Rate         float64
	FromCurrency string
	ToCurrency   string
	UpdatedAt    string
}

// NewClient cria uma nova instância do cliente
func NewClient(cfg *Config) *Client {
	return &Client{
		config: cfg,
	}
}

// GetCurrentPrices retorna os preços atuais dos recursos
func (c *Client) GetCurrentPrices(ctx context.Context) (*ResourcePrices, error) {
	// Por enquanto, retorna preços fixos baseados no GCP
	prices := &ResourcePrices{}
	prices.CPU.PerCore = 0.005425  // USD por core/hora
	prices.Memory.PerGB = 0.000729 // USD por GB/hora
	return prices, nil
}

// GetExchangeRate retorna a taxa de câmbio entre duas moedas
func (c *Client) GetExchangeRate(ctx context.Context, from, to string) (*ExchangeRate, error) {
	if from != "USD" || to != "BRL" {
		return nil, fmt.Errorf("only USD->BRL conversion is supported")
	}

	return &ExchangeRate{
		Rate:         5.78,
		FromCurrency: from,
		ToCurrency:   to,
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}, nil
}
