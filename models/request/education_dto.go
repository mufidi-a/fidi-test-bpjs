package request

import "time"

type CreateEducationRequest struct {
	ProfileCode int       `param:"profileCode" validate:"required"`
	School      string    `json:"school"`
	Degree      string    `json:"degree"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	City        string    `json:"city"`
	Description string    `json:"description"`
}

type DefaultGetDataByProfileCodeAndId struct {
	ProfileCode int `param:"profileCode" validate:"required"`
	Id          int `query:"id" validate:"required"`
}
