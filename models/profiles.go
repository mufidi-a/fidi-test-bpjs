package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Profile struct {
	bun.BaseModel `bun:"table:profile"`

	ProfileCode       int       `bun:"profile_code,pk,type:int,autoincrement"`
	WantedJobTitle    string    `bun:"wanted_job_title,notnull"`
	FirstName         string    `bun:"first_name,notnull"`
	LastName          string    `bun:"last_name"`
	Email             string    `bun:"email,notnull"`
	Phone             string    `bun:"phone,notnull"`
	Country           string    `bun:"country,notnull"`
	City              string    `bun:"city,notnull"`
	Address           string    `bun:"address,notnull"`
	PostalCode        int       `bun:"postal_code"`
	DrivingLicense    string    `bun:"driving_license"`
	Nationality       string    `bun:"nationality"`
	PlaceOfBirth      string    `bun:"place_of_birth"`
	DateOfBirth       time.Time `bun:"date_of_birth"`
	PhotoUrl          string    `bun:"photo_url"`
	WorkingExperience string    `bun:"working_experience"`
	CreatedAt         time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt         time.Time `bun:"updated_at"`
}

type ProfileDTO struct {
	ProfileCode       int       `json:"profileCode"`
	WantedJobTitle    string    `json:"wantedJobTitle"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Country           string    `json:"country"`
	City              string    `json:"city"`
	Address           string    `json:"address"`
	PostalCode        int       `json:"postalCode"`
	DrivingLicense    string    `json:"drivingLicense"`
	Nationality       string    `json:"nationality"`
	PlaceOfBirth      string    `json:"placeOfBirth"`
	DateOfBirth       time.Time `json:"dateOfBirth"`
	PhotoUrl          string    `json:"photoUrl"`
	WorkingExperience string    `json:"workingExperience"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
