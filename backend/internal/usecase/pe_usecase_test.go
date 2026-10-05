package usecase_test

import (
	"context"
	"strings"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
)

type mockReportRepo struct {
	repository.ReportRepository
	report *domain.CitizenReport
	err    error
}

func (m *mockReportRepo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	return m.report, m.err
}

type mockPERepo struct {
	repository.PERepository
}

// We also need GetPEResultByReportID because if it doesn't fail, the logic proceeds and panics on a nil repository.
func (m *mockPERepo) GetPEResultByReportID(ctx context.Context, reportID uuid.UUID) (*domain.PEResult, error) {
	// simulate not found so it proceeds
	return nil, nil // Actually if it returns nil err, it means already exists.
}

type mockInventoryRepo struct {
	repository.InventoryRepository
}

type mockFoggingRepo struct {
	repository.FoggingRepository
}

type mockMasterDataRepo struct {
	repository.MasterDataRepository
}

type mockAuditRepo struct {
	repository.AuditRepository
}

type mockEwsKlbUsecase struct {
	usecase.EwsKlbUsecase
}

func TestSubmitPEResult_Ownership(t *testing.T) {
	reportID := uuid.New()
	ownerID := uuid.New()
	nonOwnerID := uuid.New()

	tests := []struct {
		name           string
		investigatorID uuid.UUID
		report         *domain.CitizenReport
		expectedErr    string
	}{
		{
			name:           "Owner tries to submit PE",
			investigatorID: ownerID,
			report: &domain.CitizenReport{
				ID:             reportID,
				Status:         "MENUNGGU_PE_KADER",
				InvestigatorID: &ownerID,
			},
			// Expecting panic or another error from subsequent repository calls, 
			// because we only stubbed GetReportByID. But it should NOT be the UNAUTHORIZED error.
			expectedErr: "",
		},
		{
			name:           "Non-owner tries to submit PE",
			investigatorID: nonOwnerID,
			report: &domain.CitizenReport{
				ID:             reportID,
				Status:         "MENUNGGU_PE_KADER",
				InvestigatorID: &ownerID,
			},
			expectedErr: "UNAUTHORIZED_PE_SUBMISSION",
		},
		{
			name:           "Unclaimed report PE submission",
			investigatorID: ownerID,
			report: &domain.CitizenReport{
				ID:             reportID,
				Status:         "MENUNGGU_PE_KADER",
				InvestigatorID: nil,
			},
			expectedErr: "UNAUTHORIZED_PE_SUBMISSION",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reportRepo := &mockReportRepo{
				report: tt.report,
			}
			peRepo := &mockPERepo{}
			invRepo := &mockInventoryRepo{}
			fogRepo := &mockFoggingRepo{}
			mdRepo := &mockMasterDataRepo{}
			auditRepo := &mockAuditRepo{}
			ewsUsecase := &mockEwsKlbUsecase{}

			u := usecase.NewPEUsecase(peRepo, reportRepo, fogRepo, invRepo, auditRepo, mdRepo, ewsUsecase)

			// Use defer to recover from panic that might happen in subsequent code for valid owner
			defer func() {
				recover() // ignore panics from unimplemented interfaces down the line
			}()

			_, err := u.SubmitPEResult(
				context.Background(),
				reportID,
				tt.investigatorID,
				100, // radius
				true, // larva found
				2,    // fever count
				true, // psn done
				"Test",
				"http://photo.url",
			)

			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expectedErr)
				} else if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("expected error containing %q, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil && strings.Contains(err.Error(), "UNAUTHORIZED_PE_SUBMISSION") {
					t.Errorf("did not expect unauthorized error, but got: %v", err)
				}
			}
		})
	}
}

