package models

import (
	"encoding/xml"
	"time"
)

type Factura struct {
	IdFactura             int64      `json:"idfactura,omitempty"`
	IdEmpresa             int64      `json:"idempresa,omitempty"`
	RFC                   string     `json:"rfc"`
	RazonSocial           string     `json:"razon_social"`
	Subtotal              float64    `json:"subtotal,omitempty"`
	Impuestos             float64    `json:"impuestos,omitempty"`
	Estatus               string     `json:"estatus,omitempty"`
	Pagado                string     `json:"pagado,omitempty"`
	FechaPago             string     `json:"fecha_pago,omitempty"`
	Direccion             string     `json:"direccion,omitempty"`
	CodigoPostal          string     `json:"codigo_postal,omitempty"`
	Pais                  string     `json:"pais,omitempty"`
	Estado                string     `json:"estado,omitempty"`
	Localidad             string     `json:"localidad,omitempty"`
	Municipio             string     `json:"municipio,omitempty"`
	Colonia               string     `json:"colonia,omitempty"`
	Observaciones         string     `json:"observaciones,omitempty"`
	UsoCFDI               string     `json:"uso_cfdi,omitempty"`
	RegimenFiscal         string     `json:"regimen_fiscal,omitempty"`
	ClaveTicket           string     `json:"clave_ticket,omitempty"`
	Total                 float64    `json:"total,omitempty"`
	FechaEmision          string     `json:"fecha_emision,omitempty"`
	MetodoPago            string     `json:"metodo_pago,omitempty"`
	FormaPago             string     `json:"forma_pago,omitempty"`
	Descuento             float64    `json:"descuento,omitempty"`
	Moneda                string     `json:"moneda,omitempty"`
	TipoCambio            float64    `json:"tipo_cambio,omitempty"`
	NumeroCuentaPago      string     `json:"numero_cuenta_pago,omitempty"`
	CondicionesPago       string     `json:"condiciones_pago,omitempty"`
	NumeroPedido          string     `json:"numero_pedido,omitempty"`
	NumeroContrato        string     `json:"numero_contrato,omitempty"`
	NumeroCliente         string     `json:"numero_cliente,omitempty"`
	NumeroProveedor       string     `json:"numero_proveedor,omitempty"`
	FechaVencimiento      string     `json:"fecha_vencimiento,omitempty"`
	Conceptos             []Concepto `json:"conceptos,omitempty"`
	NumeroFolio           string     `json:"numero_folio,omitempty"`
	FechaFactura          time.Time  `json:"fecha_factura,omitempty"`
	LugarExpedicion       string     `json:"lugar_expedicion,omitempty"`
	ReceptorRFC           string     `json:"receptor_rfc,omitempty"`
	DomicilioFiscal       string     `json:"domicilio_fiscal,omitempty"`
	RegimenFiscalReceptor string     `json:"regimen_fiscal_receptor,omitempty"`
	IVA                   float64    `json:"iva,omitempty"`

	// Nuevos campos para el receptor
	ReceptorRazonSocial  string `json:"receptor_razon_social,omitempty"`
	ReceptorDireccion    string `json:"receptor_direccion,omitempty"`
	ReceptorCodigoPostal string `json:"receptor_codigo_postal,omitempty"`
}

type Concepto struct {
	ClaveProductoServicio string  `json:"clave_producto_servicio,omitempty"`
	ClaveUnidad           string  `json:"clave_unidad,omitempty"`
	Descripcion           string  `json:"descripcion,omitempty"`
	Cantidad              float64 `json:"cantidad,omitempty"`
	ValorUnitario         float64 `json:"valor_unitario,omitempty"`
	Importe               float64 `json:"importe,omitempty"`
	Descuento             float64 `json:"descuento,omitempty"`
}

// CFDI representa la estructura del Comprobante Fiscal Digital por Internet
type CFDI struct {
	XMLName           xml.Name `xml:"cfdi:Comprobante"`
	XmlnsCfdi         string   `xml:"xmlns:cfdi,attr"`
	Version           string   `xml:"Version,attr"`
	Serie             string   `xml:"Serie,attr"`
	Folio             string   `xml:"Folio,attr"`
	Fecha             string   `xml:"Fecha,attr"`
	SubTotal          string   `xml:"SubTotal,attr"`
	Total             string   `xml:"Total,attr"`
	Moneda            string   `xml:"Moneda,attr"`
	TipoCambio        string   `xml:"TipoCambio,attr,omitempty"`
	LugarExpedicion   string   `xml:"LugarExpedicion,attr"`
	TipoDeComprobante string   `xml:"TipoDeComprobante,attr"`
	MetodoPago        string   `xml:"MetodoPago,attr"`
	FormaPago         string   `xml:"FormaPago,attr"`
	CondicionesPago   string   `xml:"CondicionesDePago,attr,omitempty"`
	Descuento         string   `xml:"Descuento,attr,omitempty"`
	ClaveTicket       string   `xml:"ClaveTicket,attr,omitempty"`
	Emisor            struct {
		RFC           string `xml:"Rfc,attr"`
		Nombre        string `xml:"Nombre,attr"`
		RegimenFiscal string `xml:"RegimenFiscal,attr"`
	} `xml:"cfdi:Emisor"`
	Receptor struct {
		RFC                   string `xml:"Rfc,attr"`
		Nombre                string `xml:"Nombre,attr"`
		UsoCFDI               string `xml:"UsoCFDI,attr"`
		DomicilioFiscal       string `xml:"DomicilioFiscal,attr,omitempty"`
		RegimenFiscalReceptor string `xml:"RegimenFiscalReceptor,attr,omitempty"`
	} `xml:"cfdi:Receptor"`
	Conceptos []struct {
		ClaveProdServ    string `xml:"ClaveProdServ,attr"`
		NoIdentificacion string `xml:"NoIdentificacion,attr,omitempty"`
		Cantidad         string `xml:"Cantidad,attr"`
		ClaveUnidad      string `xml:"ClaveUnidad,attr"`
		Unidad           string `xml:"Unidad,attr,omitempty"`
		Descripcion      string `xml:"Descripcion,attr"`
		ValorUnitario    string `xml:"ValorUnitario,attr"`
		Importe          string `xml:"Importe,attr"`
		Descuento        string `xml:"Descuento,attr,omitempty"`
	} `xml:"cfdi:Conceptos>cfdi:Concepto"`
	Impuestos struct {
		TotalImpuestosTrasladados string `xml:"TotalImpuestosTrasladados,attr,omitempty"`
		Traslados                 []struct {
			Impuesto   string `xml:"Impuesto,attr"`
			TipoFactor string `xml:"TipoFactor,attr"`
			TasaOCuota string `xml:"TasaOCuota,attr"`
			Importe    string `xml:"Importe,attr"`
		} `xml:"cfdi:Traslados>cfdi:Traslado"`
	} `xml:"cfdi:Impuestos,omitempty"`
}
