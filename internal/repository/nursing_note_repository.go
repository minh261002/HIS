package repository

import (
	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// NursingNoteRepository handles nursing note data operations
type NursingNoteRepository struct {
	db *gorm.DB
}

// NewNursingNoteRepository creates a new nursing note repository
func NewNursingNoteRepository(db *gorm.DB) *NursingNoteRepository {
	return &NursingNoteRepository{db: db}
}

// Create creates a nursing note
func (r *NursingNoteRepository) Create(note *domain.NursingNote) error {
	return r.db.Create(note).Error
}

// FindByAdmissionID finds notes for an admission
func (r *NursingNoteRepository) FindByAdmissionID(admissionID uint) ([]*domain.NursingNote, error) {
	var notes []*domain.NursingNote
	err := r.db.Preload("Nurse").
		Where("admission_id = ?", admissionID).
		Order("note_date DESC").
		Find(&notes).Error
	return notes, err
}

// Update updates a nursing note
func (r *NursingNoteRepository) Update(note *domain.NursingNote) error {
	return r.db.Save(note).Error
}
