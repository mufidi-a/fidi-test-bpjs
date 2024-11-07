package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"test-bpjs/v2/models"
	"test-bpjs/v2/models/response"
	repository "test-bpjs/v2/repository/mocks"
	educationService "test-bpjs/v2/service/education"
	employmentService "test-bpjs/v2/service/employment"
	profileService "test-bpjs/v2/service/profile"
	skillService "test-bpjs/v2/service/skill"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var profileRepository = &repository.ProfileRepository{Mock: mock.Mock{}}
var profileServiceTest = profileService.NewProfileService(profileRepository)
var skillRepository = &repository.SkillRepository{Mock: mock.Mock{}}
var skillServiceTest = skillService.NewSkillService(skillRepository)
var educationRepository = &repository.EducationRepository{Mock: mock.Mock{}}
var educationServiceTest = educationService.NewEducationService(educationRepository)
var employmentRepository = &repository.EmploymentRepository{Mock: mock.Mock{}}
var employmentServiceTest = employmentService.NewEmploymentService(employmentRepository)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

func TestGetProfileController(t *testing.T) {
	t.Run("SuccessGetProfileController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		result := &models.ProfileDTO{
			ProfileCode: 1,
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("1")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 1).Return(result, nil)

		controller := apiHandler.GetProfileByCode()(c)
		if assert.NoError(t, controller) {
			bodyResponses := rec.Body.String()
			var response response.CreateProfileResponse
			err := json.Unmarshal([]byte(bodyResponses), &response)
			if err != nil {
				assert.Error(t, err, "error")
			}
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedGetProfileController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 2).Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("2")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetProfileByCode()(c)
		if assert.Error(t, controller) {

			var errCode int

			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])

			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedGetProfileController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 0).Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetProfileByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedGetProfileController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, "asd").Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetProfileByCode()(c)
		if assert.Error(t, controller) {

			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestCreateProfileController(t *testing.T) {
	t.Run("SuccessCreateProfileController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		result := &models.ProfileDTO{
			ProfileCode: 4,
		}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"wantedJobTitle": "Software Engineer",
			"firstName":      "Namaku",
			"lastName":       "Ukaman",
			"email":          "ukaman.namaku@gmail.com",
			"phone":          "08008880000",
			"country":        "Indonesia",
			"city":           "Jakarta",
			"address":        "Jl. Gatot Subroto",
			"postalCode":     200001,
			"drivingLicense": "1234567890123456",
			"nationality":    "Indonesia",
			"placeOfBirth":   "Maluku",
			"dateOfBirth":    "2006-01-02T00:00:00Z",
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		dob, _ := time.Parse("2006-01-02T15:04:05Z", "2006-01-02T00:00:00Z")
		profileRepository.Mock.On("CreateProfile", mock.Anything, &models.Profile{
			WantedJobTitle: "Software Engineer",
			FirstName:      "Namaku",
			LastName:       "Ukaman",
			Email:          "ukaman.namaku@gmail.com",
			Phone:          "08008880000",
			Country:        "Indonesia",
			City:           "Jakarta",
			Address:        "Jl. Gatot Subroto",
			PostalCode:     200001,
			DrivingLicense: "1234567890123456",
			Nationality:    "Indonesia",
			PlaceOfBirth:   "Maluku",
			DateOfBirth:    dob,
		}).Return(result, nil)

		controller := apiHandler.CreateProfile()(c)
		if assert.NoError(t, controller) {
			bodyResponses := rec.Body.String()
			var response response.DefaultResponse
			err := json.Unmarshal([]byte(bodyResponses), &response)
			if err != nil {
				assert.Error(t, err, "error")
			}
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedGetProfileController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		requestBody, _ := json.Marshal(map[string]interface{}{
			"firstName":      "Namaku",
			"lastName":       "Ukaman",
			"email":          "ukaman.namaku@gmail.com",
			"phone":          "08008880000",
			"country":        "Indonesia",
			"city":           "Jakarta",
			"address":        "Jl. Gatot Subroto",
			"postalCode":     200001,
			"drivingLicense": "1234567890123456",
			"nationality":    "Indonesia",
			"placeOfBirth":   "Maluku",
			"dateOfBirth":    "2006-01-02T00:00:00Z",
		})

		dob, _ := time.Parse("2006-01-02T15:04:05Z", "2006-01-02T00:00:00Z")
		profileRepository.Mock.On("CreateProfile", mock.Anything, &models.Profile{
			FirstName:      "Namaku",
			LastName:       "Ukaman",
			Email:          "ukaman.namaku@gmail.com",
			Phone:          "08008880000",
			Country:        "Indonesia",
			City:           "Jakarta",
			Address:        "Jl. Gatot Subroto",
			PostalCode:     200001,
			DrivingLicense: "1234567890123456",
			Nationality:    "Indonesia",
			PlaceOfBirth:   "Maluku",
			DateOfBirth:    dob,
		}).Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.CreateProfile()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedGetProfileController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		requestBody, _ := json.Marshal(map[string]interface{}{
			"wantedJobTitle": "Software Engineer",
			"firstName":      "Namaku",
			"lastName":       "Ukaman",
			"email":          "ukaman.namaku@gmail.com",
			"phone":          "08008880000",
			"country":        "Indonesia",
			"city":           "Jakarta",
			"address":        "Jl. Gatot Subroto",
			"postalCode":     200001,
			"drivingLicense": "1234567890123456",
			"nationality":    "Indonesia",
			"placeOfBirth":   "Maluku",
			"dateOfBirth":    "2006-01-02",
		})
		dob, _ := time.Parse("2006-01-02T15:04:05Z", "2006-01-02")
		profileRepository.Mock.On("CreateProfile", mock.Anything, &models.Profile{
			WantedJobTitle: "",
			FirstName:      "Namaku",
			LastName:       "Ukaman",
			Email:          "ukaman.namaku@gmail.com",
			Phone:          "08008880000",
			Country:        "Indonesia",
			City:           "Jakarta",
			Address:        "Jl. Gatot Subroto",
			PostalCode:     200001,
			DrivingLicense: "1234567890123456",
			Nationality:    "Indonesia",
			PlaceOfBirth:   "Maluku",
			DateOfBirth:    dob,
		}).Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.CreateProfile()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})

	// t.Run("FailedGetProfileController_ErrBind", func(t *testing.T) {
	// 	e := echo.New()
	// 	e.Validator = &CustomValidator{validator: validator.New()}
	// 	profileRepository.Mock.On("GetProfileByCode", mock.Anything, "asd").Return(nil, errors.New("sql: no rows in result set"))

	// 	req := httptest.NewRequest(http.MethodPost, "/api", nil)
	// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// 	rec := httptest.NewRecorder()
	// 	c := e.NewContext(req, rec)
	// 	c.SetPath("/profile")
	// 	c.SetParamNames("profileCode")
	// 	c.SetParamValues("asd")

	// 	apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

	// 	// apiHandler.GetProfileByCode()(c)
	// 	controller := apiHandler.GetProfileByCode()(c)
	// 	if assert.Error(t, controller) {
	// 		bodyResponses := rec.Body.String()
	// 		var response response.CreateProfileResponse
	// 		var errCode int
	// 		var errMsg string
	// 		err := json.Unmarshal([]byte(bodyResponses), &response)
	// 		if err != nil {
	// 			assert.Error(t, err, "error")
	// 		}
	// 		if controller.Error() != "" {
	// 			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
	// 			match := re.FindStringSubmatch(controller.Error())
	// 			errCode, _ = strconv.Atoi(match[1])
	// 			errMsg = match[2]
	// 		}
	// 		assert.Equal(t, http.StatusBadRequest, errCode)
	// 		assert.Equal(t, "bad request. failed to bind", errMsg)
	// 	}

	// })
}

