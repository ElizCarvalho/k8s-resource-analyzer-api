package k8s

// Config contém as configurações para conexão com o cluster Kubernetes
type Config struct {
	// KubeconfigPath é o caminho para o arquivo kubeconfig
	// Se vazio e InCluster for false, usa o padrão (~/.kube/config)
	KubeconfigPath string

	// Namespace é o namespace Kubernetes a ser usado
	// Se vazio, usa "default"
	Namespace string

	// InCluster indica se deve usar configuração in-cluster (ServiceAccount)
	InCluster bool
}
