package asset

import (
	"net/http"
	"secureasset-manager/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

var assetService = NewService(NewRepository())

// CreateAsset godoc
// @Summary Créer un asset
// @Description Ajouter un équipement informatique
// @Tags Assets
// @Accept json
// @Produce json
// @Param asset body Asset true "Asset"
// @Success 201 {object} Asset
// @Router /assets [post]
func CreateAsset(c *gin.Context) {
	var asset Asset

	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := assetService.CreateAsset(&asset); err != nil {
		logger.ErrorLogger.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{

			"error": "Impossible de créer l'asset",
		})
		return
	}
	logger.InfoLogger.Println("Asset créé :", asset.Name)
	c.JSON(http.StatusCreated, asset)
}

// GetAssets godoc
// @Summary Liste des assets
// @Description Récupérer tous les équipements
// @Tags Assets
// @Produce json
// @Success 200 {array} Asset
// @Router /assets [get]
func GetAssets(c *gin.Context) {
	assets, err := assetService.GetAssets()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les assets",
		})
		return
	}

	c.JSON(http.StatusOK, assets)
}

func GetAssetByID(c *gin.Context) {
	id := c.Param("id")

	asset, err := assetService.GetAssetByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Asset introuvable",
		})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func calculateDueDate(priority string) time.Time {
	now := time.Now()

	switch priority {
	case "Critical":
		return now.Add(4 * time.Hour)
	case "High":
		return now.Add(8 * time.Hour)
	case "Medium":
		return now.Add(24 * time.Hour)
	case "Low":
		return now.Add(72 * time.Hour)
	default:
		return now.Add(24 * time.Hour)
	}
}

func UpdateAsset(c *gin.Context) {
	id := c.Param("id")

	existingAsset, err := assetService.GetAssetByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Asset introuvable",
		})
		return
	}

	if err := c.ShouldBindJSON(existingAsset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := assetService.UpdateAsset(existingAsset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de modifier l'asset",
		})
		return
	}

	c.JSON(http.StatusOK, existingAsset)
}

func DeleteAsset(c *gin.Context) {
	id := c.Param("id")

	existingAsset, err := assetService.GetAssetByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Asset introuvable",
		})
		return
	}

	if err := assetService.DeleteAsset(existingAsset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de supprimer l'asset",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Asset supprimé avec succès",
	})
}
