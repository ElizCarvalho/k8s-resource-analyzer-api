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
	cfg := mimir.NewConfig()
	cfg.BaseURL = server.URL

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
		query := mimir.GetDeploymentCPUUsageQuery("default", "test-deployment")
		result, err := client.Query(ctx, query)
		if err != nil {
			t.Fatalf("Erro ao consultar métricas de CPU: %v", err)
		}

		if result == nil {
			t.Fatal("Resultado não deveria ser nil")
		}

		if result.Status != "success" {
			t.Errorf("Status esperado 'success', obtido '%s'", result.Status)
		}

		if len(result.Data.Result) == 0 {
			t.Error("Deveria ter encontrado métricas")
		}
	})

	// Testar consulta de métricas de memória
	t.Run("Deve consultar métricas de memória do deployment", func(t *testing.T) {
		query := mimir.GetDeploymentMemoryUsageQuery("default", "test-deployment")
		result, err := client.Query(ctx, query)
		if err != nil {
			t.Fatalf("Erro ao consultar métricas de memória: %v", err)
		}

		if result == nil {
			t.Fatal("Resultado não deveria ser nil")
		}

		if result.Status != "success" {
			t.Errorf("Status esperado 'success', obtido '%s'", result.Status)
		}

		if len(result.Data.Result) == 0 {
			t.Error("Deveria ter encontrado métricas")
		}
	})

	// Testar consulta com range de tempo
	t.Run("Deve consultar métricas com range de tempo", func(t *testing.T) {
		query := mimir.GetDeploymentCPUUsageQuery("default", "test-deployment")
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

		if result.Status != "success" {
			t.Errorf("Status esperado 'success', obtido '%s'", result.Status)
		}

		if len(result.Data.Result) == 0 {
			t.Error("Deveria ter encontrado métricas")
		}
	})
}

func TestMimirClientRetries(t *testing.T) {
	// Configurar servidor de teste
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= 2 { // Falha nas duas primeiras tentativas
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Sucesso na terceira tentativa
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[]}}`))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	// Criar configuração com retries
	cfg := mimir.NewConfig()
	cfg.BaseURL = server.URL
	cfg.Retry.MaxRetries = 3
	cfg.Retry.InitialBackoff = 100 * time.Millisecond

	client := mimir.NewClient(cfg)
	ctx := context.Background()

	// Executar query
	result, err := client.Query(ctx, "up")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Esperado 3 tentativas, obtido %d", attempts)
	}

	if result.Status != "success" {
		t.Errorf("Status esperado 'success', obtido '%s'", result.Status)
	}
}

func TestMimirClientRateLimit(t *testing.T) {
	// Configurar servidor de teste
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[]}}`))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	cfg := mimir.NewConfig()
	cfg.BaseURL = server.URL
	client := mimir.NewClient(cfg)
	ctx := context.Background()

	start := time.Now()
	_, err := client.Query(ctx, "up")
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	if duration < time.Second {
		t.Error("Rate limit não respeitado")
	}
}

func TestMimirClientCircuitBreaker(t *testing.T) {
	// Configurar servidor de teste que sempre falha
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := mimir.NewConfig()
	cfg.BaseURL = server.URL
	cfg.CircuitBreaker.MaxFailures = 2
	cfg.CircuitBreaker.ResetTimeout = 1 * time.Second
	cfg.Retry.MaxRetries = 1 // Apenas uma tentativa para facilitar o teste

	client := mimir.NewClient(cfg)
	ctx := context.Background()

	// Primeira chamada - deve falhar mas circuit breaker ainda fechado
	_, err1 := client.Query(ctx, "up")
	if err1 == nil {
		t.Error("Esperado erro na primeira chamada")
	}

	// Segunda chamada - deve falhar e abrir o circuit breaker
	_, err2 := client.Query(ctx, "up")
	if err2 == nil {
		t.Error("Esperado erro na segunda chamada")
	}

	// Terceira chamada - deve ser rejeitada pelo circuit breaker
	_, err3 := client.Query(ctx, "up")
	if err3 == nil {
		t.Error("Esperado erro de circuit breaker")
	} else if err3.Error() != "circuit breaker aberto" {
		t.Errorf("Erro esperado 'circuit breaker aberto', obtido '%s'", err3.Error())
	}

	// Aguardar reset do circuit breaker
	time.Sleep(2 * time.Second)

	// Quarta chamada - deve tentar novamente (half-open)
	_, err4 := client.Query(ctx, "up")
	if err4 == nil {
		t.Error("Esperado erro após reset do circuit breaker")
	}
}

func TestMimirClientContextTimeout(t *testing.T) {
	// Configurar servidor de teste com delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[]}}`))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	cfg := mimir.NewConfig()
	cfg.BaseURL = server.URL
	cfg.Timeouts.Query = 1 * time.Second

	client := mimir.NewClient(cfg)
	ctx := context.Background()

	_, err := client.Query(ctx, "up")
	if err == nil {
		t.Error("Esperado erro de timeout")
	}
}
