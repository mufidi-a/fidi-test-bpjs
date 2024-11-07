package repository

import (
	"context"
	"test-bpjs/v2/models"

	"github.com/uptrace/bun"
)

type EducationRepository interface {
	GetEducationByProfileCode(ctx context.Context, code int) ([]*models.EducationDTO, error)
	CreateEducation(ctx context.Context, payload *models.Education) (*models.EducationDTO, error)
	DeleteEducation(ctx context.Context, code, id int) error
}

type educationRepository struct {
	DB bun.IDB
}

func NewEducationRepository(db bun.IDB) *educationRepository {
	return &educationRepository{
		DB: db,
	}
}

func (e *educationRepository) GetEducationByProfileCode(ctx context.Context, code int) ([]*models.EducationDTO, error) {
	var education []*models.EducationDTO
	err := e.DB.NewSelect().
		Model((*models.Education)(nil)).
		Column("id", "school", "degree", "start_date", "end_date", "city", "description").
		Where("profile_code = ?", code).
		Scan(ctx, &education)
	return education, err
}

func (e *educationRepository) CreateEducation(ctx context.Context, payload *models.Education) (*models.EducationDTO, error) {
	var education models.EducationDTO
	_, err := e.DB.NewInsert().
		Model(payload).
		Returning("id").
		Exec(ctx, &education)
	return &education, err
}

func (e *educationRepository) DeleteEducation(ctx context.Context, code, id int) error {
	var education models.EducationDTO
	_, err := e.DB.NewDelete().
		Model((*models.Education)(nil)).
		Where("profile_code = ?", code).
		Where("id = ?", id).
		Returning("profile_code").
		Exec(ctx, &education)
	return err
}
