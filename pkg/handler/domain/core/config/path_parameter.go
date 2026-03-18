package config

import (
	"projeto-crud-credentials/internal/common/validation"
)

type PathParams struct {
	OrgID       string `uri:"orgId" binding:"required"`
	BenchmarkID string `uri:"benchmark_id" binding:"required"`
	ID          string `uri:"id"`
}

// Validate implementa a interface Validatable
func (p *PathParams) Validate() error {
	if err := validation.ValidateRequiredString(p.OrgID, "orgId"); err != nil {
		return err
	}
	if err := validation.ValidateRequiredString(p.BenchmarkID, "benchmark_id"); err != nil {
		return err
	}
	if p.ID != "" {
		return validation.ValidateID(p.ID)
	}
	return nil
}
