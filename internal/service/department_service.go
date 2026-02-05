package service

import (
	"errors"

	"github.com/minhtran/his/internal/domain"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

var (
	// ErrDepartmentCodeExists is returned when department code already exists
	ErrDepartmentCodeExists = errors.New("department code already exists")
	// ErrDepartmentNotFound is returned when department is not found
	ErrDepartmentNotFound = errors.New("department not found")
)

// DepartmentService handles business logic for departments
type DepartmentService struct {
	deptRepo  *repository.DepartmentRepository
	auditRepo *repository.AuditLogRepository
}

// NewDepartmentService creates a new department service
func NewDepartmentService(deptRepo *repository.DepartmentRepository, auditRepo *repository.AuditLogRepository) *DepartmentService {
	return &DepartmentService{
		deptRepo:  deptRepo,
		auditRepo: auditRepo,
	}
}

// CreateDepartment creates a new department
func (s *DepartmentService) CreateDepartment(req *dto.CreateDepartmentRequest, userID uint) (*dto.DepartmentResponse, error) {
	// Check if code exists
	existing, _ := s.deptRepo.FindByCode(req.Code)
	if existing != nil {
		return nil, ErrDepartmentCodeExists
	}

	dept := &domain.Department{
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		HeadDoctorID: req.HeadDoctorID,
	}

	dept.CreatedBy = userID
	dept.UpdatedBy = userID

	if err := s.deptRepo.Create(dept); err != nil {
		return nil, err
	}

	// Audit log
	s.auditRepo.Create(&domain.AuditLog{
		UserID:     &userID,
		Action:     domain.AuditActionCreate,
		Resource:   "Department",
		ResourceID: dept.Code, // Using Code as ID string for readability, or use strconv.Itoa(int(dept.ID))
		Details:    domain.AuditDetails{"name": dept.Name, "code": dept.Code},
	})

	return s.toDepartmentResponse(dept), nil
}

// GetDepartment returns a department by ID
func (s *DepartmentService) GetDepartment(id uint) (*dto.DepartmentResponse, error) {
	dept, err := s.deptRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, ErrDepartmentNotFound
	}
	return s.toDepartmentResponse(dept), nil
}

// UpdateDepartment updates a department
func (s *DepartmentService) UpdateDepartment(id uint, req *dto.UpdateDepartmentRequest, userID uint) (*dto.DepartmentResponse, error) {
	dept, err := s.deptRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, ErrDepartmentNotFound
	}

	dept.Name = req.Name
	dept.Description = req.Description
	dept.HeadDoctorID = req.HeadDoctorID
	dept.IsActive = req.IsActive
	dept.UpdatedBy = userID

	if err := s.deptRepo.Update(dept); err != nil {
		return nil, err
	}

	// Audit log
	s.auditRepo.Create(&domain.AuditLog{
		UserID:     &userID,
		Action:     domain.AuditActionUpdate,
		Resource:   "Department",
		ResourceID: dept.Code,
		Details:    domain.AuditDetails{"id": id, "changes": "updated details"},
	})

	return s.toDepartmentResponse(dept), nil
}

// DeleteDepartment deletes a department
func (s *DepartmentService) DeleteDepartment(id uint, userID uint) error {
	dept, err := s.deptRepo.FindByID(id)
	if err != nil {
		return err
	}
	if dept == nil {
		return ErrDepartmentNotFound
	}

	if err := s.deptRepo.Delete(id); err != nil {
		return err
	}

	// Audit log
	s.auditRepo.Create(&domain.AuditLog{
		UserID:     &userID,
		Action:     domain.AuditActionDelete,
		Resource:   "Department",
		ResourceID: dept.Code,
		Details:    domain.AuditDetails{"id": id},
	})

	return nil
}

// ListDepartments returns a list of departments
func (s *DepartmentService) ListDepartments(page, pageSize int) ([]*dto.DepartmentResponse, int64, error) {
	depts, total, err := s.deptRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.DepartmentResponse, len(depts))
	for i, d := range depts {
		responses[i] = s.toDepartmentResponse(d)
	}

	return responses, total, nil
}

// toDepartmentResponse converts domain.Department to dto.DepartmentResponse
func (s *DepartmentService) toDepartmentResponse(dept *domain.Department) *dto.DepartmentResponse {
	var headDoctor *dto.UserDetailResponse
	if dept.HeadDoctor != nil {
		headDoctor = &dto.UserDetailResponse{
			ID:          dept.HeadDoctor.ID,
			Username:    dept.HeadDoctor.Username,
			Email:       dept.HeadDoctor.Email,
			FullName:    dept.HeadDoctor.FullName,
			PhoneNumber: dept.HeadDoctor.PhoneNumber,
			IsActive:    dept.HeadDoctor.IsActive,
		}
	}

	return &dto.DepartmentResponse{
		ID:          dept.ID,
		Code:        dept.Code,
		Name:        dept.Name,
		Description: dept.Description,
		IsActive:    dept.IsActive,
		HeadDoctor:  headDoctor,
		CreatedAt:   dept.CreatedAt,
		UpdatedAt:   dept.UpdatedAt,
	}
}
