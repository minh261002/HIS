package dto

import "time"

// CreateAdmissionRequest represents request to create an admission
type CreateAdmissionRequest struct {
	VisitID            uint   `json:"visit_id" binding:"required"`
	AdmissionDiagnosis string `json:"admission_diagnosis" binding:"required"`
	BedID              *uint  `json:"bed_id" binding:"omitempty"` // Optional initial bed
}

// DischargeAdmissionRequest represents request to discharge a patient
type DischargeAdmissionRequest struct {
	DischargeDiagnosis string `json:"discharge_diagnosis" binding:"required"`
	DischargeSummary   string `json:"discharge_summary" binding:"required"`
}

// TransferBedRequest represents request to transfer bed
type TransferBedRequest struct {
	NewBedID uint   `json:"new_bed_id" binding:"required"`
	Notes    string `json:"notes" binding:"omitempty"`
}

// BedAllocationResponse represents bed allocation details
type BedAllocationResponse struct {
	ID            uint       `json:"id"`
	BedNumber     string     `json:"bed_number"`
	Ward          string     `json:"ward"`
	AllocatedDate time.Time  `json:"allocated_date"`
	ReleasedDate  *time.Time `json:"released_date,omitempty"`
	IsCurrent     bool       `json:"is_current"`
	Notes         string     `json:"notes"`
}

// AdmissionResponse represents admission details
type AdmissionResponse struct {
	ID                 uint                     `json:"id"`
	AdmissionCode      string                   `json:"admission_code"`
	VisitID            uint                     `json:"visit_id"`
	PatientID          uint                     `json:"patient_id"`
	PatientName        string                   `json:"patient_name"`
	DoctorID           uint                     `json:"doctor_id"`
	DoctorName         string                   `json:"doctor_name"`
	AdmissionDate      time.Time                `json:"admission_date"`
	DischargeDate      *time.Time               `json:"discharge_date,omitempty"`
	AdmissionDiagnosis string                   `json:"admission_diagnosis"`
	DischargeDiagnosis string                   `json:"discharge_diagnosis"`
	DischargeSummary   string                   `json:"discharge_summary"`
	Status             string                   `json:"status"`
	BedAllocations     []*BedAllocationResponse `json:"bed_allocations"`
	CreatedAt          time.Time                `json:"created_at"`
	UpdatedAt          time.Time                `json:"updated_at"`
}

// AdmissionListItem represents simplified admission for list view
type AdmissionListItem struct {
	ID               uint      `json:"id"`
	AdmissionCode    string    `json:"admission_code"`
	PatientName      string    `json:"patient_name"`
	CurrentBedNumber string    `json:"current_bed_number"`
	AdmissionDate    time.Time `json:"admission_date"`
	Status           string    `json:"status"`
}
