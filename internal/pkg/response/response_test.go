package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		message     string
		data        interface{}
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:    "Sucesso com dados",
			message: "Operação realizada com sucesso",
			data: map[string]string{
				"key": "value",
			},
			checkResult: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, "Operação realizada com sucesso", response["message"])
				assert.NotNil(t, response["data"])
				data := response["data"].(map[string]interface{})
				assert.Equal(t, "value", data["key"])
			},
		},
		{
			name:    "Sucesso sem dados",
			message: "Operação realizada",
			data:    nil,
			checkResult: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, "Operação realizada", response["message"])
				assert.Nil(t, response["data"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Success(c, tt.message, tt.data)
			tt.checkResult(t, w)
		})
	}
}

func TestSuccessWithRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		message     string
		data        interface{}
		requestID   string
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "Sucesso com RequestID",
			message:   "Operação realizada com sucesso",
			requestID: "123-456",
			data: map[string]string{
				"key": "value",
			},
			checkResult: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, "Operação realizada com sucesso", response["message"])
				assert.Equal(t, "123-456", response["request_id"])
				assert.NotNil(t, response["data"])
			},
		},
		{
			name:      "Sucesso com RequestID sem dados",
			message:   "Operação realizada",
			requestID: "789-012",
			data:      nil,
			checkResult: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, "Operação realizada", response["message"])
				assert.Equal(t, "789-012", response["request_id"])
				assert.Nil(t, response["data"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			SuccessWithRequestID(c, tt.message, tt.data, tt.requestID)
			tt.checkResult(t, w)
		})
	}
}

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		status      int
		message     string
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:    "Erro Bad Request",
			status:  http.StatusBadRequest,
			message: "Parâmetros inválidos",
			checkResult: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, "Parâmetros inválidos", response["message"])
			},
		},
		{
			name:    "Erro Internal Server",
			status:  http.StatusInternalServerError,
			message: "Erro interno do servidor",
			checkResult: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, "Erro interno do servidor", response["message"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Error(c, tt.status, tt.message)
			tt.checkResult(t, w)
		})
	}
}

func TestErrorWithRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		status      int
		message     string
		requestID   string
		checkResult func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "Erro com RequestID",
			status:    http.StatusNotFound,
			message:   "Recurso não encontrado",
			requestID: "123-456",
			checkResult: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, "Recurso não encontrado", response["message"])
				assert.Equal(t, "123-456", response["request_id"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			ErrorWithRequestID(c, tt.status, tt.message, tt.requestID)
			tt.checkResult(t, w)
		})
	}
}
