#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Variáveis
APP_NAME="k8s-resource-analyzer-api"
REQUIRED_GO_VERSION="1.22"
MIN_DOCKER_VERSION="20.10"

# Função para verificar versão do Go
check_go() {
    echo -e "${BLUE}🔍 Verificando instalação do Go...${NC}"
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go não está instalado${NC}"
        echo -e "${YELLOW}Por favor, instale o Go $REQUIRED_GO_VERSION ou superior:${NC}"
        echo "https://golang.org/doc/install"
        exit 1
    fi
    
    local version=$(go version | awk '{print $3}' | sed 's/go//')
    echo -e "Go version: $version"
    
    if [[ "$version" < "$REQUIRED_GO_VERSION" ]]; then
        echo -e "${RED}❌ Versão do Go muito antiga. Necessário $REQUIRED_GO_VERSION ou superior${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Versão do Go OK${NC}"
}

# Função para verificar Docker
check_docker() {
    echo -e "\n${BLUE}🔍 Verificando instalação do Docker...${NC}"
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}❌ Docker não está instalado${NC}"
        echo -e "${YELLOW}Por favor, instale o Docker:${NC}"
        echo "https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    local version=$(docker --version | awk '{print $3}' | cut -d'.' -f1,2)
    echo -e "Docker version: $version"
    
    if [[ "$version" < "$MIN_DOCKER_VERSION" ]]; then
        echo -e "${RED}❌ Versão do Docker muito antiga. Necessário $MIN_DOCKER_VERSION ou superior${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Versão do Docker OK${NC}"
}

# Função para verificar kubectl
check_kubectl() {
    echo -e "\n${BLUE}🔍 Verificando instalação do kubectl...${NC}"
    if ! command -v kubectl &> /dev/null; then
        echo -e "${RED}❌ kubectl não está instalado${NC}"
        echo -e "${YELLOW}Por favor, instale o kubectl:${NC}"
        echo "https://kubernetes.io/docs/tasks/tools/install-kubectl/"
        exit 1
    fi
    
    echo -e "${GREEN}✅ kubectl instalado${NC}"
}

# Função para configurar diretórios
setup_directories() {
    echo -e "\n${BLUE}📁 Configurando diretórios...${NC}"
    mkdir -p {logs,tmp,backups,docs}
    echo -e "${GREEN}✅ Diretórios criados${NC}"
}

# Função para configurar git hooks
setup_git_hooks() {
    echo -e "\n${BLUE}🔧 Configurando git hooks...${NC}"
    
    # Pre-commit hook
    cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
echo "🔍 Executando verificações pre-commit..."

# Verifica formatação
echo "Verificando formatação..."
make format

# Executa linter
echo "Executando linter..."
make lint

# Executa testes
echo "Executando testes..."
make test
EOF
    
    chmod +x .git/hooks/pre-commit
    echo -e "${GREEN}✅ Git hooks configurados${NC}"
}

# Função para configurar ambiente Go
setup_go_env() {
    echo -e "\n${BLUE}🔧 Configurando ambiente Go...${NC}"
    
    # Configura GOPATH se necessário
    if [ -z "$GOPATH" ]; then
        echo "export GOPATH=$HOME/go" >> ~/.bashrc
        echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.bashrc
        source ~/.bashrc
    fi
    
    # Inicializa módulos Go se necessário
    if [ ! -f "go.mod" ]; then
        go mod init github.com/ElizCarvalho/$APP_NAME
        go mod tidy
    fi
    
    echo -e "${GREEN}✅ Ambiente Go configurado${NC}"
}

# Função para configurar arquivo .env
setup_env() {
    echo -e "\n${BLUE}📝 Configurando variáveis de ambiente...${NC}"
    
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
            echo -e "${GREEN}✅ Arquivo .env criado a partir do .env.example${NC}"
        else
            cat > .env << EOF
# Configurações da API
PORT=9000
DEBUG_PORT=2345
ENV=development

# Configurações do Kubernetes
KUBECONFIG=~/.kube/config

# Configurações de Log
LOG_LEVEL=debug
LOG_FORMAT=json
EOF
            echo -e "${GREEN}✅ Arquivo .env criado${NC}"
        fi
    else
        echo -e "${YELLOW}ℹ️ Arquivo .env já existe${NC}"
    fi
}

# Função principal
main() {
    echo -e "${BLUE}🚀 Configurando ambiente de desenvolvimento para $APP_NAME${NC}\n"
    
    # Verificações iniciais
    check_go
    check_docker
    check_kubectl
    
    # Configurações
    setup_directories
    setup_git_hooks
    setup_go_env
    setup_env
    
    # Instala dependências de desenvolvimento
    echo -e "\n${BLUE}📦 Instalando dependências de desenvolvimento...${NC}"
    ./scripts/install-tools.sh
    
    echo -e "\n${GREEN}✅ Ambiente de desenvolvimento configurado com sucesso!${NC}"
    echo -e "\n${YELLOW}Próximos passos:${NC}"
    echo "1. Revise o arquivo .env e ajuste as configurações conforme necessário"
    echo "2. Execute 'make help' para ver todos os comandos disponíveis"
    echo "3. Execute 'make run' para iniciar a aplicação"
}

# Executa função principal
main 