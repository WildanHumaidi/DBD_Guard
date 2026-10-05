package repository

import (
	"context"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MasterDataRepository interface {
	GetPuskesmasByID(ctx context.Context, id uuid.UUID) (*domain.MasterPuskesmas, error)
	GetAllPuskesmas(ctx context.Context) ([]domain.MasterPuskesmas, error)
	
	GetKelurahanByName(ctx context.Context, name string) (*domain.MasterKelurahan, error)
	GetAllKelurahan(ctx context.Context) ([]domain.MasterKelurahan, error)
	GetKelurahanByPuskesmas(ctx context.Context, puskesmasID uuid.UUID) ([]domain.MasterKelurahan, error)
}

type masterDataRepository struct {
	db *gorm.DB
}

func NewMasterDataRepository(db *gorm.DB) MasterDataRepository {
	return &masterDataRepository{db: db}
}

func (r *masterDataRepository) GetPuskesmasByID(ctx context.Context, id uuid.UUID) (*domain.MasterPuskesmas, error) {
	var puskesmas domain.MasterPuskesmas
	err := r.db.WithContext(ctx).First(&puskesmas, "id = ?", id).Error
	return &puskesmas, err
}

func (r *masterDataRepository) GetAllPuskesmas(ctx context.Context) ([]domain.MasterPuskesmas, error) {
	var list []domain.MasterPuskesmas
	err := r.db.WithContext(ctx).Order("name asc").Find(&list).Error
	return list, err
}

func (r *masterDataRepository) GetKelurahanByName(ctx context.Context, name string) (*domain.MasterKelurahan, error) {
	var kelurahan domain.MasterKelurahan
	err := r.db.WithContext(ctx).Preload("Puskesmas").Where("name = ?", name).First(&kelurahan).Error
	return &kelurahan, err
}

func (r *masterDataRepository) GetAllKelurahan(ctx context.Context) ([]domain.MasterKelurahan, error) {
	var list []domain.MasterKelurahan
	err := r.db.WithContext(ctx).Preload("Puskesmas").Order("name asc").Find(&list).Error
	return list, err
}

func (r *masterDataRepository) GetKelurahanByPuskesmas(ctx context.Context, puskesmasID uuid.UUID) ([]domain.MasterKelurahan, error) {
	var list []domain.MasterKelurahan
	err := r.db.WithContext(ctx).Where("puskesmas_id = ?", puskesmasID).Order("name asc").Find(&list).Error
	return list, err
}
