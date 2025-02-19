# ADR 002: Zerolog para Logging Estruturado

## Status
Aceito

## Contexto
Para o K8s Resource Analyzer, precisamos de uma solução de logging que atenda aos seguintes requisitos:
- Logging estruturado em JSON
- Alta performance
- Baixa alocação de memória
- Flexibilidade de formatos
- Ergonomia para desenvolvedores

## Decisão
Decidimos utilizar o Zerolog pelos seguintes motivos:

1. **Performance**
   - Zero alocação para logs comuns
   - Otimizado para JSON
   - Benchmarks superiores
   - Baixo overhead

2. **Funcionalidades**
   - Logging estruturado nativo
   - Múltiplos formatos de saída
   - Níveis configuráveis
   - Campos contextuais

3. **Ergonomia**
   - API fluente
   - Campos estruturados
   - Boa integração Go
   - Documentação clara

## Alternativas Consideradas

1. **Logrus**
   - ✅ Estabelecido no mercado
   - ✅ Muitos plugins
   - ❌ Maior uso de memória
   - ❌ Performance inferior
   - ❌ Manutenção mais lenta

2. **Zap**
   - ✅ Alta performance
   - ✅ Bem testado
   - ❌ API mais complexa
   - ❌ Curva de aprendizado
   - ❌ Setup verboso

3. **Go standard log**
   - ✅ Simplicidade máxima
   - ✅ Nativo da linguagem
   - ❌ Sem JSON nativo
   - ❌ Sem estruturação
   - ❌ Recursos limitados

## Consequências

### Positivas
1. Logs fáceis de analisar
2. Mínimo overhead
3. Integração simplificada
4. Boa DX
5. Formato padronizado

### Negativas
1. Setup para human-readable
2. Curva inicial de aprendizado
3. Overhead de CPU
4. Necessidade de formatação

## Validação
- Testes de performance
- Facilidade de uso
- Integração com ferramentas
- Feedback da equipe

## Referências
- [Zerolog GitHub](https://github.com/rs/zerolog)
- [Zerolog vs Others](https://github.com/rs/zerolog#benchmarks)
- [Structured Logging](https://www.honeybadger.io/blog/golang-logging/)
- [Go Logging Best Practices](https://www.loggly.com/blog/logging-in-golang/) 