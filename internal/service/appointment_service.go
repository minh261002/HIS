package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

var (
	ErrAppointmentNotFound     = errors.New("appointment not found")
	ErrTimeSlotNotAvailable    = errors.New("time slot not available")
	ErrInvalidAppointmentTime  = errors.New("appointment time must be during working hours (08:00-17:00)")
	ErrPastAppointmentDate     = errors.New("appointment date must be today or in the future")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// AppointmentService handles appointment business logic
type AppointmentService struct {
	appointmentRepo *repository.AppointmentRepository
	patientRepo     *repository.PatientRepository
	userRepo        *repository.UserRepository
}

// NewAppointmentService creates a new appointment service
func NewAppointmentService(
	appointmentRepo *repository.AppointmentRepository,
	patientRepo *repository.PatientRepository,
	userRepo *repository.UserRepository,
) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: appointmentRepo,
		patientRepo:     patientRepo,
		userRepo:        userRepo,
	}
}

// ScheduleAppointment schedules a new appointment
func (s *AppointmentService) ScheduleAppointment(req *dto.CreateAppointmentRequest, createdBy uint) (*dto.AppointmentResponse, error) {
	// Validate patient exists
	patient, err := s.patientRepo.FindByID(req.PatientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find patient: %w", err)
	}
	if patient == nil {
		return nil, ErrPatientNotFound
	}

	// Validate doctor exists
	doctor, err := s.userRepo.FindByID(req.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find doctor: %w", err)
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}

	// Parse date and time
	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	appointmentTime, err := time.Parse("15:04", req.AppointmentTime)
	if err != nil {
		return nil, errors.New("invalid time format, use HH:MM")
	}

	// Validate appointment date is not in the past
	today := time.Now().Truncate(24 * time.Hour)
	if appointmentDate.Before(today) {
		return nil, ErrPastAppointmentDate
	}

	// Validate working hours (8:00 - 17:00)
	hour := appointmentTime.Hour()
	if hour < 8 || hour >= 17 {
		return nil, ErrInvalidAppointmentTime
	}

	// Set default duration if not provided
	duration := req.DurationMinutes
	if duration == 0 {
		duration = 30
	}

	// Check time slot availability
	available, err := s.appointmentRepo.CheckTimeSlotAvailable(req.DoctorID, appointmentDate, appointmentTime, duration, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check time slot: %w", err)
	}
	if !available {
		return nil, ErrTimeSlotNotAvailable
	}

	// Generate appointment code
	code, err := s.appointmentRepo.GenerateAppointmentCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate appointment code: %w", err)
	}

	// Create appointment
	appointment := &domain.Appointment{
		AppointmentCode: code,
		PatientID:       req.PatientID,
		DoctorID:        req.DoctorID,
		AppointmentDate: appointmentDate,
		AppointmentTime: appointmentTime,
		DurationMinutes: duration,
		AppointmentType: domain.AppointmentType(req.AppointmentType),
		Status:          domain.AppointmentStatusScheduled,
		Reason:          req.Reason,
		Notes:           req.Notes,
		CreatedBy:       createdBy,
	}

	if err := s.appointmentRepo.Create(appointment); err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	// Reload to get relationships
	appointment, _ = s.appointmentRepo.FindByID(appointment.ID)
	return s.toAppointmentResponse(appointment), nil
}

