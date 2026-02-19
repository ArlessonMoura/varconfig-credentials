package benchmarkschema

import (
	svchandler "projeto-crud-credentials/pkg/handler"
)


var _ svchandler.IBenchmarkSchemaService = (*Service)(nil)
