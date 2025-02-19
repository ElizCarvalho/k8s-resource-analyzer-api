# ADR 008: Estratégia de Namespaces

## Status
Aceito

## Contexto
A aplicação precisa analisar recursos Kubernetes considerando:
- Isolamento de recursos
- Multi-tenancy
- Controle de acesso
- Performance de queries
- Escalabilidade futura

## Decisão
Implementar uma estratégia de namespace único configurável com preparação para multi-namespace:

1. **Fase 1: Namespace Único**
   - Namespace definido via configuração
   - Validação de existência
   - Métricas por namespace
   - Labels para identificação

2. **Preparação Multi-namespace**
   - Interface abstraindo acesso
   - Estruturas de dados preparadas
   - Métricas agregadas
   - Cache por namespace

## Implementação

### Interface
```go
type NamespaceStrategy interface {
    // Lista namespaces disponíveis
    ListNamespaces(ctx context.Context) ([]string, error)
    
    // Valida acesso ao namespace
    ValidateAccess(ctx context.Context, namespace string) error
    
    // Obtém métricas do namespace
    GetMetrics(ctx context.Context, namespace string) (*Metrics, error)
    
    // Agrega métricas de múltiplos namespaces
    AggregateMetrics(ctx context.Context, namespaces []string) (*AggregatedMetrics, error)
}
```

### Configuração
```yaml
namespaces:
  # Modo inicial: single
  mode: single
  # Namespace padrão
  default: monitoring
  # Cache
  cache:
    enabled: true
    ttl: 5m
  # Métricas
  metrics:
    aggregation: true
    retention: 7d
```

## Alternativas Consideradas

1. **Multi-namespace desde início**
   - ✅ Mais completo
   - ✅ Sem refatoração futura
   - ❌ Complexidade inicial
   - ❌ Overhead desnecessário

2. **Namespace fixo**
   - ✅ Máxima simplicidade
   - ✅ Performance otimizada
   - ❌ Sem flexibilidade
   - ❌ Difícil evolução

3. **Namespace por tenant**
   - ✅ Isolamento total
   - ✅ Segurança melhor
   - ❌ Complexidade de gestão
   - ❌ Overhead de recursos

## Consequências

### Positivas
1. Simplicidade inicial
2. Preparado para evolução
3. Performance otimizada
4. Fácil manutenção
5. Baixo overhead

### Negativas
1. Limitação inicial multi-tenant
2. Possível refatoração futura
3. Complexidade crescente
4. Necessidade de migração

## Validação
- Testes de performance
- Validação de isolamento
- Testes de escalabilidade
- Simulação multi-tenant

## Referências
- [Kubernetes Namespaces](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
- [Multi-tenancy Best Practices](https://kubernetes.io/docs/concepts/security/multi-tenancy/)
- [Namespace Patterns](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces-overview/)
- [Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/)