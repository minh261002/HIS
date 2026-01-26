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
	ErrImagingTemplateNotFound = errors.New("imaging template not found")
	ErrImagingRequestNotFound  = errors.New("imaging request not found")
)

// ImagingRequestService handles imaging request business logic
type ImagingRequestService struct {
	requestRepo  *repository.ImagingRequestRepository
	resultRepo   *repository.ImagingResultRepository
	templateRepo *repository.ImagingTemplateRepository
	visitRepo    *repository.VisitRepository
}

// NewImagingRequestService creates a new imaging request service
func NewImagingRequestService(
	requestRepo *repository.ImagingRequestRepository,
	resultRepo *repository.ImagingResultRepository,
	templateRepo *repository.ImagingTemplateRepository,
	visitRepo *repository.VisitRepository,
) *ImagingRequestService {
	return &ImagingRequestService{
		requestRepo:  requestRepo,
		resultRepo:   resultRepo,
		templateRepo: templateRepo,
		visitRepo:    visitRepo,
	}
}

// CreateImagingRequest creates an imaging request
func (s *ImagingRequestService) CreateImagingRequest(req *dto.CreateImagingRequestRequest, requestedBy uint) (*dto.ImagingRequestResponse, error) {
	// Validate visit exists
	visit, err := s.visitRepo.FindByID(req.VisitID)
	if err != nil {
		return nil, fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}

	// Validate template exists
	template, err := s.templateRepo.FindByID(req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("failed to find template: %w", err)
	}
	if template == nil {
		return nil, ErrImagingTemplateNotFound
	}

	// Generate request code
	code, err := s.requestRepo.GenerateRequestCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate request code: %w", err)
	}

	// Create request
	request := &domain.ImagingRequest{
		RequestCode:         code,
		VisitID:             req.VisitID,
		PatientID:           visit.PatientID,
		DoctorID:            visit.DoctorID,
		TemplateID:          req.TemplateID,
		Status:              domain.ImagingRequestStatusPending,
		Priority:            domain.ImagingPriority(req.Priority),
		RequestedDate:       time.Now(),
		ClinicalIndication:  req.ClinicalIndication,
		SpecialInstructions: req.SpecialInstructions,
		CreatedBy:           requestedBy,
	}

	if err := s.requestRepo.Create(request); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Reload to get relationships
	request, _ = s.requestRepo.FindByID(request.ID)
	return s.toImagingRequestResponse(request), nil
}

// ScheduleImaging schedules an imaging request
func (s *ImagingRequestService) ScheduleImaging(id uint, scheduledDate time.Time, updatedBy uint) error {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return ErrImagingRequestNotFound
	}

	request.Status = domain.ImagingRequestStatusScheduled
	request.ScheduledDate = &scheduledDate
	request.UpdatedBy = updatedBy
	return s.requestRepo.Update(request)
}

// StartImaging marks imaging as in progress
func (s *ImagingRequestService) StartImaging(id uint, updatedBy uint) error {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return ErrImagingRequestNotFound
	}

	request.Status = domain.ImagingRequestStatusInProgress
	request.UpdatedBy = updatedBy
	return s.requestRepo.Update(request)
}

// CompleteImaging marks imaging as completed
func (s *ImagingRequestService) CompleteImaging(id uint, updatedBy uint) error {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return ErrImagingRequestNotFound
	}

	now := time.Now()
	request.Status = domain.ImagingRequestStatusCompleted
	request.CompletedAt = &now
	request.UpdatedBy = updatedBy
	return s.requestRepo.Update(request)
}

// CancelImaging cancels an imaging request
func (s *ImagingRequestService) CancelImaging(id uint, updatedBy uint) error {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return ErrImagingRequestNotFound
	}

	request.Status = domain.ImagingRequestStatusCancelled
	request.UpdatedBy = updatedBy
	return s.requestRepo.Update(request)
}

// CreateOrUpdateResult creates or updates imaging result
func (s *ImagingRequestService) CreateOrUpdateResult(requestID uint, req *dto.CreateImagingResultRequest, radiologistID uint) error {
	// Check if result exists
	existingResult, err := s.resultRepo.FindByRequestID(requestID)
	if err != nil {
		return fmt.Errorf("failed to check existing result: %w", err)
	}

	if existingResult != nil {
		// Update existing result
		existingResult.Findings = req.Findings
		existingResult.Impression = req.Impression
		existingResult.DICOMFiles = req.DICOMFiles
		existingResult.IsCritical = req.IsCritical
		existingResult.ReportDate = time.Now()
		existingResult.RadiologistID = radiologistID
		return s.resultRepo.Update(existingResult)
	}

	// Create new result
	result := &domain.ImagingResult{
		RequestID:     requestID,
		RadiologistID: radiologistID,
		Findings:      req.Findings,
		Impression:    req.Impression,
		DICOMFiles:    req.DICOMFiles,
		ReportDate:    time.Now(),
		IsCritical:    req.IsCritical,
	}

	return s.resultRepo.Create(result)
}

