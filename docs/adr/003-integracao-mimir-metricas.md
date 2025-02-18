# ADR 003: Integração com Mimir para Métricas de Recursos Kubernetes

## Status

Aceito

## Contexto

O K8s Resource Analyzer precisa acessar métricas históricas de utilização de recursos (CPU, memória) dos pods Kubernetes. No ambiente onde a prova de conceito está sendo desenvolvida, o Mimir é utilizado como fonte de dados históricos dos recursos dos pods.

## Decisão

Decidimos integrar com o Mimir existente no ambiente pelos seguintes motivos:

1. **Disponibilidade dos Dados**
   - O Mimir já está configurado e coletando métricas dos clusters
   - Contém o histórico necessário de utilização de recursos dos pods
   - Dados já estão sendo coletados e armazenados de forma confiável

2. **Compatibilidade**
   - O Mimir é compatível com a API PromQL
   - Permite consultas complexas sobre o histórico de recursos
   - Mantém compatibilidade com ferramentas do ecossistema Prometheus

3. **Praticidade**
   - Não há necessidade de configurar uma nova fonte de dados
   - Aproveita a infraestrutura existente
   - Reduz complexidade de implantação

## Consequências

### Positivas

1. Acesso imediato ao histórico de métricas
2. Infraestrutura já estabelecida e testada
3. Familiaridade da equipe com a solução
4. Sem necessidade de configurar coleta de métricas

### Negativas

1. Dependência da disponibilidade do Mimir
2. Necessidade de manter compatibilidade com a versão do Mimir utilizada
3. Possíveis limitações nas consultas baseadas na configuração existente

## Implementação

A integração é feita através de variáveis de ambiente:

```env
# URL do Mimir para consulta de métricas históricas dos pods Kubernetes
MIMIR_URL=http://localhost:9090

# Credenciais para autenticação no Mimir (se necessário)
MIMIR_USERNAME=
MIMIR_PASSWORD=
```

Exemplo de consulta para obter utilização de CPU:
```promql
rate(container_cpu_usage_seconds_total{container!=""}[5m])
```

## Referências

- [Documentação do Mimir](https://grafana.com/docs/mimir/latest/)
- [PromQL para Métricas de Kubernetes](https://prometheus.io/docs/prometheus/latest/querying/examples/#kubernetes-container-resources) 