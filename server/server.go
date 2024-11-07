package server

import (
	"context"
	"test-bpjs/v2/config"
	"test-bpjs/v2/controller"
	educationService "test-bpjs/v2/service/education"
	employmentService "test-bpjs/v2/service/employment"
	profileService "test-bpjs/v2/service/profile"
	skillService "test-bpjs/v2/service/skill"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

const TIMEOUT_MESSAGE = "request timeout"
const appVersion = "1.0"

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
func RunServer(ctx context.Context,
	cfg *config.Config,
	db bun.IDB,
	profileService profileService.ProfileService,
	skillService skillService.SkillService,
	employmentService employmentService.EmploymentService,
	educationService educationService.EducationService,
) {
	e := echo.New()
	defer e.Close()

	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogLatency:  true,
		LogProtocol: true,
		LogRemoteIP: true,
		LogHost:     true,
		LogMethod:   true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			latency := time.Duration(values.Latency) * time.Nanosecond
			log.WithContext(c.Request().Context()).Infof(
				"[RESPONSE]: URI=%s, Status=%d, Latency=%v, Protocol=%s, RemoteIP=%s, Host=%s, Method=%s",
				values.URI, values.Status, latency, values.Protocol, values.RemoteIP, values.Host, values.Method,
			)

			return nil
		},
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	go func() {
		log.Fatal(e.Start("localhost:8080"))
	}()

	// logger
	e.Pre(middleware.RemoveTrailingSlash(), middleware.Logger())

	apiController := controller.NewApiControllerHandler(e.Group("/api"), profileService, skillService, educationService, employmentService)
	apiController.MapRoutes()

	<-ctx.Done()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatalf("Error when shuting down: %v", err)
	}
}
