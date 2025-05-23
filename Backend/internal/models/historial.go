package models

import (
    "carlos/Facts/Backend/internal/db"
    "database/sql" // Añadir esta importación
    "fmt"          // Añadir esta importación
)

// HistorialFactura representa una entrada en el historial de facturas
type HistorialFactura struct {
    ID                  int       `json:"id"`
    IDUsuario           int       `json:"id_usuario"`
    RFCReceptor         string    `json:"rfc_receptor"`
    RazonSocialReceptor string    `json:"razon_social_receptor"`
    ClaveTicket         string    `json:"clave_ticket"`
    Total               float64   `json:"total"`
    UsoCFDI             string    `json:"uso_cfdi"`
    FechaGeneracion     string    `json:"fecha_generacion"` // Cambiado de time.Time a string
    Estado              string    `json:"estado"`
    Observaciones       string    `json:"observaciones"`
}

// InsertarHistorialFactura inserta una nueva entrada en el historial de facturas
func InsertarHistorialFactura(idUsuario int, rfcReceptor string, razonSocialReceptor string, 
                            claveTicket string, total float64, usoCFDI string, observaciones string) (int64, error) {
    dbConn := db.GetDB()
    
    result, err := dbConn.Exec(
        `INSERT INTO historial_facturas 
        (id_usuario, rfc_receptor, razon_social_receptor, clave_ticket, total, uso_cfdi, observaciones) 
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
        idUsuario, rfcReceptor, razonSocialReceptor, claveTicket, total, usoCFDI, observaciones,
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
        clave_ticket, total, uso_cfdi, 
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
        clave_ticket, total, uso_cfdi, 
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