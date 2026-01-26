package dto

// CreateImagingResultRequest represents request to create/update imaging result
type CreateImagingResultRequest struct {
	Findings   string   `json:"findings" binding:"required"`
	Impression string   `json:"impression" binding:"required"`
	DICOMFiles []string `json:"dicom_files" binding:"omitempty"`
	IsCritical bool     `json:"is_critical" binding:"omitempty"`
}
