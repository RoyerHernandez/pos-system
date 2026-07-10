package models

import "time"

type User struct {
	ID            int        `json:"id" db:"id"`
	Usuario       string     `json:"usuario" db:"usuario"`
	Password      string     `json:"-" db:"password"`
	Nombre        string     `json:"nombre" db:"nombre"`
	Perfil        string     `json:"perfil" db:"perfil"`
	Foto          *string    `json:"foto" db:"foto"`
	UltimoLogin   *time.Time `json:"ultimo_login" db:"ultimo_login"`
	Estado        int        `json:"estado" db:"estado"`
	FechaCreacion *time.Time `json:"fecha_creacion" db:"fecha_creacion"`
}

type UserLogin struct {
	Usuario  string `json:"usuario"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID          int        `json:"id"`
	Nombre      string     `json:"nombre"`
	Usuario     string     `json:"usuario"`
	Perfil      string     `json:"perfil"`
	Foto        *string    `json:"foto"`
	Estado      int        `json:"estado"`
	UltimoLogin *time.Time `json:"ultimo_login"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Nombre:      u.Nombre,
		Usuario:     u.Usuario,
		Perfil:      u.Perfil,
		Foto:        u.Foto,
		Estado:      u.Estado,
		UltimoLogin: u.UltimoLogin,
	}
}
