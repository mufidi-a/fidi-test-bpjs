package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"test-bpjs/v2/config"
	"test-bpjs/v2/repository"
	"test-bpjs/v2/server"
	educationService "test-bpjs/v2/service/education"
	employmentService "test-bpjs/v2/service/employment"
	profileService "test-bpjs/v2/service/profile"
	skillService "test-bpjs/v2/service/skill"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunotel"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	cfg, err := config.LoadConfig("./config", "config")
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
		cancel()
	}

	dbConn := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(cfg.DatabaseURL),
		pgdriver.WithConnParams(map[string]interface{}{
			"search_path": cfg.DatabaseSchema,
		})))
	defer dbConn.Close()
	if err := dbConn.PingContext(ctx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	dbConn.SetMaxOpenConns(5 * runtime.GOMAXPROCS(0))
	dbConn.SetMaxIdleConns(5 * runtime.GOMAXPROCS(0))

	bunDB := bun.NewDB(dbConn, pgdialect.New(), bun.WithDiscardUnknownColumns())
	bunDB.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName("test-bpjs")))

	profileRepository := repository.NewProfileRepository(bunDB)
	skillRepository := repository.NewSkillRepository(bunDB)
	employmentRepository := repository.NewEmploymentRepository(bunDB)
	educationRepository := repository.NewEducationRepository(bunDB)

	profileService := profileService.NewProfileService(profileRepository)
	skillService := skillService.NewSkillService(skillRepository)
	employmentService := employmentService.NewEmploymentService(employmentRepository)
	educationService := educationService.NewEducationService(educationRepository)

	server.RunServer(ctx,
		&cfg,
		bunDB,
		profileService,
		skillService,
		employmentService,
		educationService,
	)
}
