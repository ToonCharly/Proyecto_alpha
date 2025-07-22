package models

import (
	"Facts/internal/db"
	"fmt"
	"log"
)

type Empresa struct {
	ID            int    `json:"id"`
	IdUsuario     int    `json:"id_usuario"`
	RFC           string `json:"rfc"`
	RazonSocial   string `json:"razon_social"`
	RegimenFiscal string `json:"regimen_fiscal"`
	Direccion     string `json:"direccion"`
	CodigoPostal  string `json:"codigo_postal"`
	Pais          string `json:"pais"`
	Estado        string `json:"estado"`
	Localidad     string `json:"localidad"`
	Municipio     string `json:"municipio"`
	Colonia       string `json:"colonia"`
	CreatedAt     string `json:"created_at"`
}

// Elimina una empresa de la base de datos por su ID
func EliminarEmpresa(id int) error {
	// Primero verificamos que la empresa exista
	var exists int
	err := db.GetDB().QueryRow("SELECT COUNT(*) FROM empresas WHERE id = ?", id).Scan(&exists)
	if err != nil {
		log.Printf("Error al verificar existencia de la empresa: %v", err)
		return fmt.Errorf("error al verificar existencia de la empresa: %w", err)
	}

	if exists == 0 {
		return fmt.Errorf("la empresa con ID %d no existe en la base de datos", id)
	}

	// Eliminamos la empresa
	query := "DELETE FROM empresas WHERE id = ?"

	result, err := db.GetDB().Exec(query, id)
	if err != nil {
		log.Printf("Error al eliminar la empresa: %v", err)
		return fmt.Errorf("error al eliminar la empresa: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no se pudo eliminar la empresa con ID %d", id)
	}

	log.Printf("Empresa con ID %d eliminada exitosamente", id)
	return nil
}

// Obtiene una empresa por su ID
func ObtenerEmpresaPorID(id int) (*Empresa, error) {
	query := `
		SELECT id, id_usuario, rfc, razon_social, regimen_fiscal, direccion,
			   codigo_postal, pais, estado, localidad, municipio, colonia, created_at
		FROM empresas
		WHERE id = ?
	`

	var empresa Empresa
	err := db.GetDB().QueryRow(query, id).Scan(
		&empresa.ID, &empresa.IdUsuario, &empresa.RFC, &empresa.RazonSocial,
		&empresa.RegimenFiscal, &empresa.Direccion, &empresa.CodigoPostal,
		&empresa.Pais, &empresa.Estado, &empresa.Localidad,
		&empresa.Municipio, &empresa.Colonia, &empresa.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error al obtener la empresa: %w", err)
	}

	return &empresa, nil
}

// Verifica si el usuario con el ID proporcionado existe en la base de datos
func UsuarioExiste(idUsuario int) (bool, error) {
	var count int
	err := db.GetDB().QueryRow("SELECT COUNT(*) FROM usuarios WHERE id = ?", idUsuario).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error al verificar el usuario: %w", err)
	}
	return count > 0, nil
}

