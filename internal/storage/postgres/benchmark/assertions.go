package benchmarks

import (
	repoService "projeto-crud-credentials/internal/service"
)

var _ repoService.IRelationalRepository = (*Repository)(nil)
