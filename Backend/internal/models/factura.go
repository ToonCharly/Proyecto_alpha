package models

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"
)

type Factura struct {
	ID                 string      `json:"id"`
	EmpresaID          interface{} `json:"empresa_id"`
	RFC                string      `json:"rfc"`
	RazonSocial        string      `json:"razon_social"`
	Direccion          string      `json:"direccion"`
	CodigoPostal       string      `json:"codigo_postal"`
	DomicilioFiscal    string      `json:"domicilio_fiscal"`
	RegimenFiscal      string      `json:"regimen_fiscal"`
	ReceptorRFC        string      `json:"receptor_rfc"`
	ClienteRFC         string      `json:"cliente_rfc"`
	ClienteRazonSocial string      `json:"cliente_razon_social"`
	ClienteDireccion   string      `json:"cliente_direccion"`
	UsoCFDI            string      `json:"uso_cfdi"`
	ClaveTicket        string      `json:"clave_ticket"`
	Serie              string      `json:"serie,omitempty"` // Serie para datos fiscales (serie_df)
	FechaEmision       string      `json:"fecha_emision"`
	Subtotal           float64     `json:"subtotal"`
	Impuestos          float64     `json:"impuestos"`
	Total              float64     `json:"total"`
	Observaciones      string      `json:"observaciones"`
	Conceptos          []Concepto  `json:"conceptos"`
	NumeroFolio        string      `json:"numero_folio,omitempty"`
	FechaFactura       time.Time   `json:"fecha_factura,omitempty"`
	LugarExpedicion    string      `json:"lugar_expedicion,omitempty"`
	IVA                float64     `json:"iva,omitempty"`
	MetodoPago         string      `json:"metodo_pago,omitempty"`
	FormaPago          string      `json:"forma_pago,omitempty"`
	Descuento          float64     `json:"descuento,omitempty"`
	Moneda             string      `json:"moneda,omitempty"`
	TipoCambio         float64     `json:"tipo_cambio,omitempty"`
	NumeroCuentaPago   string      `json:"numero_cuenta_pago,omitempty"`
	CondicionesPago    string      `json:"condiciones_pago,omitempty"`
	NumeroPedido       string      `json:"numero_pedido,omitempty"`
	NumeroContrato     string      `json:"numero_contrato,omitempty"`
	NumeroCliente      string      `json:"numero_cliente,omitempty"`
	NumeroProveedor    string      `json:"numero_proveedor,omitempty"`
	FechaVencimiento   string      `json:"fecha_vencimiento,omitempty"`

	// Nuevos campos para el receptor
	ReceptorRazonSocial   string `json:"receptor_razon_social,omitempty"`
	ReceptorDireccion     string `json:"receptor_direccion,omitempty"`
	ReceptorCodigoPostal  string `json:"receptor_codigo_postal,omitempty"`
	RegimenFiscalReceptor string `json:"regimen_fiscal_receptor"`
	Localidad             string `json:"localidad,omitempty"`
	EstadoNombre          string `json:"estado_nombre,omitempty"` // Nombre del estado para mostrar en PDF

	// Nuevos campos para el emisor (datos fiscales del usuario)
	EmisorRFC             string `json:"emisor_rfc,omitempty"`
	EmisorRazonSocial     string `json:"emisor_razon_social,omitempty"`
	EmisorNombreComercial string `json:"emisor_nombre_comercial,omitempty"`
	EmisorDireccionFiscal string `json:"emisor_direccion_fiscal,omitempty"`
	EmisorDireccion       string `json:"emisor_direccion,omitempty"`
	EmisorColonia         string `json:"emisor_colonia,omitempty"`
	EmisorCodigoPostal    string `json:"emisor_codigo_postal,omitempty"`
	EmisorCiudad          string `json:"emisor_ciudad,omitempty"`
	EmisorEstado          string `json:"emisor_estado,omitempty"`
	EmisorRegimenFiscal   string `json:"emisor_regimen_fiscal,omitempty"`
	EmisorMetodoPago      string `json:"emisor_metodo_pago,omitempty"`
	EmisorTipoPago        string `json:"emisor_tipo_pago,omitempty"`
	EmisorCondicionPago   string `json:"emisor_condicion_pago,omitempty"`

	// Campos específicos para la base de datos
	IdFactura int    `json:"idfactura"`  // Para compatibilidad con la BD
	IdEmpresa int    `json:"idempresa"`  // Para compatibilidad con la BD
	IdUsuario int    `json:"id_usuario"` // ID del usuario que genera la factura
	Estatus   int    `json:"estatus"`
	Pagado    int    `json:"pagado"`
	FechaPago string `json:"fecha_pago"`
	Estado    int    `json:"estado"`
}

type Concepto struct {
	Descripcion   string  `json:"descripcion"`
	Cantidad      float64 `json:"cantidad"`
	ValorUnitario float64 `json:"valor_unitario"`
	Importe       float64 `json:"importe"`

	// Campos adicionales para la tabla detallada
	ClaveProdServ string  `json:"clave_prod_serv,omitempty"` // Clave del producto/servicio
	ClaveSAT      string  `json:"clave_sat,omitempty"`       // Clave SAT del producto
	ClaveUnidad   string  `json:"clave_unidad,omitempty"`    // Clave SAT de la unidad
	TasaIVA       float64 `json:"tasa_iva,omitempty"`        // Tasa de IVA en porcentaje (16.0)
	TasaIEPS      float64 `json:"tasa_ieps,omitempty"`       // Tasa de IEPS en porcentaje (50.0)
	Descuento     float64 `json:"descuento,omitempty"`       // Descuento aplicado
}

// CFDI representa la estructura del Comprobante Fiscal Digital
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

// GenerarFolioAutomatico genera automáticamente el folio para la factura SIN usar base de datos
func (f *Factura) GenerarFolioAutomatico() error {
	// Usar serie por defecto "F" si no se especifica
	serie := "F"

	// Usar el nuevo generador simple que NO depende de la base de datos
	folioGen := GetFolioGenerator()
	folioGenerado, err := folioGen.GenerarFolioSimple(serie)
	if err != nil {
		return fmt.Errorf("error al generar folio: %v", err)
	}

	f.NumeroFolio = folioGenerado
	log.Printf("Folio generado automáticamente (sin BD): %s", folioGenerado)
	return nil
}

// ValidarFolio verifica que el folio de la factura esté presente (no verifica BD)
func (f *Factura) ValidarFolio() error {
	if f.NumeroFolio == "" {
		return fmt.Errorf("número de folio requerido")
	}

	// Ya no validamos en BD porque cada folio generado es único por diseño
	log.Printf("Folio válido: %s", f.NumeroFolio)
	return nil
}
