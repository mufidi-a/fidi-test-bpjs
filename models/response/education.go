package response

import (
	"time"
)

type EducationResponse struct {
	Id          int       `json:"id"`
	School      string    `json:"school"`
	Degree      string    `json:"degree"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	City        string    `json:"city"`
	Description string    `json:"description"`
}
type EducationList struct {
	Data []*EducationResponse `json:"data"`
}
