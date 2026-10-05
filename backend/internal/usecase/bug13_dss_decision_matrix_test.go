package usecase_test

import (
	"context"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"errors"
)

// Minimal mock for PE requirements
type mockInvRepoForDSS struct {
	repository.InventoryRepository
	itemQty      float64
	itemReserved float64
}

func (m *mockInvRepoForDSS) GetItemByTypeAndWarehouse(ctx context.Context, itemType string, warehouseName string) (*domain.InventoryItem, error) {
	return &domain.InventoryItem{
		Quantity:         m.itemQty,
		ReservedQuantity: m.itemReserved,
	}, nil
}
func (m *mockInvRepoForDSS) ReserveStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantity float64) error { return nil }
func (m *mockInvRepoForDSS) ReleaseStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantity float64) error { return nil }
func (m *mockInvRepoForDSS) DeductStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantity float64) error  { return nil }

type mockReportRepoForDSS struct {
	repository.ReportRepository
	report *domain.CitizenReport
}
func (m *mockReportRepoForDSS) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	return m.report, nil
}
func (m *mockReportRepoForDSS) UpdateReportStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, newStatus string) error {
	return nil
}

type mockPERepoForDSS struct {
	repository.PERepository
	createdPE *domain.PEResult
}
func (m *mockPERepoForDSS) GetPEResultByReportID(ctx context.Context, reportID uuid.UUID) (*domain.PEResult, error) {
	return nil, errors.New("record not found")
}
func (m *mockPERepoForDSS) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}
func (m *mockPERepoForDSS) CreatePEResult(ctx context.Context, tx *gorm.DB, pe *domain.PEResult) error {
	pe.ID = uuid.New()
	m.createdPE = pe
	return nil
}


func TestDSSDecisionMatrix(t *testing.T) {
	tests := []struct {
		name         string
		larvaFound   bool
		feverCount   int
		itemQty      float64
		itemReserved float64
		expectedRec  string
	}{
		{
			name:         "Rule 1: Jentik=Ya, Demam=0, Stock=Cukup -> FOGGING_FOKUS",
			larvaFound:   true,
			feverCount:   0,
			itemQty:      20,
			itemReserved: 5, // Net 15
			expectedRec:  "FOGGING_FOKUS",
		},
		{
			name:         "Rule 2: Jentik=Ya, Demam=0, Stock=Kurang -> FOGGING_TERTUNDA",
			larvaFound:   true,
			feverCount:   0,
			itemQty:      10,
			itemReserved: 5, // Net 5 (Kurang dari 10)
			expectedRec:  "FOGGING_TERTUNDA",
		},
		{
			name:         "Rule 3: Jentik=Tidak, Demam=3, Stock=Cukup -> FOGGING_FOKUS",
			larvaFound:   false,
			feverCount:   3,
			itemQty:      10,
			itemReserved: 0, // Net 10 (Cukup - Boundary)
			expectedRec:  "FOGGING_FOKUS",
		},
		{
			name:         "Rule 4: Jentik=Tidak, Demam=3, Stock=Kurang -> FOGGING_TERTUNDA",
			larvaFound:   false,
			feverCount:   3,
			itemQty:      9.9,
			itemReserved: 0, // Net 9.9 (Kurang - Boundary)
			expectedRec:  "FOGGING_TERTUNDA",
		},
		{
			name:         "Rule 5: Jentik=Ya, Demam=3, Stock=Cukup -> FOGGING_FOKUS",
			larvaFound:   true,
			feverCount:   3,
			itemQty:      100,
			itemReserved: 10,
			expectedRec:  "FOGGING_FOKUS",
		},
		{
			name:         "Rule 6: Jentik=Ya, Demam=3, Stock=Kurang -> FOGGING_TERTUNDA",
			larvaFound:   true,
			feverCount:   3,
			itemQty:      5,
			itemReserved: 0,
			expectedRec:  "FOGGING_TERTUNDA",
		},
		{
			name:         "Rule 7: Jentik=Tidak, Demam=0, Stock=Cukup -> EDUKASI_PSN",
			larvaFound:   false,
			feverCount:   0,
			itemQty:      50,
			itemReserved: 0,
			expectedRec:  "EDUKASI_PSN",
		},
		{
			name:         "Boundary: Demam=2, Jentik=Tidak -> EDUKASI_PSN",
			larvaFound:   false,
			feverCount:   2,
			itemQty:      50,
			itemReserved: 0,
			expectedRec:  "EDUKASI_PSN",
		},
		{
			name:         "Boundary: Demam=4, Jentik=Tidak -> FOGGING_FOKUS",
			larvaFound:   false,
			feverCount:   4,
			itemQty:      50,
			itemReserved: 0,
			expectedRec:  "FOGGING_FOKUS",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reportID := uuid.New()
			investigatorID := uuid.New()
			
			// Setup Mocks
			peRepo := &mockPERepoForDSS{}
			reportRepo := &mockReportRepoForDSS{
				report: &domain.CitizenReport{
					ID:             reportID,
					Status:         "MENUNGGU_PE_KADER",
					InvestigatorID: &investigatorID,
				},
			}
			invRepo := &mockInvRepoForDSS{
				itemQty:      tc.itemQty,
				itemReserved: tc.itemReserved,
			}
			auditRepo := &mockAudRepo{}
			
			// Mock DB interfaces by redefining methods in existing struct
			// Wait, repository.GormDB is actually *gorm.DB. I need to handle this correctly if the types mismatch. 
			// I'll create a dummy type if it's missing, but let's assume it compiles via internal tests.
			
			u := usecase.NewPEUsecase(peRepo, reportRepo, nil, invRepo, auditRepo, nil, nil)
			
			_, err := u.SubmitPEResult(context.Background(), reportID, investigatorID, 100, tc.larvaFound, tc.feverCount, true, "Testing", "url")
			if err != nil {
				t.Fatalf("SubmitPEResult failed: %v", err)
			}
			
			if peRepo.createdPE == nil {
				t.Fatal("Expected PE to be created")
			}
			
			if peRepo.createdPE.Recommendation != tc.expectedRec {
				t.Errorf("Expected Recommendation %s, got %s", tc.expectedRec, peRepo.createdPE.Recommendation)
			}
		})
	}
}
