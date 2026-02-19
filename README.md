# VarConfig CRUD API

API REST para gerenciamento de configurações de variáveis com armazenamento em **AWS DynamoDB**.

## 🏗️ Arquitetura Limpa

Implementação seguindo **Clean Architecture** com separação rigorosa de responsabilidades:

```
Handler → Service → Storage
   ↓         ↓         ↓
HTTP      Negócio   DynamoDB
```

### 🎯 Camadas Implementadas

#### **🚀 Common Infrastructure** (`internal/common/`)

- **errors/**: Tratamento centralizado de erros estruturados
- **logger/**: Logging com context propagation e níveis
- **metrics/**: Métricas com padrão `service_<domain>_<operation>_latency_ms`

#### **🌐 Handler Layer** (`pkg/handler/varconfig/`)

- Tradução de protocolo HTTP
- Validação sintática de entrada
- Wiring de dependências
- **SEM lógica de negócio**

#### **🎯 Service Layer** (`internal/service/domain/varconfig/`)

- Lógica de negócio completa
- Validações de regras
- Injeção de repository, logger, metrics
- Context propagation

#### **💾 Storage Layer** (`internal/storage/dynamodb/varconfig/`)

- Implementação DynamoDB com PK+SK otimizado
- **ZERO Scan operations** - apenas Query eficientes
- Logging e métricas integradas
- **SEM conhecimento de handlers/DTOs**

## 📋 Pré-requisitos

- Go 1.23+
- AWS CLI configurado
- Tabela DynamoDB `var_configs` criada

## 🚀 Configuração

### 1. Criar Tabela DynamoDB

Use o script automatizado:

```bash
./setup-dynamodb.sh
```

Ou manualmente:

```bash
aws dynamodb create-table \
  --table-name var_configs \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

### 2. Configurar Credenciais AWS

```bash
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-east-1
```

### 3. Instalar Dependências

```bash
go mod download
```

## 🏃‍♂️ Execução

```bash
go run .
```

Servidor iniciará em `http://localhost:8080`

## 📊 Modelo de Dados

### VarConfig (Domínio)

```go
type VarConfig struct {
    ID          int64                  `json:"id"`
    OrgID       int64                  `json:"org_id"`
    BenchmarkID string                  `json:"benchmark_id"`
    Payload     map[string]interface{}   `json:"payload"`
    CreatedAt   time.Time              `json:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at"`
}
```

### Estrutura DynamoDB (PK+SK)

```
PK  = "ORG#{orgId}#BENCH#{benchmark_id}"
SK  = "VARCONFIG#{id}"
```

**Exemplo Real:**

```
PK  = "ORG#123#BENCH#pci-dss-v3.2.1"
SK  = "VARCONFIG#1640995200000000000"
```

## 🛤️ API Endpoints

### Coleção

```http
GET    /org/{orgId}/compliance/variables/{benchmark_id}
POST   /org/{orgId}/compliance/variables/{benchmark_id}
```

### Específico

```http
GET    /org/{orgId}/compliance/variables/{benchmark_id}/{id}
PUT    /org/{orgId}/compliance/variables/{benchmark_id}/{id}
PATCH  /org/{orgId}/compliance/variables/{benchmark_id}/{id}
DELETE /org/{orgId}/compliance/variables/{benchmark_id}/{id}
```

## 📝 Exemplos de Uso

### Criar VarConfig

```bash
curl -X POST http://localhost:8080/org/123/compliance/variables/pci-dss-v3.2.1 \
  -H "Content-Type: application/json" \
  -d '{
    "payload": {
      "max_memory": 1024,
      "allowed_types": ["web", "api"],
      "timeout": 30
    }
  }'
```

### Listar por Benchmark

```bash
curl http://localhost:8080/org/123/compliance/variables/pci-dss-v3.2.1
```

### Buscar por ID

```bash
curl http://localhost:8080/org/123/compliance/variables/pci-dss-v3.2.1/1640995200000000000
```

### Atualizar

```bash
curl -X PUT http://localhost:8080/org/123/compliance/variables/pci-dss-v3.2.1/1640995200000000000 \
  -H "Content-Type: application/json" \
  -d '{
    "payload": {
      "max_memory": 2048,
      "allowed_types": ["web", "api", "mobile"],
      "timeout": 60
    }
  }'
```

### Deletar

```bash
curl -X DELETE http://localhost:8080/org/123/compliance/variables/pci-dss-v3.2.1/1640995200000000000
```

## 🧪 Testes

### Unit Tests

```bash
go test ./internal/service/domain/varconfig/...
```

### Integration Tests

```bash
go test ./pkg/handler/varconfig/...
```

### Contract Tests

```bash
go test ./internal/service/domain/varconfig/... -run Contract
```

## 🔧 Estrutura do Projeto

```
projeto-crud-credentials/
├── 📚 DOCUMENTAÇÃO
│   ├── README.md                    # Este arquivo
│   ├── DYNAMODB_SCHEMA.md           # Schema detalhado
│   └── STRUCTURE.txt                # Arquitetura completa
├── 🚀 internal/common/              # Infraestrutura compartilhada
│   ├── errors/                    # Tratamento de erros
│   ├── logger/                    # Logging centralizado
│   └── metrics/                   # Métricas
├── 🌐 pkg/handler/varconfig/       # Camada HTTP
├── 🎯 internal/service/domain/      # Camada de negócio
└── 💾 internal/storage/dynamodb/     # Camada de dados
```

## 🐛 Debug & Troubleshooting

### Verificar Schema

```bash
./check-schema.sh
```

### Logs Estruturados

A aplicação usa logging estruturado com diferentes níveis:

```bash
# Ver logs em tempo real
go run . 2>&1 | grep -E "(ERROR|WARN)"

# Logs detalhados com contexto
export LOG_LEVEL=debug
go run .
```

### Métricas

As métricas são registradas automaticamente. Para implementação real:

```go
// Implementar Collector real em produção
type PrometheusCollector struct {
    // ... implementação Prometheus
}
```

## 🚀 Deploy

### Docker

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### AWS ECS/Fargate

1. Build da imagem Docker
2. Push para ECR
3. Configurar task definition com:
   - Environment variables para AWS credentials
   - IAM Role com acesso DynamoDB
   - Port mapping 8080

### Variáveis de Ambiente

```bash
# Configuração
AWS_REGION=us-east-1
DYNAMODB_TABLE=varconfigs
LOG_LEVEL=info
PORT=8080
```

## 🔐 Segurança

- **IAM Least Privilege**: Role com acesso apenas à tabela `varconfigs`
- **VPC Endpoints**: Acesso DynamoDB via VPC endpoints
- **Encryption**: At-rest encryption habilitado
- **TLS**: API servida com HTTPS em produção

## 📈 Monitoramento

### Health Check

```bash
curl http://localhost:8080/health
```
