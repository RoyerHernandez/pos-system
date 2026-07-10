package models

import "time"

type Category struct {
	ID            int        `json:"id" db:"id"`
	Nombre        string     `json:"nombre" db:"nombre"`
	Descripcion   *string    `json:"descripcion" db:"descripcion"`
	Estado        int        `json:"estado" db:"estado"`
	FechaCreacion *time.Time `json:"fecha_creacion" db:"fecha_creacion"`
}

type CreateCategoryRequest struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

type UpdateCategoryRequest struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}
