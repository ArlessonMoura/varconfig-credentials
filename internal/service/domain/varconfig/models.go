package varconfig

import "time"

// VarConfig representa um conjunto de configurações de variáveis para um benchmark
// Para DynamoDB, o ID será tratado como string internamente mas mantido como int64 na interface
type VarConfig struct {
	ID          int64
	OrgID       int64
	BenchmarkID string
	Payload     map[string]any // Conjunto flexível de variáveis configuradas
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