// Inserta una nueva empresa en la base de datos
func InsertarEmpresa(empresa Empresa) (int64, error) {
	// Validar si el usuario existe
	existe, err := UsuarioExiste(empresa.IdUsuario)
	if err != nil {
		log.Printf("Error al verificar el usuario: %v", err)
		return 0, fmt.Errorf("error al verificar el usuario: %w", err)
	}
	if !existe {
		// Cambio crítico: en lugar de solo advertir, devolver un error
		log.Printf("Error: El usuario con ID %d no existe. No se puede insertar la empresa.", empresa.IdUsuario)
		return 0, fmt.Errorf("el usuario con ID %d no existe en la base de datos", empresa.IdUsuario)
	}

	// Query para insertar la empresa
	query := `
		INSERT INTO empresas (
			id_usuario, rfc, razon_social, regimen_fiscal, direccion,
			codigo_postal, pais, estado, localidad, municipio, colonia
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.GetDB().Exec(query,
		empresa.IdUsuario, empresa.RFC, empresa.RazonSocial, empresa.RegimenFiscal,
		empresa.Direccion, empresa.CodigoPostal, empresa.Pais, empresa.Estado,
		empresa.Localidad, empresa.Municipio, empresa.Colonia,
	)
	if err != nil {
		log.Printf("Error al insertar empresa: %v", err)
		return 0, fmt.Errorf("error al insertar empresa: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener el ID de la empresa insertada: %w", err)
	}

	return id, nil
}

// Actualiza una empresa existente en la base de datos
func ActualizarEmpresa(empresa Empresa) error {
	// Verificar que la empresa exista
	_, err := ObtenerEmpresaPorID(empresa.ID)
	if err != nil {
		return fmt.Errorf("la empresa no existe: %w", err)
	}

	// Query para actualizar la empresa
	query := `
		UPDATE empresas 
		SET rfc = ?, razon_social = ?, regimen_fiscal = ?, direccion = ?,
			codigo_postal = ?, pais = ?, estado = ?, localidad = ?,
			municipio = ?, colonia = ?
		WHERE id = ?
	`

	_, err = db.GetDB().Exec(query,
		empresa.RFC, empresa.RazonSocial, empresa.RegimenFiscal,
		empresa.Direccion, empresa.CodigoPostal, empresa.Pais,
		empresa.Estado, empresa.Localidad, empresa.Municipio,
		empresa.Colonia, empresa.ID)

	if err != nil {
		log.Printf("Error al actualizar empresa: %v", err)
		return fmt.Errorf("error al actualizar empresa: %w", err)
	}

	log.Printf("Empresa con ID %d actualizada exitosamente", empresa.ID)
	return nil
}

// ObtenerEmpresas devuelve todas las empresas
func ObtenerEmpresas() ([]Empresa, error) {
	query := `
		SELECT id, id_usuario, rfc, razon_social, regimen_fiscal, direccion,
			   codigo_postal, pais, estado, localidad, municipio, colonia, created_at
		FROM empresas
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las empresas: %w", err)
	}
	defer rows.Close()

	var empresas []Empresa
	for rows.Next() {
		var empresa Empresa
		err := rows.Scan(
			&empresa.ID, &empresa.IdUsuario, &empresa.RFC, &empresa.RazonSocial,
			&empresa.RegimenFiscal, &empresa.Direccion, &empresa.CodigoPostal,
			&empresa.Pais, &empresa.Estado, &empresa.Localidad,
			&empresa.Municipio, &empresa.Colonia, &empresa.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear las empresas: %w", err)
		}
		empresas = append(empresas, empresa)
	}

	return empresas, nil
}

// Obtiene todas las empresas asociadas a un usuario por su ID
func ObtenerEmpresasPorUsuario(idUsuario int) ([]Empresa, error) {
	query := `
		SELECT id, id_usuario, rfc, razon_social, regimen_fiscal, direccion,
			   codigo_postal, pais, estado, localidad, municipio, colonia, created_at
		FROM empresas
		WHERE id_usuario = ?
	`

	rows, err := db.GetDB().Query(query, idUsuario)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las empresas: %w", err)
	}
	defer rows.Close()

	var empresas []Empresa
	for rows.Next() {
		var empresa Empresa
		err := rows.Scan(
			&empresa.ID, &empresa.IdUsuario, &empresa.RFC, &empresa.RazonSocial,
			&empresa.RegimenFiscal, &empresa.Direccion, &empresa.CodigoPostal,
			&empresa.Pais, &empresa.Estado, &empresa.Localidad,
			&empresa.Municipio, &empresa.Colonia, &empresa.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear las empresas: %w", err)
		}
		empresas = append(empresas, empresa)
	}

	return empresas, nil
}

// ObtenerEmpresaEmisoraPorIdEmpresa obtiene los datos de la empresa emisora para una factura
func ObtenerEmpresaEmisoraPorIdEmpresa(idEmpresa int) (*Empresa, error) {
	if idEmpresa <= 0 {
		return nil, fmt.Errorf("ID de empresa inválido: %d", idEmpresa)
	}

	return ObtenerEmpresaPorID(idEmpresa)
}
