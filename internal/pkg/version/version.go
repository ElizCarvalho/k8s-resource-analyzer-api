package version

var (
	// Version é a versão atual da aplicação
	// Será preenchida durante o build usando ldflags
	Version = "dev"
)

// GetVersion retorna a versão atual da aplicação
func GetVersion() string {
	return Version
}
