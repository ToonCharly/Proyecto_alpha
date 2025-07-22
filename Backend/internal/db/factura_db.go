package db

import (
	"database/sql"
	"log"
)

// GuardarFactura almacena TODOS los datos relevantes de la factura en la base de datos
func GuardarFactura(db *sql.DB, factura map[string]interface{}) error {
	// Extrae todos los campos relevantes del modelo Factura
	folio, _ := factura["numero_folio"].(string)
	serie, _ := factura["serie"].(string)
	fechaEmision, _ := factura["fecha_emision"].(string)
	sello, _ := factura["sello"].(string)
	noCertificado, _ := factura["no_certificado"].(string)
	certificado, _ := factura["certificado"].(string)
	subTotal, _ := factura["subtotal"].(string)
	total, _ := factura["total"].(string)
	tipoComprobante, _ := factura["tipo_comprobante"].(string)
	metodoPago, _ := factura["metodo_pago"].(string)
	formaPago, _ := factura["forma_pago"].(string)
	lugarExpedicion, _ := factura["lugar_expedicion"].(string)
	moneda, _ := factura["moneda"].(string)
	emisorRFC, _ := factura["emisor_rfc"].(string)
	emisorRazonSocial, _ := factura["emisor_razon_social"].(string)
	emisorRegimenFiscal, _ := factura["emisor_regimen_fiscal"].(string)
	receptorRFC, _ := factura["receptor_rfc"].(string)
	receptorRazonSocial, _ := factura["receptor_razon_social"].(string)
	domicilioFiscal, _ := factura["domicilio_fiscal"].(string)
	regimenFiscalReceptor, _ := factura["regimen_fiscal_receptor"].(string)
	usoCFDI, _ := factura["uso_cfdi"].(string)
	impuestos, _ := factura["impuestos"].(string)
	iva, _ := factura["iva"].(string)
	uuid, _ := factura["uuid"].(string)
	xmlStr, _ := factura["xml"].(string)
	pdfStr, _ := factura["pdf"].(string)
	logError, _ := factura["log_error"].(string)
	// Campos del timbre
	var fechaTimbrado, rfcProvCertif, selloCFD, noCertificadoSAT, selloSAT string
	if t, ok := factura["timbre"].(map[string]interface{}); ok {
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

	query := `INSERT INTO facturas (
			   numero_folio, serie, fecha_emision, sello, no_certificado, certificado, subtotal, total, tipo_comprobante,
			   metodo_pago, forma_pago, lugar_expedicion, moneda, emisor_rfc, emisor_razon_social, emisor_regimen_fiscal,
			   receptor_rfc, receptor_razon_social, domicilio_fiscal, regimen_fiscal_receptor, uso_cfdi, impuestos, iva,
			   uuid, xml, pdf, log_error, fecha_timbrado, rfc_prov_certif, sello_cfd, no_certificado_sat, sello_sat
	   ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	   ON DUPLICATE KEY UPDATE
			   serie=VALUES(serie), fecha_emision=VALUES(fecha_emision), sello=VALUES(sello), no_certificado=VALUES(no_certificado),
			   certificado=VALUES(certificado), subtotal=VALUES(subtotal), total=VALUES(total), tipo_comprobante=VALUES(tipo_comprobante),
			   metodo_pago=VALUES(metodo_pago), forma_pago=VALUES(forma_pago), lugar_expedicion=VALUES(lugar_expedicion), moneda=VALUES(moneda),
			   emisor_rfc=VALUES(emisor_rfc), emisor_razon_social=VALUES(emisor_razon_social), emisor_regimen_fiscal=VALUES(emisor_regimen_fiscal),
			   receptor_rfc=VALUES(receptor_rfc), receptor_razon_social=VALUES(receptor_razon_social), domicilio_fiscal=VALUES(domicilio_fiscal),
			   regimen_fiscal_receptor=VALUES(regimen_fiscal_receptor), uso_cfdi=VALUES(uso_cfdi), impuestos=VALUES(impuestos), iva=VALUES(iva),
			   uuid=VALUES(uuid), xml=VALUES(xml), pdf=VALUES(pdf), log_error=VALUES(log_error),
			   fecha_timbrado=VALUES(fecha_timbrado), rfc_prov_certif=VALUES(rfc_prov_certif), sello_cfd=VALUES(sello_cfd),
			   no_certificado_sat=VALUES(no_certificado_sat), sello_sat=VALUES(sello_sat)`
	_, err := db.Exec(query, folio, serie, fechaEmision, sello, noCertificado, certificado, subTotal, total, tipoComprobante,
		metodoPago, formaPago, lugarExpedicion, moneda, emisorRFC, emisorRazonSocial, emisorRegimenFiscal,
		receptorRFC, receptorRazonSocial, domicilioFiscal, regimenFiscalReceptor, usoCFDI, impuestos, iva,
		uuid, xmlStr, pdfStr, logError, fechaTimbrado, rfcProvCertif, selloCFD, noCertificadoSAT, selloSAT)
	if err != nil {
		log.Printf("Error guardando factura en BD: %v", err)
		return err
	}
	return nil
}
