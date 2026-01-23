package varconfig

import (
	svcvarconfig "projeto-crud-credencials/internal/service/domain/core/varconfig"
)

// assertions garante que Repository implementa a interface VarConfigRepository
// Estou tentando entender melhor como isso funciona
var _ svcvarconfig.IVarConfigRepository = (*Repository)(nil)
