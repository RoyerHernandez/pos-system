package models

import "time"

type Client struct {
	ID              int        `json:"id" db:"id"`
	Nombre          string     `json:"nombre" db:"nombre"`
	Documento       *string    `json:"documento" db:"documento"`
	Email           *string    `json:"email" db:"email"`
	Telefono        *string    `json:"telefono" db:"telefono"`
	Direccion       *string    `json:"direccion" db:"direccion"`
	FechaNacimiento *string    `json:"fecha_nacimiento" db:"fecha_nacimiento"`
	Compras         float64    `json:"compras" db:"compras"`
	Estado          int        `json:"estado" db:"estado"`
	FechaCreacion   *time.Time `json:"fecha_creacion" db:"fecha_creacion"`
}

type CreateClientRequest struct {
	Nombre          string  `json:"nombre"`
	Documento       *string `json:"documento"`
	Email           *string `json:"email"`
	Telefono        *string `json:"telefono"`
	Direccion       *string `json:"direccion"`
	FechaNacimiento *string `json:"fecha_nacimiento"`
}

type UpdateClientRequest struct {
	Nombre          string  `json:"nombre"`
	Documento       *string `json:"documento"`
	Email           *string `json:"email"`
	Telefono        *string `json:"telefono"`
	Direccion       *string `json:"direccion"`
	FechaNacimiento *string `json:"fecha_nacimiento"`
}
