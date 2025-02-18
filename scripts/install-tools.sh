#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${YELLOW}🔧 Instalando ferramentas de desenvolvimento e segurança...${NC}"

# Configurando GOPATH se necessário
if [ -z "$GOPATH" ]; then
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
fi

# Criando diretório bin se não existir
mkdir -p $GOPATH/bin

# Função para instalar ferramenta Go
install_go_tool() {
    local tool=$1
    local package=$2
    echo -e "${YELLOW}Verificando $tool...${NC}"
    if ! command -v $tool &> /dev/null; then
        echo -e "${GREEN}Instalando $tool...${NC}"
        go install $package@latest
    else
        echo -e "${GREEN}$tool já está instalado${NC}"
    fi
}

# Função para instalar ferramenta via apt
install_apt_tool() {
    local tool=$1
    local package=$2
    echo -e "${YELLOW}Verificando $tool...${NC}"
    if ! command -v $tool &> /dev/null; then
        echo -e "${GREEN}Instalando $tool...${NC}"
        sudo apt-get update && sudo apt-get install -y $package
    else
        echo -e "${GREEN}$tool já está instalado${NC}"
    fi
}

# Ferramentas Go
install_go_tool "swag" "github.com/swaggo/swag/cmd/swag"
install_go_tool "goimports" "golang.org/x/tools/cmd/goimports"
install_go_tool "golangci-lint" "github.com/golangci/golangci-lint/cmd/golangci-lint"
install_go_tool "gosec" "github.com/securego/gosec/v2/cmd/gosec"
install_go_tool "dlv" "github.com/go-delve/delve/cmd/dlv"
install_go_tool "go-mod-outdated" "github.com/psampaz/go-mod-outdated"
install_go_tool "nancy" "github.com/sonatype-nexus-community/nancy"
install_go_tool "modgraphviz" "golang.org/x/exp/cmd/modgraphviz"

# Git Secrets
if ! command -v git-secrets &> /dev/null; then
    echo -e "${GREEN}Instalando git-secrets...${NC}"
    git clone https://github.com/awslabs/git-secrets.git
    cd git-secrets
    sudo make install
    cd ..
    rm -rf git-secrets
else
    echo -e "${GREEN}git-secrets já está instalado${NC}"
fi

# Graphviz (para visualização de dependências)
install_apt_tool "dot" "graphviz"

# Configurando git-secrets
git secrets --install
git secrets --register-aws

# Verificando instalações
echo -e "\n${YELLOW}Verificando instalações...${NC}"
for tool in swag goimports golangci-lint gosec dlv go-mod-outdated nancy modgraphviz; do
    if ! command -v $tool &> /dev/null; then
        echo -e "${RED}⚠️ $tool não foi instalado corretamente${NC}"
    else
        echo -e "${GREEN}✅ $tool instalado com sucesso${NC}"
    fi
done

# Verificando Graphviz
if ! command -v dot &> /dev/null; then
    echo -e "${RED}⚠️ graphviz não foi instalado corretamente${NC}"
else
    echo -e "${GREEN}✅ graphviz instalado com sucesso${NC}"
fi

echo -e "${GREEN}✅ Todas as ferramentas foram instaladas com sucesso!${NC}"
echo -e "${YELLOW}Ferramentas instaladas:${NC}"
echo "- swag (Swagger)"
echo "- goimports (Formatação)"
echo "- golangci-lint (Linting)"
echo "- gosec (Análise de Segurança)"
echo "- dlv (Debugger)"
echo "- git-secrets (Verificação de Secrets)"
echo "- go-mod-outdated (Verificação de Dependências)"
echo "- nancy (Verificação de Vulnerabilidades)"
echo "- modgraphviz (Visualização de Dependências)"
echo "- graphviz (Geração de Gráficos)"

echo -e "\n${YELLOW}Para começar a usar:${NC}"
echo "1. Execute 'make setup' para configurar o ambiente"
echo "2. Execute 'make help' para ver todos os comandos disponíveis"
echo "3. Execute 'make security-check' para verificar vulnerabilidades"
echo "4. Execute './scripts/check-deps.sh' para verificar dependências" 