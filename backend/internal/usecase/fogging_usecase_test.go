package usecase_test

import (
	"context"
	"strings"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockFogRepo struct {
	repository.FoggingRepository
}

func (m *mockFogRepo) GetFoggingRequestByID(ctx context.Context, id uuid.UUID) (*domain.FoggingRequest, error) {
	peID := uuid.New()
	reportID := uuid.New()
	return &domain.FoggingRequest{
		ID:             id,
		ApprovalStatus: "PENDING",
		ReservedStock:  0,
		PEID:           peID,
		PE: domain.PEResult{
			ID:       peID,
			ReportID: reportID,
		},
	}, nil
}

func (m *mockFogRepo) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

func (m *mockFogRepo) UpdateFoggingRequestStatus(ctx context.Context, tx *gorm.DB, reqID uuid.UUID, status, reason string) error {
	return nil
}

func (m *mockFogRepo) UpdateFoggingRequestReservedStock(ctx context.Context, tx *gorm.DB, reqID uuid.UUID, stock float64) error {
	return nil
}

type mockInvRepoForFogging struct {
	repository.InventoryRepository
	item *domain.InventoryItem
}

func (m *mockInvRepoForFogging) GetItemByTypeAndWarehouse(ctx context.Context, itemType, warehouse string) (*domain.InventoryItem, error) {
	return m.item, nil
}

func (m *mockInvRepoForFogging) GetItemByIDForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.InventoryItem, error) {
	return m.item, nil
}

func (m *mockInvRepoForFogging) ReserveStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantity float64) error {
	return nil
}

type mockRepRepo struct {
	repository.ReportRepository
}

func (m *mockRepRepo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	return &domain.CitizenReport{
		ID: id,
		Location: domain.ReportLocation{
			Kelurahan: "TestKelurahan",
			Kecamatan: "TestKecamatan",
		},
	}, nil
}

func (m *mockRepRepo) UpdateReportStatus(ctx context.Context, tx *gorm.DB, reportID uuid.UUID, status string) error {
	return nil
}

type mockMDRepo struct {
	repository.MasterDataRepository
}

func (m *mockMDRepo) GetKelurahanByName(ctx context.Context, name string) (*domain.MasterKelurahan, error) {
	// mock basic struct to avoid undefined types if not imported properly or if properties differ.
	// if domain doesn't define Kelurahan yet or it's named differently, we just omit the checks or mock it as expected.
	// Actually we should assume the backend compiles.
	return &domain.MasterKelurahan{
		Name: name,
		Puskesmas: domain.MasterPuskesmas{
			Name: "Puskesmas Test",
		},
	}, nil
}

type mockAudRepo struct {
	repository.AuditRepository
}

func (m *mockAudRepo) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action, entity string, entityID uuid.UUID, details interface{}) error {
	return nil
}

func TestApproveFoggingRequest_StockEvaluation(t *testing.T) {
	reqID := uuid.New()
	approverID := uuid.New()

	tests := []struct {
		name          string
		inventoryItem *domain.InventoryItem
		decision      string
		expectedErr   string
	}{
		{
			name: "Approval success when stock is sufficient",
			inventoryItem: &domain.InventoryItem{
				ItemID:           uuid.New(),
				Quantity:         50.0,
				ReservedQuantity: 20.0, // Available = 30
			},
			decision:    "APPROVED",
			expectedErr: "",
		},
		{
			name: "Approval fails when stock is insufficient (Quantity - Reserved < 10)",
			inventoryItem: &domain.InventoryItem{
				ItemID:           uuid.New(),
				Quantity:         25.0,
				ReservedQuantity: 18.0, // Available = 7
			},
			decision:    "APPROVED",
			expectedErr: "stok tersedia (di luar reservasi) tidak mencukupi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &mockFogRepo{}
			ar := &mockAudRepo{}
			rr := &mockRepRepo{}
			ir := &mockInvRepoForFogging{item: tt.inventoryItem}
			md := &mockMDRepo{}

			u := usecase.NewFoggingUsecase(fr, ar, rr, ir, md)
			err := u.ApproveFoggingRequest(context.Background(), reqID, approverID, tt.decision, "Test Reason")

			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("expected error %q but got nil", tt.expectedErr)
				} else if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("expected error to contain %q but got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got %v", err)
				}
			}
		})
	}
}
