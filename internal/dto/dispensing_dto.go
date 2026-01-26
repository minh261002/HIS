package dto

import "time"

// DispensingItem represents an item to dispense
type DispensingItem struct {
	PrescriptionItemID uint `json:"prescription_item_id" binding:"required"`
	InventoryID        uint `json:"inventory_id" binding:"required"`
	Quantity           int  `json:"quantity" binding:"required,min=1"`
}

// DispensePrescriptionRequest represents request to dispense prescription
type DispensePrescriptionRequest struct {
	PrescriptionID uint              `json:"prescription_id" binding:"required"`
	Items          []*DispensingItem `json:"items" binding:"required,min=1,dive"`
}

// DispensingResponse represents dispensing details
type DispensingResponse struct {
	ID                 uint      `json:"id"`
	DispensingCode     string    `json:"dispensing_code"`
	PrescriptionID     uint      `json:"prescription_id"`
	PrescriptionItemID uint      `json:"prescription_item_id"`
	MedicationID       uint      `json:"medication_id"`
	MedicationName     string    `json:"medication_name"`
	PatientID          uint      `json:"patient_id"`
	PatientName        string    `json:"patient_name"`
	PharmacistID       uint      `json:"pharmacist_id"`
	PharmacistName     string    `json:"pharmacist_name"`
	QuantityDispensed  int       `json:"quantity_dispensed"`
	BatchNumber        string    `json:"batch_number"`
	DispensedDate      time.Time `json:"dispensed_date"`
	Notes              string    `json:"notes"`
	CreatedAt          time.Time `json:"created_at"`
}

// DispensingListItem represents simplified dispensing for list view
type DispensingListItem struct {
	ID                uint      `json:"id"`
	DispensingCode    string    `json:"dispensing_code"`
	MedicationName    string    `json:"medication_name"`
	QuantityDispensed int       `json:"quantity_dispensed"`
	DispensedDate     time.Time `json:"dispensed_date"`
	PharmacistName    string    `json:"pharmacist_name"`
}
