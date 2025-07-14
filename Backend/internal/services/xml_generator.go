package services

import (
	"bytes"
	"carlos/Facts/Backend/internal/models"
	"encoding/xml"
	"fmt"
	"time"
)

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
	NoCertificado     string        `xml:"NoCertificado,attr,omitempty"`
	Certificado       string        `xml:"Certificado,attr,omitempty"`
	SubTotal          string        `xml:"SubTotal,attr"`
	Descuento         string        `xml:"Descuento,attr,omitempty"`
	Moneda            string        `xml:"Moneda,attr"`
	Total             string        `xml:"Total,attr"`
	TipoDeComprobante string        `xml:"TipoDeComprobante,attr"`
	Exportacion       string        `xml:"Exportacion,attr"`
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
func ifEmpty(value, def string) string {
	if value == "" {
		return def
	} else {
		return value
	}
}

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
		cp = factura.CodigoPostal
	}
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
		Certificado: factura.Certificado,
		SubTotal:          formatFloat(subtotal),
		Moneda:            ifEmpty(factura.Moneda, "MXN"),
		Total:             formatFloat(subtotal + totalImpuestos - totalDescuento),
		TipoDeComprobante: "I",
		Exportacion:       "01",
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