package models

import (
	"time"
)

// LogoEmpresa representa un logo almacenado en la base de datos
type LogoEmpresa struct {
	ID                 int       `json:"id" db:"id"`
	IdUsuario          int       `json:"id_usuario" db:"id_usuario"`
	IdEmpresa          *int      `json:"id_empresa" db:"id_empresa"` // Nullable
	NombreLogo         string    `json:"nombre_logo" db:"nombre_logo"`
	ImagenLogo         []byte    `json:"imagen_logo" db:"imagen_logo"`
	TipoMime           string    `json:"tipo_mime" db:"tipo_mime"`
	TamañoArchivo      int       `json:"tamaño_archivo" db:"tamaño_archivo"`
	EstadoLogo         int       `json:"estado_logo" db:"estado_logo"` // 1 = activo, 0 = inactivo
	FechaCreacion      time.Time `json:"fecha_creacion" db:"fecha_creacion"`
	FechaActualizacion time.Time `json:"fecha_actualizacion" db:"fecha_actualizacion"`
}

// LogoRequest representa la estructura para subir un logo
type LogoRequest struct {
	IdUsuario    int    `json:"id_usuario" validate:"required"`
	IdEmpresa    *int   `json:"id_empresa"`
	NombreLogo   string `json:"nombre_logo" validate:"required"`
	ImagenBase64 string `json:"imagen_base64" validate:"required"`
	TipoMime     string `json:"tipo_mime" validate:"required"`
}

// LogoResponse representa la respuesta al consultar un logo
type LogoResponse struct {
	ID            int    `json:"id"`
	IdUsuario     int    `json:"id_usuario"`
	IdEmpresa     *int   `json:"id_empresa"`
	NombreLogo    string `json:"nombre_logo"`
	TipoMime      string `json:"tipo_mime"`
	TamañoArchivo int    `json:"tamaño_archivo"`
	EstadoLogo    int    `json:"estado_logo"`
	FechaCreacion string `json:"fecha_creacion"`
}
