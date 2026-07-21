package domain

// Roles
const (
	RoleAdmin    = "Administrador"
	RoleEspecial = "Especial"
	RoleVendedor = "Vendedor"
)

// Sale statuses
const (
	SaleStatusCompleted = "completada"
	SaleStatusCancelled = "cancelada"
)

// Payment methods
const (
	PaymentCash     = "Efectivo"
	PaymentCard     = "Tarjeta"
	PaymentTransfer = "Transferencia"
)

// ValidPaymentMethods lists all accepted payment methods for validation.
var ValidPaymentMethods = []string{PaymentCash, PaymentCard, PaymentTransfer}

// Cash register statuses
const (
	RegisterOpen   = "abierta"
	RegisterClosed = "cerrada"
)

// Inventory movement types
const (
	MovementIn  = "entrada"
	MovementOut = "salida"
)
