package k8s_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/k8s"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

const (
	testDeploymentName = "test-deployment"
	testNamespace      = "default"
	testImage          = "nginx:latest"
)

// createTestDeployment cria um deployment de teste com recursos definidos
func createTestDeployment(ctx context.Context, clientset *kubernetes.Clientset) error {
	replicas := int32(2)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testDeploymentName,
			Namespace: testNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": testDeploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": testDeploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: testImage,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("128Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("200m"),
									corev1.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := clientset.AppsV1().Deployments(testNamespace).Create(ctx, deployment, metav1.CreateOptions{})
	return err
}

// waitForDeploymentReady espera até que o deployment esteja pronto
func waitForDeploymentReady(ctx context.Context, clientset *kubernetes.Clientset, timeout time.Duration) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-timeoutCh:
			return fmt.Errorf("timeout waiting for deployment to be ready")
		case <-ticker.C:
			deployment, err := clientset.AppsV1().Deployments(testNamespace).Get(ctx, testDeploymentName, metav1.GetOptions{})
			if err != nil {
				return err
			}

			if deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
				return nil
			}
		}
	}
}

func TestGetDeploymentMetrics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Configurar K3s container
	k3sContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "rancher/k3s:latest",
			Cmd:          []string{"server", "--disable-network-policy", "--disable", "traefik"},
			ExposedPorts: []string{"6443/tcp"},
			Privileged:   true,
			Binds: []string{
				"/tmp:/output",
			},
			Env: map[string]string{
				"K3S_KUBECONFIG_OUTPUT": "/output/kubeconfig.yaml",
				"K3S_KUBECONFIG_MODE":   "666",
			},
			WaitingFor: wait.ForLog("Running kube-apiserver"),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("failed to start k3s: %v", err)
	}
	defer func() {
		if err := k3sContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate k3s container: %v", err)
		}
	}()

	// Obter o endereço do servidor
	endpoint, err := k3sContainer.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("failed to get endpoint: %v", err)
	}

	// Esperar um pouco para o kubeconfig ser gerado
	time.Sleep(5 * time.Second)

	// Copiar kubeconfig do container
	kubeconfigBytes, err := k3sContainer.CopyFileFromContainer(ctx, "/output/kubeconfig.yaml")
	if err != nil {
		t.Fatalf("failed to copy kubeconfig: %v", err)
	}

	// Criar arquivo temporário para o kubeconfig
	kubeconfigFile, err := os.CreateTemp("", "kubeconfig-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(kubeconfigFile.Name())

	// Copiar conteúdo do kubeconfig
	_, err = io.Copy(kubeconfigFile, kubeconfigBytes)
	if err != nil {
		t.Fatalf("failed to write kubeconfig: %v", err)
	}
	kubeconfigFile.Close()

	// Atualizar o kubeconfig com o endereço correto do servidor
	kubeConfig, err := clientcmd.LoadFromFile(kubeconfigFile.Name())
	if err != nil {
		t.Fatalf("failed to load kubeconfig: %v", err)
	}

	// Atualizar o endereço do servidor no kubeconfig
	for _, cluster := range kubeConfig.Clusters {
		cluster.Server = fmt.Sprintf("https://%s", endpoint)
	}

	// Salvar o kubeconfig atualizado
	if err := clientcmd.WriteToFile(*kubeConfig, kubeconfigFile.Name()); err != nil {
		t.Fatalf("failed to write updated kubeconfig: %v", err)
	}

	// Criar clientset para setup do teste
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigFile.Name())
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("failed to create clientset: %v", err)
	}

	// Criar deployment de teste
	if err := createTestDeployment(ctx, clientset); err != nil {
		t.Fatalf("failed to create test deployment: %v", err)
	}

	// Esperar deployment estar pronto
	if err := waitForDeploymentReady(ctx, clientset, 2*time.Minute); err != nil {
		t.Fatalf("deployment not ready: %v", err)
	}

	// Verificar se o metrics-server está funcionando
	t.Log("Verificando metrics-server...")
	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		t.Fatalf("failed to create metrics client: %v", err)
	}

	// Esperar até que as métricas estejam disponíveis
	if err := waitForMetrics(ctx, metricsClient, testNamespace, 2*time.Minute); err != nil {
		t.Fatalf("metrics not available: %v", err)
	}

	// Criar cliente de teste
	cfg := &k8s.Config{
		KubeconfigPath: kubeconfigFile.Name(),
		Namespace:      testNamespace,
		InCluster:      false,
	}

	client, err := k8s.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Testar obtenção de métricas
	metrics, err := client.GetDeploymentMetrics(ctx, testDeploymentName)
	if err != nil {
		t.Fatalf("failed to get metrics: %v", err)
	}

	if len(metrics.Pods) == 0 {
		t.Error("expected pod metrics, got none")
	}

	// Validar métricas
	for _, pod := range metrics.Pods {
		if pod.CPU == "" {
			t.Error("expected CPU metrics")
		}
		if pod.Memory == "" {
			t.Error("expected Memory metrics")
		}
		t.Logf("Pod %s: CPU=%s, Memory=%s", pod.Name, pod.CPU, pod.Memory)
	}

	// Instalar metrics-server
	t.Log("Instalando metrics-server...")
	metricsServerManifest := `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: metrics-server
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:metrics-server
rules:
- apiGroups: [""]
  resources: ["pods", "nodes", "nodes/stats", "namespaces", "configmaps"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:metrics-server
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:metrics-server
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system
---
apiVersion: v1
kind: Service
metadata:
  name: metrics-server
  namespace: kube-system
spec:
  ports:
  - port: 443
    protocol: TCP
    targetPort: 4443
  selector:
    k8s-app: metrics-server
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: metrics-server
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: metrics-server
  template:
    metadata:
      labels:
        k8s-app: metrics-server
    spec:
      serviceAccountName: metrics-server
      volumes:
      - name: tmp-dir
        emptyDir: {}
      containers:
      - name: metrics-server
        image: k8s.gcr.io/metrics-server/metrics-server:v0.6.3
        args:
          - --cert-dir=/tmp
          - --secure-port=4443
          - --kubelet-insecure-tls
          - --kubelet-preferred-address-types=InternalIP
        ports:
        - name: tcp-4443
          containerPort: 4443
          protocol: TCP
        volumeMounts:
        - name: tmp-dir
          mountPath: /tmp
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          runAsUser: 1000
---
apiVersion: apiregistration.k8s.io/v1
kind: APIService
metadata:
  name: v1beta1.metrics.k8s.io
spec:
  service:
    name: metrics-server
    namespace: kube-system
  group: metrics.k8s.io
  version: v1beta1
  insecureSkipTLSVerify: true
  groupPriorityMinimum: 100
  versionPriority: 100`

	// Criar arquivo temporário para o manifesto do metrics-server
	metricsServerFile, err := os.CreateTemp("", "metrics-server-*.yaml")
	if err != nil {
		t.Fatalf("failed to create metrics-server manifest file: %v", err)
	}
	defer os.Remove(metricsServerFile.Name())

	if _, err := metricsServerFile.WriteString(metricsServerManifest); err != nil {
		t.Fatalf("failed to write metrics-server manifest: %v", err)
	}
	metricsServerFile.Close()

	// Aplicar manifesto do metrics-server usando kubectl
	kubectlApply := fmt.Sprintf("kubectl --kubeconfig=%s apply -f %s", kubeconfigFile.Name(), metricsServerFile.Name())
	if err := exec.Command("sh", "-c", kubectlApply).Run(); err != nil {
		t.Fatalf("failed to apply metrics-server manifest: %v", err)
	}

	// Esperar mais tempo para o metrics-server iniciar
	t.Log("Aguardando metrics-server iniciar...")
	time.Sleep(45 * time.Second)
}

