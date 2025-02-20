package collector_test

import (
	"context"
	"testing"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/resource/collector"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
	"github.com/stretchr/testify/assert"
)

// MockK8sClient é um mock do cliente Kubernetes
type MockK8sClient struct{}

func (m *MockK8sClient) GetDeploymentMetrics(ctx context.Context, namespace, name string) (*types.K8sMetrics, error) {
	metrics := &types.K8sMetrics{}
	metrics.CPU.Usage = 100
	metrics.CPU.Average = 90
	metrics.CPU.Peak = 110
	metrics.CPU.Utilization = 0.8
	metrics.Memory.Usage = 128 * 1024 * 1024 // 128Mi em bytes
	metrics.Memory.Average = 120 * 1024 * 1024
	metrics.Memory.Peak = 150 * 1024 * 1024
	metrics.Memory.Utilization = 0.7
	metrics.Pods.Running = 3
	metrics.Pods.Utilization = 0.6
	return metrics, nil
}

func (m *MockK8sClient) GetDeploymentConfig(ctx context.Context, namespace, name string) (*types.K8sDeploymentConfig, error) {
	config := &types.K8sDeploymentConfig{}
	config.CPU.Request = 200
	config.CPU.Limit = 400
	config.Memory.Request = 256 * 1024 * 1024 // 256Mi em bytes
	config.Memory.Limit = 512 * 1024 * 1024
	config.Pods.Replicas = 3
	config.Pods.MinReplicas = 2
	config.Pods.MaxReplicas = 5
	config.Pods.TargetCPU = 0.8
	config.ClusterName = "test-cluster"
	return config, nil
}

func (m *MockK8sClient) CheckConnection(ctx context.Context) error {
	return nil
}

// MockMimirClient é um mock do cliente Mimir
type MockMimirClient struct{}

func (m *MockMimirClient) Query(ctx context.Context, query string) (*types.QueryResult, error) {
	return &types.QueryResult{
		Value:     42.0,
		Timestamp: time.Now(),
	}, nil
}

func (m *MockMimirClient) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*types.QueryRangeResult, error) {
	return &types.QueryRangeResult{
		Values: []types.QueryResult{
			{Value: 42.0, Timestamp: start},
			{Value: 43.0, Timestamp: start.Add(step)},
		},
		StartTime: start,
		EndTime:   end,
	}, nil
}

func (m *MockMimirClient) CheckConnection(ctx context.Context) error {
	return nil
}

type testCase struct {
	name     string
	testFunc func(*collector.K8sMimirCollector) (interface{}, error)
	validate func(t *testing.T, result interface{})
}

func TestK8sMimirCollector(t *testing.T) {
	c := collector.NewK8sMimirCollector(&MockK8sClient{}, &MockMimirClient{})
	ctx := context.Background()
	now := time.Now()

	tests := []testCase{
		{
			name: "GetDeploymentMetrics",
			testFunc: func(c *collector.K8sMimirCollector) (interface{}, error) {
				return c.GetDeploymentMetrics(ctx, "default", "test-deployment")
			},
			validate: func(t *testing.T, result interface{}) {
				metrics := result.(*types.K8sMetrics)
				assert.Equal(t, float64(100), metrics.CPU.Usage)
				assert.Equal(t, float64(128*1024*1024), metrics.Memory.Usage)
				assert.Equal(t, 3, metrics.Pods.Running)
			},
		},
		{
			name: "GetDeploymentConfig",
			testFunc: func(c *collector.K8sMimirCollector) (interface{}, error) {
				return c.GetDeploymentConfig(ctx, "default", "test-deployment")
			},
			validate: func(t *testing.T, result interface{}) {
				config := result.(*types.K8sDeploymentConfig)
				assert.Equal(t, float64(200), config.CPU.Request)
				assert.Equal(t, float64(256*1024*1024), config.Memory.Request)
				assert.Equal(t, 3, config.Pods.Replicas)
				assert.Equal(t, "test-cluster", config.ClusterName)
			},
		},
		{
			name: "Query",
			testFunc: func(c *collector.K8sMimirCollector) (interface{}, error) {
				return c.Query(ctx, "test_query")
			},
			validate: func(t *testing.T, result interface{}) {
				queryResult := result.(*types.QueryResult)
				assert.Equal(t, 42.0, queryResult.Value)
			},
		},
		{
			name: "QueryRange",
			testFunc: func(c *collector.K8sMimirCollector) (interface{}, error) {
				return c.QueryRange(ctx, "test_query", now, now.Add(time.Hour), time.Minute)
			},
			validate: func(t *testing.T, result interface{}) {
				rangeResult := result.(*types.QueryRangeResult)
				assert.Equal(t, 42.0, rangeResult.Values[0].Value)
				assert.Equal(t, 43.0, rangeResult.Values[1].Value)
				assert.Equal(t, now, rangeResult.StartTime)
				assert.Equal(t, now.Add(time.Hour), rangeResult.EndTime)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.testFunc(c)
			assert.NoError(t, err)
			tc.validate(t, result)
		})
	}
}

func TestK8sMimirCollector_GetDeploymentMetrics(t *testing.T) {
	c := collector.NewK8sMimirCollector(&MockK8sClient{}, &MockMimirClient{})

	metrics, err := c.GetDeploymentMetrics(context.Background(), "default", "test-deployment")

	assert.NoError(t, err)
	assert.Equal(t, float64(100), metrics.CPU.Usage)
	assert.Equal(t, float64(128*1024*1024), metrics.Memory.Usage)
	assert.Equal(t, 3, metrics.Pods.Running)
}

func TestK8sMimirCollector_GetDeploymentConfig(t *testing.T) {
	c := collector.NewK8sMimirCollector(&MockK8sClient{}, &MockMimirClient{})

	config, err := c.GetDeploymentConfig(context.Background(), "default", "test-deployment")

	assert.NoError(t, err)
	assert.Equal(t, float64(200), config.CPU.Request)
	assert.Equal(t, float64(256*1024*1024), config.Memory.Request)
	assert.Equal(t, 3, config.Pods.Replicas)
	assert.Equal(t, "test-cluster", config.ClusterName)
}

func TestK8sMimirCollector_Query(t *testing.T) {
	c := collector.NewK8sMimirCollector(&MockK8sClient{}, &MockMimirClient{})

	result, err := c.Query(context.Background(), "test_query")

	assert.NoError(t, err)
	assert.Equal(t, 42.0, result.Value)
}

func TestK8sMimirCollector_QueryRange(t *testing.T) {
	c := collector.NewK8sMimirCollector(&MockK8sClient{}, &MockMimirClient{})

	now := time.Now()
	result, err := c.QueryRange(context.Background(), "test_query", now, now.Add(time.Hour), time.Minute)

	assert.NoError(t, err)
	assert.Equal(t, 42.0, result.Values[0].Value)
	assert.Equal(t, 43.0, result.Values[1].Value)
	assert.Equal(t, now, result.StartTime)
	assert.Equal(t, now.Add(time.Hour), result.EndTime)
}
