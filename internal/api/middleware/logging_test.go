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
			name:           "Success - Simple Request",
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
			name:           "Error - Status 404",
			path:           "/not-found",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			setupRouter:    func(r *gin.Engine) {},
		},
		{
			name:           "Success - With query params",
			path:           "/test",
			method:         http.MethodGet,
			query:          "param1=value1&param2=value2",
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
			name: "With error",
			path: "/error",
			setupRouter: func(r *gin.Engine) {
				r.GET("/error", func(c *gin.Context) {
					_ = c.Error(assert.AnError)
					c.Status(http.StatusInternalServerError)
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name: "Without error",
			path: "/success",
			setupRouter: func(r *gin.Engine) {
				r.GET("/success", func(c *gin.Context) {
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
			name: "With panic",
			path: "/panic",
			setupRouter: func(r *gin.Engine) {
				r.GET("/panic", func(c *gin.Context) {
					panic("panic test")
				})
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Without panic",
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
