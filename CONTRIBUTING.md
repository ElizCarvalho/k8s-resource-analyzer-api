# Contribuindo com o K8s Resource Analyzer

## Índice
1. [Como Contribuir](#como-contribuir)
2. [Processo de Desenvolvimento](#processo-de-desenvolvimento)
3. [Processo de Release](#processo-de-release)
4. [Padrões de Código](#padrões-de-código)

## Como Contribuir

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nome-da-feature`)
3. Commit suas mudanças (`git commit -m 'tipo: descrição'`)
4. Push para a branch (`git push origin feature/nome-da-feature`)
5. Abra um Pull Request

## Processo de Desenvolvimento

### Branches
- `main`: Código em produção
- `feature/*`: Novas funcionalidades
- `bugfix/*`: Correções de bugs
- `hotfix/*`: Correções urgentes em produção

### Commits
Usamos Conventional Commits:
- `feat`: Nova funcionalidade
- `fix`: Correção de bug
- `docs`: Documentação
- `style`: Formatação
- `refactor`: Refatoração
- `test`: Testes
- `chore`: Manutenção

## Processo de Release

### Pré-requisitos
- Ser mantenedor do projeto
- Ter acesso aos secrets do repositório
- Ter permissão de escrita no Docker Hub

### Checklist de Release
1. Garantir que todos os testes passam
2. Verificar se a documentação está atualizada
3. Revisar o CHANGELOG.md
4. Verificar se todas as dependências estão atualizadas

### Criando uma Release

1. Atualizar a main:
   ```bash
   git checkout main
   git pull origin main
   ```

2. Criar e publicar tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. Monitorar o workflow de release no GitHub Actions:
   - Build do binário
   - Publicação da imagem Docker
   - Criação da release no GitHub

### Versionamento
MAJOR.MINOR.PATCH:
- MAJOR: Mudanças incompatíveis
- MINOR: Novas funcionalidades
- PATCH: Correções de bugs

### Pós-Release
1. Verificar se a imagem Docker foi publicada
2. Validar a documentação da release no GitHub
3. Comunicar a equipe sobre a nova versão

## Padrões de Código

### Go
- Use `gofmt` para formatação
- Siga as convenções do Go
- Mantenha 80% de cobertura de testes
- Documente funções públicas

### Documentação
- Mantenha o README.md atualizado
- Documente mudanças na API
- Atualize o Swagger quando necessário

### Qualidade
- Todos os testes devem passar
- Não deve haver warnings do linter
- Mantenha a complexidade ciclomática baixa 