package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"Facts/internal/db"
	"Facts/internal/models"
	"Facts/internal/services"
)

func GenerarNombreArchivoFactura(serieDF, numeroFolio, extension string) string {
	re := regexp.MustCompile(`\d+$`)
	soloNumero := re.FindString(numeroFolio)
	// Quitar ceros a la izquierda:
	numeroSinCeros := soloNumero
	if n, err := strconv.Atoi(soloNumero); err == nil {
		numeroSinCeros = fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("Factura %s%s.%s", serieDF, numeroSinCeros, extension)
}

// obtenerConceptosDesdeVentas obtiene los productos desde la tabla ventas_det usando la serie
// NOTA: Actualmente la tabla ventas_det no tiene columna 'serie', usando como fallback los productos más recientes
func obtenerConceptosDesdeVentas(claveTicket string) ([]models.Concepto, error) {
	log.Printf("🔍 BD_DEBUG - Iniciando búsqueda de conceptos para ticket: '%s'", claveTicket)
	database := db.GetDB()

	// Primero intentar buscar por serie (cuando la columna exista)
	query := `
		SELECT 
			id,
			COALESCE(clave_producto, '') as clave_producto,
			descripcion,
			clave_sat,
			unidad_sat,
			cantidad,
			precio_unitario,
			descuento,
			total,
			COALESCE(iva, 16.0) as iva
		FROM ventas_det
		WHERE serie = ?
		ORDER BY id
	`

	log.Printf("🔍 BD_DEBUG - Ejecutando query con serie: '%s'", claveTicket)

	// Primero verificar si existen datos con esta serie
	var count int
	countQuery := "SELECT COUNT(*) FROM ventas_det WHERE serie = ?"
	err := database.QueryRow(countQuery, claveTicket).Scan(&count)
	if err != nil {
		log.Printf("❌ BD_DEBUG - Error al contar registros: %v", err)
	} else {
		log.Printf("🔍 BD_DEBUG - Registros encontrados en ventas_det con serie '%s': %d", claveTicket, count)
	}

	rows, err := database.Query(query, claveTicket)
	if err != nil {
		log.Printf("❌ BD_DEBUG - Error en query con serie: %v", err)
		log.Printf("🔄 BD_DEBUG - Intentando fallback: productos más recientes")

		// Fallback: obtener los productos más recientes (últimos 10 minutos)
		queryFallback := `
			SELECT 
				id,
				COALESCE(clave_producto, '') as clave_producto,
				descripcion,
				clave_sat,
				unidad_sat,
				cantidad,
				precio_unitario,
				descuento,
				total,
				COALESCE(iva, 16.0) as iva
			FROM ventas_det
			WHERE fecha_venta >= DATE_SUB(NOW(), INTERVAL 10 MINUTE)
			ORDER BY id DESC
		`

		rows, err = database.Query(queryFallback)
		if err != nil {
			log.Printf("❌ BD_DEBUG - Error en query fallback: %v", err)
			return nil, fmt.Errorf("error al consultar ventas_det: %v", err)
		}
	}
	defer rows.Close()

	var conceptos []models.Concepto
	// Usar mapa para deduplicar conceptos por descripción + clave SAT
	conceptosUnicos := make(map[string]models.Concepto)
	rowCount := 0

	for rows.Next() {
		rowCount++
		var concepto models.Concepto
		var id int
		var claveProducto, claveSat, unidadSat string
		var descuento, total, iva float64

		err := rows.Scan(
			&id,
			&claveProducto,
			&concepto.Descripcion,
			&claveSat,
			&unidadSat,
			&concepto.Cantidad,
			&concepto.ValorUnitario,
			&descuento,
			&total,
			&iva,
		)
		if err != nil {
			log.Printf("❌ BD_DEBUG - Error al escanear fila %d: %v", rowCount, err)
			continue
		}

		// Crear clave única usando descripción + cantidad + precio para mayor precisión
		claveUnica := fmt.Sprintf("%s_%.2f_%.2f", concepto.Descripcion, concepto.Cantidad, concepto.ValorUnitario)

		// Solo agregar si no existe este concepto único
		if _, existe := conceptosUnicos[claveUnica]; !existe {
			log.Printf("✅ BD_DEBUG - Fila %d: ID=%d, ClaveProd='%s', Desc='%s', ClaveSAT='%s', UnidadSAT='%s', Cant=%.2f, Precio=%.2f, IVA=%.2f",
				rowCount, id, claveProducto, concepto.Descripcion, claveSat, unidadSat, concepto.Cantidad, concepto.ValorUnitario, iva)

			// Mapear campos usando los datos reales de la tabla
			concepto.ClaveProdServ = claveProducto // Usar clave_producto como clave del producto/servicio
			concepto.ClaveSAT = claveSat           // Usar clave_sat como clave SAT
			concepto.ClaveUnidad = unidadSat       // Usar unidad_sat como clave de unidad
			concepto.Importe = total               // El total ya viene calculado
			concepto.Descuento = descuento         // Usar el descuento de la tabla
			concepto.TasaIVA = iva                 // Usar el IVA real de la base de datos
			concepto.TasaIEPS = 0.0                // Sin IEPS por defecto

			conceptosUnicos[claveUnica] = concepto
		} else {
			log.Printf("⚠️ BD_DEBUG - Concepto duplicado ignorado: Desc='%s', Cant=%.2f, Precio=%.2f",
				concepto.Descripcion, concepto.Cantidad, concepto.ValorUnitario)
		}
	}

	// Convertir mapa a slice
	for _, concepto := range conceptosUnicos {
		conceptos = append(conceptos, concepto)
	}

	log.Printf("🔍 BD_DEBUG - Total filas procesadas: %d", rowCount)
	log.Printf("🔍 BD_DEBUG - Conceptos únicos después de deduplicación: %d", len(conceptosUnicos))
	log.Printf("🔍 BD_DEBUG - Total conceptos finales: %d", len(conceptos))

	if len(conceptos) == 0 {
		log.Printf("❌ BD_DEBUG - No se encontraron productos para la serie '%s'", claveTicket)
		return nil, fmt.Errorf("no se encontraron productos para la serie %s", claveTicket)
	}

	log.Printf("✅ BD_DEBUG - Retornando %d conceptos para serie '%s'", len(conceptos), claveTicket)
	return conceptos, nil
}

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

// procesarMultipartForm procesa datos multipart/form-data con validación mejorada
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

	log.Printf("[DEBUG] Datos de factura recibidos: %s", facturaData)

	// Validar formato antes de hacer Unmarshal
	if err := validarFormatoJSON(facturaData); err != nil {
		return factura, nil, err
	}

	// Intentar decodificar con manejo de errores específico
	if err := json.Unmarshal([]byte(facturaData), &factura); err != nil {
		if strings.Contains(err.Error(), "invalid character '-'") {
			return factura, nil, fmt.Errorf("JSON en form-data contiene un valor numérico inválido con guión. Los campos con guiones deben ser strings. Error: %v", err)
		}
		return factura, nil, fmt.Errorf("error al decodificar JSON: %v", err)
	}

	// *** LÓGICA DE OBTENCIÓN DE CONCEPTOS DESDE BD ***
	log.Printf("🔍 DEBUG_CONCEPTOS - Verificando conceptos recibidos:")
	log.Printf("🔍 DEBUG_CONCEPTOS - len(factura.Conceptos) = %d", len(factura.Conceptos))
	log.Printf("🔍 DEBUG_CONCEPTOS - ClaveTicket = '%s'", factura.ClaveTicket)

	// Si no hay conceptos, intentar obtenerlos desde la base de datos usando clave_pedido
	if len(factura.Conceptos) == 0 {
		log.Printf("⚠️ FLUJO REAL - No se recibieron conceptos en el JSON")
		log.Printf("⚠️ FLUJO REAL - ClaveTicket recibida: '%s'", factura.ClaveTicket)

		if factura.ClaveTicket != "" {
			log.Printf("🔍 FLUJO REAL - Intentando obtener conceptos desde BD para ticket: '%s'", factura.ClaveTicket)
			conceptosBD, err := obtenerConceptosDesdeVentas(factura.ClaveTicket)
			if err != nil {
				log.Printf("❌ FLUJO REAL - Error al obtener conceptos desde BD: %v", err)
				log.Printf("🔄 FLUJO REAL - Continuando sin conceptos")
			} else {
				log.Printf("✅ Obtenidos %d conceptos desde BD usando clave_pedido", len(conceptosBD))
				for i, concepto := range conceptosBD {
					log.Printf("  [BD] Concepto %d: ClaveProdServ='%s', Desc='%s', Cant=%.2f, Precio=%.2f, Importe=%.2f",
						i+1, concepto.ClaveProdServ, concepto.Descripcion, concepto.Cantidad, concepto.ValorUnitario, concepto.Importe)
				}
				factura.Conceptos = conceptosBD
			}
		} else {
			log.Printf("⚠️ No hay clave_ticket disponible")
		}
	}

	log.Printf("🔍 FINAL_DEBUG - Total conceptos finales: %d", len(factura.Conceptos))

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

// procesarJSON procesa datos JSON con validación mejorada
func procesarJSON(r *http.Request) (models.Factura, []byte, error) {
	var factura models.Factura

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return factura, nil, fmt.Errorf("error al leer el cuerpo de la solicitud: %v", err)
	}

	// Log para depuración
	jsonStr := string(body)
	log.Printf("[DEBUG] JSON recibido: %s", jsonStr)

	// Validaciones mejoradas para detectar problemas con guiones
	if err := validarFormatoJSON(jsonStr); err != nil {
		// Hacer debugging adicional si hay error
		debugJSON(jsonStr)
		return factura, nil, err
	}

	// Intentar decodificar JSON con manejo de errores específico
	if err := json.Unmarshal(body, &factura); err != nil {
		// Si el error contiene información sobre guiones, hacer debugging y proporcionar mensaje específico
		if strings.Contains(err.Error(), "invalid character '-'") {
			log.Printf("[ERROR] Detectado error de guión en JSON:")
			debugJSON(jsonStr)

			// Opción temporal: intentar sanitizar el JSON
			log.Printf("[DEBUG] Intentando sanitizar JSON...")
			jsonSanitizado := sanitizarJSONProblematico(jsonStr)
			log.Printf("[DEBUG] JSON sanitizado: %s", jsonSanitizado)

			// Intentar de nuevo con JSON sanitizado
			if err2 := json.Unmarshal([]byte(jsonSanitizado), &factura); err2 == nil {
				log.Printf("[DEBUG] ✅ JSON sanitizado exitosamente")
				return factura, nil, nil
			}

			return factura, nil, fmt.Errorf("JSON contiene un valor numérico inválido con guión. Ejemplo: '001-002' debe ser string (con comillas): \"001-002\". Error original: %v", err)
		}
		return factura, nil, fmt.Errorf("error al decodificar JSON: %v", err)
	}

	return factura, nil, nil
}

