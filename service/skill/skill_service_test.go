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

var skillRepository = &repository.SkillRepository{Mock: mock.Mock{}}
var skillServiceTest = skillService{skillRepo: skillRepository}

func TestInitSkillService(t *testing.T) {
	t.Run("SuccessInitSkillService", func(t *testing.T) {
		assert.NotNil(t, NewSkillService(skillRepository))
	})
}

func TestGetSkill(t *testing.T) {
	t.Run("SuccessGetSkill", func(t *testing.T) {
		var response = []*models.SkillDTO{}
		result := &models.SkillDTO{
			Id:    1,
			Skill: "Golang",
			Level: "Beginner",
		}
		response = append(response, result)
		skillRepository.Mock.On("GetSkillsByProfileCode", mock.Anything, 1).Return(response, nil)

		skills, err := skillServiceTest.GetSkillsByCode(context.Background(), 1)
		assert.Nil(t, err)
		assert.NotNil(t, skills)
	})
	t.Run("FailedGetSkill", func(t *testing.T) {
		// program mock
		skillRepository.Mock.On("GetSkillsByProfileCode", mock.Anything, 2).Return(nil, errors.New("sql: no rows in result set"))

		skill, err := skillServiceTest.GetSkillsByCode(context.Background(), 2)
		assert.Nil(t, skill)
		assert.NotNil(t, err)
		// assert.Equal(t, "failed to get profile: sql: no rows in result set", err.Error())
		assert.Contains(t, err.Error(), "failed to get skills")
	})
}

func TestCreateSkill(t *testing.T) {
	t.Run("SuccessCreateSkill", func(t *testing.T) {
		result := &models.SkillDTO{
			Id:    1,
			Skill: "Golang",
			Level: "Beginner",
		}
		skillRepository.Mock.On("CreateSkill", mock.Anything, &models.Skill{
			ProfileCode: 1,
			Skill:       "Golang",
			Level:       "Beginner",
		}).Return(result, nil)

		skills, err := skillServiceTest.CreateSkill(context.Background(), request.CreateSkillRequest{
			ProfileCode: 1,
			Skill:       "Golang",
			Level:       "Beginner",
		})
		assert.Nil(t, err)
		assert.NotNil(t, skills)
	})
	t.Run("FailedCreateSkill", func(t *testing.T) {
		// program mock
		skillRepository.Mock.On("CreateSkill", mock.Anything, &models.Skill{
			ProfileCode: 2,
			Skill:       "Golang",
			Level:       "Beginner",
		}).Return(nil, errors.New(""))

		skill, err := skillServiceTest.CreateSkill(context.Background(), request.CreateSkillRequest{
			ProfileCode: 2,
			Skill:       "Golang",
			Level:       "Beginner",
		})
		assert.Nil(t, skill)
		assert.NotNil(t, err)
		// assert.Equal(t, "failed to get profile: sql: no rows in result set", err.Error())
		assert.Contains(t, err.Error(), "failed to create skill")
	})
}

func TestDeleteSkill(t *testing.T) {
	t.Run("SuccessDeleteSkill", func(t *testing.T) {
		skillRepository.Mock.On("DeleteSkill", mock.Anything, 1, 1).Return(nil)

		skills, err := skillServiceTest.DeleteSkill(context.Background(), 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, skills)
	})
	t.Run("FailedCreateSkill", func(t *testing.T) {
		// program mock
		skillRepository.Mock.On("DeleteSkill", mock.Anything, 1, 2).Return(errors.New(""))

		skill, err := skillServiceTest.DeleteSkill(context.Background(), 1, 2)
		assert.Nil(t, skill)
		assert.NotNil(t, err)
		// assert.Equal(t, "failed to get profile: sql: no rows in result set", err.Error())
		assert.Contains(t, err.Error(), "failed to delete skill")
	})
}
