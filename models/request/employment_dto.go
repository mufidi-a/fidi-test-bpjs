package request

import (
	"time"
)

type CreateEmploymentRequest struct {
	ProfileCode int       `param:"profileCode" validate:"required"`
	JobTitle    string    `json:"jobTitle"`
	Employer    string    `json:"employer"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	City        string    `json:"city"`
	Description string    `json:"description"`
}
