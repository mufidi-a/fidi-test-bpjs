package transform

import (
	"test-bpjs/v2/models"
	"test-bpjs/v2/models/response"
)

func TransformEducation(education *models.EducationDTO) *response.EducationResponse {
	return &response.EducationResponse{
		Id:          education.Id,
		School:      education.School,
		Degree:      education.Degree,
		StartDate:   education.StartDate,
		EndDate:     education.EndDate,
		City:        education.City,
		Description: education.Description,
	}
}
