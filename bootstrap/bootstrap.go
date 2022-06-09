package bootstrap

import (
	"context"
	"net/http"
	"time"

	"github.com/Aguztinus/petty-cash-backend/api/controllers"
	"github.com/Aguztinus/petty-cash-backend/api/middlewares"
	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/api/routes"
	"github.com/Aguztinus/petty-cash-backend/api/services"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"

	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	controllers.Module,
	routes.Module,
	lib.Module,
	services.Module,
	middlewares.Module,
	repository.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	handler lib.HttpHandler,
	routes routes.Routes,
	logger lib.Logger,
	config lib.Config,
	middlewares middlewares.Middlewares,
	database lib.Database,
) {
	db, err := database.ORM.DB()
	if err != nil {
		logger.Zap.Fatalf("Error to get database connection: %v", err)
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Zap.Info("Starting Application")

			if err := db.Ping(); err != nil {
				logger.Zap.Fatalf("Error to ping database connection: %v", err)
			}

			// set conn
			db.SetMaxOpenConns(config.Database.MaxOpenConns)
			db.SetMaxIdleConns(config.Database.MaxIdleConns)
			db.SetConnMaxLifetime(time.Duration(config.Database.MaxLifetime) * time.Second)

			go func() {
				middlewares.Setup()
				routes.Setup()

				if err := handler.Engine.Start(config.Http.ListenAddr()); err != nil {
					if errors.Is(err, http.ErrServerClosed) {
						logger.Zap.Debug("Shutting down the Application")
					} else {
						logger.Zap.Fatalf("Error to Start Application: %v", err)
					}
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			logger.Zap.Info("Stopping Application")

			handler.Engine.Close()
			db.Close()
			return nil
		},
	})
}
