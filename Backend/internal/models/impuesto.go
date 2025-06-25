package models

import (
	"database/sql"
	"fmt"
	"log"
)

// Impuesto representa un impuesto en la base de datos
type Impuesto struct {
	IDIVA       int     `json:"idiva"`
	IDEmpresa   int     `json:"idempresa"`
	Descripcion string  `json:"descripcion"`
	IEPS1       float64 `json:"ieps1"`
	TipoIEPS1   string  `json:"tipo_ieps1"`
	IEPS2       float64 `json:"ieps2"`
	TipoIEPS2   string  `json:"tipo_ieps2"`
	IEPS3       float64 `json:"ieps3"`
	TipoIEPS3   string  `json:"tipo_ieps3"`
	IVA         float64 `json:"iva"`
	TipoIVA     string  `json:"tipo_iva"`
}

// Producto representa un producto con información de impuestos
type ProductoConImpuesto struct {
	IDProducto  int     `json:"idproducto"`
	IDEmpresa   int     `json:"idempresa"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	Clave       string  `json:"clave"`
	SATClave    string  `json:"sat_clave"`
	SATMedida   string  `json:"sat_medida"`
	// Información del impuesto
	IDIVA          int     `json:"idiva"`
	DescripcionImp string  `json:"descripcion_impuesto"`
	IEPS1          float64 `json:"ieps1"`
	TipoIEPS1      string  `json:"tipo_ieps1"`
	IEPS2          float64 `json:"ieps2"`
	TipoIEPS2      string  `json:"tipo_ieps2"`
	IEPS3          float64 `json:"ieps3"`
	TipoIEPS3      string  `json:"tipo_ieps3"`
	IVA            float64 `json:"iva"`
	TipoIVA        string  `json:"tipo_iva"`
}

// GetImpuestosByEmpresa obtiene todos los impuestos de una empresa
func GetImpuestosByEmpresa(db *sql.DB, idEmpresa int) ([]Impuesto, error) {
	query := `
		SELECT idiva, idempresa, descripcion, ieps1, tipo_ieps1, ieps2, tipo_ieps2, ieps3, tipo_ieps3, iva, tipo_iva
		FROM crm_impuestos 
		WHERE idempresa = ?
		ORDER BY descripcion
	`

	rows, err := db.Query(query, idEmpresa)
	if err != nil {
		log.Printf("Error al obtener impuestos: %v", err)
		return nil, fmt.Errorf("error al obtener impuestos: %w", err)
	}
	defer rows.Close()

	var impuestos []Impuesto
	for rows.Next() {
		var impuesto Impuesto
		err := rows.Scan(
			&impuesto.IDIVA,
			&impuesto.IDEmpresa,
			&impuesto.Descripcion,
			&impuesto.IEPS1,
			&impuesto.TipoIEPS1,
			&impuesto.IEPS2,
			&impuesto.TipoIEPS2,
			&impuesto.IEPS3,
			&impuesto.TipoIEPS3,
			&impuesto.IVA,
			&impuesto.TipoIVA,
		)
		if err != nil {
			log.Printf("Error al escanear impuesto: %v", err)
			continue
		}
		impuestos = append(impuestos, impuesto)
	}

	return impuestos, nil
}

// GetProductosConImpuestos obtiene productos con información de impuestos
func GetProductosConImpuestos(db *sql.DB, idEmpresa int) ([]ProductoConImpuesto, error) {
	query := `
		SELECT 
			p.idproducto,
			p.idempresa,
			p.descripcion,
			p.precio,
			p.clave,
			COALESCE(p.sat_clave, '') as sat_clave,
			COALESCE(p.sat_medida, '') as sat_medida,
			COALESCE(i.idiva, 0) as idiva,
			COALESCE(i.descripcion, '') as descripcion_impuesto,
			COALESCE(i.ieps1, 0) as ieps1,
			COALESCE(i.tipo_ieps1, '') as tipo_ieps1,
			COALESCE(i.ieps2, 0) as ieps2,
			COALESCE(i.tipo_ieps2, '') as tipo_ieps2,
			COALESCE(i.ieps3, 0) as ieps3,
			COALESCE(i.tipo_ieps3, '') as tipo_ieps3,
			COALESCE(i.iva, 0) as iva,
			COALESCE(i.tipo_iva, '') as tipo_iva
		FROM crm_productos p
		LEFT JOIN crm_impuestos i ON p.idempresa = i.idempresa
		WHERE p.idempresa = ?
		ORDER BY p.descripcion
	`

	rows, err := db.Query(query, idEmpresa)
	if err != nil {
		log.Printf("Error al obtener productos con impuestos: %v", err)
		return nil, fmt.Errorf("error al obtener productos con impuestos: %w", err)
	}
	defer rows.Close()

	var productos []ProductoConImpuesto
	for rows.Next() {
		var producto ProductoConImpuesto
		err := rows.Scan(
			&producto.IDProducto,
			&producto.IDEmpresa,
			&producto.Descripcion,
			&producto.Precio,
			&producto.Clave,
			&producto.SATClave,
			&producto.SATMedida,
			&producto.IDIVA,
			&producto.DescripcionImp,
			&producto.IEPS1,
			&producto.TipoIEPS1,
			&producto.IEPS2,
			&producto.TipoIEPS2,
			&producto.IEPS3,
			&producto.TipoIEPS3,
			&producto.IVA,
			&producto.TipoIVA,
		)
		if err != nil {
			log.Printf("Error al escanear producto: %v", err)
			continue
		}
		productos = append(productos, producto)
	}

	return productos, nil
}

// GetImpuestoByID obtiene un impuesto específico por su ID
func GetImpuestoByID(db *sql.DB, idIVA int) (*Impuesto, error) {
	query := `
		SELECT idiva, idempresa, descripcion, ieps1, tipo_ieps1, ieps2, tipo_ieps2, ieps3, tipo_ieps3, iva, tipo_iva
		FROM crm_impuestos 
		WHERE idiva = ?
	`

	var impuesto Impuesto
	err := db.QueryRow(query, idIVA).Scan(
		&impuesto.IDIVA,
		&impuesto.IDEmpresa,
		&impuesto.Descripcion,
		&impuesto.IEPS1,
		&impuesto.TipoIEPS1,
		&impuesto.IEPS2,
		&impuesto.TipoIEPS2,
		&impuesto.IEPS3,
		&impuesto.TipoIEPS3,
		&impuesto.IVA,
		&impuesto.TipoIVA,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("impuesto no encontrado")
		}
		log.Printf("Error al obtener impuesto: %v", err)
		return nil, fmt.Errorf("error al obtener impuesto: %w", err)
	}

	return &impuesto, nil
}

// CreateImpuesto crea un nuevo impuesto
func CreateImpuesto(db *sql.DB, impuesto *Impuesto) error {
	query := `
		INSERT INTO crm_impuestos (idempresa, descripcion, ieps1, tipo_ieps1, ieps2, tipo_ieps2, ieps3, tipo_ieps3, iva, tipo_iva)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query,
		impuesto.IDEmpresa,
		impuesto.Descripcion,
		impuesto.IEPS1,
		impuesto.TipoIEPS1,
		impuesto.IEPS2,
		impuesto.TipoIEPS2,
		impuesto.IEPS3,
		impuesto.TipoIEPS3,
		impuesto.IVA,
		impuesto.TipoIVA)
	if err != nil {
		log.Printf("Error al crear impuesto: %v", err)
		return fmt.Errorf("error al crear impuesto: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error al obtener ID del impuesto creado: %w", err)
	}

	impuesto.IDIVA = int(id)
	return nil
}

