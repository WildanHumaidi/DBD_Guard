package usecase

import (
	"context"

	"github.com/dbd-guard/backend/internal/repository"
)

type GisUsecase interface {
	GetHeatmapData(ctx context.Context) (interface{}, error)
	GetRiskScores(ctx context.Context) (interface{}, error)
}

type gisUsecase struct {
	gisRepo repository.GisRepository
}

func NewGisUsecase(gisRepo repository.GisRepository) GisUsecase {
	return &gisUsecase{gisRepo: gisRepo}
}

// GetHeatmapData returns GeoJSON FeatureCollection
func (u *gisUsecase) GetHeatmapData(ctx context.Context) (interface{}, error) {
	locations, err := u.gisRepo.GetPositiveCaseLocations(ctx)
	if err != nil {
		return nil, err
	}

	type Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"` // [longitude, latitude]
	}

	type Feature struct {
		Type       string                 `json:"type"`
		Geometry   Geometry               `json:"geometry"`
		Properties map[string]interface{} `json:"properties"`
	}

	var features []Feature
	for _, loc := range locations {
		features = append(features, Feature{
			Type: "Feature",
			Geometry: Geometry{
				Type:        "Point",
				Coordinates: []float64{loc.Longitude, loc.Latitude},
			},
			Properties: map[string]interface{}{
				"id":        loc.ReportID,
				"address":   loc.Address,
				"rt":        loc.RT,
				"rw":        loc.RW,
				"kelurahan": loc.Kelurahan,
				"kecamatan": loc.Kecamatan,
				"status":    loc.Status,
			},
		})
	}

	geoJSON := map[string]interface{}{
		"type":     "FeatureCollection",
		"features": features,
	}

	return geoJSON, nil
}

// GetRiskScores calculates risk level per area
func (u *gisUsecase) GetRiskScores(ctx context.Context) (interface{}, error) {
	counts, err := u.gisRepo.GetAreaCaseCounts(ctx)
	if err != nil {
		return nil, err
	}

	type AreaRisk struct {
		Kelurahan string `json:"kelurahan"`
		Kecamatan string `json:"kecamatan"`
		CaseCount int    `json:"case_count"`
		RiskLevel string `json:"risk_level"` // Aman, Waspada, Bahaya
	}

	var riskScores []AreaRisk
	for _, c := range counts {
		level := "Aman"
		if c.CaseCount >= 1 && c.CaseCount <= 3 {
			level = "Waspada"
		} else if c.CaseCount > 3 {
			level = "Bahaya"
		}

		riskScores = append(riskScores, AreaRisk{
			Kelurahan: c.Kelurahan,
			Kecamatan: c.Kecamatan,
			CaseCount: c.CaseCount,
			RiskLevel: level,
		})
	}

	return riskScores, nil
}
