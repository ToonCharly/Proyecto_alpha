package services

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"Facts/internal/db"
	"Facts/internal/models"
)

const errConectarDB = "error al conectar a la base de datos: %w"

// LogoService maneja las operaciones relacionadas con logos
type LogoService struct{}

// NewLogoService crea una nueva instancia del servicio de logos
func NewLogoService() *LogoService {
	return &LogoService{}
}

// GuardarLogo guarda un logo en la base de datos
func (s *LogoService) GuardarLogo(logoReq models.LogoRequest) (*models.LogoEmpresa, error) {
	// Conectar a la base de datos
	dbConn, err := db.ConnectUserDB()
	if err != nil {
		return nil, fmt.Errorf(errConectarDB, err)
	}
	defer dbConn.Close()

	// Decodificar la imagen base64
	imageData, err := base64.StdEncoding.DecodeString(logoReq.ImagenBase64)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar imagen base64: %w", err)
	}

	// Comenzar transacción
	tx, err := dbConn.Begin()
	if err != nil {
		return nil, fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	// Desactivar logos anteriores del usuario
	_, err = tx.Exec("UPDATE logos_empresas SET estado_logo = 0 WHERE id_usuario = ?", logoReq.IdUsuario)
	if err != nil {
		return nil, fmt.Errorf("error al desactivar logos anteriores: %w", err)
	}

	// Insertar nuevo logo
	query := `
		INSERT INTO logos_empresas (id_usuario, nombre_logo, imagen_logo, tipo_mime, tamaño_archivo, estado_logo)
		VALUES (?, ?, ?, ?, ?, 1)
	`

	result, err := tx.Exec(query, logoReq.IdUsuario, logoReq.NombreLogo, imageData, logoReq.TipoMime, len(imageData))
	if err != nil {
		return nil, fmt.Errorf("error al insertar logo: %w", err)
	}

	// Obtener ID del logo insertado
	logoID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ID del logo: %w", err)
	}

	// Confirmar transacción
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error al confirmar transacción: %w", err)
	}

	// Buscar y devolver el logo guardado
	return s.ObtenerLogoPorID(int(logoID))
}

// ObtenerLogoPorID obtiene un logo por su ID
func (s *LogoService) ObtenerLogoPorID(id int) (*models.LogoEmpresa, error) {
	// Conectar a la base de datos
	dbConn, err := db.ConnectUserDB()
	if err != nil {
		return nil, fmt.Errorf(errConectarDB, err)
	}
	defer dbConn.Close()

	query := `
		SELECT id, id_usuario, nombre_logo, imagen_logo, tipo_mime, tamaño_archivo, estado_logo, fecha_creacion, fecha_actualizacion
		FROM logos_empresas 
		WHERE id = ?
	`

	logo := &models.LogoEmpresa{
		IdEmpresa: nil, // Establecer nil por defecto
	}
	err = dbConn.QueryRow(query, id).Scan(
		&logo.ID,
		&logo.IdUsuario,
		&logo.NombreLogo,
		&logo.ImagenLogo,
		&logo.TipoMime,
		&logo.TamañoArchivo,
		&logo.EstadoLogo,
		&logo.FechaCreacion,
		&logo.FechaActualizacion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("logo no encontrado")
		}
		return nil, fmt.Errorf("error al obtener logo: %w", err)
	}

	return logo, nil
}

// ObtenerLogoActivoPorUsuario obtiene el logo activo de un usuario
func (s *LogoService) ObtenerLogoActivoPorUsuario(idUsuario int) (*models.LogoEmpresa, error) {
	// Conectar a la base de datos
	dbConn, err := db.ConnectUserDB()
	if err != nil {
		return nil, fmt.Errorf(errConectarDB, err)
	}
	defer dbConn.Close()

	query := `
		SELECT id, id_usuario, nombre_logo, imagen_logo, tipo_mime, tamaño_archivo, estado_logo, fecha_creacion, fecha_actualizacion
		FROM logos_empresas 
		WHERE id_usuario = ? AND estado_logo = 1
		ORDER BY fecha_creacion DESC
		LIMIT 1
	`

	logo := &models.LogoEmpresa{
		IdEmpresa: nil, // Establecer nil por defecto
	}
	err = dbConn.QueryRow(query, idUsuario).Scan(
		&logo.ID,
		&logo.IdUsuario,
		&logo.NombreLogo,
		&logo.ImagenLogo,
		&logo.TipoMime,
		&logo.TamañoArchivo,
		&logo.EstadoLogo,
		&logo.FechaCreacion,
		&logo.FechaActualizacion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no se encontró logo activo para el usuario")
		}
		return nil, fmt.Errorf("error al obtener logo activo: %w", err)
	}

	return logo, nil
}

