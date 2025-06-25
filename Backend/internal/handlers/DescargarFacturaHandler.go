package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"carlos/Facts/Backend/internal/models"
	"carlos/Facts/Backend/internal/services"
	"carlos/Facts/Backend/internal/utils"
)

// DescargarFacturaHandler maneja la descarga de una factura del historial
func DescargarFacturaHandler(w http.ResponseWriter, r *http.Request, facturaID int) {
	// Obtener información de la factura desde el historial
	factura, err := models.ObtenerFacturaPorID(facturaID)
	if err != nil {
		log.Printf("Error al obtener factura con ID %d: %v", facturaID, err)
		utils.RespondWithError(w, fmt.Sprintf("Error al obtener factura: %v", err))
		return
	}

	// Definir directorios y nombres de archivo
	directorioFacturas := filepath.Join("facturas", strconv.Itoa(factura.IDUsuario))
	nombreArchivo := fmt.Sprintf("factura_%d.zip", factura.ID)
	rutaArchivo := filepath.Join(directorioFacturas, nombreArchivo)

	// Verificar si el directorio existe, si no, crearlo
	if err := os.MkdirAll(directorioFacturas, 0755); err != nil {
		log.Printf("Error al crear directorio para facturas: %v", err)
		utils.RespondWithError(w, "Error al preparar la descarga de la factura")
		return
	}

	// PRIMERA FORMA: Verificar si el archivo ya existe y enviarlo
	if _, err := os.Stat(rutaArchivo); err == nil {
		// El archivo existe, lo enviamos directamente
		servirArchivoFactura(w, r, rutaArchivo, nombreArchivo)
		return
	}

	// SEGUNDA FORMA: Regenerar la factura a partir de los datos del historial
	if err := regenerarFactura(factura, rutaArchivo); err != nil {
		log.Printf("Error al regenerar factura: %v", err)
		utils.RespondWithError(w, "Error al regenerar la factura para descarga")
		return
	}

	// Enviar el archivo regenerado
	servirArchivoFactura(w, r, rutaArchivo, nombreArchivo)
}

// servirArchivoFactura envía el archivo de factura como respuesta HTTP
func servirArchivoFactura(w http.ResponseWriter, r *http.Request, rutaArchivo, nombreArchivo string) {
	// Configurar encabezados para la descarga
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", nombreArchivo))
	w.Header().Set("Content-Type", "application/zip")

	// Servir el archivo
	http.ServeFile(w, r, rutaArchivo)
}

// regenerarFactura crea nuevamente el archivo de factura a partir de los datos del historial
func regenerarFactura(factura *models.HistorialFactura, rutaArchivo string) error {
	// Crear un buffer para el archivo ZIP
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Crear factura XML
	xmlData, err := generarXMLFactura(factura)
	if err != nil {
		return fmt.Errorf("error al generar XML: %v", err)
	}

	// Añadir XML al ZIP
	xmlFile, err := zipWriter.Create("factura.xml")
	if err != nil {
		return fmt.Errorf("error al crear archivo XML en ZIP: %v", err)
	}
	if _, err = xmlFile.Write(xmlData); err != nil {
		return fmt.Errorf("error al escribir XML: %v", err)
	}

	// Generar PDF de la factura
	pdfData, err := generarPDFFactura(factura)
	if err != nil {
		return fmt.Errorf("error al generar PDF: %v", err)
	}

	// Añadir PDF al ZIP
	pdfFile, err := zipWriter.Create("factura.pdf")
	if err != nil {
		return fmt.Errorf("error al crear archivo PDF en ZIP: %v", err)
	}
	if _, err = pdfFile.Write(pdfData); err != nil {
		return fmt.Errorf("error al escribir PDF: %v", err)
	}

	// Añadir archivo JSON con los datos
	jsonFile, err := zipWriter.Create("datos.json")
	if err != nil {
		return fmt.Errorf("error al crear archivo JSON en ZIP: %v", err)
	}
	jsonData, err := json.MarshalIndent(factura, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar JSON: %v", err)
	}
	if _, err = jsonFile.Write(jsonData); err != nil {
		return fmt.Errorf("error al escribir JSON: %v", err)
	}

	// Cerrar el ZIP writer para finalizar
	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("error al cerrar ZIP: %v", err)
	}

	// Guardar el archivo ZIP en el sistema de archivos
	file, err := os.Create(rutaArchivo)
	if err != nil {
		return fmt.Errorf("error al crear archivo ZIP: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, buf); err != nil {
		return fmt.Errorf("error al escribir archivo ZIP: %v", err)
	}

	return nil
}

