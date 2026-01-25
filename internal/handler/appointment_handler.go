package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/middleware"
	"github.com/minhtran/his/internal/pkg/response"
	"github.com/minhtran/his/internal/service"
)

// AppointmentHandler handles appointment HTTP requests
type AppointmentHandler struct {
	appointmentService *service.AppointmentService
}

// NewAppointmentHandler creates a new appointment handler
func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

// ScheduleAppointment handles scheduling a new appointment
func (h *AppointmentHandler) ScheduleAppointment(c *gin.Context) {
	var req dto.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	appointment, err := h.appointmentService.ScheduleAppointment(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrPatientNotFound) {
			response.NotFound(c, "Patient not found")
			return
		}
		if errors.Is(err, service.ErrTimeSlotNotAvailable) {
			response.BadRequest(c, "Time slot not available", nil)
			return
		}
		if errors.Is(err, service.ErrInvalidAppointmentTime) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		if errors.Is(err, service.ErrPastAppointmentDate) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		if errors.Is(err, service.ErrInvalidDateFormat) {
			response.BadRequest(c, "Invalid date format, use YYYY-MM-DD", nil)
			return
		}
		response.InternalServerError(c, "Failed to schedule appointment")
		return
	}

	response.Created(c, "Appointment scheduled successfully", appointment)
}

// ListAppointments handles listing appointments with filters
func (h *AppointmentHandler) ListAppointments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	filters := make(map[string]interface{})
	if patientID := c.Query("patient_id"); patientID != "" {
		filters["patient_id"] = patientID
	}
	if doctorID := c.Query("doctor_id"); doctorID != "" {
		filters["doctor_id"] = doctorID
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if appointmentType := c.Query("appointment_type"); appointmentType != "" {
		filters["appointment_type"] = appointmentType
	}
	if fromDate := c.Query("from_date"); fromDate != "" {
		filters["from_date"] = fromDate
	}
	if toDate := c.Query("to_date"); toDate != "" {
		filters["to_date"] = toDate
	}

	appointments, total, err := h.appointmentService.SearchAppointments(filters, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to list appointments")
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response.SuccessPaginated(c, "Appointments retrieved successfully", appointments, response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}

// GetUpcomingAppointments handles getting upcoming appointments
func (h *AppointmentHandler) GetUpcomingAppointments(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	appointments, err := h.appointmentService.GetUpcomingAppointments(limit)
	if err != nil {
		response.InternalServerError(c, "Failed to get upcoming appointments")
		return
	}

	response.Success(c, "Upcoming appointments retrieved successfully", appointments)
}

// GetAppointment handles getting appointment details by ID
func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appointment ID", nil)
		return
	}

	appointment, err := h.appointmentService.GetAppointmentByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrAppointmentNotFound) {
			response.NotFound(c, "Appointment not found")
			return
		}
		response.InternalServerError(c, "Failed to get appointment")
		return
	}

	response.Success(c, "Appointment retrieved successfully", appointment)
}

// GetAppointmentByCode handles getting appointment by code
func (h *AppointmentHandler) GetAppointmentByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.BadRequest(c, "Appointment code is required", nil)
		return
	}

	appointment, err := h.appointmentService.GetAppointmentByCode(code)
	if err != nil {
		if errors.Is(err, service.ErrAppointmentNotFound) {
			response.NotFound(c, "Appointment not found")
			return
		}
		response.InternalServerError(c, "Failed to get appointment")
		return
	}

	response.Success(c, "Appointment retrieved successfully", appointment)
}

// RescheduleAppointment handles rescheduling an appointment
func (h *AppointmentHandler) RescheduleAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appointment ID", nil)
		return
	}

	var req dto.UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	appointment, err := h.appointmentService.RescheduleAppointment(uint(id), &req, userID)
	if err != nil {
		if errors.Is(err, service.ErrAppointmentNotFound) {
			response.NotFound(c, "Appointment not found")
			return
		}
		if errors.Is(err, service.ErrTimeSlotNotAvailable) {
			response.BadRequest(c, "Time slot not available", nil)
			return
		}
		if errors.Is(err, service.ErrInvalidAppointmentTime) || errors.Is(err, service.ErrPastAppointmentDate) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Failed to reschedule appointment")
		return
	}

	response.Success(c, "Appointment rescheduled successfully", appointment)
}

