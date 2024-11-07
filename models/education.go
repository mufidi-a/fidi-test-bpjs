package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Education struct {
	bun.BaseModel `bun:"table:education"`

	ProfileCode int       `bun:"profile_code"`
	Id          int       `bun:"id,pk,type:int,autoincrement"`
	School      string    `bun:"school"`
	Degree      string    `bun:"degree"`
	StartDate   time.Time `bun:"start_date"`
	EndDate     time.Time `bun:"end_date"`
	City        string    `bun:"city"`
	Description string    `bun:"description"`
	CreatedAt   time.Time `bun:"created_at,default:current_timestamp"`
}

type EducationDTO struct {
	ProfileCode int       `json:"profileCode"`
	Id          int       `json:"id"`
	School      string    `json:"school"`
	Degree      string    `json:"degree"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	City        string    `json:"city"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
