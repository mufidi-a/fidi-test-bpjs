package service

import (
	"context"
	"errors"
	"fmt"
	"test-bpjs/v2/models"
	"test-bpjs/v2/models/request"
	repository "test-bpjs/v2/repository/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var profileRepository = &repository.ProfileRepository{Mock: mock.Mock{}}
var profileServiceTest = profileService{profileRepo: profileRepository}

func TestInitProfileService(t *testing.T) {
	t.Run("SuccessInitSProfileService", func(t *testing.T) {
		assert.NotNil(t, NewProfileService(profileRepository))
	})
}

func TestGetProfile(t *testing.T) {
	t.Run("SuccessGetProfile", func(t *testing.T) {
		result := &models.ProfileDTO{
			ProfileCode: 1,
		}
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 2).Return(result, nil)

		profile, err := profileServiceTest.GetProfileByCode(context.Background(), 2)
		assert.Nil(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, result.ProfileCode, profile.ProfileCode)
	})
	t.Run("FailedGetProfile", func(t *testing.T) {
		// program mock
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 1).Return(nil, errors.New("sql: no rows in result set"))

		profile, err := profileServiceTest.GetProfileByCode(context.Background(), 1)
		assert.Nil(t, profile)
		assert.NotNil(t, err)
		// assert.Equal(t, "failed to get profile: sql: no rows in result set", err.Error())
		assert.Contains(t, err.Error(), "failed to get profile:")
	})
}

func TestGetWorkingExperiences(t *testing.T) {
	t.Run("SuccessGetWorkingExperiences", func(t *testing.T) {
		workingExperience := &models.ProfileDTO{
			WorkingExperience: "test",
		}
		profileRepository.Mock.On("GetWorkingExperienceByCode", mock.Anything, 2).Return(workingExperience, nil)

		result, err := profileServiceTest.GetWorkingExperienceByCode(context.Background(), 2)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, workingExperience.WorkingExperience, result.WorkingExperience)
	})
	t.Run("FailedGetWorkingExperiences", func(t *testing.T) {
		// program mock
		profileRepository.Mock.On("GetWorkingExperienceByCode", mock.Anything, 1).Return(nil, errors.New("sql: no rows in result set"))

		profile, err := profileServiceTest.GetWorkingExperienceByCode(context.Background(), 1)
		assert.Nil(t, profile)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to get working experience:")
	})
}

func TestCreateProfile(t *testing.T) {
	t.Run("SuccessCreateProfile", func(t *testing.T) {
		profile := &models.ProfileDTO{
			ProfileCode: 123456,
		}
		profileRepository.Mock.On("CreateProfile", mock.Anything,
			&models.Profile{
				WantedJobTitle: "test",
				FirstName:      "test",
				Email:          "test",
				Phone:          "0888888888",
				Country:        "test",
				City:           "test",
				Address:        "test"}).Return(profile, nil)

		result, err := profileServiceTest.CreateProfile(context.Background(),
			request.CreateProfileRequest{
				WantedJobTitle: "test",
				FirstName:      "test",
				Email:          "test",
				Phone:          "0888888888",
				Country:        "test",
				City:           "test",
				Address:        "test",
			})
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, profile.ProfileCode, result.ProfileCode)
	})
	t.Run("FailedCreateProfile", func(t *testing.T) {
		// program mock
		profileRepository.Mock.On("CreateProfile", mock.Anything, &models.Profile{}).Return(nil, errors.New("NOT NULL VIOLATION"))

		profile, err := profileServiceTest.CreateProfile(context.Background(), request.CreateProfileRequest{})
		assert.Nil(t, profile)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to create profile:")
	})
}

func TestUpdateProfile(t *testing.T) {
	t.Run("SuccessUpdateProfile", func(t *testing.T) {
		profile := &models.ProfileDTO{
			ProfileCode: 3,
		}
		profileRepository.Mock.On("UpdateProfile", context.Background(), 3,
			&models.Profile{
				WantedJobTitle: "test",
				FirstName:      "test",
				Email:          "test",
				Phone:          "0888888889",
				Country:        "test",
				City:           "test",
				Address:        "test"}).Return(profile, nil)

		result, err := profileServiceTest.UpdateProfile(context.Background(),
			request.UpdateProfileRequest{
				ProfileCode:    3,
				WantedJobTitle: "test",
				FirstName:      "test",
				Email:          "test",
				Phone:          "0888888889",
				Country:        "test",
				City:           "test",
				Address:        "test",
			})
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, profile.ProfileCode, result.ProfileCode)
	})
	t.Run("FailedUpdateProfile", func(t *testing.T) {
		// program mock
		profileRepository.Mock.On("UpdateProfile", mock.Anything, 2, &models.Profile{}).Return(nil, errors.New("NOT NULL VIOLATION"))

		profile, err := profileServiceTest.UpdateProfile(context.Background(), request.UpdateProfileRequest{ProfileCode: 2})
		assert.Nil(t, profile)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to update profile:")
	})
}

