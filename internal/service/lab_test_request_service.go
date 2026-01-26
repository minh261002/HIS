package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/minhtran/his/internal/domain"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

var (
	ErrLabTestTemplateNotFound = errors.New("lab test template not found")
	ErrLabTestRequestNotFound  = errors.New("lab test request not found")
)

// LabTestRequestService handles lab test request business logic
type LabTestRequestService struct {
	requestRepo  *repository.LabTestRequestRepository
	resultRepo   *repository.LabTestResultRepository
	templateRepo *repository.LabTestTemplateRepository
	visitRepo    *repository.VisitRepository
}

// NewLabTestRequestService creates a new lab test request service
func NewLabTestRequestService(
	requestRepo *repository.LabTestRequestRepository,
	resultRepo *repository.LabTestResultRepository,
	templateRepo *repository.LabTestTemplateRepository,
	visitRepo *repository.VisitRepository,
) *LabTestRequestService {
	return &LabTestRequestService{
		requestRepo:  requestRepo,
		resultRepo:   resultRepo,
		templateRepo: templateRepo,
		visitRepo:    visitRepo,
	}
}

// CreateLabTestRequest creates a lab test request with auto-generated result placeholders
func (s *LabTestRequestService) CreateLabTestRequest(req *dto.CreateLabTestRequestRequest, requestedBy uint) (*dto.LabTestRequestResponse, error) {
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
		return nil, ErrLabTestTemplateNotFound
	}

	// Generate request code
	code, err := s.requestRepo.GenerateRequestCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate request code: %w", err)
	}

	// Create request
	request := &domain.LabTestRequest{
		RequestCode:   code,
		VisitID:       req.VisitID,
		PatientID:     visit.PatientID,
		DoctorID:      visit.DoctorID,
		TemplateID:    req.TemplateID,
		Status:        domain.LabTestRequestStatusPending,
		Priority:      domain.LabTestPriority(req.Priority),
		RequestedDate: time.Now(),
		ClinicalNotes: req.ClinicalNotes,
		CreatedBy:     requestedBy,
	}

	if err := s.requestRepo.Create(request); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Auto-create result placeholders from template parameters
	results := make([]*domain.LabTestResult, len(template.Parameters))
	for i, param := range template.Parameters {
		results[i] = &domain.LabTestResult{
			RequestID:       request.ID,
			ParameterName:   param.ParameterName,
			Unit:            param.Unit,
			NormalRangeText: param.NormalRangeText,
			IsAbnormal:      false,
		}
	}

	if len(results) > 0 {
		if err := s.resultRepo.CreateBatch(results); err != nil {
			return nil, fmt.Errorf("failed to create result placeholders: %w", err)
		}
	}

	// Reload to get relationships
	request, _ = s.requestRepo.FindByID(request.ID)
	return s.toLabTestRequestResponse(request), nil
}

// CollectSample marks sample as collected
func (s *LabTestRequestService) CollectSample(id uint, updatedBy uint) error {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return ErrLabTestRequestNotFound
	}

	now := time.Now()
	request.Status = domain.LabTestRequestStatusSampleCollected
	request.SampleCollectedAt = &now
	request.UpdatedBy = updatedBy
	return s.requestRepo.Update(request)
}

// StartProcessing marks test as in progress
func (s *LabTestRequestService) StartProcessing(id uint, updatedBy uint) error {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return ErrLabTestRequestNotFound
	}

	request.Status = domain.LabTestRequestStatusInProgress
	request.UpdatedBy = updatedBy
	return s.requestRepo.Update(request)
}

// CompleteTest marks test as completed
func (s *LabTestRequestService) CompleteTest(id uint, updatedBy uint) error {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return ErrLabTestRequestNotFound
	}

	now := time.Now()
	request.Status = domain.LabTestRequestStatusCompleted
	request.CompletedAt = &now
	request.UpdatedBy = updatedBy
	return s.requestRepo.Update(request)
}

// CancelTest cancels a test
func (s *LabTestRequestService) CancelTest(id uint, updatedBy uint) error {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return ErrLabTestRequestNotFound
	}

	request.Status = domain.LabTestRequestStatusCancelled
	request.UpdatedBy = updatedBy
	return s.requestRepo.Update(request)
}

