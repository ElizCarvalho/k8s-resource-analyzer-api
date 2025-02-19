#!/bin/bash
# ==============================================================================
# Script de Build - K8s Resource Analyzer API
# ==============================================================================
# Autor: Elizabeth Carvalho
# Data: Fevereiro 2024
#
# Descrição:
#   Script para compilar a aplicação com informações de versão.
#   Injeta a versão do git (tag ou commit hash) no binário.
#
# Uso:
#   ./build.sh [output_dir]
#
# Argumentos:
#   output_dir - Diretório de saída para o binário (opcional, padrão: bin/)
#
# Exemplos:
#   ./build.sh
#   ./build.sh /tmp/output
# ==============================================================================

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Diretório de saída padrão
OUTPUT_DIR=${1:-"bin"}

# Função para log com cores
log() {
    local level=$1
    local message=$2
    case $level in
        "info")
            echo -e "${GREEN}[INFO]${NC} $message"
            ;;
        "warn")
            echo -e "${YELLOW}[WARN]${NC} $message"
            ;;
        "error")
            echo -e "${RED}[ERROR]${NC} $message"
            ;;
    esac
}

# Cria diretório de saída se não existir
if [ ! -d "$OUTPUT_DIR" ]; then
    log "info" "Criando diretório de saída: $OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR"
fi

# Obtém a versão do git tag ou usa o hash do commit atual
VERSION=$(git describe --tags 2>/dev/null || git rev-parse --short HEAD)
log "info" "Versão detectada: $VERSION"

# Compila a aplicação injetando a versão
log "info" "Iniciando build da aplicação..."
go build -ldflags "-X github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/version.Version=${VERSION}" \
    -o "$OUTPUT_DIR/api" cmd/api/main.go

if [ $? -eq 0 ]; then
    log "info" "Build concluído com sucesso! Binário disponível em: $OUTPUT_DIR/api"
else
    log "error" "Erro durante o build da aplicação"
    exit 1
fi 