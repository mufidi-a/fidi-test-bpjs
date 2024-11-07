package response

import (
	"time"
)

type EmploymentResponse struct {
	Id          int       `json:"id"`
	JobTitle    string    `json:"jobTitle"`
	Employer    string    `json:"employer"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	City        string    `json:"city"`
	Description string    `json:"description"`
}

type EmploymentList struct {
	Data []*EmploymentResponse `json:"data"`
}
