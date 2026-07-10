package models

import "time"

type InventoryMovement struct {
	ID            int       `json:"id" db:"id"`
	IDProducto    int       `json:"id_producto" db:"id_producto"`
	IDUsuario     int       `json:"id_usuario" db:"id_usuario"`
	Tipo          string    `json:"tipo" db:"tipo"`
	Motivo        string    `json:"motivo" db:"motivo"`
	Cantidad      int       `json:"cantidad" db:"cantidad"`
	Observaciones *string   `json:"observaciones" db:"observaciones"`
	IDReferencia  *int      `json:"id_referencia" db:"id_referencia"`
	Fecha         time.Time `json:"fecha" db:"fecha"`
}

type CreateMovementRequest struct {
	IDProducto    int     `json:"id_producto"`
	Tipo          string  `json:"tipo"`
	Motivo        string  `json:"motivo"`
	Cantidad      int     `json:"cantidad"`
	Observaciones *string `json:"observaciones"`
}
