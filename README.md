# 🚀 K8s Resource Analyzer API

> API REST em Go para análise de recursos Kubernetes com foco em FinOps.

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat-square&logo=go)
![Kubernetes](https://img.shields.io/badge/Kubernetes-Analyzer-326CE5?style=flat-square&logo=kubernetes)
![Swagger](https://img.shields.io/badge/Swagger-Documentation-85EA2D?style=flat-square&logo=swagger)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)
![Status](https://img.shields.io/badge/Status-In%20Development-yellow?style=flat-square)

</div>

<hr>

<p align="center">
  <a href="#-sobre">Sobre</a> •
  <a href="#-status-do-projeto">Status</a> •
  <a href="#-funcionalidades">Funcionalidades</a> •
  <a href="#-tecnologias">Tecnologias</a> •
  <a href="#-início-rápido">Início Rápido</a> •
  <a href="#-api-endpoints">API</a>
</p>

<hr>

## 📌 Sobre

O K8s Resource Analyzer é uma API desenvolvida em Go que permite analisar recursos do Kubernetes com foco em FinOps. A ferramenta fornece insights valiosos sobre utilização de recursos e custos em clusters Kubernetes.

## ⚡ Status do Projeto

🚧 **Em Desenvolvimento** 🚧

- [x] Configuração inicial do projeto
- [x] Implementação do health check
- [x] Documentação Swagger
- [ ] Análise de recursos Kubernetes
- [ ] Integração com Prometheus/Mimir
- [ ] Dashboard de métricas

## 🎯 Funcionalidades

- Documentação Swagger interativa
- Endpoints RESTful
- Health Check e monitoramento
- Suporte a múltiplos ambientes via variáveis de ambiente

## 🛠️ Tecnologias

- [Go 1.22+](https://go.dev/) - Linguagem de programação
- [Gin](https://gin-gonic.com/) - Web Framework
- [Swagger](https://swagger.io/) - Documentação API
- [Docker](https://www.docker.com/) - Containerização

## 📋 Pré-requisitos

- Go 1.22 ou superior
- Docker
- Make (opcional, para comandos de desenvolvimento)

## 🚀 Início Rápido

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

## 📚 API Endpoints

### Health Check
- `GET /api/v1/ping` - Verifica o status da API
  - **Resposta de Sucesso**: `200 OK`
  - **Corpo**: `{"message": "pong", "status": "ok", "timestamp": "2024-02-18T00:00:00Z"}`

Documentação completa disponível em `/swagger/index.html`

## 🐳 Docker

### Build
```bash
docker build -t ecarvalho2020/k8s-resource-analyzer-api:latest .
```

### Run
```bash
docker run -p 9000:9000 ecarvalho2020/k8s-resource-analyzer-api:latest
```

### Docker Hub
```bash
docker pull ecarvalho2020/k8s-resource-analyzer-api:latest
```

## 🧪 Testes

```bash
# Roda testes unitários
make test

# Roda testes com cobertura
make test-cover
```

## 📦 Estrutura do Projeto

```
k8s-resource-analyzer-api/
├── cmd/                    # Ponto de entrada da aplicação
│   └── api/               # Arquivo main.go e configurações
├── internal/              # Código privado da aplicação
│   └── api/              # Handlers e rotas da API
├── docs/                 # Documentação (Swagger)
├── .env.example         # Exemplo de variáveis de ambiente
├── Dockerfile          # Configuração Docker
├── Makefile           # Comandos de desenvolvimento
└── README.md         # Este arquivo
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