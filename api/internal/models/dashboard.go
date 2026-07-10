package models

import "time"

type DashboardKPIs struct {
	TodaySalesCount int            `json:"today_sales_count"`
	TodayRevenue    float64        `json:"today_revenue"`
	TotalProducts   int            `json:"total_products"`
	LowStockCount   int            `json:"low_stock_count"`
	TopProducts     []TopProduct   `json:"top_products"`
	RecentSales     []RecentSale   `json:"recent_sales"`
}

type TopProduct struct {
	ID          int     `json:"id" db:"id"`
	Descripcion string  `json:"descripcion" db:"descripcion"`
	PrecioVenta float64 `json:"precio_venta" db:"precio_venta"`
	Stock       int     `json:"stock" db:"stock"`
	TotalSold   int     `json:"total_sold" db:"total_sold"`
}

type RecentSale struct {
	ID          int     `json:"id" db:"id"`
	CodigoVenta string  `json:"codigo_venta" db:"codigo_venta"`
	Total       float64 `json:"total" db:"total"`
	MetodoPago  string  `json:"metodo_pago" db:"metodo_pago"`
	Estado      string  `json:"estado" db:"estado"`
	Fecha       *time.Time `json:"fecha" db:"fecha"`
}
