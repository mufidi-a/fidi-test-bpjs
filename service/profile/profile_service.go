package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	transform "test-bpjs/v2/helper/transform"
	"test-bpjs/v2/models"
	"test-bpjs/v2/models/request"
	"test-bpjs/v2/models/response"
	"test-bpjs/v2/repository"
	"time"
)

type ProfileService interface {
	GetProfileByCode(ctx context.Context, code int) (*response.CreateProfileResponse, error)                //v
	GetWorkingExperienceByCode(ctx context.Context, code int) (*response.WorkingExperiencesResponse, error) //v
	CreateProfile(ctx context.Context, payload request.CreateProfileRequest) (*response.DefaultResponse, error)
	UpdateProfile(ctx context.Context, payload request.UpdateProfileRequest) (*response.DefaultResponse, error)
	DeletePhotoByCode(ctx context.Context, code int) (*response.DefaultResponse, error)
	UploadPhotoByCode(ctx context.Context, payload request.UploadPhotoRequest) (*response.UploadPhotoResponse, error)
	DownloadPhotoByCode(ctx context.Context, code int) (string, error)
}

type profileService struct {
	profileRepo repository.ProfileRepository
}

func NewProfileService(profileRepo repository.ProfileRepository) *profileService {
	return &profileService{profileRepo: profileRepo}
}

func (p *profileService) GetProfileByCode(ctx context.Context, code int) (*response.CreateProfileResponse, error) {
	profile, err := p.profileRepo.GetProfileByCode(ctx, code)
	if err != nil {
		errString := fmt.Errorf("failed to get profile: %v", err)
		return nil, errString
	}
	return transform.TransformProfile(profile), err
}

func (p *profileService) GetWorkingExperienceByCode(ctx context.Context, code int) (*response.WorkingExperiencesResponse, error) {
	workingExperiences, err := p.profileRepo.GetWorkingExperienceByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get working experience: %v", err)
	}
	return &response.WorkingExperiencesResponse{
		WorkingExperience: workingExperiences.WorkingExperience,
	}, nil
}

func (p *profileService) CreateProfile(ctx context.Context, payload request.CreateProfileRequest) (*response.DefaultResponse, error) {
	profile, err := p.profileRepo.CreateProfile(ctx, &models.Profile{
		WantedJobTitle: payload.WantedJobTitle,
		FirstName:      payload.FirstName,
		LastName:       payload.LastName,
		Email:          payload.Email,
		Phone:          payload.Phone,
		Country:        payload.Country,
		City:           payload.City,
		Address:        payload.Address,
		PostalCode:     payload.PostalCode,
		DrivingLicense: payload.DrivingLicense,
		Nationality:    payload.Nationality,
		PlaceOfBirth:   payload.PlaceOfBirth,
		DateOfBirth:    payload.DateOfBirth,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create profile: %v", err)
	}
	return &response.DefaultResponse{
		ProfileCode: profile.ProfileCode,
	}, nil
}

func (p *profileService) UpdateProfile(ctx context.Context, payload request.UpdateProfileRequest) (*response.DefaultResponse, error) {
	profile, err := p.profileRepo.UpdateProfile(ctx, payload.ProfileCode, &models.Profile{
		WantedJobTitle:    payload.WantedJobTitle,
		FirstName:         payload.FirstName,
		LastName:          payload.LastName,
		Email:             payload.Email,
		Phone:             payload.Phone,
		Country:           payload.Country,
		City:              payload.City,
		Address:           payload.Address,
		PostalCode:        payload.PostalCode,
		DrivingLicense:    payload.DrivingLicense,
		Nationality:       payload.Nationality,
		PlaceOfBirth:      payload.PlaceOfBirth,
		DateOfBirth:       payload.DateOfBirth,
		WorkingExperience: payload.WorkingExperience,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %v", err)
	}
	return &response.DefaultResponse{
		ProfileCode: profile.ProfileCode,
	}, nil
}

func (p *profileService) DeletePhotoByCode(ctx context.Context, code int) (*response.DefaultResponse, error) {
	profileCode, err := p.profileRepo.DeletePhotoByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to delete photo: %v", err)
	}
	return &response.DefaultResponse{
		ProfileCode: profileCode.ProfileCode,
	}, nil
}

func (p *profileService) UploadPhotoByCode(ctx context.Context, payload request.UploadPhotoRequest) (*response.UploadPhotoResponse, error) {
	imgPath := filepath.Join("public/image", fmt.Sprintf("%d-%d.png", payload.ProfileCode, time.Now().Unix()))
	imgFullPath := filepath.Join("../../", imgPath)
	imgFullPath2 := filepath.Join("../", imgPath)

	b64data := payload.Base64Img[strings.IndexByte(payload.Base64Img, ',')+1:]
	imgData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %v", err)
	}

	imgReader := bytes.NewReader(imgData)

	img, _, err := image.Decode(imgReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	outFile, err := os.Create(imgFullPath)
	if err != nil {
		outFile, err = os.Create(imgFullPath2)
		if err != nil {
			outFile, err = os.Create(imgPath)
			if err != nil {
				return nil, fmt.Errorf("failed creating output file: %v", err)
			}
		}
	}

	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image to file: %v", err)
	}

	profile, err := p.profileRepo.UpdateProfile(ctx, payload.ProfileCode, &models.Profile{
		PhotoUrl: imgPath,
	})
	if err != nil {
		return nil, err
	}

	return &response.UploadPhotoResponse{
		ProfileCode: profile.ProfileCode,
		PhotoUrl:    imgPath,
	}, nil
}

func (p *profileService) DownloadPhotoByCode(ctx context.Context, code int) (string, error) {
	var buf bytes.Buffer

	profile, err := p.profileRepo.GetProfileByCode(ctx, code)
	if err != nil {
		return "", err
	}

	imgFullPath := filepath.Join("../../", profile.PhotoUrl)
	imgFile, err := os.Open(imgFullPath)
	if err != nil {
		imgFile, err = os.Open(profile.PhotoUrl)
		if err != nil {
			return "", fmt.Errorf("failed to open image file: %v", err)
		}
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	err = png.Encode(&buf, img)
	if err != nil {
		return "", fmt.Errorf("failed to encode image: %v", err)
	}

	// Encode the buffer to base64 string
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	res := fmt.Sprintf("data:image/png;base64,%s", base64Str)

	return res, nil
}
