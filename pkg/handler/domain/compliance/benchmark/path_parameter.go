package benchmark

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type PathParameter struct {
	BenchmarkID string `json:"benchmarkId"`
	SchemaID    string `json:"id"`
}

// Validate validates the PathParameter struct
func (pp *PathParameter) Validate() error {
	// Validate benchmarkId
	if err := validateRequiredString(pp.BenchmarkID, "benchmarkId"); err != nil {
		return err
	}

	// Validate schemaId if present
	if pp.SchemaID != "" {
		if err := validateRequiredString(pp.SchemaID, "id"); err != nil {
			return err
		}
		if strings.Contains(pp.SchemaID, " ") {
			return errors.New("id inválido")
		}
	}

	return nil
}

// validateRequiredString validates that a string parameter is not empty after trimming
func validateRequiredString(value, paramName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(paramName + " parameter is required and cannot be empty")
	}
	return nil
}

// ValidateBenchmarkID validates the benchmark ID parameter
func ValidateBenchmarkID(c *gin.Context) (string, bool) {
	benchmarkID := strings.TrimSpace(c.Param("benchmarkId"))

	if benchmarkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"details": "benchmarkId parameter is required and cannot be empty",
		})
		return "", false
	}

	return benchmarkID, true
}

// ValidateSchemaID validates the schema ID parameter
func ValidateSchemaID(c *gin.Context) (string, bool) {
	schemaID := strings.TrimSpace(c.Param("id"))

	if schemaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"details": "id parameter is required and cannot be empty",
		})
		return "", false
	}

	if strings.Contains(schemaID, " ") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"details": "id parameter cannot contain spaces",
		})
		return "", false
	}

	return schemaID, true
}
