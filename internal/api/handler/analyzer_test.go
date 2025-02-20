package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockResourceAnalyzer implementa a interface ResourceAnalyzer para testes
type MockResourceAnalyzer struct {
	GetMetricsFunc       func(ctx context.Context, namespace, deployment string, period time.Duration) (*types.MetricsResponse, error)
	GetTrendsFunc        func(ctx context.Context, namespace, deployment string, period time.Duration) (*types.TrendsResponse, error)
	AnalyzeResourcesFunc func(current *types.CurrentMetrics, historical *types.HistoricalMetrics) *types.ResourceAnalysis
	CalculateCostsFunc   func(ctx context.Context, current *types.CurrentMetrics, analysis *types.ResourceRecommendationAnalysis) (*types.CostAnalysis, error)
	GenerateAlertsFunc   func(current *types.CurrentMetrics, historical *types.HistoricalMetrics) []types.Alert
}

func (m *MockResourceAnalyzer) GetMetrics(ctx context.Context, namespace, deployment string, period time.Duration) (*types.MetricsResponse, error) {
	if m.GetMetricsFunc != nil {
		return m.GetMetricsFunc(ctx, namespace, deployment, period)
	}
	return nil, nil
}

func (m *MockResourceAnalyzer) GetTrends(ctx context.Context, namespace, deployment string, period time.Duration) (*types.TrendsResponse, error) {
	if m.GetTrendsFunc != nil {
		return m.GetTrendsFunc(ctx, namespace, deployment, period)
	}
	return nil, nil
}

func (m *MockResourceAnalyzer) AnalyzeResources(current *types.CurrentMetrics, historical *types.HistoricalMetrics) *types.ResourceAnalysis {
	if m.AnalyzeResourcesFunc != nil {
		return m.AnalyzeResourcesFunc(current, historical)
	}
	return nil
}

func (m *MockResourceAnalyzer) CalculateCosts(ctx context.Context, current *types.CurrentMetrics, analysis *types.ResourceRecommendationAnalysis) (*types.CostAnalysis, error) {
	if m.CalculateCostsFunc != nil {
		return m.CalculateCostsFunc(ctx, current, analysis)
	}
	return nil, nil
}

func (m *MockResourceAnalyzer) GenerateAlerts(current *types.CurrentMetrics, historical *types.HistoricalMetrics) []types.Alert {
	if m.GenerateAlertsFunc != nil {
		return m.GenerateAlertsFunc(current, historical)
	}
	return nil
}

func TestAnalyzerHandler_AnalyzeResources(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		namespace      string
		deployment     string
		period         string
		setupMock      func(*MockResourceAnalyzer)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:       "Sucesso - Análise completa",
			namespace:  "default",
			deployment: "test-app",
			period:     "24h",
			setupMock: func(m *MockResourceAnalyzer) {
				m.GetMetricsFunc = func(ctx context.Context, namespace, deployment string, period time.Duration) (*types.MetricsResponse, error) {
					return &types.MetricsResponse{
						Current: &types.CurrentMetrics{
							CPU:    &types.ResourceMetrics{Usage: 500, Request: 1000},
							Memory: &types.ResourceMetrics{Usage: 256, Request: 512},
						},
						Historical: &types.HistoricalMetrics{},
					}, nil
				}
				m.GetTrendsFunc = func(ctx context.Context, namespace, deployment string, period time.Duration) (*types.TrendsResponse, error) {
					return &types.TrendsResponse{
						CPU:    &types.TrendMetrics{Trend: 0.5},
						Memory: &types.TrendMetrics{Trend: 0.3},
					}, nil
				}
				m.AnalyzeResourcesFunc = func(current *types.CurrentMetrics, historical *types.HistoricalMetrics) *types.ResourceAnalysis {
					return &types.ResourceAnalysis{Status: "normal"}
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response["current"])
				assert.NotNil(t, response["trends"])
				assert.NotNil(t, response["analysis"])
			},
		},
		{
			name:           "Erro - Namespace não informado",
			namespace:      "",
			deployment:     "test-app",
			period:         "24h",
			setupMock:      func(m *MockResourceAnalyzer) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:           "Erro - Período inválido",
			namespace:      "default",
			deployment:     "test-app",
			period:         "invalid",
			setupMock:      func(m *MockResourceAnalyzer) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "período inválido")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mock := &MockResourceAnalyzer{}
			tt.setupMock(mock)
			handler := NewAnalyzerHandler(mock)

			router := gin.New()
			router.GET("/resources/:deployment/analysis", handler.AnalyzeResources)

			// Criar request
			url := "/resources/" + tt.deployment + "/analysis?namespace=" + tt.namespace
			if tt.period != "" {
				url += "&period=" + tt.period
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			// Executar request
			router.ServeHTTP(w, req)

			// Verificar resposta
			tt.checkResponse(t, w)
		})
	}
}
