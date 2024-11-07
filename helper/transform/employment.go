package transform

import (
	"test-bpjs/v2/models"
	"test-bpjs/v2/models/response"
)

func TransformEmployment(employment *models.EmploymentDTO) *response.EmploymentResponse {
	return &response.EmploymentResponse{
		Id:          employment.Id,
		JobTitle:    employment.JobTitle,
		Employer:    employment.Employer,
		StartDate:   employment.StartDate,
		EndDate:     employment.EndDate,
		City:        employment.City,
		Description: employment.Description,
	}
}
