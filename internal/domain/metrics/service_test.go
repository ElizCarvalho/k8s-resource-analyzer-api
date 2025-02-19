package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/pricing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockK8sClient é um mock do cliente Kubernetes
type MockK8sClient struct {
	mock.Mock
}

func (m *MockK8sClient) GetDeploymentMetrics(ctx context.Context, namespace, name string) (*types.K8sMetrics, error) {
	args := m.Called(ctx, namespace, name)
	return args.Get(0).(*types.K8sMetrics), args.Error(1)
}

func (m *MockK8sClient) GetDeploymentConfig(ctx context.Context, namespace, name string) (*types.K8sDeploymentConfig, error) {
	args := m.Called(ctx, namespace, name)
	return args.Get(0).(*types.K8sDeploymentConfig), args.Error(1)
}

// MockMimirClient é um mock do cliente Mimir
type MockMimirClient struct {
	mock.Mock
}

func (m *MockMimirClient) Query(ctx context.Context, query string) (*types.QueryResult, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*types.QueryResult), args.Error(1)
}

func (m *MockMimirClient) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error) {
	args := m.Called(ctx, query, start, end, step)
	return args.Get(0).(*types.QueryRangeResult), args.Error(1)
}

// MockPricingClient é um mock do cliente de preços
type MockPricingClient struct {
	mock.Mock
}

func (m *MockPricingClient) GetResourcePricing(ctx context.Context) (*pricing.ResourcePricing, error) {
	args := m.Called(ctx)
	return args.Get(0).(*pricing.ResourcePricing), args.Error(1)
}

func (m *MockPricingClient) GetExchangeRate(ctx context.Context) (*pricing.ExchangeRate, error) {
	args := m.Called(ctx)
	return args.Get(0).(*pricing.ExchangeRate), args.Error(1)
}

func TestGetCurrentMetrics(t *testing.T) {
	// Arrange
	k8sClient := new(MockK8sClient)
	mimirClient := new(MockMimirClient)
	pricingClient := pricing.NewClient(&pricing.Config{
		ExchangeURL: "http://mock.test",
		Timeout:     time.Second,
	})
	service := NewService(k8sClient, mimirClient, pricingClient)

	ctx := context.Background()
	params := types.MetricsParams{
		Namespace:  "default",
		Deployment: "test-app",
		Period:     time.Hour,
	}

	// Mock das métricas do Kubernetes
	k8sMetrics := &types.K8sMetrics{}
	k8sMetrics.CPU.Usage = 0.5
	k8sMetrics.CPU.Average = 0.4
	k8sMetrics.CPU.Peak = 0.6
	k8sMetrics.Memory.Usage = 1.5
	k8sMetrics.Memory.Average = 1.2
	k8sMetrics.Memory.Peak = 1.8
	k8sMetrics.Pods.Running = 3

	k8sClient.On("GetDeploymentMetrics", ctx, params.Namespace, params.Deployment).Return(k8sMetrics, nil)

	// Mock da configuração do deployment
	deployConfig := &types.K8sDeploymentConfig{}
	deployConfig.CPU.Request = 1.0
	deployConfig.CPU.Limit = 2.0
	deployConfig.Memory.Request = 2.0
	deployConfig.Memory.Limit = 4.0
	deployConfig.Pods.Replicas = 3
	deployConfig.Pods.MinReplicas = 2
	deployConfig.Pods.MaxReplicas = 5

	k8sClient.On("GetDeploymentConfig", ctx, params.Namespace, params.Deployment).Return(deployConfig, nil)

	// Act
	result, err := service.GetCurrentMetrics(ctx, params)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verifica métricas de CPU
	assert.Equal(t, k8sMetrics.CPU.Average, result.Current.CPU.Average)
	assert.Equal(t, k8sMetrics.CPU.Peak, result.Current.CPU.Peak)
	assert.Equal(t, k8sMetrics.CPU.Usage, result.Current.CPU.Usage)
	assert.Equal(t, deployConfig.CPU.Request, result.Current.CPU.Request)
	assert.Equal(t, deployConfig.CPU.Limit, result.Current.CPU.Limit)
	assert.Equal(t, float64(50), result.Current.CPU.Utilization) // 0.5/1.0 * 100

	// Verifica métricas de memória
	assert.Equal(t, k8sMetrics.Memory.Average, result.Current.Memory.Average)
	assert.Equal(t, k8sMetrics.Memory.Peak, result.Current.Memory.Peak)
	assert.Equal(t, k8sMetrics.Memory.Usage, result.Current.Memory.Usage)
	assert.Equal(t, deployConfig.Memory.Request, result.Current.Memory.Request)
	assert.Equal(t, deployConfig.Memory.Limit, result.Current.Memory.Limit)
	assert.Equal(t, float64(75), result.Current.Memory.Utilization) // 1.5/2.0 * 100

	// Verifica métricas de pods
	assert.Equal(t, k8sMetrics.Pods.Running, result.Current.Pods.Running)
	assert.Equal(t, deployConfig.Pods.Replicas, result.Current.Pods.Replicas)
	assert.Equal(t, deployConfig.Pods.MinReplicas, result.Current.Pods.MinReplicas)
	assert.Equal(t, deployConfig.Pods.MaxReplicas, result.Current.Pods.MaxReplicas)

	// Verifica metadados
	assert.Equal(t, params.Period, result.Metadata.TimeWindow)
	assert.NotZero(t, result.Metadata.CollectedAt)

	// Verifica se todos os mocks foram chamados como esperado
	k8sClient.AssertExpectations(t)
	mimirClient.AssertExpectations(t)
}

func TestCalculateUtilization(t *testing.T) {
	tests := []struct {
		name     string
		usage    float64
		request  float64
		expected float64
	}{
		{
			name:     "Utilização normal",
			usage:    0.5,
			request:  1.0,
			expected: 50.0,
		},
		{
			name:     "Sobre-utilização",
			usage:    2.0,
			request:  1.0,
			expected: 200.0,
		},
		{
			name:     "Request zero",
			usage:    1.0,
			request:  0.0,
			expected: 0.0,
		},
		{
			name:     "Sem utilização",
			usage:    0.0,
			request:  1.0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateUtilization(tt.usage, tt.request)
			assert.Equal(t, tt.expected, result)
		})
	}
}
