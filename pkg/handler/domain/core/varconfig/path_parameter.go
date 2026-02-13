package varconfig

import (
	"fmt"
	"strings"
)

type PathParameter struct {
	Payload map[string]any `json:"payload" binding:"required"`
}

// Validate validates the PathParameter struct
func (pp *PathParameter) Validate() error {
	return ValidateCreateAndUpdateRequest(pp)
}

// ValidateCreateAndUpdateRequest valida a requisição de criar VarConfig
func ValidateCreateAndUpdateRequest(req *PathParameter) error {
	if len(req.Payload) == 0 {
		return fmt.Errorf("payload não pode estar vazio")
	}

	return nil
}


// ValidatePathParams validates path parameters for VarConfig operations
func ValidatePathParams(orgID, benchmarkID, id string) error {
	if strings.TrimSpace(orgID) == "" {
		return fmt.Errorf("orgId não pode estar vazio")
	}
	if strings.TrimSpace(benchmarkID) == "" {
		return fmt.Errorf("benchmark_id não pode estar vazio")
	}
	if id != "" {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("id não pode ser vazio")
		}
		if strings.Contains(id, " ") {
			return fmt.Errorf("id inválido")
		}
	}

	return nil
}
}
