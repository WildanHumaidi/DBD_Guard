package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
)

type mockClaimAtomicRepo struct {
	repository.ReportRepository
	mu             sync.Mutex
	status         string
	investigatorID *uuid.UUID
}

func (m *mockClaimAtomicRepo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return &domain.CitizenReport{
		ID:             id,
		Status:         m.status,
		InvestigatorID: m.investigatorID,
	}, nil
}

func (m *mockClaimAtomicRepo) VerifyReportLocation(ctx context.Context, reportID uuid.UUID, locationID uuid.UUID) (bool, error) {
	return true, nil
}

func (m *mockClaimAtomicRepo) ClaimReportAtomic(ctx context.Context, id uuid.UUID, kaderID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.status != "MENUNGGU_PE_KADER" || m.investigatorID != nil {
		return errors.New("laporan sudah diklaim oleh kader lain atau status tidak valid")
	}
	m.investigatorID = &kaderID
	return nil
}

func TestBUGNEW008_ClaimReport_AtomicRaceCondition(t *testing.T) {
	reportID := uuid.New()
	kader1 := uuid.New()
	kader2 := uuid.New()

	repo := &mockClaimAtomicRepo{
		status:         "MENUNGGU_PE_KADER",
		investigatorID: nil,
	}

	u := usecase.NewReportUsecase(repo, nil, nil, &storage.MockStorageService{})
	ctx := context.Background()

	var wg sync.WaitGroup
	successCount := 0
	failCount := 0
	var mu sync.Mutex

	locID := uuid.New()

	claim := func(kaderID uuid.UUID) {
		defer wg.Done()
		err := u.ClaimReport(ctx, reportID, kaderID, &locID)
		mu.Lock()
		defer mu.Unlock()
		if err == nil {
			successCount++
		} else {
			failCount++
		}
	}

	wg.Add(2)
	go claim(kader1)
	go claim(kader2)
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("BUG-NEW-008 Detected! Expected exactly 1 successful claim, got: %d", successCount)
	}
	if failCount != 1 {
		t.Fatalf("Expected exactly 1 failed claim, got: %d", failCount)
	}

	if repo.investigatorID == nil || (*repo.investigatorID != kader1 && *repo.investigatorID != kader2) {
		t.Fatalf("Expected valid assigned investigator, got: %v", repo.investigatorID)
	}
}
