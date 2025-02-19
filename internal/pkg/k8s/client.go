package k8s

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// DeploymentMetrics contém as métricas de um deployment
type DeploymentMetrics struct {
	Pods []PodMetrics
}

// PodMetrics contém as métricas de um pod
type PodMetrics struct {
	Name      string
	CPU       string
	Memory    string
	Timestamp time.Time
}

// Client é o cliente para interagir com o cluster Kubernetes
type Client struct {
	clientset     *kubernetes.Clientset
	metricsClient *versioned.Clientset
	config        *Config
}

// NewClient cria um novo cliente Kubernetes
func NewClient(cfg *Config) (*Client, error) {
	var k8sConfig *rest.Config
	var err error

	if cfg.InCluster {
		k8sConfig, err = rest.InClusterConfig()
	} else {
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", cfg.KubeconfigPath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	metricsClient, err := versioned.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	return &Client{
		clientset:     clientset,
		metricsClient: metricsClient,
		config:        cfg,
	}, nil
}

// GetDeploymentMetrics obtém métricas de um deployment específico
func (c *Client) GetDeploymentMetrics(ctx context.Context, name string) (*DeploymentMetrics, error) {
	// Obtém o deployment
	deployment, err := c.clientset.AppsV1().Deployments(c.config.Namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	// Obtém os pods do deployment
	labelSelector := metav1.FormatLabelSelector(deployment.Spec.Selector)
	pods, err := c.clientset.CoreV1().Pods(c.config.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	// Obtém métricas dos pods
	podMetrics := make([]PodMetrics, 0, len(pods.Items))
	for _, pod := range pods.Items {
		metrics, err := c.metricsClient.MetricsV1beta1().PodMetricses(c.config.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			continue // Skip if metrics not available for this pod
		}

		for _, container := range metrics.Containers {
			podMetrics = append(podMetrics, PodMetrics{
				Name:      pod.Name,
				CPU:       container.Usage.Cpu().String(),
				Memory:    container.Usage.Memory().String(),
				Timestamp: metrics.Timestamp.Time,
			})
		}
	}

	return &DeploymentMetrics{
		Pods: podMetrics,
	}, nil
}
