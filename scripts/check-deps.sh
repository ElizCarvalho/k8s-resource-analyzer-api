#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configurando GOPATH se necessário
if [ -z "$GOPATH" ]; then
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
fi

# Verificando ferramentas necessárias
check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}❌ Ferramenta $1 não encontrada. Execute 'make install-tools' primeiro.${NC}"
        exit 1
    fi
}

check_tool "go-mod-outdated"
check_tool "nancy"
check_tool "modgraphviz"
check_tool "dot"

echo -e "${BLUE}🔍 Verificando dependências...${NC}"

# Verifica atualizações disponíveis
echo -e "\n${YELLOW}Atualizações disponíveis:${NC}"
go list -u -m -json all | $GOPATH/bin/go-mod-outdated -update -direct

# Verifica vulnerabilidades
echo -e "\n${YELLOW}Verificando vulnerabilidades:${NC}"
go list -json -m all | $GOPATH/bin/nancy sleuth

# Mostra um resumo das dependências
echo -e "\n${YELLOW}Resumo das dependências:${NC}"
go mod graph | $GOPATH/bin/modgraphviz | dot -Tsvg -o deps-graph.svg

echo -e "\n${GREEN}✅ Verificação concluída!${NC}"
echo -e "Para atualizar todas as dependências, execute: ${BLUE}go get -u ./...${NC}"
echo -e "Para atualizar uma dependência específica, execute: ${BLUE}go get -u nome-do-pacote${NC}" 