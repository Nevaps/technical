package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"os"
	"os/signal"
	"syscall"
	"techical/internal/config"
	"techical/internal/handlers"
	repository "techical/internal/repository/currency"
	"techical/internal/service"

	_ "github.com/lib/pq"
)

type App struct {
	log    *zerolog.Logger
	config *config.Config
}

func NewApp(log *zerolog.Logger, config *config.Config) *App {
	return &App{
		log:    log,
		config: config,
	}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", a.config.DatabaseUser, a.config.DatabasePass, a.config.DatabaseHost, a.config.DatabasePort, a.config.DatabaseName)

	dsn := fmt.Sprintf("%s", connection)
	dbConn, err := sql.Open(`postgres`, dsn)
	if err != nil {
		a.log.Fatal().Err(err).Msg("failed to open connection to database")
	}
	err = dbConn.Ping()
	if err != nil {
		a.log.Fatal().Err(err).Msg("failed to ping database")
	}

	currencyRepo := repository.NewCurrencyRepository(dbConn)

	svc := service.NewCurrencyService(a.config, a.log, currencyRepo)
	svc.Run(ctx)

	app := fiber.New()

	handler := handlers.NewHandler(svc)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Get("/rates", handler.GetRate)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		if err := app.Shutdown(); err != nil {
			a.log.Error().Err(err).Msg("server shutdown failed")
		}
	}()

	a.log.Info().Msgf("starting server on port %s", a.config.ListenPort)
	if err := app.Listen(a.config.ListenPort); err != nil {
		a.log.Fatal().Err(err).Msg("failed to start server")
	}
}