// UpdateImpuesto actualiza un impuesto existente
func UpdateImpuesto(db *sql.DB, impuesto *Impuesto) error {
	query := `
		UPDATE crm_impuestos 
		SET descripcion = ?, ieps1 = ?, tipo_ieps1 = ?, ieps2 = ?, tipo_ieps2 = ?, ieps3 = ?, tipo_ieps3 = ?, iva = ?, tipo_iva = ?
		WHERE idiva = ? AND idempresa = ?
	`

	result, err := db.Exec(query,
		impuesto.Descripcion,
		impuesto.IEPS1,
		impuesto.TipoIEPS1,
		impuesto.IEPS2,
		impuesto.TipoIEPS2,
		impuesto.IEPS3,
		impuesto.TipoIEPS3,
		impuesto.IVA,
		impuesto.TipoIVA,
		impuesto.IDIVA,
		impuesto.IDEmpresa)
	if err != nil {
		log.Printf("Error al actualizar impuesto: %v", err)
		return fmt.Errorf("error al actualizar impuesto: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar actualización: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("impuesto no encontrado o no autorizado")
	}

	return nil
}

// DeleteImpuesto elimina un impuesto (eliminación física)
func DeleteImpuesto(db *sql.DB, idIVA, idEmpresa int) error {
	query := `
		DELETE FROM crm_impuestos 
		WHERE idiva = ? AND idempresa = ?
	`

	result, err := db.Exec(query, idIVA, idEmpresa)
	if err != nil {
		log.Printf("Error al eliminar impuesto: %v", err)
		return fmt.Errorf("error al eliminar impuesto: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar eliminación: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("impuesto no encontrado o no autorizado")
	}

	return nil
}
