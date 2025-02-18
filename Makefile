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
DOCKER_IMAGE=ecarvalho2020/$(APP_NAME)
VERSION?=latest
PORT?=9000
DEBUG_PORT?=2345
APP_DIR=/usr/app

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
# Configuração do Make
# ==============================================================================
.PHONY: all build clean test coverage deps run docker-build docker-run docker-push swagger help colors welcome setup format lint debug debug-docker check-deps validate ci cd env-check env-setup security metrics health backup snapshot monitor install-tools analyze version tag release release-docker ci-check cd-check clean-all docs deps-check
.DEFAULT_GOAL := help

# ==============================================================================
# Comandos de Desenvolvimento
# ==============================================================================
setup-dev: ## Configura ambiente de desenvolvimento completo
	@chmod +x ./scripts/setup-dev.sh
	@./scripts/setup-dev.sh

setup: welcome ## Configura ambiente básico
	@echo -e "🔧 $(BLUE)Configurando ambiente básico...$(NC)"
	@if [ ! -f .env ]; then \
		cp .env.example .env 2>/dev/null || echo -e "$(YELLOW)Arquivo .env.example não encontrado$(NC)"; \
	fi
	@$(MAKE) check-deps
	@$(MAKE) deps
	@echo -e "✅ $(GREEN)Ambiente básico configurado!$(NC)"
	@echo -e "\n$(YELLOW)💡 Para uma configuração completa, execute:$(NC) make setup-dev"

