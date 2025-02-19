#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}🔍 Verificando documentação...${NC}"

# Verifica se README.md e README.en.md existem
if [ ! -f "README.md" ] || [ ! -f "README.en.md" ]; then
    echo -e "${RED}❌ README.md ou README.en.md não encontrado${NC}"
    exit 1
fi

# Verifica se os READMEs têm tamanhos similares (tolerância de 20%)
pt_size=$(wc -l < README.md)
en_size=$(wc -l < README.en.md)
diff=$((pt_size - en_size))
if [ ${diff#-} -gt $((pt_size * 20 / 100)) ]; then
    echo -e "${RED}❌ Diferença significativa no tamanho dos READMEs${NC}"
    echo "README.md: $pt_size linhas"
    echo "README.en.md: $en_size linhas"
    exit 1
fi

# Verifica links quebrados
echo -e "\n${YELLOW}🔍 Verificando links...${NC}"
for file in README.md README.en.md docs/*.md; do
    if [ -f "$file" ]; then
        echo -e "Verificando links em $file..."
        # Extrai URLs e verifica cada uma
        grep -o '\[.*\]([^)]*)' "$file" | grep -o '([^)]*)' | tr -d '()' | while read -r url; do
            if [[ $url == http* ]]; then
                if ! curl --output /dev/null --silent --head --fail "$url"; then
                    echo -e "${RED}❌ Link quebrado em $file: $url${NC}"
                    exit 1
                fi
            fi
        done
    fi
done

# Valida Swagger
if [ -f "docs/swagger.json" ]; then
    echo -e "\n${YELLOW}🔍 Validando Swagger...${NC}"
    if command -v swagger-cli &> /dev/null; then
        swagger-cli validate docs/swagger.json || {
            echo -e "${RED}❌ Erro na validação do Swagger${NC}"
            exit 1
        }
    else
        echo -e "${YELLOW}⚠️ swagger-cli não encontrado. Instale com: npm install -g swagger-cli${NC}"
    fi
fi

echo -e "\n${GREEN}✅ Documentação verificada com sucesso!${NC}" 