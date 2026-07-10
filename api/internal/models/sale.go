package models

import "time"

type Sale struct {
	ID         int        `json:"id" db:"id"`
	IDUsuario  int        `json:"id_usuario" db:"id_usuario"`
	IDCliente  *int       `json:"id_cliente" db:"id_cliente"`
	CodigoVenta string   `json:"codigo_venta" db:"codigo_venta"`
	Subtotal   float64    `json:"subtotal" db:"subtotal"`
	Impuesto   float64    `json:"impuesto" db:"impuesto"`
	Descuento  float64    `json:"descuento" db:"descuento"`
	Total      float64    `json:"total" db:"total"`
	MetodoPago string     `json:"metodo_pago" db:"metodo_pago"`
	Estado     string     `json:"estado" db:"estado"`
	Fecha      *time.Time `json:"fecha" db:"fecha"`
	IDCaja     *int       `json:"id_caja" db:"id_caja"`
}

type SaleDetail struct {
	ID             int     `json:"id" db:"id"`
	IDVenta        int     `json:"id_venta" db:"id_venta"`
	IDProducto     int     `json:"id_producto" db:"id_producto"`
	Cantidad       int     `json:"cantidad" db:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario" db:"precio_unitario"`
	Descuento      float64 `json:"descuento" db:"descuento"`
	Subtotal       float64 `json:"subtotal" db:"subtotal"`
}

type SaleItem struct {
	IDProducto int `json:"id_producto"`
	Cantidad   int `json:"cantidad"`
}

type CreateSaleRequest struct {
	IDCliente  int        `json:"id_cliente"`
	MetodoPago string     `json:"metodo_pago"`
	Productos  []SaleItem `json:"productos"`
}

type SaleResponse struct {
	Sale
	Details []SaleDetail `json:"details"`
}
