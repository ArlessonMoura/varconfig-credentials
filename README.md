# VarConfig CRUD API

API REST para gerenciamento de configurações de variáveis com armazenamento em **AWS DynamoDB**.

## 🏗️ Arquitetura

- **Go 1.23** com **Gin** para HTTP
- **AWS SDK v2** para DynamoDB
- **Chave composta PK+SK** para performance otimizada
- **Zero Scan operations** - apenas GetItem/Query
- **Payload JSON flexível** sem validação de conteúdo

## 📋 Pré-requisitos

- Go 1.23+
- AWS CLI configurado
- Tabela DynamoDB `VarConfigs` criada

## 🚀 Configuração

### 1. Criar Tabela DynamoDB

```bash
aws dynamodb create-table \
  --table-name VarConfigs \
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
# ou compilado
go build -o varconfig-api && ./varconfig-api
```

Servidor iniciado em `http://localhost:8080`

## 📡 Endpoints da API

### Coleção

- **GET** `/org/{orgId}/compliance/variables/{benchmark_id}` — lista todos os VarConfigs do benchmark
- **POST** `/org/{orgId}/compliance/variables/{benchmark_id}` — cria um VarConfig

### Específico

- **GET** `/org/{orgId}/compliance/variables/{benchmark_id}/{id}` — obtém um VarConfig
- **PUT** `/org/{orgId}/compliance/variables/{benchmark_id}/{id}` — atualiza (completo)
- **PATCH** `/org/{orgId}/compliance/variables/{benchmark_id}/{id}` — atualiza (parcial)
- **DELETE** `/org/{orgId}/compliance/variables/{benchmark_id}/{id}` — remove

## 💡 Exemplos de Uso

### Criar VarConfig

```bash
curl -X POST "http://localhost:8080/org/123/compliance/variables/pci-dss-v3.2.1" \
  -H "Content-Type: application/json" \
  -d '{
    "payload": {
      "max_login_attempts": 3,
      "session_timeout": 1800,
      "allowed_types": ["admin", "user"]
    }
  }'
```

### Listar VarConfigs

```bash
curl -X GET "http://localhost:8080/org/123/compliance/variables/pci-dss-v3.2.1"
```

### Obter VarConfig Específico

```bash
curl -X GET "http://localhost:8080/org/123/compliance/variables/pci-dss-v3.2.1/1640995200000000000"
```

### Atualizar VarConfig

```bash
curl -X PUT "http://localhost:8080/org/123/compliance/variables/pci-dss-v3.2.1/1640995200000000000" \
  -H "Content-Type: application/json" \
  -d '{
    "payload": {
      "max_login_attempts": 5,
      "session_timeout": 3600
    }
  }'
```

### Deletar VarConfig

```bash
curl -X DELETE "http://localhost:8080/org/123/compliance/variables/pci-dss-v3.2.1/1640995200000000000"
```

## 📊 Modelo de Dados

### Domain Model

```go
type VarConfig struct {
    ID          int64                  `json:"id"`
    OrgID       int64                  `json:"org_id"`
    BenchmarkID string                 `json:"benchmark_id"`
    Payload     map[string]interface{} `json:"payload"`
    CreatedAt   time.Time              `json:"created_at"`
    UpdatedAt   time.Time              `json:"updated_at"`
}
```

### DynamoDB Schema

```json
{
  "PK": "ORG#123#BENCH#pci-dss-v3.2.1",
  "SK": "VARCONFIG#1640995200000000000",
  "id": "1640995200000000000",
  "org_id": "123",
  "benchmark_id": "pci-dss-v3.2.1",
  "payload": "{\"max_login_attempts\": 3}",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

## ⚡ Performance

### Operações Otimizadas

| Operação | Método DynamoDB | Performance |
|----------|-----------------|-------------|
| Create   | PutItem         | O(1) |
| Read     | GetItem         | O(1) |
| List     | Query           | O(log N) |
| Update   | UpdateItem      | O(1) |
| Delete   | DeleteItem      | O(1) |

### Estrutura de Chaves

```
PK  = "ORG#{orgId}#BENCH#{benchmark_id}"
SK  = "VARCONFIG#{id}"
```

**Benefícios:**

- ✅ **Zero Scan** - apenas operações otimizadas
- ✅ **Particionamento** por organização
- ✅ **Escalabilidade** linear
- ✅ **Custo otimizado** com pay-per-request

## 🧪 Testes

### Testes Unitários

```bash
go test ./...
```

### Testes de Integração

```bash
# Script completo de testes
chmod +x test.sh
./test.sh
```

## 📁 Estrutura do Projeto

```
projeto-crud-credencials/
├── 📚 docs/
│   ├── README.md
│   ├── DYNAMODB_SCHEMA.md
│   └── STRUCTURE.txt
├── 🎯 pkg/handler/varconfig/     # HTTP Layer
├── 🔧 internal/service/domain/   # Business Logic
├── 💾 internal/storage/dynamodb/ # DynamoDB Layer
└── 🚀 main.go                    # Application Entry
```

[Ver estrutura completa →](STRUCTURE.txt)

## 🔧 Configuração Avançada

### Variáveis de Ambiente

```bash
export TABLE_NAME=VarConfigs
export PORT=8080
export LOG_LEVEL=info
```

### Configuração AWS

```go
// Custom endpoint para desenvolvimento
cfg, err := config.LoadDefaultConfig(context.TODO(),
    config.WithRegion("us-east-1"),
    config.WithEndpointResolver(aws.EndpointResolverWithOptionsFunc(
        func(service, region string, options ...interface{}) (aws.Endpoint, error) {
            return aws.Endpoint{URL: "http://localhost:8000"}, nil
        },
    )),
)
```

## 🛡️ Segurança

- **Validação de entrada** nos handlers
- **Tratamento de erros** sem expor detalhes internos
- **CORS** configurável para produção
- **Rate limiting** recomendado via API Gateway

## 📈 Monitoramento

### Logs Estruturados

```go
fmt.Printf("[%s] %s %s - %v\n", 
    time.Now().Format("2006-01-02T15:04:05Z"),
    c.Request.Method,
    c.Request.URL.Path,
    err,
)
```

### Métricas Sugeridas

- Latência por operação
- Taxa de erro por endpoint
- Consumo de capacidade DynamoDB
- Contagem de requisições por organização

## 🚀 Deploy

### Docker

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o varconfig-api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/varconfig-api .
EXPOSE 8080
CMD ["./varconfig-api"]
```

### AWS ECS/Fargate

1. Build da imagem Docker
2. Push para ECR
3. Configurar task definition com IAM role para DynamoDB
4. Deploy via ECS service

## 📝 Licença

MIT License - ver arquivo [LICENSE](LICENSE)

## 🤝 Contribuição

1. Fork do projeto
2. Feature branch (`git checkout -b feature/amazing-feature`)
3. Commit (`git commit -m 'Add amazing feature'`)
4. Push (`git push origin feature/amazing-feature`)
5. Pull Request

---

## 📞 Suporte

- 📧 Email: <support@projeto-crud-credencials.com>
- 📖 Docs: [DYNAMODB_SCHEMA.md](DYNAMODB_SCHEMA.md)
- 🏗️ Arquitetura: [STRUCTURE.txt](STRUCTURE.txt)
