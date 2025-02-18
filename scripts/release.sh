#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Variáveis
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
CHANGELOG_FILE="CHANGELOG.md"
RELEASE_NOTES_FILE="release_notes.md"

# Função para validar versão semântica
validate_version() {
    local version=$1
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo -e "${RED}❌ Formato de versão inválido. Use o formato vX.Y.Z (ex: v1.0.0)${NC}"
        exit 1
    fi
}

# Função para verificar se há mudanças não commitadas
check_git_status() {
    if [[ -n $(git status -s) ]]; then
        echo -e "${RED}❌ Existem mudanças não commitadas. Faça commit ou stash antes de criar uma release.${NC}"
        exit 1
    fi
}

# Função para gerar changelog
generate_changelog() {
    local new_version=$1
    local last_tag=$CURRENT_VERSION
    local date=$(date '+%Y-%m-%d')
    
    echo -e "${BLUE}📝 Gerando changelog...${NC}"
    
    # Cria arquivo de release notes temporário
    echo "# Release $new_version ($date)" > $RELEASE_NOTES_FILE
    echo "" >> $RELEASE_NOTES_FILE
    echo "## Mudanças" >> $RELEASE_NOTES_FILE
    echo "" >> $RELEASE_NOTES_FILE
    
    # Obtém commits desde a última tag
    git log --pretty=format:"* %s" $last_tag..HEAD >> $RELEASE_NOTES_FILE
    
    # Abre o arquivo para edição
    echo -e "${YELLOW}📝 O arquivo de release notes foi gerado.${NC}"
    echo -e "${YELLOW}Por favor, revise e edite o arquivo $RELEASE_NOTES_FILE${NC}"
    echo -e "${YELLOW}Pressione ENTER quando terminar...${NC}"
    read -r
    
    # Atualiza CHANGELOG.md
    if [ ! -f $CHANGELOG_FILE ]; then
        echo "# Changelog" > $CHANGELOG_FILE
        echo "" >> $CHANGELOG_FILE
    fi
    
    # Insere as novas release notes no início do changelog
    cat $RELEASE_NOTES_FILE | cat - $CHANGELOG_FILE > temp && mv temp $CHANGELOG_FILE
}

# Função para atualizar versão
update_version() {
    local new_version=$1
    
    # Atualiza versão no código
    echo -e "${BLUE}📦 Atualizando versão nos arquivos...${NC}"
    
    # Aqui você pode adicionar comandos para atualizar a versão em outros arquivos do projeto
    # Por exemplo:
    # sed -i "s/version = \".*\"/version = \"$new_version\"/" config.go
}

# Função para criar release
create_release() {
    local new_version=$1
    
    echo -e "${BLUE}🚀 Criando release $new_version...${NC}"
    
    # Commit das mudanças
    git add $CHANGELOG_FILE
    git commit -m "chore: atualiza changelog para versão $new_version"
    
    # Cria tag
    git tag -a $new_version -m "Release $new_version"
    
    echo -e "${GREEN}✅ Release $new_version criada com sucesso!${NC}"
    echo -e "\n${YELLOW}Próximos passos:${NC}"
    echo "1. Execute 'git push origin main' para enviar as mudanças"
    echo "2. Execute 'git push --tags' para enviar a tag"
    echo "3. Execute 'make docker-build' para criar a imagem Docker"
    echo "4. Execute 'make docker-push' para publicar a imagem"
}

# Menu principal
echo -e "${BLUE}🚀 Assistente de Release${NC}"
echo -e "Versão atual: ${YELLOW}$CURRENT_VERSION${NC}\n"

# Verifica status do git
check_git_status

# Solicita nova versão
echo -e "Digite a nova versão (formato vX.Y.Z):"
read -r new_version

# Valida formato da versão
validate_version $new_version

# Confirma ação
echo -e "\n${YELLOW}Você está prestes a criar a release $new_version${NC}"
echo -e "Isso irá:"
echo "1. Gerar um changelog com as mudanças desde $CURRENT_VERSION"
echo "2. Atualizar a versão nos arquivos do projeto"
echo "3. Criar uma tag git"
echo -e "\nContinuar? [y/N]"
read -r response

if [[ "$response" =~ ^[Yy]$ ]]; then
    generate_changelog $new_version
    update_version $new_version
    create_release $new_version
    
    # Limpa arquivo temporário
    rm -f $RELEASE_NOTES_FILE
else
    echo -e "${YELLOW}❌ Operação cancelada${NC}"
    exit 0
fi 