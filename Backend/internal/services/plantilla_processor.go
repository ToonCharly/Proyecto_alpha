package services

import (
    "bytes"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"

    "baliance.com/gooxml/document"
    "carlos/Facts/Backend/internal/models"
)

func ifEmpty(value, defaultValue string) string {
    if strings.TrimSpace(value) == "" {
        return defaultValue
    }
    return value
}
func formatearFecha(fechaISO string) string {
    // Parsear la fecha ISO
    t, err := time.Parse("2006-01-02T15:04:05", fechaISO)
    if err != nil {
        return fechaISO // Si hay error, retornar la fecha original
    }

    // Formatear fecha a un formato más legible
    return t.Format("02/01/2006")
}

func obtenerDescripcionRegimenFiscal(clave string) string {
    descripciones := map[string]string{
        "601": "General de Ley Personas Morales",
        "603": "Personas Morales con Fines no Lucrativos",
        "605": "Sueldos y Salarios e Ingresos Asimilados a Salarios",
        "606": "Arrendamiento",
        "608": "Demás ingresos",
        "612": "Personas Físicas con Actividades Empresariales y Profesionales",
        "621": "Incorporación Fiscal",
        "622": "Actividades Agrícolas, Ganaderas, Silvícolas y Pesqueras",
        "626": "Régimen Simplificado de Confianza",
    }

    if descripcion, ok := descripciones[clave]; ok {
        return descripcion
    }
    return clave
}

func ProcesarPlantilla(factura models.Factura, plantillaBytes []byte) (*bytes.Buffer, error) {
    log.Printf("Iniciando procesamiento de plantilla para factura de %s", factura.RazonSocial)

    tmpDir := os.TempDir()
    timestamp := time.Now().UnixNano()
    plantillaPath := filepath.Join(tmpDir, fmt.Sprintf("plantilla-%d.docx", timestamp))
    outputPath := filepath.Join(tmpDir, fmt.Sprintf("factura_rellena-%d.docx", timestamp))

    // Guardar plantilla temporal
    if err := os.WriteFile(plantillaPath, plantillaBytes, 0644); err != nil {
        return nil, fmt.Errorf("error al guardar plantilla temporal: %w", err)
    }
    defer os.Remove(plantillaPath)

    // Abrir el DOCX
    doc, err := document.Open(plantillaPath)
    if err != nil {
        return nil, fmt.Errorf("error al abrir documento DOCX: %w", err)
    }

    // Reemplazar placeholders
    placeholders := map[string]string{
        "{{DIRECCION}}":               ifEmpty(factura.Direccion, "Campo no completo"),
        "{{RFC}}":                     ifEmpty(factura.RFC, "Campo no completo"),
        "{{NUMERO_FOLIO}}":            ifEmpty(factura.ClaveTicket, "Campo no completo"),
        "{{FECHA_FACTURA}}":           ifEmpty(formatearFecha(factura.FechaEmision), "Campo no completo"),
        "{{REGIMEN_FISCAL}}":          ifEmpty(obtenerDescripcionRegimenFiscal(factura.RegimenFiscal), "Campo no completo"),
        "{{LUGAR_EXPEDICION}}":        ifEmpty(factura.CodigoPostal, "Campo no completo"),
        "{{RECEPTOR_RFC}}":            ifEmpty(factura.RFC, "Campo no completo"),
        "{{DOMICILIO_FISCAL}}":        ifEmpty(factura.Direccion, "Campo no completo"),
        "{{USO_CFDI}}":                ifEmpty(obtenerDescripcionUsoCfdi(factura.UsoCFDI), "Campo no completo"),
        "{{REGIMEN_FISCAL_RECEPTOR}}": ifEmpty(obtenerDescripcionRegimenFiscal(factura.RegimenFiscal), "Campo no completo"),
        "{{SUBTOTAL}}":                fmt.Sprintf("$%.2f", factura.Subtotal),
        "{{IVA}}":                     fmt.Sprintf("$%.2f", factura.Impuestos),
        "{{TOTAL}}":                   fmt.Sprintf("$%.2f", factura.Total),
    }

    for _, para := range doc.Paragraphs() {
        for _, run := range para.Runs() {
            text := run.Text()
            for placeholder, valor := range placeholders {
                if strings.Contains(text, placeholder) {
                    text = strings.ReplaceAll(text, placeholder, valor)
                }
            }
            run.ClearContent()
            run.AddText(text)
        }
    }

    // Guardar el archivo modificado
    if err := doc.SaveToFile(outputPath); err != nil {
        return nil, fmt.Errorf("error al guardar documento modificado: %w", err)
    }
    defer os.Remove(outputPath)

    // Localizar LibreOffice
    libreOfficePath := "C:/Program Files/LibreOffice/program/soffice.exe"
    if _, err := os.Stat(libreOfficePath); os.IsNotExist(err) {
        libreOfficePath = "C:/Program Files (x86)/LibreOffice/program/soffice.exe"
        if _, err := os.Stat(libreOfficePath); os.IsNotExist(err) {
            libreOfficePath, err = exec.LookPath("libreoffice")
            if err != nil {
                return nil, fmt.Errorf("LibreOffice no está instalado o no se encuentra en el PATH")
            }
        }
    }

    // Convertir DOCX a PDF
    cmd := exec.Command(libreOfficePath, "--headless", "--convert-to", "pdf", "--outdir", tmpDir, outputPath)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("error al convertir DOCX a PDF. Salida: %s. Error: %w", string(output), err)
    }

    // Buscar el archivo PDF más reciente en el directorio temporal
    files, err := os.ReadDir(tmpDir)
    if err != nil {
        return nil, fmt.Errorf("error al leer el directorio temporal: %w", err)
    }

    var latestPDF string
    var latestModTime time.Time

    for _, file := range files {
        if strings.HasSuffix(file.Name(), ".pdf") {
            info, err := os.Stat(filepath.Join(tmpDir, file.Name()))
            if err == nil && info.ModTime().After(latestModTime) {
                latestModTime = info.ModTime()
                latestPDF = filepath.Join(tmpDir, file.Name())
            }
        }
    }

    if latestPDF == "" {
        return nil, fmt.Errorf("no se encontró ningún archivo PDF generado")
    }
    defer os.Remove(latestPDF)

    // Leer PDF
    pdfBytes, err := os.ReadFile(latestPDF)
    if err != nil {
        return nil, fmt.Errorf("error al leer archivo PDF generado: %w", err)
    }

    return bytes.NewBuffer(pdfBytes), nil
}
