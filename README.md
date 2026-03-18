# 🔐 Projeto CRUD Credentials

**API REST para gerenciamento de configurações de variáveis e schemas de compliance** com arquitetura limpa, persistência 100% PostgreSQL e GORM.

> **Status:** ✅ Refatoração Completa | **Versão:** 1.0 | **Data:** Feb 2026

---

## 📖 Visão Geral

Sistema de gerenciamento de credenciais e configurações para frameworks de compliance (PCI-DSS, ISO 27001, etc). Implementa:

- ✅ **2 domínios de negócio**: Benchmark Schemas (compliance) + Variable Configs (core)
- ✅ **Clean Architecture**: Handler → Service → Repository
- ✅ **PostgreSQL + GORM**: Transações ACID, JSONB, índices otimizados
- ✅ **Sem dependências AWS**: Zero DynamoDB, código portável
- ✅ **API RESTful**: CRUD completo com validação de entrada

---

## 🏗️ Arquitetura

### Modelo de Fluxo

```
┌─────────────────────────────────────────────────────┐
│ HTTP Request (Gin)                                  │
└──────────────┬──────────────────────────────────────┘
               ↓
┌─────────────────────────────────────────────────────┐
│ Handler Layer (pkg/handler)                         │
│ • Tradução HTTP → DTO                               │
│ • Validação sintática                               │
│ • Wiring de dependências                            │
└──────────────┬──────────────────────────────────────┘
               ↓
┌─────────────────────────────────────────────────────┐
│ Service Layer (internal/service)                    │
│ • Lógica de negócio                                 │
│ • Validações de domínio                             │
│ • Context propagation                               │
└──────────────┬──────────────────────────────────────┘
               ↓
┌─────────────────────────────────────────────────────┐
│ Repository Layer (internal/storage/postgres)        │
│ • CRUD com GORM                                     │
│ • Transações automáticas                            │
│ • Queries otimizadas                                │
└──────────────┬──────────────────────────────────────┘
               ↓
┌─────────────────────────────────────────────────────┐
│ PostgreSQL Database                                 │
│ • benchmark_schemas (JSONB)                         │
│ • var_configs (JSONB)                               │
└─────────────────────────────────────────────────────┘
```

### Separação de Responsabilidades

| Camada         | Responsabilidade   | **Conhece**         |
| -------------- | ------------------ | ------------------- |
| **Handler**    | Tradução HTTP      | DTOs, status codes  |
| **Service**    | Lógica de negócio  | Domínio, validações |
| **Repository** | Persistência       | GORM, SQL           |
| **Model**      | Estrutura de dados | GORM tags, tipos    |

---

## 🚀 Quick Start

### Pré-requisitos

```bash
# Verificar versões
go version        # 1.23.0+
psql --version    # 12.0+
```

### 1️⃣ Clonar & Instalar

```bash
cd /home/valcann/Documentos/projeto-crud-credentials
go mod download
```

### 2️⃣ Configurar PostgreSQL

```bash
# Criar database
createdb -U postgres projeto-crud-credentials

# Executar migrations
psql -U postgres -d projeto-crud-credentials < migrations/001_create_benchmark_schemas.sql
psql -U postgres -d projeto-crud-credentials < migrations/002_create_varconfig_table.sql

# Verificar
psql -U postgres -d projeto-crud-credentials -c "\dt"
```

### 3️⃣ Configurar Ambiente

```bash
export POSTGRES_CONNECTION_STRING="postgresql://postgres:password@localhost:5432/projeto-crud-credentials?sslmode=disable"
export PORT=8080
```

### 4️⃣ Executar

```bash
go run ./cmd/api/
# Server iniciado em http://localhost:8080
```

---

## 📚 Documentação

---

## 📚 Documentação

Para detalhes específicos, consulte:

- 📖 [REFACTORING_SUMMARY.md](REFACTORING_SUMMARY.md) - Histórico completo da refatoração
- 🚀 [SETUP_POSTGRESQL.md](SETUP_POSTGRESQL.md) - Guia detalhado de configuração PostgreSQL
- 📝 [CHANGES.md](CHANGES.md) - Lista de mudanças e estatísticas
- 🏗️ [ARCHITECTURE.md](ARCHITECTURE.md) - Padrões arquiteturais adotados

---