// CancelAppointment handles cancelling an appointment
func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appointment ID", nil)
		return
	}

	var req dto.CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]interface{}{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)

	appointment, err := h.appointmentService.CancelAppointment(uint(id), req.Reason, userID)
	if err != nil {
		if errors.Is(err, service.ErrAppointmentNotFound) {
			response.NotFound(c, "Appointment not found")
			return
		}
		response.InternalServerError(c, "Failed to cancel appointment")
		return
	}

	response.Success(c, "Appointment cancelled successfully", appointment)
}

// ConfirmAppointment handles confirming an appointment
func (h *AppointmentHandler) ConfirmAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appointment ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.appointmentService.ConfirmAppointment(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrAppointmentNotFound) {
			response.NotFound(c, "Appointment not found")
			return
		}
		response.InternalServerError(c, "Failed to confirm appointment")
		return
	}

	response.Success(c, "Appointment confirmed successfully", nil)
}

// StartAppointment handles starting an appointment
func (h *AppointmentHandler) StartAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appointment ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.appointmentService.StartAppointment(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrAppointmentNotFound) {
			response.NotFound(c, "Appointment not found")
			return
		}
		response.InternalServerError(c, "Failed to start appointment")
		return
	}

	response.Success(c, "Appointment started successfully", nil)
}

// CompleteAppointment handles completing an appointment
func (h *AppointmentHandler) CompleteAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appointment ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.appointmentService.CompleteAppointment(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrAppointmentNotFound) {
			response.NotFound(c, "Appointment not found")
			return
		}
		response.InternalServerError(c, "Failed to complete appointment")
		return
	}

	response.Success(c, "Appointment completed successfully", nil)
}

// MarkNoShow handles marking appointment as no-show
func (h *AppointmentHandler) MarkNoShow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid appointment ID", nil)
		return
	}

	userID, _ := middleware.GetUserID(c)

	if err := h.appointmentService.MarkNoShow(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrAppointmentNotFound) {
			response.NotFound(c, "Appointment not found")
			return
		}
		response.InternalServerError(c, "Failed to mark as no-show")
		return
	}

	response.Success(c, "Appointment marked as no-show", nil)
}

// GetPatientAppointments handles getting patient's appointments
func (h *AppointmentHandler) GetPatientAppointments(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid patient ID", nil)
		return
	}

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if fromDate := c.Query("from_date"); fromDate != "" {
		filters["from_date"] = fromDate
	}
	if toDate := c.Query("to_date"); toDate != "" {
		filters["to_date"] = toDate
	}

	appointments, err := h.appointmentService.GetPatientAppointments(uint(patientID), filters)
	if err != nil {
		response.InternalServerError(c, "Failed to get patient appointments")
		return
	}

	response.Success(c, "Patient appointments retrieved successfully", appointments)
}

// GetDoctorSchedule handles getting doctor's schedule
func (h *AppointmentHandler) GetDoctorSchedule(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid doctor ID", nil)
		return
	}

	date := c.Query("date")
	if date == "" {
		response.BadRequest(c, "Date parameter is required (YYYY-MM-DD)", nil)
		return
	}

	appointments, err := h.appointmentService.GetDoctorSchedule(uint(doctorID), date)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			response.BadRequest(c, "Invalid date format, use YYYY-MM-DD", nil)
			return
		}
		response.InternalServerError(c, "Failed to get doctor schedule")
		return
	}

	response.Success(c, "Doctor schedule retrieved successfully", appointments)
}

// GetAvailableTimeSlots handles getting available time slots
func (h *AppointmentHandler) GetAvailableTimeSlots(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid doctor ID", nil)
		return
	}

	date := c.Query("date")
	if date == "" {
		response.BadRequest(c, "Date parameter is required (YYYY-MM-DD)", nil)
		return
	}

	duration, _ := strconv.Atoi(c.DefaultQuery("duration", "30"))

	slots, err := h.appointmentService.GetAvailableTimeSlots(uint(doctorID), date, duration)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateFormat) {
			response.BadRequest(c, "Invalid date format, use YYYY-MM-DD", nil)
			return
		}
		response.InternalServerError(c, "Failed to get available time slots")
		return
	}

	response.Success(c, "Available time slots retrieved successfully", slots)
}
