# ADR 004: Princípios de Design de Código

## Status

Aceito

## Contexto

Durante o desenvolvimento do K8s Resource Analyzer, precisamos estabelecer princípios claros de design de código que guiem nossas decisões arquiteturais e de implementação. Os princípios escolhidos precisam:

- Manter o código simples e manutenível
- Evitar complexidade desnecessária
- Facilitar testes e manutenção
- Promover boas práticas de desenvolvimento

## Decisão

Decidimos adotar três princípios fundamentais de design:

1. **DRY (Don't Repeat Yourself)**
   - Cada conhecimento/lógica deve ter uma representação única no sistema
   - Implementado através de:
     * Pacote `pkg` para código compartilhado
     * Middlewares reutilizáveis
     * Tipos e interfaces comuns

2. **KISS (Keep It Simple, Stupid)**
   - Manter a simplicidade como prioridade
   - Implementado através de:
     * Funções pequenas e focadas
     * Estrutura de diretórios clara
     * Configuração via variáveis de ambiente
     * Respostas HTTP padronizadas

3. **YAGNI (You Aren't Gonna Need It)**
   - Não adicionar funcionalidades até que sejam realmente necessárias
   - Implementado através de:
     * Foco em requisitos atuais
     * Evitar abstrações prematuras
     * MVP primeiro, complexidade depois

## Consequências

### Positivas

1. Código mais limpo e manutenível
2. Menor complexidade acidental
3. Mais fácil de testar
4. Desenvolvimento mais rápido
5. Onboarding mais simples para novos desenvolvedores

### Negativas

1. Possível necessidade de refatoração quando novos requisitos surgirem
2. Potencial resistência a mudanças por manter código "simples demais"
3. Necessidade de balancear entre simplicidade e flexibilidade

## Exemplos de Aplicação

### DRY
```go
// pkg/response/response.go - Estrutura única para respostas HTTP
type Response struct {
    Success   bool        `json:"success"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Error     string      `json:"error,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
    RequestID string      `json:"request_id"`
}
```

### KISS
```go
// Configuração simples via variáveis de ambiente
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
```

### YAGNI
```go
// Apenas o necessário para o health check
func PingHandler(c *gin.Context) {
    timestamp := time.Now()
    response.Success(c, "pong", PingResponse{
        Status:    "ok",
        Timestamp: timestamp,
    })
}
```

## Referências

- [DRY Principle](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)
- [KISS Principle](https://en.wikipedia.org/wiki/KISS_principle)
- [YAGNI Principle](https://en.wikipedia.org/wiki/You_aren%27t_gonna_need_it)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) 