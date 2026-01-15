package varconfig

import (
	svcvarconfig "projeto-crud-credencials/internal/service/domain/varconfig"
)

// assertions garante que Repository implementa a interface VarConfigRepository
// Estou tentando entender melhor como isso funciona
var _ svcvarconfig.VarConfigRepository = (*Repository)(nil)
