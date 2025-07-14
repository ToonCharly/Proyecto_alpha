package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"carlos/Facts/Backend/internal/db"
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

	// *** ANTES DE VERIFICAR CONCEPTOS - LOG COMPLETO ***
	log.Printf("🔍 DEBUG_CONCEPTOS - Verificando conceptos recibidos:")
	log.Printf("🔍 DEBUG_CONCEPTOS - len(factura.Conceptos) = %d", len(factura.Conceptos))
	log.Printf("🔍 DEBUG_CONCEPTOS - ClaveTicket = '%s'", factura.ClaveTicket)
	log.Printf("🔍 DEBUG_CONCEPTOS - factura.Conceptos = %+v", factura.Conceptos)

	// *** Si no hay conceptos, intentar obtenerlos desde la base de datos usando clave_ticket ***
	if len(factura.Conceptos) == 0 {
		log.Printf("⚠️ FLUJO REAL - No se recibieron conceptos en el JSON")
		log.Printf("⚠️ FLUJO REAL - ClaveTicket recibida: '%s'", factura.ClaveTicket)
		log.Printf("⚠️ FLUJO REAL - IdUsuario: %d", factura.IdUsuario)
		log.Printf("⚠️ FLUJO REAL - Total factura: %.2f", factura.Total)

		if factura.ClaveTicket != "" {
			log.Printf("🔍 FLUJO REAL - Intentando obtener conceptos desde BD para ticket: '%s'", factura.ClaveTicket)
			conceptosBD, err := obtenerConceptosDesdeVentas(factura.ClaveTicket)
			if err != nil {
				log.Printf("❌ FLUJO REAL - Error al obtener conceptos desde BD: %v", err)
				log.Printf("🔄 FLUJO REAL - Usando productos de ejemplo como fallback")
				factura.Conceptos = []models.Concepto{}
			} else {
				log.Printf("✅ Obtenidos %d conceptos desde BD usando clave_ticket", len(conceptosBD))
				for i, concepto := range conceptosBD {
					log.Printf("  [BD] Concepto %d: ClaveProdServ='%s', Desc='%s', Cant=%.2f, Precio=%.2f, Importe=%.2f", i+1, concepto.ClaveProdServ, concepto.Descripcion, concepto.Cantidad, concepto.ValorUnitario, concepto.Importe)
				}
				factura.Conceptos = conceptosBD
			}
		} else {
			log.Printf("⚠️ No hay clave_ticket disponible, usando productos de ejemplo")
			factura.Conceptos = []models.Concepto{}
		}
	}

	// *** LOG FINAL - VERIFICAR CONCEPTOS ANTES DE GENERAR PDF ***
	log.Printf("🔍 FINAL_DEBUG - Conceptos finales para PDF:")
	log.Printf("🔍 FINAL_DEBUG - Total conceptos: %d", len(factura.Conceptos))
	for i, concepto := range factura.Conceptos {
		log.Printf("🔍 FINAL_DEBUG - Concepto %d: ClaveProdServ='%s', Desc='%s', Cant=%.2f, Precio=%.2f, Importe=%.2f",
			i+1, concepto.ClaveProdServ, concepto.Descripcion, concepto.Cantidad, concepto.ValorUnitario, concepto.Importe)
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

	// Llenar datos del emisor automáticamente desde los datos fiscales del usuario
	if factura.IdUsuario > 0 {
		err := LlenarDatosEmisor(&factura, factura.IdUsuario)
		if err != nil {
			log.Printf("INFO - No se llenaron datos del emisor: %v", err)
			log.Printf("INFO - La factura se generará sin datos del emisor predefinidos")
			// No devolvemos error, continuamos con los datos disponibles
		} else {
			log.Printf("SUCCESS - Datos del emisor llenados correctamente")
		}
	} else {
		log.Printf("WARNING: No se encontró ID de usuario válido en la factura (IdUsuario=%d)", factura.IdUsuario)
	}

	// Convertir el ID del régimen fiscal al código del SAT
	if factura.RegimenFiscal != "" {
		codigo, err := ObtenerCodigoRegimenFiscal(factura.RegimenFiscal)
		if err != nil {
			log.Printf("Error al obtener código de régimen fiscal: %v", err)
			// No devolvemos error, usamos el valor original
		} else {
			factura.RegimenFiscal = codigo
		}
	}

	if factura.RegimenFiscalReceptor != "" {
		codigo, err := ObtenerCodigoRegimenFiscal(factura.RegimenFiscalReceptor)
		if err != nil {
			log.Printf("Error al obtener código de régimen fiscal receptor: %v", err)
			// No devolvemos error, usamos el valor original
		} else {
			factura.RegimenFiscalReceptor = codigo
		}
	}

	// Convertir el ID del estado al nombre del estado si viene como número Y no tenemos ya el nombre
	if factura.EstadoNombre == "" && factura.Estado != 0 {
		estadoStr := fmt.Sprintf("%d", factura.Estado)
		nombreEstado, err := ObtenerNombreEstado(estadoStr)
		if err != nil {
			log.Printf("Error al obtener nombre de estado: %v", err)
			// No devolvemos error, usamos el valor original
		} else if nombreEstado != "" {
			factura.EstadoNombre = nombreEstado
		}
	}

	// Cargar logo del usuario admin (ID=1) - se usa en todas las facturas
	logoBytes, err := services.CargarLogoPlantilla("1") // Siempre usar el logo del admin
	if err != nil {
		log.Printf("Error al cargar logo del admin: %v", err)
		// Continuar sin logo en caso de error
		logoBytes = nil
	}

	// Procesar la plantilla si se proporcionó
	var pdfBuffer *bytes.Buffer
	var err2 error

	if len(plantillaBytes) > 0 {
		pdfBuffer, err2 = services.ProcesarPlantilla(factura, plantillaBytes)
		if err2 != nil {
			log.Printf("Error al procesar plantilla: %v", err2)
			http.Error(w, "Error al procesar la plantilla", http.StatusInternalServerError)
			return
		}
	} else {
		// Obtener UUID y NoCertificado de la base de datos usando el folio generado
		uuid, noCert, err := db.ObtenerUUIDyNoCertificado(factura.NumeroFolio)
		if err != nil {
			log.Printf("No se pudo obtener UUID o NoCertificado para el folio %s: %v", factura.NumeroFolio, err)
			// Puedes decidir continuar, pero los campos estarán vacíos
		}
		factura.UUID = uuid
		factura.NoCertificado = noCert

		// Ahora sí, genera el PDF y los datos se mostrarán correctamente
		pdfBuffer, _, err2 = services.GenerarPDF(factura, nil, logoBytes)
		if err2 != nil {
			log.Printf("Error al generar PDF: %v", err2)
			http.Error(w, "Error al generar la factura", http.StatusInternalServerError)
			return
		}
	}

	// Generar XML para la factura
	xmlBytes, err := services.GenerarXML(factura)
	if err2 != nil {
		log.Printf("Error al generar XML: %v", err2)
		http.Error(w, "Error al generar XML de la factura", http.StatusInternalServerError)
		return
	}

	// Obtener los bytes del PDF y XML
	pdfBytes := pdfBuffer.Bytes()
	// xmlBytes ya es []byte

	// Obtener la serie y el folio
	serieDF := factura.Serie
	numeroFolio := factura.NumeroFolio
	nombrePDF := GenerarNombreArchivoFactura(serieDF, numeroFolio, "pdf")
	nombreXML := GenerarNombreArchivoFactura(serieDF, numeroFolio, "xml")

	// Usar la función CrearZIP para empaquetar PDF y XML con nombre correcto
	zipBuffer, err := services.CrearZIPConNombres(pdfBytes, xmlBytes, nombrePDF, nombreXML)
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
	if factura.IdUsuario > 0 { // Solo si tenemos un ID de usuario válido
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
	}

	// Guardar el archivo ZIP temporalmente (opcional, para descarga posterior)
	// Aquí podrías guardar el ZIP en un almacenamiento temporal si es necesario

	utils.RespondWithJSON(w, http.StatusOK, response)
	log.Printf("Factura generada exitosamente con folio: %s", numeroFolio)
}

