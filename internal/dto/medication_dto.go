package dto

// MedicationResponse represents medication details
type MedicationResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	GenericName  string `json:"generic_name"`
	DosageForm   string `json:"dosage_form"`
	Strength     string `json:"strength"`
	Unit         string `json:"unit"`
	Manufacturer string `json:"manufacturer"`
}

// MedicationListItem represents simplified medication for search results
type MedicationListItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	GenericName string `json:"generic_name"`
	DosageForm  string `json:"dosage_form"`
	Strength    string `json:"strength"`
}
