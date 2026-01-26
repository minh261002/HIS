package repository

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// MedicationRepository handles medication data operations
type MedicationRepository struct {
	db *gorm.DB
}

// NewMedicationRepository creates a new medication repository
func NewMedicationRepository(db *gorm.DB) *MedicationRepository {
	return &MedicationRepository{db: db}
}

// FindByID finds a medication by ID
func (r *MedicationRepository) FindByID(id uint) (*domain.Medication, error) {
	var medication domain.Medication
	err := r.db.Where("is_active = ?", true).First(&medication, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &medication, nil
}

// Search searches medications by name or generic name
func (r *MedicationRepository) Search(query string, limit int) ([]*domain.Medication, error) {
	var medications []*domain.Medication

	if limit == 0 {
		limit = 20
	}

	err := r.db.Where("is_active = ? AND (name LIKE ? OR generic_name LIKE ?)",
		true, "%"+query+"%", "%"+query+"%").
		Order("name ASC").
		Limit(limit).
		Find(&medications).Error

	return medications, err
}

// FindByDosageForm finds medications by dosage form
func (r *MedicationRepository) FindByDosageForm(form string) ([]*domain.Medication, error) {
	var medications []*domain.Medication
	err := r.db.Where("dosage_form = ? AND is_active = ?", form, true).
		Order("name ASC").
		Find(&medications).Error
	return medications, err
}