// DebugVentasDetHandler - endpoint de debug para ver datos en ventas_det
func DebugVentasDetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	database := db.GetDB()

	query := `
		SELECT id, clave_producto, descripcion, clave_sat, unidad_sat, cantidad, precio_unitario, descuento, total, fecha_venta
		FROM ventas_det 
		ORDER BY id DESC 
		LIMIT 10
	`

	rows, err := database.Query(query)
	if err != nil {
		log.Printf("Error al consultar ventas_det: %v", err)
		http.Error(w, "Error al consultar base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ventas []map[string]interface{}
	for rows.Next() {
		var id int
		var claveProducto, descripcion, claveSat, unidadSat, fechaVenta string
		var cantidad, precioUnitario, descuento, total float64

		err := rows.Scan(&id, &claveProducto, &descripcion, &claveSat, &unidadSat, &cantidad, &precioUnitario, &descuento, &total, &fechaVenta)
		if err != nil {
			log.Printf("Error al escanear fila: %v", err)
			continue
		}

		venta := map[string]interface{}{
			"id":              id,
			"clave_producto":  claveProducto,
			"descripcion":     descripcion,
			"clave_sat":       claveSat,
			"unidad_sat":      unidadSat,
			"cantidad":        cantidad,
			"precio_unitario": precioUnitario,
			"descuento":       descuento,
			"total":           total,
			"fecha_venta":     fechaVenta,
		}
		ventas = append(ventas, venta)
	}

	response := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Encontradas %d ventas en la tabla ventas_det", len(ventas)),
		"ventas":  ventas,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}
