package benchmark

import (
	"projeto-crud-credentials/internal/common/validation"
)

type PathParams struct {
	BenchmarkID string `uri:"benchmarkId" binding:"required"`
}

func (p *PathParams) Validate() error {
	return validation.ValidateRequiredString(p.BenchmarkID, "benchmarkId")
}
