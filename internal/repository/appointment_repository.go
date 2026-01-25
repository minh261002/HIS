package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"gorm.io/gorm"
)

// AppointmentRepository handles appointment data operations
type AppointmentRepository struct {
	db *gorm.DB
}

// NewAppointmentRepository creates a new appointment repository
func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

// Create creates a new appointment
func (r *AppointmentRepository) Create(appointment *domain.Appointment) error {
	return r.db.Create(appointment).Error
}

// FindByID finds an appointment by ID
func (r *AppointmentRepository) FindByID(id uint) (*domain.Appointment, error) {
	var appointment domain.Appointment
	err := r.db.Preload("Patient").Preload("Doctor").First(&appointment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &appointment, nil
}

// FindByCode finds an appointment by appointment code
func (r *AppointmentRepository) FindByCode(code string) (*domain.Appointment, error) {
	var appointment domain.Appointment
	err := r.db.Preload("Patient").Preload("Doctor").Where("appointment_code = ?", code).First(&appointment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &appointment, nil
}

// FindByPatientID finds appointments for a patient
func (r *AppointmentRepository) FindByPatientID(patientID uint, filters map[string]interface{}) ([]*domain.Appointment, error) {
	var appointments []*domain.Appointment
	query := r.db.Preload("Doctor").Where("patient_id = ?", patientID)

	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate, ok := filters["from_date"]; ok && fromDate != "" {
		query = query.Where("appointment_date >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"]; ok && toDate != "" {
		query = query.Where("appointment_date <= ?", toDate)
	}

	err := query.Order("appointment_date DESC, appointment_time DESC").Find(&appointments).Error
	return appointments, err
}

// FindByDoctorID finds appointments for a doctor on a specific date
func (r *AppointmentRepository) FindByDoctorID(doctorID uint, date time.Time) ([]*domain.Appointment, error) {
	var appointments []*domain.Appointment
	err := r.db.Preload("Patient").
		Where("doctor_id = ? AND appointment_date = ?", doctorID, date.Format("2006-01-02")).
		Order("appointment_time ASC").
		Find(&appointments).Error
	return appointments, err
}

// CheckTimeSlotAvailable checks if a time slot is available for a doctor
func (r *AppointmentRepository) CheckTimeSlotAvailable(doctorID uint, date time.Time, appointmentTime time.Time, duration int, excludeID *uint) (bool, error) {
	dateStr := date.Format("2006-01-02")
	timeStr := appointmentTime.Format("15:04:05")

	// Calculate end time
	endTime := appointmentTime.Add(time.Duration(duration) * time.Minute)
	endTimeStr := endTime.Format("15:04:05")

	query := r.db.Model(&domain.Appointment{}).
		Where("doctor_id = ?", doctorID).
		Where("appointment_date = ?", dateStr).
		Where("status NOT IN ?", []string{"CANCELLED", "NO_SHOW"}).
		Where("(appointment_time < ? AND ADDTIME(appointment_time, SEC_TO_TIME(duration_minutes * 60)) > ?) OR (appointment_time >= ? AND appointment_time < ?)",
			endTimeStr, timeStr, timeStr, endTimeStr)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// Update updates an appointment
func (r *AppointmentRepository) Update(appointment *domain.Appointment) error {
	return r.db.Save(appointment).Error
}

// Delete soft deletes an appointment
func (r *AppointmentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Appointment{}, id).Error
}

// GetUpcomingAppointments gets upcoming appointments
func (r *AppointmentRepository) GetUpcomingAppointments(limit int) ([]*domain.Appointment, error) {
	var appointments []*domain.Appointment
	now := time.Now()
	today := now.Format("2006-01-02")
	currentTime := now.Format("15:04:05")

	err := r.db.Preload("Patient").Preload("Doctor").
		Where("(appointment_date > ? OR (appointment_date = ? AND appointment_time >= ?))", today, today, currentTime).
		Where("status IN ?", []string{"SCHEDULED", "CONFIRMED"}).
		Order("appointment_date ASC, appointment_time ASC").
		Limit(limit).
		Find(&appointments).Error

	return appointments, err
}

// Search searches appointments with filters
func (r *AppointmentRepository) Search(filters map[string]interface{}, page, pageSize int) ([]*domain.Appointment, int64, error) {
	var appointments []*domain.Appointment
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&domain.Appointment{}).Preload("Patient").Preload("Doctor")

	if patientID, ok := filters["patient_id"]; ok && patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}
	if doctorID, ok := filters["doctor_id"]; ok && doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if appointmentType, ok := filters["appointment_type"]; ok && appointmentType != "" {
		query = query.Where("appointment_type = ?", appointmentType)
	}
	if fromDate, ok := filters["from_date"]; ok && fromDate != "" {
		query = query.Where("appointment_date >= ?", fromDate)
	}
	if toDate, ok := filters["to_date"]; ok && toDate != "" {
		query = query.Where("appointment_date <= ?", toDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).
		Limit(pageSize).
		Order("appointment_date DESC, appointment_time DESC").
		Find(&appointments).Error

	return appointments, total, err
}

// GenerateAppointmentCode generates a unique appointment code
func (r *AppointmentRepository) GenerateAppointmentCode() (string, error) {
	today := time.Now().Format("20060102") // YYYYMMDD
	prefix := fmt.Sprintf("APT-%s-", today)

	var lastAppointment domain.Appointment
	err := r.db.Where("appointment_code LIKE ?", prefix+"%").
		Order("appointment_code DESC").
		First(&lastAppointment).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	sequence := 1
	if lastAppointment.AppointmentCode != "" {
		var lastSeq int
		fmt.Sscanf(lastAppointment.AppointmentCode, prefix+"%d", &lastSeq)
		sequence = lastSeq + 1
	}

	return fmt.Sprintf("%s%04d", prefix, sequence), nil
}
