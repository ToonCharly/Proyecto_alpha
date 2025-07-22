package services

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"Facts/internal/db"
	"Facts/internal/models"

	"github.com/phpdave11/gofpdf"
)

// Función auxiliar para dividir texto de manera segura manteniendo caracteres UTF-8
func splitTextSafely(pdf *gofpdf.Fpdf, texto string, ancho float64) []string {
	// Dividir texto manualmente para evitar problemas con SplitText
	if len(texto) <= 80 {
		return []string{texto}
	}

	var lineas []string
	palabras := strings.Fields(texto)
	lineaActual := ""
	longitudMaxima := 80

	for _, palabra := range palabras {
		testLinea := lineaActual
		if testLinea != "" {
			testLinea += " "
		}
		testLinea += palabra

		// Si la línea se vuelve muy larga, guardar la anterior y empezar nueva
		if len(testLinea) > longitudMaxima {
			if lineaActual != "" {
				lineas = append(lineas, lineaActual)
			}
			lineaActual = palabra
		} else {
			lineaActual = testLinea
		}
	}

	if lineaActual != "" {
		lineas = append(lineas, lineaActual)
	}

	return lineas
}

// Función auxiliar para obtener la descripción del uso de CFDI
func obtenerDescripcionUsoCfdi(clave string) string {
	descripciones := map[string]string{
		"G01": "Adquisición de mercancías",
		"G02": "Devoluciones, descuentos o bonificaciones",
		"G03": "Gastos en general",
		"I01": "Construcciones",
		"D01": "Honorarios médicos",
		"P01": "Por definir",
	}
	if desc, ok := descripciones[clave]; ok {
		return desc
	}
	return clave
}

// Función auxiliar para obtener la descripción del régimen fiscal
func obtenerDescripcionRegimenFiscalPDF(clave string) string {
	if desc, ok := regimenFiscalDescripciones[clave]; ok {
		return desc
	}
	return clave
}

