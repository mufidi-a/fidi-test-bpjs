package service

import (
	"context"
	"errors"
	"test-bpjs/v2/models"
	"test-bpjs/v2/models/request"
	repository "test-bpjs/v2/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var educationRepository = &repository.EducationRepository{Mock: mock.Mock{}}
var educationServiceTest = educationService{educationRepo: educationRepository}

func TestInitEducationService(t *testing.T) {
	t.Run("SuccessInitEducationService", func(t *testing.T) {
		assert.NotNil(t, NewEducationService(educationRepository))
	})
}

func TestGetEducation(t *testing.T) {
	t.Run("SuccessGetEducation", func(t *testing.T) {
		var response = []*models.EducationDTO{}
		result := &models.EducationDTO{
			Id:          1,
			School:      "UGM",
			Degree:      "S1",
			City:        "Jogja",
			Description: "I'm Programmer",
		}
		response = append(response, result)
		educationRepository.Mock.On("GetEducationByProfileCode", mock.Anything, 1).Return(response, nil)

		education, err := educationServiceTest.GetEducationByCode(context.Background(), 1)
		assert.Nil(t, err)
		assert.NotNil(t, education)
	})
	t.Run("FailedGetEducation", func(t *testing.T) {
		// program mock
		educationRepository.Mock.On("GetEducationByProfileCode", mock.Anything, 2).Return(nil, errors.New("sql: no rows in result set"))

		education, err := educationServiceTest.GetEducationByCode(context.Background(), 2)
		assert.Nil(t, education)
		assert.NotNil(t, err)
		// assert.Equal(t, "failed to get profile: sql: no rows in result set", err.Error())
		assert.Contains(t, err.Error(), "failed to get education")
	})
}

func TestCreateEducation(t *testing.T) {
	t.Run("SuccessCreateEducation", func(t *testing.T) {
		result := &models.EducationDTO{
			Id:          1,
			School:      "UGM",
			Degree:      "S1",
			City:        "Jogja",
			Description: "I'm Programmer",
		}
		educationRepository.Mock.On("CreateEducation", mock.Anything, &models.Education{
			ProfileCode: 1,
			School:      "UGM",
			Degree:      "S1",
			City:        "Jogja",
			Description: "I'm Programmer",
		}).Return(result, nil)

		education, err := educationServiceTest.CreateEducation(context.Background(), request.CreateEducationRequest{
			ProfileCode: 1,
			School:      "UGM",
			Degree:      "S1",
			City:        "Jogja",
			Description: "I'm Programmer",
		})
		assert.Nil(t, err)
		assert.NotNil(t, education)
	})
	t.Run("FailedCreateEducation", func(t *testing.T) {
		// program mock
		educationRepository.Mock.On("CreateEducation", mock.Anything, &models.Education{
			ProfileCode: 2,
			School:      "UGM",
			Degree:      "S1",
			City:        "Jakarta",
			Description: "I'm Programmer",
		}).Return(nil, errors.New(""))

		education, err := educationServiceTest.CreateEducation(context.Background(), request.CreateEducationRequest{
			ProfileCode: 2,
			School:      "UGM",
			Degree:      "S1",
			City:        "Jakarta",
			Description: "I'm Programmer",
		})
		assert.Nil(t, education)
		assert.NotNil(t, err)
		// assert.Equal(t, "failed to get profile: sql: no rows in result set", err.Error())
		assert.Contains(t, err.Error(), "failed to create education")
	})
}

func TestDeleteEducation(t *testing.T) {
	t.Run("SuccessDeleteEducation", func(t *testing.T) {
		educationRepository.Mock.On("DeleteEducation", mock.Anything, 1, 1).Return(nil)

		education, err := educationServiceTest.DeleteEducation(context.Background(), 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, education)
	})
	t.Run("FailedDeleteEducation", func(t *testing.T) {
		// program mock
		educationRepository.Mock.On("DeleteEducation", mock.Anything, 1, 2).Return(errors.New(""))

		education, err := educationServiceTest.DeleteEducation(context.Background(), 1, 2)
		assert.Nil(t, education)
		assert.NotNil(t, err)
		// assert.Equal(t, "failed to get profile: sql: no rows in result set", err.Error())
		assert.Contains(t, err.Error(), "failed to delete education")
	})
}
