# ADR 006: Framework de Testes Kubernetes

## Status
Aceito

## Contexto
Para garantir a qualidade e confiabilidade da integração com Kubernetes, precisamos de uma solução de testes que:
- Simule um ambiente real
- Suporte métricas (kubectl top)
- Permita testes de integração completos
- Seja reproduzível e confiável
- Tenha boa performance

## Decisão
Utilizar Testcontainers com K3s pelos seguintes motivos:

1. **Realismo**
   - Kubernetes real em container
   - Métricas reais via metrics-server
   - Comportamento idêntico à produção
   - Validação de recursos real

2. **Isolamento**
   - Container novo para cada teste
   - Sem interferência entre testes
   - Limpeza automática
   - Estado inicial consistente

3. **Integração**
   - API Go nativa
   - Suporte a Docker compose
   - Gestão automática de ciclo de vida
   - Paralelização de testes

4. **Observabilidade**
   - Logs completos do cluster
   - Métricas reais
   - Debug facilitado
   - Inspeção de estado

## Alternativas Consideradas

1. **fake client**
   - ✅ Testes muito rápidos
   - ✅ Sem dependências externas
   - ❌ Não testa métricas reais
   - ❌ Comportamento simulado
   - ❌ Sem validação real

2. **kind**
   - ✅ Cluster real
   - ✅ Bom para CI
   - ❌ Setup manual
   - ❌ Mais pesado
   - ❌ Difícil paralelização

3. **minikube**
   - ✅ Ambiente completo
   - ✅ Mais próximo de produção
   - ❌ Muito pesado para testes
   - ❌ Setup complexo
   - ❌ Lento para iniciar

## Consequências

### Positivas
1. Testes mais confiáveis e realistas
2. Validação de métricas real
3. Ambiente isolado por teste
4. Fácil debug
5. CI/CD amigável

### Negativas
1. Testes mais lentos que fake client
2. Necessidade de Docker
3. Maior consumo de recursos
4. Setup inicial mais complexo

## Validação
- Tempo de execução dos testes
- Cobertura de cenários reais
- Estabilidade em CI/CD
- Facilidade de manutenção

## Referências
- [Testcontainers for Go](https://golang.testcontainers.org/)
- [K3s in Docker](https://k3s.io/)
- [Testing Kubernetes Operators](https://sdk.operatorframework.io/docs/building-operators/golang/testing/)
- [Kubernetes Testing Tools Comparison](https://kubernetes.io/docs/tasks/debug/debug-application/debug-running-pod/)