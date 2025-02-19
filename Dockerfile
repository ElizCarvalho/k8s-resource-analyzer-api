# Build stage
FROM golang:1.21-alpine AS builder

# Instala dependências de build
RUN apk add --no-cache git

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o k8s-resource-analyzer-api ./cmd/api

# Final stage
FROM alpine:3.19

# Instala certificados CA
RUN apk add --no-cache ca-certificates

# Define o diretório de trabalho
WORKDIR /app

# Copia o binário compilado
COPY --from=builder /app/k8s-resource-analyzer-api .

# Expõe a porta da API
EXPOSE 8080

# Define as variáveis de ambiente padrão
ENV GIN_MODE=release
ENV IN_CLUSTER=true

# Executa a aplicação
CMD ["./k8s-resource-analyzer-api"] 