package services

import (
	"Facts/internal/models"
	"Facts/internal/pac"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	outformPEM       = "-outform"
	pemErrorBlockMsg = "[PEM_ERROR] Bloque PEM completo:\n%s\n"
)

// Extrae el certificado y número de certificado desde base64 o archivo local
func asignarCertificadoCFDI(factura *models.Factura) error {
	var cerBytes []byte
	var err error
	// Si viene el certificado en base64, úsalo
	if factura.CerBase64 != "" {
		cerBytes, err = base64.StdEncoding.DecodeString(factura.CerBase64)
		if err != nil {
			return fmt.Errorf("no se pudo decodificar el certificado base64: %v", err)
		}
	} else {
		// Si no, usa la ruta local
		cerPath := strings.ReplaceAll(factura.CerPath, "\\", "/")
		absCerPath, err := filepath.Abs(cerPath)
		if err != nil {
			fmt.Printf("[CER_DEBUG] Error obteniendo ruta absoluta: %v\n", err)
			absCerPath = cerPath
		}
		fmt.Printf("[CER_DEBUG] Intentando leer archivo .cer en ruta: %s\n", absCerPath)
		cerBytes, err = os.ReadFile(absCerPath)
		if err != nil {
			fmt.Printf("[CER_DEBUG] Error al leer archivo .cer: %v\n", err)
			return fmt.Errorf("no se pudo leer el archivo .cer en %s: %v", absCerPath, err)
		}
	}
	var cert *x509.Certificate
	var certBlock *pem.Block
	certBlock, _ = pem.Decode(cerBytes)
	if certBlock != nil {
		cert, err = x509.ParseCertificate(certBlock.Bytes)
	} else {
		cert, err = x509.ParseCertificate(cerBytes)
	}
	if err != nil {
		return fmt.Errorf("no se pudo parsear el certificado: %v", err)
	}
	noCert := fmt.Sprintf("%X", cert.SerialNumber)
	if len(noCert) < 20 {
		noCert = fmt.Sprintf("%020s", noCert)
	}
	factura.NoCertificado = noCert
	certBase64 := base64.StdEncoding.EncodeToString(cert.Raw)
	factura.Certificado = certBase64
	fmt.Printf("[CER_DEBUG] NoCertificado: %s\n", noCert)
	fmt.Printf("[CER_DEBUG] Certificado (base64, primeros 100): %s\n", certBase64[:100])
	return nil
}

