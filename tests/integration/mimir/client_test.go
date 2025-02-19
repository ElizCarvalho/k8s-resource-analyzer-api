package mimir_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/clients/mimir"
)

const (
	testDeploymentName = "travelernotifierbyevent"
	testNamespace      = "default"
)

func TestMimirIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de integração em modo short")
	}

	// Configurar servidor de teste
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var response string
		if r.URL.Path == "/prometheus/api/v1/query" {
			response = `{
				"status": "success",
				"data": {
					"resultType": "vector",
					"result": [
						{
							"metric": {
								"pod": "travelernotifierbyevent-7cb945c95d-hzhsw"
							},
							"value": [1613765411.781, "0.071"]
						}
					]
				}
			}`
		} else {
			response = `{
				"status": "success",
				"data": {
					"resultType": "matrix",
					"result": [
						{
							"metric": {
								"pod": "travelernotifierbyevent-7cb945c95d-hzhsw"
							},
							"values": [
								[1613765411.781, "0.071"],
								[1613765711.781, "0.075"]
							]
						}
					]
				}
			}`
		}

		_, err := w.Write([]byte(response))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	// Criar configuração do Mimir
	cfg := &mimir.ClientConfig{
		BaseURL: server.URL,
		Timeout: 30 * time.Second,
		OrgID:   "anonymous",
	}

	// Criar cliente
	client := mimir.NewClient(cfg)

	// Contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Testar conexão
	t.Run("Deve conectar ao Mimir", func(t *testing.T) {
		err := client.CheckConnection(ctx)
		if err != nil {
			t.Fatalf("Erro ao conectar ao Mimir: %v", err)
		}
	})

	// Testar consulta de métricas de CPU
	t.Run("Deve consultar métricas de CPU do deployment", func(t *testing.T) {
		query := mimir.GetDeploymentCPUUsageQuery(testNamespace, testDeploymentName)
		result, err := client.Query(ctx, query)
		if err != nil {
			t.Fatalf("Erro ao consultar métricas de CPU: %v", err)
		}

		if result == nil {
			t.Fatal("Resultado não deveria ser nil")
		}

		if result.Value == 0 {
			t.Error("Valor da métrica não deveria ser zero")
		}

		if result.Timestamp.IsZero() {
			t.Error("Timestamp não deveria estar vazio")
		}
	})

	// Testar consulta com range de tempo
	t.Run("Deve consultar métricas com range de tempo", func(t *testing.T) {
		query := mimir.GetDeploymentCPUUsageQuery(testNamespace, testDeploymentName)
		end := time.Now()
		start := end.Add(-1 * time.Hour)
		step := 5 * time.Minute

		result, err := client.QueryRange(ctx, query, start, end, step)
		if err != nil {
			t.Fatalf("Erro ao consultar métricas com range: %v", err)
		}

		if result == nil {
			t.Fatal("Resultado não deveria ser nil")
		}

		if len(result.Values) == 0 {
			t.Error("Deveria ter encontrado métricas")
		}
	})
}
