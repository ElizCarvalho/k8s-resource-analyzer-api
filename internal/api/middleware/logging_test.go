package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	// Configura o modo de teste do Gin
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		method         string
		query          string
		expectedStatus int
		setupRouter    func(*gin.Engine)
	}{
		{
			name:           "Sucesso - Request simples",
			path:           "/test",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			setupRouter: func(r *gin.Engine) {
				r.GET("/test", func(c *gin.Context) {
					c.Status(http.StatusOK)
				})
			},
		},
		{
			name:           "Erro - Status 404",
			path:           "/nao-existe",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			setupRouter:    func(r *gin.Engine) {},
		},
		{
			name:           "Sucesso - Com query params",
			path:           "/test",
			method:         http.MethodGet,
			query:          "param1=valor1&param2=valor2",
			expectedStatus: http.StatusOK,
			setupRouter: func(r *gin.Engine) {
				r.GET("/test", func(c *gin.Context) {
					c.Status(http.StatusOK)
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(RequestLogger())
			tt.setupRouter(router)

			w := httptest.NewRecorder()
			url := tt.path
			if tt.query != "" {
				url += "?" + tt.query
			}
			req := httptest.NewRequest(tt.method, url, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestErrorLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		setupRouter    func(*gin.Engine)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Com erro",
			path: "/erro",
			setupRouter: func(r *gin.Engine) {
				r.GET("/erro", func(c *gin.Context) {
					_ = c.Error(assert.AnError)
					c.Status(http.StatusInternalServerError)
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name: "Sem erro",
			path: "/sucesso",
			setupRouter: func(r *gin.Engine) {
				r.GET("/sucesso", func(c *gin.Context) {
					c.Status(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(ErrorLogger())
			tt.setupRouter(router)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRecoveryLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		setupRouter    func(*gin.Engine)
		expectedStatus int
	}{
		{
			name: "Com panic",
			path: "/panic",
			setupRouter: func(r *gin.Engine) {
				r.GET("/panic", func(c *gin.Context) {
					panic("teste de panic")
				})
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Sem panic",
			path: "/normal",
			setupRouter: func(r *gin.Engine) {
				r.GET("/normal", func(c *gin.Context) {
					c.Status(http.StatusOK)
				})
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(RecoveryLogger())
			tt.setupRouter(router)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