// Recibe el archivo .key cifrado y la clave, lo descifra, guarda el PEM en la base de datos y genera el XML CFDI firmado
func ProcesarKeyYGenerarCFDI(factura models.Factura, keyPath string, claveCSD string, xsltPath string) ([]byte, error) {
	// Asignar certificado y número de certificado
	err := asignarCertificadoCFDI(&factura)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo certificado: %w", err)
	}

	// Descifrar la llave .key como antes
	keyPathNorm := strings.ReplaceAll(keyPath, "\\", "/")
	absPath, err := filepath.Abs(keyPath)
	if err != nil {
		fmt.Printf("[KEY_DEBUG] Error obteniendo ruta absoluta: %v\n", err)
		absPath = keyPath
	}
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo acceder al archivo .key en %s: %v", absPath, err)
	}
	fmt.Printf("[KEY_DEBUG] Ruta absoluta .key: %s, tamaño: %d bytes\n", absPath, fileInfo.Size())
	fmt.Printf("[KEY_DEBUG] Ruta recibida para .key: %s\n", keyPath)
	fmt.Printf("[KEY_DEBUG] Ruta normalizada para OpenSSL: %s\n", keyPathNorm)
	pemTempPKCS8 := keyPath + "_pkcs8.pem"
	pemTempPKCS8Norm := strings.ReplaceAll(pemTempPKCS8, "\\", "/")
	pemTempRSA := keyPath + "_rsa.pem"
	pemTempRSANorm := strings.ReplaceAll(pemTempRSA, "\\", "/")
	fmt.Printf("[KEY_DEBUG] Archivo temporal PKCS#8: %s\n", pemTempPKCS8)
	fmt.Printf("[KEY_DEBUG] Archivo temporal RSA: %s\n", pemTempRSA)

	cmd1 := exec.Command("openssl", "pkcs8", "-in", keyPathNorm, "-inform", "DER", "-passin", "pass:"+claveCSD, "-out", pemTempPKCS8Norm, outformPEM, "PEM")
	errCmd1 := cmd1.Run()
	if errCmd1 != nil {
		fmt.Printf("[KEY_DEBUG] Falló descifrado DER, intentando PEM...\n")
		cmd1pem := exec.Command("openssl", "pkcs8", "-in", keyPathNorm, "-inform", "PEM", "-passin", "pass:"+claveCSD, "-out", pemTempPKCS8Norm, outformPEM, "PEM")
		errCmd1pem := cmd1pem.Run()
		if errCmd1pem != nil {
			return nil, fmt.Errorf("error descifrando .key a PKCS#8 con OpenSSL (DER y PEM): %v | %v", errCmd1, errCmd1pem)
		}
	}

	cmd2 := exec.Command("openssl", "rsa", "-in", pemTempPKCS8Norm, "-out", pemTempRSANorm, outformPEM, "PEM")
	if err := cmd2.Run(); err != nil {
		os.Remove(pemTempPKCS8)
		return nil, fmt.Errorf("error convirtiendo PKCS#8 a RSA PRIVATE KEY con OpenSSL: %v", err)
	}
	defer os.Remove(pemTempPKCS8)
	defer os.Remove(pemTempRSA)

	pemBytes, err := os.ReadFile(pemTempRSA)
	if err != nil {
		return nil, fmt.Errorf("error leyendo PEM RSA descifrado: %v", err)
	}
	keyPEM := string(pemBytes)
	fmt.Printf("[KEY_DEBUG] === CONTENIDO COMPLETO PEM DESCIFRADO (RSA PRIVATE KEY) ===\n%s\n[KEY_DEBUG] === FIN CONTENIDO PEM DESCIFRADO ===\n", keyPEM)
	preview := keyPEM
	if len(preview) > 100 {
		preview = preview[:100]
	}
	fmt.Printf("[KEY_DEBUG] Inicio PEM descifrado: %s\n", preview)
	block, _ := pem.Decode(pemBytes)
	if block != nil {
		fmt.Printf("[KEY_DEBUG] Tipo de bloque PEM descifrado: %s\n", block.Type)
	}
	return FlujoCFDIFirmado(factura, keyPEM, xsltPath)
}

// Estructura exacta según el XML de ejemplo CFDI 4.0
type CFDIComprobante struct {
	XMLName           xml.Name      `xml:"cfdi:Comprobante"`
	XMLNS             string        `xml:"xmlns:cfdi,attr"`
	XMLNSXSI          string        `xml:"xmlns:xsi,attr"`
	XSISchemaLocation string        `xml:"xsi:schemaLocation,attr"`
	Version           string        `xml:"Version,attr"`
	Serie             string        `xml:"Serie,attr,omitempty"`
	Folio             string        `xml:"Folio,attr"`
	Fecha             string        `xml:"Fecha,attr"`
	Sello             string        `xml:"Sello,attr,omitempty"`
	FormaPago         string        `xml:"FormaPago,attr"`
	NoCertificado     string        `xml:"NoCertificado,attr"`
	Certificado       string        `xml:"Certificado,attr"`
	SubTotal          string        `xml:"SubTotal,attr"`
	Descuento         string        `xml:"Descuento,attr,omitempty"`
	Moneda            string        `xml:"Moneda,attr"`
	Total             string        `xml:"Total,attr"`
	TipoDeComprobante string        `xml:"TipoDeComprobante,attr"`
	MetodoPago        string        `xml:"MetodoPago,attr"`
	LugarExpedicion   string        `xml:"LugarExpedicion,attr"`
	Emisor            CFDIEmisor    `xml:"cfdi:Emisor"`
	Receptor          CFDIReceptor  `xml:"cfdi:Receptor"`
	Conceptos         CFDIConceptos `xml:"cfdi:Conceptos"`
	Impuestos         CFDIImpuestos `xml:"cfdi:Impuestos"`
}

type CFDIEmisor struct {
	Rfc           string `xml:"Rfc,attr"`
	Nombre        string `xml:"Nombre,attr"`
	RegimenFiscal string `xml:"RegimenFiscal,attr"`
}