// ListarLogosPorUsuario lista todos los logos de un usuario
func (s *LogoService) ListarLogosPorUsuario(idUsuario int) ([]models.LogoResponse, error) {
	// Conectar a la base de datos
	dbConn, err := db.ConnectUserDB()
	if err != nil {
		return nil, fmt.Errorf(errConectarDB, err)
	}
	defer dbConn.Close()

	query := `
		SELECT id, id_usuario, nombre_logo, tipo_mime, tamaño_archivo, estado_logo, fecha_creacion
		FROM logos_empresas 
		WHERE id_usuario = ?
		ORDER BY fecha_creacion DESC
	`

	rows, err := dbConn.Query(query, idUsuario)
	if err != nil {
		return nil, fmt.Errorf("error al consultar logos: %w", err)
	}
	defer rows.Close()

	var logos []models.LogoResponse
	for rows.Next() {
		var logo models.LogoResponse
		err := rows.Scan(
			&logo.ID,
			&logo.IdUsuario,
			&logo.NombreLogo,
			&logo.TipoMime,
			&logo.TamañoArchivo,
			&logo.EstadoLogo,
			&logo.FechaCreacion,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear logo: %w", err)
		}
		// Establecer id_empresa como nil
		logo.IdEmpresa = nil
		logos = append(logos, logo)
	}

	return logos, nil
}

// ActivarLogo activa un logo específico y desactiva los demás del usuario
func (s *LogoService) ActivarLogo(idLogo, idUsuario int) error {
	// Conectar a la base de datos
	dbConn, err := db.ConnectUserDB()
	if err != nil {
		return fmt.Errorf(errConectarDB, err)
	}
	defer dbConn.Close()

	tx, err := dbConn.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	// Desactivar todos los logos del usuario
	_, err = tx.Exec("UPDATE logos_empresas SET estado_logo = 0 WHERE id_usuario = ?", idUsuario)
	if err != nil {
		return fmt.Errorf("error al desactivar logos: %w", err)
	}

	// Activar el logo específico
	result, err := tx.Exec("UPDATE logos_empresas SET estado_logo = 1 WHERE id = ? AND id_usuario = ?", idLogo, idUsuario)
	if err != nil {
		return fmt.Errorf("error al activar logo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("logo no encontrado o no pertenece al usuario")
	}

	return tx.Commit()
}

// EliminarLogo elimina un logo de la base de datos
func (s *LogoService) EliminarLogo(idLogo, idUsuario int) error {
	// Conectar a la base de datos
	dbConn, err := db.ConnectUserDB()
	if err != nil {
		return fmt.Errorf(errConectarDB, err)
	}
	defer dbConn.Close()

	result, err := dbConn.Exec("DELETE FROM logos_empresas WHERE id = ? AND id_usuario = ?", idLogo, idUsuario)
	if err != nil {
		return fmt.Errorf("error al eliminar logo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("logo no encontrado o no pertenece al usuario")
	}

	return nil
}

// ValidarTipoImagen valida si el tipo MIME es una imagen válida
func ValidarTipoImagen(tipoMime string) bool {
	tiposValidos := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	}

	tipoMime = strings.ToLower(tipoMime)
	for _, tipo := range tiposValidos {
		if tipo == tipoMime {
			return true
		}
	}
	return false
}

// ObtenerExtensionPorTipo obtiene la extensión de archivo basada en el tipo MIME
func ObtenerExtensionPorTipo(tipoMime string) string {
	switch strings.ToLower(tipoMime) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	default:
		return ".png"
	}
}

// CargarLogoDesdeBaseDatos carga el logo activo de un usuario desde la base de datos
func (s *LogoService) CargarLogoDesdeBaseDatos(idUsuario int) ([]byte, error) {
	logo, err := s.ObtenerLogoActivoPorUsuario(idUsuario)
	if err != nil {
		log.Printf("No se encontró logo activo para el usuario %d: %v", idUsuario, err)
		return nil, err
	}

	if len(logo.ImagenLogo) == 0 {
		return nil, fmt.Errorf("imagen de logo vacía")
	}

	log.Printf("Logo cargado desde BD para usuario %d: %s (%d bytes)", idUsuario, logo.NombreLogo, len(logo.ImagenLogo))
	return logo.ImagenLogo, nil
}