// validarFormatoJSON realiza validaciones específicas del formato JSON
func validarFormatoJSON(jsonStr string) error {
	// Detectar números con guion en contexto JSON (más preciso)
	// Busca patrones como: "campo": 001-002 o "campo":001-002
	reNumGuion := regexp.MustCompile(`"[^"]*"\s*:\s*[0-9]+-[0-9]+`)
	if reNumGuion.MatchString(jsonStr) {
		matches := reNumGuion.FindAllString(jsonStr, -1)
		return fmt.Errorf("JSON contiene números con guión (no válidos): %v. Los campos con guiones deben ser strings (entre comillas). Ejemplo: '001-002' debe ser \"001-002\"", matches)
	}

	// Detectar valores numéricos con guion fuera de comillas
	reNumCeroGuion := regexp.MustCompile(`:\s*\d+-\d+`)
	if reNumCeroGuion.MatchString(jsonStr) {
		matches := reNumCeroGuion.FindAllString(jsonStr, -1)
		return fmt.Errorf("JSON contiene valores numéricos con guión: %v. Deben ser strings (entre comillas). Ejemplo: '001-002' debe ser \"001-002\"", matches)
	}

	// Detectar casos específicos como folios mal formateados
	reFolioMalFormato := regexp.MustCompile(`"(numero_folio|folio|serie)"\s*:\s*[0-9]+-[0-9]+`)
	if reFolioMalFormato.MatchString(jsonStr) {
		matches := reFolioMalFormato.FindAllString(jsonStr, -1)
		return fmt.Errorf("Campo de folio mal formateado: %v. Los folios deben ser strings (entre comillas). Ejemplo: \"numero_folio\": \"001-002\"", matches)
	}

	return nil
}

