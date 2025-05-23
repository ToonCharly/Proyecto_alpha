package services

import (
    "encoding/xml"
    "time"

    "carlos/Facts/Backend/internal/models"
)

// FacturaXML representa la estructura XML de una factura
type FacturaXML struct {
    XMLName      xml.Name  `xml:"Factura"`
    Version      string    `xml:"version,attr"`
    FechaEmision string    `xml:"FechaEmision"`
    Emisor       Emisor    `xml:"Emisor"`
    Receptor     Receptor  `xml:"Receptor"`
    Conceptos    Conceptos `xml:"Conceptos"`
    Total        float64   `xml:"Total"`
    Observacion  string    `xml:"Observacion,omitempty"`
    ClaveTicket  string    `xml:"ClaveTicket,omitempty"`
}

type Emisor struct {
    RFC          string `xml:"RFC"`
    RazonSocial  string `xml:"RazonSocial"`
    RegimenFiscal string `xml:"RegimenFiscal"`
}

type Receptor struct {
    RFC         string `xml:"RFC"`
    RazonSocial string `xml:"RazonSocial"`
    Direccion   string `xml:"Direccion"`
    CodigoPostal string `xml:"CodigoPostal"`
    UsoCFDI     string `xml:"UsoCFDI"`
    RegimenFiscal string `xml:"RegimenFiscal,omitempty"`
}

type Conceptos struct {
    Concepto []Concepto `xml:"Concepto"`
}

type Concepto struct {
    Descripcion string  `xml:"Descripcion"`
    Cantidad    float64 `xml:"Cantidad"`
    ValorUnitario float64 `xml:"ValorUnitario"`
    Importe     float64 `xml:"Importe"`
}

// GenerarXML convierte los datos de la factura en un archivo XML
func GenerarXML(factura models.Factura) ([]byte, error) {
    // Crear estructura XML
    facturaXML := FacturaXML{
        Version:      "1.0",
        FechaEmision: time.Now().Format("2006-01-02T15:04:05"),
        Emisor: Emisor{
            RFC:          factura.RFC,
            RazonSocial:  factura.RazonSocial,
            RegimenFiscal: factura.RegimenFiscal,
        },
        Receptor: Receptor{
            // CAMBIOS AQUÍ - Usar información clara del receptor
            RFC:         factura.ReceptorRFC,
            // Usar nombre inferido del RFC o un texto genérico
            RazonSocial: "Receptor: " + factura.ReceptorRFC,
            // Usar DomicilioFiscal específicamente para el receptor
            Direccion:   factura.DomicilioFiscal, 
            // Código postal genérico o usar un valor por defecto
            CodigoPostal: "00000", // Código postal genérico
            UsoCFDI:     factura.UsoCFDI,
            RegimenFiscal: factura.RegimenFiscalReceptor,
        },
        Total:       factura.Total,
        Observacion: factura.Observaciones,
        ClaveTicket: factura.ClaveTicket,
    }

    // Usar los conceptos de la factura si están disponibles
    if len(factura.Conceptos) > 0 {
        facturaXML.Conceptos.Concepto = make([]Concepto, len(factura.Conceptos))
        for i, c := range factura.Conceptos {
            facturaXML.Conceptos.Concepto[i] = Concepto{
                Descripcion:   c.Descripcion,
                Cantidad:      c.Cantidad,
                ValorUnitario: c.ValorUnitario,
                Importe:       c.Importe,
            }
        }
    } else {
        // Si no hay conceptos, añadir uno genérico
        facturaXML.Conceptos.Concepto = []Concepto{
            {
                Descripcion:   "Venta general",
                Cantidad:      1,
                ValorUnitario: factura.Total,
                Importe:       factura.Total,
            },
        }
    }

    // Convertir a XML con formato
    return xml.MarshalIndent(facturaXML, "", "  ")
}

// Ya no necesitamos la función getDireccion porque estamos usando directamente
// el campo DomicilioFiscal en lugar de construir una dirección compleja