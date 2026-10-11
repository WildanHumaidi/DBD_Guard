package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/domain/dto"
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/google/uuid"
)

type ReportUsecase interface {
	UploadPhoto(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
	SubmitCitizenReport(ctx context.Context, req dto.SubmitReportRequest) (*dto.ReportResponse, error)
	GetReports(ctx context.Context, userRole string, locationID *uuid.UUID) ([]domain.CitizenReport, error)
	TrackReport(ctx context.Context, refNo string) (*domain.CitizenReport, error)
	TrackReportsByPhone(ctx context.Context, phone string) ([]dto.ReportSummaryResponse, error)
	SubmitFeedback(ctx context.Context, refNo string, rating int, notes string) error
	GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error)
	SubmitAppeal(ctx context.Context, id uuid.UUID, reason string) error
	ClaimReport(ctx context.Context, id uuid.UUID, kaderID uuid.UUID, locationID *uuid.UUID) error
	CancelReport(ctx context.Context, refNo string) error
	ReturnReportTask(ctx context.Context, id uuid.UUID, kaderID uuid.UUID, reason string) error
	ForceUnassignReport(ctx context.Context, id uuid.UUID) error
	TakeoverReport(ctx context.Context, id uuid.UUID, officerID uuid.UUID, reason string) error
	AutoReleaseExpiredTasks(ctx context.Context) (int, error)
	GetSLAMetrics(ctx context.Context) (map[string]interface{}, error)
	SendFeedbackReminders(ctx context.Context) (int, error)
	GetProtectedPhotoPath(ctx context.Context, userRole string, userLocationID *uuid.UUID, filename string) (string, error)
	MergeDuplicateReport(ctx context.Context, childID, primaryID, kaderID uuid.UUID, reason string) error
}

type reportUsecase struct {
	reportRepo repository.ReportRepository
	auditRepo  repository.AuditRepository
	notifUsecase NotificationUsecase
	storage    storage.StorageService
}

func NewReportUsecase(reportRepo repository.ReportRepository, auditRepo repository.AuditRepository, notifUsecase NotificationUsecase, storageService storage.StorageService) ReportUsecase {
	return &reportUsecase{
		reportRepo:   reportRepo,
		auditRepo:    auditRepo,
		notifUsecase: notifUsecase,
		storage:      storageService,
	}
}

func (u *reportUsecase) SubmitCitizenReport(ctx context.Context, req dto.SubmitReportRequest) (*dto.ReportResponse, error) {
	// 1. Generate Reference Number (e.g. REP-YYYYMMDD-UUID)
	refNo := generateReferenceNo()

	// 2. Prepare Reporter Entity
	reporter := &domain.Reporter{
		Name:          req.Name,
		Contact:       req.Contact,
		AnonymousFlag: req.AnonymousFlag,
	}
	if req.AnonymousFlag {
		reporter.Name = "Anonymous"
		reporter.Contact = ""
	}

	report := &domain.CitizenReport{
		ReferenceNo: refNo,
		Category:    req.Category,
		Description: req.Description,
		Status:      "MENUNGGU_PE_KADER", // Initial status
	}

	// PRD Gap: Clustering Laporan Duplikat (<100m, 7 hari)
	isDuplicate, err := u.reportRepo.CheckDuplicateRadius(ctx, req.Latitude, req.Longitude, 100, 7)
	if err == nil && isDuplicate {
		report.IsPossibleDuplicate = true
	}

	// PRD Gap: Spam Clustering Check (<50m, 1 jam, >= 3 laporan)
	spamCount, err := u.reportRepo.CheckSpamClustering(ctx, req.Latitude, req.Longitude, 50, 1)
	if err == nil && spamCount >= 3 {
		report.Status = "REJECTED" // Automatically reject if it's considered spam
	}

	// 4. Prepare ReportLocation Entity
	geomStr := fmt.Sprintf("POINT(%f %f)", req.Longitude, req.Latitude)
	location := &domain.ReportLocation{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Geom:      geomStr,
		Address:   req.Address,
		RT:        req.RT,
		RW:        req.RW,
		Kelurahan: req.Kelurahan,
		Kecamatan: req.Kecamatan,
	}

	// 5. Process Photos
	var photos []domain.ReportPhoto
	for _, fileHeader := range req.Photos {
		storageKey, err := u.UploadPhoto(ctx, fileHeader)
		if err != nil {
			return nil, fmt.Errorf("failed to save photo: %w", err)
		}
		photos = append(photos, domain.ReportPhoto{
			StorageKey: storageKey,
		})
	}

	// 6. Save using Transaction in Repository
	err = u.reportRepo.CreateReportTransaction(ctx, reporter, report, location, photos)
	if err != nil {
		return nil, fmt.Errorf("failed to save citizen report: %w", err)
	}

	if report.Status == "REJECTED" {
		if u.auditRepo != nil {
			_ = u.auditRepo.LogAction(ctx, nil, uuid.Nil, "AUTO_REJECT_SPAM", "CitizenReport", report.ID, map[string]interface{}{
				"reason":     "SPAM_CLUSTERING_DETECTED",
				"spam_count": spamCount,
			})
		}
		return &dto.ReportResponse{
			ReferenceNo: refNo,
			Message:     "Laporan terdeteksi sebagai spam (duplikasi berulang dalam waktu singkat) dan ditolak otomatis oleh sistem.",
		}, nil
	}

	if u.notifUsecase != nil {
		_ = u.notifUsecase.NotifyNewReport(ctx, report.ID.String(), location.RT, location.RW, location.Kelurahan)
	}

	return &dto.ReportResponse{
		ReferenceNo: refNo,
		Message:     "Laporan berhasil dikirim dan sedang dalam proses penyelidikan (PE) oleh Kader.",
	}, nil
}

