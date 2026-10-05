package dto

import "mime/multipart"

// SubmitReportRequest represents the expected payload (multipart or JSON) from clients.
type SubmitReportRequest struct {
	Name          string                  `form:"name" json:"name"`
	Contact       string                  `form:"contact" json:"contact"`
	AnonymousFlag bool                    `form:"anonymous_flag" json:"anonymous_flag"`
	Category      string                  `form:"category" json:"category"`
	Description   string                  `form:"description" json:"description"`
	Latitude      float64                 `form:"latitude" json:"latitude"`
	Longitude     float64                 `form:"longitude" json:"longitude"`
	Address       string                  `form:"address" json:"address"`
	RT            string                  `form:"rt" json:"rt"`
	RW            string                  `form:"rw" json:"rw"`
	Kelurahan     string                  `form:"kelurahan" json:"kelurahan"`
	Kecamatan     string                  `form:"kecamatan" json:"kecamatan"`
	Photos        []*multipart.FileHeader `form:"photos" json:"-"` // Handled explicitly in controller
}

// ReportResponse represents the successful API response containing the reference number.
type ReportResponse struct {
	ReferenceNo string `json:"reference_no"`
	Message     string `json:"message"`
}

// ReportSummaryResponse is used for public tracking queries to prevent exposing full details.
type ReportSummaryResponse struct {
	ReferenceNo string `json:"reference_no"`
	CreatedAt   string `json:"created_at"`
	Status      string `json:"status"`
	Category    string `json:"category"`
	Kecamatan   string `json:"kecamatan"`
}

// SLAStatsResponse represents SLA breach statistics aggregated by area
type SLAStatsResponse struct {
	Kelurahan   string `json:"kelurahan"`
	RW          string `json:"rw"`
	RT          string `json:"rt"`
	BreachCount int    `json:"breach_count"`
}
