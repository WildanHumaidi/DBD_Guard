package http

// SubmitPEResultRequest represents the input from Puskesmas staff for PE results.
type SubmitPEResultRequest struct {
	PeRadius         int    `json:"pe_radius" validate:"required,min=0,max=200"`
	LarvaFound       bool   `json:"larva_found"`
	FeverCount       int    `json:"fever_count" validate:"min=0"`
	PsnEdukasiDone   bool   `json:"psn_edukasi_done"`
	Notes            string `json:"notes" validate:"required"`
	EvidencePhotoUrl string `json:"evidence_photo_url"`
}
