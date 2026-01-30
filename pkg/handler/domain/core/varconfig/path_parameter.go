package varconfig

import (
	"fmt"
)



type PathParameter struct {
	Payload map[string]any `json:"payload"`
}

// ValidateCreateAndUpdateRequest valida a requisição de criar VarConfig
func ValidateCreateAndUpdateRequest(req *PathParameter) error {
	if len(req.Payload) == 0 {
		return fmt.Errorf("payload não pode estar vazio")
	}

	return nil
}