func (u *reportUsecase) GetReports(ctx context.Context, userRole string, locationID *uuid.UUID) ([]domain.CitizenReport, error) {
	// If role is Petugas or Kepala Puskesmas, they should only see reports in their Puskesmas
	// If role is Kader, they should only see reports in their assigned area
	if userRole == domain.RolePetugasPuskesmas || userRole == domain.RoleKepalaPuskesmas || userRole == domain.RoleKader {
		return u.reportRepo.GetReports(ctx, locationID)
	}
	
	// Other roles (e.g., Admin Dinkes, Admin Developer) see all reports
	return u.reportRepo.GetReports(ctx, nil)
}

func (u *reportUsecase) TrackReport(ctx context.Context, refNo string) (*domain.CitizenReport, error) {
	report, err := u.reportRepo.GetReportByReference(ctx, refNo)
	if err != nil {
		return nil, err
	}
	if report != nil && report.Reporter.AnonymousFlag {
		report.Reporter.Name = "Anonymous"
		report.Reporter.Contact = ""
	}
	return report, nil
}

func (u *reportUsecase) TrackReportsByPhone(ctx context.Context, phone string) ([]dto.ReportSummaryResponse, error) {
	reports, err := u.reportRepo.GetReportsByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	var summaries []dto.ReportSummaryResponse
	for _, r := range reports {
		summary := dto.ReportSummaryResponse{
			ReferenceNo: r.ReferenceNo,
			CreatedAt:   r.CreatedAt.Format(time.RFC3339),
			Status:      r.Status,
			Category:    r.Category,
		}
		if r.Location.Kecamatan != "" {
			summary.Kecamatan = r.Location.Kecamatan
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (u *reportUsecase) SubmitFeedback(ctx context.Context, refNo string, rating int, notes string) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating harus antara 1 sampai 5")
	}

	report, err := u.reportRepo.GetReportByReference(ctx, refNo)
	if err != nil {
		return err
	}
	if report.Status != "COMPLETED" {
		return fmt.Errorf("feedback hanya dapat diberikan setelah fogging selesai (status: COMPLETED)")
	}
	if report.FeedbackRating > 0 {
		return fmt.Errorf("feedback/ulasan sudah pernah diberikan untuk laporan ini")
	}

	if err := u.reportRepo.SubmitFeedback(ctx, refNo, rating, notes); err != nil {
		return err
	}

	auditDetails := map[string]interface{}{
		"rating": rating,
		"notes":  notes,
	}
	return u.auditRepo.LogAction(ctx, nil, uuid.Nil, "SUBMIT_FEEDBACK", "CitizenReport", report.ID, auditDetails)
}

func (u *reportUsecase) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	return u.reportRepo.GetReportByID(ctx, id)
}

