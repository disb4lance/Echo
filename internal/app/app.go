package app

import (
	"context"
	"echo/internal/adapter/minio"
	"echo/internal/adapter/postgres"
	"echo/internal/config"
	handler "echo/internal/handler/http"
	"echo/internal/service"
	transport "echo/internal/transport/http"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Server *http.Server
	DB     *pgxpool.Pool
}

func New(cfg *config.Config, logger *log.Logger) (*App, error) {
	// Изменено: cfg.Database.* вместо cfg.DB.*
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.DBUser,
		cfg.Database.DBPassword,
		cfg.Database.DBHost,
		cfg.Database.DBPort,
		cfg.Database.DBName,
		cfg.Database.DBSSLMode,
	)

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	// Репозитории
	profileRepo := postgres.NewUserProfileRepo(db)
	photoRepo := postgres.NewUserPhotoRepo(db)
	txManaget := postgres.NewTxManager(db)
	minIO, err := minio.NewMinioStorage(cfg.MinIO)

	if err != nil {
		return nil, err
	}

	// Сервис
	photoService := service.NewUserPhotoService(photoRepo, txManaget, minIO)
	profileService := service.NewUserProfileService(profileRepo)

	// Хендлер
	photoHandler := handler.NewUserPhotoHandler(photoService)
	profileHandler := handler.NewUserProfileHandler(profileService)

	// router
	router := transport.NewRouter(photoHandler, profileHandler, cfg)

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%s", cfg.App.APIPort),
		Handler:     router,
		IdleTimeout: time.Minute,
	}

	return &App{
		Server: srv,
		DB:     db,
	}, nil
}