type CFDIReceptor struct {
	Rfc                     string `xml:"Rfc,attr"`
	Nombre                  string `xml:"Nombre,attr"`
	DomicilioFiscalReceptor string `xml:"DomicilioFiscalReceptor,attr"`
	RegimenFiscalReceptor   string `xml:"RegimenFiscalReceptor,attr"`
	UsoCFDI                 string `xml:"UsoCFDI,attr"`
}

type CFDIConceptos struct {
	Concepto []CFDIConcepto `xml:"cfdi:Concepto"`
}

type CFDIConcepto struct {
	ClaveProdServ    string                `xml:"ClaveProdServ,attr"`
	NoIdentificacion string                `xml:"NoIdentificacion,attr,omitempty"`
	Cantidad         string                `xml:"Cantidad,attr"`
	ClaveUnidad      string                `xml:"ClaveUnidad,attr"`
	Unidad           string                `xml:"Unidad,attr,omitempty"`
	Descripcion      string                `xml:"Descripcion,attr"`
	ValorUnitario    string                `xml:"ValorUnitario,attr"`
	Importe          string                `xml:"Importe,attr"`
	ObjetoImp        string                `xml:"ObjetoImp,attr"`
	Impuestos        CFDIConceptoImpuestos `xml:"cfdi:Impuestos"`
}

type CFDIConceptoImpuestos struct {
	Traslados CFDIConceptoTraslados `xml:"cfdi:Traslados"`
}

type CFDIConceptoTraslados struct {
	Traslado []CFDIConceptoTraslado `xml:"cfdi:Traslado"`
}

type CFDIConceptoTraslado struct {
	Base       string `xml:"Base,attr"`
	Impuesto   string `xml:"Impuesto,attr"`
	TipoFactor string `xml:"TipoFactor,attr"`
	TasaOCuota string `xml:"TasaOCuota,attr"`
	Importe    string `xml:"Importe,attr"`
}

type CFDIImpuestos struct {
	TotalImpuestosTrasladados string        `xml:"TotalImpuestosTrasladados,attr"`
	Traslados                 CFDITraslados `xml:"cfdi:Traslados"`
}

type CFDITraslados struct {
	Traslado []CFDITraslado `xml:"cfdi:Traslado"`
}

type CFDITraslado struct {
	Impuesto   string `xml:"Impuesto,attr"`
	TipoFactor string `xml:"TipoFactor,attr"`
	TasaOCuota string `xml:"TasaOCuota,attr"`
	Importe    string `xml:"Importe,attr"`
}

// Auxiliares
func formatFloat(f float64) string { return fmt.Sprintf("%.2f", f) }
func formatTasa(f float64) string  { return fmt.Sprintf("%.6f", f/100) }

// Valores seguros para receptor, para que nunca queden vacíos o inválidos
func safeReceptor(factura models.Factura) CFDIReceptor {
	// RFC receptor
	rfc := factura.ReceptorRFC
	if rfc == "" && factura.ClienteRFC != "" {
		rfc = factura.ClienteRFC
	}
	if rfc == "" {
		rfc = "XAXX010101000"
	}

	// Nombre receptor
	nombre := factura.ReceptorRazonSocial
	if nombre == "" && factura.ClienteRazonSocial != "" {
		nombre = factura.ClienteRazonSocial
	}
	if nombre == "" {
		nombre = "PUBLICO EN GENERAL"
	}

	// Código postal receptor
	cp := factura.ReceptorCodigoPostal
	if cp == "" {
		cp = "00000"
	}

	// Régimen fiscal receptor
	regimen := factura.RegimenFiscalReceptor
	if regimen == "" {
		regimen = factura.RegimenFiscal
	}
	if regimen == "" {
		regimen = "601"
	}

	// Uso CFDI
	uso := factura.UsoCFDI
	if uso == "" {
		uso = "G03"
	}

	return CFDIReceptor{
		Rfc:                     rfc,
		Nombre:                  nombre,
		DomicilioFiscalReceptor: cp,
		RegimenFiscalReceptor:   regimen,
		UsoCFDI:                 uso,
	}
}

