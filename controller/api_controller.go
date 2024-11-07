package controller

import (
	"net/http"
	"test-bpjs/v2/models/request"
	educationService "test-bpjs/v2/service/education"
	employmentService "test-bpjs/v2/service/employment"
	profileService "test-bpjs/v2/service/profile"
	skillService "test-bpjs/v2/service/skill"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var apiTracer = otel.Tracer("apiController")

type apiControllerHandler struct {
	group             *echo.Group
	profileService    profileService.ProfileService
	skillService      skillService.SkillService
	educationService  educationService.EducationService
	employmentService employmentService.EmploymentService
}

func NewApiControllerHandler(
	group *echo.Group,
	profileService profileService.ProfileService,
	skillService skillService.SkillService,
	educationService educationService.EducationService,
	employmentService employmentService.EmploymentService,
) *apiControllerHandler {
	return &apiControllerHandler{
		group:             group,
		profileService:    profileService,
		skillService:      skillService,
		educationService:  educationService,
		employmentService: employmentService,
	}
}

func (h *apiControllerHandler) MapRoutes() {
	//profile
	h.group.GET("/profile/:profileCode", h.GetProfileByCode())
	h.group.POST("/profile", h.CreateProfile())
	h.group.PUT("/profile/:profileCode", h.UpdateProfile())

	//photo
	h.group.GET("/photo/:profileCode", h.DownloadPhoto())
	h.group.PUT("/photo/:profileCode", h.UploadPhoto())
	h.group.DELETE("/photo/:profileCode", h.DeletePhoto())

	//working experiences
	h.group.GET("/working-experience/:profileCode", h.GetWorkingExperienceByCode())
	h.group.PUT("/working-experience/:profileCode", h.UpdateWorkingExperienceByCode())

	//education
	h.group.GET("/education/:profileCode", h.GetEducationListByCode())
	h.group.POST("/education/:profileCode", h.AddEducationByCode())
	h.group.DELETE("/education/:profileCode", h.DeleteEducationByCodeAndId())

	//employment
	h.group.GET("/employment/:profileCode", h.GetEmploymentListByCode())
	h.group.POST("/employment/:profileCode", h.AddEmploymentByCode())
	h.group.DELETE("/employment/:profileCode", h.DeleteEmploymentByCodeAndId())

	//skill
	h.group.GET("/skill/:profileCode", h.GetSkillListByCode())
	h.group.POST("/skill/:profileCode", h.AddSkillByCode())
	h.group.DELETE("/skill/:profileCode", h.DeleteSkillByCodeAndId())
}

func (h *apiControllerHandler) GetProfileByCode() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "GetProfileByCode", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.GetProfileRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.profileService.GetProfileByCode(ctx, request.ProfileCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) CreateProfile() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "CreateProfile", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.CreateProfileRequest
		// Read the raw JSON data from the request body
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		res, err := h.profileService.CreateProfile(ctx, request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) UpdateProfile() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "UpdateProfile", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.UpdateProfileRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.profileService.UpdateProfile(ctx, request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) DownloadPhoto() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "DownloadPhoto", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.GetProfileRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.profileService.DownloadPhotoByCode(ctx, request.ProfileCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) UploadPhoto() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "DownloadPhoto", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.UploadPhotoRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.profileService.UploadPhotoByCode(ctx, request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) DeletePhoto() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "DeletePhoto", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.GetProfileRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.profileService.DeletePhotoByCode(ctx, request.ProfileCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) GetWorkingExperienceByCode() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "GetWorkingExperienceByCode", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.GetProfileRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.profileService.GetWorkingExperienceByCode(ctx, request.ProfileCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) UpdateWorkingExperienceByCode() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "UpdateWorkingExperienceByCode", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.UpdateProfileRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.profileService.UpdateProfile(ctx, request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) GetEducationListByCode() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "GetEducationListByCode", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.GetProfileRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.educationService.GetEducationByCode(ctx, request.ProfileCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) AddEducationByCode() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "AddEducationByCode", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.CreateEducationRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.educationService.CreateEducation(ctx, request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) DeleteEducationByCodeAndId() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "DeleteEducationByCodeAndId", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.DefaultGetDataByProfileCodeAndId
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.educationService.DeleteEducation(ctx, request.ProfileCode, request.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) GetEmploymentListByCode() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "GetEmploymentListByCode", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.GetProfileRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.employmentService.GetEmploymentByCode(ctx, request.ProfileCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) AddEmploymentByCode() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "AddEmploymentByCode", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.CreateEmploymentRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.employmentService.CreateEmployment(ctx, request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) DeleteEmploymentByCodeAndId() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "DeleteEmploymentByCodeAndId", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.DefaultGetDataByProfileCodeAndId
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.employmentService.DeleteEmployment(ctx, request.ProfileCode, request.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) GetSkillListByCode() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "GetSkillListByCode", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.GetProfileRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.skillService.GetSkillsByCode(ctx, request.ProfileCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) AddSkillByCode() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "AddSkillByCode", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.CreateSkillRequest
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.skillService.CreateSkill(ctx, request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *apiControllerHandler) DeleteSkillByCodeAndId() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, span := apiTracer.Start(c.Request().Context(), "DeleteSkillByCodeAndId", trace.WithTimestamp(time.Now()), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End(trace.WithStackTrace(true))

		var request request.DefaultGetDataByProfileCodeAndId
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to bind")
		}

		if err := c.Validate(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request. failed to validate")
		}

		res, err := h.skillService.DeleteSkill(ctx, request.ProfileCode, request.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, res)
	}
}
