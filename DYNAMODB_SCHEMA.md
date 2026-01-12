# DynamoDB Table Schema

## Table: VarConfigs

### Attributes

- **PK** (String, Partition Key) - "ORG#{orgId}#BENCH#{benchmark_id}"
- **SK** (String, Sort Key) - "VARCONFIG#{id}"
- **id** (String) - Unique identifier for the VarConfig
- **org_id** (String) - Organization identifier
- **benchmark_id** (String) - Benchmark identifier
- **payload** (String) - JSON string containing configuration variables
- **created_at** (String) - ISO 8601 timestamp of creation
- **updated_at** (String) - ISO 8601 timestamp of last update

### Key Structure

```
PK = "ORG#{orgId}#BENCH#{benchmark_id}"
SK = "VARCONFIG#{id}"
```

### Example Item

```json
{
  "PK": "ORG#123#BENCH#pci-dss-v3.2.1",
  "SK": "VARCONFIG#1640995200000000000",
  "id": "1640995200000000000",
  "org_id": "123",
  "benchmark_id": "pci-dss-v3.2.1",
  "payload": "{\"max_login_attempts\": 3, \"session_timeout\": 1800}",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}

### Create Table
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

### Query Examples

#### List all VarConfigs for a specific org and benchmark

```bash
aws dynamodb query \
  --table-name VarConfigs \
  --key-condition-expression "PK = :pk" \
  --expression-attribute-values '{":pk":{"S":"ORG#123#BENCH#pci-dss-v3.2.1"}}' \
  --region us-east-1
```

#### Get specific VarConfig

```bash
aws dynamodb get-item \
  --table-name VarConfigs \
  --key '{"PK":{"S":"ORG#123#BENCH#pci-dss-v3.2.1"},"SK":{"S":"VARCONFIG#1640995200000000000"}}' \
  --region us-east-1
```

## Access Patterns

### 1. Get VarConfig by ID

- **Operation**: GetItem
- **Key**: PK + SK
- **Performance**: O(1)

### 2. List VarConfigs by Organization + Benchmark

- **Operation**: Query
- **Key**: PK only
- **Performance**: O(log N) where N = items per org/benchmark

### 3. Create/Update/Delete VarConfig

- **Operation**: PutItem/UpdateItem/DeleteItem
- **Key**: PK + SK
- **Performance**: O(1)

## Benefits of This Design

✅ **No Scan Operations** - All access patterns use optimized GetItem/Query
✅ **Efficient Partitioning** - Data distributed by org+benchmark
✅ **Strong Consistency** - Direct key access guarantees consistency
✅ **Cost Effective** - Minimal read/write capacity consumption
✅ **Scalable** - Handles high throughput per organization/benchmark
