package http

import (
	"echo/internal/config"
	handler "echo/internal/handler/http"
	mid "echo/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	photoHandler *handler.UserPhotoHandler,
	profileHandler *handler.UserProfileHandler,
	cfg *config.Config,
) *chi.Mux {

	authMiddleware := mid.NewAuthMiddleware(cfg.JWT.JWTSecret)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		r.Get("/user/photos", photoHandler.GetUserPhotos)
		r.Put("/user/photos", photoHandler.ReplaceUserPhotos)

		r.Get("/profiles/{id}", profileHandler.Get)
		r.Post("/profiles", profileHandler.Create)
		r.Put("/profiles/{id}", profileHandler.Update)
	})

	return r
}
