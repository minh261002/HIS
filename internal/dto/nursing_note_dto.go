package dto

import (
	"time"

	"github.com/minhtran/his/internal/domain"
)

// CreateNursingNoteRequest represents request to create a nursing note
type CreateNursingNoteRequest struct {
	VitalSigns    domain.VitalSigns `json:"vital_signs" binding:"required"`
	Observations  string            `json:"observations" binding:"required"`
	Interventions string            `json:"interventions" binding:"omitempty"`
}

// NursingNoteResponse represents nursing note details
type NursingNoteResponse struct {
	ID            uint              `json:"id"`
	NurseID       uint              `json:"nurse_id"`
	NurseName     string            `json:"nurse_name"`
	NoteDate      time.Time         `json:"note_date"`
	VitalSigns    domain.VitalSigns `json:"vital_signs"`
	Observations  string            `json:"observations"`
	Interventions string            `json:"interventions"`
	CreatedAt     time.Time         `json:"created_at"`
}
