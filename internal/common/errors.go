package common

import "errors"

var (
	// Erros de Config Service
	ErrConfigNotFound          = errors.New("configuração não encontrada")
	ErrConfigInvalidInput      = errors.New("entrada inválida: orgID e benchmarkID são obrigatórios")
	ErrConfigPayloadValidation = errors.New("validação de payload falhou")
	ErrConfigDuplicateName     = errors.New("configuração com este nome já existe para esta organização e benchmark")

	// Erros de Benchmark Service
	ErrBenchmarkInvalidInput   = errors.New("entrada inválida: nome é obrigatório")
	ErrBenchmarkEmptySchema    = errors.New("schema não pode estar vazio")
	ErrBenchmarkSchemaNotFound = errors.New("schema não encontrado")
	ErrBenchmarkDuplicateName  = errors.New("schema com este nome já existe")
)