// GetImagingRequestByID gets request by ID
func (s *ImagingRequestService) GetImagingRequestByID(id uint) (*dto.ImagingRequestResponse, error) {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return nil, ErrImagingRequestNotFound
	}
	return s.toImagingRequestResponse(request), nil
}

// GetImagingRequestByCode gets request by code
func (s *ImagingRequestService) GetImagingRequestByCode(code string) (*dto.ImagingRequestResponse, error) {
	request, err := s.requestRepo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return nil, ErrImagingRequestNotFound
	}
	return s.toImagingRequestResponse(request), nil
}

// GetVisitImagingRequests gets visit's imaging requests
func (s *ImagingRequestService) GetVisitImagingRequests(visitID uint) ([]*dto.ImagingRequestListItem, error) {
	requests, err := s.requestRepo.FindByVisitID(visitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requests: %w", err)
	}

	items := make([]*dto.ImagingRequestListItem, len(requests))
	for i, r := range requests {
		items[i] = s.toImagingRequestListItem(r)
	}
	return items, nil
}

// GetPatientImagingRequests gets patient's imaging request history
func (s *ImagingRequestService) GetPatientImagingRequests(patientID uint, filters map[string]interface{}) ([]*dto.ImagingRequestListItem, error) {
	requests, err := s.requestRepo.FindByPatientID(patientID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get requests: %w", err)
	}

	items := make([]*dto.ImagingRequestListItem, len(requests))
	for i, r := range requests {
		items[i] = s.toImagingRequestListItem(r)
	}
	return items, nil
}

// Helper functions
func (s *ImagingRequestService) toImagingRequestResponse(r *domain.ImagingRequest) *dto.ImagingRequestResponse {
	resp := &dto.ImagingRequestResponse{
		ID:                  r.ID,
		RequestCode:         r.RequestCode,
		VisitID:             r.VisitID,
		PatientID:           r.PatientID,
		DoctorID:            r.DoctorID,
		TemplateID:          r.TemplateID,
		Status:              string(r.Status),
		Priority:            string(r.Priority),
		RequestedDate:       r.RequestedDate,
		ScheduledDate:       r.ScheduledDate,
		CompletedAt:         r.CompletedAt,
		ClinicalIndication:  r.ClinicalIndication,
		SpecialInstructions: r.SpecialInstructions,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
	}

	if r.Patient != nil {
		resp.PatientName = r.Patient.FullName
	}
	if r.Doctor != nil {
		resp.DoctorName = r.Doctor.FullName
	}
	if r.Template != nil {
		resp.TemplateName = r.Template.Name
		resp.TemplateCode = r.Template.Code
		resp.Modality = string(r.Template.Modality)
	}

	// Add result if exists
	if r.Result != nil {
		resp.Result = &dto.ImagingResultResponse{
			ID:            r.Result.ID,
			RadiologistID: r.Result.RadiologistID,
			Findings:      r.Result.Findings,
			Impression:    r.Result.Impression,
			DICOMFiles:    r.Result.DICOMFiles,
			ReportDate:    r.Result.ReportDate.Format("2006-01-02"),
			IsCritical:    r.Result.IsCritical,
			CreatedAt:     r.Result.CreatedAt,
		}
		if r.Result.Radiologist != nil {
			resp.Result.RadiologistName = r.Result.Radiologist.FullName
		}
	}

	return resp
}

func (s *ImagingRequestService) toImagingRequestListItem(r *domain.ImagingRequest) *dto.ImagingRequestListItem {
	item := &dto.ImagingRequestListItem{
		ID:            r.ID,
		RequestCode:   r.RequestCode,
		Status:        string(r.Status),
		Priority:      string(r.Priority),
		RequestedDate: r.RequestedDate,
		HasResult:     r.Result != nil,
	}

	if r.Patient != nil {
		item.PatientName = r.Patient.FullName
	}
	if r.Template != nil {
		item.TemplateName = r.Template.Name
		item.Modality = string(r.Template.Modality)
	}

	return item
}
