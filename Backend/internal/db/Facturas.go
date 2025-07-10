package db

import (
    "database/sql"
    "log"
)

// Devuelve el UUID y el No. de Certificado para un folio de factura dado
func ObtenerUUIDyNoCertificado(numeroFolio string) (uuid string, noCertificado string, err error) {
    db := GetDB()
    err = db.QueryRow(
        "SELECT uuid, no_certificado FROM facturas WHERE numero_folio = ?",
        numeroFolio,
    ).Scan(&uuid, &noCertificado)
    if err != nil && err != sql.ErrNoRows {
        log.Printf("Error consultando UUID y NoCertificado para folio %s: %v", numeroFolio, err)
    }
    return
}