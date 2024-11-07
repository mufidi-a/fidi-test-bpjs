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

type SkillService interface {
	GetSkillsByCode(ctx context.Context, code int) (*response.SkillList, error)
	CreateSkill(ctx context.Context, payload request.CreateSkillRequest) (*response.DefaultResponseWithId, error)
	DeleteSkill(ctx context.Context, code, id int) (*response.DefaultResponse, error)
}

type skillService struct {
	skillRepo repository.SkillRepository
}

func NewSkillService(skillRepo repository.SkillRepository) *skillService {
	return &skillService{skillRepo: skillRepo}
}

func (s *skillService) GetSkillsByCode(ctx context.Context, code int) (*response.SkillList, error) {
	var skillList []*response.SkillResponse

	skills, err := s.skillRepo.GetSkillsByProfileCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %v", err)
	}

	for _, skill := range skills {
		skillList = append(skillList, transform.TransformSkill(skill))
	}
	return &response.SkillList{
		Data: skillList,
	}, nil
}

func (s *skillService) CreateSkill(ctx context.Context, payload request.CreateSkillRequest) (*response.DefaultResponseWithId, error) {
	skill, err := s.skillRepo.CreateSkill(ctx, &models.Skill{
		ProfileCode: payload.ProfileCode,
		Skill:       payload.Skill,
		Level:       payload.Level,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create skill: %v", err)
	}
	return &response.DefaultResponseWithId{
		ProfileCode: payload.ProfileCode,
		Id:          skill.Id,
	}, nil
}

func (s *skillService) DeleteSkill(ctx context.Context, code, id int) (*response.DefaultResponse, error) {
	err := s.skillRepo.DeleteSkill(ctx, code, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete skill: %v", err)
	}
	return &response.DefaultResponse{
		ProfileCode: code,
	}, nil
}
