#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Variáveis
APP_NAME="k8s-resource-analyzer-api"
PORT=${PORT:-9000}
INTERVAL=${INTERVAL:-5}

# Função para verificar saúde da aplicação
check_health() {
    local response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:${PORT}/api/v1/ping)
    if [ "$response" == "200" ]; then
        echo -e "${GREEN}✅ API está saudável (HTTP 200)${NC}"
    else
        echo -e "${RED}❌ API não está respondendo corretamente (HTTP $response)${NC}"
    fi
}

# Função para coletar métricas do sistema
collect_system_metrics() {
    echo -e "\n${BLUE}📊 Métricas do Sistema:${NC}"
    echo -e "${YELLOW}CPU:${NC}"
    top -bn1 | grep "Cpu(s)" | awk '{print "  Uso: " $2 "%"}'
    
    echo -e "\n${YELLOW}Memória:${NC}"
    free -m | awk 'NR==2{printf "  Total: %s MB, Usado: %s MB, Livre: %s MB\n", $2,$3,$4}'
    
    echo -e "\n${YELLOW}Disco:${NC}"
    df -h | awk '$NF=="/"{printf "  Total: %s, Usado: %s, Livre: %s\n", $2,$3,$4}'
}

# Função para coletar métricas da aplicação
collect_app_metrics() {
    echo -e "\n${BLUE}📈 Métricas da Aplicação:${NC}"
    
    # Processos
    local process_count=$(ps aux | grep ${APP_NAME} | grep -v grep | wc -l)
    echo -e "${YELLOW}Processos:${NC}"
    echo "  Ativos: $process_count"
    
    # Portas
    echo -e "\n${YELLOW}Portas:${NC}"
    netstat -tlpn 2>/dev/null | grep ${PORT} || echo "  Nenhuma porta encontrada"
    
    # Logs recentes
    echo -e "\n${YELLOW}Logs Recentes:${NC}"
    tail -n 5 ./logs/${APP_NAME}.log 2>/dev/null || echo "  Nenhum log encontrado"
}

# Função para mostrar estatísticas do projeto
show_project_stats() {
    echo -e "\n${BLUE}📁 Estatísticas do Projeto:${NC}"
    
    echo -e "${YELLOW}Código:${NC}"
    echo "  Arquivos Go: $(find . -name '*.go' -not -path './vendor/*' | wc -l)"
    echo "  Linhas de código: $(find . -name '*.go' -not -path './vendor/*' | xargs wc -l 2>/dev/null | tail -n 1 | awk '{print $1}')"
    
    echo -e "\n${YELLOW}Git:${NC}"
    echo "  Commits: $(git rev-list --count HEAD 2>/dev/null || echo 'N/A')"
    echo "  Branch atual: $(git branch --show-current 2>/dev/null || echo 'N/A')"
    echo "  Último commit: $(git log -1 --format=%cd --date=relative 2>/dev/null || echo 'N/A')"
}

# Função principal de monitoramento
monitor() {
    while true; do
        clear
        echo -e "${BLUE}🔍 Monitoramento do ${APP_NAME}${NC}"
        echo -e "Hora: $(date '+%Y-%m-%d %H:%M:%S')\n"
        
        check_health
        collect_system_metrics
        collect_app_metrics
        show_project_stats
        
        echo -e "\n${YELLOW}Atualizando em ${INTERVAL} segundos...${NC}"
        sleep ${INTERVAL}
    done
}

# Menu de opções
show_menu() {
    echo -e "${BLUE}🔧 Ferramentas de Monitoramento${NC}\n"
    echo "1. Iniciar monitoramento contínuo"
    echo "2. Verificar saúde da aplicação"
    echo "3. Mostrar métricas do sistema"
    echo "4. Mostrar métricas da aplicação"
    echo "5. Mostrar estatísticas do projeto"
    echo "6. Sair"
    echo -e "\nEscolha uma opção: "
}

# Loop principal
while true; do
    show_menu
    read -r opt
    
    case $opt in
        1) monitor ;;
        2) check_health ;;
        3) collect_system_metrics ;;
        4) collect_app_metrics ;;
        5) show_project_stats ;;
        6) echo -e "${GREEN}Encerrando...${NC}"; exit 0 ;;
        *) echo -e "${RED}Opção inválida${NC}" ;;
    esac
    
    if [ "$opt" != "1" ]; then
        echo -e "\nPressione ENTER para continuar..."
        read -r
    fi
done 