// Función auxiliar para debuggear el JSON problemático
func debugJSON(jsonStr string) {
	log.Printf("[DEBUG] Analizando JSON problemático:")
	log.Printf("[DEBUG] Longitud: %d caracteres", len(jsonStr))

	// Buscar patrones problemáticos
	reNumGuion := regexp.MustCompile(`"[^"]*"\s*:\s*[0-9]+-[0-9]+`)
	matches := reNumGuion.FindAllString(jsonStr, -1)
	if len(matches) > 0 {
		log.Printf("[DEBUG] Números con guión encontrados: %v", matches)
	}

	// Buscar líneas específicas con problemas
	lines := strings.Split(jsonStr, "\n")
	for i, line := range lines {
		if strings.Contains(line, "-") && (strings.Contains(line, ":") || strings.Contains(line, "folio")) {
			log.Printf("[DEBUG] Línea %d sospechosa: %s", i+1, strings.TrimSpace(line))
		}
	}
}

// Función para sanitizar JSON problemático (uso temporal para debugging)
func sanitizarJSONProblematico(jsonStr string) string {
	// Convertir números con guión a strings
	re := regexp.MustCompile(`("numero_folio"|"folio"|"serie")\s*:\s*([0-9]+-[0-9]+)`)
	jsonStr = re.ReplaceAllString(jsonStr, `$1: "$2"`)

	// Patrón más general para otros campos numéricos con guión
	re2 := regexp.MustCompile(`("\w+")\s*:\s*([0-9]+-[0-9]+)`)
	jsonStr = re2.ReplaceAllString(jsonStr, `$1: "$2"`)

	return jsonStr
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

	// Cargar logo del usuario admin (ID=1) - se usa en todas las facturas
	logoBytes, err := services.CargarLogoPlantilla("1") // Siempre usar el logo del admin
	if err != nil {
		log.Printf("Error al cargar logo del admin: %v", err)
		// Continuar sin logo en caso de error
		logoBytes = nil
	}

	// Generar PDF
	if len(plantillaBytes) > 0 {
		pdfBuffer, err = services.ProcesarPlantilla(factura, plantillaBytes)
	} else {
		// Usar el logo de plantillas cargado
		pdfBuffer, _, err = services.GenerarPDF(factura, nil, logoBytes)
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
			factura.NumeroFolio,
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

	// Usar procesarDatosFactura para unificar validación y decodificación
	factura, plantillaBytes, err := procesarDatosFactura(r)
	if err != nil {
		log.Printf("Error al decodificar solicitud: %v", err)
		http.Error(w, "Error al procesar los datos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// --- Mapear datos de empresa si viene EmpresaID, IdEmpresa o EmpresaRFC ---
	// Si tienes el ID (o RFC) de la empresa en la factura recibida
	if factura.EmpresaID != nil {
		// Si viene como float64 del JSON
		if empresaID, ok := factura.EmpresaID.(float64); ok {
			empresa, err := models.ObtenerEmpresaPorID(int(empresaID))
			if err == nil && empresa != nil {
				factura.EmpresaRFC = empresa.RFC
				factura.RazonSocial = empresa.RazonSocial
				factura.Direccion = empresa.Direccion
				factura.CodigoPostal = empresa.CodigoPostal
				factura.RegimenFiscal = empresa.RegimenFiscal
				// Otros campos que quieras mapear
			}
		}
	} else if factura.IdEmpresa > 0 {
		empresa, err := models.ObtenerEmpresaPorID(factura.IdEmpresa)
		if err == nil && empresa != nil {
			factura.EmpresaRFC = empresa.RFC
			factura.RazonSocial = empresa.RazonSocial
			factura.Direccion = empresa.Direccion
			factura.CodigoPostal = empresa.CodigoPostal
			factura.RegimenFiscal = empresa.RegimenFiscal
		}
	}

	if len(factura.Conceptos) == 0 && factura.ClaveTicket != "" {
		conceptosBD, err := obtenerConceptosDesdeVentas(factura.ClaveTicket)
		if err == nil {
			factura.Conceptos = conceptosBD
		}
	}

	if factura.NumeroFolio == "" {
		err := factura.GenerarFolioAutomatico()
		if err != nil {
			log.Printf("Error al generar folio automático: %v", err)
			http.Error(w, "Error al generar folio de factura", http.StatusInternalServerError)
			return
		}
	} else {
		err := factura.ValidarFolio()
		if err != nil {
			log.Printf("Error al validar folio: %v", err)
			http.Error(w, "Folio duplicado o inválido: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	if factura.IdUsuario > 0 {
		err := LlenarDatosEmisor(&factura, factura.IdUsuario)
		if err != nil {
			log.Printf("INFO - No se llenaron datos del emisor: %v", err)
		}
		// Log de depuración para KeyPath y ClaveCSD
		log.Printf("DEBUG - KeyPath: %s, ClaveCSD: %s", factura.KeyPath, factura.ClaveCSD)
	}
	if factura.RegimenFiscal != "" {
		codigo, err := ObtenerCodigoRegimenFiscal(factura.RegimenFiscal)
		if err == nil {
			factura.RegimenFiscal = codigo
		}
	}
	if factura.RegimenFiscalReceptor != "" {
		codigo, err := ObtenerCodigoRegimenFiscal(factura.RegimenFiscalReceptor)
		if err == nil {
			factura.RegimenFiscalReceptor = codigo
		}
	}
	if factura.EstadoNombre == "" && factura.Estado != 0 {
		estadoStr := fmt.Sprintf("%d", factura.Estado)
		nombreEstado, err := ObtenerNombreEstado(estadoStr)
		if err == nil && nombreEstado != "" {
			factura.EstadoNombre = nombreEstado
		}
	}

	// Cargar logo del usuario admin (ID=1)
	logoBytes, err := services.CargarLogoPlantilla("1")
	if err != nil {
		log.Printf("Error al cargar logo del admin: %v", err)
		logoBytes = nil
	}

	var pdfBuffer *bytes.Buffer
	if len(plantillaBytes) > 0 {
		pdfBuffer, err = services.ProcesarPlantilla(factura, plantillaBytes)
		if err != nil {
			log.Printf("Error al procesar plantilla: %v", err)
			http.Error(w, "Error al procesar la plantilla", http.StatusInternalServerError)
			return
		}
	} else {
		pdfBuffer, _, err = services.GenerarPDF(factura, nil, logoBytes)
		if err != nil {
			log.Printf("Error al generar PDF: %v", err)
			http.Error(w, "Error al generar la factura", http.StatusInternalServerError)
			return
		}
	}

	keyPath := factura.KeyPath
	claveCSD := factura.ClaveCSD
	if keyPath == "" || claveCSD == "" {
		log.Printf("Error: No se proporcionó la ruta al archivo .key o la clave CSD")
		http.Error(w, "Faltan datos para la firma digital (archivo .key o clave CSD)", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(keyPath); err != nil {
		log.Printf("Error: El archivo .key no existe en la ruta proporcionada: %s", keyPath)
		http.Error(w, "El archivo .key no existe en la ruta proporcionada: "+keyPath, http.StatusBadRequest)
		return
	}
	xmlBytes, err := services.ProcesarKeyYGenerarCFDI(factura, keyPath, claveCSD, "")
	if err != nil {
		log.Printf("Error al generar XML firmado CFDI: %v", err)
		http.Error(w, "Error al generar XML firmado CFDI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	serieDF := factura.Serie
	numeroFolio := factura.NumeroFolio
	nombrePDF := GenerarNombreArchivoFactura(serieDF, numeroFolio, "pdf")
	nombreXML := GenerarNombreArchivoFactura(serieDF, numeroFolio, "xml")
	pdfBytes := pdfBuffer.Bytes()

	zipBuffer, err := services.CrearZIPConNombres(pdfBytes, xmlBytes, nombrePDF, nombreXML)
	if err != nil {
		log.Printf("Error al crear ZIP: %v", err)
		http.Error(w, "Error al crear archivo ZIP", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=factura_%s.zip", numeroFolio))
	_, err = w.Write(zipBuffer.Bytes())
	if err != nil {
		log.Printf("Error al enviar archivo ZIP: %v", err)
	}
	log.Printf("Factura generada exitosamente con folio: %s", numeroFolio)
}