// generarXMLFactura crea el XML de la factura basado en los datos del historial
func generarXMLFactura(factura *models.HistorialFactura) ([]byte, error) {
	// Obtener datos adicionales necesarios (empresa emisora, etc.)
	empresa, err := obtenerEmpresaEmisoraParaFactura(factura)
	if err != nil {
		return nil, fmt.Errorf("error al obtener datos de empresa emisora: %v", err)
	}

	// Parsear la fecha de string a time.Time para poder formatearla correctamente
	var fechaXML string
	fechaParsed, err := time.Parse("2006-01-02 15:04:05", factura.FechaGeneracion)
	if err != nil {
		// Si no podemos parsear la fecha, usamos el string original
		log.Printf("Advertencia: No se pudo parsear la fecha: %v. Usando fecha original.", err)
		fechaXML = factura.FechaGeneracion
	} else {
		// Si podemos parsear, la formateamos en RFC3339
		fechaXML = fechaParsed.Format(time.RFC3339)
	}

	// Crear estructura XML CFDI
	cfdi := models.CFDI{
		XmlnsCfdi:         "http://www.sat.gob.mx/cfd/4",
		Version:           "4.0",
		Serie:             "A",
		Folio:             strconv.Itoa(factura.ID),
		Fecha:             fechaXML,
		SubTotal:          fmt.Sprintf("%.2f", factura.Total/1.16), // Asumiendo IVA del 16%
		Total:             fmt.Sprintf("%.2f", factura.Total),
		Moneda:            "MXN",
		LugarExpedicion:   empresa.CodigoPostal,
		TipoDeComprobante: "I",   // Ingreso
		MetodoPago:        "PUE", // Pago en una sola exhibición
		FormaPago:         "01",  // Efectivo (por defecto)
	}

	// Datos del emisor
	cfdi.Emisor.RFC = empresa.RFC
	cfdi.Emisor.Nombre = empresa.RazonSocial
	cfdi.Emisor.RegimenFiscal = empresa.RegimenFiscal

	// Datos del receptor
	cfdi.Receptor.RFC = factura.RFCReceptor
	cfdi.Receptor.Nombre = factura.RazonSocialReceptor
	cfdi.Receptor.UsoCFDI = factura.UsoCFDI

	// Conceptos de factura (simulados o recuperados de otra tabla)
	concepto := struct {
		ClaveProdServ    string `xml:"ClaveProdServ,attr"`
		NoIdentificacion string `xml:"NoIdentificacion,attr,omitempty"`
		Cantidad         string `xml:"Cantidad,attr"`
		ClaveUnidad      string `xml:"ClaveUnidad,attr"`
		Unidad           string `xml:"Unidad,attr,omitempty"`
		Descripcion      string `xml:"Descripcion,attr"`
		ValorUnitario    string `xml:"ValorUnitario,attr"`
		Importe          string `xml:"Importe,attr"`
		Descuento        string `xml:"Descuento,attr,omitempty"` // Añadir este campo
	}{
		ClaveProdServ: "01010101", // Código genérico
		Cantidad:      "1",
		ClaveUnidad:   "ACT", // Actividad
		Descripcion:   fmt.Sprintf("Venta relacionada con ticket %s", factura.ClaveTicket),
		ValorUnitario: fmt.Sprintf("%.2f", factura.Total/1.16),
		Importe:       fmt.Sprintf("%.2f", factura.Total/1.16),
		Descuento:     "0.00", // Valor por defecto
	}

	cfdi.Conceptos = append(cfdi.Conceptos, concepto)

	// Impuestos
	cfdi.Impuestos.TotalImpuestosTrasladados = fmt.Sprintf("%.2f", factura.Total-(factura.Total/1.16))
	traslado := struct {
		Impuesto   string `xml:"Impuesto,attr"`
		TipoFactor string `xml:"TipoFactor,attr"`
		TasaOCuota string `xml:"TasaOCuota,attr"`
		Importe    string `xml:"Importe,attr"`
	}{
		Impuesto:   "002", // IVA
		TipoFactor: "Tasa",
		TasaOCuota: "0.160000", // 16%
		Importe:    fmt.Sprintf("%.2f", factura.Total-(factura.Total/1.16)),
	}

	cfdi.Impuestos.Traslados = append(cfdi.Impuestos.Traslados, traslado)

	// Convertir estructura a XML
	return xml.MarshalIndent(cfdi, "", "  ")
}

// generarPDFFactura crea el PDF de la factura usando los datos del historial
func generarPDFFactura(factura *models.HistorialFactura) ([]byte, error) {
	// Convertir HistorialFactura a la estructura Factura que espera el generador de PDF
	facturaParaPDF := models.Factura{
		ID:          strconv.Itoa(factura.ID), // Convertir int a string
		RFC:         factura.RFCReceptor,
		RazonSocial: factura.RazonSocialReceptor,
		ClaveTicket: factura.ClaveTicket,
		Total:       factura.Total,
		UsoCFDI:     factura.UsoCFDI,
		Subtotal:    factura.Total / 1.16, // Asumiendo IVA 16%
		Impuestos:   factura.Total - (factura.Total / 1.16),
		// Añadir un concepto genérico
		Conceptos: []models.Concepto{
			{
				Descripcion:   fmt.Sprintf("Venta relacionada con ticket %s", factura.ClaveTicket),
				Cantidad:      1,
				ValorUnitario: factura.Total / 1.16,
				Importe:       factura.Total / 1.16,
			},
		},
	}

	// Generar el PDF usando el servicio existente
	pdfBuffer, err := services.GenerarPDF(facturaParaPDF, nil) // Sin logo
	if err != nil {
		return nil, fmt.Errorf("error al generar PDF: %v", err)
	}

	return pdfBuffer.Bytes(), nil
}

// obtenerEmpresaEmisoraParaFactura recupera los datos de la empresa emisora
func obtenerEmpresaEmisoraParaFactura(factura *models.HistorialFactura) (*models.Empresa, error) {
	// Simplificamos la función para obtener directamente las empresas del usuario
	// Ya que la tabla relacion_factura_empresa no existe

	empresas, err := models.ObtenerEmpresasPorUsuario(factura.IDUsuario)
	if err != nil {
		return nil, fmt.Errorf("error al obtener empresas del usuario: %v", err)
	}

	if len(empresas) == 0 {
		return nil, fmt.Errorf("no se encontraron empresas para el usuario ID %d", factura.IDUsuario)
	}

	// Usar la primera empresa del usuario como emisora de la factura
	return &empresas[0], nil
}
