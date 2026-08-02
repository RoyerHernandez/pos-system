package server

import (
	"github.com/RoyerHernandez/pos-system/api/cmd/api/server/middleware"
	"github.com/RoyerHernandez/pos-system/api/internal/controller/web"
	"github.com/go-chi/chi/v5"
)

func mapURLs(
	r chi.Router,
	jwtSecret string,
	health *web.HealthController,
	auth *web.AuthController,
	user *web.UserController,
	category *web.CategoryController,
	product *web.ProductController,
	client *web.ClientController,
	sale *web.SaleController,
	inventory *web.InventoryController,
	cash *web.CashRegisterController,
	dashboard *web.DashboardController,
	reports *web.ReportsController,
	table *web.TableController,
) {
	r.Route("/api/pos/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/login", auth.Login)
		r.Post("/auth/refresh", auth.Refresh)
		r.Get("/health", health.Health)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtSecret))

			// All authenticated users
			r.Post("/sales", sale.Create)
			r.Post("/cashregister/open", cash.Open)
			r.Put("/cashregister/{id}/close", cash.Close)
			r.Get("/cashregister/current", cash.GetCurrent)

			// Admin + Especial
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("Administrador", "Especial"))

				// Products CRUD
				r.Get("/products", product.GetAll)
				r.Get("/products/{id}", product.GetByID)
				r.Post("/products", product.Create)
				r.Put("/products/{id}", product.Update)
				r.Delete("/products/{id}", product.Delete)

				// Clients CRUD
				r.Get("/clients", client.GetAll)
				r.Get("/clients/{id}", client.GetByID)
				r.Post("/clients", client.Create)
				r.Put("/clients/{id}", client.Update)
				r.Delete("/clients/{id}", client.Delete)

				// Sales listing and detail
				r.Get("/sales", sale.GetAll)
				r.Get("/sales/{id}", sale.GetByID)

				// Inventory
				r.Get("/inventory", inventory.GetAll)
				r.Post("/inventory", inventory.Create)

				// Cash register list
				r.Get("/cashregister", cash.GetAll)

				// Dashboard + Reports
				r.Get("/dashboard/kpis", dashboard.GetKPIs)
				r.Get("/reports/sales", reports.GetSalesReport)
				r.Get("/reports/products", reports.GetProductReport)
				r.Get("/reports/clients", reports.GetClientReport)
			})

			// Admin only
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("Administrador"))

				// Users CRUD
				r.Get("/users", user.GetAll)
				r.Get("/users/{id}", user.GetByID)
				r.Post("/users", user.Create)
				r.Put("/users/{id}", user.Update)
				r.Delete("/users/{id}", user.Delete)

				// Categories CRUD
				r.Get("/categories", category.GetAll)
				r.Get("/categories/{id}", category.GetByID)
				r.Post("/categories", category.Create)
				r.Put("/categories/{id}", category.Update)
				r.Delete("/categories/{id}", category.Delete)

				// Tables CRUD
				r.Get("/tables", table.GetAll)
				r.Get("/tables/{id}", table.GetByID)
				r.Post("/tables", table.Create)
				r.Put("/tables/{id}", table.Update)
				r.Delete("/tables/{id}", table.Delete)

				// Cancel sale
				r.Put("/sales/{id}/cancel", sale.Cancel)
			})
		})
	})
}
