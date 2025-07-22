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
	"strconv"
	"time"

	"Facts/internal/db"
	"Facts/internal/models"
	"Facts/internal/services"
	"Facts/internal/utils"
)

func DescargarFacturaHandler(w http.ResponseWriter, r *http.Request, facturaID int) {
	factura, err := models.ObtenerFacturaPorID(facturaID)
	if err != nil {
		log.Printf("Error al obtener factura con ID %d: %v", facturaID, err)
		utils.RespondWithError(w, fmt.Sprintf("Error al obtener factura: %v", err))
		return
	}

	// Obtener serie_df desde datos fiscales
	serieDF := ""
	datosFiscales, err := db.ObtenerDatosFiscales(factura.IDUsuario)
	if err == nil {
		if s, ok := datosFiscales["serie_df"].(string); ok {
			serieDF = s
		}
	}

	// Si no hay NumeroFolio, generar uno automáticamente con la serie
	if factura.NumeroFolio == "" {
		folio, _ := models.GetFolioGenerator().GenerarFolioSimple(serieDF)
		factura.NumeroFolio = folio
	}

	// Siempre regenerar el ZIP en un archivo temporal, servirlo y eliminarlo
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("factura_%d_*.zip", factura.ID))
	if err != nil {
		log.Printf("Error al crear archivo temporal para factura: %v", err)
		utils.RespondWithError(w, "Error al preparar la descarga de la factura")
		return
	}
	tmpFilePath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpFilePath)

	if err := regenerarFactura(factura, tmpFilePath, serieDF); err != nil {
		log.Printf("Error al regenerar factura: %v", err)
		utils.RespondWithError(w, "Error al regenerar la factura para descarga")
		return
	}

	servirArchivoFactura(w, r, tmpFilePath, fmt.Sprintf("factura_%d.zip", factura.ID))
}

func servirArchivoFactura(w http.ResponseWriter, r *http.Request, rutaArchivo, nombreArchivo string) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", nombreArchivo))
	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, rutaArchivo)
}

func regenerarFactura(factura *models.HistorialFactura, rutaArchivo, serieDF string) error {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	xmlData, xmlFileName, err := generarXMLFacturaConNombre(factura, serieDF)
	if err != nil {
		return fmt.Errorf("error al generar XML: %v", err)
	}

	xmlFile, err := zipWriter.Create(xmlFileName)
	if err != nil {
		return fmt.Errorf("error al crear archivo XML en ZIP: %v", err)
	}
	if _, err = xmlFile.Write(xmlData); err != nil {
		return fmt.Errorf("error al escribir XML: %v", err)
	}

	pdfData, pdfFileName, err := generarPDFFacturaConNombre(factura, serieDF)
	if err != nil {
		return fmt.Errorf("error al generar PDF: %v", err)
	}

	pdfFile, err := zipWriter.Create(pdfFileName)
	if err != nil {
		return fmt.Errorf("error al crear archivo PDF en ZIP: %v", err)
	}
	if _, err = pdfFile.Write(pdfData); err != nil {
		return fmt.Errorf("error al escribir PDF: %v", err)
	}

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

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("error al cerrar ZIP: %v", err)
	}

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

