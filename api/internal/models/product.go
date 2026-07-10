package models

import "time"

type Product struct {
	ID           int       `json:"id" db:"id"`
	IDCategoria  int       `json:"id_categoria" db:"id_categoria"`
	Descripcion  string    `json:"descripcion" db:"descripcion"`
	Imagen       *string   `json:"imagen" db:"imagen"`
	Stock        int       `json:"stock" db:"stock"`
	PrecioCompra float64   `json:"precio_compra" db:"precio_compra"`
	PrecioVenta  float64   `json:"precio_venta" db:"precio_venta"`
	Ventas       int       `json:"ventas" db:"ventas"`
	Estado       int       `json:"estado" db:"estado"`
	Fecha        time.Time `json:"fecha" db:"fecha"`
}

type CreateProductRequest struct {
	IDCategoria  int     `json:"id_categoria"`
	Descripcion  string  `json:"descripcion"`
	PrecioCompra float64 `json:"precio_compra"`
	PrecioVenta  float64 `json:"precio_venta"`
	Stock        int     `json:"stock"`
}

type UpdateProductRequest struct {
	IDCategoria  int     `json:"id_categoria"`
	Descripcion  string  `json:"descripcion"`
	PrecioCompra float64 `json:"precio_compra"`
	PrecioVenta  float64 `json:"precio_venta"`
	Stock        int     `json:"stock"`
}
