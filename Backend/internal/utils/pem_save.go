package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

// SavePEMFromDER guarda archivos PEM a partir de binarios DER (.cer/.key) y retorna el contenido PEM
func SavePEMFromDER(derData []byte, pemType string, filePath string) (string, error) {
	// Codifica a base64
	b64 := base64.StdEncoding.EncodeToString(derData)
	var pem string
	var err error
	if pemType == "CERTIFICATE" {
		pem, err = Base64ToPEMCert(b64)
	} else if pemType == "PRIVATE KEY" {
		pem, err = Base64ToPEMKey(b64)
	} else {
		return "", fmt.Errorf("tipo PEM no soportado: %s", pemType)
	}
	if err != nil {
		return "", err
	}
	if filePath != "" {
		if err := ioutil.WriteFile(filePath, []byte(pem), 0600); err != nil {
			return "", fmt.Errorf("error guardando archivo PEM: %v", err)
		}
	}
	return pem, nil
}

// Ejemplo de uso en tu endpoint de alta de datos fiscales:
//
// 1. Recibes el archivo .cer y .key como []byte (binario)
// 2. Llamas a esta función para obtener el PEM y lo guardas en la base de datos
//
// cerPEM, err := utils.SavePEMFromDER(archivoCer, "CERTIFICATE", "")
// keyPEM, err := utils.SavePEMFromDER(archivoKey, "PRIVATE KEY", "")
//
// Luego guarda cerPEM y keyPEM en los campos nuevos de la tabla datos_fiscales
