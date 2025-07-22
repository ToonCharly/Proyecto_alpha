package handlers

import (
	"Facts/internal/db"
	"Facts/internal/models"
	"Facts/internal/pac"
	"Facts/internal/services"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
)

// TimbrarFacturaHandler recibe una FacturaCFDI, timbra y retorna el resultado
func TimbrarFacturaHandler(factura models.FacturaCFDI) (map[string]interface{}, error) {
	// 1. Llenar datos fiscales del emisor
	err := LlenarDatosEmisor(&factura, factura.IdUsuario)
	if err != nil {
		return nil, err
	}

	// 2. Generar folio automático si no existe
	if factura.NumeroFolio == "" {
		if err := factura.GenerarFolioAutomatico(); err != nil {
			return nil, err
		}
	}
	if err := factura.ValidarFolio(); err != nil {
		return nil, err
	}

	// 3. Generar XML CFDI usando el generador
	xmlCFDI, err := services.ProcesarKeyYGenerarCFDI(factura, factura.KeyPath, factura.ClaveCSD, "")
	if err != nil {
		return nil, err
	}

	// 4. Timbrar el XML con el PAC y extraer el timbre fiscal digital usando la función integrada
	timbre, xmlTimbrado, err := services.TimbrarFactura(
		xmlCFDI,
		factura.EmisorRFC,
		factura.ClaveCSD,
		"https://api.pac.com/timbrar",
	)
	if err != nil {
		return nil, err
	}
	factura.Timbre = timbre

	// 6. Generar PDF usando el generador
	pdfBuf, _, err := services.GenerarPDF(factura, nil, nil)
	if err != nil {
		return nil, err
	}
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBuf.Bytes())

	// 7. Guardar en la base de datos
	resultado := map[string]interface{}{
		"uuid":   timbre.UUID,
		"xml":    string(xmlTimbrado),
		"pdf":    pdfBase64,
		"timbre": timbre,
		"folio":  factura.NumeroFolio,
	}
	if err := db.GuardarFacturaTimbrada(db.GetDB(), resultado); err != nil {
		return nil, err
	}
	return resultado, nil
}

