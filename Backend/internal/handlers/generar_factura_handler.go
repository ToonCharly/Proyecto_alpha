package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"carlos/Facts/Backend/internal/models"
	"carlos/Facts/Backend/internal/services"
)

// procesarDatosFactura extrae los datos de la factura del request
func procesarDatosFactura(r *http.Request) (models.Factura, []byte, error) {
	var factura models.Factura
	var plantillaBytes []byte

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return factura, nil, fmt.Errorf("Content-Type no especificado")
	}

	if strings.Contains(contentType, "multipart/form-data") {
		return procesarMultipartForm(r)
	} else if strings.Contains(contentType, "application/json") {
		return procesarJSON(r)
	}

	return factura, plantillaBytes, fmt.Errorf("tipo de contenido no soportado")
}

// procesarMultipartForm procesa datos multipart/form-data
func procesarMultipartForm(r *http.Request) (models.Factura, []byte, error) {
	var factura models.Factura
	var plantillaBytes []byte

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return factura, nil, fmt.Errorf("error al parsear multipart form: %v", err)
	}

	// Leer datos de factura
	facturaData := r.FormValue("datos")
	if facturaData == "" {
		return factura, nil, fmt.Errorf("no se encontraron datos de factura")
	}

	if err := json.Unmarshal([]byte(facturaData), &factura); err != nil {
		return factura, nil, fmt.Errorf("error al decodificar JSON: %v", err)
	}

	// Leer plantilla si existe
	plantillaFile, _, err := r.FormFile("plantilla")
	if err == nil {
		defer plantillaFile.Close()
		plantillaBytes, err = io.ReadAll(plantillaFile)
		if err != nil {
			return factura, nil, fmt.Errorf("error al leer la plantilla: %v", err)
		}
	}

	return factura, plantillaBytes, nil
}

// procesarJSON procesa datos JSON
func procesarJSON(r *http.Request) (models.Factura, []byte, error) {
	var factura models.Factura

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return factura, nil, fmt.Errorf("error al leer el cuerpo de la solicitud: %v", err)
	}

	if err := json.Unmarshal(body, &factura); err != nil {
		return factura, nil, fmt.Errorf("error al decodificar JSON: %v", err)
	}

	return factura, nil, nil
}

// manejarFolio genera o valida el folio de la factura
func manejarFolio(factura *models.Factura) error {
	if factura.NumeroFolio == "" {
		err := factura.GenerarFolioAutomatico()
		if err != nil {
			return fmt.Errorf("error al generar folio automático: %v", err)
		}
		log.Printf("Folio generado automáticamente: %s", factura.NumeroFolio)
	} else {
		err := factura.ValidarFolio()
		if err != nil {
			return fmt.Errorf("error al validar folio: %v", err)
		}
	}
	return nil
}

// generarArchivos genera PDF y XML de la factura
func generarArchivos(factura models.Factura, plantillaBytes []byte) (*bytes.Buffer, []byte, error) {
	var pdfBuffer *bytes.Buffer
	var err error

	// Generar PDF
	if len(plantillaBytes) > 0 {
		pdfBuffer, err = services.ProcesarPlantilla(factura, plantillaBytes)
	} else {
		pdfBuffer, err = services.GenerarPDF(factura, nil)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("error al generar PDF: %v", err)
	}

	// Generar XML
	xmlBytes, err := services.GenerarXML(factura)
	if err != nil {
		return nil, nil, fmt.Errorf("error al generar XML: %v", err)
	}

	return pdfBuffer, xmlBytes, nil
}

// guardarEnHistorial guarda la factura en el historial
func guardarEnHistorial(factura models.Factura) {
	if factura.IdUsuario > 0 {
		_, err := models.InsertarHistorialFactura(
			factura.IdUsuario,
			factura.ReceptorRFC,
			factura.ReceptorRazonSocial,
			factura.ClaveTicket,
			factura.NumeroFolio, // Incluir el folio generado
			factura.Total,
			factura.UsoCFDI,
			factura.Observaciones,
		)

		if err != nil {
			log.Printf("Error al guardar en historial (no crítico): %v", err)
		} else {
			log.Printf("Factura guardada en historial con folio: %s", factura.NumeroFolio)
		}
	}
}

func GenerarFacturaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Procesar datos del request
	factura, plantillaBytes, err := procesarDatosFactura(r)
	if err != nil {
		log.Printf("Error al procesar datos: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Manejar folio
	if err := manejarFolio(&factura); err != nil {
		log.Printf("Error con folio: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generar archivos
	pdfBuffer, xmlBytes, err := generarArchivos(factura, plantillaBytes)
	if err != nil {
		log.Printf("Error al generar archivos: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Crear ZIP
	zipBuffer, err := services.CrearZIP(pdfBuffer.Bytes(), xmlBytes)
	if err != nil {
		log.Printf("Error al crear ZIP: %v", err)
		http.Error(w, "Error al crear archivo ZIP", http.StatusInternalServerError)
		return
	}

	// Guardar en historial
	guardarEnHistorial(factura)

	// Enviar respuesta
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=factura_%s.zip", factura.NumeroFolio))
	_, err = w.Write(zipBuffer.Bytes())
	if err != nil {
		log.Printf("Error al enviar archivo ZIP: %v", err)
	}

	log.Printf("Factura generada exitosamente con folio: %s", factura.NumeroFolio)
}
