package benchmarkschemas

import (
	repoService "projeto-crud-credentials/internal/service"
)

var _ repoService.IRelationalRepository = (*Repository)(nil)