## 📊 Modelo de Dados

### Domínio 1: Benchmark Schemas

Armazena definições de schemas para compliance frameworks.

```go
type BenchmarkSchema struct {
    ID        int64           `gorm:"primaryKey;autoIncrement"`
    Name      string          `gorm:"not null"`
    Schema    json.RawMessage `gorm:"type:jsonb"`           // Definição do schema
    CreatedAt time.Time       `gorm:"autoCreateTime"`
    UpdatedAt time.Time       `gorm:"autoUpdateTime"`
}
```

**Endpoints:**

```
GET    /compliance/benchmarks              # Listar todos
GET    /compliance/benchmarks/{id}         # Buscar por ID
POST   /compliance/benchmarks              # Criar novo
```

---

### Domínio 2: Variable Configs

Armazena configurações de variáveis ligadas a organizações e benchmarks.

```go
type VarConfig struct {
    ID          int64           `gorm:"primaryKey;autoIncrement"`
    OrgID       string          `gorm:"not null;index:idx_varcfg_org_bench"`
    BenchmarkID string          `gorm:"not null;index:idx_varcfg_org_bench"`
    Payload     json.RawMessage `gorm:"type:jsonb"`         // Configuração JSON
    CreatedAt   time.Time       `gorm:"autoCreateTime"`
    UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
}
```

**Endpoints:**

```
GET    /org/{orgId}/compliance/variables/{benchmarkId}              # Listar
POST   /org/{orgId}/compliance/variables/{benchmarkId}              # Criar
GET    /org/{orgId}/compliance/variables/{benchmarkId}/{id}         # Buscar
PUT    /org/{orgId}/compliance/variables/{benchmarkId}/{id}         # Atualizar
DELETE /org/{orgId}/compliance/variables/{benchmarkId}/{id}         # Deletar
```

---

## 🛠️ Exemplos de Uso

### 1️⃣ Criar um Benchmark Schema

```bash
curl -X POST http://localhost:8080/compliance/benchmarks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PCI-DSS-3.2.1",
    "schema": {
      "card_number": {"type": "string", "required": true},
      "expiry": {"type": "string", "required": true},
      "cvv": {"type": "string", "required": true}
    }
  }'
```

**Response:**

```json
{
  "id": "1",
  "name": "PCI-DSS-3.2.1",
  "created_at": "2026-02-19T10:30:00Z"
}
```

---

### 2️⃣ Listar Schemas

```bash
curl http://localhost:8080/compliance/benchmarks
```

**Response:**

```json
{
  "data": [
    {
      "id": "1",
      "name": "PCI-DSS-3.2.1",
      "schema_body": null,
      "created_at": "2026-02-19T10:30:00Z"
    }
  ],
  "count": 1
}
```

---

### 3️⃣ Buscar Schema Completo (com schema_body)

```bash
curl http://localhost:8080/compliance/benchmarks/1
```

**Response:**

```json
{
  "id": "1",
  "name": "PCI-DSS-3.2.1",
  "schema_body": {
    "card_number": { "type": "string", "required": true },
    "expiry": { "type": "string", "required": true },
    "cvv": { "type": "string", "required": true }
  },
  "created_at": "2026-02-19T10:30:00Z"
}
```

---

### 4️⃣ Criar Configuração de Variável

```bash
curl -X POST http://localhost:8080/org/org-123/compliance/variables/pci-dss-v3.2.1 \
  -H "Content-Type: application/json" \
  -d '{
    "payload": {
      "environment": "production",
      "encryption_enabled": true,
      "rotation_days": 90,
      "allowed_ips": ["192.168.1.0/24", "10.0.0.0/8"]
    }
  }'
```

**Response:**

```json
{
  "id": "1",
  "org_id": "org-123",
  "benchmark_id": "pci-dss-v3.2.1",
  "payload": {
    "environment": "production",
    "encryption_enabled": true,
    "rotation_days": 90,
    "allowed_ips": ["192.168.1.0/24", "10.0.0.0/8"]
  },
  "created_at": "2026-02-19T10:35:00Z",
  "updated_at": "2026-02-19T10:35:00Z"
}
```

---

### 5️⃣ Listar Configurações de uma Organização

```bash
curl http://localhost:8080/org/org-123/compliance/variables/pci-dss-v3.2.1
```

