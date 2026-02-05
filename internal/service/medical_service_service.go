package service

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

var (
	// ErrServiceCodeExists is returned when medical service code already exists
	ErrServiceCodeExists = errors.New("service code already exists")
	// ErrServiceNotFound is returned when medical service is not found
	ErrServiceNotFound = errors.New("service not found")
)

// MedicalServiceService handles business logic for medical services
type MedicalServiceService struct {
	serviceRepo *repository.MedicalServiceRepository
	auditRepo   *repository.AuditLogRepository
}

// NewMedicalServiceService creates a new medical service service
func NewMedicalServiceService(serviceRepo *repository.MedicalServiceRepository, auditRepo *repository.AuditLogRepository) *MedicalServiceService {
	return &MedicalServiceService{
		serviceRepo: serviceRepo,
		auditRepo:   auditRepo,
	}
}

// CreateService creates a new medical service
func (s *MedicalServiceService) CreateService(req *dto.CreateMedicalServiceRequest, userID uint) (*dto.MedicalServiceResponse, error) {
	existing, _ := s.serviceRepo.FindByCode(req.Code)
	if existing != nil {
		return nil, ErrServiceCodeExists
	}

	service := &domain.MedicalService{
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		ServiceType:  domain.ServiceType(req.ServiceType),
		BasePrice:    req.BasePrice,
		DepartmentID: req.DepartmentID,
	}

	if err := s.serviceRepo.Create(service); err != nil {
		return nil, err
	}

	s.auditRepo.Create(&domain.AuditLog{
		UserID:     &userID,
		Action:     domain.AuditActionCreate,
		Resource:   "MedicalService",
		ResourceID: service.Code,
		Details:    domain.AuditDetails{"name": service.Name, "price": service.BasePrice},
	})

	return s.toMedicalServiceResponse(service), nil
}

// UpdateService updates a medical service
func (s *MedicalServiceService) UpdateService(id uint, req *dto.UpdateMedicalServiceRequest, userID uint) (*dto.MedicalServiceResponse, error) {
	service, err := s.serviceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrServiceNotFound
	}

	service.Name = req.Name
	service.Description = req.Description
	if req.ServiceType != "" {
		service.ServiceType = domain.ServiceType(req.ServiceType)
	}
	if req.BasePrice >= 0 {
		service.BasePrice = req.BasePrice
	}
	service.IsActive = req.IsActive
	service.DepartmentID = req.DepartmentID

	if err := s.serviceRepo.Update(service); err != nil {
		return nil, err
	}

	s.auditRepo.Create(&domain.AuditLog{
		UserID:     &userID,
		Action:     domain.AuditActionUpdate,
		Resource:   "MedicalService",
		ResourceID: service.Code,
		Details:    domain.AuditDetails{"id": id},
	})

	return s.toMedicalServiceResponse(service), nil
}

// ListServices returns a list of services
func (s *MedicalServiceService) ListServices(page, pageSize int, departmentID *uint) ([]*dto.MedicalServiceResponse, int64, error) {
	services, total, err := s.serviceRepo.List(page, pageSize, departmentID)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.MedicalServiceResponse, len(services))
	for i, svc := range services {
		responses[i] = s.toMedicalServiceResponse(svc)
	}

	return responses, total, nil
}

// GetService returns a service by ID
func (s *MedicalServiceService) GetService(id uint) (*dto.MedicalServiceResponse, error) {
	service, err := s.serviceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrServiceNotFound
	}
	return s.toMedicalServiceResponse(service), nil
}

// toMedicalServiceResponse converts domain.MedicalService to dto.MedicalServiceResponse
func (s *MedicalServiceService) toMedicalServiceResponse(service *domain.MedicalService) *dto.MedicalServiceResponse {
	var deptResp *dto.DepartmentResponse
	if service.Department != nil {
		deptResp = &dto.DepartmentResponse{
			ID:          service.Department.ID,
			Code:        service.Department.Code,
			Name:        service.Department.Name,
			Description: service.Department.Description,
			IsActive:    service.Department.IsActive,
			CreatedAt:   service.Department.CreatedAt,
			UpdatedAt:   service.Department.UpdatedAt,
		}
	}

	return &dto.MedicalServiceResponse{
		ID:          service.ID,
		Code:        service.Code,
		Name:        service.Name,
		Description: service.Description,
		ServiceType: string(service.ServiceType),
		BasePrice:   service.BasePrice,
		IsActive:    service.IsActive,
		Department:  deptResp,
		CreatedAt:   service.CreatedAt,
		UpdatedAt:   service.UpdatedAt,
	}
}
