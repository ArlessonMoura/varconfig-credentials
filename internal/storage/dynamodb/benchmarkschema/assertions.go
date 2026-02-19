package benchmarkschema

import (
	repoService "projeto-crud-credentials/internal/service"
)

var _ repoService.INoSQLRepository = (*Repository)(nil)
