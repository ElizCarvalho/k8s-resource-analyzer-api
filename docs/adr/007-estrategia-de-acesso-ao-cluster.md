# ADR 007: Estratégia de Acesso ao Cluster

## Status
Aceito

## Contexto
A aplicação precisa acessar o cluster Kubernetes de forma segura e flexível, considerando:
- Diferentes ambientes (dev, staging, prod)
- Segurança e princípio do menor privilégio
- Facilidade de desenvolvimento
- Compatibilidade com CI/CD
- Auditoria de acessos

## Decisão
Implementar uma estratégia dual de autenticação:

1. **Desenvolvimento (kubeconfig)**
   - Uso do arquivo ~/.kube/config
   - Suporte a múltiplos contextos
   - Override via KUBECONFIG
   - Configuração via flags

2. **Produção (ServiceAccount)**
   - ServiceAccount dedicada
   - RBAC com permissões mínimas
   - Montagem automática de token
   - Namespace isolado

## Implementação

### ServiceAccount
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: resource-analyzer
  namespace: monitoring

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: resource-analyzer
rules:
- apiGroups: [""]
  resources: ["pods", "nodes"]
  verbs: ["get", "list"]
- apiGroups: ["metrics.k8s.io"]
  resources: ["pods", "nodes"]
  verbs: ["get", "list"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: resource-analyzer
subjects:
- kind: ServiceAccount
  name: resource-analyzer
  namespace: monitoring
roleRef:
  kind: ClusterRole
  name: resource-analyzer
  apiGroup: rbac.authorization.k8s.io
```

## Alternativas Consideradas

1. **Apenas ServiceAccount**
   - ✅ Mais seguro
   - ✅ Padrão Kubernetes
   - ❌ Desenvolvimento mais complexo
   - ❌ Setup local trabalhoso

2. **Apenas kubeconfig**
   - ✅ Familiar para desenvolvedores
   - ✅ Fácil de configurar
   - ❌ Risco de vazamento de credenciais
   - ❌ Difícil rotação de credenciais

3. **Certificate-based**
   - ✅ Mais seguro
   - ✅ Suporte a rotação
   - ❌ Complexidade de gestão
   - ❌ Setup complicado

## Consequências

### Positivas
1. Flexibilidade no desenvolvimento
2. Segurança em produção
3. Auditoria clara
4. Fácil revogação
5. CI/CD simplificado

### Negativas
1. Manutenção de duas formas de autenticação
2. Configuração adicional em produção
3. Necessidade de documentação clara
4. Possível confusão inicial

## Validação
- Testes em diferentes ambientes
- Verificação de permissões
- Auditoria de acessos
- Testes de revogação

## Referências
- [Kubernetes Authentication](https://kubernetes.io/docs/reference/access-authn-authz/authentication/)
- [RBAC Best Practices](https://kubernetes.io/docs/concepts/security/rbac-good-practices/)
- [ServiceAccount Security](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/)
- [Kubernetes Security Best Practices](https://kubernetes.io/docs/concepts/security/security-checklist/)