package cmdb

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAssetImpact(c *gin.Context) {
	id := c.Param("id")

	assetID, err := strconv.Atoi(id)
	if err != nil || assetID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID asset invalide",
		})
		return
	}

	result, err := GetImpact(uint(assetID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de calculer l'impact",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
