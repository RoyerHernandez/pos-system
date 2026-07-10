package models

import "time"

type CashRegister struct {
	ID                 int        `json:"id" db:"id"`
	IDUsuario          int        `json:"id_usuario" db:"id_usuario"`
	MontoApertura      float64    `json:"monto_apertura" db:"monto_apertura"`
	MontoCierre        *float64   `json:"monto_cierre" db:"monto_cierre"`
	TotalVentas        float64    `json:"total_ventas" db:"total_ventas"`
	TotalEfectivo      float64    `json:"total_efectivo" db:"total_efectivo"`
	TotalTarjeta       float64    `json:"total_tarjeta" db:"total_tarjeta"`
	TotalTransferencia float64    `json:"total_transferencia" db:"total_transferencia"`
	Estado             string     `json:"estado" db:"estado"`
	FechaApertura      *time.Time `json:"fecha_apertura" db:"fecha_apertura"`
	FechaCierre        *time.Time `json:"fecha_cierre" db:"fecha_cierre"`
}

type OpenCashRegisterRequest struct {
	MontoApertura float64 `json:"monto_apertura"`
}

type CloseCashRegisterRequest struct {
	MontoCierre float64 `json:"monto_cierre"`
}
