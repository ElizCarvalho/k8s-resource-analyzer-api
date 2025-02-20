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
	logger.Info("Executando query instantânea",
		logger.NewField("base_url", c.baseURL),
		logger.NewField("query", query),
	)

	u, err := url.Parse(c.baseURL + "/prometheus/api/v1/query")
	if err != nil {
		logger.Error("Erro ao fazer parse da URL", err,
			logger.NewField("base_url", c.baseURL),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao fazer parse da URL")
	}

	q := u.Query()
	q.Set("query", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		logger.Error("Erro ao criar requisição", err,
			logger.NewField("url", u.String()),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao criar requisição")
	}

	req.Header.Set("X-Scope-OrgID", c.config.OrgID)
	logger.Info("Enviando requisição",
		logger.NewField("url", req.URL.String()),
		logger.NewField("org_id", c.config.OrgID),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("Erro ao executar requisição", err,
			logger.NewField("url", req.URL.String()),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao executar requisição")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Erro ao ler resposta", err)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao ler resposta")
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error("Status code inesperado", nil,
			logger.NewField("status_code", resp.StatusCode),
			logger.NewField("body", string(body)),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "status code inesperado")
	}

	var queryResp QueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		logger.Error("Erro ao fazer unmarshal da resposta", err,
			logger.NewField("body", string(body)),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao fazer unmarshal da resposta")
	}

	if queryResp.Status != "success" {
		logger.Error("Query falhou", nil,
			logger.NewField("status", queryResp.Status),
			logger.NewField("body", string(body)),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "query falhou")
	}

	if len(queryResp.Data.Result) == 0 {
		logger.Info("Nenhum resultado encontrado")
		return &types.QueryResult{
			Value:     0,
			Timestamp: time.Now(),
		}, nil
	}

	// Extrai o valor e timestamp
	result := queryResp.Data.Result[0]
	if len(result.Value) != 2 {
		logger.Error("Formato de valor inesperado", nil,
			logger.NewField("value", result.Value),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "formato de valor inesperado")
	}

	timestamp, ok := result.Value[0].(float64)
	if !ok {
		logger.Error("Formato de timestamp inválido", nil,
			logger.NewField("timestamp", result.Value[0]),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "formato de timestamp inválido")
	}

	value, ok := result.Value[1].(string)
	if !ok {
		logger.Error("Formato de valor inválido", nil,
			logger.NewField("value", result.Value[1]),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "formato de valor inválido")
	}

	floatValue, err := parseValue(value)
	if err != nil {
		logger.Error("Erro ao fazer parse do valor", err,
			logger.NewField("value", value),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao fazer parse do valor")
	}

	logger.Info("Query executada com sucesso",
		logger.NewField("value", floatValue),
		logger.NewField("timestamp", time.Unix(int64(timestamp), 0)),
	)

	return &types.QueryResult{
		Value:     floatValue,
		Timestamp: time.Unix(int64(timestamp), 0),
	}, nil
}

// QueryRange executa uma query de intervalo no Mimir
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error) {
	logger.Info("Executando query com range",
		logger.NewField("base_url", c.baseURL),
		logger.NewField("query", query),
		logger.NewField("start", start),
		logger.NewField("end", end),
		logger.NewField("step", step),
	)

	u, err := url.Parse(c.baseURL + "/prometheus/api/v1/query_range")
	if err != nil {
		logger.Error("Erro ao fazer parse da URL", err,
			logger.NewField("base_url", c.baseURL),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao fazer parse da URL")
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("start", start.Format(time.RFC3339))
	q.Set("end", end.Format(time.RFC3339))
	q.Set("step", step.String())
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		logger.Error("Erro ao criar requisição", err,
			logger.NewField("url", u.String()),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao criar requisição")
	}

	req.Header.Set("X-Scope-OrgID", c.config.OrgID)
	logger.Info("Enviando requisição",
		logger.NewField("url", req.URL.String()),
		logger.NewField("org_id", c.config.OrgID),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("Erro ao executar requisição", err,
			logger.NewField("url", req.URL.String()),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao executar requisição")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Erro ao ler resposta", err)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao ler resposta")
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error("Status code inesperado", nil,
			logger.NewField("status_code", resp.StatusCode),
			logger.NewField("body", string(body)),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "status code inesperado")
	}

	var queryResp QueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		logger.Error("Erro ao fazer unmarshal da resposta", err,
			logger.NewField("body", string(body)),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "erro ao fazer unmarshal da resposta")
	}

	if queryResp.Status != "success" {
		logger.Error("Query falhou", nil,
			logger.NewField("status", queryResp.Status),
			logger.NewField("body", string(body)),
		)
		return nil, errors.NewInvalidConfigurationError("mimir", "query falhou")
	}

	if len(queryResp.Data.Result) == 0 {
		logger.Info("Nenhum resultado encontrado")
		return &types.QueryRangeResult{
			Values:    []types.QueryResult{},
			StartTime: start,
			EndTime:   end,
		}, nil
	}

	logger.Info("Processando resultados",
		logger.NewField("count", len(queryResp.Data.Result)),
	)

	// Processa os valores
	result := queryResp.Data.Result[0]
	values := make([]types.QueryResult, 0, len(result.Values))

	for _, v := range result.Values {
		if len(v) != 2 {
			logger.Error("Formato de valor inesperado", nil,
				logger.NewField("value", v),
			)
			continue
		}

		timestamp, ok := v[0].(float64)
		if !ok {
			logger.Error("Formato de timestamp inválido", nil,
				logger.NewField("timestamp", v[0]),
			)
			continue
		}

		value, ok := v[1].(string)
		if !ok {
			logger.Error("Formato de valor inválido", nil,
				logger.NewField("value", v[1]),
			)
			continue
		}

		floatValue, err := parseValue(value)
		if err != nil {
			logger.Error("Erro ao fazer parse do valor", err,
				logger.NewField("value", value),
			)
			continue
		}

		values = append(values, types.QueryResult{
			Value:     floatValue,
			Timestamp: time.Unix(int64(timestamp), 0),
		})
	}

	logger.Info("Query com range executada com sucesso",
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
	logger.Info("Verificando conexão com o Mimir",
		logger.NewField("base_url", c.baseURL),
	)

	u, err := url.Parse(c.baseURL + "/prometheus/api/v1/status/config")
	if err != nil {
		logger.Error("Erro ao fazer parse da URL", err,
			logger.NewField("base_url", c.baseURL),
		)
		return errors.NewInvalidConfigurationError("mimir", "erro ao fazer parse da URL")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		logger.Error("Erro ao criar requisição", err,
			logger.NewField("url", u.String()),
		)
		return errors.NewInvalidConfigurationError("mimir", "erro ao criar requisição")
	}

	req.Header.Set("X-Scope-OrgID", c.config.OrgID)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		logger.Error("Erro ao executar requisição", err,
			logger.NewField("url", req.URL.String()),
		)
		return errors.NewInvalidConfigurationError("mimir", "erro ao executar requisição")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Status code inesperado", nil,
			logger.NewField("status_code", resp.StatusCode),
		)
		return errors.NewInvalidConfigurationError("mimir", "status code inesperado")
	}

	logger.Info("Conexão com o Mimir estabelecida com sucesso")
	return nil
}
