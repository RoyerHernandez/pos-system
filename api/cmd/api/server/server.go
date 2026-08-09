package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RoyerHernandez/pos-system/api/cmd/api/server/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"os/signal"
	"syscall"
)

// Server holds the HTTP server and chi router.
type Server struct {
	router chi.Router
	server *http.Server
	db     *sqlx.DB
}

// New builds the entire application: config, database, dependencies, routes.
func New() *Server {
	// Load .env (non-fatal if missing)
	_ = godotenv.Load()

	// Read env vars
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvIntOrDefault("DB_PORT", 3306)
	dbName := getEnvOrDefault("DB_NAME", "pos")
	dbUser := getEnvOrDefault("DB_USER", "root")
	dbPass := getEnvOrDefault("DB_PASS", "")
	serverPort := getEnvIntOrDefault("SERVER_PORT", 8080)
	jwtSecret := getEnvOrDefault("JWT_SECRET", "pos-secret-change-me")
	corsOrigins := getEnvOrDefault("CORS_ORIGINS", "http://localhost:8000,http://localhost:8080")

	// Validate JWT secret
	if jwtSecret == "" || jwtSecret == "pos-secret-change-me" {
		log.Fatal("JWT_SECRET environment variable must be set to a strong random value")
	}

	// Database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("DB connected")

	// Router with middleware
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(middleware.RequestID)

	rateLimiter := middleware.NewRateLimiter(100, time.Minute)
	r.Use(rateLimiter.Handler)

	origins := strings.Split(corsOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: origins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	// Resolve dependencies: repositories -> services -> controllers
	userRepo := resolveUserRepository(db)
	categoryRepo := resolveCategoryRepository(db)
	productRepo := resolveProductRepository(db)
	clientRepo := resolveClientRepository(db)
	saleRepo := resolveSaleRepository(db)
	inventoryRepo := resolveInventoryRepository(db)
	cashRepo := resolveCashRegisterRepository(db)
	tableRepo := resolveTableRepository(db)
	dashRepo := resolveDashboardRepository(db)
	reportsRepo := resolveReportsRepository(db)

	authSvc := resolveAuthService(userRepo, jwtSecret)
	userSvc := resolveUserService(userRepo)
	categorySvc := resolveCategoryService(categoryRepo)
	productSvc := resolveProductService(productRepo)
	clientSvc := resolveClientService(clientRepo)
	saleSvc := resolveSaleService(db, saleRepo, inventoryRepo)
	inventorySvc := resolveInventoryService(inventoryRepo, db)
	cashSvc := resolveCashRegisterService(cashRepo, db)
	tableSvc := resolveTableService(tableRepo)
	tableOpSvc := resolveTableOperationService(db, tableRepo, saleRepo, inventoryRepo, cashRepo)
	dashSvc := resolveDashboardService(dashRepo)
	reportsSvc := resolveReportsService(reportsRepo)

	healthCtrl := resolveHealthController(db)
	authCtrl := resolveAuthController(authSvc)
	userCtrl := resolveUserController(userSvc)
	categoryCtrl := resolveCategoryController(categorySvc)
	productCtrl := resolveProductController(productSvc)
	clientCtrl := resolveClientController(clientSvc)
	saleCtrl := resolveSaleController(saleSvc)
	inventoryCtrl := resolveInventoryController(inventorySvc)
	cashCtrl := resolveCashRegisterController(cashSvc)
	tableCtrl := resolveTableController(tableSvc, tableOpSvc)
	dashCtrl := resolveDashboardController(dashSvc)
	reportsCtrl := resolveReportsController(reportsSvc)

	// Map all URLs
	mapURLs(r, jwtSecret,
		healthCtrl,
		authCtrl,
		userCtrl,
		categoryCtrl,
		productCtrl,
		clientCtrl,
		saleCtrl,
		inventoryCtrl,
		cashCtrl,
		tableCtrl,
		dashCtrl,
		reportsCtrl,
	)

	// HTTP server with timeouts
	addr := fmt.Sprintf(":%d", serverPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		router: r,
		server: srv,
		db:     db,
	}
}

// Run starts the HTTP server and blocks until a shutdown signal is received.
func (s *Server) Run() {
	defer s.db.Close()

	go func() {
		log.Printf("API server starting on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvIntOrDefault(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}
