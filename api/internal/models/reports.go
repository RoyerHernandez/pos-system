package models

type SalesReport struct {
	TotalSales    int     `json:"total_sales" db:"total_sales"`
	TotalRevenue  float64 `json:"total_revenue" db:"total_revenue"`
	TotalEfectivo float64 `json:"total_efectivo" db:"total_efectivo"`
	TotalTarjeta  float64 `json:"total_tarjeta" db:"total_tarjeta"`
	TotalTransfer float64 `json:"total_transferencia" db:"total_transferencia"`
	AvgTicket     float64 `json:"avg_ticket" db:"avg_ticket"`
}

type ProductReport struct {
	ID          int     `json:"id" db:"id"`
	Descripcion string  `json:"descripcion" db:"descripcion"`
	TotalSold   int     `json:"total_sold" db:"total_sold"`
	TotalRevenue float64 `json:"total_revenue" db:"total_revenue"`
}

type ClientReport struct {
	ID           int     `json:"id" db:"id"`
	Nombre       string  `json:"nombre" db:"nombre"`
	Documento    *string `json:"documento" db:"documento"`
	TotalCompras float64 `json:"total_compras" db:"total_compras"`
	NumCompras   int     `json:"num_compras" db:"num_compras"`
}
