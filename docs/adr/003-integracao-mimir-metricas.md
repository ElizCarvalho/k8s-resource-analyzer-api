# ADR 003: Integração com Mimir para Métricas

## Status
Aceito

## Contexto
Para o K8s Resource Analyzer, precisamos de uma fonte de métricas que atenda aos seguintes requisitos:
- Histórico de métricas
- Compatibilidade Kubernetes
- Confiabilidade dos dados
- Performance em queries
- Facilidade de integração

## Decisão
Decidimos integrar com o Mimir existente pelos seguintes motivos:

1. **Disponibilidade**
   - Já coleta métricas
   - Histórico existente
   - Dados confiáveis
   - Infraestrutura pronta

2. **Compatibilidade**
   - API PromQL
   - Queries complexas
   - Ecossistema Prometheus
   - Ferramentas existentes

3. **Praticidade**
   - Setup mínimo
   - Infraestrutura pronta
   - Implantação simples
   - Equipe familiarizada

## Alternativas Consideradas

1. **Prometheus Direto**
   - ✅ Mais simples
   - ✅ Menor latência
   - ❌ Sem long-term storage
   - ❌ Limitação de escala
   - ❌ Setup adicional

2. **Thanos**
   - ✅ Alta disponibilidade
   - ✅ Global view
   - ❌ Complexidade extra
   - ❌ Overhead de recursos
   - ❌ Setup complexo

3. **VictoriaMetrics**
   - ✅ Alta performance
   - ✅ Boa compressão
   - ❌ Menos maduro
   - ❌ Comunidade menor
   - ❌ Migração necessária

## Consequências

### Positivas
1. Acesso imediato a dados
2. Infraestrutura testada
3. Equipe familiarizada
4. Setup simplificado
5. Manutenção mínima

### Negativas
1. Dependência do Mimir
2. Compatibilidade de versões
3. Limitações de queries
4. Overhead de rede

## Validação
- Verificar existência e disponibilidade do Mimir no ambiente
- Validar período de retenção das métricas (mínimo 30 dias)
- Confirmar disponibilidade das métricas necessárias:
  * container_cpu_usage_seconds_total
  * container_memory_working_set_bytes
  * kube_pod_container_resource_requests
  * kube_pod_container_resource_limits
- Testar latência das queries em diferentes cenários:
  * Última hora
  * Último dia
  * Última semana
  * Último mês
- Verificar limites de rate limiting e quotas
- Validar formato e precisão dos dados retornados
- Testar recuperação em caso de falhas de conexão

## Referências
- [Documentação Mimir](https://grafana.com/docs/mimir/latest/)
- [PromQL para K8s](https://prometheus.io/docs/prometheus/latest/querying/examples/#kubernetes-container-resources)
- [Mimir vs Alternatives](https://grafana.com/docs/mimir/latest/mimir-vs-alternatives/)
- [Kubernetes Metrics](https://kubernetes.io/docs/concepts/cluster-administration/monitoring/) 