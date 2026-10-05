package repository

import (
	"context"

	"github.com/dbd-guard/backend/internal/domain"
	"gorm.io/gorm"
)

type GisRepository interface {
	GetPositiveCaseLocations(ctx context.Context) ([]LocationWithStatus, error)
	GetAreaCaseCounts(ctx context.Context) ([]AreaCaseCount, error)
}

type AreaCaseCount struct {
	Kelurahan string
	Kecamatan string
	CaseCount int
}

type gisRepository struct {
	db *gorm.DB
}

func NewGisRepository(db *gorm.DB) GisRepository {
	return &gisRepository{db: db}
}

type LocationWithStatus struct {
	domain.ReportLocation
	Status string `json:"status"`
}

func (r *gisRepository) GetPositiveCaseLocations(ctx context.Context) ([]LocationWithStatus, error) {
	var locations []LocationWithStatus
	
	// Join with reports to only get locations of positive/active cases (VERIFIED, PE_SCHEDULED, FOGGING_SCHEDULED, COMPLETED)
	// For MVP heatmap, we can include all valid ones, or specifically ones after PE (positive).
	// Let's include everything except REJECTED for now to populate the map.
	err := r.db.WithContext(ctx).
		Table("report_locations").
		Select("report_locations.*, citizen_reports.status").
		Joins("JOIN citizen_reports ON citizen_reports.report_id = report_locations.report_id").
		Where("citizen_reports.status != ?", "REJECTED").
		Find(&locations).Error
		
	return locations, err
}

func (r *gisRepository) GetAreaCaseCounts(ctx context.Context) ([]AreaCaseCount, error) {
	var counts []AreaCaseCount
	
	err := r.db.WithContext(ctx).
		Table("report_locations").
		Select("kelurahan, kecamatan, count(report_locations.report_id) as case_count").
		Joins("JOIN citizen_reports ON citizen_reports.report_id = report_locations.report_id").
		Where("citizen_reports.status != ?", "REJECTED").
		Group("kelurahan, kecamatan").
		Scan(&counts).Error

	return counts, err
}
