package varconfig

import (
	repoVarconfig "projeto-crud-credencials/internal/service"
)

// assertions garante que Repository implementa a interface IVarConfigRepository
var _ repoVarconfig.IVarConfigRepository = (*Repository)(nil)
