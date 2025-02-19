package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned"
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

// Client implementa a interface metrics.K8sClient
type Client struct {
	clientset     *kubernetes.Clientset
	metricsClient *metricsv1beta1.Clientset
}

// ClientConfig contém as configurações para o cliente Kubernetes
type ClientConfig struct {
	KubeconfigPath string
	InCluster      bool
}

// NewClient cria uma nova instância do cliente Kubernetes
func NewClient(cfg *ClientConfig) (*Client, error) {
	var config *rest.Config
	var err error

	if cfg.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("erro ao criar configuração in-cluster: %w", err)
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", cfg.KubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar configuração a partir do kubeconfig: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar clientset: %w", err)
	}

	metricsClient, err := metricsv1beta1.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar metrics client: %w", err)
	}

	return &Client{
		clientset:     clientset,
		metricsClient: metricsClient,
	}, nil
}

// GetDeploymentMetrics retorna as métricas atuais de um deployment
func (c *Client) GetDeploymentMetrics(ctx context.Context, namespace, name string) (*types.K8sMetrics, error) {
	fmt.Printf("Obtendo métricas para deployment %s no namespace %s\n", name, namespace)

	// Obtém os pods do deployment
	deployment, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Erro ao obter deployment: %v\n", err)
		return nil, fmt.Errorf("erro ao obter deployment: %w", err)
	}
	fmt.Printf("Deployment encontrado: %s\n", deployment.Name)

	// Obtém os pods usando o selector do deployment
	selector := deployment.Spec.Selector.MatchLabels
	fmt.Printf("Usando selector: %v\n", selector)

	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: selector}),
	})
	if err != nil {
		fmt.Printf("Erro ao listar pods: %v\n", err)
		return nil, fmt.Errorf("erro ao listar pods: %w", err)
	}
	fmt.Printf("Encontrados %d pods\n", len(pods.Items))

	// Inicializa as métricas
	result := &types.K8sMetrics{}
	var totalCPUUsage, peakCPUUsage, totalMemoryUsage, peakMemoryUsage float64
	runningPods := 0

	// Coleta métricas de cada pod
	for _, pod := range pods.Items {
		fmt.Printf("Verificando pod %s (status: %s)\n", pod.Name, pod.Status.Phase)
		if pod.Status.Phase != "Running" {
			continue
		}
		runningPods++

		// Obtém métricas do pod
		podMetrics, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Erro ao obter métricas do pod %s: %v\n", pod.Name, err)
			continue
		}
		fmt.Printf("Métricas obtidas para o pod %s\n", pod.Name)

		// Soma métricas de todos os containers do pod
		for _, container := range podMetrics.Containers {
			cpuUsage := float64(container.Usage.Cpu().MilliValue()) / 1000
			memoryUsage := float64(container.Usage.Memory().Value()) / (1024 * 1024 * 1024) // Converte para GB

			fmt.Printf("Container %s: CPU=%.3f cores, Memory=%.3f GB\n", container.Name, cpuUsage, memoryUsage)

			totalCPUUsage += cpuUsage
			totalMemoryUsage += memoryUsage

			if cpuUsage > peakCPUUsage {
				peakCPUUsage = cpuUsage
			}
			if memoryUsage > peakMemoryUsage {
				peakMemoryUsage = memoryUsage
			}
		}
	}

	// Calcula médias
	if runningPods > 0 {
		result.CPU.Average = totalCPUUsage / float64(runningPods)
		result.CPU.Peak = peakCPUUsage
		result.CPU.Usage = totalCPUUsage

		result.Memory.Average = totalMemoryUsage / float64(runningPods)
		result.Memory.Peak = peakMemoryUsage
		result.Memory.Usage = totalMemoryUsage

		fmt.Printf("Métricas calculadas:\n")
		fmt.Printf("CPU: Average=%.3f, Peak=%.3f, Total=%.3f\n", result.CPU.Average, result.CPU.Peak, result.CPU.Usage)
		fmt.Printf("Memory: Average=%.3f, Peak=%.3f, Total=%.3f\n", result.Memory.Average, result.Memory.Peak, result.Memory.Usage)
	}

	result.Pods.Running = runningPods
	fmt.Printf("Total de pods em execução: %d\n", runningPods)

	return result, nil
}

// GetDeploymentConfig retorna a configuração de um deployment
func (c *Client) GetDeploymentConfig(ctx context.Context, namespace, name string) (*types.K8sDeploymentConfig, error) {
	// Obtém o deployment
	deployment, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("erro ao obter deployment: %w", err)
	}

	// Obtém o HPA, se existir
	hpa, err := c.clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		// Ignora erro se HPA não existir
		hpa = nil
	}

	result := &types.K8sDeploymentConfig{}

	// Obtém requests e limits do primeiro container
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		container := deployment.Spec.Template.Spec.Containers[0]

		// CPU
		if cpu := container.Resources.Requests.Cpu(); cpu != nil {
			result.CPU.Request = float64(cpu.MilliValue()) / 1000
		}
		if cpu := container.Resources.Limits.Cpu(); cpu != nil {
			result.CPU.Limit = float64(cpu.MilliValue()) / 1000
		}

		// Memória
		if memory := container.Resources.Requests.Memory(); memory != nil {
			result.Memory.Request = float64(memory.Value()) / (1024 * 1024 * 1024) // Converte para GB
		}
		if memory := container.Resources.Limits.Memory(); memory != nil {
			result.Memory.Limit = float64(memory.Value()) / (1024 * 1024 * 1024) // Converte para GB
		}
	}

	// Configuração de pods
	result.Pods.Replicas = int(*deployment.Spec.Replicas)

	if hpa != nil {
		result.Pods.MinReplicas = int(*hpa.Spec.MinReplicas)
		result.Pods.MaxReplicas = int(hpa.Spec.MaxReplicas)
	} else {
		// Se não houver HPA, usa o número de réplicas do deployment
		result.Pods.MinReplicas = result.Pods.Replicas
		result.Pods.MaxReplicas = result.Pods.Replicas
	}

	return result, nil
}

// CheckConnection verifica a conexão com o cluster Kubernetes
func (c *Client) CheckConnection(ctx context.Context) error {
	_, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("erro ao conectar ao cluster: %w", err)
	}
	return nil
}
