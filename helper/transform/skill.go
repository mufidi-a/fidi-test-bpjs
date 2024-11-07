package transform

import (
	"test-bpjs/v2/models"
	"test-bpjs/v2/models/response"
)

func TransformSkill(skill *models.SkillDTO) *response.SkillResponse {
	return &response.SkillResponse{
		Id:    skill.Id,
		Skill: skill.Skill,
		Level: skill.Level,
	}
}
