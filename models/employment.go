package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Employment struct {
	bun.BaseModel `bun:"table:employment"`

	ProfileCode int       `bun:"profile_code"`
	Id          int       `bun:"id,pk,type:int,autoincrement"`
	JobTitle    string    `bun:"job_title"`
	Employer    string    `bun:"employer"`
	StartDate   time.Time `bun:"start_date"`
	EndDate     time.Time `bun:"end_date"`
	City        string    `bun:"city"`
	Description string    `bun:"description"`
	CreatedAt   time.Time `bun:"created_at,default:current_timestamp"`
}

type EmploymentDTO struct {
	ProfileCode int       `json:"profileCode"`
	Id          int       `json:"id"`
	JobTitle    string    `json:"jobTitle"`
	Employer    string    `json:"employer"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	City        string    `json:"city"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
