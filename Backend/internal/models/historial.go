package models

import (
	"Facts/internal/db"
	"database/sql" // Añadir esta importación
	"fmt"          // Añadir esta importación
)

// HistorialFactura representa una entrada en el historial de facturas
type HistorialFactura struct {
	ID                  int     `json:"id"`
	IDUsuario           int     `json:"id_usuario"`
	RFCReceptor         string  `json:"rfc_receptor"`
	RazonSocialReceptor string  `json:"razon_social_receptor"`
	ClaveTicket         string  `json:"clave_ticket"`
	NumeroFolio         string  `json:"numero_folio"` // Ahora es NumeroFolio en vez de Folio
	Total               float64 `json:"total"`
	UsoCFDI             string  `json:"uso_cfdi"`
	FechaGeneracion     string  `json:"fecha_generacion"`
	Estado              string  `json:"estado"`
	Observaciones       string  `json:"observaciones"`
}

// InsertarHistorialFactura inserta una nueva entrada en el historial de facturas
func InsertarHistorialFactura(idUsuario int, rfcReceptor string, razonSocialReceptor string,
	claveTicket string, numeroFolio string, total float64, usoCFDI string, observaciones string) (int64, error) {
	dbConn := db.GetDB()

	result, err := dbConn.Exec(
		`INSERT INTO historial_facturas 
		(id_usuario, rfc_receptor, razon_social_receptor, clave_ticket, folio, total, uso_cfdi, observaciones) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		idUsuario, rfcReceptor, razonSocialReceptor, claveTicket, numeroFolio, total, usoCFDI, observaciones,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// ObtenerHistorialFacturasPorUsuario obtiene todas las facturas generadas por un usuario
func ObtenerHistorialFacturasPorUsuario(idUsuario int) ([]HistorialFactura, error) {
	dbConn := db.GetDB()

	rows, err := dbConn.Query(
		`SELECT id, id_usuario, rfc_receptor, razon_social_receptor, 
		clave_ticket, folio AS numero_folio, total, uso_cfdi, 
		DATE_FORMAT(fecha_generacion, '%Y-%m-%d %H:%i:%s') as fecha_generacion, 
		estado, observaciones 
		FROM historial_facturas 
		WHERE id_usuario = ? 
		ORDER BY fecha_generacion DESC`,
		idUsuario,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var facturas []HistorialFactura

	for rows.Next() {
		var factura HistorialFactura

		err := rows.Scan(
			&factura.ID,
			&factura.IDUsuario,
			&factura.RFCReceptor,
			&factura.RazonSocialReceptor,
			&factura.ClaveTicket,
			&factura.NumeroFolio,
			&factura.Total,
			&factura.UsoCFDI,
			&factura.FechaGeneracion,
			&factura.Estado,
			&factura.Observaciones,
		)

		if err != nil {
			return nil, err
		}

		facturas = append(facturas, factura)
	}

	return facturas, nil
}

// ObtenerFacturaPorID obtiene una factura específica del historial por su ID
func ObtenerFacturaPorID(id int) (*HistorialFactura, error) {
	dbConn := db.GetDB()

	var factura HistorialFactura

	err := dbConn.QueryRow(
		`SELECT id, id_usuario, rfc_receptor, razon_social_receptor, 
		clave_ticket, folio AS numero_folio, total, uso_cfdi, 
		DATE_FORMAT(fecha_generacion, '%Y-%m-%d %H:%i:%s') as fecha_generacion, 
		estado, observaciones 
		FROM historial_facturas 
		WHERE id = ?`,
		id,
	).Scan(
		&factura.ID,
		&factura.IDUsuario,
		&factura.RFCReceptor,
		&factura.RazonSocialReceptor,
		&factura.ClaveTicket,
		&factura.NumeroFolio,
		&factura.Total,
		&factura.UsoCFDI,
		&factura.FechaGeneracion,
		&factura.Estado,
		&factura.Observaciones,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no se encontró la factura con ID %d", id)
		}
		return nil, err
	}

	return &factura, nil
}

// BuscarHistorialFacturas busca facturas en el historial por múltiples criterios
func BuscarHistorialFacturas(idUsuario int, folio, rfcReceptor, razonSocial string) ([]HistorialFactura, error) {
	dbConn := db.GetDB()

	// Construir la consulta dinámicamente basándose en los criterios proporcionados
	query := `SELECT id, id_usuario, rfc_receptor, razon_social_receptor, 
		clave_ticket, folio, total, uso_cfdi, 
		DATE_FORMAT(fecha_generacion, '%Y-%m-%d %H:%i:%s') as fecha_generacion, 
		estado, observaciones 
		FROM historial_facturas 
		WHERE id_usuario = ?`

	var args []interface{}
	args = append(args, idUsuario)

	// Agregar condiciones de búsqueda si se proporcionan
	if folio != "" {
		query += " AND folio LIKE ?"
		args = append(args, "%"+folio+"%")
	}

	if rfcReceptor != "" {
		query += " AND rfc_receptor LIKE ?"
		args = append(args, "%"+rfcReceptor+"%")
	}

	if razonSocial != "" {
		query += " AND razon_social_receptor LIKE ?"
		args = append(args, "%"+razonSocial+"%")
	}

	query += " ORDER BY fecha_generacion DESC"

	// Log para debugging
	fmt.Printf("Ejecutando consulta: %s\nCon argumentos: %v\n", query, args)

	rows, err := dbConn.Query(query, args...)
	if err != nil {
		fmt.Printf("Error en la consulta: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var facturas []HistorialFactura

	for rows.Next() {
		var factura HistorialFactura

		err := rows.Scan(
			&factura.ID,
			&factura.IDUsuario,
			&factura.RFCReceptor,
			&factura.RazonSocialReceptor,
			&factura.ClaveTicket,
			&factura.NumeroFolio,
			&factura.Total,
			&factura.UsoCFDI,
			&factura.FechaGeneracion,
			&factura.Estado,
			&factura.Observaciones,
		)

		if err != nil {
			fmt.Printf("Error al escanear fila: %v\n", err)
			return nil, err
		}

		facturas = append(facturas, factura)
	}

	fmt.Printf("Búsqueda completada: %d facturas encontradas\n", len(facturas))
	return facturas, nil
}

// ObtenerHistorialFacturasPorUsuarioConPaginacion obtiene facturas con paginación
func ObtenerHistorialFacturasPorUsuarioConPaginacion(idUsuario, pagina, limite int) ([]HistorialFactura, int, error) {
	dbConn := db.GetDB()

	// Calcular el offset
	offset := (pagina - 1) * limite

	// Obtener el total de facturas para este usuario
	var totalFacturas int
	err := dbConn.QueryRow("SELECT COUNT(*) FROM historial_facturas WHERE id_usuario = ?", idUsuario).Scan(&totalFacturas)
	if err != nil {
		return nil, 0, err
	}

	// Obtener las facturas con paginación
	rows, err := dbConn.Query(
		`SELECT id, id_usuario, rfc_receptor, razon_social_receptor, 
		clave_ticket, folio, total, uso_cfdi, 
		DATE_FORMAT(fecha_generacion, '%Y-%m-%d %H:%i:%s') as fecha_generacion, 
		estado, observaciones 
		FROM historial_facturas 
		WHERE id_usuario = ? 
		ORDER BY fecha_generacion DESC
		LIMIT ? OFFSET ?`,
		idUsuario, limite, offset,
	)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var facturas []HistorialFactura

	for rows.Next() {
		var factura HistorialFactura

		err := rows.Scan(
			&factura.ID,
			&factura.IDUsuario,
			&factura.RFCReceptor,
			&factura.RazonSocialReceptor,
			&factura.ClaveTicket,
			&factura.NumeroFolio,
			&factura.Total,
			&factura.UsoCFDI,
			&factura.FechaGeneracion,
			&factura.Estado,
			&factura.Observaciones,
		)

		if err != nil {
			return nil, 0, err
		}

		facturas = append(facturas, factura)
	}

	return facturas, totalFacturas, nil
}

// BuscarHistorialFacturasConPaginacion busca facturas en el historial por múltiples criterios con paginación
func BuscarHistorialFacturasConPaginacion(idUsuario int, folio, rfcReceptor, razonSocial string, pagina, limite int) ([]HistorialFactura, int, error) {
	dbConn := db.GetDB()

	// Calcular el offset
	offset := (pagina - 1) * limite

	// Construir la consulta base para el conteo
	queryCount := `SELECT COUNT(*) FROM historial_facturas WHERE id_usuario = ?`
	querySelect := `SELECT id, id_usuario, rfc_receptor, razon_social_receptor, 
		clave_ticket, folio, total, uso_cfdi, 
		DATE_FORMAT(fecha_generacion, '%Y-%m-%d %H:%i:%s') as fecha_generacion, 
		estado, observaciones 
		FROM historial_facturas 
		WHERE id_usuario = ?`

	var args []interface{}
	args = append(args, idUsuario)

	// Agregar condiciones de búsqueda si se proporcionan
	if folio != "" {
		queryCount += " AND folio LIKE ?"
		querySelect += " AND folio LIKE ?"
		args = append(args, "%"+folio+"%")
	}

	if rfcReceptor != "" {
		queryCount += " AND rfc_receptor LIKE ?"
		querySelect += " AND rfc_receptor LIKE ?"
		args = append(args, "%"+rfcReceptor+"%")
	}

	if razonSocial != "" {
		queryCount += " AND razon_social_receptor LIKE ?"
		querySelect += " AND razon_social_receptor LIKE ?"
		args = append(args, "%"+razonSocial+"%")
	}

	// Obtener el total de facturas que coinciden con los criterios
	var totalFacturas int
	err := dbConn.QueryRow(queryCount, args...).Scan(&totalFacturas)
	if err != nil {
		return nil, 0, err
	}

	// Agregar ORDER BY y LIMIT a la consulta de selección
	querySelect += " ORDER BY fecha_generacion DESC LIMIT ? OFFSET ?"
	argsWithPagination := append(args, limite, offset)

	rows, err := dbConn.Query(querySelect, argsWithPagination...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var facturas []HistorialFactura

	for rows.Next() {
		var factura HistorialFactura

		err := rows.Scan(
			&factura.ID,
			&factura.IDUsuario,
			&factura.RFCReceptor,
			&factura.RazonSocialReceptor,
			&factura.ClaveTicket,
			&factura.NumeroFolio,
			&factura.Total,
			&factura.UsoCFDI,
			&factura.FechaGeneracion,
			&factura.Estado,
			&factura.Observaciones,
		)

		if err != nil {
			return nil, 0, err
		}

		facturas = append(facturas, factura)
	}

	return facturas, totalFacturas, nil
}
