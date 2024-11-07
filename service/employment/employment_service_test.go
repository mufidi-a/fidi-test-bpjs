package service

import (
	"context"
	"errors"
	"test-bpjs/v2/models"
	"test-bpjs/v2/models/request"
	repository "test-bpjs/v2/repository/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var employmentRepository = &repository.EmploymentRepository{Mock: mock.Mock{}}
var employmentServiceTest = employmentService{employmentRepo: employmentRepository}

func TestInitEmploymentService(t *testing.T) {
	t.Run("SuccessInitEmploymentService", func(t *testing.T) {
		assert.NotNil(t, NewEmploymentService(employmentRepository))
	})
}

func TestGetEmployment(t *testing.T) {
	t.Run("SuccessGetEmployment", func(t *testing.T) {
		var response = []*models.EmploymentDTO{}
		result := &models.EmploymentDTO{
			Id:          1,
			JobTitle:    "Programmer",
			Employer:    "PT. ABC",
			StartDate:   time.Now(),
			EndDate:     time.Now(),
			City:        "Jakarta",
			Description: "I'm Programmer",
		}
		response = append(response, result)
		employmentRepository.Mock.On("GetEmploymentByProfileCode", mock.Anything, 1).Return(response, nil)

		employment, err := employmentServiceTest.GetEmploymentByCode(context.Background(), 1)
		assert.Nil(t, err)
		assert.NotNil(t, employment)
	})
	t.Run("FailedGetEmployment", func(t *testing.T) {
		// program mock
		employmentRepository.Mock.On("GetEmploymentByProfileCode", mock.Anything, 2).Return(nil, errors.New("sql: no rows in result set"))

		skill, err := employmentServiceTest.GetEmploymentByCode(context.Background(), 2)
		assert.Nil(t, skill)
		assert.NotNil(t, err)
		// assert.Equal(t, "failed to get profile: sql: no rows in result set", err.Error())
		assert.Contains(t, err.Error(), "failed to get employment")
	})
}

func TestCreateEmployment(t *testing.T) {
	t.Run("SuccessCreateEmployment", func(t *testing.T) {
		result := &models.EmploymentDTO{
			Id:          1,
			JobTitle:    "Programmer",
			Employer:    "PT. ABC",
			City:        "Jakarta",
			Description: "I'm Programmer",
		}
		employmentRepository.Mock.On("CreateEmployment", mock.Anything, &models.Employment{
			ProfileCode: 1,
			JobTitle:    "Programmer",
			Employer:    "PT. ABC",
			City:        "Jakarta",
			Description: "I'm Programmer",
		}).Return(result, nil)

		skills, err := employmentServiceTest.CreateEmployment(context.Background(), request.CreateEmploymentRequest{
			ProfileCode: 1,
			JobTitle:    "Programmer",
			Employer:    "PT. ABC",
			City:        "Jakarta",
			Description: "I'm Programmer",
		})
		assert.Nil(t, err)
		assert.NotNil(t, skills)
	})
	t.Run("FailedCreateEmployment", func(t *testing.T) {
		// program mock
		employmentRepository.Mock.On("CreateEmployment", mock.Anything, &models.Employment{
			ProfileCode: 2,
			JobTitle:    "Programmer",
			Employer:    "PT. ABC",
			City:        "Jakarta",
			Description: "I'm Programmer",
		}).Return(nil, errors.New(""))

		skill, err := employmentServiceTest.CreateEmployment(context.Background(), request.CreateEmploymentRequest{
			ProfileCode: 2,
			JobTitle:    "Programmer",
			Employer:    "PT. ABC",
			City:        "Jakarta",
			Description: "I'm Programmer",
		})
		assert.Nil(t, skill)
		assert.NotNil(t, err)
		// assert.Equal(t, "failed to get profile: sql: no rows in result set", err.Error())
		assert.Contains(t, err.Error(), "failed to create employment")
	})
}

func TestDeleteEmployment(t *testing.T) {
	t.Run("SuccessDeleteEmployment", func(t *testing.T) {
		employmentRepository.Mock.On("DeleteEmployment", mock.Anything, 1, 1).Return(nil)

		skills, err := employmentServiceTest.DeleteEmployment(context.Background(), 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, skills)
	})
	t.Run("FailedCreateEmployment", func(t *testing.T) {
		// program mock
		employmentRepository.Mock.On("DeleteEmployment", mock.Anything, 1, 2).Return(errors.New(""))

		skill, err := employmentServiceTest.DeleteEmployment(context.Background(), 1, 2)
		assert.Nil(t, skill)
		assert.NotNil(t, err)
		// assert.Equal(t, "failed to get profile: sql: no rows in result set", err.Error())
		assert.Contains(t, err.Error(), "failed to delete employment")
	})
}