// RescheduleAppointment reschedules an appointment
func (s *AppointmentService) RescheduleAppointment(id uint, req *dto.UpdateAppointmentRequest, updatedBy uint) (*dto.AppointmentResponse, error) {
	appointment, err := s.appointmentRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find appointment: %w", err)
	}
	if appointment == nil {
		return nil, ErrAppointmentNotFound
	}

	// Can't reschedule completed or cancelled appointments
	if appointment.Status == domain.AppointmentStatusCompleted || appointment.Status == domain.AppointmentStatusCancelled {
		return nil, errors.New("cannot reschedule completed or cancelled appointment")
	}

	// Update date if provided
	if req.AppointmentDate != "" {
		newDate, err := time.Parse("2006-01-02", req.AppointmentDate)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		today := time.Now().Truncate(24 * time.Hour)
		if newDate.Before(today) {
			return nil, ErrPastAppointmentDate
		}
		appointment.AppointmentDate = newDate
	}

	// Update time if provided
	if req.AppointmentTime != "" {
		newTime, err := time.Parse("15:04", req.AppointmentTime)
		if err != nil {
			return nil, errors.New("invalid time format, use HH:MM")
		}
		hour := newTime.Hour()
		if hour < 8 || hour >= 17 {
			return nil, ErrInvalidAppointmentTime
		}
		appointment.AppointmentTime = newTime
	}

	// Update duration if provided
	if req.DurationMinutes > 0 {
		appointment.DurationMinutes = req.DurationMinutes
	}

	// Check new time slot availability
	available, err := s.appointmentRepo.CheckTimeSlotAvailable(
		appointment.DoctorID,
		appointment.AppointmentDate,
		appointment.AppointmentTime,
		appointment.DurationMinutes,
		&appointment.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to check time slot: %w", err)
	}
	if !available {
		return nil, ErrTimeSlotNotAvailable
	}

	// Update other fields
	if req.Reason != "" {
		appointment.Reason = req.Reason
	}
	if req.Notes != "" {
		appointment.Notes = req.Notes
	}

	appointment.UpdatedBy = updatedBy

	if err := s.appointmentRepo.Update(appointment); err != nil {
		return nil, fmt.Errorf("failed to update appointment: %w", err)
	}

	appointment, _ = s.appointmentRepo.FindByID(appointment.ID)
	return s.toAppointmentResponse(appointment), nil
}

// CancelAppointment cancels an appointment
func (s *AppointmentService) CancelAppointment(id uint, reason string, cancelledBy uint) (*dto.AppointmentResponse, error) {
	appointment, err := s.appointmentRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find appointment: %w", err)
	}
	if appointment == nil {
		return nil, ErrAppointmentNotFound
	}

	if appointment.Status == domain.AppointmentStatusCompleted {
		return nil, errors.New("cannot cancel completed appointment")
	}
	if appointment.Status == domain.AppointmentStatusCancelled {
		return nil, errors.New("appointment already cancelled")
	}

	now := time.Now()
	appointment.Status = domain.AppointmentStatusCancelled
	appointment.CancelledReason = reason
	appointment.CancelledAt = &now
	appointment.CancelledBy = &cancelledBy
	appointment.UpdatedBy = cancelledBy

	if err := s.appointmentRepo.Update(appointment); err != nil {
		return nil, fmt.Errorf("failed to cancel appointment: %w", err)
	}

	appointment, _ = s.appointmentRepo.FindByID(appointment.ID)
	return s.toAppointmentResponse(appointment), nil
}

// ConfirmAppointment confirms an appointment
func (s *AppointmentService) ConfirmAppointment(id uint, updatedBy uint) error {
	return s.updateStatus(id, domain.AppointmentStatusConfirmed, updatedBy)
}

// StartAppointment starts an appointment
func (s *AppointmentService) StartAppointment(id uint, updatedBy uint) error {
	return s.updateStatus(id, domain.AppointmentStatusInProgress, updatedBy)
}

// CompleteAppointment completes an appointment
func (s *AppointmentService) CompleteAppointment(id uint, updatedBy uint) error {
	return s.updateStatus(id, domain.AppointmentStatusCompleted, updatedBy)
}

// MarkNoShow marks appointment as no-show
func (s *AppointmentService) MarkNoShow(id uint, updatedBy uint) error {
	return s.updateStatus(id, domain.AppointmentStatusNoShow, updatedBy)
}

// updateStatus updates appointment status
func (s *AppointmentService) updateStatus(id uint, newStatus domain.AppointmentStatus, updatedBy uint) error {
	appointment, err := s.appointmentRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find appointment: %w", err)
	}
	if appointment == nil {
		return ErrAppointmentNotFound
	}

	appointment.Status = newStatus
	appointment.UpdatedBy = updatedBy

	return s.appointmentRepo.Update(appointment)
}

// GetAppointmentByID gets appointment by ID
func (s *AppointmentService) GetAppointmentByID(id uint) (*dto.AppointmentResponse, error) {
	appointment, err := s.appointmentRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find appointment: %w", err)
	}
	if appointment == nil {
		return nil, ErrAppointmentNotFound
	}
	return s.toAppointmentResponse(appointment), nil
}

// GetAppointmentByCode gets appointment by code
func (s *AppointmentService) GetAppointmentByCode(code string) (*dto.AppointmentResponse, error) {
	appointment, err := s.appointmentRepo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find appointment: %w", err)
	}
	if appointment == nil {
		return nil, ErrAppointmentNotFound
	}
	return s.toAppointmentResponse(appointment), nil
}

