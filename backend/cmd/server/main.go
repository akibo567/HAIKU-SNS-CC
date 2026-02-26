package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/goshichigo/backend/internal/config"
	"github.com/goshichigo/backend/internal/db"
	"github.com/goshichigo/backend/internal/handler"
	"github.com/goshichigo/backend/internal/middleware"
	"github.com/goshichigo/backend/internal/repository"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")

	userRepo := repository.NewUserRepository(pool)
	haikuRepo := repository.NewHaikuRepository(pool)

	authHandler := handler.NewAuthHandler(userRepo, cfg.JWTSecret, cfg.JWTRefreshSecret)
	haikuHandler := handler.NewHaikuHandler(haikuRepo)
	userHandler := handler.NewUserHandler(userRepo, haikuRepo)

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(30 * time.Second))
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// 認証
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)
			r.With(middleware.Auth(cfg.JWTSecret)).Post("/logout", authHandler.Logout)
		})

		// 俳句 (タイムライン取得・詳細はOptionalAuth)
		r.Route("/posts", func(r chi.Router) {
			r.With(middleware.OptionalAuth(cfg.JWTSecret)).Get("/", haikuHandler.ListTimeline)
			r.With(middleware.Auth(cfg.JWTSecret)).Post("/", haikuHandler.Create)
			r.With(middleware.OptionalAuth(cfg.JWTSecret)).Get("/{id}", haikuHandler.GetByID)
			r.With(middleware.Auth(cfg.JWTSecret)).Delete("/{id}", haikuHandler.Delete)
			r.With(middleware.Auth(cfg.JWTSecret)).Post("/{id}/like", haikuHandler.Like)
			r.With(middleware.Auth(cfg.JWTSecret)).Delete("/{id}/like", haikuHandler.Unlike)
		})

		// ユーザー
		r.Route("/users", func(r chi.Router) {
			r.With(middleware.Auth(cfg.JWTSecret)).Put("/me", userHandler.UpdateProfile)
			r.With(middleware.OptionalAuth(cfg.JWTSecret)).Get("/{username}", userHandler.GetProfile)
			r.With(middleware.OptionalAuth(cfg.JWTSecret)).Get("/{username}/posts", userHandler.GetPosts)
		})
	})

	addr := ":" + cfg.Port
	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
