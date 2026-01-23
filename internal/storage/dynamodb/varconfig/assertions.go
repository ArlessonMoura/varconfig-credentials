package varconfig

import (
	svcvarconfig "projeto-crud-credencials/internal/service/domain/core/varconfig"
)

// assertions garante que Repository implementa a interface IVarConfigRepository
var _ svcvarconfig.IVarConfigRepository = (*Repository)(nil)
