package services

import (
    "bytes"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/phpdave11/gofpdf"
    "carlos/Facts/Backend/internal/models"
)

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

func GenerarPDF(factura models.Factura, logoBytes []byte) (*bytes.Buffer, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.SetAuthor("Sistema de Facturación", true)
    pdf.SetTitle("Factura Electrónica", true)
    pdf.AddPage()
    pdf.SetMargins(15, 15, 15)

    // Configurar la fuente principal
    tr := pdf.UnicodeTranslatorFromDescriptor("")
    pdf.SetFont("Arial", "B", 20)
    pdf.SetTextColor(50, 50, 50)
    pdf.SetXY(20, 20)
    pdf.Cell(130, 10, tr("FACTURA ELECTRÓNICA"))

    // Agregar logo si está disponible
    if len(logoBytes) > 0 {
        tmpDir := os.TempDir()
        tmpfileName := filepath.Join(tmpDir, fmt.Sprintf("logo-%d.png", time.Now().UnixNano()))
        err := os.WriteFile(tmpfileName, logoBytes, 0644)
        if err == nil {
            defer os.Remove(tmpfileName)
            pdf.Image(tmpfileName, 155, 17, 38, 36, false, "", 0, "")
        }
    }

    // Agregar datos del emisor
    pdf.SetFont("Arial", "B", 12)
    pdf.SetTextColor(30, 80, 150)
    pdf.SetXY(15, 50)
    pdf.Cell(180, 8, tr("DATOS DEL EMISOR"))
    pdf.SetFont("Arial", "", 10)
    pdf.SetTextColor(50, 50, 50)
    pdf.SetXY(15, 60)
    pdf.Cell(180, 6, tr("Empresa Ejemplo SA de CV"))
    pdf.SetXY(15, 66)
    pdf.Cell(180, 6, tr("RFC: AAA010101AAA"))
    pdf.SetXY(15, 72)
    pdf.Cell(180, 6, tr("Régimen Fiscal: 601 - General de Ley Personas Morales"))

    // Agregar datos del receptor
    pdf.SetFont("Arial", "B", 12)
    pdf.SetTextColor(30, 80, 150)
    pdf.SetXY(15, 80)
    pdf.Cell(180, 8, tr("DATOS DEL RECEPTOR"))
    pdf.SetFont("Arial", "", 10)
    pdf.SetTextColor(50, 50, 50)
    pdf.SetXY(15, 90)
    pdf.Cell(180, 6, tr(fmt.Sprintf("Razón Social: %s", factura.RazonSocial)))
    pdf.SetXY(15, 96)
    pdf.Cell(180, 6, tr(fmt.Sprintf("RFC: %s", factura.RFC)))
    pdf.SetXY(15, 102)
    pdf.Cell(180, 6, tr(fmt.Sprintf("Uso CFDI: %s", obtenerDescripcionUsoCfdi(factura.UsoCFDI))))

    // Agregar detalles de la factura
    pdf.SetFont("Arial", "B", 10)
    pdf.SetFillColor(230, 230, 230)
    pdf.SetXY(15, 110)
    pdf.CellFormat(120, 8, tr("Concepto"), "1", 0, "C", true, 0, "")
    pdf.CellFormat(30, 8, tr("Cantidad"), "1", 0, "C", true, 0, "")
    pdf.CellFormat(30, 8, tr("Importe"), "1", 1, "C", true, 0, "")

    pdf.SetFont("Arial", "", 10)
    pdf.SetFillColor(255, 255, 255)
    for _, concepto := range factura.Conceptos {
        pdf.SetX(15)
        pdf.CellFormat(120, 8, tr(concepto.Descripcion), "1", 0, "L", true, 0, "")
        pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", concepto.Cantidad), "1", 0, "C", true, 0, "")
        pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", concepto.Importe), "1", 1, "R", true, 0, "")
    }

    // Agregar totales
    pdf.SetFont("Arial", "B", 10)
    pdf.SetXY(15, pdf.GetY()+10)
    pdf.CellFormat(150, 8, tr("Subtotal:"), "0", 0, "R", false, 0, "")
    pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", factura.Subtotal), "0", 1, "R", false, 0, "")

    pdf.SetXY(15, pdf.GetY())
    pdf.CellFormat(150, 8, tr("IVA (16%):"), "0", 0, "R", false, 0, "")
    pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", factura.Impuestos), "0", 1, "R", false, 0, "")

    pdf.SetXY(15, pdf.GetY())
    pdf.SetFont("Arial", "B", 12)
    pdf.SetTextColor(30, 80, 150)
    pdf.CellFormat(150, 8, tr("Total:"), "0", 0, "R", false, 0, "")
    pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", factura.Total), "0", 1, "R", false, 0, "")

    // Guardar el PDF en un buffer
    var pdfBuffer bytes.Buffer
    err := pdf.Output(&pdfBuffer)
    if err != nil {
        return nil, fmt.Errorf("error al generar PDF: %w", err)
    }

    return &pdfBuffer, nil
}