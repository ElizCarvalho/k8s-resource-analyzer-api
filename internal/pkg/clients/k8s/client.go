package k8s

import (
	"context"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/errors"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/domain/types"
	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/logger"
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

	logger.Info("Criando cliente Kubernetes",
		logger.NewField("in_cluster", cfg.InCluster),
		logger.NewField("kubeconfig_path", cfg.KubeconfigPath),
	)

	if cfg.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			logger.Error("Erro ao criar configuração in-cluster", err)
			return nil, errors.NewInvalidConfigurationError("kubernetes", "erro ao criar configuração in-cluster")
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", cfg.KubeconfigPath)
		if err != nil {
			logger.Error("Erro ao criar configuração a partir do kubeconfig", err,
				logger.NewField("kubeconfig_path", cfg.KubeconfigPath),
			)
			return nil, errors.NewInvalidConfigurationError("kubernetes", "erro ao criar configuração a partir do kubeconfig")
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error("Erro ao criar clientset", err)
		return nil, errors.NewInvalidConfigurationError("kubernetes", "erro ao criar clientset")
	}

	metricsClient, err := metricsv1beta1.NewForConfig(config)
	if err != nil {
		logger.Error("Erro ao criar metrics client", err)
		return nil, errors.NewInvalidConfigurationError("kubernetes", "erro ao criar metrics client")
	}

	logger.Info("Cliente Kubernetes criado com sucesso")
	return &Client{
		clientset:     clientset,
		metricsClient: metricsClient,
	}, nil
}

// GetDeploymentMetrics retorna as métricas atuais de um deployment
func (c *Client) GetDeploymentMetrics(ctx context.Context, namespace, name string) (*types.K8sMetrics, error) {
	logger.Info("Obtendo métricas do deployment",
		logger.NewField("namespace", namespace),
		logger.NewField("deployment", name),
	)

	// Obtém os pods do deployment
	deployment, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		logger.Error("Erro ao obter deployment", err,
			logger.NewField("namespace", namespace),
			logger.NewField("deployment", name),
		)
		return nil, errors.NewResourceNotFoundError("deployment", "erro ao obter deployment")
	}
	logger.Info("Deployment encontrado",
		logger.NewField("name", deployment.Name),
	)

	// Obtém os pods usando o selector do deployment
	selector := deployment.Spec.Selector.MatchLabels
	logger.Info("Listando pods",
		logger.NewField("selector", selector),
	)

	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: selector}),
	})
	if err != nil {
		logger.Error("Erro ao listar pods", err,
			logger.NewField("namespace", namespace),
			logger.NewField("selector", selector),
		)
		return nil, errors.NewResourceNotFoundError("pods", "erro ao listar pods")
	}
	logger.Info("Pods encontrados",
		logger.NewField("count", len(pods.Items)),
	)

	// Inicializa as métricas
	result := &types.K8sMetrics{
		CPU: struct {
			Usage       float64 `json:"usage"`
			Average     float64 `json:"average"`
			Peak        float64 `json:"peak"`
			Utilization float64 `json:"utilization"`
		}{},
		Memory: struct {
			Usage       float64 `json:"usage"`
			Average     float64 `json:"average"`
			Peak        float64 `json:"peak"`
			Utilization float64 `json:"utilization"`
		}{},
		Pods: struct {
			Running     int     `json:"running"`
			Utilization float64 `json:"utilization"`
		}{},
	}
	var totalCPUUsage, peakCPUUsage, totalMemoryUsage, peakMemoryUsage float64
	runningPods := 0

	// Coleta métricas de cada pod
	for _, pod := range pods.Items {
		logger.Info("Verificando pod",
			logger.NewField("name", pod.Name),
			logger.NewField("status", pod.Status.Phase),
		)
		if pod.Status.Phase != "Running" {
			continue
		}
		runningPods++

		// Obtém métricas do pod
		podMetrics, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			logger.Error("Erro ao obter métricas do pod", err,
				logger.NewField("pod", pod.Name),
			)
			continue
		}
		logger.Info("Métricas obtidas para o pod",
			logger.NewField("pod", pod.Name),
		)

		// Soma métricas de todos os containers do pod
		for _, container := range podMetrics.Containers {
			cpuUsage := float64(container.Usage.Cpu().MilliValue())                  // Já está em milicores
			memoryUsage := float64(container.Usage.Memory().Value()) / (1024 * 1024) // Converte para Mi

			logger.Info("Métricas do container",
				logger.NewField("container", container.Name),
				logger.NewField("cpu_usage", cpuUsage),
				logger.NewField("memory_usage", memoryUsage),
			)

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

	// Calcula médias e utilização
	if runningPods > 0 {
		result.CPU.Usage = totalCPUUsage
		result.CPU.Average = totalCPUUsage / float64(runningPods)
		result.CPU.Peak = peakCPUUsage
		result.CPU.Utilization = totalCPUUsage / float64(deployment.Status.Replicas) * 100

		result.Memory.Usage = totalMemoryUsage
		result.Memory.Average = totalMemoryUsage / float64(runningPods)
		result.Memory.Peak = peakMemoryUsage
		result.Memory.Utilization = totalMemoryUsage / float64(deployment.Status.Replicas) * 100

		result.Pods.Running = runningPods
		result.Pods.Utilization = float64(runningPods) / float64(deployment.Status.Replicas) * 100
	}

	logger.Info("Métricas coletadas com sucesso",
		logger.NewField("cpu_usage", result.CPU.Usage),
		logger.NewField("memory_usage", result.Memory.Usage),
		logger.NewField("running_pods", result.Pods.Running),
	)

	return result, nil
}

