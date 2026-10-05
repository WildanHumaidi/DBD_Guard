package http

// ApproveFoggingRequest represents the payload for approving/rejecting a fogging request.
type ApproveFoggingRequest struct {
	Decision string `json:"decision" validate:"required,oneof=APPROVED REJECTED"`
	Reason   string `json:"reason"`
}
