package cmdb

import (
	"fmt"
	"net/http"

	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

type GraphNode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Position map[string]float64     `json:"position"`
	Data     map[string]interface{} `json:"data"`
}

type GraphEdge struct {
	ID        string `json:"id"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	Label     string `json:"label"`
	Animated  bool   `json:"animated"`
	MarkerEnd string `json:"markerEnd,omitempty"`
}

type GraphResponse struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

func GetGraphHandler(c *gin.Context) {
	var businessServices []BusinessService
	var applications []Application
	var databases []Database
	var assets []asset.Asset
	var relations []ConfigurationRelation

	if err := database.DB.Preload("Site").Find(&businessServices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de charger les services métiers"})
		return
	}

	if err := database.DB.Preload("Site").Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de charger les applications"})
		return
	}

	if err := database.DB.Preload("Site").Find(&databases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de charger les bases de données"})
		return
	}

	if err := database.DB.Preload("Site").Find(&assets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de charger les assets"})
		return
	}

	if err := database.DB.Find(&relations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de charger les relations CMDB"})
		return
	}

	nodes := []GraphNode{}
	edges := []GraphEdge{}

	addedNodes := map[string]bool{}

	addNode := func(id string, nodeType string, x float64, y float64, data map[string]interface{}) {
		if addedNodes[id] {
			return
		}

		nodes = append(nodes, GraphNode{
			ID:   id,
			Type: "default",
			Position: map[string]float64{
				"x": x,
				"y": y,
			},
			Data: data,
		})

		addedNodes[id] = true
	}

	businessServiceByID := map[uint]BusinessService{}
	applicationByID := map[uint]Application{}
	databaseByID := map[uint]Database{}
	assetByID := map[uint]asset.Asset{}

	for _, item := range businessServices {
		businessServiceByID[item.ID] = item
	}

	for _, item := range applications {
		applicationByID[item.ID] = item
	}

	for _, item := range databases {
		databaseByID[item.ID] = item
	}

	for _, item := range assets {
		assetByID[item.ID] = item
	}

	bsIndex := 0
	appIndex := 0
	dbIndex := 0
	assetIndex := 0

	for _, rel := range relations {
		var sourceNodeID string
		var targetNodeID string

		switch rel.SourceType {
		case "BusinessService":
			if item, ok := businessServiceByID[rel.SourceID]; ok {
				sourceNodeID = fmt.Sprintf("bs-%d", item.ID)
				addNode(
					sourceNodeID,
					"business_service",
					0,
					float64(bsIndex*120),
					map[string]interface{}{
						"label":       "🟣 " + item.Name,
						"ci_type":     "Business Service",
						"name":        item.Name,
						"criticality": item.Criticality,
						"status":      item.Status,
						"site":        item.Site.Name,
					},
				)
				bsIndex++
			}

		case "Application":
			if item, ok := applicationByID[rel.SourceID]; ok {
				sourceNodeID = fmt.Sprintf("app-%d", item.ID)
				addNode(
					sourceNodeID,
					"application",
					350,
					float64(appIndex*120),
					map[string]interface{}{
						"label":       "🔵 " + item.Name,
						"ci_type":     "Application",
						"name":        item.Name,
						"criticality": item.Criticality,
						"status":      item.Status,
						"version":     item.Version,
						"site":        item.Site.Name,
					},
				)
				appIndex++
			}

		case "Database":
			if item, ok := databaseByID[rel.SourceID]; ok {
				sourceNodeID = fmt.Sprintf("db-%d", item.ID)
				addNode(
					sourceNodeID,
					"database",
					700,
					float64(dbIndex*120),
					map[string]interface{}{
						"label":       "🟢 " + item.Name,
						"ci_type":     "Database",
						"name":        item.Name,
						"engine":      item.Engine,
						"version":     item.Version,
						"environment": item.Environment,
						"site":        item.Site.Name,
					},
				)
				dbIndex++
			}

		case "Asset":
			if item, ok := assetByID[rel.SourceID]; ok {
				sourceNodeID = fmt.Sprintf("asset-%d", item.ID)
				addNode(
					sourceNodeID,
					"asset",
					1050,
					float64(assetIndex*120),
					map[string]interface{}{
						"label":   "🟠 " + item.Name,
						"ci_type": "Asset",
						"name":    item.Name,
						"type":    item.Type,
						"status":  item.Status,
						"ip":      item.IPAddress,
						"site":    item.Site.Name,
					},
				)
				assetIndex++
			}
		}

		switch rel.TargetType {
		case "BusinessService":
			if item, ok := businessServiceByID[rel.TargetID]; ok {
				targetNodeID = fmt.Sprintf("bs-%d", item.ID)
				addNode(
					targetNodeID,
					"business_service",
					0,
					float64(bsIndex*120),
					map[string]interface{}{
						"label":       "🟣 " + item.Name,
						"ci_type":     "Business Service",
						"name":        item.Name,
						"criticality": item.Criticality,
						"status":      item.Status,
						"site":        item.Site.Name,
					},
				)
				bsIndex++
			}

		case "Application":
			if item, ok := applicationByID[rel.TargetID]; ok {
				targetNodeID = fmt.Sprintf("app-%d", item.ID)
				addNode(
					targetNodeID,
					"application",
					350,
					float64(appIndex*120),
					map[string]interface{}{
						"label":       "🔵 " + item.Name,
						"ci_type":     "Application",
						"name":        item.Name,
						"criticality": item.Criticality,
						"status":      item.Status,
						"version":     item.Version,
						"site":        item.Site.Name,
					},
				)
				appIndex++
			}

		case "Database":
			if item, ok := databaseByID[rel.TargetID]; ok {
				targetNodeID = fmt.Sprintf("db-%d", item.ID)
				addNode(
					targetNodeID,
					"database",
					700,
					float64(dbIndex*120),
					map[string]interface{}{
						"label":       "🟢 " + item.Name,
						"ci_type":     "Database",
						"name":        item.Name,
						"engine":      item.Engine,
						"version":     item.Version,
						"environment": item.Environment,
						"site":        item.Site.Name,
					},
				)
				dbIndex++
			}

		case "Asset":
			if item, ok := assetByID[rel.TargetID]; ok {
				targetNodeID = fmt.Sprintf("asset-%d", item.ID)
				addNode(
					targetNodeID,
					"asset",
					1050,
					float64(assetIndex*120),
					map[string]interface{}{
						"label":   "🟠 " + item.Name,
						"ci_type": "Asset",
						"name":    item.Name,
						"type":    item.Type,
						"status":  item.Status,
						"ip":      item.IPAddress,
						"site":    item.Site.Name,
					},
				)
				assetIndex++
			}
		}

		if sourceNodeID != "" && targetNodeID != "" {
			edges = append(edges, GraphEdge{
				ID:        fmt.Sprintf("edge-%d", rel.ID),
				Source:    sourceNodeID,
				Target:    targetNodeID,
				Label:     rel.RelationType,
				Animated:  false,
				MarkerEnd: "arrowclosed",
			})
		}
	}

	c.JSON(http.StatusOK, GraphResponse{
		Nodes: nodes,
		Edges: edges,
	})
}
