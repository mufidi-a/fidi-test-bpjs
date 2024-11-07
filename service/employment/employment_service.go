package service

import (
	"context"
	"fmt"
	transform "test-bpjs/v2/helper/transform"
	"test-bpjs/v2/models"
	"test-bpjs/v2/models/request"
	"test-bpjs/v2/models/response"
	"test-bpjs/v2/repository"
)

type EmploymentService interface {
	GetEmploymentByCode(ctx context.Context, code int) (*response.EmploymentList, error)
	CreateEmployment(ctx context.Context, payload request.CreateEmploymentRequest) (*response.DefaultResponseWithId, error)
	DeleteEmployment(ctx context.Context, code, id int) (*response.DefaultResponse, error)
}

type employmentService struct {
	employmentRepo repository.EmploymentRepository
}

func NewEmploymentService(employmentRepo repository.EmploymentRepository) *employmentService {
	return &employmentService{employmentRepo: employmentRepo}
}

func (e *employmentService) GetEmploymentByCode(ctx context.Context, code int) (*response.EmploymentList, error) {
	var employmentList []*response.EmploymentResponse

	employments, err := e.employmentRepo.GetEmploymentByProfileCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get employments: %v", err)
	}

	for _, employment := range employments {
		employmentList = append(employmentList, transform.TransformEmployment(employment))
	}
	return &response.EmploymentList{
		Data: employmentList,
	}, nil
}

func (e *employmentService) CreateEmployment(ctx context.Context, payload request.CreateEmploymentRequest) (*response.DefaultResponseWithId, error) {
	employment, err := e.employmentRepo.CreateEmployment(ctx, &models.Employment{
		ProfileCode: payload.ProfileCode,
		JobTitle:    payload.JobTitle,
		Employer:    payload.Employer,
		StartDate:   payload.StartDate,
		EndDate:     payload.EndDate,
		City:        payload.City,
		Description: payload.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create employment: %v", err)
	}
	return &response.DefaultResponseWithId{
		ProfileCode: payload.ProfileCode,
		Id:          employment.Id,
	}, nil
}

func (s *employmentService) DeleteEmployment(ctx context.Context, code, id int) (*response.DefaultResponse, error) {
	err := s.employmentRepo.DeleteEmployment(ctx, code, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete employment: %v", err)
	}
	return &response.DefaultResponse{
		ProfileCode: code,
	}, nil
}
