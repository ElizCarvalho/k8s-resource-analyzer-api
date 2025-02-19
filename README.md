# 🚀 K8s Resource Analyzer API

[🇺🇸 English Version](README.en.md)

> API HTTP em Go para análise de recursos Kubernetes com foco em FinOps.

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat-square&logo=go)
![Kubernetes](https://img.shields.io/badge/Kubernetes-Analyzer-326CE5?style=flat-square&logo=kubernetes)
![Swagger](https://img.shields.io/badge/Swagger-Documentation-85EA2D?style=flat-square&logo=swagger)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)
![Status](https://img.shields.io/badge/Status-In%20Development-yellow?style=flat-square)
[![CI](https://github.com/ElizCarvalho/k8s-resource-analyzer-api/actions/workflows/ci.yml/badge.svg)](https://github.com/ElizCarvalho/k8s-resource-analyzer-api/actions/workflows/ci.yml)
[![Release](https://github.com/ElizCarvalho/k8s-resource-analyzer-api/actions/workflows/release.yml/badge.svg)](https://github.com/ElizCarvalho/k8s-resource-analyzer-api/actions/workflows/release.yml)

<p align="center">
  <a href="#-sobre">Sobre</a> •
  <a href="#-status-do-projeto">Status</a> •
  <a href="#-funcionalidades">Funcionalidades</a> •
  <a href="#-tecnologias">Tecnologias</a> •
  <a href="#-início-rápido">Início Rápido</a> •
  <a href="#-api-endpoints">API</a>
</p>

</div>

<hr>

## 📌 Sobre

<div align="center">

```mermaid
graph LR
    A[Kubernetes Cluster] --> B[Resource Analyzer]
    B --> C[Métricas & Custos]
    C --> D[Insights FinOps]
    style A fill:#326CE5,stroke:#fff,stroke-width:2px,color:#fff
    style B fill:#00ADD8,stroke:#fff,stroke-width:2px,color:#fff
    style C fill:#85EA2D,stroke:#fff,stroke-width:2px,color:#fff
    style D fill:#2496ED,stroke:#fff,stroke-width:2px,color:#fff
```

</div>

O K8s Resource Analyzer é uma API desenvolvida em Go que permite analisar recursos do Kubernetes com foco em FinOps. A ferramenta fornece insights valiosos sobre utilização de recursos e custos em clusters Kubernetes.

## ⚡ Status do Projeto

| Status | Funcionalidade | Descrição |
|--------|----------------|-----------|
| ✅ | **Configuração Inicial** | Estrutura base do projeto implementada |
| ✅ | **Health Check** | Endpoint de verificação de saúde da API |
| ✅ | **Documentação** | OpenAPI/Swagger implementado |
| 🚧 | **Análise de Recursos** | Coleta e análise de recursos K8s |
| 🚧 | **Integração Metrics** | Conexão com Prometheus/Mimir |
| 🚧 | **Dashboard** | Visualização de métricas e custos |

## 🛠️ Stack Tecnológica

<table>
  <tr>
    <td align="center">
      <b>Core & API</b><br/>
      <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg" width="40" height="40"/><br/>
      <a href="https://go.dev/"><b>Go 1.22+ & Gin</b></a>
      <p align="center">
        • Integração nativa com client-go<br/>
        • Alta performance e baixa alocação<br/>
        • Middleware robusto e flexível<br/>
        • Execução concorrente
      </p>
      <p align="center">
        <code>Framework web de alta performance</code>
      </p>
    </td>
    <td align="center">
      <b>Observabilidade</b><br/>
      <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/prometheus/prometheus-original.svg" width="40" height="40"/><br/>
      <a href="https://grafana.com/oss/mimir/"><b>Mimir & Zerolog</b></a>
      <p align="center">
        • Métricas históricas K8s<br/>
        • Logs estruturados em JSON<br/>
        • Rastreamento por Request ID<br/>
        • Zero alocação em logs
      </p>
      <p align="center">
        <code>Monitoramento completo e eficiente</code>
      </p>
    </td>
    <td align="center">
      <b>Qualidade</b><br/>
      <img src="https://raw.githubusercontent.com/golangci/golangci-lint/master/assets/go.png" width="40" height="40"/><br/>
      <a href="https://golangci-lint.run/"><b>Ferramentas & Padrões</b></a>
      <p align="center">
        • Linting (golangci-lint)<br/>
        • Formatação (goimports)<br/>
        • Segurança (nancy)<br/>
        • Automação (Make)
      </p>
      <p align="center">
        <code>Garantia de qualidade de código</code>
      </p>
    </td>
  </tr>
  <tr>
    <td align="center">
      <b>Infraestrutura</b><br/>
      <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/docker/docker-original.svg" width="40" height="40"/><br/>
      <a href="https://www.docker.com/"><b>Container & CI/CD</b></a>
      <p align="center">
        • Docker multi-stage build<br/>
        • GitHub Actions Workflows<br/>
        • Deploy automatizado<br/>
        • Isolamento seguro
      </p>
      <p align="center">
        <code>Pipeline e deploy consistentes</code>
      </p>
    </td>
    <td align="center">
      <b>Documentação</b><br/>
      <img src="https://raw.githubusercontent.com/swagger-api/swagger.io/wordpress/images/assets/SW-logo-clr.png" width="40" height="40"/><br/>
      <a href="https://swagger.io/"><b>OpenAPI/Swagger</b></a>
      <p align="center">
        • Documentação interativa<br/>
        • Schemas bem definidos<br/>
        • Exemplos práticos<br/>
        • ADRs detalhadas
      </p>
      <p align="center">
        <code>Documentação clara e atualizada</code>
      </p>
    </td>
    <td align="center">
      <b>Ambiente</b><br/>
      <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/kubernetes/kubernetes-plain.svg" width="40" height="40"/><br/>
      <a href="https://kubernetes.io/"><b>Kubernetes & Cloud</b></a>
      <p align="center">
        • Análise de recursos K8s<br/>
        • Métricas de custos<br/>
        • Insights FinOps<br/>
        • Otimização de recursos
      </p>
      <p align="center">
        <code>Foco em eficiência e custos</code>
      </p>
    </td>
  </tr>
</table>

> **Nota**: Cada tecnologia foi escolhida considerando as necessidades específicas de análise de recursos Kubernetes e FinOps, priorizando performance, observabilidade e manutenibilidade.

## 📦 Estrutura do Projeto

```
k8s-resource-analyzer-api/
├── cmd/                    # Binários da aplicação
│   └── api/               # Ponto de entrada da API HTTP
├── internal/              # Código privado não exportável
│   ├── api/              # Implementação dos endpoints
│   └── pkg/              # Pacotes compartilhados
├── docs/                 # Documentação OpenAPI/Swagger
├── .env.example         # Template de configuração
├── Dockerfile          # Instruções de containerização
├── Makefile           # Automação de tarefas
└── README.md         # Documentação principal
```

## 🤝 Contribuindo

1. Fork o projeto
2. Crie sua branch de feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add: nova funcionalidade'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👩‍💻 Autora

Feito com ❤️ por Elizabeth Carvalho

[![LinkedIn](https://img.shields.io/badge/-Elizabeth%20Carvalho-blue?style=flat-square&logo=linkedin&logoColor=white&link=https://br.linkedin.com/in/elizcarvalho)](https://br.linkedin.com/in/elizcarvalho)
[![GitHub](https://img.shields.io/badge/-ElizCarvalho-gray?style=flat-square&logo=github&logoColor=white&link=https://github.com/ElizCarvalho)](https://github.com/ElizCarvalho)

## 📋 Pré-requisitos

<table>
  <tr>
    <td align="center">
      <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg" width="40" height="40"/><br/>
      <b>Go 1.22+</b>
    </td>
    <td align="center">
      <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/docker/docker-original.svg" width="40" height="40"/><br/>
      <b>Docker</b>
    </td>
    <td align="center">
      <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/kubernetes/kubernetes-plain.svg" width="40" height="40"/><br/>
      <b>Kubernetes</b>
    </td>
    <td align="center">
      <img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/prometheus/prometheus-original.svg" width="40" height="40"/><br/>
      <b>Prometheus</b>
    </td>
  </tr>
</table>

## 🚀 Início Rápido

```mermaid
graph LR
    A[Clone] --> B[Setup]
    B --> C[Configure]
    C --> D[Execute]
    style A fill:#00ADD8,stroke:#fff,stroke-width:2px,color:#fff
    style B fill:#2496ED,stroke:#fff,stroke-width:2px,color:#fff
    style C fill:#85EA2D,stroke:#fff,stroke-width:2px,color:#fff
    style D fill:#326CE5,stroke:#fff,stroke-width:2px,color:#fff
```

1. **Clone o repositório:**
```bash
git clone https://github.com/ElizCarvalho/k8s-resource-analyzer-api.git
cd k8s-resource-analyzer-api
```

2. **Instale as dependências:**
```bash
go mod download
```

3. **Configure as variáveis de ambiente:**
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

4. **Execute localmente:**
```bash
make run
```

5. **Ou com Docker:**
```bash
make docker-build
make docker-run
```

## 🔧 Configuração

### Variáveis de Ambiente

| Variável    | Descrição                   | Padrão  | Obrigatório |
|-------------|-----------------------------|---------|-------------|
| PORT        | Porta da API                | 9000    | Não         |
| GIN_MODE    | Modo do Gin (debug/release) | debug   | Não         |
| LOG_LEVEL   | Nível de log               | info    | Não         |
| LOG_FORMAT  | Formato dos logs (json/text)| json    | Não         |

## 📚 API Endpoints

### Health Check
- `GET /api/v1/ping` - Verifica o status da API
  - **Resposta de Sucesso**: `200 OK`
  - **Corpo**: `{"message": "pong", "status": "ok", "timestamp": "2024-02-18T00:00:00Z"}`

Documentação completa disponível em `/swagger/index.html`

## 🐳 Docker

### Build
```bash
docker build -t eliscarvalho/k8s-resource-analyzer-api:latest .
```

### Run
```bash
docker run -p 9000:9000 eliscarvalho/k8s-resource-analyzer-api:latest
```

### Docker Hub
```bash
docker pull eliscarvalho/k8s-resource-analyzer-api:latest
```

## 🧪 Testes

```bash
# Roda testes unitários
make test

# Roda testes com cobertura
make test-cover
```

## Funcionalidades

- Coleta de métricas atuais de CPU, memória e pods
- Histórico de utilização de recursos
- Análise de tendências de uso
- Integração com Mimir para armazenamento de métricas de longo prazo

## Requisitos

- Go 1.21 ou superior
- Kubernetes 1.19 ou superior
- Metrics Server instalado no cluster
- Mimir para armazenamento de métricas históricas

## Configuração

### Variáveis de Ambiente

- `KUBECONFIG`: Caminho para o arquivo kubeconfig (opcional, usado apenas fora do cluster)
- `IN_CLUSTER`: Define se a API está rodando dentro do cluster (`true` ou `false`)
- `MIMIR_URL`: URL do servidor Mimir
- `GIN_MODE`: Modo de execução do Gin (`debug` ou `release`)

## Instalação

### Local

1. Clone o repositório:
```bash
git clone https://github.com/ElizCarvalho/k8s-resource-analyzer-api.git
cd k8s-resource-analyzer-api
```

2. Instale as dependências:
```bash
go mod download
```

3. Execute a aplicação:
```bash
go run cmd/api/main.go
```

### Docker

1. Construa a imagem:
```bash
docker build -t k8s-resource-analyzer-api:latest .
```

2. Execute o container:
```bash
docker run -p 8080:8080 k8s-resource-analyzer-api:latest
```

### Kubernetes

1. Aplique os manifestos:
```bash
kubectl apply -f k8s/deployment.yaml
```

## Uso

### Obter Métricas

```bash
curl "http://localhost:8080/metrics?namespace=default&deployment=my-app&period=24h"
```

### Exemplo de Resposta

```json
{
  "current": {
    "cpu": {
      "average": 0.25,
      "peak": 0.45,
      "usage": 0.35,
      "request": 0.5,
      "limit": 1.0,
      "utilization": 70.0
    },
    "memory": {
      "average": 256.0,
      "peak": 512.0,
      "usage": 384.0,
      "request": 512.0,
      "limit": 1024.0,
      "utilization": 75.0
    },
    "pods": {
      "running": 3,
      "replicas": 3,
      "minReplicas": 2,
      "maxReplicas": 5
    }
  },
  "historical": {
    "cpu": [...],
    "memory": [...],
    "pods": [...],
    "period": "24h"
  },
  "trends": {
    "cpu": {
      "trend": 0.15,
      "confidence": 0.95,
      "period": "24h"
    },
    "memory": {
      "trend": 0.08,
      "confidence": 0.92,
      "period": "24h"
    },
    "pods": {
      "trend": 0.0,
      "confidence": 1.0,
      "period": "24h"
    }
  },
  "metadata": {
    "collectedAt": "2024-03-20T10:30:00Z",
    "timeWindow": "24h"
  }
}
```

## Documentação da API

A documentação da API está disponível em formato OpenAPI/Swagger em `/docs/swagger.yaml`.

## Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Crie um Pull Request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.