// GetDeploymentConfig retorna a configuração de um deployment
func (c *Client) GetDeploymentConfig(ctx context.Context, namespace, name string) (*types.K8sDeploymentConfig, error) {
	logger.Info("Obtendo configuração do deployment",
		logger.NewField("namespace", namespace),
		logger.NewField("deployment", name),
	)

	// Obtém o deployment
	deployment, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		logger.Error("Erro ao obter deployment", err,
			logger.NewField("namespace", namespace),
			logger.NewField("deployment", name),
		)
		return nil, errors.NewResourceNotFoundError("deployment", "erro ao obter deployment")
	}

	// Obtém o HPA, se existir
	hpa, err := c.clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		logger.Info("HPA não encontrado",
			logger.NewField("namespace", namespace),
			logger.NewField("deployment", name),
		)
		hpa = nil
	} else {
		logger.Info("HPA encontrado",
			logger.NewField("namespace", namespace),
			logger.NewField("deployment", name),
		)
	}

	result := &types.K8sDeploymentConfig{}

	// Obtém requests e limits do primeiro container
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		container := deployment.Spec.Template.Spec.Containers[0]

		// CPU
		if cpu := container.Resources.Requests.Cpu(); cpu != nil {
			result.CPU.Request = float64(cpu.MilliValue()) // Já está em milicores
		}
		if cpu := container.Resources.Limits.Cpu(); cpu != nil {
			result.CPU.Limit = float64(cpu.MilliValue()) // Já está em milicores
		}

		// Memória
		if memory := container.Resources.Requests.Memory(); memory != nil {
			result.Memory.Request = float64(memory.Value()) / (1024 * 1024) // Converte bytes para Mi
		}
		if memory := container.Resources.Limits.Memory(); memory != nil {
			result.Memory.Limit = float64(memory.Value()) / (1024 * 1024) // Converte bytes para Mi
		}

		logger.Info("Recursos do container",
			logger.NewField("cpu_request", result.CPU.Request),
			logger.NewField("cpu_limit", result.CPU.Limit),
			logger.NewField("memory_request", result.Memory.Request),
			logger.NewField("memory_limit", result.Memory.Limit),
		)
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

	logger.Info("Configuração de pods",
		logger.NewField("replicas", result.Pods.Replicas),
		logger.NewField("min_replicas", result.Pods.MinReplicas),
		logger.NewField("max_replicas", result.Pods.MaxReplicas),
	)

	return result, nil
}

// CheckConnection verifica a conexão com o cluster Kubernetes
func (c *Client) CheckConnection(ctx context.Context) error {
	logger.Info("Verificando conexão com o cluster Kubernetes")
	_, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error("Erro ao conectar ao cluster", err)
		return errors.NewInvalidConfigurationError("kubernetes", "erro ao conectar ao cluster")
	}
	logger.Info("Conexão com o cluster estabelecida com sucesso")
	return nil
}