func TestUpdateProfileController(t *testing.T) {
	t.Run("SuccessUpdateProfileController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		result := &models.ProfileDTO{
			ProfileCode: 4,
		}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"wantedJobTitle": "Software Engineer",
			"firstName":      "Namaku",
			"lastName":       "Ukaman",
			"email":          "ukaman.namaku@gmail.com",
			"phone":          "08008880000",
			"country":        "Indonesia",
			"city":           "Jakarta",
			"address":        "Jl. Gatot Subroto",
			"postalCode":     200001,
			"drivingLicense": "1234567890123456",
			"nationality":    "Indonesia",
			"placeOfBirth":   "Maluku",
			"dateOfBirth":    "2006-01-02T00:00:00Z",
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("5")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		dob, _ := time.Parse("2006-01-02T15:04:05Z", "2006-01-02T00:00:00Z")
		profileRepository.Mock.On("UpdateProfile", mock.Anything, 5, &models.Profile{
			WantedJobTitle: "Software Engineer",
			FirstName:      "Namaku",
			LastName:       "Ukaman",
			Email:          "ukaman.namaku@gmail.com",
			Phone:          "08008880000",
			Country:        "Indonesia",
			City:           "Jakarta",
			Address:        "Jl. Gatot Subroto",
			PostalCode:     200001,
			DrivingLicense: "1234567890123456",
			Nationality:    "Indonesia",
			PlaceOfBirth:   "Maluku",
			DateOfBirth:    dob,
		}).Return(result, nil)

		controller := apiHandler.UpdateProfile()(c)
		if assert.NoError(t, controller) {
			bodyResponses := rec.Body.String()
			var response response.DefaultResponse
			err := json.Unmarshal([]byte(bodyResponses), &response)
			if err != nil {
				assert.Error(t, err, "error")
			}
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedGetUpdateController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		requestBody, _ := json.Marshal(map[string]interface{}{})
		profileRepository.Mock.On("UpdateProfile", mock.Anything, 6, &models.Profile{}).Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodPut, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("6")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.UpdateProfile()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedGetProfileController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("UpdateProfile", mock.Anything, "0", &models.Profile{}).Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodPut, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.UpdateProfile()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedGetProfileController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		requestBody, _ := json.Marshal(map[string]interface{}{
			"wantedJobTitle": "Software Engineer",
			"firstName":      "Namaku",
			"lastName":       "Ukaman",
			"email":          "ukaman.namaku@gmail.com",
			"phone":          "08008880000",
			"country":        "Indonesia",
			"city":           "Jakarta",
			"address":        "Jl. Gatot Subroto",
			"postalCode":     200001,
			"drivingLicense": "1234567890123456",
			"nationality":    "Indonesia",
			"placeOfBirth":   "Maluku",
			"dateOfBirth":    "2006-01-02",
		})
		dob, _ := time.Parse("2006-01-02T15:04:05Z", "2006-01-02")
		profileRepository.Mock.On("UpdateProfile", mock.Anything, "asd", &models.Profile{
			WantedJobTitle: "Software Engineer",
			FirstName:      "Namaku",
			LastName:       "Ukaman",
			Email:          "ukaman.namaku@gmail.com",
			Phone:          "08008880000",
			Country:        "Indonesia",
			City:           "Jakarta",
			Address:        "Jl. Gatot Subroto",
			PostalCode:     200001,
			DrivingLicense: "1234567890123456",
			Nationality:    "Indonesia",
			PlaceOfBirth:   "Maluku",
			DateOfBirth:    dob,
		}).Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodPut, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.UpdateProfile()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestDownloadPhotoController(t *testing.T) {
	t.Run("SuccessDownloadPhotoController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/photo")
		c.SetParamNames("profileCode")
		c.SetParamValues("8")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 8).Return(&models.ProfileDTO{
			WantedJobTitle: "test",
			FirstName:      "test",
			Email:          "test",
			Phone:          "0888888889",
			Country:        "test",
			City:           "test",
			Address:        "test",
			PhotoUrl:       "../public/image/1-1730888286.png"}, nil)

		controller := apiHandler.DownloadPhoto()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedDownloadPhotoController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		requestBody, _ := json.Marshal(map[string]interface{}{})
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 9).Return(&models.ProfileDTO{
			WantedJobTitle: "test",
			FirstName:      "test",
			Email:          "test",
			Phone:          "0888888889",
			Country:        "test",
			City:           "test",
			Address:        "test",
			PhotoUrl:       "../publizc/image/1-1730888286.png"}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("9")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DownloadPhoto()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedDownloadPhotoController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, 0).Return(&models.ProfileDTO{}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DownloadPhoto()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedDownloadPhotoController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("GetProfileByCode", mock.Anything, "asd").Return(&models.ProfileDTO{}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DownloadPhoto()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestUploadController(t *testing.T) {
	t.Run("SuccessUploadPhotoController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"base64img": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAApgAAAKYB3X3/OAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAANCSURBVEiJtZZPbBtFFMZ/M7ubXdtdb1xSFyeilBapySVU8h8OoFaooFSqiihIVIpQBKci6KEg9Q6H9kovIHoCIVQJJCKE1ENFjnAgcaSGC6rEnxBwA04Tx43t2FnvDAfjkNibxgHxnWb2e/u992bee7tCa00YFsffekFY+nUzFtjW0LrvjRXrCDIAaPLlW0nHL0SsZtVoaF98mLrx3pdhOqLtYPHChahZcYYO7KvPFxvRl5XPp1sN3adWiD1ZAqD6XYK1b/dvE5IWryTt2udLFedwc1+9kLp+vbbpoDh+6TklxBeAi9TL0taeWpdmZzQDry0AcO+jQ12RyohqqoYoo8RDwJrU+qXkjWtfi8Xxt58BdQuwQs9qC/afLwCw8tnQbqYAPsgxE1S6F3EAIXux2oQFKm0ihMsOF71dHYx+f3NND68ghCu1YIoePPQN1pGRABkJ6Bus96CutRZMydTl+TvuiRW1m3n0eDl0vRPcEysqdXn+jsQPsrHMquGeXEaY4Yk4wxWcY5V/9scqOMOVUFthatyTy8QyqwZ+kDURKoMWxNKr2EeqVKcTNOajqKoBgOE28U4tdQl5p5bwCw7BWquaZSzAPlwjlithJtp3pTImSqQRrb2Z8PHGigD4RZuNX6JYj6wj7O4TFLbCO/Mn/m8R+h6rYSUb3ekokRY6f/YukArN979jcW+V/S8g0eT/N3VN3kTqWbQ428m9/8k0P/1aIhF36PccEl6EhOcAUCrXKZXXWS3XKd2vc/TRBG9O5ELC17MmWubD2nKhUKZa26Ba2+D3P+4/MNCFwg59oWVeYhkzgN/JDR8deKBoD7Y+ljEjGZ0sosXVTvbc6RHirr2reNy1OXd6pJsQ+gqjk8VWFYmHrwBzW/n+uMPFiRwHB2I7ih8ciHFxIkd/3Omk5tCDV1t+2nNu5sxxpDFNx+huNhVT3/zMDz8usXC3ddaHBj1GHj/As08fwTS7Kt1HBTmyN29vdwAw+/wbwLVOJ3uAD1wi/dUH7Qei66PfyuRj4Ik9is+hglfbkbfR3cnZm7chlUWLdwmprtCohX4HUtlOcQjLYCu+fzGJH2QRKvP3UNz8bWk1qMxjGTOMThZ3kvgLI5AzFfo379UAAAAASUVORK5CYII=",
		})

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/photo")
		c.SetParamNames("profileCode")
		c.SetParamValues("10")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		profileRepository.Mock.On("UpdateProfile", mock.Anything, 10, &models.Profile{
			PhotoUrl: fmt.Sprintf("public/image/10-%d.png", time.Now().Unix())}).
			Return(&models.ProfileDTO{ProfileCode: 10, PhotoUrl: fmt.Sprintf("public/image/10-%d.png", time.Now().Unix())}, nil)

		controller := apiHandler.UploadPhoto()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedUploadPhotoController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		requestBody, _ := json.Marshal(map[string]interface{}{})
		profileRepository.Mock.On("UpdateProfile", mock.Anything, 11, &models.Profile{}).Return(nil, errors.New(""))

		req := httptest.NewRequest(http.MethodPut, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("11")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.UploadPhoto()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedUploadPhotoController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		requestBody, _ := json.Marshal(map[string]interface{}{})
		profileRepository.Mock.On("UpdateProfile", mock.Anything, 0, &models.Profile{}).Return(nil, errors.New("a"))

		req := httptest.NewRequest(http.MethodPut, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.UploadPhoto()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, "bad request. failed to validate", errMsg)
			assert.Equal(t, http.StatusBadRequest, errCode)
		}

	})

	t.Run("FailedUploadPhotoController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		requestBody, _ := json.Marshal(map[string]interface{}{
			"base64img": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAApgAAAKYB3X3/OAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAANCSURBVEiJtZZPbBtFFMZ/M7ubXdtdb1xSFyeilBapySVU8h8OoFaooFSqiihIVIpQBKci6KEg9Q6H9kovIHoCIVQJJCKE1ENFjnAgcaSGC6rEnxBwA04Tx43t2FnvDAfjkNibxgHxnWb2e/u992bee7tCa00YFsffekFY+nUzFtjW0LrvjRXrCDIAaPLlW0nHL0SsZtVoaF98mLrx3pdhOqLtYPHChahZcYYO7KvPFxvRl5XPp1sN3adWiD1ZAqD6XYK1b/dvE5IWryTt2udLFedwc1+9kLp+vbbpoDh+6TklxBeAi9TL0taeWpdmZzQDry0AcO+jQ12RyohqqoYoo8RDwJrU+qXkjWtfi8Xxt58BdQuwQs9qC/afLwCw8tnQbqYAPsgxE1S6F3EAIXux2oQFKm0ihMsOF71dHYx+f3NND68ghCu1YIoePPQN1pGRABkJ6Bus96CutRZMydTl+TvuiRW1m3n0eDl0vRPcEysqdXn+jsQPsrHMquGeXEaY4Yk4wxWcY5V/9scqOMOVUFthatyTy8QyqwZ+kDURKoMWxNKr2EeqVKcTNOajqKoBgOE28U4tdQl5p5bwCw7BWquaZSzAPlwjlithJtp3pTImSqQRrb2Z8PHGigD4RZuNX6JYj6wj7O4TFLbCO/Mn/m8R+h6rYSUb3ekokRY6f/YukArN979jcW+V/S8g0eT/N3VN3kTqWbQ428m9/8k0P/1aIhF36PccEl6EhOcAUCrXKZXXWS3XKd2vc/TRBG9O5ELC17MmWubD2nKhUKZa26Ba2+D3P+4/MNCFwg59oWVeYhkzgN/JDR8deKBoD7Y+ljEjGZ0sosXVTvbc6RHirr2reNy1OXd6pJsQ+gqjk8VWFYmHrwBzW/n+uMPFiRwHB2I7ih8ciHFxIkd/3Omk5tCDV1t+2nNu5sxxpDFNx+huNhVT3/zMDz8usXC3ddaHBj1GHj/As08fwTS7Kt1HBTmyN29vdwAw+/wbwLVOJ3uAD1wi/dUH7Qei66PfyuRj4Ik9is+hglfbkbfR3cnZm7chlUWLdwmprtCohX4HUtlOcQjLYCu+fzGJH2QRKvP3UNz8bWk1qMxjGTOMThZ3kvgLI5AzFfo379UAAAAASUVORK5CYII=",
		})
		profileRepository.Mock.On("UpdateProfile", mock.Anything, "asd", &models.Profile{
			PhotoUrl: fmt.Sprintf("public/image/12-%d.png", time.Now().Unix())}).
			Return(nil, errors.New(""))

		req := httptest.NewRequest(http.MethodPut, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.UploadPhoto()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestDeletePhotoController(t *testing.T) {
	t.Run("SuccessDeletePhotoController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/photo")
		c.SetParamNames("profileCode")
		c.SetParamValues("13")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		profileRepository.Mock.On("DeletePhotoByCode", mock.Anything, 13).
			Return(&models.DefaultResponse{ProfileCode: 13}, nil)

		controller := apiHandler.DeletePhoto()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedDeletePhotoController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("DeletePhotoByCode", mock.Anything, 14).
			Return(nil, errors.New(""))

		req := httptest.NewRequest(http.MethodDelete, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("14")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeletePhoto()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedDeletePhotoController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("DeletePhotoByCode", mock.Anything, 0).Return(nil, errors.New("a"))

		req := httptest.NewRequest(http.MethodDelete, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeletePhoto()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, "bad request. failed to validate", errMsg)
			assert.Equal(t, http.StatusBadRequest, errCode)
		}

	})

	t.Run("FailedDeletePhotoController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("DeletePhotoByCode", mock.Anything, "asd").Return(nil, errors.New("a"))

		req := httptest.NewRequest(http.MethodDelete, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeletePhoto()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestGetWorkingExperienceController(t *testing.T) {
	t.Run("SuccessGetWorkingExperienceController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		workingExperience := &models.ProfileDTO{
			WorkingExperience: "test",
		}
		profileRepository.Mock.On("GetWorkingExperienceByCode", mock.Anything, 15).Return(workingExperience, nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/working-experience")
		c.SetParamNames("profileCode")
		c.SetParamValues("15")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		controller := apiHandler.GetWorkingExperienceByCode()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedGetWorkingExperienceController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		requestBody, _ := json.Marshal(map[string]interface{}{})
		profileRepository.Mock.On("GetWorkingExperienceByCode", mock.Anything, 16).Return(nil, errors.New(""))

		req := httptest.NewRequest(http.MethodGet, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/working_experience")
		c.SetParamNames("profileCode")
		c.SetParamValues("16")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetWorkingExperienceByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedGetWorkingExperienceController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("GetWorkingExperienceByCode", mock.Anything, 0).Return(&models.ProfileDTO{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetWorkingExperienceByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedGetWorkingExperienceByCodeController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("GetWorkingExperienceByCode", mock.Anything, "asd").Return(&models.ProfileDTO{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetWorkingExperienceByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestUpdateWorkingExperienceController(t *testing.T) {
	t.Run("SuccessUpdateWorkingExperienceController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		result := &models.ProfileDTO{
			ProfileCode: 17,
		}
		requestBody, _ := json.Marshal(map[string]interface{}{
			"workingExperience": "software engineer",
		})
		profileRepository.Mock.On("UpdateProfile", mock.Anything, 17,
			&models.Profile{
				WorkingExperience: "software engineer",
			}).Return(result, nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/working-experience")
		c.SetParamNames("profileCode")
		c.SetParamValues("17")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		controller := apiHandler.UpdateWorkingExperienceByCode()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedUpdateWorkingExperienceController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		requestBody, _ := json.Marshal(map[string]interface{}{})
		profileRepository.Mock.On("UpdateProfile", mock.Anything, 18, &models.Profile{}).Return(nil, errors.New(""))

		req := httptest.NewRequest(http.MethodPut, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/working_experience")
		c.SetParamNames("profileCode")
		c.SetParamValues("18")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.UpdateWorkingExperienceByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedUpdateWorkingExperienceController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("UpdateProfile", mock.Anything, 0, &models.Profile{
			WorkingExperience: "software engineer",
		}).Return(&models.ProfileDTO{}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.UpdateWorkingExperienceByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedUpdateWorkingExperienceByCodeController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		profileRepository.Mock.On("UpdateProfile", mock.Anything, "asd", &models.Profile{
			WorkingExperience: "software engineer",
		}).Return(&models.ProfileDTO{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.UpdateWorkingExperienceByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestGetEducationListController(t *testing.T) {
	t.Run("SuccessGetEducationListController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		var response = []*models.EducationDTO{}
		result := &models.EducationDTO{
			Id:          1,
			School:      "UGM",
			Degree:      "S1",
			City:        "Jogja",
			Description: "I'm Programmer",
		}
		response = append(response, result)
		educationRepository.Mock.On("GetEducationByProfileCode", mock.Anything, 19).Return(response, nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("19")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		controller := apiHandler.GetEducationListByCode()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedGetEducationListController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		educationRepository.Mock.On("GetEducationByProfileCode", mock.Anything, 20).Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("20")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetEducationListByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedGetEducationListController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		educationRepository.Mock.On("GetEducationByProfileCode", mock.Anything, 0).Return(&models.ProfileDTO{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetEducationListByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedGetEducationByProfileCodeController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		educationRepository.Mock.On("GetEducationByProfileCode", mock.Anything, "asd").Return(&models.ProfileDTO{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetEducationListByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestCreateEducationController(t *testing.T) {
	t.Run("SuccessCreateEducationController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		result := &models.EducationDTO{
			Id:          1,
			School:      "UGM",
			Degree:      "S1",
			City:        "Jogja",
			Description: "I'm Programmer",
		}
		requestBody, _ := json.Marshal(map[string]interface{}{
			"school":      "UGM",
			"degree":      "S1",
			"city":        "Jogja",
			"description": "I'm Programmer",
		})
		educationRepository.Mock.On("CreateEducation", mock.Anything, &models.Education{
			ProfileCode: 21,
			School:      "UGM",
			Degree:      "S1",
			City:        "Jogja",
			Description: "I'm Programmer",
		}).Return(result, nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("21")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		controller := apiHandler.AddEducationByCode()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedCreateEducationController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		educationRepository.Mock.On("CreateEducation", mock.Anything, &models.Education{
			ProfileCode: 22,
		}).Return(nil, errors.New(""))

		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("22")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.AddEducationByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedCreateEducationController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		educationRepository.Mock.On("CreateEducation", mock.Anything, 0, &models.Education{}).Return(&models.EducationDTO{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.AddEducationByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedCreateEducationController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		educationRepository.Mock.On("CreateEducation", mock.Anything, "asd", &models.Education{}).Return(&models.ProfileDTO{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/profile")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.AddEducationByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestDeleteEducationController(t *testing.T) {
	t.Run("SuccessDeleteEducationController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}

		educationRepository.Mock.On("DeleteEducation", mock.Anything, 23, 1).
			Return(nil)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api?id=1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("23")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)
		fmt.Println(c.ParamValues())
		fmt.Println(c.ParamNames())
		// apiHandler.GetProfileByCode()(c)

		controller := apiHandler.DeleteEducationByCodeAndId()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedDeleteEducationController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		educationRepository.Mock.On("DeleteEducation", mock.Anything, 23, 2).
			Return(errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodDelete, "/api?id=2", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("23")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeleteEducationByCodeAndId()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedDeletePhotoController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		educationRepository.Mock.On("DeleteEducation", mock.Anything, 0, 0).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api?id=0", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeleteEducationByCodeAndId()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, "bad request. failed to validate", errMsg)
			assert.Equal(t, http.StatusBadRequest, errCode)
		}

	})

	t.Run("FailedDeleteEducationController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		educationRepository.Mock.On("DeleteEducation", mock.Anything, "asd", 0).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/education")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeleteEducationByCodeAndId()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestGetEmploymentListController(t *testing.T) {
	t.Run("SuccessGetEmploymentListController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		var response = []*models.EmploymentDTO{}
		result := &models.EmploymentDTO{
			Id:          1,
			JobTitle:    "Programmer",
			Employer:    "PT. ABC",
			City:        "Jogja",
			Description: "I'm Programmer",
		}
		response = append(response, result)
		employmentRepository.Mock.On("GetEmploymentByProfileCode", mock.Anything, 1).Return(response, nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("1")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		controller := apiHandler.GetEmploymentListByCode()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedGetEmploymentListController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		employmentRepository.Mock.On("GetEmploymentByProfileCode", mock.Anything, 2).Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("2")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetEmploymentListByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedGetEmploymentListController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		employmentRepository.Mock.On("GetEmploymentByProfileCode", mock.Anything, 0).Return([]*models.EmploymentDTO{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetEmploymentListByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedGetEmploymentListController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		employmentRepository.Mock.On("GetEmploymentByProfileCode", mock.Anything, "asd").Return([]*models.EmploymentDTO{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetEmploymentListByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestCreateEmploymentController(t *testing.T) {
	t.Run("SuccessCreateEmploymentController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		result := &models.EmploymentDTO{
			Id:          1,
			JobTitle:    "Programmer",
			Employer:    "PT. ABC",
			City:        "Jogja",
			Description: "I'm Programmer",
		}
		requestBody, _ := json.Marshal(map[string]interface{}{
			"jobTitle":    "Programmer",
			"employer":    "PT. ABC",
			"city":        "Jogja",
			"description": "I'm Programmer",
		})
		employmentRepository.Mock.On("CreateEmployment", mock.Anything, &models.Employment{
			ProfileCode: 1,
			JobTitle:    "Programmer",
			Employer:    "PT. ABC",
			City:        "Jogja",
			Description: "I'm Programmer",
		}).Return(result, nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("1")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		controller := apiHandler.AddEmploymentByCode()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedCreateEmploymentController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		employmentRepository.Mock.On("CreateEmployment", mock.Anything, &models.Employment{
			ProfileCode: 22,
		}).Return(nil, errors.New(""))

		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("22")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.AddEmploymentByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedCreateEmploymentController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		employmentRepository.Mock.On("CreateEmployment", mock.Anything, &models.Employment{
			ProfileCode: 0,
		}).Return(&models.EducationDTO{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.AddEmploymentByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedCreateEmploymentController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		employmentRepository.Mock.On("CreateEmployment", mock.Anything, "asd").Return(&models.ProfileDTO{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.AddEmploymentByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestDeleteEmploymentController(t *testing.T) {
	t.Run("SuccessDeleteEmploymentController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}

		employmentRepository.Mock.On("DeleteEmployment", mock.Anything, 1, 1).Return(nil)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api?id=1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("1")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)
		// apiHandler.GetProfileByCode()(c)

		controller := apiHandler.DeleteEmploymentByCodeAndId()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedDeleteEmploymentController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		employmentRepository.Mock.On("DeleteEmployment", mock.Anything, 1, 2).Return(errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodDelete, "/api?id=2", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("1")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeleteEmploymentByCodeAndId()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedDeleteEmploymentController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		employmentRepository.Mock.On("DeleteEmployment", mock.Anything, 0, 0).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api?id=0", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeleteEmploymentByCodeAndId()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, "bad request. failed to validate", errMsg)
			assert.Equal(t, http.StatusBadRequest, errCode)
		}

	})

	t.Run("FailedDeleteEmploymentController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		employmentRepository.Mock.On("DeleteEmployment", mock.Anything, "asd", 0).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/employment")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeleteEmploymentByCodeAndId()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestGetSkillListController(t *testing.T) {
	t.Run("SuccessGetSkillListController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		var response = []*models.SkillDTO{}
		result := &models.SkillDTO{
			Id:    1,
			Skill: "Golang",
			Level: "Beginner",
		}
		response = append(response, result)
		skillRepository.Mock.On("GetSkillsByProfileCode", mock.Anything, 1).Return(response, nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("1")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		controller := apiHandler.GetSkillListByCode()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedGetSkillListController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		skillRepository.Mock.On("GetSkillsByProfileCode", mock.Anything, 2).Return(nil, errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("2")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetSkillListByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedGetSkillListController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		skillRepository.Mock.On("GetSkillsByProfileCode", mock.Anything, 0).Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetSkillListByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedGetSkillListController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		skillRepository.Mock.On("GetSkillsByProfileCode", mock.Anything, "asd").Return([]*models.EmploymentDTO{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.GetSkillListByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestCreateSkillController(t *testing.T) {
	t.Run("SuccessCreateSkillController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		result := &models.SkillDTO{
			Id:    1,
			Skill: "Golang",
			Level: "Beginner",
		}
		requestBody, _ := json.Marshal(map[string]interface{}{
			"skill": "Golang",
			"level": "Beginner",
		})
		skillRepository.Mock.On("CreateSkill", mock.Anything, &models.Skill{
			ProfileCode: 1,
			Skill:       "Golang",
			Level:       "Beginner",
		}).Return(result, nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("1")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		controller := apiHandler.AddSkillByCode()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedCreateSkillController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		skillRepository.Mock.On("CreateSkill", mock.Anything, &models.Skill{
			ProfileCode: 2,
		}).Return(nil, errors.New(""))

		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("2")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.AddSkillByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedCreateSkillController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		skillRepository.Mock.On("CreateSkill", mock.Anything, &models.Employment{
			ProfileCode: 0,
		}).Return(nil, nil)

		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.AddSkillByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to validate", errMsg)
		}

	})

	t.Run("FailedCreateSkillController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		skillRepository.Mock.On("CreateSkill", mock.Anything, "asd").Return(nil, nil)

		req := httptest.NewRequest(http.MethodPost, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.AddSkillByCode()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]

			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}

func TestDeleteSkillController(t *testing.T) {
	t.Run("SuccessDeleteSkillController", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}

		skillRepository.Mock.On("DeleteSkill", mock.Anything, 1, 1).Return(nil)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api?id=1", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("1")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)
		// apiHandler.GetProfileByCode()(c)

		controller := apiHandler.DeleteSkillByCodeAndId()(c)
		if assert.NoError(t, controller) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedDeleteSkillController_Err500", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		skillRepository.Mock.On("DeleteSkill", mock.Anything, 1, 2).Return(errors.New("sql: no rows in result set"))

		req := httptest.NewRequest(http.MethodDelete, "/api?id=2", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("1")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeleteSkillByCodeAndId()(c)
		if assert.Error(t, controller) {
			var errCode int
			re := regexp.MustCompile(`code=(\d+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			assert.Equal(t, http.StatusInternalServerError, errCode)
		}

	})

	t.Run("FailedDeleteSkillController_ErrValidate", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		skillRepository.Mock.On("DeleteSkill", mock.Anything, 0, 0).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api?id=0", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("0")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeleteSkillByCodeAndId()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, "bad request. failed to validate", errMsg)
			assert.Equal(t, http.StatusBadRequest, errCode)
		}

	})

	t.Run("FailedDeleteSkillController_ErrBind", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}
		skillRepository.Mock.On("DeleteSkill", mock.Anything, "asd", 0).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/skill")
		c.SetParamNames("profileCode")
		c.SetParamValues("asd")

		apiHandler := NewApiControllerHandler(e.Group("api"), profileServiceTest, skillServiceTest, educationServiceTest, employmentServiceTest)

		// apiHandler.GetProfileByCode()(c)
		controller := apiHandler.DeleteSkillByCodeAndId()(c)
		if assert.Error(t, controller) {
			var errCode int
			var errMsg string

			re := regexp.MustCompile(`code=(\d+), message=(.+)`)
			match := re.FindStringSubmatch(controller.Error())
			errCode, _ = strconv.Atoi(match[1])
			errMsg = match[2]
			assert.Equal(t, http.StatusBadRequest, errCode)
			assert.Equal(t, "bad request. failed to bind", errMsg)
		}

	})
}
