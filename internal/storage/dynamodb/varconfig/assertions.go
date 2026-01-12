package varconfig

import (
	svcvarconfig "projeto-crud-credencials/internal/service/domain/varconfig"
)

// assertions garante que Repository implementa a interface VarConfigRepository
var _ svcvarconfig.VarConfigRepository = (*Repository)(nil)
