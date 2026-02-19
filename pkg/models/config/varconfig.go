package models

import "time"

// VarConfig é a entidade de domínio e resposta da API
type VarConfig struct {
	ID          string         `json:"id"`
	OrgID       string         `json:"org_id"`
	BenchmarkID string         `json:"benchmark_id"`
	Payload     map[string]any `json:"payload"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// VarConfigCreateRequest representa os dados para criação de uma nova configuração
type VarConfigCreateRequest struct {
	OrgID       string         `json:"org_id"`
	BenchmarkID string         `json:"benchmark_id"`
	Payload     map[string]any `json:"payload"`
}

// VarConfigUpdateRequest representa os dados para atualização de uma configuração existente
type VarConfigUpdateRequest struct {
	ID          string         `json:"id"`
	OrgID       string         `json:"org_id"`
	BenchmarkID string         `json:"benchmark_id"`
	Payload     map[string]any `json:"payload"`
}

// VarConfigItem é como o dado é salvo no DynamoDB
type VarConfigItem struct {
	PK string `dynamodbav:"PK"` // ORG#{orgId}#BENCH#{benchmark_id}
	SK string `dynamodbav:"SK"` // VARCONFIG#{id}

	//Voltei a Repetir os IDs como atributos para facilitar a leitura, filtros e possiveis indices com a provavel volta do projeto full relacional
	ID          string         `dynamodbav:"id"`
	OrgID       string         `dynamodbav:"org_id"`
	BenchmarkID string         `dynamodbav:"benchmark_id"`
	Payload     map[string]any `dynamodbav:"payload"`
	CreatedAt   string         `dynamodbav:"created_at"`
	UpdatedAt   string         `dynamodbav:"updated_at"`
}

// VarConfigListItem representa os campos essenciais para listagem (sem o payload pesado)
type VarConfigListItem struct {
	ID          string `dynamodbav:"id"`
	OrgID       string `dynamodbav:"org_id"`
	BenchmarkID string `dynamodbav:"benchmark_id"`
	CreatedAt   string `dynamodbav:"created_at"`
	UpdatedAt   string `dynamodbav:"updated_at"`
}
