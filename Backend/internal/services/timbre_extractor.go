package services

import (
	"carlos/Facts/Backend/internal/models"
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
func ExtraerTimbreFiscalDigital(xmlTimbrado []byte) (*models.TimbreFiscalDigital, error) {
	var tfd tfdXML
	if err := xml.Unmarshal(xmlTimbrado, &tfd); err != nil {
		return nil, err
	}
	return &models.TimbreFiscalDigital{
		UUID:             tfd.UUID,
		FechaTimbrado:    tfd.FechaTimbrado,
		RfcProvCertif:    tfd.RfcProvCertif,
		SelloCFD:         tfd.SelloCFD,
		NoCertificadoSAT: tfd.NoCertificadoSAT,
		SelloSAT:         tfd.SelloSAT,
	}, nil
}
