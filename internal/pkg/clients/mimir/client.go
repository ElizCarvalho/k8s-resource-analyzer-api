package mimir

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/errors"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
)

// Client é o cliente para o Mimir
type Client struct {
	baseURL    string
	httpClient *http.Client
	config     *ClientConfig
}

// ClientConfig contém as configurações do cliente
type ClientConfig struct {
	BaseURL     string
	Timeout     time.Duration
	ServiceName string
	Namespace   string
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

// NewClient cria uma nova instância do cliente Mimir
func NewClient(cfg *ClientConfig) *Client {
	logger.Info("Criando cliente Mimir",
		logger.NewField("base_url", cfg.BaseURL),
		logger.NewField("service_name", cfg.ServiceName),
		logger.NewField("namespace", cfg.Namespace),
		logger.NewField("timeout", cfg.Timeout),
	)

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
	logger.Info("Executing instant query",
		logger.NewField("base_url", c.baseURL),
		logger.NewField("query", query),
	)

	u, err := url.Parse(c.baseURL + "/prometheus/api/v1/query")
	if err != nil {
		logger.Error("Failed to parse URL", err,
			logger.NewField("url", c.baseURL),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to parse URL")
	}

	q := u.Query()
	q.Set("query", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		logger.Error("Failed to create request", err)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to create request")
	}

	req.Header.Set("X-Scope-OrgID", c.config.OrgID)
	logger.Info("Sending request",
		logger.NewField("url", req.URL.String()),
		logger.NewField("org_id", c.config.OrgID),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("Failed to execute request", err)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to execute request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response", err)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to read response")
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error("Unexpected status code", nil,
			logger.NewField("status_code", resp.StatusCode),
			logger.NewField("body", string(body)),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "unexpected status code")
	}

	var queryResp QueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		logger.Error("Failed to unmarshal response", err)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to unmarshal response")
	}

	if queryResp.Status != "success" {
		logger.Error("Query failed", nil,
			logger.NewField("status", queryResp.Status),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "query failed")
	}

	if len(queryResp.Data.Result) == 0 {
		logger.Info("No results found")
		return &types.QueryResult{
			Value:     0,
			Timestamp: time.Now(),
		}, nil
	}

	// Extrai o valor e timestamp
	result := queryResp.Data.Result[0]
	if len(result.Value) != 2 {
		logger.Error("Unexpected value format", nil,
			logger.NewField("value", result.Value),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "unexpected value format")
	}

	timestamp, ok := result.Value[0].(float64)
	if !ok {
		logger.Error("Invalid timestamp format", nil,
			logger.NewField("timestamp", result.Value[0]),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "invalid timestamp format")
	}

	value, ok := result.Value[1].(string)
	if !ok {
		logger.Error("Invalid value format", nil,
			logger.NewField("value", result.Value[1]),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "invalid value format")
	}

	floatValue, err := parseValue(value)
	if err != nil {
		logger.Error("Failed to parse value", err,
			logger.NewField("value", value),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to parse value")
	}

	logger.Info("Query executed successfully",
		logger.NewField("query", query),
		logger.NewField("value", floatValue),
	)

	return &types.QueryResult{
		Value:     floatValue,
		Timestamp: time.Unix(int64(timestamp), 0),
	}, nil
}

// QueryRange executa uma query de intervalo no Mimir
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error) {
	logger.Info("Executing range query",
		logger.NewField("base_url", c.baseURL),
		logger.NewField("query", query),
		logger.NewField("start", start),
		logger.NewField("end", end),
		logger.NewField("step", step),
	)

	u, err := url.Parse(c.baseURL + "/prometheus/api/v1/query_range")
	if err != nil {
		logger.Error("Failed to parse URL", err,
			logger.NewField("url", c.baseURL),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to parse URL")
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("start", start.Format(time.RFC3339))
	q.Set("end", end.Format(time.RFC3339))
	q.Set("step", step.String())
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		logger.Error("Failed to create request", err)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to create request")
	}

	req.Header.Set("X-Scope-OrgID", c.config.OrgID)
	logger.Info("Sending request",
		logger.NewField("url", req.URL.String()),
		logger.NewField("org_id", c.config.OrgID),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("Failed to execute request", err)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to execute request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response", err)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to read response")
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error("Unexpected status code", nil,
			logger.NewField("status_code", resp.StatusCode),
			logger.NewField("body", string(body)),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "unexpected status code")
	}

	var queryResp QueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		logger.Error("Failed to unmarshal response", err)
		return nil, errors.NewInvalidConfigurationError("mimir", "failed to unmarshal response")
	}

	if queryResp.Status != "success" {
		logger.Error("Query failed", nil,
			logger.NewField("status", queryResp.Status),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "query failed")
	}

	if len(queryResp.Data.Result) == 0 {
		logger.Info("No results found")
		return &types.QueryRangeResult{
			Values:    []types.QueryResult{},
			StartTime: start,
			EndTime:   end,
		}, nil
	}

	logger.Info("Processing results",
		logger.NewField("count", len(queryResp.Data.Result)),
	)

	// Processa os valores
	result := queryResp.Data.Result[0]
	values := make([]types.QueryResult, 0, len(result.Values))

	for _, v := range result.Values {
		if len(v) != 2 {
			logger.Error("Unexpected value format", nil,
				logger.NewField("value", v),
			)
			continue
		}

		timestamp, ok := v[0].(float64)
		if !ok {
			logger.Error("Invalid timestamp format", nil,
				logger.NewField("timestamp", v[0]),
			)
			continue
		}

		value, ok := v[1].(string)
		if !ok {
			logger.Error("Invalid value format", nil,
				logger.NewField("value", v[1]),
			)
			continue
		}

		floatValue, err := parseValue(value)
		if err != nil {
			logger.Error("Failed to parse value", err,
				logger.NewField("value", value),
			)
			continue
		}

		values = append(values, types.QueryResult{
			Value:     floatValue,
			Timestamp: time.Unix(int64(timestamp), 0),
		})
	}

	logger.Info("Range query executed successfully",
		logger.NewField("values_count", len(values)),
	)

	return &types.QueryRangeResult{
		Values:    values,
		StartTime: start,
		EndTime:   end,
	}, nil
}

// CheckConnection verifica a conexão com o Mimir
func (c *Client) CheckConnection(ctx context.Context) error {
	logger.Info("Verifying connection with Mimir",
		logger.NewField("base_url", c.baseURL),
	)

	u, err := url.Parse(c.baseURL + "/prometheus/api/v1/query")
	if err != nil {
		logger.Error("Failed to parse URL", err,
			logger.NewField("base_url", c.baseURL),
		)
		return errors.NewInvalidConfigurationError("mimir", "failed to parse URL")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		logger.Error("Failed to create request", err,
			logger.NewField("url", u.String()),
		)
		return errors.NewInvalidConfigurationError("mimir", "failed to create request")
	}

	req.Header.Set("X-Scope-OrgID", c.config.OrgID)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("Failed to execute request", err,
			logger.NewField("url", req.URL.String()),
		)
		return errors.NewInvalidConfigurationError("mimir", "failed to execute request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Unexpected status code", nil,
			logger.NewField("status_code", resp.StatusCode),
		)
		return errors.NewInvalidConfigurationError("mimir", "unexpected status code")
	}

	logger.Info("Successfully connected to Mimir")
	return nil
}