func (u *reportUsecase) SubmitAppeal(ctx context.Context, id uuid.UUID, reason string) error {
	// Atomic conditional update ensuring status is REJECTED or REJECTED_HOAX and appeal_count < 1
	if err := u.reportRepo.SubmitAppealAtomic(ctx, id); err != nil {
		return err
	}

	// Log Audit
	systemActorID := uuid.Nil // or get user ID if authenticated
	auditDetails := map[string]interface{}{
		"reason": reason,
	}
	return u.auditRepo.LogAction(ctx, nil, systemActorID, "SUBMIT_APPEAL", "CitizenReport", id, auditDetails)
}


// generateReferenceNo generates a reference number like REP-20260910-A1B2
func generateReferenceNo() string {
	dateStr := time.Now().Format("20060102")
	uuidStr := strings.Split(uuid.New().String(), "-")[0] // take first part of UUID
	return fmt.Sprintf("REP-%s-%s", dateStr, strings.ToUpper(uuidStr))
}

// UploadPhoto uploads the multipart file using StorageService and returns the path
func (u *reportUsecase) UploadPhoto(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Generate safe filename
	ext := filepath.Ext(fileHeader.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	
	// Upload to default bucket (e.g., "reports")
	bucketName := "dbd-guard-photos"
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	objectName, err := u.storage.Upload(ctx, bucketName, newFilename, file, fileHeader.Size, contentType)
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (u *reportUsecase) ClaimReport(ctx context.Context, id uuid.UUID, kaderID uuid.UUID, locationID *uuid.UUID) error {
	if locationID == nil {
		return domain.ErrMissingLocationContext
	}

	report, err := u.reportRepo.GetReportByID(ctx, id)
	if err != nil {
		return err
	}

	allowed, err := u.reportRepo.VerifyReportLocation(ctx, id, *locationID)
	if err != nil {
		return err
	}
	if !allowed {
		return domain.ErrUnauthorizedJurisdiction
	}

	if report.Status != "MENUNGGU_PE_KADER" {
		return fmt.Errorf("laporan tidak dalam status MENUNGGU_PE_KADER")
	}

	if report.InvestigatorID != nil {
		return fmt.Errorf("laporan sudah diklaim oleh kader lain")
	}

	if err := u.reportRepo.ClaimReportAtomic(ctx, id, kaderID); err != nil {
		return err
	}

	if u.auditRepo != nil {
		_ = u.auditRepo.LogAction(ctx, nil, kaderID, "CLAIM_REPORT", "CitizenReport", id, map[string]interface{}{
			"claimed_by":  kaderID,
			"location_id": *locationID,
		})
	}

	return nil
}

func (u *reportUsecase) CancelReport(ctx context.Context, refNo string) error {
	report, err := u.reportRepo.GetReportByReference(ctx, refNo)
	if err != nil {
		return err
	}

	if report.Status != "MENUNGGU_PE_KADER" {
		return fmt.Errorf("laporan sudah diproses dan tidak bisa dibatalkan")
	}

	if report.InvestigatorID != nil {
		return fmt.Errorf("laporan sudah diambil oleh kader dan tidak bisa dibatalkan")
	}

	err = u.reportRepo.UpdateReportStatus(ctx, nil, report.ID, "CANCELLED_BY_USER")
	if err != nil {
		return err
	}

	if u.auditRepo != nil {
		auditDetails := map[string]interface{}{
			"reference_no": refNo,
			"status":       "CANCELLED_BY_USER",
		}
		_ = u.auditRepo.LogAction(ctx, nil, uuid.Nil, "CANCEL_REPORT_BY_USER", "CitizenReport", report.ID, auditDetails)
	}

	return nil
}

func (u *reportUsecase) ReturnReportTask(ctx context.Context, id uuid.UUID, kaderID uuid.UUID, reason string) error {
	report, err := u.reportRepo.GetReportByID(ctx, id)
	if err != nil {
		return err
	}

	if report.InvestigatorID == nil || *report.InvestigatorID != kaderID {
		return fmt.Errorf("laporan ini tidak diklaim oleh anda")
	}

	err = u.reportRepo.UnassignInvestigator(ctx, id)
	if err != nil {
		return err
	}

	// Log reason to audit trail
	auditDetails := map[string]interface{}{
		"reason": reason,
	}
	return u.auditRepo.LogAction(ctx, nil, kaderID, "RETURN_TASK", "CitizenReport", id, auditDetails)
}

func (u *reportUsecase) ForceUnassignReport(ctx context.Context, id uuid.UUID) error {
	report, err := u.reportRepo.GetReportByID(ctx, id)
	if err != nil {
		return err
	}

	if report.InvestigatorID == nil {
		return fmt.Errorf("laporan belum diklaim oleh kader manapun")
	}

	prevInvestigator := *report.InvestigatorID
	if err := u.reportRepo.UnassignInvestigator(ctx, id); err != nil {
		return err
	}

	if u.auditRepo != nil {
		auditDetails := map[string]interface{}{
			"previous_investigator_id": prevInvestigator,
			"action":                   "ADMIN_FORCE_UNASSIGN",
			"timestamp":                time.Now(),
		}
		_ = u.auditRepo.LogAction(ctx, nil, uuid.Nil, "FORCE_UNASSIGN_REPORT", "CitizenReport", id, auditDetails)
	}

	return nil
}

func (u *reportUsecase) TakeoverReport(ctx context.Context, id uuid.UUID, officerID uuid.UUID, reason string) error {
	report, err := u.reportRepo.GetReportByID(ctx, id)
	if err != nil {
		return err
	}

	// Status validation: report must be in actionable pre-PE status
	if report.Status != "MENUNGGU_PE_KADER" && report.Status != "ESKALASI_PUSKESMAS" && report.Status != "APPEALED" {
		return fmt.Errorf("laporan dengan status %s tidak dapat diambil alih", report.Status)
	}

	if reason == "" {
		reason = "Puskesmas direct investigation takeover"
	}

	// Record audit trail with previous investigator and reason
	auditDetails := map[string]interface{}{
		"previous_investigator_id": report.InvestigatorID,
		"new_investigator_id":      officerID,
		"reason":                   reason,
		"timestamp":                time.Now(),
	}
	if err := u.auditRepo.LogAction(ctx, nil, officerID, "TAKEOVER_REPORT", "CitizenReport", id, auditDetails); err != nil {
		return err
	}

	return u.reportRepo.AssignInvestigator(ctx, id, officerID)
}

func (u *reportUsecase) AutoReleaseExpiredTasks(ctx context.Context) (int, error) {
	// PRD 5.1 & Business Rule #1: SLA 3x24 jam (72 hours) sejak laporan masuk
	cutoff := time.Now().Add(-72 * time.Hour)
	reports, err := u.reportRepo.GetUnverifiedOverdueReports(ctx, cutoff)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, report := range reports {
		// Escalate report to ESKALASI_PUSKESMAS and unassign stagnant investigator
		if err := u.reportRepo.EscalateReport(ctx, report.ID); err == nil {
			count++
			// Log the audit for SLA escalation
			auditDetails := map[string]interface{}{
				"reason":     "SLA_72H_EXPIRED_RT_UNRESPONSIVE",
				"kader_id":   report.InvestigatorID,
				"created_at": report.CreatedAt,
			}
			_ = u.auditRepo.LogAction(ctx, nil, uuid.Nil, "SLA_ESCALATION", "CitizenReport", report.ID, auditDetails)

			if u.notifUsecase != nil && report.Location.Kelurahan != "" {
				_ = u.notifUsecase.NotifyNewReport(ctx, report.ID.String(), report.Location.RT, report.Location.RW, report.Location.Kelurahan)
			}
		}
	}
	return count, nil
}

func (u *reportUsecase) GetSLAMetrics(ctx context.Context) (map[string]interface{}, error) {
	return u.reportRepo.GetSLAMetrics(ctx)
}

func (u *reportUsecase) SendFeedbackReminders(ctx context.Context) (int, error) {
	// Find reports completed at least 14 days ago without feedback requested yet
	reports, err := u.reportRepo.GetCompletedReportsForFeedback(ctx, 14)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, report := range reports {
		// Skip feedback reminders for anonymous reports to protect privacy (PRD Rule 5)
		if report.Reporter.AnonymousFlag {
			continue
		}

		// Send notification
		// We'll use a dummy FCM token since it's an MVP, in real world we need to get user's token or send SMS/Email
		if u.notifUsecase != nil {
			_ = fmt.Sprintf("Penanganan kasus di %s telah selesai 14 hari yang lalu. Mohon luangkan waktu untuk mengisi evaluasi efektivitas fogging.", report.Location.Address)
			_ = u.notifUsecase.NotifyNewReport(ctx, report.ID.String(), report.Location.RT, report.Location.RW, report.Location.Kelurahan) // Reusing existing method loosely
		}

		err = u.reportRepo.MarkFeedbackRequested(ctx, report.ID)
		if err == nil {
			count++
			
			// Audit log
			auditDetails := map[string]interface{}{
				"reason": "14_DAYS_POST_FOGGING_EVALUATION",
				"reference_no": report.ReferenceNo,
			}
			_ = u.auditRepo.LogAction(ctx, nil, uuid.Nil, "SEND_FEEDBACK_REMINDER", "CitizenReport", report.ID, auditDetails)
		}
	}
	
	return count, nil
}

func (u *reportUsecase) GetProtectedPhotoPath(ctx context.Context, userRole string, userLocationID *uuid.UUID, filename string) (string, error) {
	// 1. Sanitize filename to prevent directory traversal
	if filename == "" || filepath.Base(filename) != filename || strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return "", errors.New("nama file tidak valid atau terdeteksi path traversal")
	}

	// 2. Check existence using StorageService
	bucketName := "dbd-guard-photos"
	exists, err := u.storage.Exists(ctx, bucketName, filename)
	if err != nil || !exists {
		return "", errors.New("file foto tidak ditemukan")
	}

	generateURL := func() (string, error) {
		return u.storage.GetPresignedURL(ctx, bucketName, filename, 15*time.Minute)
	}

	// 3. Super Admin & Dinkes have city-wide surveillance authority
	if userRole == domain.RoleAdminDeveloper || userRole == domain.RoleDinkes {
		return generateURL()
	}

	// 4. Role & Location scoping for other users
	if userLocationID == nil {
		return "", errors.New("akses ditolak: akun tidak memiliki wilayah kerja yang valid")
	}

	// Try lookup in ReportPhoto
	photo, err := u.reportRepo.GetPhotoByStorageKey(ctx, filename)
	if err == nil && photo != nil {
		allowed, err := u.reportRepo.VerifyReportLocation(ctx, photo.ReportID, *userLocationID)
		if err != nil || !allowed {
			return "", errors.New("akses ditolak: wilayah kerja Anda tidak memiliki izin mengakses foto pada laporan ini")
		}
		return generateURL()
	}

	// Try lookup in PEResult
	peResult, err := u.reportRepo.GetPEResultByPhotoURL(ctx, filename)
	if err == nil && peResult != nil {
		allowed, err := u.reportRepo.VerifyReportLocation(ctx, peResult.ReportID, *userLocationID)
		if err != nil || !allowed {
			return "", errors.New("akses ditolak: wilayah kerja Anda tidak memiliki izin mengakses foto PE pada laporan ini")
		}
		return generateURL()
	}

	// File exists on disk but is not associated with any authorized report or PE
	return "", errors.New("akses ditolak: tidak memiliki izin mengakses dokumen atau foto sensitif ini")
}

func (u *reportUsecase) MergeDuplicateReport(ctx context.Context, childID, primaryID, kaderID uuid.UUID, reason string) error {
	if childID == primaryID {
		return errors.New("laporan tidak dapat digabungkan dengan dirinya sendiri")
	}

	// Validate reports exist
	childReport, err := u.reportRepo.GetReportByID(ctx, childID)
	if err != nil {
		return err
	}
	primaryReport, err := u.reportRepo.GetReportByID(ctx, primaryID)
	if err != nil {
		return err
	}

	if childReport.Status != "MENUNGGU_PE_KADER" {
		return errors.New("hanya laporan berstatus MENUNGGU_PE_KADER yang dapat di-merge")
	}

	if primaryReport.Status == "MERGED_DUPLICATE" || primaryReport.Status == "REJECTED" || primaryReport.Status == "REJECTED_HOAX" || primaryReport.Status == "CANCELLED_BY_USER" {
		return errors.New("laporan utama tidak valid atau sudah digabungkan/ditutup")
	}

	// Do merge
	if err := u.reportRepo.MergeDuplicateReportsAtomic(ctx, childID, primaryID, kaderID); err != nil {
		return err
	}

	if u.auditRepo != nil {
		auditDetails := map[string]interface{}{
			"primary_report_id": primaryID,
			"child_report_id":   childID,
			"reason":            reason,
			"merged_by":         kaderID,
		}
		_ = u.auditRepo.LogAction(ctx, nil, kaderID, "MERGE_DUPLICATE", "CitizenReport", childID, auditDetails)
	}

	return nil
}

