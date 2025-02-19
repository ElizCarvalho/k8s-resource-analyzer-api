# ==============================================================================
# K8s Resource Analyzer API - Makefile
# ==============================================================================
# Autor: Elizabeth Carvalho
# Repositório: https://github.com/ElizCarvalho/k8s-resource-analyzer-api
# 
# Este Makefile contém todos os comandos necessários para desenvolvimento,
# teste, build e deploy da aplicação.
#
# Uso: make <comando>
# Execute 'make help' para ver todos os comandos disponíveis
# ==============================================================================

# ==============================================================================
# Variáveis do Projeto
# ==============================================================================
APP_NAME=k8s-resource-analyzer-api
DOCKER_IMAGE=eliscarvalho/$(APP_NAME)
VERSION=$(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD || echo "dev")
PORT?=9000
DEBUG_PORT?=2345
APP_DIR=/usr/app

# Versão e Build Info
LDFLAGS=-X github.com/ElizCarvalho/k8s-resource-analyzer-api/internal/pkg/version.Version=$(VERSION)

# ==============================================================================
# Variáveis de Ambiente
# ==============================================================================
INTERACTIVE:=$(shell [ -t 0 ] && echo i || echo d)
PROJECT_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
PWD=$(shell pwd)
SHELL:=/bin/bash

# ==============================================================================
# Variáveis Go
# ==============================================================================
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
CMD_PATH=./cmd/api
MAIN_FILE=$(CMD_PATH)/main.go

