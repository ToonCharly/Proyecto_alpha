package services

import (
    "archive/zip"
    "bytes"
    "fmt"
)

func CrearZIP(pdfBytes, xmlBytes []byte) (*bytes.Buffer, error) {
    var zipBuffer bytes.Buffer
    zipWriter := zip.NewWriter(&zipBuffer)

    // Agregar el archivo PDF
    pdfFile, err := zipWriter.Create("factura.pdf")
    if err != nil {
        return nil, fmt.Errorf("error al crear archivo PDF en ZIP: %w", err)
    }
    _, err = pdfFile.Write(pdfBytes)
    if err != nil {
        return nil, fmt.Errorf("error al escribir archivo PDF en ZIP: %w", err)
    }

    // Agregar el archivo XML
    xmlFile, err := zipWriter.Create("factura.xml")
    if err != nil {
        return nil, fmt.Errorf("error al crear archivo XML en ZIP: %w", err)
    }
    _, err = xmlFile.Write(xmlBytes)
    if err != nil {
        return nil, fmt.Errorf("error al escribir archivo XML en ZIP: %w", err)
    }

    // Cerrar el ZIP
    if err := zipWriter.Close(); err != nil {
        return nil, fmt.Errorf("error al cerrar archivo ZIP: %w", err)
    }

    return &zipBuffer, nil
}