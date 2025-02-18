# ADR 002: Zerolog para Logging Estruturado

## Status

Aceito

## Contexto

Para o K8s Resource Analyzer, precisávamos de uma solução de logging que atendesse aos seguintes requisitos:

- Logging estruturado em formato JSON para fácil integração com ferramentas de análise
- Alta performance para não impactar o processamento de métricas
- Baixa alocação de memória
- Flexibilidade para diferentes formatos de saída
- Facilidade de uso e boa ergonomia para desenvolvedores

## Decisão

Decidimos utilizar o [Zerolog](https://github.com/rs/zerolog) como biblioteca de logging pelos seguintes motivos:

1. **Performance**
   - Zero alocação de memória para logs comuns
   - Otimizado para logging em JSON
   - Benchmarks superiores comparado a outras bibliotecas

2. **Funcionalidades**
   - Logging estruturado nativo em JSON
   - Suporte a diferentes formatos de saída (JSON, console)
   - Níveis de log configuráveis
   - Campos contextuais
   - Sampling de logs

3. **Ergonomia**
   - API fluente e intuitiva
   - Fácil adição de campos estruturados
   - Boa integração com o ecossistema Go
   - Documentação clara e exemplos práticos

4. **Flexibilidade**
   - Customização de formatadores
   - Hooks para processamento de logs
   - Integração com writers padrão do Go

## Consequências

### Positivas

1. Logs estruturados facilitam análise e busca
2. Mínimo overhead de performance
3. Fácil integração com ferramentas de análise de logs
4. Boa experiência de desenvolvimento
5. Formato JSON padronizado

### Negativas

1. Necessidade de configuração adicional para formato human-readable
2. Curva de aprendizado inicial para logging estruturado
3. Overhead de CPU ligeiramente maior que logging simples

## Alternativas Consideradas

1. **Logrus**
   - Mais antigo e estabelecido
   - Maior overhead de memória
   - Menos performático

2. **Zap**
   - Performance similar
   - API mais complexa
   - Maior curva de aprendizado

3. **Go standard log**
   - Mais simples
   - Sem suporte nativo a JSON
   - Sem recursos avançados

## Implementação

Exemplo de uso no projeto:

```go
logger.Info().
    Str("handler", "ping").
    Str("method", "GET").
    Str("path", "/ping").
    Str("ip", c.ClientIP()).
    Msg("Recebida requisição ping")
```

Saída JSON:
```json
{
  "level": "info",
  "app": "k8s-resource-analyzer",
  "env": "development",
  "handler": "ping",
  "method": "GET",
  "path": "/ping",
  "ip": "127.0.0.1",
  "message": "Recebida requisição ping",
  "time": "2024-02-18T15:04:05Z"
}
```

## Referências

- [Zerolog GitHub](https://github.com/rs/zerolog)
- [Zerolog vs Other Loggers Benchmark](https://github.com/rs/zerolog#benchmarks)
- [Structured Logging in Go](https://www.honeybadger.io/blog/golang-logging/) 