// GetPatientAppointments gets patient's appointments
func (s *AppointmentService) GetPatientAppointments(patientID uint, filters map[string]interface{}) ([]*dto.AppointmentListItem, error) {
	appointments, err := s.appointmentRepo.FindByPatientID(patientID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %w", err)
	}

	items := make([]*dto.AppointmentListItem, len(appointments))
	for i, apt := range appointments {
		items[i] = s.toAppointmentListItem(apt)
	}
	return items, nil
}

// GetDoctorSchedule gets doctor's schedule for a date
func (s *AppointmentService) GetDoctorSchedule(doctorID uint, dateStr string) ([]*dto.AppointmentListItem, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	appointments, err := s.appointmentRepo.FindByDoctorID(doctorID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	items := make([]*dto.AppointmentListItem, len(appointments))
	for i, apt := range appointments {
		items[i] = s.toAppointmentListItem(apt)
	}
	return items, nil
}

// GetAvailableTimeSlots gets available time slots for a doctor on a date
func (s *AppointmentService) GetAvailableTimeSlots(doctorID uint, dateStr string, duration int) ([]*dto.TimeSlot, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	if duration == 0 {
		duration = 30
	}

	slots := []*dto.TimeSlot{}
	startHour := 8
	endHour := 17

	for hour := startHour; hour < endHour; hour++ {
		for minute := 0; minute < 60; minute += 30 {
			timeStr := fmt.Sprintf("%02d:%02d", hour, minute)
			slotTime, _ := time.Parse("15:04", timeStr)

			available, _ := s.appointmentRepo.CheckTimeSlotAvailable(doctorID, date, slotTime, duration, nil)
			slots = append(slots, &dto.TimeSlot{
				Time:      timeStr,
				Available: available,
			})
		}
	}

	return slots, nil
}

// SearchAppointments searches appointments
func (s *AppointmentService) SearchAppointments(filters map[string]interface{}, page, pageSize int) ([]*dto.AppointmentListItem, int64, error) {
	appointments, total, err := s.appointmentRepo.Search(filters, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search appointments: %w", err)
	}

	items := make([]*dto.AppointmentListItem, len(appointments))
	for i, apt := range appointments {
		items[i] = s.toAppointmentListItem(apt)
	}
	return items, total, nil
}

// GetUpcomingAppointments gets upcoming appointments
func (s *AppointmentService) GetUpcomingAppointments(limit int) ([]*dto.AppointmentListItem, error) {
	appointments, err := s.appointmentRepo.GetUpcomingAppointments(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming appointments: %w", err)
	}

	items := make([]*dto.AppointmentListItem, len(appointments))
	for i, apt := range appointments {
		items[i] = s.toAppointmentListItem(apt)
	}
	return items, nil
}

// Helper functions
func (s *AppointmentService) toAppointmentResponse(apt *domain.Appointment) *dto.AppointmentResponse {
	resp := &dto.AppointmentResponse{
		ID:              apt.ID,
		AppointmentCode: apt.AppointmentCode,
		PatientID:       apt.PatientID,
		DoctorID:        apt.DoctorID,
		AppointmentDate: apt.AppointmentDate.Format("2006-01-02"),
		AppointmentTime: apt.AppointmentTime.Format("15:04"),
		DurationMinutes: apt.DurationMinutes,
		AppointmentType: string(apt.AppointmentType),
		Status:          string(apt.Status),
		Reason:          apt.Reason,
		Notes:           apt.Notes,
		CancelledReason: apt.CancelledReason,
		CancelledAt:     apt.CancelledAt,
		CreatedAt:       apt.CreatedAt,
		UpdatedAt:       apt.UpdatedAt,
	}

	if apt.Patient != nil {
		resp.PatientName = apt.Patient.FullName
	}
	if apt.Doctor != nil {
		resp.DoctorName = apt.Doctor.FullName
	}

	return resp
}

func (s *AppointmentService) toAppointmentListItem(apt *domain.Appointment) *dto.AppointmentListItem {
	item := &dto.AppointmentListItem{
		ID:              apt.ID,
		AppointmentCode: apt.AppointmentCode,
		AppointmentDate: apt.AppointmentDate.Format("2006-01-02"),
		AppointmentTime: apt.AppointmentTime.Format("15:04"),
		AppointmentType: string(apt.AppointmentType),
		Status:          string(apt.Status),
	}

	if apt.Patient != nil {
		item.PatientName = apt.Patient.FullName
	}
	if apt.Doctor != nil {
		item.DoctorName = apt.Doctor.FullName
	}

	return item
}
