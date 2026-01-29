package varconfig

import (
	"fmt"

	varconfigdto "projeto-crud-credencials/dto/varconfig_dto"
)

// ValidateCreateRequest valida a requisição de criar VarConfig
func ValidateCreateRequest(req varconfigdto.CreateVarConfigRequest) error {
	if len(req.Payload) == 0 {
		return fmt.Errorf("payload não pode estar vazio")
	}

	return nil
}

// ValidateUpdateRequest valida a requisição de atualizar VarConfig
func ValidateUpdateRequest(req varconfigdto.UpdateVarConfigRequest) error {
	if len(req.Payload) == 0 {
		return fmt.Errorf("payload não pode estar vazio")
	}

	return nil
}
