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
	"carlos/Facts/Backend/internal/utils"
)

// GenerarFacturaConInfoHandler genera una factura y devuelve información sobre ella
func GenerarFacturaConInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var factura models.Factura
	var plantillaBytes []byte

	// Detectar el tipo de contenido
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		http.Error(w, "Content-Type no especificado", http.StatusBadRequest)
		return
	}

	if strings.Contains(contentType, "multipart/form-data") {
		// Procesar como multipart/form-data
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Printf("Error al parsear multipart form: %v", err)
			http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
			return
		}

		// Leer el campo "datos" que contiene el JSON de la factura
		facturaData := r.FormValue("datos")
		if facturaData == "" {
			http.Error(w, "No se encontraron datos de factura", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal([]byte(facturaData), &factura); err != nil {
			log.Printf("Error al decodificar JSON en multipart: %v", err)
			http.Error(w, "Error al procesar los datos: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Leer el archivo de la plantilla
		plantillaFile, _, err := r.FormFile("plantilla")
		if err == nil {
			defer plantillaFile.Close()
			plantillaBytes, err = io.ReadAll(plantillaFile)
			if err != nil {
				log.Printf("Error al leer la plantilla: %v", err)
				http.Error(w, "Error al leer la plantilla", http.StatusInternalServerError)
				return
			}
		}
	} else if strings.Contains(contentType, "application/json") {
		// Procesar como JSON simple
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
			http.Error(w, "Error al leer la solicitud", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &factura); err != nil {
			log.Printf("Error al decodificar JSON: %v", err)
			http.Error(w, "Error al procesar los datos: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Generar folio automáticamente si no se proporcionó uno
	if factura.NumeroFolio == "" {
		err := factura.GenerarFolioAutomatico()
		if err != nil {
			log.Printf("Error al generar folio automático: %v", err)
			http.Error(w, "Error al generar folio de factura", http.StatusInternalServerError)
			return
		}
		log.Printf("Folio generado automáticamente: %s", factura.NumeroFolio)
	} else {
		// Si se proporcionó un folio, validar que sea único
		err := factura.ValidarFolio()
		if err != nil {
			log.Printf("Error al validar folio: %v", err)
			http.Error(w, "Folio duplicado o inválido: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Procesar la plantilla si se proporcionó
	var pdfBuffer *bytes.Buffer
	var err error
	if len(plantillaBytes) > 0 {
		pdfBuffer, err = services.ProcesarPlantilla(factura, plantillaBytes)
		if err != nil {
			log.Printf("Error al procesar plantilla: %v", err)
			http.Error(w, "Error al procesar la plantilla", http.StatusInternalServerError)
			return
		}
	} else {
		// Generar el PDF sin plantilla
		pdfBuffer, err = services.GenerarPDF(factura, nil)
		if err != nil {
			log.Printf("Error al generar PDF: %v", err)
			http.Error(w, "Error al generar la factura", http.StatusInternalServerError)
			return
		}
	}

	// Generar XML para la factura
	xmlBytes, err := services.GenerarXML(factura)
	if err != nil {
		log.Printf("Error al generar XML: %v", err)
		http.Error(w, "Error al generar XML de la factura", http.StatusInternalServerError)
		return
	}

	// Usar la función CrearZIP para empaquetar PDF y XML
	zipBuffer, err := services.CrearZIP(pdfBuffer.Bytes(), xmlBytes)
	if err != nil {
		log.Printf("Error al crear ZIP: %v", err)
		http.Error(w, "Error al crear archivo ZIP", http.StatusInternalServerError)
		return
	}

	// Devolver información sobre la factura generada
	response := map[string]interface{}{
		"success":       true,
		"message":       "Factura generada exitosamente",
		"numero_folio":  factura.NumeroFolio,
		"rfc_emisor":    factura.RFC,
		"rfc_receptor":  factura.ReceptorRFC,
		"total":         factura.Total,
		"fecha_emision": factura.FechaEmision,
		"uso_cfdi":      factura.UsoCFDI,
		"archivo_size":  zipBuffer.Len(),
		"download_url":  fmt.Sprintf("/api/descargar-factura-zip?folio=%s", factura.NumeroFolio),
	}

	// Guardar automáticamente en el historial de facturas (DESPUÉS de generar todo)
	log.Printf("DEBUG - Intentando guardar en historial:")
	log.Printf("  IdEmpresa: %d", factura.IdEmpresa)
	log.Printf("  ReceptorRFC: '%s'", factura.ReceptorRFC)
	log.Printf("  ReceptorRazonSocial: '%s'", factura.ReceptorRazonSocial)
	log.Printf("  ClaveTicket: '%s'", factura.ClaveTicket)
	log.Printf("  NumeroFolio: '%s'", factura.NumeroFolio)
	log.Printf("  Total: %f", factura.Total)
	log.Printf("  UsoCFDI: '%s'", factura.UsoCFDI)
	log.Printf("  Observaciones: '%s'", factura.Observaciones)

	if factura.IdUsuario > 0 { // Solo si tenemos un ID de usuario válido
		log.Printf("DEBUG - Llamando a InsertarHistorialFactura...")
		_, err = models.InsertarHistorialFactura(
			factura.IdUsuario, // Usar el ID del usuario que genera la factura
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
			// No devolvemos error aquí porque la factura se generó correctamente
		} else {
			log.Printf("Factura guardada en historial con folio: %s", factura.NumeroFolio)
		}
	} else {
		log.Printf("DEBUG - No se guarda en historial porque IdEmpresa <= 0: %d", factura.IdEmpresa)
	}

	// Guardar el archivo ZIP temporalmente (opcional, para descarga posterior)
	// Aquí podrías guardar el ZIP en un almacenamiento temporal si es necesario

	utils.RespondWithJSON(w, http.StatusOK, response)
	log.Printf("Factura generada exitosamente con folio: %s", factura.NumeroFolio)
}
