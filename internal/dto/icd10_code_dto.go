package dto

// ICD10CodeResponse represents ICD-10 code details
type ICD10CodeResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// ICD10CodeListItem represents simplified ICD-10 code for search results
type ICD10CodeListItem struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Category    string `json:"category"`
}
