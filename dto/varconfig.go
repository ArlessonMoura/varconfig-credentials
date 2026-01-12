package varconfig

// CreateVarConfigRequest DTO para criar um novo VarConfig
type CreateVarConfigRequest struct {
Payload map[string]any `json:"payload" binding:"required"`
}

// UpdateVarConfigRequest DTO para atualizar um VarConfig
type UpdateVarConfigRequest struct {
Payload map[string]any `json:"payload" binding:"required"`
}

// VarConfigResponse DTO para responder com um VarConfig
type VarConfigResponse struct {
ID          int64          `json:"id"`
OrgID       int64          `json:"org_id"`
BenchmarkID string         `json:"benchmark_id"`
Payload     map[string]any `json:"payload"`
CreatedAt   string         `json:"created_at"`
UpdatedAt   string         `json:"updated_at"`
}

// ListVarConfigResponse DTO para listar VarConfigs
type ListVarConfigResponse struct {
Data []VarConfigResponse `json:"data"`
}
