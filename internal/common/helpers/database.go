package helpers

// isUniqueViolationError verifica se o erro é uma violação de constraint unique
func IsUniqueViolationError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return Contains(errStr, "unique constraint") || Contains(errStr, "duplicate key") || Contains(errStr, "UNIQUE violation")
}

// Contains verifica se uma substring existe em uma string (case-insensitive)
func Contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
			(s[:len(substr)] == substr || 
				s[len(s)-len(substr):] == substr || 
				IndexOf(s, substr) >= 0)))
}

// IndexOf retorna o índice da primeira ocorrência de substr em s
func IndexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// ValidateUniqueViolationError valida se um erro é de violação de unique constraint
// e retorna um erro específico se for o caso
func ValidateUniqueViolationError(err error, specificError error) error {
	if IsUniqueViolationError(err) {
		return specificError
	}
	return err
}
