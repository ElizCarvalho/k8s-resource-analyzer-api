package collector

import (
	"context"
	"testing"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/clients/k8s"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/clients/mimir"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name     string
	testFunc func(*K8sMimirCollector) error
}

func TestK8sMimirCollector(t *testing.T) {
	collector := &K8sMimirCollector{
		K8sClient:   &k8s.Client{},
		MimirClient: &mimir.Client{},
	}
	ctx := context.Background()

	tests := []testCase{
		{
			name: "GetDeploymentMetrics",
			testFunc: func(c *K8sMimirCollector) error {
				_, err := c.GetDeploymentMetrics(ctx, "default", "test-deployment")
				return err
			},
		},
		{
			name: "GetDeploymentConfig",
			testFunc: func(c *K8sMimirCollector) error {
				_, err := c.GetDeploymentConfig(ctx, "default", "test-deployment")
				return err
			},
		},
		{
			name: "Query",
			testFunc: func(c *K8sMimirCollector) error {
				_, err := c.Query(ctx, "test_query")
				return err
			},
		},
		{
			name: "QueryRange",
			testFunc: func(c *K8sMimirCollector) error {
				now := time.Now()
				_, err := c.QueryRange(ctx, "test_query", now.Add(-time.Hour), now, time.Minute)
				return err
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.testFunc(collector)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "not implemented")
		})
	}
}
