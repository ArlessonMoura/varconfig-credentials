package varconfig

// Service assertions são verificados nos testes unitários
// para evitar ciclos de importação

import (
	svcvarconfig "projeto-crud-credentials/pkg/handler"
)

// assertions garante que Service implementa a interface IVarConfigService
var _ svcvarconfig.IVarConfigService = (*Service)(nil)