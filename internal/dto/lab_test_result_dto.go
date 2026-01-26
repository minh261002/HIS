package dto

// ResultValueRequest represents a single result value
type ResultValueRequest struct {
	ParameterName string `json:"parameter_name" binding:"required"`
	Value         string `json:"value" binding:"required"`
	Remarks       string `json:"remarks" binding:"omitempty"`
}

// EnterLabTestResultsRequest represents request to enter multiple results
type EnterLabTestResultsRequest struct {
	Results []*ResultValueRequest `json:"results" binding:"required,min=1,dive"`
}

// UpdateLabTestResultRequest represents request to update a single result
type UpdateLabTestResultRequest struct {
	Value   string `json:"value" binding:"omitempty"`
	Remarks string `json:"remarks" binding:"omitempty"`
}
