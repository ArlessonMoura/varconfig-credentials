package config

import (
	repoService "projeto-crud-credentials/internal/service"
)

var _ repoService.IVarConfigRepository = (*Repository)(nil)
