package models

type DefaultResponse struct {
	ProfileCode int `json:"profileCode"`
}

type DefaultResponseWithId struct {
	ProfileCode int `json:"profileCode"`
	Id          int `json:"id"`
}
