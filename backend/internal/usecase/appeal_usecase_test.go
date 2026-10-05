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
	"gorm.io/gorm"
)

type mockAppealReportRepo struct {
	repository.ReportRepository
	mu          sync.Mutex
	status      string
	appealCount int
}

func (m *mockAppealReportRepo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return &domain.CitizenReport{
		ID:          id,
		Status:      m.status,
		AppealCount: m.appealCount,
	}, nil
}

func (m *mockAppealReportRepo) SubmitAppealAtomic(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if (m.status != "REJECTED" && m.status != "REJECTED_HOAX") || m.appealCount >= 1 {
		return errors.New("laporan tidak dapat dibanding: status bukan REJECTED/REJECTED_HOAX atau batas 1x banding telah tercapai")
	}
	m.status = "APPEALED"
	m.appealCount++
	return nil
}

type mockAppealAuditRepo struct {
	repository.AuditRepository
}

func (m *mockAppealAuditRepo) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action string, entity string, entityID uuid.UUID, details interface{}) error {
	return nil
}

func TestSubmitAppeal_1xLimitAndHoaxStatus(t *testing.T) {
	reportID := uuid.New()
	mockRepo := &mockAppealReportRepo{
		status:      "REJECTED_HOAX",
		appealCount: 0,
	}
	mockAudit := &mockAppealAuditRepo{}
	u := usecase.NewReportUsecase(mockRepo, mockAudit, nil, &storage.MockStorageService{})

	ctx := context.Background()

	// 1. First appeal on REJECTED_HOAX -> must PASS
	err := u.SubmitAppeal(ctx, reportID, "Alasan banding pertama")
	if err != nil {
		t.Fatalf("expected first appeal to pass, got: %v", err)
	}

	// 2. Second appeal on same report -> must FAIL (1x limit)
	err2 := u.SubmitAppeal(ctx, reportID, "Alasan banding kedua")
	if err2 == nil {
		t.Fatalf("expected second appeal to fail due to 1x appeal limit, but passed")
	}

	if mockRepo.appealCount != 1 {
		t.Fatalf("expected appealCount to be 1, got: %d", mockRepo.appealCount)
	}
}

func TestSubmitAppeal_ConcurrentRequests(t *testing.T) {
	reportID := uuid.New()
	mockRepo := &mockAppealReportRepo{
		status:      "REJECTED",
		appealCount: 0,
	}
	mockAudit := &mockAppealAuditRepo{}
	u := usecase.NewReportUsecase(mockRepo, mockAudit, nil, &storage.MockStorageService{})

	ctx := context.Background()

	// Simulate 2 concurrent appeal requests
	var wg sync.WaitGroup
	successCount := 0
	failCount := 0
	var mu sync.Mutex

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := u.SubmitAppeal(ctx, reportID, "Concurrent appeal")
			mu.Lock()
			if err == nil {
				successCount++
			} else {
				failCount++
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	if successCount != 1 || failCount != 1 {
		t.Fatalf("expected exactly 1 success and 1 fail in concurrent appeal, got success: %d, fail: %d", successCount, failCount)
	}
}
