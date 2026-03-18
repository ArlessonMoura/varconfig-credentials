package validation

import (
	"errors"
	"strings"
)

func ValidateRequiredString(value, paramName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(paramName + " não pode estar vazio")
	}
	return nil
}

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("nome é obrigatório")
	}
	return nil
}

func ValidateNoSpaces(value, paramName string) error {
	if strings.Contains(value, " ") {
		return errors.New(paramName + " inválido")
	}
	return nil
}

func ValidateID(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id é obrigatório e não pode estar vazio")
	}
	if strings.Contains(id, " ") {
		return errors.New("id inválido")
	}
	return nil
}

func ValidateRequiredField(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(fieldName + " é obrigatório")
	}
	return nil
}
