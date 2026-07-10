package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RoyerHernandez/pos-system/api/internal/config"
	"github.com/RoyerHernandez/pos-system/api/internal/handlers"
	"github.com/RoyerHernandez/pos-system/api/internal/middleware"
	"github.com/RoyerHernandez/pos-system/api/internal/repositories"
	"github.com/RoyerHernandez/pos-system/api/internal/router"
	"github.com/RoyerHernandez/pos-system/api/internal/services"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	db, err := config.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("DB connected")

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	clientRepo := repositories.NewClientRepository(db)
	saleRepo := repositories.NewSaleRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)
	cashRepo := repositories.NewCashRegisterRepository(db)
	dashRepo := repositories.NewDashboardRepository(db)
	reportsRepo := repositories.NewReportsRepository(db)

	// Services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo)
	clientService := services.NewClientService(clientRepo)
	saleService := services.NewSaleService(db, saleRepo, inventoryRepo)
	inventoryService := services.NewInventoryService(db, inventoryRepo)
	cashService := services.NewCashRegisterService(cashRepo)
	dashService := services.NewDashboardService(dashRepo)
	reportsService := services.NewReportsService(reportsRepo)

	// Handlers
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)
	clientHandler := handlers.NewClientHandler(clientService)
	saleHandler := handlers.NewSaleHandler(saleService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	cashHandler := handlers.NewCashRegisterHandler(cashService)
	dashboardHandler := handlers.NewDashboardHandler(dashService)
	reportsHandler := handlers.NewReportsHandler(reportsService)

	// Rate limiter: 100 requests per minute per IP
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	// Router
	r := router.NewRouter(cfg.CORSOrigins)
	r.Use(middleware.RequestID)
	r.Use(rateLimiter.Handler)

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Get("/health", healthHandler.Health)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

			// All authenticated users
			r.Post("/sales", saleHandler.Create)
			r.Post("/cashregister/open", cashHandler.Open)
			r.Put("/cashregister/{id}/close", cashHandler.Close)
			r.Get("/cashregister/current", cashHandler.GetCurrent)

			// Admin + Especial
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("Administrador", "Especial"))

				// Products CRUD
				r.Get("/products", productHandler.GetAll)
				r.Get("/products/{id}", productHandler.GetByID)
				r.Post("/products", productHandler.Create)
				r.Put("/products/{id}", productHandler.Update)
				r.Delete("/products/{id}", productHandler.Delete)

				// Clients CRUD
				r.Get("/clients", clientHandler.GetAll)
				r.Get("/clients/{id}", clientHandler.GetByID)
				r.Post("/clients", clientHandler.Create)
				r.Put("/clients/{id}", clientHandler.Update)
				r.Delete("/clients/{id}", clientHandler.Delete)

				// Sales listing and detail
				r.Get("/sales", saleHandler.GetAll)
				r.Get("/sales/{id}", saleHandler.GetByID)

				// Inventory
				r.Get("/inventory", inventoryHandler.GetAll)
				r.Post("/inventory", inventoryHandler.Create)

				// Cash register list
				r.Get("/cashregister", cashHandler.GetAll)

				// Dashboard + Reports
				r.Get("/dashboard/kpis", dashboardHandler.GetKPIs)
				r.Get("/reports/sales", reportsHandler.GetSalesReport)
				r.Get("/reports/products", reportsHandler.GetProductReport)
				r.Get("/reports/clients", reportsHandler.GetClientReport)
			})

			// Admin only
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("Administrador"))

				// Users CRUD
				r.Get("/users", userHandler.GetAll)
				r.Get("/users/{id}", userHandler.GetByID)
				r.Post("/users", userHandler.Create)
				r.Put("/users/{id}", userHandler.Update)
				r.Delete("/users/{id}", userHandler.Delete)

				// Categories CRUD
				r.Get("/categories", categoryHandler.GetAll)
				r.Get("/categories/{id}", categoryHandler.GetByID)
				r.Post("/categories", categoryHandler.Create)
				r.Put("/categories/{id}", categoryHandler.Update)
				r.Delete("/categories/{id}", categoryHandler.Delete)

				// Cancel sale
				r.Put("/sales/{id}/cancel", saleHandler.Cancel)
			})
		})
	})

	// Server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("API server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
