package request

import "time"

type UploadPhotoRequest struct {
	ProfileCode int    `param:"profileCode" validate:"required"`
	Base64Img   string `json:"base64img"`
}

type GetProfileRequest struct {
	ProfileCode int `param:"profileCode" validate:"required"`
}

type CreateProfileRequest struct {
	WantedJobTitle string    `json:"wantedJobTitle"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Country        string    `json:"country"`
	City           string    `json:"city"`
	Address        string    `json:"address"`
	PostalCode     int       `json:"postalCode"`
	DrivingLicense string    `json:"drivingLicense"`
	Nationality    string    `json:"nationality"`
	PlaceOfBirth   string    `json:"placeOfBirth"`
	DateOfBirth    time.Time `json:"dateOfBirth"`
}

type UpdateProfileRequest struct {
	ProfileCode       int       `param:"profileCode" validate:"required"`
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
	WorkingExperience string    `json:"workingExperience"`
}
