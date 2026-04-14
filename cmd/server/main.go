package main

import (
	"context"
	"log"
	"net/http"
	"sup-anapa/internal/config"
	"sup-anapa/internal/handlers"
	"sup-anapa/internal/middleware"
	"sup-anapa/internal/repository"
	"sup-anapa/internal/services"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	// Подключение к БД
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer pool.Close()

	// Проверка подключения
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Unable to ping database:", err)
	}
	log.Println("✓ Successfully connected to database")

	// Инициализация репозиториев
	bookingRepo := repository.NewBookingRepository(pool)
	instructorRepo := repository.NewInstructorRepository(pool)
	slotRepo := repository.NewSlotRepository(pool)
	adminRepo := repository.NewAdminRepository(pool)

	// Инициализация сервисов
	notificationService := services.NewNotificationService(cfg.VKBotToken)
	bookingService := services.NewBookingService(bookingRepo, notificationService)
	weatherService := services.NewWeatherService(cfg.WeatherAPIKey)
	authService := services.NewAuthService(adminRepo)

	// Инициализация handlers
	handlers.Init(cfg.SessionSecret, authService)
	handlers.SetRepositories(bookingRepo, instructorRepo, slotRepo)

	// Инициализация middleware
	middleware.InitAuth(handlers.GetStore())

	// Инициализация API handlers
	instructorHandler := handlers.NewInstructorHandler(instructorRepo)
	slotHandler := handlers.NewSlotHandler(slotRepo, instructorRepo)
	bookingHandler := handlers.NewBookingHandler(bookingRepo, slotRepo, instructorRepo)

	// TODO: Передать сервисы в handlers
	_ = bookingService
	_ = weatherService

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Static files
	fileServer := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Public routes
	r.Get("/", handlers.Home)
	r.Get("/booking", handlers.BookingPage)
	r.Get("/instructors", handlers.InstructorsPage)
	r.Post("/booking", handlers.CreateBooking)

	// Public API (no auth required)
	r.Get("/api/instructors", instructorHandler.List)
	r.Get("/api/slots", slotHandler.List)

	// Admin routes
	r.Route("/admin", func(r chi.Router) {
		r.Get("/login", handlers.AdminLogin)
		r.Post("/login", handlers.AdminLoginPost)
		r.Get("/logout", handlers.AdminLogout)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Get("/", handlers.AdminDashboard)
			r.Get("/instructors", handlers.AdminInstructors)
			r.Get("/slots", handlers.AdminSlots)
			r.Get("/bookings", handlers.AdminBookings)

			// API routes
			r.Route("/api/instructors", func(r chi.Router) {
				r.Get("/", instructorHandler.List)
				r.Post("/", instructorHandler.Create)
				r.Get("/{id}", instructorHandler.Get)
				r.Put("/{id}", instructorHandler.Update)
				r.Delete("/{id}", instructorHandler.Delete)
			})

			r.Route("/api/slots", func(r chi.Router) {
				r.Get("/", slotHandler.List)
				r.Post("/", slotHandler.Create)
				r.Get("/{id}", slotHandler.Get)
				r.Put("/{id}", slotHandler.Update)
				r.Delete("/{id}", slotHandler.Delete)
			})

			r.Route("/api/bookings", func(r chi.Router) {
				r.Get("/", bookingHandler.List)
				r.Get("/{id}", bookingHandler.Get)
				r.Put("/{id}/status", bookingHandler.UpdateStatus)
			})
		})
	})

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			ctx := context.Background()
			n, err := slotRepo.ExpireHolds(ctx)
			if err != nil {
				log.Printf("Error expiring slot holds: %v", err)
			} else if n > 0 {
				log.Printf("Expired %d pending slot holds", n)
			}
		}
	}()

	addr := ":" + cfg.Port
	log.Printf("✓ Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