// Firma la cadena original usando la llave privada PEM (archivo_key_pem)
func firmarCadenaOriginal(cadenaOriginal string, keyPEM string) (string, error) {
	// Log de depuración para ver el contenido recibido
	fmt.Println("[PEM_DEBUG] keyPEM recibido:\n" + keyPEM)
	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		fmt.Println("[PEM_ERROR] No se pudo decodificar el PEM de la llave privada.\nContenido recibido:\n" + keyPEM)
		return "", errors.New("no se pudo decodificar el PEM de la llave privada")
	}
	fmt.Printf("[PEM_DEBUG] Tipo de bloque PEM: %s\n", block.Type)
	var rsaPriv *rsa.PrivateKey
	var err error
	if block.Type == "PRIVATE KEY" {
		var key interface{}
		key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			fmt.Printf("[PEM_ERROR] Error en x509.ParsePKCS8PrivateKey: %v\n", err)
			fmt.Printf("[PEM_ERROR] Bytes del bloque PRIVATE KEY (hex): %x\n", block.Bytes)
			fmt.Printf(pemErrorBlockMsg, string(pem.EncodeToMemory(block)))
			return "", errors.New("la llave privada no es RSA o está dañada (PKCS8)")
		}
		var ok bool
		rsaPriv, ok = key.(*rsa.PrivateKey)
		if !ok {
			fmt.Printf("[PEM_ERROR] La llave PKCS#8 no es RSA\n")
			fmt.Printf(pemErrorBlockMsg, string(pem.EncodeToMemory(block)))
			return "", errors.New("la llave PKCS#8 no es RSA")
		}
	} else if block.Type == "RSA PRIVATE KEY" {
		rsaPriv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			fmt.Printf("[PEM_ERROR] Error en x509.ParsePKCS1PrivateKey: %v\n", err)
			fmt.Printf("[PEM_ERROR] Bytes del bloque RSA PRIVATE KEY (hex): %x\n", block.Bytes)
			fmt.Printf(pemErrorBlockMsg, string(pem.EncodeToMemory(block)))
			return "", errors.New("la llave privada no es RSA o está dañada (PKCS1)")
		}
	}
	if rsaPriv == nil {
		fmt.Printf("[PEM_ERROR] La llave privada no es RSA o está dañada. Bloque tipo: %s\n", block.Type)
		fmt.Printf(pemErrorBlockMsg, string(pem.EncodeToMemory(block)))
		return "", errors.New("la llave privada no es RSA o está dañada")
	}
	hash := crypto.SHA256.New()
	hash.Write([]byte(cadenaOriginal))
	hashed := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPriv, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// GenerarCadenaOriginal ejecuta xsltproc para obtener la cadena original CFDI 4.0
func GenerarCadenaOriginal(xmlPath, xsltPath string) (string, error) {
	fmt.Printf("[XSLT_DEBUG] Ejecutando xsltproc con XSLT: %s y XML: %s\n", xsltPath, xmlPath)
	cmd := exec.Command("xsltproc", xsltPath, xmlPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[XSLT_ERROR] xsltproc error: %v\n", err)
		fmt.Printf("[XSLT_ERROR] xsltproc output: %s\n", string(out))
		return "", fmt.Errorf("xsltproc error: %v, output: %s", err, string(out))
	}
	// Limpia saltos de línea y espacios extra
	return strings.TrimSpace(string(out)), nil
}