// EnterResults enters test results with auto-abnormal flagging
func (s *LabTestRequestService) EnterResults(requestID uint, req *dto.EnterLabTestResultsRequest) error {
	// Get request with template parameters
	request, err := s.requestRepo.FindByID(requestID)
	if err != nil {
		return fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return ErrLabTestRequestNotFound
	}

	// Get existing results
	existingResults, err := s.resultRepo.FindByRequestID(requestID)
	if err != nil {
		return fmt.Errorf("failed to get results: %w", err)
	}

	// Create map of parameter name to template parameter for normal range lookup
	paramMap := make(map[string]*domain.LabTestTemplateParameter)
	for _, p := range request.Template.Parameters {
		paramMap[p.ParameterName] = p
	}

	// Update results with values and check for abnormal
	resultsToUpdate := make([]*domain.LabTestResult, 0)
	for _, resultReq := range req.Results {
		// Find existing result
		var result *domain.LabTestResult
		for _, r := range existingResults {
			if r.ParameterName == resultReq.ParameterName {
				result = r
				break
			}
		}

		if result != nil {
			result.Value = resultReq.Value
			result.Remarks = resultReq.Remarks

			// Auto-flag abnormal values
			if param, ok := paramMap[resultReq.ParameterName]; ok {
				if value, err := strconv.ParseFloat(resultReq.Value, 64); err == nil {
					if value < param.NormalRangeMin || value > param.NormalRangeMax {
						result.IsAbnormal = true
					} else {
						result.IsAbnormal = false
					}
				}
			}

			resultsToUpdate = append(resultsToUpdate, result)
		}
	}

	if len(resultsToUpdate) > 0 {
		return s.resultRepo.UpdateBatch(resultsToUpdate)
	}

	return nil
}

// GetLabTestRequestByID gets request by ID
func (s *LabTestRequestService) GetLabTestRequestByID(id uint) (*dto.LabTestRequestResponse, error) {
	request, err := s.requestRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return nil, ErrLabTestRequestNotFound
	}
	return s.toLabTestRequestResponse(request), nil
}

// GetLabTestRequestByCode gets request by code
func (s *LabTestRequestService) GetLabTestRequestByCode(code string) (*dto.LabTestRequestResponse, error) {
	request, err := s.requestRepo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find request: %w", err)
	}
	if request == nil {
		return nil, ErrLabTestRequestNotFound
	}
	return s.toLabTestRequestResponse(request), nil
}

// GetVisitLabTests gets visit's lab tests
func (s *LabTestRequestService) GetVisitLabTests(visitID uint) ([]*dto.LabTestRequestListItem, error) {
	requests, err := s.requestRepo.FindByVisitID(visitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get requests: %w", err)
	}

	items := make([]*dto.LabTestRequestListItem, len(requests))
	for i, r := range requests {
		items[i] = s.toLabTestRequestListItem(r)
	}
	return items, nil
}

// GetPatientLabTests gets patient's lab test history
func (s *LabTestRequestService) GetPatientLabTests(patientID uint, filters map[string]interface{}) ([]*dto.LabTestRequestListItem, error) {
	requests, err := s.requestRepo.FindByPatientID(patientID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get requests: %w", err)
	}

	items := make([]*dto.LabTestRequestListItem, len(requests))
	for i, r := range requests {
		items[i] = s.toLabTestRequestListItem(r)
	}
	return items, nil
}

// Helper functions
func (s *LabTestRequestService) toLabTestRequestResponse(r *domain.LabTestRequest) *dto.LabTestRequestResponse {
	resp := &dto.LabTestRequestResponse{
		ID:                r.ID,
		RequestCode:       r.RequestCode,
		VisitID:           r.VisitID,
		PatientID:         r.PatientID,
		DoctorID:          r.DoctorID,
		TemplateID:        r.TemplateID,
		Status:            string(r.Status),
		Priority:          string(r.Priority),
		RequestedDate:     r.RequestedDate,
		SampleCollectedAt: r.SampleCollectedAt,
		CompletedAt:       r.CompletedAt,
		ClinicalNotes:     r.ClinicalNotes,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
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
	}

	// Add results
	resp.Results = make([]*dto.LabTestResultResponse, len(r.Results))
	for i, result := range r.Results {
		resp.Results[i] = &dto.LabTestResultResponse{
			ID:              result.ID,
			ParameterName:   result.ParameterName,
			Value:           result.Value,
			Unit:            result.Unit,
			NormalRangeText: result.NormalRangeText,
			IsAbnormal:      result.IsAbnormal,
			Remarks:         result.Remarks,
		}
	}

	return resp
}

func (s *LabTestRequestService) toLabTestRequestListItem(r *domain.LabTestRequest) *dto.LabTestRequestListItem {
	item := &dto.LabTestRequestListItem{
		ID:            r.ID,
		RequestCode:   r.RequestCode,
		Status:        string(r.Status),
		Priority:      string(r.Priority),
		RequestedDate: r.RequestedDate,
	}

	if r.Patient != nil {
		item.PatientName = r.Patient.FullName
	}
	if r.Template != nil {
		item.TemplateName = r.Template.Name
	}

	return item
}
