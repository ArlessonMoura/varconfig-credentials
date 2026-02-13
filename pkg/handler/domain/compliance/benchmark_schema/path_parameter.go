package benchmark_schema

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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