// Modificada para incluir nombre de archivo con serie_df y folio sin ceros a la izquierda
func GenerarPDF(factura models.Factura, empresa *models.Empresa, logoBytes []byte) (*bytes.Buffer, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAuthor("Sistema de Facturación", true)
	pdf.SetTitle("Factura Electrónica", true)
	pdf.AddPage()
	pdf.SetMargins(10, 10, 10)

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// ========== HEADER LIMPIO ==========
	// Agregar logo si está disponible (posición izquierda)
	if len(logoBytes) > 0 {
		tmpDir := os.TempDir()
		tmpfileName := filepath.Join(tmpDir, fmt.Sprintf("logo-%d.png", time.Now().UnixNano()))
		err := os.WriteFile(tmpfileName, logoBytes, 0644)
		if err == nil {
			defer os.Remove(tmpfileName)
			pdf.Image(tmpfileName, 10, 10, 30, 30, false, "", 0, "")
		}
	}

	// === OBTENER SERIE_DF DEL USUARIO (idUsuario de la factura) ===
	serieDF := ""
	idUsuario := 1
	if factura.EmpresaID != nil {
		if id, ok := factura.EmpresaID.(int); ok {
			idUsuario = id
		}
	}
	datosFiscales, err := db.ObtenerDatosFiscales(idUsuario)
	if err == nil {
		if s, ok := datosFiscales["serie_df"].(string); ok {
			serieDF = s
		}
	}

	// === OBTENER FOLIO COMO INT (sin ceros a la izquierda) ===
	folioInt := 0
	if factura.NumeroFolio != "" {
		folioInt, _ = strconv.Atoi(factura.NumeroFolio)
	}

	// === GENERAR NOMBRE DE ARCHIVO ===
	nombreArchivo := fmt.Sprintf("Factura_%s%d.pdf", serieDF, folioInt)

	y := 45.0
	pdf.SetTextColor(0, 0, 0)

	// ========== DATOS ESENCIALES DE LA FACTURA (COLUMNA IZQUIERDA) ==========
	pdf.SetFont("Arial", "B", 8) // Negrita para las etiquetas

	// RFC emisor
	if factura.EmisorRFC != "" {
		pdf.SetXY(10, y)
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(40, 4, tr("RFC emisor:"))
		pdf.SetFont("Arial", "", 8) // Normal para los datos
		pdf.SetTextColor(64, 64, 64)
		pdf.Cell(70, 4, tr(factura.EmisorRFC))
		y += 4
	}

	// Nombre emisor
	if factura.EmisorRazonSocial != "" {
		pdf.SetXY(10, y)
		pdf.SetFont("Arial", "B", 8) // Negrita para la etiqueta
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(40, 4, tr("Nombre emisor:"))
		pdf.SetFont("Arial", "", 8) // Normal para los datos
		pdf.SetTextColor(64, 64, 64)
		lineas := splitTextSafely(pdf, tr(factura.EmisorRazonSocial), 70)
		pdf.Cell(70, 4, tr(lineas[0]))
		y += 4
		for i := 1; i < len(lineas); i++ {
			pdf.SetXY(50, y)
			pdf.SetFont("Arial", "", 8) // Normal para los datos
			pdf.SetTextColor(64, 64, 64)
			pdf.Cell(70, 4, tr(lineas[i]))
			y += 4
		}
	}

	// Folio
	folio := ""
	if factura.NumeroFolio != "" {
		folio = factura.NumeroFolio
	} else if factura.Serie != "" {
		folio = factura.Serie
	}
	if folio != "" {
		pdf.SetXY(10, y)
		pdf.SetFont("Arial", "B", 8) // Negrita para la etiqueta
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(40, 4, tr("Folio:"))
		pdf.SetFont("Arial", "", 8) // Normal para los datos
		pdf.SetTextColor(64, 64, 64)
		pdf.Cell(70, 4, tr(folio))
		y += 4
	}

	// RFC receptor
	receptorRFC := ""
	if factura.ReceptorRFC != "" {
		receptorRFC = factura.ReceptorRFC
	} else if factura.ClienteRFC != "" {
		receptorRFC = factura.ClienteRFC
	}
	if receptorRFC != "" {
		pdf.SetXY(10, y)
		pdf.SetFont("Arial", "B", 8) // Negrita para la etiqueta
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(40, 4, tr("RFC receptor:"))
		pdf.SetFont("Arial", "", 8) // Normal para los datos
		pdf.SetTextColor(64, 64, 64)
		pdf.Cell(70, 4, tr(receptorRFC))
		y += 4
	}

	// Nombre receptor
	receptorRazon := ""
	if factura.ReceptorRazonSocial != "" {
		receptorRazon = factura.ReceptorRazonSocial
	} else if factura.ClienteRazonSocial != "" {
		receptorRazon = factura.ClienteRazonSocial
	}
	if receptorRazon != "" {
		pdf.SetXY(10, y)
		pdf.SetFont("Arial", "B", 8) // Negrita para la etiqueta
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(40, 4, tr("Nombre receptor:"))
		pdf.SetFont("Arial", "", 8) // Normal para los datos
		pdf.SetTextColor(64, 64, 64)
		lineas := splitTextSafely(pdf, tr(receptorRazon), 70)
		pdf.Cell(70, 4, tr(lineas[0]))
		y += 4
		for i := 1; i < len(lineas); i++ {
			pdf.SetXY(50, y)
			pdf.SetFont("Arial", "", 8) // Normal para los datos
			pdf.SetTextColor(64, 64, 64)
			pdf.Cell(70, 4, tr(lineas[i]))
			y += 4
		}
	}

	// Código postal del receptor (etiqueta en una sola línea, con acento)
	cpReceptor := ""
	if factura.CodigoPostal != "" {
		cpReceptor = factura.CodigoPostal
	} else if factura.ReceptorCodigoPostal != "" {
		cpReceptor = factura.ReceptorCodigoPostal
	}
	if cpReceptor != "" {
		pdf.SetXY(10, y)
		pdf.SetFont("Arial", "B", 8)
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(40, 4, tr("Código postal del receptor:"))
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(64, 64, 64)
		pdf.SetXY(50, y)
		pdf.Cell(70, 4, cpReceptor)
		y += 4
	}

	// Régimen fiscal receptor
	regimenReceptor := ""
	if factura.RegimenFiscalReceptor != "" {
		regimenReceptor = factura.RegimenFiscalReceptor
	} else if factura.RegimenFiscal != "" {
		regimenReceptor = factura.RegimenFiscal
	}
	if regimenReceptor != "" {
		pdf.SetXY(10, y)
		pdf.SetFont("Arial", "B", 8) // Negrita para la etiqueta
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(40, 4, tr("Régimen fiscal"))
		y += 4
		pdf.SetXY(10, y)
		pdf.SetFont("Arial", "B", 8) // Negrita para la etiqueta
		pdf.Cell(40, 4, tr("receptor:"))
		pdf.SetFont("Arial", "", 8) // Normal para los datos
		pdf.SetTextColor(64, 64, 64)
		regimenDesc := obtenerDescripcionRegimenFiscalPDF(regimenReceptor)
		pdf.Cell(70, 4, tr(regimenDesc))
		y += 4
	}

	// Uso CFDI
	pdf.SetXY(10, y)
	pdf.SetFont("Arial", "B", 8) // Negrita para la etiqueta
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(40, 4, tr("Uso CFDI:"))
	pdf.SetFont("Arial", "", 8) // Normal para los datos
	pdf.SetTextColor(64, 64, 64)
	var usoCfdiTexto string
	if factura.UsoCFDI != "" {
		usoCfdiDesc := obtenerDescripcionUsoCfdi(factura.UsoCFDI)
		usoCfdiTexto = usoCfdiDesc
	} else {
		usoCfdiTexto = "Gastos en general"
	}
	pdf.Cell(70, 4, tr(usoCfdiTexto))
	y += 4

	// ========== INFORMACIÓN FISCAL (COLUMNA DERECHA) ==========
	yFiscal := 45.0

	// Folio fiscal
	// Folio fiscal (UUID del timbre)
	pdf.SetXY(110, yFiscal)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(50, 4, tr("Folio fiscal:"))
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(64, 64, 64)
	pdf.Cell(80, 4, tr(factura.UUID))
	yFiscal += 4

	// Etiqueta y valor de No. de serie del CSD en varias líneas si es necesario
	leftX := 110.0
	labelWidth := 50.0
	lineHeight := 4.0

	noSerie := factura.NoCertificado
	splitByLength := func(text string, length int) []string {
		var out []string
		for len(text) > length {
			out = append(out, text[:length])
			text = text[length:]
		}
		if len(text) > 0 {
			out = append(out, text)
		}
		return out
	}
	lines := splitByLength(noSerie, 22)
	if len(lines) > 0 {
		// Imprime la etiqueta solo en la primera línea
		pdf.SetXY(leftX, yFiscal)
		pdf.SetFont("Arial", "B", 8)
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(labelWidth, lineHeight, "No. de serie del CSD:")

		// Imprime el valor, que puede ocupar varias líneas
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(64, 64, 64)
		pdf.SetXY(leftX+labelWidth, yFiscal)
		pdf.Cell(80, lineHeight, lines[0])

		for i := 1; i < len(lines); i++ {
			pdf.SetXY(leftX+labelWidth, yFiscal+float64(i)*lineHeight)
			pdf.Cell(80, lineHeight, lines[i])
		}
		// Ahora, actualiza yFiscal para dejar espacio después de todas las líneas
		yFiscal += float64(len(lines)) * lineHeight
	} else {
		log.Printf("PDF_WARNING - No se pudo mostrar el número de serie del CSD: valor vacío")
		yFiscal += lineHeight
	}

	// === DATOS DE TIMBRADO FISCAL DIGITAL ===
	if factura.Timbre != nil {
		pdf.SetXY(110, yFiscal)
		pdf.SetFont("Arial", "B", 8)
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(50, 4, tr("Datos de timbrado fiscal"))
		yFiscal += 4

		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(64, 64, 64)

		pdf.SetXY(110, yFiscal)
		pdf.Cell(50, 4, "UUID:")
		pdf.Cell(80, 4, factura.Timbre.UUID)
		yFiscal += 4

		pdf.SetXY(110, yFiscal)
		pdf.Cell(50, 4, "Fecha timbrado:")
		pdf.Cell(80, 4, factura.Timbre.FechaTimbrado)
		yFiscal += 4

		pdf.SetXY(110, yFiscal)
		pdf.Cell(50, 4, "RFC Proveedor Certif.:")
		pdf.Cell(80, 4, factura.Timbre.RfcProvCertif)
		yFiscal += 4

		pdf.SetXY(110, yFiscal)
		pdf.Cell(50, 4, "No. Certificado SAT:")
		pdf.Cell(80, 4, factura.Timbre.NoCertificadoSAT)
		yFiscal += 4

		pdf.SetXY(110, yFiscal)
		pdf.Cell(50, 4, "Sello CFD:")
		pdf.SetXY(160, yFiscal)
		pdf.MultiCell(130, 4, factura.Timbre.SelloCFD, "", "", false)
		yFiscal += 8

		pdf.SetXY(110, yFiscal)
		pdf.Cell(50, 4, "Sello SAT:")
		pdf.SetXY(160, yFiscal)
		pdf.MultiCell(130, 4, factura.Timbre.SelloSAT, "", "", false)
		yFiscal += 8
	}

	// ---- Etiqueta en dos líneas y formato con 'T' ----
	pdf.SetXY(110, yFiscal)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(50, 4, tr("Código postal, fecha y hora"))
	yFiscal += 4
	pdf.SetXY(110, yFiscal)
	pdf.Cell(50, 4, tr("de emisión:"))

	// Valor en la segunda columna
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(64, 64, 64)
	pdf.SetXY(160, yFiscal)
	codigoPostalEmisor := ""
	if factura.EmisorCodigoPostal != "" {
		codigoPostalEmisor = factura.EmisorCodigoPostal
	}
	fechaEmision := ""
	if factura.FechaEmision != "" {
		t, err := time.Parse(time.RFC3339, factura.FechaEmision)
		if err == nil {
			fechaEmision = t.Format("02/01/2006T15:04:05") // <-- Separador T aquí
		} else {
			fechaEmision = factura.FechaEmision
		}
	} else {
		fechaEmision = time.Now().Format("02/01/2006T15:04:05") // <-- Separador T aquí
	}
	var datosEmision string
	if codigoPostalEmisor != "" && fechaEmision != "" {
		datosEmision = fmt.Sprintf("CP: %s | %s", codigoPostalEmisor, fechaEmision)
	} else if codigoPostalEmisor != "" {
		datosEmision = fmt.Sprintf("CP: %s", codigoPostalEmisor)
	} else if fechaEmision != "" {
		datosEmision = fechaEmision
	} else {
		datosEmision = "-"
	}
	pdf.Cell(80, 4, tr(datosEmision))
	yFiscal += 4

	// Efecto de comprobante
	pdf.SetXY(110, yFiscal)
	pdf.SetFont("Arial", "B", 8) // Negrita para la etiqueta
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(50, 4, tr("Efecto de comprobante:"))
	pdf.SetFont("Arial", "", 8) // Normal para los datos
	pdf.SetTextColor(64, 64, 64)
	pdf.Cell(80, 4, tr("Ingreso"))
	yFiscal += 4

	// Régimen fiscal (del emisor)
	pdf.SetXY(110, yFiscal)
	pdf.SetFont("Arial", "B", 8) // Negrita para la etiqueta
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(50, 4, tr("Régimen fiscal:"))
	pdf.SetFont("Arial", "", 8) // Normal para los datos
	pdf.SetTextColor(64, 64, 64)
	regimenFiscalTexto := ""
	if factura.EmisorRegimenFiscal != "" {
		regimenFiscalTexto = obtenerDescripcionRegimenFiscalPDF(factura.EmisorRegimenFiscal)
	} else {
		regimenFiscalTexto = "General de Ley Personas Morales"
	}
	pdf.Cell(80, 4, tr(regimenFiscalTexto))
	yFiscal += 4

	// Ajustar Y para continuar después de ambas columnas
	if yFiscal > y {
		y = yFiscal
	}

	y += 2 // Espacio adicional reducido

	// Sin línea divisoria antes de la tabla

	// ========== TABLA DE PRODUCTOS MEJORADA ==========
	log.Printf("PDF_DEBUG - Generando tabla de conceptos")
	log.Printf("PDF_DEBUG - Total conceptos en factura: %d", len(factura.Conceptos))

	// Deduplicar conceptos
	conceptosUnicos := make(map[string]models.Concepto)
	for _, concepto := range factura.Conceptos {
		claveUnica := fmt.Sprintf("%s_%.2f_%.2f", concepto.Descripcion, concepto.Cantidad, concepto.ValorUnitario)
		if _, existe := conceptosUnicos[claveUnica]; !existe {
			conceptosUnicos[claveUnica] = concepto
		}
	}

	conceptosParaPDF := make([]models.Concepto, 0, len(conceptosUnicos))
	for _, concepto := range conceptosUnicos {
		conceptosParaPDF = append(conceptosParaPDF, concepto)
	}

	// Verificar espacio para la tabla - filas con altura fija
	spaceNeeded := float64(len(conceptosParaPDF)*8) + 20 // Altura fija por fila
	spaceAvailable := 280 - y
	if spaceNeeded > spaceAvailable {
		pdf.AddPage()
		y = 20
	}

	// Título de la tabla
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(30, 80, 150)
	pdf.SetXY(15, y)
	pdf.Cell(180, 8, tr("DETALLE DE PRODUCTOS"))
	y += 12

	// NUEVOS anchos de columna, quitando IEPS y el viejo IVA (%)
	colWidths := []float64{
		18, // Clave Prod/Ser
		62, // Producto
		13, // Clave SAT
		16, // Unidad SAT
		15, // Cantidad
		15, // Precio
		16, // IVA ($)
		20, // TOTAL
	}

	headers := []string{
		"Clave Prod/Ser", "Producto", "Clave SAT", "Unidad SAT",
		"Cantidad", "Precio", "IVA ($)", "TOTAL",
	}

	// Encabezados
	pdf.SetFont("Arial", "B", 7)
	pdf.SetFillColor(60, 120, 180)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetDrawColor(0, 0, 0)
	x := 15.0
	for i, header := range headers {
		pdf.SetXY(x, y)
		pdf.CellFormat(colWidths[i], 8, tr(header), "1", 0, "C", true, 0, "")
		x += colWidths[i]
	}
	y += 8

	// Configurar para el contenido de la tabla
	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(40, 40, 40)
	pdf.SetFillColor(248, 248, 248)
	pdf.SetDrawColor(0, 0, 0)

	// Dibujar filas de la tabla
	for i, concepto := range conceptosParaPDF {
		fillColor := i%2 == 0

		// Verificar si necesitamos una nueva página
		if y > 265 {
			pdf.AddPage()
			y = 20
			// Redibujar encabezados en la nueva página
			pdf.SetFont("Arial", "B", 7)
			pdf.SetFillColor(60, 120, 180)
			pdf.SetTextColor(255, 255, 255)
			pdf.SetDrawColor(0, 0, 0)
			x = 15.0
			for j, header := range headers {
				pdf.SetXY(x, y)
				pdf.CellFormat(colWidths[j], 8, tr(header), "1", 0, "C", true, 0, "")
				x += colWidths[j]
			}
			y += 8
			pdf.SetFont("Arial", "", 7)
			pdf.SetTextColor(40, 40, 40)
			pdf.SetFillColor(248, 248, 248)
			pdf.SetDrawColor(0, 0, 0)
		}

		// Calcular altura de la fila basada en el texto más largo
		maxLines := 1
		descripcionLines := splitTextSafely(pdf, tr(concepto.Descripcion), colWidths[1]-4)
		if len(descripcionLines) > maxLines {
			maxLines = len(descripcionLines)
		}

		rowHeight := float64(maxLines) * 3.5
		if rowHeight < 10 {
			rowHeight = 10
		}

		x = 15.0
		pdf.SetXY(x, y)
		pdf.CellFormat(colWidths[0], rowHeight, concepto.ClaveProdServ, "1", 0, "C", fillColor, 0, "")
		x += colWidths[0]

		// Producto (descripción multilínea igual que antes)
		pdf.SetXY(x, y)
		pdf.CellFormat(colWidths[1], rowHeight, "", "1", 0, "L", fillColor, 0, "")
		for j, line := range splitTextSafely(pdf, tr(concepto.Descripcion), colWidths[1]-4) {
			lineY := y + float64(j)*3.5 + 2.5
			pdf.SetXY(x+2, lineY)
			pdf.Cell(colWidths[1]-4, 3.5, line)
		}
		x += colWidths[1]

		claveSAT := concepto.ClaveSAT
		if claveSAT == "" || claveSAT == "0" {
			claveSAT = "N/A"
		}
		pdf.SetXY(x, y)
		pdf.CellFormat(colWidths[2], rowHeight, claveSAT, "1", 0, "C", fillColor, 0, "")
		x += colWidths[2]

		unidadSAT := concepto.ClaveUnidad
		if unidadSAT == "" || unidadSAT == "0" {
			unidadSAT = "PZA"
		}
		pdf.SetXY(x, y)
		pdf.CellFormat(colWidths[3], rowHeight, unidadSAT, "1", 0, "C", fillColor, 0, "")
		x += colWidths[3]

		pdf.SetXY(x, y)
		pdf.CellFormat(colWidths[4], rowHeight, fmt.Sprintf("%.0f", concepto.Cantidad), "1", 0, "C", fillColor, 0, "")
		x += colWidths[4]

		pdf.SetXY(x, y)
		pdf.CellFormat(colWidths[5], rowHeight, fmt.Sprintf("$%.2f", concepto.ValorUnitario), "1", 0, "C", fillColor, 0, "")
		x += colWidths[5]

		// IVA ($) SOLO IMPORTE
		importeConcepto := concepto.Cantidad * concepto.ValorUnitario
		importeIVA := importeConcepto * (concepto.TasaIVA / 100)
		pdf.SetXY(x, y)
		pdf.CellFormat(colWidths[6], rowHeight, fmt.Sprintf("$%.2f", importeIVA), "1", 0, "C", fillColor, 0, "")
		x += colWidths[6]

		// TOTAL
		pdf.SetXY(x, y)
		pdf.CellFormat(colWidths[7], rowHeight, fmt.Sprintf("$%.2f", importeConcepto), "1", 0, "C", fillColor, 0, "")

		y += rowHeight
	}

	// Verificar espacio para el total
	if y > 260 {
		pdf.AddPage()
		y = 20
	}

	// Calcular subtotal e impuestos
	subtotal := 0.0
	totalIVA := 0.0
	totalIEPS := 0.0

	for _, concepto := range conceptosParaPDF {
		importeConcepto := concepto.Cantidad * concepto.ValorUnitario
		subtotal += importeConcepto
		totalIVA += importeConcepto * (concepto.TasaIVA / 100)
		totalIEPS += importeConcepto * (concepto.TasaIEPS / 100)
	}

	// Mostrar desglose de totales
	y += 10
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(0, 0, 0) // Texto negro

	// SUBTOTAL
	pdf.SetXY(15, y)
	pdf.CellFormat(150, 6, tr("SUBTOTAL:"), "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 6, fmt.Sprintf("$%.2f", subtotal), "0", 1, "R", false, 0, "")
	y += 6

	// IVA (16%)
	pdf.SetXY(15, y)
	pdf.CellFormat(150, 6, tr("IVA (16%):"), "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 6, fmt.Sprintf("$%.2f", totalIVA), "0", 1, "R", false, 0, "")
	y += 6

	// TOTAL (usando el valor original de la factura)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(15, y)
	pdf.CellFormat(150, 8, tr("TOTAL:"), "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", factura.Total), "0", 1, "R", false, 0, "")

	// Información adicional si existe
	if factura.Observaciones != "" {
		// Calcular el espacio necesario para las observaciones
		observaciones := factura.Observaciones
		lineWidth := 170.0
		lines := splitTextSafely(pdf, tr(observaciones), lineWidth)
		spaceNeededForObservaciones := float64(len(lines)*5) + 20 // 5mm por línea + título + margen

		// Verificar si hay suficiente espacio en la página actual
		currentY := y + 30
		if currentY+spaceNeededForObservaciones > 280 { // Si no cabe en la página actual
			pdf.AddPage()
			currentY = 20
		}

		y = currentY
		pdf.SetFont("Arial", "B", 10)
		pdf.SetXY(10, y)
		pdf.CellFormat(180, 6, tr("OBSERVACIONES:"), "", 0, "L", false, 0, "")
		y += 8

		pdf.SetFont("Arial", "", 9)
		for _, line := range lines {
			pdf.SetXY(10, y)
			pdf.CellFormat(180, 5, line, "", 0, "L", false, 0, "")
			y += 5
		}
	}

	// Guardar el PDF en un buffer
	var pdfBuffer bytes.Buffer
	err = pdf.Output(&pdfBuffer)
	if err != nil {
		return nil, "", fmt.Errorf("error al generar PDF: %w", err)
	}

	// Regresa el buffer y el nombre de archivo sugerido
	return &pdfBuffer, nombreArchivo, nil
}

// CargarLogoDesdeBaseDatos carga el logo de un usuario desde la base de datos
func CargarLogoDesdeBaseDatos(idUsuario int, logoService *LogoService) ([]byte, error) {
	if logoService == nil {
		return nil, fmt.Errorf("servicio de logo no inicializado")
	}

	logoBytes, err := logoService.CargarLogoDesdeBaseDatos(idUsuario)
	if err != nil {
		log.Printf("Error al cargar logo desde BD para usuario %d: %v", idUsuario, err)
		return nil, err
	}

	return logoBytes, nil
}

// CargarLogoPlantilla carga el logo desde la base de datos
func CargarLogoPlantilla(idUsuario string) ([]byte, error) {
	if idUsuario == "" {
		return nil, nil
	}

	// Convertir idUsuario a int
	idUsuarioInt, err := strconv.Atoi(idUsuario)
	if err != nil {
		log.Printf("Error al convertir idUsuario a entero: %v", err)
		return nil, nil
	}

	// Crear servicio de logos
	logoService := NewLogoService()

	// Obtener logo activo del usuario desde la base de datos
	logoBytes, err := logoService.CargarLogoDesdeBaseDatos(idUsuarioInt)
	if err != nil {
		log.Printf("No se encontró logo activo para el usuario %s: %v", idUsuario, err)
		return nil, nil // No es un error fatal, simplemente no hay logo
	}

	return logoBytes, nil
}

// Mapa de descripciones de régimen fiscal SAT (puedes completarlo con todos los códigos oficiales)
var regimenFiscalDescripciones = map[string]string{
	// Códigos oficiales del SAT
	"601": "General de Ley Personas Morales",
	"603": "Personas Morales con Fines no Lucrativos",
	"605": "Sueldos y Salarios e Ingresos Asimilados a Salarios",
	"606": "Arrendamiento",
	"607": "Régimen de Enajenación o Adquisición de Bienes",
	"608": "Demás ingresos",
	"609": "Consolidación",
	"610": "Residentes en el Extranjero sin Establecimiento Permanente en México",
	"611": "Ingresos por Dividendos (socios y accionistas)",
	"612": "Personas Físicas con Actividades Empresariales y Profesionales",
	"614": "Ingresos por intereses",
	"615": "Régimen de los ingresos por obtención de premios",
	"616": "Sin obligaciones fiscales",
	"620": "Sociedades Cooperativas de Producción que optan por diferir sus ingresos",
	"621": "Incorporación Fiscal",
	"622": "Actividades Agrícolas, Ganaderas, Silvícolas y Pesqueras",
	"623": "Opcional para Grupos de Sociedades",
	"624": "Coordinados",
	"628": "Hidrocarburos",
	"629": "De los Regímenes Fiscales Preferentes y de las Empresas Multinacionales",
	"630": "Enajenación de acciones en bolsa de valores",
	"626": "Régimen Simplificado de Confianza",

	// Códigos adicionales que podrían aparecer en sistemas locales
	"1":  "Régimen General",
	"2":  "Personas Físicas con Actividades Empresariales",
	"3":  "Régimen de Incorporación Fiscal",
	"4":  "Arrendamiento",
	"5":  "Honorarios",
	"6":  "Sueldos y Salarios",
	"7":  "Sin obligaciones fiscales",
	"8":  "Actividades Agrícolas",
	"9":  "Demás ingresos",
	"10": "Residentes en el Extranjero",
}
