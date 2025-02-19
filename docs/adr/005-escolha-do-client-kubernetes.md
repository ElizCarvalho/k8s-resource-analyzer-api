# ADR 005: Escolha do Client Kubernetes

## Status
Aceito

## Contexto
Para interagir com o cluster Kubernetes, precisamos escolher uma biblioteca cliente que atenda aos seguintes requisitos:
- Performance em operações de leitura
- Facilidade de manutenção
- Suporte a métricas (kubectl top)
- Baixo overhead
- Compatibilidade com versões do Kubernetes

## Decisão
Utilizar `client-go` diretamente pelos seguintes motivos:

1. **Performance**
   - Menos camadas de abstração
   - Otimizado para operações de leitura
   - Menor overhead de memória
   - Cache eficiente

2. **Simplicidade**
   - API mais próxima do Kubernetes
   - Menos dependências
   - Curva de aprendizado mais suave
   - Segue princípio KISS

3. **Flexibilidade**
   - Controle fino sobre operações
   - Facilidade de mock para testes
   - Configuração granular
   - Melhor observabilidade

## Alternativas Consideradas

1. **controller-runtime**
   - ✅ Mais abstrações prontas
   - ✅ Melhor para operadores
   - ❌ Overhead desnecessário
   - ❌ Muitas dependências
   - ❌ Complexidade adicional

2. **k8s.io/api direto**
   - ✅ Máxima performance
   - ✅ Mínimo de dependências
   - ❌ Muito código boilerplate
   - ❌ Difícil manutenção
   - ❌ Sem abstrações úteis

3. **kubernetes/apimachinery**
   - ✅ Baixo nível e flexível
   - ✅ Bom controle
   - ❌ Complexidade desnecessária
   - ❌ Muito manual
   - ❌ Pouca documentação

## Consequências

### Positivas
1. Menor footprint de memória
2. Performance otimizada para leitura
3. Facilidade de debug
4. Melhor testabilidade
5. Documentação abundante

### Negativas
1. Mais código boilerplate
2. Menos abstrações prontas
3. Necessidade de mais testes unitários
4. Manutenção de código de conexão

## Validação
- Testes de performance
- Comparação de uso de memória
- Facilidade de implementação
- Cobertura de testes

## Referências
- [client-go Examples](https://github.com/kubernetes/client-go/tree/master/examples)
- [Kubernetes API Concepts](https://kubernetes.io/docs/reference/using-api/api-concepts/)
- [client-go vs controller-runtime](https://cloud.redhat.com/blog/kubernetes-clients-comparing-client-go-and-controller-runtime)
- [client-go Best Practices](https://github.com/kubernetes/client-go/blob/master/INSTALL.md#best-practices)