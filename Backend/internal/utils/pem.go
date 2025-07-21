package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Base64ToPEMCert convierte un certificado en base64 plano a formato PEM
func Base64ToPEMCert(base64Cert string) (string, error) {
	clean := strings.ReplaceAll(base64Cert, "\n", "")
	clean = strings.ReplaceAll(clean, "\r", "")
	clean = strings.TrimSpace(clean)

	_, err := base64.StdEncoding.DecodeString(clean)
	if err != nil {
		return "", fmt.Errorf("certificado base64 inválido: %v", err)
	}

	var pemLines []string
	for i := 0; i < len(clean); i += 64 {
		end := i + 64
		if end > len(clean) {
			end = len(clean)
		}
		pemLines = append(pemLines, clean[i:end])
	}

	pem := "-----BEGIN CERTIFICATE-----\n" +
		strings.Join(pemLines, "\n") +
		"\n-----END CERTIFICATE-----\n"
	return pem, nil
}

// Base64ToPEMKey convierte una llave privada en base64 plano a formato PEM
func Base64ToPEMKey(base64Key string) (string, error) {
	clean := strings.ReplaceAll(base64Key, "\n", "")
	clean = strings.ReplaceAll(clean, "\r", "")
	clean = strings.TrimSpace(clean)

	_, err := base64.StdEncoding.DecodeString(clean)
	if err != nil {
		return "", fmt.Errorf("llave base64 inválida: %v", err)
	}

	var pemLines []string
	for i := 0; i < len(clean); i += 64 {
		end := i + 64
		if end > len(clean) {
			end = len(clean)
		}
		pemLines = append(pemLines, clean[i:end])
	}

	pem := "-----BEGIN PRIVATE KEY-----\n" +
		strings.Join(pemLines, "\n") +
		"\n-----END PRIVATE KEY-----\n"
	return pem, nil
}