// waitForMetrics espera até que as métricas estejam disponíveis
func waitForMetrics(ctx context.Context, client *versioned.Clientset, namespace string, timeout time.Duration) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-timeoutCh:
			return fmt.Errorf("timeout waiting for metrics")
		case <-ticker.C:
			_, err := client.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
			if err == nil {
				return nil
			}
		}
	}
}

func TestGetTravelerNotifierMetrics(t *testing.T) {
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
	cfg := &k8s.Config{
		KubeconfigPath: kubeconfigPath,
		Namespace:      "default", // ajuste se necessário
		InCluster:      false,
	}

	client, err := k8s.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Testar obtenção de métricas
	metrics, err := client.GetDeploymentMetrics(ctx, "travelernotifierbyevent")
	if err != nil {
		t.Fatalf("failed to get metrics: %v", err)
	}

	if len(metrics.Pods) == 0 {
		t.Error("expected pod metrics, got none")
	}

	// Validar métricas
	for _, pod := range metrics.Pods {
		if pod.CPU == "" {
			t.Error("expected CPU metrics")
		}
		if pod.Memory == "" {
			t.Error("expected Memory metrics")
		}
		t.Logf("Pod %s: CPU=%s, Memory=%s", pod.Name, pod.CPU, pod.Memory)
	}
}
