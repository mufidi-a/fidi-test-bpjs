package repository

import (
	"context"
	"test-bpjs/v2/models"

	"github.com/uptrace/bun"
)

type SkillRepository interface {
	GetSkillsByProfileCode(ctx context.Context, code int) ([]*models.SkillDTO, error)
	CreateSkill(ctx context.Context, payload *models.Skill) (*models.SkillDTO, error)
	DeleteSkill(ctx context.Context, code, id int) error
}

type skillRepository struct {
	DB bun.IDB
}

func NewSkillRepository(db bun.IDB) *skillRepository {
	return &skillRepository{
		DB: db,
	}
}

func (s *skillRepository) GetSkillsByProfileCode(ctx context.Context, code int) ([]*models.SkillDTO, error) {
	var skill []*models.SkillDTO
	err := s.DB.NewSelect().
		Model((*models.Skill)(nil)).
		Column("id", "skill", "level").
		Where("profile_code = ?", code).
		Scan(ctx, &skill)
	return skill, err
}

func (s *skillRepository) CreateSkill(ctx context.Context, payload *models.Skill) (*models.SkillDTO, error) {
	var skill models.SkillDTO
	_, err := s.DB.NewInsert().
		Model(payload).
		Returning("id").
		Exec(ctx, &skill)
	return &skill, err
}

func (s *skillRepository) DeleteSkill(ctx context.Context, code, id int) error {
	var skill models.SkillDTO
	_, err := s.DB.NewDelete().
		Model((*models.Skill)(nil)).
		Where("profile_code = ?", code).
		Where("id = ?", id).
		Returning("profile_code").
		Exec(ctx, &skill)
	return err
}