// Handler HTTP para timbrar una factura CFDI 4.0 y devolver todos los datos relevantes
func TimbrarCFDIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer datos de la factura CFDI desde el body
	var factura models.FacturaCFDI
	if err := json.NewDecoder(r.Body).Decode(&factura); err != nil {
		http.Error(w, "Error al leer los datos de la factura", http.StatusBadRequest)
		return
	}

	// 1. Generar XML firmado CFDI
	xmlFirmado, err := services.ProcesarKeyYGenerarCFDI(factura, factura.KeyPath, factura.ClaveCSD, "")
	if err != nil {
		factura.LogError = "Error generando XML firmado: " + err.Error()
		resultado := map[string]interface{}{
			"error":     factura.LogError,
			"folio":     factura.NumeroFolio,
			"log_error": factura.LogError,
		}
		db.GuardarFacturaTimbrada(db.GetDB(), resultado)
		http.Error(w, factura.LogError, http.StatusInternalServerError)
		return
	}

	// 2. Timbrar el XML con el PAC (configuración real)
	pacURL := os.Getenv("PAC_URL")
	pacUser := os.Getenv("PAC_USER")
	pacPass := os.Getenv("PAC_PASS")
	if pacURL == "" || pacUser == "" || pacPass == "" {
		factura.LogError = "Configuración de PAC incompleta. Verifica las variables de entorno PAC_URL, PAC_USER y PAC_PASS."
		resultado := map[string]interface{}{
			"error":     factura.LogError,
			"folio":     factura.NumeroFolio,
			"log_error": factura.LogError,
		}
		db.GuardarFacturaTimbrada(db.GetDB(), resultado)
		http.Error(w, factura.LogError, http.StatusInternalServerError)
		return
	}
	xmlTimbrado, err := pac.TimbrarConPAC(string(xmlFirmado), pacURL, pacUser, pacPass)
	if err != nil {
		factura.LogError = "Error al timbrar con PAC: " + err.Error()
		resultado := map[string]interface{}{
			"error":     factura.LogError,
			"folio":     factura.NumeroFolio,
			"log_error": factura.LogError,
		}
		db.GuardarFacturaTimbrada(db.GetDB(), resultado)
		http.Error(w, factura.LogError, http.StatusInternalServerError)
		return
	}

	// 3. Extraer timbre fiscal digital
	timbre, err := services.ExtraerTimbreFiscalDigital(xmlTimbrado)
	if err != nil {
		factura.LogError = "Error extrayendo timbre fiscal: " + err.Error()
		resultado := map[string]interface{}{
			"error":     factura.LogError,
			"folio":     factura.NumeroFolio,
			"log_error": factura.LogError,
		}
		db.GuardarFacturaTimbrada(db.GetDB(), resultado)
		http.Error(w, factura.LogError, http.StatusInternalServerError)
		return
	}
	// Convertir timbre fiscal digital de services a models
	factura.Timbre = &models.TimbreFiscalDigital{
		UUID:             timbre.UUID,
		FechaTimbrado:    timbre.FechaTimbrado,
		RfcProvCertif:    timbre.NoCertificadoSAT, // Si el campo correcto es SelloSAT, ajústalo aquí
		SelloCFD:         timbre.SelloCFD,
		NoCertificadoSAT: timbre.NoCertificadoSAT,
		SelloSAT:         timbre.SelloSAT,
	}

	// 4. Generar PDF con el timbre fiscal digital
	pdfBuf, _, err := services.GenerarPDF(factura, nil, nil)
	if err != nil {
		factura.LogError = "Error generando PDF: " + err.Error()
		resultado := map[string]interface{}{
			"error":     factura.LogError,
			"folio":     factura.NumeroFolio,
			"log_error": factura.LogError,
		}
		db.GuardarFacturaTimbrada(db.GetDB(), resultado)
		http.Error(w, factura.LogError, http.StatusInternalServerError)
		return
	}
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBuf.Bytes())

	// 5. Guardar en la base de datos (todos los datos relevantes)
	resultado := map[string]interface{}{
		"uuid":                  timbre.UUID,
		"xml":                   string(xmlTimbrado),
		"pdf":                   pdfBase64,
		"timbre":                timbre,
		"numero_folio":          factura.NumeroFolio,
		"log_error":             factura.LogError,
		"emisor_rfc":            factura.EmisorRFC,
		"emisor_razon_social":   factura.EmisorRazonSocial,
		"receptor_rfc":          factura.ReceptorRFC,
		"receptor_razon_social": factura.ReceptorRazonSocial,
		"serie":                 factura.Serie,
		"fecha_emision":         factura.FechaEmision,
		"moneda":                factura.Moneda,
		"metodo_pago":           factura.MetodoPago,
		"forma_pago":            factura.FormaPago,
		// Agrega aquí más campos si lo necesitas
	}

	// Log de timbrado exitoso
	if timbre != nil && timbre.UUID != "" {
		println("TIMBRADO EXITOSO | UUID:", timbre.UUID, "| Fecha:", timbre.FechaTimbrado, "| SelloSAT:", timbre.SelloSAT, "| SelloCFD:", timbre.SelloCFD, "| CertificadoSAT:", timbre.NoCertificadoSAT)
	}
	if err := db.GuardarFactura(db.GetDB(), resultado); err != nil {
		factura.LogError = "Error guardando factura: " + err.Error()
		resultado["error"] = factura.LogError
		db.GuardarFactura(db.GetDB(), resultado)
		http.Error(w, factura.LogError, http.StatusInternalServerError)
		return
	}

	// 6. Responder con todos los datos relevantes
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}
