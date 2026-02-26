package config

import (
	"errors"
	"fmt"
	"strings"
)

type PathParams struct {
	OrgID       string `uri:"orgId" binding:"required"`
	BenchmarkID string `uri:"benchmark_id" binding:"required"`
	ID          string `uri:"id"`
}

// Validate implementa a interface Validatable
func (p *PathParams) Validate() error {
	if err := validateRequiredString(p.OrgID, "orgId"); err != nil {
		return err
	}
	if err := validateRequiredString(p.BenchmarkID, "benchmark_id"); err != nil {
		return err
	}
	if p.ID != "" {
		if err := validateRequiredString(p.ID, "id"); err != nil {
			return err
		}
		if strings.Contains(p.ID, " ") {
			return errors.New("id inválido")
		}
	}
	return nil
}


func validateRequiredString(value, paramName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s não pode estar vazio", paramName)
	}
	return nil
}