**Response:**

```json
{
  "data": [
    {
      "id": "1",
      "org_id": "org-123",
      "benchmark_id": "pci-dss-v3.2.1",
      "payload": { ... },
      "created_at": "2026-02-19T10:35:00Z",
      "updated_at": "2026-02-19T10:35:00Z"
    }
  ]
}
```

---

### 6️⃣ Atualizar Configuração

```bash
curl -X PUT http://localhost:8080/org/org-123/compliance/variables/pci-dss-v3.2.1/1 \
  -H "Content-Type: application/json" \
  -d '{
    "payload": {
      "environment": "production",
      "encryption_enabled": true,
      "rotation_days": 60
    }
  }'
```

---

## 🗄️ Schema PostgreSQL

### benchmark_schemas

```sql
CREATE TABLE benchmark_schemas (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  schema_body JSONB NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_benchmark_schema_body_gin ON benchmark_schemas USING gin (schema_body);
```

### var_configs

```sql
CREATE TABLE var_configs (
  id BIGSERIAL PRIMARY KEY,
  org_id TEXT NOT NULL,
  benchmark_id TEXT NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_varcfg_org_bench ON var_configs (org_id, benchmark_id);
CREATE INDEX idx_varcfg_payload_gin ON var_configs USING gin (payload);
```

---

## 📁 Estrutura do Projeto

```
projeto-crud-credentials/
│
├── 📖 Documentação
│   ├── README.md                           # Este arquivo
│   ├── ARCHITECTURE.md                     # Padrões de arquitetura
│   ├── REFACTORING_SUMMARY.md             # Histórico da refatoração
│   ├── SETUP_POSTGRESQL.md                # Guia de setup
│   ├── CHANGES.md                         # Lista de mudanças
│   └── DYNAMODB_SCHEMA.md                 # (Legado - referência)
│
├── 🗄️ Migrations
│   ├── 001_create_benchmark_schemas.sql   # Tabela benchmark_schemas
│   └── 002_create_varconfig_table.sql     # Tabela var_configs
│
├── 🏃 cmd/api/
│   └── main.go                            # Bootstrap da aplicação
│
├── 📦 pkg/
│   ├── handler/domain/
│   │   ├── compliance/benchmark/          # Handler: Benchmark Schemas
│   │   └── core/config/                   # Handler: Variable Configs
│   └── models/
│       ├── benchmark/
│       │   └── benchmark_schema.go        # DTO e Model
│       └── config/
│           └── varconfig.go               # DTO e Model
│
├── 🔧 internal/
│   ├── common/
│   │   └── erros.go                       # Tratamento de erros
│   │
│   ├── service/domain/
│   │   ├── compliance/benchmark/          # Service: Benchmark
│   │   │   ├── service.go
│   │   │   ├── ports.go
│   │   │   └── helper.go
│   │   └── core/config/                   # Service: Config
│   │       ├── service.go
│   │       └── ports.go
│   │
│   └── storage/postgres/
│       ├── benchmark/                     # Repository: Benchmark
│       │   ├── repository.go
│       │   └── assertions.go
│       └── config/                        # Repository: Config
│           └── repository.go
│
├── 🛣️ routes/
│   ├── router.go                          # Setup de rotas
│   ├── benchmark_schema.go                # Rotas: Benchmark
│   └── varconfig.go                       # Rotas: Config
│
├── 📋 dto/
│   ├── benchmark/
│   │   └── benchmark_schemas.go           # DTOs: Benchmark
│   └── config/
│       └── varconfig.go                   # DTOs: Config
│
├── go.mod                                 # Dependências
├── go.sum                                 # Lock de dependências

```

---

## 🧪 Testes

### Executar Todos os Testes

```bash
go test ./...
```

### Testes por Camada

```bash
# Testes de serviço (lógica de negócio)
go test ./internal/service/domain/...

# Testes de handler (validação HTTP)
go test ./pkg/handler/domain/...

# Testes de repository (persistência)
go test ./internal/storage/postgres/...
```

### Testes com Cobertura

