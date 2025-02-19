package mimir_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/mimir"
)

func TestMimirIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando teste de integração em modo short")
	}

	// Configurar servidor de teste
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"pod":"test-pod"},"value":[1613765411.781,"42.5"]}]}}`))
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
		query := "rate(container_cpu_usage_seconds_total{namespace=\"default\",pod=~\"test-deployment.*\"}[5m])"
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
		query := "rate(container_cpu_usage_seconds_total{namespace=\"default\",pod=~\"test-deployment.*\"}[5m])"
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

		if result.StartTime.IsZero() || result.EndTime.IsZero() {
			t.Error("Timestamps não deveriam estar vazios")
		}
	})
}
