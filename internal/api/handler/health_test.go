package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockK8sClient implementa a interface K8sClient para testes
type MockK8sClient struct {
	CheckConnectionFunc func(ctx context.Context) error
}

func (m *MockK8sClient) CheckConnection(ctx context.Context) error {
	if m.CheckConnectionFunc != nil {
		return m.CheckConnectionFunc(ctx)
	}
	return nil
}

// MockMimirClient implementa a interface MimirClient para testes
type MockMimirClient struct {
	CheckConnectionFunc func(ctx context.Context) error
}

func (m *MockMimirClient) CheckConnection(ctx context.Context) error {
	if m.CheckConnectionFunc != nil {
		return m.CheckConnectionFunc(ctx)
	}
	return nil
}

func TestHealthHandler_Check(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMocks     func(*MockK8sClient, *MockMimirClient)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Sucesso - Todos os serviços saudáveis",
			setupMocks: func(k8s *MockK8sClient, mimir *MockMimirClient) {
				k8s.CheckConnectionFunc = func(ctx context.Context) error {
					return nil
				}
				mimir.CheckConnectionFunc = func(ctx context.Context) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				data := response["data"].(map[string]interface{})
				assert.Equal(t, "healthy", data["status"])
				assert.NotEmpty(t, data["timestamp"])
				assert.NotEmpty(t, data["uptime"])
				assert.NotNil(t, data["system"])
				assert.NotNil(t, data["dependencies"])

				deps := data["dependencies"].(map[string]interface{})
				k8s := deps["kubernetes"].(map[string]interface{})
				assert.Equal(t, "healthy", k8s["status"])
				assert.Equal(t, "conectado ao cluster", k8s["message"])

				mimir := deps["mimir"].(map[string]interface{})
				assert.Equal(t, "healthy", mimir["status"])
				assert.Equal(t, "conectado ao serviço", mimir["message"])
			},
		},
		{
			name: "Degradado - Kubernetes indisponível",
			setupMocks: func(k8s *MockK8sClient, mimir *MockMimirClient) {
				k8s.CheckConnectionFunc = func(ctx context.Context) error {
					return context.DeadlineExceeded
				}
				mimir.CheckConnectionFunc = func(ctx context.Context) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				data := response["data"].(map[string]interface{})
				assert.Equal(t, "degraded", data["status"])

				deps := data["dependencies"].(map[string]interface{})
				k8s := deps["kubernetes"].(map[string]interface{})
				assert.Equal(t, "unhealthy", k8s["status"])
				assert.NotEmpty(t, k8s["error"])
			},
		},
		{
			name: "Degradado - Mimir indisponível",
			setupMocks: func(k8s *MockK8sClient, mimir *MockMimirClient) {
				k8s.CheckConnectionFunc = func(ctx context.Context) error {
					return nil
				}
				mimir.CheckConnectionFunc = func(ctx context.Context) error {
					return context.DeadlineExceeded
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				data := response["data"].(map[string]interface{})
				assert.Equal(t, "degraded", data["status"])

				deps := data["dependencies"].(map[string]interface{})
				mimir := deps["mimir"].(map[string]interface{})
				assert.Equal(t, "unhealthy", mimir["status"])
				assert.NotEmpty(t, mimir["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			k8sMock := &MockK8sClient{}
			mimirMock := &MockMimirClient{}
			tt.setupMocks(k8sMock, mimirMock)

			handler := NewHealthHandler(k8sMock, mimirMock)

			router := gin.New()
			router.GET("/health", handler.Check)

			// Criar request
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()

			// Executar request
			router.ServeHTTP(w, req)

			// Verificar resposta
			tt.checkResponse(t, w)
		})
	}
}
