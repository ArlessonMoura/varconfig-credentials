package varconfig

type CreateVarConfigRequest struct {
	Payload map[string]any `json:"payload" binding:"required"`
}

type UpdateVarConfigRequest struct {
	Payload map[string]any `json:"payload" binding:"required"`
}

type VarConfigResponse struct {
	ID          string         `json:"id"`           // Alterado para string (SK)
	OrgID       string         `json:"org_id"`       // Alterado para string (PK)
	BenchmarkID string         `json:"benchmark_id"` // Parte da PK
	Payload     map[string]any `json:"payload"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type ListVarConfigResponse struct {
	Data []VarConfigResponse `json:"data"`
}