# ==============================================================================
# Cores e Formatação
# ==============================================================================
BLUE=\033[0;34m
GREEN=\033[0;32m
RED=\033[0;31m
YELLOW=\033[0;33m
BOLD=\033[1m
NC=\033[0m # No Color

# ==============================================================================
# ASCII Art Header
# ==============================================================================
define HEADER

$(BLUE)██╗  ██╗ █████╗ ███████╗    $(YELLOW)█████╗ ███╗   ██╗ █████╗ ██╗  ██╗   ██╗███████╗███████╗██████╗ 
$(BLUE)██║ ██╔╝██╔══██╗██╔════╝   $(YELLOW)██╔══██╗████╗  ██║██╔══██╗██║  ╚██╗ ██╔╝╚══███╔╝██╔════╝██╔══██╗
$(BLUE)█████╔╝ ╚█████╔╝███████╗   $(YELLOW)███████║██╔██╗ ██║███████║██║   ╚████╔╝   ███╔╝ █████╗  ██████╔╝
$(BLUE)██╔═██╗ ██╔══██╗╚════██║   $(YELLOW)██╔══██║██║╚██╗██║██╔══██║██║    ╚██╔╝   ███╔╝  ██╔══╝  ██╔══██╗
$(BLUE)██║  ██╗╚█████╔╝███████║   $(YELLOW)██║  ██║██║ ╚████║██║  ██║███████╗██║   ███████╗███████╗██║  ██║
$(BLUE)╚═╝  ╚═╝ ╚════╝ ╚══════╝   $(YELLOW)╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝╚═╝   ╚══════╝╚══════╝╚═╝  ╚═╝
$(GREEN)
╔══════════════════════════════════════════════════════════════════════════════╗
║  Resource Analyzer for Kubernetes - Monitore e otimize seus recursos         ║
╚══════════════════════════════════════════════════════════════════════════════╝
$(NC)
endef
export HEADER

# ==============================================================================
# Variáveis de Versão e Ambiente
# ==============================================================================
ENV?=dev
COMMIT_SHA=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
VERSION_TAG=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
VERSION_INFO="$(VERSION_TAG) ($(COMMIT_SHA)) - Built on $(BUILD_TIME)"

# ==============================================================================
# Variáveis de Segurança
# ==============================================================================
GOSEC_VERSION=v2.18.2
NANCY_VERSION=v1.0.45

# ==============================================================================
# Configuração do Make
# ==============================================================================
.PHONY: all analyze build build-version ci clean coverage debug debug-docker deps deps-check docker-build docker-push docker-run docs format health install-hooks install-tools lint logs metrics run run-dev run-prod setup setup-dev swagger test test-k8s test-mimir tidy validate version welcome help security-check security-deps
.DEFAULT_GOAL := help

# ==============================================================================
# Comandos de Desenvolvimento
# ==============================================================================
all: setup deps build ## Executa setup, deps e build em sequência

setup-dev: ## Configura ambiente de desenvolvimento completo
	@chmod +x ./scripts/dev/setup-dev.sh
	@./scripts/dev/setup-dev.sh

setup: welcome ## Configura ambiente básico
	@echo -e "🔧 $(BLUE)Configurando ambiente básico...$(NC)"
	@if [ ! -f .env ]; then \
		cp .env.example .env 2>/dev/null || echo -e "$(YELLOW)Arquivo .env.example não encontrado$(NC)"; \
	fi
	@$(MAKE) check-deps
	@$(MAKE) deps
	@$(MAKE) install-hooks
	@echo -e "✅ $(GREEN)Ambiente básico configurado!$(NC)"
	@echo -e "\n$(YELLOW)💡 Para uma configuração completa, execute:$(NC) make setup-dev"

install-hooks: ## Instala os git hooks
	@echo -e "🔧 $(BLUE)Instalando git hooks...$(NC)"
	@cp -f scripts/git/hooks/* .git/hooks/
	@chmod +x .git/hooks/*
	@echo -e "✅ $(GREEN)Git hooks instalados!$(NC)"

check-deps: ## Verifica e instala dependências de desenvolvimento
	@echo -e "🔍 $(BLUE)Verificando dependências de desenvolvimento...$(NC)"
	@command -v swag >/dev/null 2>&1 || { echo -e "$(YELLOW)Instalando swag...$(NC)" && go install github.com/swaggo/swag/cmd/swag@latest; }
	@command -v goimports >/dev/null 2>&1 || { echo -e "$(YELLOW)Instalando goimports...$(NC)" && go install golang.org/x/tools/cmd/goimports@latest; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo -e "$(YELLOW)Instalando golangci-lint...$(NC)" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	@command -v dlv >/dev/null 2>&1 || { echo -e "$(YELLOW)Instalando delve...$(NC)" && go install github.com/go-delve/delve/cmd/dlv@latest; }
	@echo -e "✅ $(GREEN)Todas as dependências estão instaladas!$(NC)"

build: welcome ## Build a versão local para desenvolvimento
	@echo -e "🔨 $(BLUE)Building$(NC) $(BOLD)$(APP_NAME)$(NC) versão $(VERSION)..."
	@$(GOBUILD) -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) $(MAIN_FILE)
	@echo -e "✅ $(GREEN)Build completed!$(NC)"

run: ## Roda a aplicação com a versão do git ou dev
	@echo -e "🚀 $(BLUE)Iniciando aplicação com versão $(VERSION)$(NC)..."
	@$(GORUN) -ldflags "$(LDFLAGS)" $(MAIN_FILE)

run-version: ## Roda a aplicação com uma versão específica (make run-version VERSION=1.0.0)
	@echo -e "🚀 $(BLUE)Iniciando aplicação com versão $(VERSION)$(NC)..."
	@$(GORUN) -ldflags "$(LDFLAGS)" $(MAIN_FILE)

run-dev: ## Executa a aplicação em modo desenvolvimento
	@echo -e "🚀 $(BLUE)Iniciando em modo desenvolvimento...$(NC)"
	@APP_ENV=development $(GORUN) -ldflags "$(LDFLAGS)" $(MAIN_FILE)

run-prod: ## Executa a aplicação em modo produção
	@echo -e "🚀 $(BLUE)Iniciando em modo produção...$(NC)"
	@APP_ENV=production $(GORUN) -ldflags "$(LDFLAGS)" $(MAIN_FILE)

clean: ## Limpa os arquivos de build
	@echo -e "🧹 $(YELLOW)Cleaning$(NC) build files..."
	@$(GOCLEAN)
	@rm -f $(APP_NAME)
	@echo -e "✨ $(GREEN)Cleanup completed!$(NC)"

# ==============================================================================
# Comandos de Teste e Qualidade
# ==============================================================================
test: ## Roda os testes
	@echo -e "🔍 $(BLUE)Running tests$(NC)..."
	@$(GOTEST) -v -race -cover ./...
	@echo -e "✅ $(GREEN)Tests completed!$(NC)"

test-k8s: ## Testa integração com Kubernetes
	@echo -e "🔍 $(BLUE)Testando integração com Kubernetes$(NC)..."
	@$(GOTEST) -v ./tests/integration/k8s/...
	@echo -e "✅ $(GREEN)Testes K8s completados!$(NC)"

test-mimir: ## Testa integração com Mimir
	@echo -e "🔍 $(BLUE)Testando integração com Mimir$(NC)..."
	@$(GOTEST) -v ./tests/integration/mimir/...
	@echo -e "✅ $(GREEN)Testes Mimir completados!$(NC)"

coverage: ## Roda os testes com cobertura
	@echo -e "📊 $(BLUE)Generating coverage report$(NC)..."
	@$(GOTEST) -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out
	@echo -e "📈 $(GREEN)Coverage report generated!$(NC)"

format: ## Formata o código
	@echo -e "✨ $(BLUE)Formatting code$(NC)..."
	@if ! command -v goimports >/dev/null 2>&1; then \
		go install golang.org/x/tools/cmd/goimports@latest; \
		export PATH="$$PATH:$$(go env GOPATH)/bin"; \
	fi
	@$$(go env GOPATH)/bin/goimports -l -w -d $(PROJECT_FILES)
	@gofmt -l -s -w $(PROJECT_FILES)
	@echo -e "✅ $(GREEN)Code formatted!$(NC)"

lint: ## Executa o linter
	@echo -e "🔍 $(BLUE)Running linter$(NC)..."
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run --timeout=180s
	@echo -e "✅ $(GREEN)Lint completed!$(NC)"

# ==============================================================================
# Comandos de Dependências
# ==============================================================================
deps: ## Baixa as dependências
	@echo -e "📦 $(BLUE)Downloading dependencies$(NC)..."
	@$(GOMOD) download
	@echo -e "✅ $(GREEN)Dependencies downloaded!$(NC)"

tidy: ## Organiza as dependências
	@echo -e "🔄 $(BLUE)Tidying dependencies$(NC)..."
	@$(GOMOD) tidy
	@echo -e "✨ $(GREEN)Dependencies organized!$(NC)"

# ==============================================================================
# Comandos Docker
# ==============================================================================
docker-build: ## Builda a imagem Docker
	@echo -e "🐳 $(BLUE)Building Docker image$(NC)..."
	@docker build -t $(DOCKER_IMAGE):$(VERSION) .
	@echo -e "✅ $(GREEN)Docker image built:$(NC) $(BOLD)$(DOCKER_IMAGE):$(VERSION)$(NC)"

docker-run: ## Roda o container Docker
	@echo -e "🐳 $(BLUE)Running Docker container$(NC)..."
	@docker run -t${INTERACTIVE} --rm \
		--name $(APP_NAME) \
		-p $(PORT):$(PORT) \
		-v $(PWD):$(APP_DIR):delegated \
		-v $(HOME)/.kube:/root/.kube:ro \
		$(DOCKER_IMAGE):$(VERSION)

docker-push: ## Push da imagem para o Docker Hub
	@echo -e "🚀 $(BLUE)Pushing$(NC) $(BOLD)$(DOCKER_IMAGE):$(VERSION)$(NC) to Docker Hub..."
	@docker push $(DOCKER_IMAGE):$(VERSION)
	@echo -e "✅ $(GREEN)Push completed!$(NC)"

# ==============================================================================
# Comandos de Debug
# ==============================================================================
debug: ## Roda a aplicação em modo debug localmente
	@echo -e "🔍 $(BLUE)Starting debugger$(NC) na porta $(YELLOW)$(DEBUG_PORT)$(NC)..."
	@echo -e "$(GREEN)➜ Conecte sua IDE ao debugger na porta $(DEBUG_PORT)$(NC)"
	@echo -e "$(GREEN)➜ VS Code: Pressione F5 ou use o menu de Debug$(NC)"
	@echo -e "$(GREEN)➜ GoLand: Use o botão de Debug$(NC)"
	@dlv debug $(MAIN_FILE) --headless --listen=:$(DEBUG_PORT) --api-version=2 --accept-multiclient

debug-docker: ## Roda a aplicação em modo debug no container
	@echo -e "🔍 $(BLUE)Starting debugger in container$(NC) na porta $(YELLOW)$(DEBUG_PORT)$(NC)..."
	@echo -e "$(GREEN)➜ Conecte sua IDE ao debugger na porta $(DEBUG_PORT)$(NC)"
	@docker run -t${INTERACTIVE} --rm \
		--name $(APP_NAME)-debug \
		-p $(PORT):$(PORT) \
		-p $(DEBUG_PORT):$(DEBUG_PORT) \
		-v $(PWD):/app:delegated \
		-v $(HOME)/.kube:/root/.kube:ro \
		--security-opt="apparmor=unconfined" \
		--security-opt="seccomp=unconfined" \
		--cap-add=SYS_PTRACE \
		$(DOCKER_IMAGE):$(VERSION) \
		dlv debug --headless --listen=:$(DEBUG_PORT) --api-version=2 --accept-multiclient --continue /app/cmd/api/main.go

# ==============================================================================
# Comandos de Documentação
# ==============================================================================
swagger: ## Gera a documentação Swagger
	@echo -e "📚 $(BLUE)Generating Swagger documentation$(NC)..."
	@command -v swag >/dev/null 2>&1 || { echo -e "$(RED)Swagger não encontrado. Instalando...$(NC)" && go install github.com/swaggo/swag/cmd/swag@latest; }
	@swag init -g $(MAIN_FILE) -o ./docs
	@echo -e "✅ $(GREEN)Swagger documentation generated!$(NC)"

# ==============================================================================
# Comandos de Utilidade
# ==============================================================================
welcome: ## Mostra o banner de boas-vindas
	@printf "$$HEADER"
	@echo -e "$(BLUE)Bem-vindo ao $(BOLD)$(APP_NAME)$(NC)"

help: welcome ## Mostra essa ajuda
	@echo -e "$(BLUE)Comandos disponíveis:$(NC)"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==============================================================================
# Comandos de Validação e CI/CD
# ==============================================================================
validate: format lint test ## Valida o código (formato, lint e testes)
	@echo -e "✅ $(GREEN)Validação completa com sucesso!$(NC)"

ci: validate build docker-build ## Pipeline de CI
	@echo -e "✅ $(GREEN)Pipeline de CI completado com sucesso!$(NC)"

# ==============================================================================
# Comandos de Monitoramento
# ==============================================================================
logs: ## Mostra logs da aplicação
	@echo -e "📋 $(BLUE)Buscando logs da aplicação...$(NC)"
	@if [ -f $(APP_NAME) ]; then \
		tail -f /var/log/$(APP_NAME).log 2>/dev/null || echo -e "$(YELLOW)Arquivo de log não encontrado$(NC)"; \
	else \
		docker logs -f $(APP_NAME) 2>/dev/null || echo -e "$(YELLOW)Container não encontrado$(NC)"; \
	fi

metrics: ## Mostra métricas básicas
	@echo -e "📊 $(BLUE)Coletando métricas...$(NC)"
	@echo -e "$(YELLOW)Estatísticas do Projeto:$(NC)"
	@echo -e "- Linhas de código Go: $$(find . -name '*.go' -not -path './vendor/*' | xargs wc -l | tail -n 1 | awk '{print $$1}')"
	@echo -e "- Número de arquivos Go: $$(find . -name '*.go' -not -path './vendor/*' | wc -l)"
	@echo -e "- Tamanho do binário: $$(ls -lh $(APP_NAME) 2>/dev/null | awk '{print $$5}' || echo 'N/A')"
	@echo -e "- Versão: $(shell git describe --tags --abbrev=0 2>/dev/null || echo 'v0.1.0') ($(shell git rev-parse --short HEAD)) - Built on $(shell date -u '+%Y-%m-%d_%H:%M:%S')"

health: ## Verifica saúde da aplicação
	@echo -e "🏥 $(BLUE)Verificando saúde da aplicação...$(NC)"
	@curl -s http://localhost:$(PORT)/api/v1/ping || echo -e "$(RED)Aplicação não está respondendo$(NC)"

# ==============================================================================
# Comandos de Análise
# ==============================================================================
analyze: lint test security-check ## Executa todas as análises (lint, test, security)
	@echo -e "🔍 $(BLUE)Iniciando análise completa...$(NC)"
	@$(MAKE) lint
	@$(MAKE) test
	@$(MAKE) security-check
	@$(MAKE) metrics
	@echo -e "✅ $(GREEN)Análise completa finalizada!$(NC)"

# ==============================================================================
# Comandos de Versionamento
# ==============================================================================
version: ## Mostra a versão atual
	@echo -e "📦 $(BLUE)Informações de Versão:$(NC)"
	@echo -e "$(YELLOW)Version Tag:$(NC) $(VERSION_TAG)"
	@echo -e "$(YELLOW)Commit:$(NC) $(COMMIT_SHA)"
	@echo -e "$(YELLOW)Build Time:$(NC) $(BUILD_TIME)"

# ==============================================================================
# Comandos de Documentação
# ==============================================================================
docs: swagger ## Gera toda a documentação
	@echo -e "📚 $(BLUE)Gerando documentação...$(NC)"
	@if [ -f "./docs/swagger.json" ]; then \
		echo -e "$(GREEN)✅ Swagger gerado em ./docs/swagger.json$(NC)"; \
	else \
		echo -e "$(RED)❌ Erro ao gerar Swagger$(NC)"; \
		exit 1; \
	fi

deps-check: ## Verifica dependências desatualizadas e vulnerabilidades
	@echo -e "🔍 $(BLUE)Verificando dependências...$(NC)"
	@chmod +x ./scripts/dev/check-deps.sh
	@./scripts/dev/check-deps.sh

# ==============================================================================
# Comandos de Desenvolvimento
# ==============================================================================
build-version: ## Build com uma versão específica (make build-version VERSION=1.0.0)
	@echo -e "🔨 $(BLUE)Building$(NC) $(BOLD)$(APP_NAME)$(NC) versão $(VERSION)..."
	@$(GOBUILD) -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) $(MAIN_FILE)
	@echo -e "✅ $(GREEN)Build completed!$(NC)"

# ==============================================================================
# Comandos de Segurança
# ==============================================================================
security-deps: ## Instala dependências de segurança
	@echo "🔒 Instalando ferramentas de segurança..."
	@go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)
	@go install github.com/sonatype-nexus-community/nancy@$(NANCY_VERSION)
	@echo "✅ Ferramentas de segurança instaladas!"

security-check: security-deps ## Executa verificações de segurança
	@echo "🔍 Executando análise de segurança..."
	@echo "🔒 Verificando código com gosec..."
	@gosec -quiet ./...
	@echo "📦 Verificando dependências com nancy..."
	@go list -json -deps | nancy sleuth
	@echo "✅ Verificações de segurança concluídas!"

install-tools: ## Instala todas as ferramentas de desenvolvimento
	@chmod +x ./scripts/dev/install-tools.sh
	@./scripts/dev/install-tools.sh 