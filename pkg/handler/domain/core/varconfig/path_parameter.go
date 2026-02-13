package varconfig

import (
	"errors"
	"fmt"
	"strings"
)

type PathParameter struct {
	Payload map[string]any `json:"payload" binding:"required"`
	OrgID    string       `json:"org_id"`
	BenchmarkID string       `json:"benchmark_id"`
	ID        string       `json:"id"`
}

// Validate validates the PathParameter struct
func (pp *PathParameter) Validate() error {
	// Validate payload if present
	if pp.Payload != nil {
		if len(pp.Payload) == 0 {
			return errors.New("payload não pode estar vazio")
		}
	}
	
	// Validate orgId
	if err := validateRequiredString(pp.OrgID, "orgId"); err != nil {
		return err
	}
	
	// Validate benchmarkId
	if err := validateRequiredString(pp.BenchmarkID, "benchmark_id"); err != nil {
		return err
	}
	
	// Validate id if present
	if pp.ID != "" {
		if err := validateRequiredString(pp.ID, "id"); err != nil {
			return err
		}
		if strings.Contains(pp.ID, " ") {
			return errors.New("id inválido")
		}
	}
	
	return nil
}

// validateRequiredString validates that a string parameter is not empty after trimming
func validateRequiredString(value, paramName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s não pode estar vazio", paramName)
	}
	return nil
}

// ValidateCreateAndUpdateRequest valida a requisição de criar VarConfig
func ValidateCreateAndUpdateRequest(req *PathParameter) error {
	return req.Validate()
}

// ValidatePathParams validates path parameters for VarConfig operations
func ValidatePathParams(orgID, benchmarkID, id string) error {
	param := &PathParameter{
		OrgID:       orgID,
		BenchmarkID: benchmarkID,
		ID:          id,
	}
	return param.Validate()
}