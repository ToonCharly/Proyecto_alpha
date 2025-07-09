package services

import (
	"archive/zip"
	"bytes"
	"fmt"
)

// CrearZIPConNombres crea un ZIP con nombres personalizados para los archivos PDF y XML
func CrearZIPConNombres(pdfBytes, xmlBytes []byte, nombrePDF, nombreXML string) (*bytes.Buffer, error) {
	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	pdfFile, err := zipWriter.Create(nombrePDF)
	if err != nil {
		return nil, fmt.Errorf("error al crear archivo PDF en ZIP: %w", err)
	}
	_, err = pdfFile.Write(pdfBytes)
	if err != nil {
		return nil, fmt.Errorf("error al escribir archivo PDF en ZIP: %w", err)
	}

	xmlFile, err := zipWriter.Create(nombreXML)
	if err != nil {
		return nil, fmt.Errorf("error al crear archivo XML en ZIP: %w", err)
	}
	_, err = xmlFile.Write(xmlBytes)
	if err != nil {
		return nil, fmt.Errorf("error al escribir archivo XML en ZIP: %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("error al cerrar archivo ZIP: %w", err)
	}

	return &zipBuffer, nil
}