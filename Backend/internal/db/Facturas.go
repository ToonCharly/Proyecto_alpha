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

// GuardarFacturaTimbrada guarda el resultado del timbrado en la base de datos
func GuardarFacturaTimbrada(db *sql.DB, resultado map[string]interface{}) error {
	// Guarda XML, PDF, UUID, timbre y folio en la tabla facturas, y cada campo del timbre en columnas separadas
	uuid, _ := resultado["uuid"].(string)
	xmlStr, _ := resultado["xml"].(string)
	pdfStr, _ := resultado["pdf"].(string)
	folio, _ := resultado["folio"].(string)
	var (
		fechaTimbrado, rfcProvCertif, selloCFD, noCertificadoSAT, selloSAT string
	)
	if t, ok := resultado["timbre"].(map[string]interface{}); ok {
		if v, ok := t["FechaTimbrado"].(string); ok {
			fechaTimbrado = v
		}
		if v, ok := t["RfcProvCertif"].(string); ok {
			rfcProvCertif = v
		}
		if v, ok := t["SelloCFD"].(string); ok {
			selloCFD = v
		}
		if v, ok := t["NoCertificadoSAT"].(string); ok {
			noCertificadoSAT = v
		}
		if v, ok := t["SelloSAT"].(string); ok {
			selloSAT = v
		}
	}

	// Inserta o actualiza la factura timbrada con los campos del timbre
	query := `INSERT INTO facturas (
		numero_folio, uuid, xml, pdf,
		fecha_timbrado, rfc_prov_certif, sello_cfd, no_certificado_sat, sello_sat
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
		uuid=VALUES(uuid), xml=VALUES(xml), pdf=VALUES(pdf),
		fecha_timbrado=VALUES(fecha_timbrado), rfc_prov_certif=VALUES(rfc_prov_certif),
		sello_cfd=VALUES(sello_cfd), no_certificado_sat=VALUES(no_certificado_sat), sello_sat=VALUES(sello_sat)`
	_, err := db.Exec(query, folio, uuid, xmlStr, pdfStr,
		fechaTimbrado, rfcProvCertif, selloCFD, noCertificadoSAT, selloSAT)
	if err != nil {
		log.Printf("Error guardando factura timbrada en BD: %v", err)
		return err
	}
	return nil
}
