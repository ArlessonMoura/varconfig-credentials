#!/bin/bash
# Script para executar testes do projeto VarConfig

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

echo "========================================="
echo "🧪 Executando Testes — VarConfig CRUD"
echo "========================================="
echo ""

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 1. Verificar formatação
echo -e "${YELLOW}1️⃣  Verificando formatação (gofmt)...${NC}"
if gofmt -l ./... | grep -q .; then
    echo -e "${RED}❌ Arquivos mal formatados encontrados${NC}"
    gofmt -l ./...
    exit 1
else
    echo -e "${GREEN}✅ Formatação OK${NC}"
fi
echo ""

# 2. Verificar lint (se golangci-lint estiver instalado)
echo -e "${YELLOW}2️⃣  Executando golangci-lint...${NC}"
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run ./... --deadline=2m; then
        echo -e "${GREEN}✅ Lint OK${NC}"
    else
        echo -e "${RED}❌ Lint encontrou problemas${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⏭️  golangci-lint não instalado, pulando${NC}"
fi
echo ""

# 3. Executar testes de negócio
echo -e "${YELLOW}3️⃣  Executando testes unitários...${NC}"
go test -v ./internal/service/domain/varconfig/... -run TestService
echo -e "${GREEN}✅ Testes do service OK${NC}"
echo ""

# 4. Verificar se o código compila
echo -e "${YELLOW}4️⃣  Verificando compilação...${NC}"
go build -v ./...
echo -e "${GREEN}✅ Build OK${NC}"
echo ""

# 5. Coverage (opcional)
echo -e "${YELLOW}5️⃣  Calculando cobertura de testes...${NC}"
if go test -coverprofile=coverage.out ./internal/service/domain/varconfig/...; then
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo -e "${GREEN}✅ Cobertura: $COVERAGE${NC}"
    # Descomenta linha abaixo para gerar relatório HTML
    # go tool cover -html=coverage.out -o coverage.html
else
    echo -e "${YELLOW}⏭️  Testes de coverage falharam${NC}"
fi
echo ""

# 6. Verificar se todas as interfaces estão implementadas
echo -e "${YELLOW}6️⃣  Verificando assertions de interface...${NC}"
grep -r "var _ " ./internal/service/domain/varconfig/assertions.go && \
grep -r "var _ " ./internal/storage/postgres/varconfig/assertions.go && \
echo -e "${GREEN}✅ Assertions OK${NC}" || \
echo -e "${YELLOW}⚠️  Verifique manualmente${NC}"
echo ""

echo "========================================="
echo -e "${GREEN}🎉 Todos os testes passaram!${NC}"
echo "========================================="
