package varconfig

import (
	repovarconfig "projeto-crud-credencials/internal/service"
)

// assertions garante que Repository implementa a interface IVarConfigRepository
var _ repovarconfig.IVarConfigRepository = (*Repository)(nil)