func generarXMLFacturaConNombre(factura *models.HistorialFactura, serieDF string) ([]byte, string, error) {
	empresa, err := obtenerEmpresaEmisoraParaFactura(factura)
	if err != nil {
		return nil, "", fmt.Errorf("error al obtener datos de empresa emisora: %v", err)
	}

	var fechaXML string
	fechaParsed, err := time.Parse("2006-01-02 15:04:05", factura.FechaGeneracion)
	if err != nil {
		fechaXML = factura.FechaGeneracion
	} else {
		fechaXML = fechaParsed.Format(time.RFC3339)
	}

	folio := factura.NumeroFolio
	if folio == "" {
		folio = strconv.Itoa(factura.ID)
	}

	cfdi := models.CFDI{
		XmlnsCfdi:         "http://www.sat.gob.mx/cfd/4",
		Version:           "4.0",
		Serie:             serieDF,
		Folio:             folio,
		Fecha:             fechaXML,
		SubTotal:          fmt.Sprintf("%.2f", factura.Total/1.16),
		Total:             fmt.Sprintf("%.2f", factura.Total),
		Moneda:            "MXN",
		LugarExpedicion:   empresa.CodigoPostal,
		TipoDeComprobante: "I",
		MetodoPago:        "PUE",
		FormaPago:         "01",
	}

	cfdi.Emisor.RFC = empresa.RFC
	cfdi.Emisor.Nombre = empresa.RazonSocial
	cfdi.Emisor.RegimenFiscal = empresa.RegimenFiscal

	cfdi.Receptor.RFC = factura.RFCReceptor
	cfdi.Receptor.Nombre = factura.RazonSocialReceptor
	cfdi.Receptor.UsoCFDI = factura.UsoCFDI

	concepto := struct {
		ClaveProdServ    string `xml:"ClaveProdServ,attr"`
		NoIdentificacion string `xml:"NoIdentificacion,attr,omitempty"`
		Cantidad         string `xml:"Cantidad,attr"`
		ClaveUnidad      string `xml:"ClaveUnidad,attr"`
		Unidad           string `xml:"Unidad,attr,omitempty"`
		Descripcion      string `xml:"Descripcion,attr"`
		ValorUnitario    string `xml:"ValorUnitario,attr"`
		Importe          string `xml:"Importe,attr"`
		Descuento        string `xml:"Descuento,attr,omitempty"`
	}{
		ClaveProdServ: "01010101",
		Cantidad:      "1",
		ClaveUnidad:   "ACT",
		Descripcion:   fmt.Sprintf("Venta relacionada con ticket %s", factura.ClaveTicket),
		ValorUnitario: fmt.Sprintf("%.2f", factura.Total/1.16),
		Importe:       fmt.Sprintf("%.2f", factura.Total/1.16),
		Descuento:     "0.00",
	}

	cfdi.Conceptos = append(cfdi.Conceptos, concepto)

	cfdi.Impuestos.TotalImpuestosTrasladados = fmt.Sprintf("%.2f", factura.Total-(factura.Total/1.16))
	traslado := struct {
		Impuesto   string `xml:"Impuesto,attr"`
		TipoFactor string `xml:"TipoFactor,attr"`
		TasaOCuota string `xml:"TasaOCuota,attr"`
		Importe    string `xml:"Importe,attr"`
	}{
		Impuesto:   "002",
		TipoFactor: "Tasa",
		TasaOCuota: "0.160000",
		Importe:    fmt.Sprintf("%.2f", factura.Total-(factura.Total/1.16)),
	}

	cfdi.Impuestos.Traslados = append(cfdi.Impuestos.Traslados, traslado)

	xmlBytes, err := xml.MarshalIndent(cfdi, "", "  ")
	if err != nil {
		return nil, "", fmt.Errorf("error al serializar XML: %v", err)
	}

	nombreArchivo := fmt.Sprintf("Factura_%s%s.xml", serieDF, folio)
	return xmlBytes, nombreArchivo, nil
}

func generarPDFFacturaConNombre(factura *models.HistorialFactura, serieDF string) ([]byte, string, error) {
	folio := factura.NumeroFolio
	if folio == "" {
		folio = strconv.Itoa(factura.ID)
	}

	facturaParaPDF := models.Factura{
		ID:                  strconv.Itoa(factura.ID),
		Serie:               serieDF,
		NumeroFolio:         folio,
		ReceptorRFC:         factura.RFCReceptor,
		ReceptorRazonSocial: factura.RazonSocialReceptor,
		ClaveTicket:         factura.ClaveTicket,
		Total:               factura.Total,
		UsoCFDI:             factura.UsoCFDI,
		Subtotal:            factura.Total / 1.16,
		Impuestos:           factura.Total - (factura.Total / 1.16),
		Conceptos: []models.Concepto{
			{
				Descripcion:   fmt.Sprintf("Venta relacionada con ticket %s", factura.ClaveTicket),
				Cantidad:      1,
				ValorUnitario: factura.Total / 1.16,
				Importe:       factura.Total / 1.16,
			},
		},
	}

	empresaEmisora, err := obtenerEmpresaEmisoraParaFactura(factura)
	if err != nil {
		return nil, "", fmt.Errorf("error al obtener empresa emisora: %v", err)
	}

	var logoBytes []byte
	logoBytes, err = services.CargarLogoPlantilla("1")
	if err != nil {
		logoBytes = nil
	}

	pdfBuffer, nombreArchivo, err := services.GenerarPDF(facturaParaPDF, empresaEmisora, logoBytes)
	if err != nil {
		return nil, "", fmt.Errorf("error al generar PDF: %v", err)
	}

	return pdfBuffer.Bytes(), nombreArchivo, nil
}

func obtenerEmpresaEmisoraParaFactura(factura *models.HistorialFactura) (*models.Empresa, error) {
	empresas, err := models.ObtenerEmpresasPorUsuario(factura.IDUsuario)
	if err != nil {
		return nil, fmt.Errorf("error al obtener empresas del usuario: %v", err)
	}
	if len(empresas) == 0 {
		return nil, fmt.Errorf("no se encontraron empresas para el usuario ID %d", factura.IDUsuario)
	}
	return &empresas[0], nil
}
