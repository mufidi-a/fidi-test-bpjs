package transform

import (
	"test-bpjs/v2/models"
	"test-bpjs/v2/models/response"
)

func TransformProfile(profile *models.ProfileDTO) *response.CreateProfileResponse {
	return &response.CreateProfileResponse{
		ProfileCode:    profile.ProfileCode,
		WantedJobTitle: profile.WantedJobTitle,
		FirstName:      profile.FirstName,
		LastName:       profile.LastName,
		Email:          profile.Email,
		Phone:          profile.Phone,
		Country:        profile.Country,
		City:           profile.City,
		Address:        profile.Address,
		PostalCode:     profile.PostalCode,
		DrivingLicense: profile.DrivingLicense,
		Nationality:    profile.Nationality,
		PlaceOfBirth:   profile.PlaceOfBirth,
		DateOfBirth:    profile.DateOfBirth,
		PhotoUrl:       profile.PhotoUrl,
	}
}