func TestDeletePhoto(t *testing.T) {
	t.Run("SuccessDeletePhotoByCode", func(t *testing.T) {
		profileRepository.Mock.On("DeletePhotoByCode", context.Background(), 1).Return(&models.DefaultResponse{ProfileCode: 1}, nil)

		result, err := profileServiceTest.DeletePhotoByCode(context.Background(), 1)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.ProfileCode)
	})
	t.Run("FailedDeletePhotoByCode", func(t *testing.T) {
		// program mock
		profileRepository.Mock.On("DeletePhotoByCode", mock.Anything, 2).Return(nil, errors.New("a"))

		profile, err := profileServiceTest.DeletePhotoByCode(context.Background(), 2)
		assert.Nil(t, profile)
		// assert.Equal(t, 0, profile.ProfileCode)
		assert.NotNil(t, err)
		// assert.Contains(t, err.Error(), "failed to delete photo:")
	})
}

func TestUploadPhoto(t *testing.T) {
	t.Run("SuccessUploadPhotoByCode", func(t *testing.T) {
		profileRepository.Mock.On("UpdateProfile", context.Background(), 1, &models.Profile{
			PhotoUrl: fmt.Sprintf("public/image/1-%d.png", time.Now().Unix())}).Return(&models.ProfileDTO{ProfileCode: 1, PhotoUrl: fmt.Sprintf("public/image/1-%d.png", time.Now().Unix())}, nil)

		result, err := profileServiceTest.UploadPhotoByCode(context.Background(), request.UploadPhotoRequest{
			ProfileCode: 1,
			Base64Img:   "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAApgAAAKYB3X3/OAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAANCSURBVEiJtZZPbBtFFMZ/M7ubXdtdb1xSFyeilBapySVU8h8OoFaooFSqiihIVIpQBKci6KEg9Q6H9kovIHoCIVQJJCKE1ENFjnAgcaSGC6rEnxBwA04Tx43t2FnvDAfjkNibxgHxnWb2e/u992bee7tCa00YFsffekFY+nUzFtjW0LrvjRXrCDIAaPLlW0nHL0SsZtVoaF98mLrx3pdhOqLtYPHChahZcYYO7KvPFxvRl5XPp1sN3adWiD1ZAqD6XYK1b/dvE5IWryTt2udLFedwc1+9kLp+vbbpoDh+6TklxBeAi9TL0taeWpdmZzQDry0AcO+jQ12RyohqqoYoo8RDwJrU+qXkjWtfi8Xxt58BdQuwQs9qC/afLwCw8tnQbqYAPsgxE1S6F3EAIXux2oQFKm0ihMsOF71dHYx+f3NND68ghCu1YIoePPQN1pGRABkJ6Bus96CutRZMydTl+TvuiRW1m3n0eDl0vRPcEysqdXn+jsQPsrHMquGeXEaY4Yk4wxWcY5V/9scqOMOVUFthatyTy8QyqwZ+kDURKoMWxNKr2EeqVKcTNOajqKoBgOE28U4tdQl5p5bwCw7BWquaZSzAPlwjlithJtp3pTImSqQRrb2Z8PHGigD4RZuNX6JYj6wj7O4TFLbCO/Mn/m8R+h6rYSUb3ekokRY6f/YukArN979jcW+V/S8g0eT/N3VN3kTqWbQ428m9/8k0P/1aIhF36PccEl6EhOcAUCrXKZXXWS3XKd2vc/TRBG9O5ELC17MmWubD2nKhUKZa26Ba2+D3P+4/MNCFwg59oWVeYhkzgN/JDR8deKBoD7Y+ljEjGZ0sosXVTvbc6RHirr2reNy1OXd6pJsQ+gqjk8VWFYmHrwBzW/n+uMPFiRwHB2I7ih8ciHFxIkd/3Omk5tCDV1t+2nNu5sxxpDFNx+huNhVT3/zMDz8usXC3ddaHBj1GHj/As08fwTS7Kt1HBTmyN29vdwAw+/wbwLVOJ3uAD1wi/dUH7Qei66PfyuRj4Ik9is+hglfbkbfR3cnZm7chlUWLdwmprtCohX4HUtlOcQjLYCu+fzGJH2QRKvP3UNz8bWk1qMxjGTOMThZ3kvgLI5AzFfo379UAAAAASUVORK5CYII=",
		})
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.ProfileCode)
	})
	t.Run("FailedUploadPhoto_DecodeToString", func(t *testing.T) {
		// program mock
		profileRepository.Mock.On("UpdateProfile", mock.Anything, 2, &models.Profile{}).Return(nil, errors.New("a"))

		result, err := profileServiceTest.UploadPhotoByCode(context.Background(), request.UploadPhotoRequest{
			ProfileCode: 2,
			Base64Img:   "data:image/png;base64,i",
		})
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to decode base64 string:")
	})
	// t.Run("FailedUploadPhoto_FormatNotFound", func(t *testing.T) {
	// 	// program mock
	// 	profileRepository.Mock.On("UpdateProfile", context.Background(), 101, &models.Profile{}).Return(nil, errors.New(""))

	// 	_, err := profileServiceTest.UploadPhotoByCode(context.Background(), request.UploadPhotoRequest{
	// 		ProfileCode: 101,
	// 		Base64Img:   "data:image/wep;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAApgAAAKYB3X3/OAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAANCSURBVEiJtZZPbBtFFMZ/M7ubXdtdb1xSFyeilBapySVU8h8OoFaooFSqiihIVIpQBKci6KEg9Q6H9kovIHoCIVQJJCKE1ENFjnAgcaSGC6rEnxBwA04Tx43t2FnvDAfjkNibxgHxnWb2e/u992bee7tCa00YFsffekFY+nUzFtjW0LrvjRXrCDIAaPLlW0nHL0SsZtVoaF98mLrx3pdhOqLtYPHChahZcYYO7KvPFxvRl5XPp1sN3adWiD1ZAqD6XYK1b/dvE5IWryTt2udLFedwc1+9kLp+vbbpoDh+6TklxBeAi9TL0taeWpdmZzQDry0AcO+jQ12RyohqqoYoo8RDwJrU+qXkjWtfi8Xxt58BdQuwQs9qC/afLwCw8tnQbqYAPsgxE1S6F3EAIXux2oQFKm0ihMsOF71dHYx+f3NND68ghCu1YIoePPQN1pGRABkJ6Bus96CutRZMydTl+TvuiRW1m3n0eDl0vRPcEysqdXn+jsQPsrHMquGeXEaY4Yk4wxWcY5V/9scqOMOVUFthatyTy8QyqwZ+kDURKoMWxNKr2EeqVKcTNOajqKoBgOE28U4tdQl5p5bwCw7BWquaZSzAPlwjlithJtp3pTImSqQRrb2Z8PHGigD4RZuNX6JYj6wj7O4TFLbCO/Mn/m8R+h6rYSUb3ekokRY6f/YukArN979jcW+V/S8g0eT/N3VN3kTqWbQ428m9/8k0P/1aIhF36PccEl6EhOcAUCrXKZXXWS3XKd2vc/TRBG9O5ELC17MmWubD2nKhUKZa26Ba2+D3P+4/MNCFwg59oWVeYhkzgN/JDR8deKBoD7Y+ljEjGZ0sosXVTvbc6RHirr2reNy1OXd6pJsQ+gqjk8VWFYmHrwBzW/n+uMPFiRwHB2I7ih8ciHFxIkd/3Omk5tCDV1t+2nNu5sxxpDFNx+huNhVT3/zMDz8usXC3ddaHBj1GHj/As08fwTS7Kt1HBTmyN29vdwAw+/wbwLVOJ3uAD1wi/dUH7Qei66PfyuRj4Ik9is+hglfbkbfR3cnZm7chlUWLdwmprtCohX4HUtlOcQjLYCu+fzGJH2QRKvP3UNz8bWk1qMxjGTOMThZ3kvgLI5AzFfo379UAAAAASUVORK5CYII=",
	// 	})
	// 	// assert.Nil(t, result)
	// 	// assert.NotNil(t, err)
	// 	assert.Contains(t, err.Error(), "sql: no rows in result set")
	// })
}