check-deps: ## Verifica e instala dependências de desenvolvimento
	@echo -e "🔍 $(BLUE)Verificando dependências de desenvolvimento...$(NC)"
	@command -v swag >/dev/null 2>&1 || { echo -e "$(YELLOW)Instalando swag...$(NC)" && go install github.com/swaggo/swag/cmd/swag@latest; }
	@command -v goimports >/dev/null 2>&1 || { echo -e "$(YELLOW)Instalando goimports...$(NC)" && go install golang.org/x/tools/cmd/goimports@latest; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo -e "$(YELLOW)Instalando golangci-lint...$(NC)" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	@command -v dlv >/dev/null 2>&1 || { echo -e "$(YELLOW)Instalando delve...$(NC)" && go install github.com/go-delve/delve/cmd/dlv@latest; }
	@echo -e "✅ $(GREEN)Todas as dependências estão instaladas!$(NC)"

build: welcome ## Build a versão local para desenvolvimento
	@echo -e "🔨 $(BLUE)Building$(NC) $(BOLD)$(APP_NAME)$(NC)..."
	@$(GOBUILD) -o $(APP_NAME) $(MAIN_FILE)
	@echo -e "✅ $(GREEN)Build completed!$(NC)"

run: ## Roda a aplicação localmente
	@echo -e "🚀 $(BLUE)Starting$(NC) $(BOLD)$(APP_NAME)$(NC) on port $(YELLOW)$(PORT)$(NC)..."
	@PORT=$(PORT) $(GORUN) $(MAIN_FILE)

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

coverage: ## Roda os testes com cobertura
	@echo -e "📊 $(BLUE)Generating coverage report$(NC)..."
	@$(GOTEST) -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out
	@echo -e "📈 $(GREEN)Coverage report generated!$(NC)"

format: ## Formata o código
	@echo -e "✨ $(BLUE)Formatting code$(NC)..."
	@command -v goimports >/dev/null 2>&1 || go install golang.org/x/tools/cmd/goimports@latest
	@goimports -l -w -d $(PROJECT_FILES)
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

colors: ## Demonstração de cores disponíveis
	@echo -e "=== 🎨 $(BLUE)Demonstração de Cores$(NC) ==="
	@echo -e "$(RED)Texto em Vermelho$(NC)"
	@echo -e "$(GREEN)Texto em Verde$(NC)"
	@echo -e "$(YELLOW)Texto em Amarelo$(NC)"
	@echo -e "$(BLUE)Texto em Azul$(NC)"
	@echo -e "$(BOLD)Texto em Negrito$(NC)"
	@echo -e "$(RED)$(BOLD)Texto em Vermelho e Negrito$(NC)"
	@echo -e "$(GREEN)$(BOLD)Texto em Verde e Negrito$(NC)"

help: welcome ## Mostra essa ajuda
	@echo -e "$(BLUE)Comandos disponíveis:$(NC)"
	@echo
	@printf "Uso:\n  make \033[36m<target>\033[0m\n\nTargets:\n"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ==============================================================================
# Comandos de Validação e CI/CD
# ==============================================================================
validate: format lint test ## Valida o código (formato, lint e testes)
	@echo -e "✅ $(GREEN)Validação completa com sucesso!$(NC)"

ci: validate build docker-build ## Pipeline de CI
	@echo -e "✅ $(GREEN)Pipeline de CI completado com sucesso!$(NC)"

cd: docker-push ## Pipeline de CD
	@echo -e "🚀 $(BLUE)Iniciando deploy...$(NC)"
	@echo -e "📦 Version: $(VERSION_INFO)"
	@echo -e "🌍 Environment: $(ENV)"
	@echo -e "✅ $(GREEN)Deploy completado com sucesso!$(NC)"

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
	@echo -e "- Versão: $(VERSION_INFO)"

health: ## Verifica saúde da aplicação
	@echo -e "🏥 $(BLUE)Verificando saúde da aplicação...$(NC)"
	@curl -s http://localhost:$(PORT)/api/v1/ping || echo -e "$(RED)Aplicação não está respondendo$(NC)"

# ==============================================================================
# Comandos de Backup
# ==============================================================================
backup: ## Backup de configurações
	@echo -e "💾 $(BLUE)Criando backup...$(NC)"
	@mkdir -p ./backups/$(shell date +%Y%m%d_%H%M%S)
	@cp .env* ./backups/$(shell date +%Y%m%d_%H%M%S)/ 2>/dev/null || true
	@cp config.* ./backups/$(shell date +%Y%m%d_%H%M%S)/ 2>/dev/null || true
	@echo -e "✅ $(GREEN)Backup criado em ./backups/$(shell date +%Y%m%d_%H%M%S)$(NC)"

snapshot: ## Snapshot do estado atual
	@echo -e "📸 $(BLUE)Criando snapshot do projeto...$(NC)"
	@tar -czf ./backups/snapshot_$(shell date +%Y%m%d_%H%M%S).tar.gz \
		--exclude='.git' \
		--exclude='vendor' \
		--exclude='node_modules' \
		--exclude='*.log' \
		--exclude='backups' \
		.
	@echo -e "✅ $(GREEN)Snapshot criado em ./backups/snapshot_$(shell date +%Y%m%d_%H%M%S).tar.gz$(NC)"

# ==============================================================================
# Comandos de Monitoramento
# ==============================================================================
monitor: ## Inicia o monitoramento interativo
	@chmod +x ./scripts/monitor.sh
	@./scripts/monitor.sh

install-tools: ## Instala todas as ferramentas de desenvolvimento
	@chmod +x ./scripts/install-tools.sh
	@./scripts/install-tools.sh

# ==============================================================================
# Comandos de Análise
# ==============================================================================
analyze: ## Executa todas as análises (lint, segurança, testes)
	@echo -e "🔍 $(BLUE)Iniciando análise completa...$(NC)"
	@$(MAKE) lint
	@$(MAKE) security-check
	@$(MAKE) test
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

tag: ## Cria uma nova tag de versão
	@echo -e "🏷️ $(BLUE)Criando nova tag...$(NC)"
	@read -p "Nova versão (atual: $(VERSION_TAG)): " version; \
	git tag -a $$version -m "Release $$version"
	@echo -e "✅ $(GREEN)Tag criada com sucesso!$(NC)"
	@echo -e "💡 Execute 'git push --tags' para publicar a tag"

# ==============================================================================
# Comandos de Release
# ==============================================================================
release: ## Inicia o processo de release
	@chmod +x ./scripts/release.sh
	@./scripts/release.sh

release-docker: ## Cria e publica uma nova versão Docker
	@echo -e "🐳 $(BLUE)Iniciando release Docker...$(NC)"
	@$(MAKE) docker-build
	@$(MAKE) docker-push
	@echo -e "✅ $(GREEN)Release Docker completada!$(NC)"

# ==============================================================================
# Comandos de CI/CD
# ==============================================================================
ci-check: ## Verifica se o código está pronto para CI
	@echo -e "🔍 $(BLUE)Verificando código para CI...$(NC)"
	@$(MAKE) format
	@$(MAKE) lint
	@$(MAKE) test
	@$(MAKE) security-check
	@echo -e "✅ $(GREEN)Código pronto para CI!$(NC)"

cd-check: ## Verifica se está pronto para deploy
	@echo -e "🚀 $(BLUE)Verificando pré-requisitos para deploy...$(NC)"
	@$(MAKE) version
	@$(MAKE) health
	@echo -e "✅ $(GREEN)Pronto para deploy!$(NC)"

# ==============================================================================
# Comandos de Limpeza
# ==============================================================================
clean-all: clean ## Limpa todos os arquivos gerados
	@echo -e "🧹 $(BLUE)Limpeza profunda...$(NC)"
	@rm -rf ./backups/* ./logs/* ./tmp/*
	@rm -f coverage.out release_notes.md
	@docker rmi $(DOCKER_IMAGE):$(VERSION) 2>/dev/null || true
	@echo -e "✨ $(GREEN)Limpeza completa!$(NC)"

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
	@chmod +x ./scripts/check-deps.sh
	@./scripts/check-deps.sh 