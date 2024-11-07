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

type EducationService interface {
	GetEducationByCode(ctx context.Context, code int) (*response.EducationList, error)
	CreateEducation(ctx context.Context, payload request.CreateEducationRequest) (*response.DefaultResponseWithId, error)
	DeleteEducation(ctx context.Context, code, id int) (*response.DefaultResponse, error)
}

type educationService struct {
	educationRepo repository.EducationRepository
}

func NewEducationService(educationRepo repository.EducationRepository) *educationService {
	return &educationService{educationRepo: educationRepo}
}

func (s *educationService) GetEducationByCode(ctx context.Context, code int) (*response.EducationList, error) {
	var educationList []*response.EducationResponse

	educations, err := s.educationRepo.GetEducationByProfileCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get education: %v", err)
	}

	for _, education := range educations {
		educationList = append(educationList, transform.TransformEducation(education))
	}
	return &response.EducationList{
		Data: educationList,
	}, nil
}

func (s *educationService) CreateEducation(ctx context.Context, payload request.CreateEducationRequest) (*response.DefaultResponseWithId, error) {
	education, err := s.educationRepo.CreateEducation(ctx, &models.Education{
		ProfileCode: payload.ProfileCode,
		School:      payload.School,
		Degree:      payload.Degree,
		StartDate:   payload.StartDate,
		EndDate:     payload.EndDate,
		City:        payload.City,
		Description: payload.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create education: %v", err)
	}
	return &response.DefaultResponseWithId{
		ProfileCode: payload.ProfileCode,
		Id:          education.Id,
	}, nil
}

func (s *educationService) DeleteEducation(ctx context.Context, code, id int) (*response.DefaultResponse, error) {
	err := s.educationRepo.DeleteEducation(ctx, code, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete education: %v", err)
	}
	return &response.DefaultResponse{
		ProfileCode: code,
	}, nil
}
