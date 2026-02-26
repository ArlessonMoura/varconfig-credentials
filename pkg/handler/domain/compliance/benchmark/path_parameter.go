package benchmark

import (
	"errors"
	"strings"
)

type PathParams struct {
	BenchmarkID string `uri:"benchmarkId" binding:"required"`
}

// Validate implementa a interface Validatable
func (p *PathParams) Validate() error {
	if err := validateRequiredString(p.BenchmarkID, "benchmarkId"); err != nil {
		return err
	}
	return nil
}

func validateRequiredString(value, paramName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(paramName + " parameter is required and cannot be empty")
	}
	return nil
}
