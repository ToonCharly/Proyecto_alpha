package utils

import (
	"crypto/x509"
	"encoding/hex"
	"errors"
	"strings"
)

// ObtenerNoSerieCER extrae el número de serie del certificado CSD (.cer) recibido como []byte
func ObtenerNoSerieCER(cerBytes []byte) (string, error) {
	if len(cerBytes) == 0 {
		return "", errors.New("el archivo .cer está vacío")
	}
	cert, err := x509.ParseCertificate(cerBytes)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(cert.SerialNumber.Bytes())), nil
}
