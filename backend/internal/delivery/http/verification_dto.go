package http

// VerifyReportRequest represents the payload from Puskesmas when verifying a report.
type VerifyReportRequest struct {
	Decision string `json:"decision" validate:"required,oneof=APPROVED REJECTED"`
	Reason   string `json:"reason"`
}
