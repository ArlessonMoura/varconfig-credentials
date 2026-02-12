//USAR MÉTODO VALIDATE!!!!!!
package benchmark_schema

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateBenchmarkID(c *gin.Context) (string, bool) {
	benchmarkID := c.Param("benchmarkId")
	
	if benchmarkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"details": "benchmarkId parameter is required and cannot be empty",
		})
		return "", false
	}
	
	return benchmarkID, true
}