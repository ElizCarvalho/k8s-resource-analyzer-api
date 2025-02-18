#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}🔧 Instalando ferramentas de desenvolvimento e segurança...${NC}"

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

# Configurando git-secrets
git secrets --install
git secrets --register-aws

echo -e "${GREEN}✅ Todas as ferramentas foram instaladas com sucesso!${NC}"
echo -e "${YELLOW}Ferramentas instaladas:${NC}"
echo "- swag (Swagger)"
echo "- goimports (Formatação)"
echo "- golangci-lint (Linting)"
echo "- gosec (Análise de Segurança)"
echo "- dlv (Debugger)"
echo "- git-secrets (Verificação de Secrets)"

echo -e "\n${YELLOW}Para começar a usar:${NC}"
echo "1. Execute 'make setup' para configurar o ambiente"
echo "2. Execute 'make help' para ver todos os comandos disponíveis"
echo "3. Execute 'make security-check' para verificar vulnerabilidades" 