func TestDownloadPhoto(t *testing.T) {
	t.Run("SuccessDownloadPhotoByCode", func(t *testing.T) {
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 8).Return(&models.ProfileDTO{
			WantedJobTitle: "test",
			FirstName:      "test",
			Email:          "test",
			Phone:          "0888888889",
			Country:        "test",
			City:           "test",
			Address:        "test",
			PhotoUrl:       "public/image/1-1730888286.png"}, nil)

		result, err := profileServiceTest.DownloadPhotoByCode(context.Background(), 8)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		// assert.Contains(t, result, "iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAADJUlEQVR4nLSWTUwbRxTH/zPeZdc2ZjEFtwaVUqgEXFokG/fQqlXVVlRCVP2QWglVqD1RtT3QSr23B7iWS6NwShRFSImUKIqUAwrH5BAwUsglipQPQmKcmIAx+JO1Z6IZy47tXWODlL+02vG8t795b+btWyucc9gpOvnnN0Tlvyrugqb2ZExjLJYFQVAaOcKJJZ9uRpxqPuXIcZOc9i/+d8WOQ0oLRKenXUpS7+lqzW7Ecq4fmInzlY6ej3bh/jAux6lbXhzc7KgCURU/+bT0xe2k3pdvzUb8Cwvp8gKxyZkvGSGXBAeU71CNGyxDldpoOn/ZlPcXZ3otkVIny7McSYCRNwAcUM6/9y3OXyfRyb8+A9gSANV2ryrU8WNE3ncv9DRyFTIBOqYALNAMXKZLm/EqSxVsBYR4UOegq+mAo90sj9HEI4JNOcEymlihpTsL6izIS4wbi3PBpv6/N+55PtlljdxdIwnbcT0JpmBTmIVRd3DP4fl0B0SxT0QfSkIfTL76PZiUc3YSDMESTMFWQFgQnMAd2IPWn0JqxYvchgss5ZAPODx5GJ9vW0BizozoKBwUq5m6C9D60nCH4lC8pbNiQQWMBOShAdJgjMXk2IxpOHzkgvp2BkSz7qCYM8afw3ziRMu7aai+nDUdRgKEhye2APgbburJFD1eZZ9AVDSu10bnCCugfA2cTNTa/j+3ggeP4/C26Wg3dHgNJ7yGLm3xRBbxRAZ74r6fxcA7Xvw+FbIJn68p4DRs91puRhJIpQ/l9fTZ/pGBCl/7DGiYQnWsFhtTtYYGOo+ENuFrCjbF+5dj4GSu1vrtV8No82gN4cJH+FpE+KxgF6uIvDULYL3S3t6m44+pEN7sdNeFC5vwEb41WgfvnkPlFw2r4yOgjpXa1p3PMyzfeIi797exuVXc695uA8PvdeGLj/uhKJZKN8EKIYxeu129gNDa178BmG/2+2AjcZYzCFw9VZoglo9+MZOzAD44JnwdrPBzKfKSrG+ydPCPgpN/7arLNmrC/5HP1MDtM6jUne98ouXKjltsiuW/LcUXlIZlmYtKrKOXAQAA//8+2DMY6mBorgAAAABJRU5ErkJggg==")
		// assert.Equal(t, "public/image/1-1730888286.png", result)
	})
	t.Run("FailedDownloadPhoto_UserNotFound", func(t *testing.T) {
		// program mock
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 12).Return(&models.ProfileDTO{}, errors.New("sql: no rows in result set"))

		result, err := profileServiceTest.DownloadPhotoByCode(context.Background(), 12)
		assert.Equal(t, "", result)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "sql: no rows in result set")
		// assert.Contains(t, err.Error(), "failed to open image file")
	})
	t.Run("FailedDownloadPhoto_FileNotFound", func(t *testing.T) {
		// program mock
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 10).Return(&models.ProfileDTO{ProfileCode: 10,
			WantedJobTitle: "test",
			FirstName:      "test",
			Email:          "test",
			Phone:          "0888888889",
			Country:        "test",
			City:           "test",
			Address:        "test",
			PhotoUrl:       "publizc/image/2-1730888286.png"}, nil)

		result, err := profileServiceTest.DownloadPhotoByCode(context.Background(), 10)
		assert.Equal(t, "", result)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
		assert.Contains(t, err.Error(), "failed to open image file")
	})

	t.Run("FailedDownloadPhoto_FailedToDecode", func(t *testing.T) {
		// program mock
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 3).Return(&models.ProfileDTO{ProfileCode: 2, PhotoUrl: "public/image/asdaa.webp"}, nil)

		result, err := profileServiceTest.DownloadPhotoByCode(context.Background(), 3)
		assert.Equal(t, "", result)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to decode image:")
	})

	// t.Run("FailedDownloadPhoto_FailedToEncode", func(t *testing.T) {
	// 	// program mock
	// 	profileRepository.Mock.On("GetProfileByCode", mock.Anything, 2).Return(&models.ProfileDTO{ProfileCode: 2, PhotoUrl: "public/image/as.jpg"}, nil)

	// 	result, err := profileServiceTest.DownloadPhotoByCode(context.Background(), 2)
	// 	assert.Equal(t, "", result)
	// 	assert.NotNil(t, err)
	// 	assert.Contains(t, err.Error(), "failed to encode image:")
	// })
}
