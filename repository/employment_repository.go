package repository

import (
	"context"
	"test-bpjs/v2/models"

	"github.com/uptrace/bun"
)

type EmploymentRepository interface {
	GetEmploymentByProfileCode(ctx context.Context, code int) ([]*models.EmploymentDTO, error)
	CreateEmployment(ctx context.Context, payload *models.Employment) (*models.EmploymentDTO, error)
	DeleteEmployment(ctx context.Context, id, code int) error
}

type employmentRepository struct {
	DB bun.IDB
}

func NewEmploymentRepository(db bun.IDB) *employmentRepository {
	return &employmentRepository{
		DB: db,
	}
}

func (e *employmentRepository) GetEmploymentByProfileCode(ctx context.Context, code int) ([]*models.EmploymentDTO, error) {
	var employment []*models.EmploymentDTO
	err := e.DB.NewSelect().
		Model((*models.Employment)(nil)).
		Column("id", "job_title", "employer", "start_date", "end_date", "city", "description").
		Where("profile_code = ?", code).
		Scan(ctx, &employment)
	return employment, err
}

func (e *employmentRepository) CreateEmployment(ctx context.Context, payload *models.Employment) (*models.EmploymentDTO, error) {
	var employment models.EmploymentDTO
	_, err := e.DB.NewInsert().
		Model(payload).
		Returning("id").
		Exec(ctx, &employment)
	return &employment, err
}

func (e *employmentRepository) DeleteEmployment(ctx context.Context, code, id int) error {
	var employment models.EmploymentDTO
	_, err := e.DB.NewDelete().
		Model((*models.Employment)(nil)).
		Where("profile_code = ?", code).
		Where("id = ?", id).
		Returning("profile_code").
		Exec(ctx, &employment)
	return err
}
