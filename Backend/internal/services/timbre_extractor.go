package services

import (
	"encoding/xml"
)

// Estructura para parsear el nodo TimbreFiscalDigital del XML timbrado
// Puedes agregar más campos si lo necesitas

type tfdXML struct {
	XMLName          xml.Name `xml:"TimbreFiscalDigital"`
	UUID             string   `xml:"UUID,attr"`
	FechaTimbrado    string   `xml:"FechaTimbrado,attr"`
	RfcProvCertif    string   `xml:"RfcProvCertif,attr"`
	SelloCFD         string   `xml:"SelloCFD,attr"`
	NoCertificadoSAT string   `xml:"NoCertificadoSAT,attr"`
	SelloSAT         string   `xml:"SelloSAT,attr"`
}

// Extrae los datos del timbre fiscal digital del XML timbrado
// ...existing code...
