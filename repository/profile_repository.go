package repository

import (
	"context"
	"test-bpjs/v2/models"
	"time"

	"github.com/uptrace/bun"
)

type ProfileRepository interface {
	GetProfileByCode(ctx context.Context, code int) (*models.ProfileDTO, error)
	GetWorkingExperienceByCode(ctx context.Context, code int) (*models.ProfileDTO, error)
	CreateProfile(ctx context.Context, payload *models.Profile) (*models.ProfileDTO, error)
	UpdateProfile(ctx context.Context, code int, payload *models.Profile) (*models.ProfileDTO, error)
	DeletePhotoByCode(ctx context.Context, code int) (*models.DefaultResponse, error)
}

type profileRepository struct {
	DB bun.IDB
}

func NewProfileRepository(db bun.IDB) *profileRepository {
	return &profileRepository{
		DB: db,
	}
}

func (p *profileRepository) GetProfileByCode(ctx context.Context, code int) (*models.ProfileDTO, error) {
	var profile models.ProfileDTO
	err := p.DB.NewSelect().
		Model((*models.Profile)(nil)).
		Column("profile_code", "wanted_job_title", "first_name", "last_name", "email", "phone", "country", "city", "address", "postal_code", "driving_license", "nationality", "place_of_birth", "date_of_birth", "photo_url").
		Where("profile_code = ?", code).
		Scan(ctx, &profile)
	return &profile, err
}

func (p *profileRepository) GetWorkingExperienceByCode(ctx context.Context, code int) (*models.ProfileDTO, error) {
	var profile models.ProfileDTO
	err := p.DB.NewSelect().
		Model((*models.Profile)(nil)).
		Column("working_experience").
		Where("profile_code = ?", code).
		Scan(ctx, &profile)
	return &profile, err
}

func (p *profileRepository) CreateProfile(ctx context.Context, payload *models.Profile) (*models.ProfileDTO, error) {
	var profile models.ProfileDTO
	_, err := p.DB.NewInsert().
		Model(payload).
		Returning("profile_code").
		Exec(ctx, &profile)
	return &profile, err
}

func (p *profileRepository) UpdateProfile(ctx context.Context, code int, payload *models.Profile) (*models.ProfileDTO, error) {
	var profile models.ProfileDTO
	payload.UpdatedAt = time.Now()
	_, err := p.DB.NewUpdate().
		Model(payload).
		OmitZero().
		Where("profile_code = ?", code).
		Returning("profile_code").
		Exec(ctx, &profile)
	return &profile, err
}

func (p *profileRepository) DeletePhotoByCode(ctx context.Context, code int) (*models.DefaultResponse, error) {
	var profile models.ProfileDTO
	_, err := p.DB.NewUpdate().
		Model((*models.Profile)(nil)).
		Set("photo_url = NULL").
		Where("profile_code = ?", code).
		Returning("profile_code").
		Exec(ctx, &profile)
	return &models.DefaultResponse{
		ProfileCode: profile.ProfileCode,
	}, err
}
