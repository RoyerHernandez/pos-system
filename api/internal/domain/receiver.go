package domain

import "time"

// UserResponse is a cross-layer DTO that excludes the password field.
// No json tags here; controller-layer DTOs add serialization concerns.
type UserResponse struct {
	ID          int
	Usuario     string
	Nombre      string
	Perfil      string
	Foto        *string
	UltimoLogin *time.Time
	Estado      int
	FechaCreacion *time.Time
}

// ToResponse converts a User to a UserResponse, stripping the password.
func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Usuario:       u.Usuario,
		Nombre:        u.Nombre,
		Perfil:        u.Perfil,
		Foto:          u.Foto,
		UltimoLogin:   u.UltimoLogin,
		Estado:        u.Estado,
		FechaCreacion: u.FechaCreacion,
	}
}