```bash
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 🔧 Variáveis de Ambiente

| Variável                     | Padrão | Descrição                                  |
| ---------------------------- | ------ | ------------------------------------------ |
| `POSTGRES_CONNECTION_STRING` | -      | Connection string PostgreSQL (obrigatório) |
| `PORT`                       | `8080` | Porta da aplicação                         |
| `LOG_LEVEL`                  | `info` | Nível de log (debug, info, warn, error)    |

**Exemplo de .env:**

```bash
POSTGRES_CONNECTION_STRING="postgresql://postgres:password@localhost:5432/projeto-crud-credentials?sslmode=disable"
PORT=8080
LOG_LEVEL=debug
```

---

## 🐛 Troubleshooting

### Erro: Connection Refused

```bash
# Verificar se PostgreSQL está rodando
psql -U postgres -d postgres -c "SELECT 1"

# Testar connection string
psql "postgresql://postgres:password@localhost:5432/projeto-crud-credentials?sslmode=disable"
```

### Erro: Tabelas não Encontradas

```bash
# Executar migrations novamente
psql -U postgres -d projeto-crud-credentials < migrations/001_create_benchmark_schemas.sql
psql -U postgres -d projeto-crud-credentials < migrations/002_create_varconfig_table.sql

# Verificar tabelas
psql -U postgres -d projeto-crud-credentials -c "\dt"
```

### Erro: Permission Denied (GORM)

```bash
# Verificar permissões do usuário PostgreSQL
psql -U postgres -c "\du"

# Se necessário, resetar senha
ALTER USER postgres WITH PASSWORD 'new_password';
```

---

## 📈 Monitoramento

### Logs em Tempo Real

```bash
# Ver apenas erros
go run ./cmd/api/ 2>&1 | grep ERROR

# Ver com timestamp
go run ./cmd/api/ 2>&1 | while read line; do echo "[$(date '+%H:%M:%S')] $line"; done
```

### Verificar Saúde da Aplicação

```bash
# Health check (se implementado)
curl http://localhost:8080/health

# Verificar tabelas
psql -U postgres -d projeto-crud-credentials -c "SELECT COUNT(*) FROM benchmark_schemas"
```

---

## 🚀 Deploy

### Docker

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main ./cmd/api/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

**Build e Run:**

```bash
docker build -t projeto-crud-credentials .
docker run -e POSTGRES_CONNECTION_STRING="..." -p 8080:8080 projeto-crud-credentials
```

### Docker Compose

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: projeto-crud-credentials
    ports:
      - '5432:5432'
    volumes:
      - ./migrations:/docker-entrypoint-initdb.d

  api:
    build: .
    environment:
      POSTGRES_CONNECTION_STRING: 'postgresql://postgres:password@postgres:5432/projeto-crud-credentials?sslmode=disable'
      PORT: 8080
    ports:
      - '8080:8080'
    depends_on:
      - postgres
```

**Executar:**

```bash
docker-compose up
```

---

## 🔐 Segurança

- ✅ **Validação de entrada**: Todos os handlers validam requests
- ✅ **ACID Transactions**: GORM gerencia transações automaticamente
- ✅ **SQL Injection Prevention**: Prepared statements via GORM
- ✅ **Context Propagation**: Cancellation e timeouts propagados
- ✅ **Error Handling**: Erros sensíveis não vazados para cliente

---

## 📊 Dependências

| Pacote                    | Versão  | Propósito          |
| ------------------------- | ------- | ------------------ |
| `gin-gonic/gin`           | v1.11.0 | HTTP Framework     |
| `gorm.io/gorm`            | v1.31.1 | ORM                |
| `gorm.io/driver/postgres` | v1.6.0  | PostgreSQL Driver  |
| `lib/pq`                  | v1.11.1 | PostgreSQL Adapter |

**Nota:** Todas as dependências AWS (DynamoDB SDK) foram removidas na refatoração.

---

## 📝 Licença

MIT License - Veja LICENSE file para detalhes.

---

## 👥 Contribuindo

Para contribuir, crie uma branch e faça um pull request:

```bash
git checkout -b feature/nova-feature
git commit -am 'Adiciona nova feature'
git push origin feature/nova-feature
```

---

## 📞 Suporte

Para dúvidas ou issues:

1. Consulte a documentação em [docs/](docs/)
2. Verifique [REFACTORING_SUMMARY.md](REFACTORING_SUMMARY.md)
3. Abra uma issue no repositório

---

**Última Atualização:** 19 de fevereiro de 2026 | **Status:** ✅ Pronto para Produção
