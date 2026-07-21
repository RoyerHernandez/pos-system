package server

import (
	"log"

	"github.com/RoyerHernandez/pos-system/api/internal/controller/web"
	"github.com/RoyerHernandez/pos-system/api/internal/repository"
	"github.com/RoyerHernandez/pos-system/api/internal/service"
	"github.com/jmoiron/sqlx"
)

// --- Repositories ---

func resolveUserRepository(db *sqlx.DB) *repository.UserRepository {
	repo, err := repository.NewUserRepository(db)
	panicOnError(err)
	return repo
}

func resolveCategoryRepository(db *sqlx.DB) *repository.CategoryRepository {
	repo, err := repository.NewCategoryRepository(db)
	panicOnError(err)
	return repo
}

func resolveProductRepository(db *sqlx.DB) *repository.ProductRepository {
	repo, err := repository.NewProductRepository(db)
	panicOnError(err)
	return repo
}

func resolveClientRepository(db *sqlx.DB) *repository.ClientRepository {
	repo, err := repository.NewClientRepository(db)
	panicOnError(err)
	return repo
}

func resolveSaleRepository(db *sqlx.DB) *repository.SaleRepository {
	repo, err := repository.NewSaleRepository(db)
	panicOnError(err)
	return repo
}

func resolveInventoryRepository(db *sqlx.DB) *repository.InventoryRepository {
	repo, err := repository.NewInventoryRepository(db)
	panicOnError(err)
	return repo
}

func resolveCashRegisterRepository(db *sqlx.DB) *repository.CashRegisterRepository {
	repo, err := repository.NewCashRegisterRepository(db)
	panicOnError(err)
	return repo
}

func resolveDashboardRepository(db *sqlx.DB) *repository.DashboardRepository {
	repo, err := repository.NewDashboardRepository(db)
	panicOnError(err)
	return repo
}

func resolveReportsRepository(db *sqlx.DB) *repository.ReportsRepository {
	repo, err := repository.NewReportsRepository(db)
	panicOnError(err)
	return repo
}

// --- Services ---

func resolveAuthService(userRepo *repository.UserRepository, jwtSecret string) *service.AuthService {
	svc, err := service.NewAuthService(userRepo, jwtSecret)
	panicOnError(err)
	return svc
}

func resolveUserService(repo *repository.UserRepository) *service.UserService {
	svc, err := service.NewUserService(repo)
	panicOnError(err)
	return svc
}

func resolveCategoryService(repo *repository.CategoryRepository) *service.CategoryService {
	svc, err := service.NewCategoryService(repo)
	panicOnError(err)
	return svc
}

func resolveProductService(repo *repository.ProductRepository) *service.ProductService {
	svc, err := service.NewProductService(repo)
	panicOnError(err)
	return svc
}

func resolveClientService(repo *repository.ClientRepository) *service.ClientService {
	svc, err := service.NewClientService(repo)
	panicOnError(err)
	return svc
}

func resolveSaleService(db *sqlx.DB, saleRepo *repository.SaleRepository, inventoryRepo *repository.InventoryRepository) *service.SaleService {
	svc, err := service.NewSaleService(db, saleRepo, inventoryRepo)
	panicOnError(err)
	return svc
}

func resolveInventoryService(repo *repository.InventoryRepository, db *sqlx.DB) *service.InventoryService {
	svc, err := service.NewInventoryService(repo, db)
	panicOnError(err)
	return svc
}

func resolveCashRegisterService(repo *repository.CashRegisterRepository, db *sqlx.DB) *service.CashRegisterService {
	svc, err := service.NewCashRegisterService(repo, db)
	panicOnError(err)
	return svc
}

func resolveDashboardService(repo *repository.DashboardRepository) *service.DashboardService {
	svc, err := service.NewDashboardService(repo)
	panicOnError(err)
	return svc
}

func resolveReportsService(repo *repository.ReportsRepository) *service.ReportsService {
	svc, err := service.NewReportsService(repo)
	panicOnError(err)
	return svc
}

// --- Controllers ---

func resolveHealthController(db *sqlx.DB) *web.HealthController {
	ctrl, err := web.NewHealthController(db)
	panicOnError(err)
	return ctrl
}

func resolveAuthController(svc *service.AuthService) *web.AuthController {
	ctrl, err := web.NewAuthController(svc)
	panicOnError(err)
	return ctrl
}

func resolveUserController(svc *service.UserService) *web.UserController {
	ctrl, err := web.NewUserController(svc)
	panicOnError(err)
	return ctrl
}

func resolveCategoryController(svc *service.CategoryService) *web.CategoryController {
	ctrl, err := web.NewCategoryController(svc)
	panicOnError(err)
	return ctrl
}

func resolveProductController(svc *service.ProductService) *web.ProductController {
	ctrl, err := web.NewProductController(svc)
	panicOnError(err)
	return ctrl
}

func resolveClientController(svc *service.ClientService) *web.ClientController {
	ctrl, err := web.NewClientController(svc)
	panicOnError(err)
	return ctrl
}

func resolveSaleController(svc *service.SaleService) *web.SaleController {
	ctrl, err := web.NewSaleController(svc)
	panicOnError(err)
	return ctrl
}

func resolveInventoryController(svc *service.InventoryService) *web.InventoryController {
	ctrl, err := web.NewInventoryController(svc)
	panicOnError(err)
	return ctrl
}

func resolveCashRegisterController(svc *service.CashRegisterService) *web.CashRegisterController {
	ctrl, err := web.NewCashRegisterController(svc)
	panicOnError(err)
	return ctrl
}

func resolveDashboardController(svc *service.DashboardService) *web.DashboardController {
	ctrl, err := web.NewDashboardController(svc)
	panicOnError(err)
	return ctrl
}

func resolveReportsController(svc *service.ReportsService) *web.ReportsController {
	ctrl, err := web.NewReportsController(svc)
	panicOnError(err)
	return ctrl
}

func panicOnError(err error) {
	if err != nil {
		log.Fatalf("dependency resolution failed: %v", err)
	}
}
