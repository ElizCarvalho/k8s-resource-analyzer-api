package k8s_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/clients/k8s"
)

const (
	testDeploymentName = "travelernotifierbyevent"
	testNamespace      = "default"
)

func TestGetDeploymentMetrics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Usar o kubeconfig do ambiente
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	kubeconfigPath := filepath.Join(homeDir, ".kube", "config")

	// Criar cliente usando o kubeconfig do ambiente
	cfg := &k8s.ClientConfig{
		KubeconfigPath: kubeconfigPath,
		InCluster:      false,
	}

	client, err := k8s.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Testar obtenção de métricas
	metrics, err := client.GetDeploymentMetrics(ctx, testNamespace, testDeploymentName)
	if err != nil {
		t.Fatalf("failed to get metrics: %v", err)
	}

	if metrics.Pods.Running == 0 {
		t.Error("expected pod metrics, got none")
	}

	// Validar valores de CPU e memória
	if metrics.CPU.Usage == 0 {
		t.Error("expected CPU metrics")
	}
	if metrics.Memory.Usage == 0 {
		t.Error("expected Memory metrics")
	}

	t.Logf("Métricas: CPU=%.2f, Memory=%.2f, Pods=%d",
		metrics.CPU.Usage,
		metrics.Memory.Usage,
		metrics.Pods.Running)
}
