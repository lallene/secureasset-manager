package change

import (
	"net/http"

	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/cmdb"
	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

type ChangeImpactResponse struct {
	Change          ChangeRequest         `json:"change"`
	BusinessService *cmdb.BusinessService `json:"business_service"`
	Applications    []cmdb.Application    `json:"applications"`
	Databases       []cmdb.Database       `json:"databases"`
	Assets          []asset.Asset         `json:"assets"`
}

func GetChangeImpact(c *gin.Context) {
	id := c.Param("id")

	var change ChangeRequest

	if err := database.DB.
		Preload("BusinessService").
		First(&change, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Changement introuvable",
		})
		return
	}

	if change.BusinessServiceID == nil {
		c.JSON(http.StatusOK, ChangeImpactResponse{
			Change:          change,
			BusinessService: nil,
			Applications:    []cmdb.Application{},
			Databases:       []cmdb.Database{},
			Assets:          []asset.Asset{},
		})
		return
	}

	var businessService cmdb.BusinessService

	if err := database.DB.
		Preload("Site").
		First(&businessService, *change.BusinessServiceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Service métier introuvable",
		})
		return
	}

	var applications []cmdb.Application

	if err := database.DB.
		Model(&cmdb.Application{}).
		Preload("Site").
		Joins(`
			JOIN configuration_relations cr
			ON cr.source_type = ?
			AND cr.source_id = ?
			AND cr.target_type = ?
			AND cr.target_id = applications.id
			AND cr.deleted_at IS NULL
		`, "BusinessService", businessService.ID, "Application").
		Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de charger les applications impactées",
		})
		return
	}

	appIDs := make([]uint, 0, len(applications))
	for _, app := range applications {
		appIDs = append(appIDs, app.ID)
	}

	var databases []cmdb.Database
	var assets []asset.Asset

	if len(appIDs) > 0 {
		if err := database.DB.
			Model(&cmdb.Database{}).
			Preload("Site").
			Joins(`
				JOIN configuration_relations cr
				ON cr.source_type = ?
				AND cr.source_id IN ?
				AND cr.target_type = ?
				AND cr.target_id = databases.id
				AND cr.deleted_at IS NULL
			`, "Application", appIDs, "Database").
			Find(&databases).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Impossible de charger les bases impactées",
			})
			return
		}

		if err := database.DB.
			Model(&asset.Asset{}).
			Preload("Site").
			Joins(`
				JOIN configuration_relations cr
				ON cr.source_type = ?
				AND cr.source_id IN ?
				AND cr.target_type = ?
				AND cr.target_id = assets.id
				AND cr.deleted_at IS NULL
			`, "Application", appIDs, "Asset").
			Find(&assets).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Impossible de charger les assets impactés",
			})
			return
		}
	}

	c.JSON(http.StatusOK, ChangeImpactResponse{
		Change:          change,
		BusinessService: &businessService,
		Applications:    applications,
		Databases:       databases,
		Assets:          assets,
	})
}
