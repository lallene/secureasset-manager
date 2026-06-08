package cmdb

import (
	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/database"
)

type ImpactResponse struct {
	Asset            asset.Asset       `json:"asset"`
	Applications     []Application     `json:"applications"`
	Databases        []Database        `json:"databases"`
	BusinessServices []BusinessService `json:"business_services"`
}

func GetImpact(assetID uint) (*ImpactResponse, error) {
	var impactedAsset asset.Asset
	var applications []Application
	var databases []Database
	var businessServices []BusinessService

	if err := database.DB.
		Preload("Site").
		First(&impactedAsset, assetID).Error; err != nil {
		return nil, err
	}

	if err := database.DB.
		Model(&Application{}).
		Preload("Site").
		Preload("Owner").
		Joins(`
			JOIN configuration_relations cr
			ON cr.source_type = ?
			AND cr.source_id = applications.id
			AND cr.target_type = ?
			AND cr.target_id = ?
			AND cr.deleted_at IS NULL
		`, "Application", "Asset", assetID).
		Find(&applications).Error; err != nil {
		return nil, err
	}

	if len(applications) > 0 {
		appIDs := make([]uint, 0, len(applications))

		for _, app := range applications {
			appIDs = append(appIDs, app.ID)
		}

		if err := database.DB.
			Model(&Database{}).
			Preload("Site").
			Joins(`
				JOIN configuration_relations cr
				ON cr.target_type = ?
				AND cr.target_id = databases.id
				AND cr.source_type = ?
				AND cr.source_id IN ?
				AND cr.deleted_at IS NULL
			`, "Database", "Application", appIDs).
			Find(&databases).Error; err != nil {
			return nil, err
		}

		if err := database.DB.
			Model(&BusinessService{}).
			Preload("Site").
			Preload("Owner").
			Joins(`
				JOIN configuration_relations cr
				ON cr.source_type = ?
				AND cr.source_id = business_services.id
				AND cr.target_type = ?
				AND cr.target_id IN ?
				AND cr.deleted_at IS NULL
			`, "BusinessService", "Application", appIDs).
			Find(&businessServices).Error; err != nil {
			return nil, err
		}
	}

	for i := range applications {
		if applications[i].Owner != nil {
			applications[i].Owner.Password = ""
		}
	}

	for i := range businessServices {
		if businessServices[i].Owner != nil {
			businessServices[i].Owner.Password = ""
		}
	}

	return &ImpactResponse{
		Asset:            impactedAsset,
		Applications:     applications,
		Databases:        databases,
		BusinessServices: businessServices,
	}, nil
}
