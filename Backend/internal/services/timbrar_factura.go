package services

import (
	"Facts/internal/models"
)

// TimbrarFactura envía el XML firmado al PAC, imprime el XML timbrado recibido y extrae el timbre fiscal digital
func TimbrarFactura(xmlFirmado []byte, usuarioPAC, clavePAC, endpoint string) (*models.TimbreFiscalDigital, []byte, error) {
	// 1. Envía el XML firmado al PAC
	xmlTimbrado, err := EnviarXMLAlPAC(xmlFirmado, usuarioPAC, clavePAC, endpoint)
	if err != nil {
		return nil, nil, err
	}

	// Log: imprimir el XML timbrado recibido
	println("[PAC_DEBUG] XML timbrado recibido:\n", string(xmlTimbrado))

	// 2. Extrae el timbre fiscal digital
	tfd, err := ExtraerTimbreFiscalDigital(xmlTimbrado)
	if err != nil {
		return nil, xmlTimbrado, err
	}
	// Convertir a models.TimbreFiscalDigital
	timbre := &models.TimbreFiscalDigital{
		UUID:             tfd.UUID,
		FechaTimbrado:    tfd.FechaTimbrado,
		RfcProvCertif:    "", // No disponible en tfd, dejar vacío o asignar otro campo si aplica
		SelloCFD:         tfd.SelloCFD,
		NoCertificadoSAT: tfd.NoCertificadoSAT,
		SelloSAT:         tfd.SelloSAT,
	}
	return timbre, xmlTimbrado, nil
}