// Ejemplo de flujo completo: genera XML preliminar, obtiene cadena original y genera XML final firmado
// Nota: Este es un ejemplo, puedes adaptarlo a tu flujo real
func FlujoCFDIFirmado(factura models.Factura, keyPEM string, xsltPath string) ([]byte, error) {
	// Si xsltPath está vacío, usar la ruta por defecto del XSLT oficial CFDI 4.0
	if xsltPath == "" {
		xsltPath = "internal/services/cadenaoriginal_4_0.xslt"
	}
	// Convertir a ruta absoluta para evitar errores de xsltproc
	absXSLT, err := filepath.Abs(xsltPath)
	if err == nil {
		xsltPath = absXSLT
	}

	// 1. Generar XML preliminar (sin sello)
	tmpFile, err := os.CreateTemp("", "cfdi_pre_*.xml")
	if err != nil {
		return nil, fmt.Errorf("error creando archivo temporal: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	xmlPre, err := GenerarXML(factura)
	if err != nil {
		return nil, fmt.Errorf("error generando XML preliminar: %w", err)
	}
	if _, err := tmpFile.Write(xmlPre); err != nil {
		return nil, fmt.Errorf("error escribiendo XML temporal: %w", err)
	}
	tmpFile.Close()

	// 2. Generar cadena original usando xsltproc
	cadenaOriginal, err := GenerarCadenaOriginal(tmpFile.Name(), xsltPath)
	if err != nil {
		return nil, fmt.Errorf("error generando cadena original: %w", err)
	}

	// 3. Generar XML final con sello
	xmlFinal, err := GenerarXMLConSello(factura, keyPEM, cadenaOriginal)
	if err != nil {
		return nil, fmt.Errorf("error generando XML final con sello: %w", err)
	}
	return xmlFinal, nil
}

// GenerarXML convierte los datos de la factura en XML compatible con CFDI 4.0
func GenerarXML(factura models.Factura) ([]byte, error) {
	// Si la fecha de emisión no viene, asígnala en formato RFC3339 recortado a 19 caracteres
	if factura.FechaEmision == "" {
		t := time.Now().Format(time.RFC3339)
		if len(t) > 19 {
			t = t[:19]
		}
		factura.FechaEmision = t
	}
	// Serie nunca debe ser "undefined", "null" o vacía
	serie := factura.Serie
	if serie == "" || serie == "undefined" || serie == "null" {
		serie = "A"
	}

	var subtotal, totalImpuestos, totalDescuento float64

	conceptos := make([]CFDIConcepto, len(factura.Conceptos))
	for i, c := range factura.Conceptos {
		importe := c.Cantidad * c.ValorUnitario
		impuestoConcepto := importe * (c.TasaIVA / 100)
		subtotal += importe
		totalImpuestos += impuestoConcepto
		if c.Descuento > 0 {
			totalDescuento += c.Descuento
		}
		conceptos[i] = CFDIConcepto{
			ClaveProdServ:    c.ClaveProdServ,
			NoIdentificacion: "", // no existe en tu modelo
			Cantidad:         formatFloat(c.Cantidad),
			ClaveUnidad:      c.ClaveUnidad,
			Unidad:           "", // no existe en tu modelo
			Descripcion:      c.Descripcion,
			ValorUnitario:    formatFloat(c.ValorUnitario),
			Importe:          formatFloat(importe),
			ObjetoImp:        "02",
			Impuestos: CFDIConceptoImpuestos{
				Traslados: CFDIConceptoTraslados{
					Traslado: []CFDIConceptoTraslado{
						{
							Base:       formatFloat(importe),
							Impuesto:   "002",
							TipoFactor: "Tasa",
							TasaOCuota: formatTasa(c.TasaIVA),
							Importe:    formatFloat(impuestoConcepto),
						},
					},
				},
			},
		}
	}

	if factura.Descuento > 0 {
		totalDescuento += factura.Descuento
	}

	comprobante := CFDIComprobante{
		XMLNS:             "http://www.sat.gob.mx/cfd/4",
		XMLNSXSI:          "http://www.w3.org/2001/XMLSchema-instance",
		XSISchemaLocation: "http://www.sat.gob.mx/cfd/4 http://www.sat.gob.mx/sitio_internet/cfd/4/cfdv40.xsd",
		Version:           "4.0",
		Serie:             serie,
		Folio:             factura.NumeroFolio,
		Fecha:             factura.FechaEmision,
		Sello:             "",
		FormaPago:         ifEmpty(factura.FormaPago, "01"),
		NoCertificado:     factura.NoCertificado,
		Certificado:       factura.Certificado,
		SubTotal:          formatFloat(subtotal),
		Moneda:            ifEmpty(factura.Moneda, "MXN"),
		Total:             formatFloat(subtotal + totalImpuestos - totalDescuento),
		TipoDeComprobante: "I",
		MetodoPago:        ifEmpty(factura.MetodoPago, "PUE"),
		LugarExpedicion:   factura.EmisorCodigoPostal,
		Emisor: CFDIEmisor{
			Rfc:           factura.EmisorRFC,
			Nombre:        factura.EmisorRazonSocial,
			RegimenFiscal: factura.EmisorRegimenFiscal,
		},
		Receptor:  safeReceptor(factura),
		Conceptos: CFDIConceptos{Concepto: conceptos},
		Impuestos: CFDIImpuestos{
			TotalImpuestosTrasladados: formatFloat(totalImpuestos),
			Traslados: CFDITraslados{
				Traslado: []CFDITraslado{
					{
						Impuesto:   "002",
						TipoFactor: "Tasa",
						TasaOCuota: "0.160000",
						Importe:    formatFloat(totalImpuestos),
					},
				},
			},
		},
	}

	if totalDescuento > 0 {
		comprobante.Descuento = formatFloat(totalDescuento)
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")
	if err := encoder.Encode(comprobante); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerarXMLConSello genera el XML CFDI 4.0 y firma el sello usando el archivo_key_pem y la cadena original
func GenerarXMLConSello(factura models.Factura, keyPEM string, cadenaOriginal string) ([]byte, error) {
	// Si la fecha de emisión no viene, asígnala en formato RFC3339 recortado a 19 caracteres
	if factura.FechaEmision == "" {
		t := time.Now().Format(time.RFC3339)
		if len(t) > 19 {
			t = t[:19]
		}
		factura.FechaEmision = t
	}
	// Serie nunca debe ser "undefined", "null" o vacía
	serie := factura.Serie
	if serie == "" || serie == "undefined" || serie == "null" {
		serie = "A"
	}

	var subtotal, totalImpuestos, totalDescuento float64
	conceptos := make([]CFDIConcepto, len(factura.Conceptos))
	for i, c := range factura.Conceptos {
		importe := c.Cantidad * c.ValorUnitario
		impuestoConcepto := importe * (c.TasaIVA / 100)
		subtotal += importe
		totalImpuestos += impuestoConcepto
		if c.Descuento > 0 {
			totalDescuento += c.Descuento
		}
		conceptos[i] = CFDIConcepto{
			ClaveProdServ:    c.ClaveProdServ,
			NoIdentificacion: "",
			Cantidad:         formatFloat(c.Cantidad),
			ClaveUnidad:      c.ClaveUnidad,
			Unidad:           "",
			Descripcion:      c.Descripcion,
			ValorUnitario:    formatFloat(c.ValorUnitario),
			Importe:          formatFloat(importe),
			ObjetoImp:        "02",
			Impuestos: CFDIConceptoImpuestos{
				Traslados: CFDIConceptoTraslados{
					Traslado: []CFDIConceptoTraslado{{
						Base:       formatFloat(importe),
						Impuesto:   "002",
						TipoFactor: "Tasa",
						TasaOCuota: formatTasa(c.TasaIVA),
						Importe:    formatFloat(impuestoConcepto),
					}},
				},
			},
		}
	}
	if factura.Descuento > 0 {
		totalDescuento += factura.Descuento
	}
	comprobante := CFDIComprobante{
		XMLNS:             "http://www.sat.gob.mx/cfd/4",
		XMLNSXSI:          "http://www.w3.org/2001/XMLSchema-instance",
		XSISchemaLocation: "http://www.sat.gob.mx/cfd/4 http://www.sat.gob.mx/sitio_internet/cfd/4/cfdv40.xsd",
		Version:           "4.0",
		Serie:             serie,
		Folio:             factura.NumeroFolio,
		Fecha:             factura.FechaEmision,
		Sello:             "", // Se asigna después
		FormaPago:         ifEmpty(factura.FormaPago, "01"),
		NoCertificado:     factura.NoCertificado,
		Certificado:       factura.Certificado,
		SubTotal:          formatFloat(subtotal),
		Moneda:            ifEmpty(factura.Moneda, "MXN"),
		Total:             formatFloat(subtotal + totalImpuestos - totalDescuento),
		TipoDeComprobante: "I",
		MetodoPago:        ifEmpty(factura.MetodoPago, "PUE"),
		LugarExpedicion:   factura.EmisorCodigoPostal,
		Emisor: CFDIEmisor{
			Rfc:           factura.EmisorRFC,
			Nombre:        factura.EmisorRazonSocial,
			RegimenFiscal: factura.EmisorRegimenFiscal,
		},
		Receptor:  safeReceptor(factura),
		Conceptos: CFDIConceptos{Concepto: conceptos},
		Impuestos: CFDIImpuestos{
			TotalImpuestosTrasladados: formatFloat(totalImpuestos),
			Traslados: CFDITraslados{
				Traslado: []CFDITraslado{{
					Impuesto:   "002",
					TipoFactor: "Tasa",
					TasaOCuota: "0.160000",
					Importe:    formatFloat(totalImpuestos),
				}},
			},
		},
	}
	if totalDescuento > 0 {
		comprobante.Descuento = formatFloat(totalDescuento)
	}

	// --- Generar el sello digital ---
	sello, err := firmarCadenaOriginal(cadenaOriginal, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("error generando sello digital: %w", err)
	}
	comprobante.Sello = sello

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")
	if err := encoder.Encode(comprobante); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// --- PAC Integration & Timbre Extraction ---
// Estructura para el Timbre Fiscal Digital
type TimbreFiscalDigital struct {
	UUID             string
	SelloSAT         string
	SelloCFD         string
	NoCertificadoSAT string
	FechaTimbrado    string
	Version          string
}

// Genera el XML firmado, lo timbra con Solución Factible y retorna el XML timbrado y el timbre fiscal digital
func GenerarYTimbrarCFDIConPAC(factura models.Factura, keyPath string, claveCSD string, xsltPath string, pacUser string, pacPass string) ([]byte, *TimbreFiscalDigital, error) {
	// 1. Genera el XML firmado
	xmlFirmado, err := ProcesarKeyYGenerarCFDI(factura, keyPath, claveCSD, xsltPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error generando XML firmado: %w", err)
	}

	// 2. Timbrar con Solución Factible (PAC)
	endpoint := "https://demo-facturacion.solucionfactible.com/ws/rest/timbrado"
	xmlTimbradoBytes, err := pac.TimbrarConPAC(string(xmlFirmado), factura.EmisorRFC, pacPass, endpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("error timbrando con PAC: %w", err)
	}
	xmlTimbrado := string(xmlTimbradoBytes)

	// 3. Extraer el timbre fiscal digital
	timbre, err := ExtraerTimbreFiscalDigital([]byte(xmlTimbrado))
	if err != nil {
		return []byte(xmlTimbrado), nil, fmt.Errorf("error extrayendo timbre fiscal digital: %w", err)
	}

	return []byte(xmlTimbrado), timbre, nil
}

// Extrae el Timbre Fiscal Digital del XML timbrado
func ExtraerTimbreFiscalDigital(xmlTimbrado []byte) (*TimbreFiscalDigital, error) {
	type Timbre struct {
		XMLName          xml.Name `xml:"tfd:TimbreFiscalDigital"`
		UUID             string   `xml:"UUID,attr"`
		SelloSAT         string   `xml:"SelloSAT,attr"`
		SelloCFD         string   `xml:"SelloCFD,attr"`
		NoCertificadoSAT string   `xml:"NoCertificadoSAT,attr"`
		FechaTimbrado    string   `xml:"FechaTimbrado,attr"`
		Version          string   `xml:"Version,attr"`
	}
	// Buscar el nodo TimbreFiscalDigital en el XML
	var t Timbre
	err := xml.Unmarshal(xmlTimbrado, &t)
	if err == nil && t.UUID != "" {
		return &TimbreFiscalDigital{
			UUID:             t.UUID,
			SelloSAT:         t.SelloSAT,
			SelloCFD:         t.SelloCFD,
			NoCertificadoSAT: t.NoCertificadoSAT,
			FechaTimbrado:    t.FechaTimbrado,
			Version:          t.Version,
		}, nil
	}
	// Si falla el Unmarshal directo, buscar el nodo manualmente
	tfdStart := bytes.Index(xmlTimbrado, []byte("<tfd:TimbreFiscalDigital"))
	tfdEnd := bytes.Index(xmlTimbrado[tfdStart:], []byte("/>"))
	if tfdStart >= 0 && tfdEnd > 0 {
		tfdNode := xmlTimbrado[tfdStart : tfdStart+tfdEnd+2]
		var t2 Timbre
		err2 := xml.Unmarshal(tfdNode, &t2)
		if err2 == nil && t2.UUID != "" {
			return &TimbreFiscalDigital{
				UUID:             t2.UUID,
				SelloSAT:         t2.SelloSAT,
				SelloCFD:         t2.SelloCFD,
				NoCertificadoSAT: t2.NoCertificadoSAT,
				FechaTimbrado:    t2.FechaTimbrado,
				Version:          t2.Version,
			}, nil
		}
	}
	return nil, errors.New("No se encontró el nodo TimbreFiscalDigital en el XML timbrado")
}
