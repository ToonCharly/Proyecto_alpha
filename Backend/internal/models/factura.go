package models

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"
)

// Alias para compatibilidad con handlers y main.go
type FacturaCFDI = Factura

type Factura struct {
	ID                 string      `json:"id"`
	EmpresaID          interface{} `json:"empresa_id"`
	RFC                string      `json:"rfc"`
	RazonSocial        string      `json:"razon_social"`
	Direccion          string      `json:"direccion"`
	CodigoPostal       string      `json:"codigo_postal"`
	DomicilioFiscal    string      `json:"domicilio_fiscal"`
	RegimenFiscal      string      `json:"regimen_fiscal"`
	UUID               string      `json:"uuid"`
	NoCertificado      string      `json:"no_certificado"`
	Certificado        string      `json:"certificado,omitempty"`
	ReceptorRFC        string      `json:"receptor_rfc"`
	EmpresaRFC         string
	ClienteRFC         string     `json:"cliente_rfc"`
	ClienteRazonSocial string     `json:"cliente_razon_social"`
	ClienteDireccion   string     `json:"cliente_direccion"`
	UsoCFDI            string     `json:"uso_cfdi"`
	ClaveTicket        string     `json:"clave_ticket"`
	Serie              string     `json:"serie,omitempty"` // Serie para datos fiscales (serie_df)
	FechaEmision       string     `json:"fecha_emision"`
	Subtotal           float64    `json:"subtotal"`
	Impuestos          float64    `json:"impuestos"`
	Total              float64    `json:"total"`
	Observaciones      string     `json:"observaciones"`
	Conceptos          []Concepto `json:"conceptos"`
	NumeroFolio        string     `json:"numero_folio,omitempty"`
	FechaFactura       time.Time  `json:"fecha_factura,omitempty"`
	LugarExpedicion    string     `json:"lugar_expedicion,omitempty"`
	IVA                float64    `json:"iva,omitempty"`
	MetodoPago         string     `json:"metodo_pago,omitempty"`
	FormaPago          string     `json:"forma_pago,omitempty"`
	Descuento          float64    `json:"descuento,omitempty"`
	Moneda             string     `json:"moneda,omitempty"`
	TipoCambio         float64    `json:"tipo_cambio,omitempty"`
	NumeroCuentaPago   string     `json:"numero_cuenta_pago,omitempty"`
	CondicionesPago    string     `json:"condiciones_pago,omitempty"`
	NumeroPedido       string     `json:"numero_pedido,omitempty"`
	NumeroContrato     string     `json:"numero_contrato,omitempty"`
	NumeroCliente      string     `json:"numero_cliente,omitempty"`
	NumeroProveedor    string     `json:"numero_proveedor,omitempty"`
	FechaVencimiento   string     `json:"fecha_vencimiento,omitempty"`

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

	// Campos para la firma digital CFDI
	KeyPath   string `json:"key_path"`             // Ruta al archivo .key del CSD
	ClaveCSD  string `json:"clave_csd"`            // Contraseña de la llave privada del CSD
	CerPath   string `json:"cer_path"`             // Ruta al archivo .cer del CSD
	CerBase64 string `json:"cer_base64,omitempty"` // Certificado en base64 (opcional, para frontend)
	Timbre    *TimbreFiscalDigital

	// Estado de la generación/timbrado de la factura
	EstatusFac string `json:"estatus_fac"` // F: fallo, P: pendiente, T: timbrando, G: generado
	LogError   string `json:"log_error"`   // Mensaje de error si ocurre
}

// Configuración del PAC y CSD para timbrado
type PACConfig struct {
	UsuarioPAC string // RFC del emisor (usuario PAC)
	ClavePAC   string // Clave del PAC
	Produccion bool   // Modo producción
	CerPath    string // Ruta al .cer
	KeyPath    string // Ruta al .key
	ClaveCSD   string // Contraseña del CSD
	Endpoint   string // URL del PAC
}

// Proceso principal de timbrado CFDI
func (f *Factura) TimbrarCFDI(pac PACConfig) error {
	// 1. Generar el XML CFDI (solo estructura, no firmado)
	xmlCFDI, err := f.GenerarXMLCFDI()
	if err != nil {
		f.LogError = "Error generando XML: " + err.Error()
		return err
	}

	// 2. Firmar el XML CFDI (usar CSD)
	signedXML, err := FirmarXMLCFDI(xmlCFDI, pac.CerPath, pac.KeyPath, pac.ClaveCSD)
	if err != nil {
		f.LogError = "Error firmando XML: " + err.Error()
		return err
	}

	// 3. Enviar el XML firmado al PAC
	timbradoXML, err := EnviarXMLAlPAC(signedXML, pac)
	if err != nil {
		f.LogError = "Error enviando al PAC: " + err.Error()
		return err
	}

	// 4. Extraer el timbre fiscal digital del XML timbrado
	timbre, err := ExtraerTimbreFiscal(timbradoXML)
	if err != nil {
		f.LogError = "Error extrayendo timbre fiscal: " + err.Error()
		return err
	}
	f.Timbre = timbre
	f.EstatusFac = "T" // Timbrado exitoso
	return nil
}

// Genera el XML CFDI a partir de la estructura Factura
func (f *Factura) GenerarXMLCFDI() (string, error) {
	cfdi := CFDI{
		XmlnsCfdi:         "http://www.sat.gob.mx/cfd/4",
		Version:           "4.0",
		Serie:             f.Serie,
		Folio:             f.NumeroFolio,
		Fecha:             f.FechaEmision,
		SubTotal:          fmt.Sprintf("%.2f", f.Subtotal),
		Total:             fmt.Sprintf("%.2f", f.Total),
		Moneda:            f.Moneda,
		TipoCambio:        fmt.Sprintf("%.2f", f.TipoCambio),
		LugarExpedicion:   f.LugarExpedicion,
		TipoDeComprobante: "I",
		MetodoPago:        f.MetodoPago,
		FormaPago:         f.FormaPago,
		CondicionesPago:   f.CondicionesPago,
		Descuento:         fmt.Sprintf("%.2f", f.Descuento),
		ClaveTicket:       f.ClaveTicket,
		Emisor: struct {
			RFC           string `xml:"Rfc,attr"`
			Nombre        string `xml:"Nombre,attr"`
			RegimenFiscal string `xml:"RegimenFiscal,attr"`
		}{
			RFC:           f.EmisorRFC,
			Nombre:        f.EmisorRazonSocial,
			RegimenFiscal: f.EmisorRegimenFiscal,
		},
		Receptor: struct {
			RFC                   string `xml:"Rfc,attr"`
			Nombre                string `xml:"Nombre,attr"`
			UsoCFDI               string `xml:"UsoCFDI,attr"`
			DomicilioFiscal       string `xml:"DomicilioFiscal,attr,omitempty"`
			RegimenFiscalReceptor string `xml:"RegimenFiscalReceptor,attr,omitempty"`
		}{
			RFC:                   f.ReceptorRFC,
			Nombre:                f.ReceptorRazonSocial,
			UsoCFDI:               f.UsoCFDI,
			DomicilioFiscal:       f.ReceptorCodigoPostal,
			RegimenFiscalReceptor: f.RegimenFiscalReceptor,
		},
	}
	// Conceptos
	for _, c := range f.Conceptos {
		cfdi.Conceptos = append(cfdi.Conceptos, struct {
			ClaveProdServ    string `xml:"ClaveProdServ,attr"`
			NoIdentificacion string `xml:"NoIdentificacion,attr,omitempty"`
			Cantidad         string `xml:"Cantidad,attr"`
			ClaveUnidad      string `xml:"ClaveUnidad,attr"`
			Unidad           string `xml:"Unidad,attr,omitempty"`
			Descripcion      string `xml:"Descripcion,attr"`
			ValorUnitario    string `xml:"ValorUnitario,attr"`
			Importe          string `xml:"Importe,attr"`
			Descuento        string `xml:"Descuento,attr,omitempty"`
		}{
			ClaveProdServ:    c.ClaveProdServ,
			NoIdentificacion: "",
			Cantidad:         fmt.Sprintf("%.2f", c.Cantidad),
			ClaveUnidad:      c.ClaveUnidad,
			Unidad:           "",
			Descripcion:      c.Descripcion,
			ValorUnitario:    fmt.Sprintf("%.2f", c.ValorUnitario),
			Importe:          fmt.Sprintf("%.2f", c.Importe),
			Descuento:        fmt.Sprintf("%.2f", c.Descuento),
		})
	}
	output, err := xml.MarshalIndent(cfdi, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// Firmar el XML CFDI usando el CSD (.cer, .key, clave)
func FirmarXMLCFDI(xmlCFDI, cerPath, keyPath, claveCSD string) (string, error) {
	// Aquí deberías implementar la firma digital del XML usando el CSD
	// Puedes usar una librería externa o llamar a un ejecutable que firme el XML
	// Por ahora, se retorna el XML sin firmar (solo para ejemplo)
	return xmlCFDI, nil
}

// Enviar el XML firmado al PAC y recibir el XML timbrado
func EnviarXMLAlPAC(xmlFirmado string, pac PACConfig) (string, error) {
	// Implementa el envío HTTP POST al endpoint del PAC
	// El cuerpo debe incluir el XML firmado y los datos de autenticación
	// Ejemplo genérico:
	/*
		resp, err := http.Post(pac.Endpoint, "application/xml", strings.NewReader(xmlFirmado))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	*/
	return xmlFirmado, nil // Solo ejemplo, aquí deberías poner la respuesta real del PAC
}

// Extraer el timbre fiscal digital del XML timbrado
func ExtraerTimbreFiscal(xmlTimbrado string) (*TimbreFiscalDigital, error) {
	// Aquí deberías parsear el XML timbrado y extraer el nodo TimbreFiscalDigital
	// Ejemplo genérico:
	var tfd TimbreFiscalDigital
	// ...parsear xmlTimbrado y llenar tfd...
	tfd.UUID = "EJEMPLO-UUID"
	tfd.FechaTimbrado = time.Now().Format("2006-01-02T15:04:05")
	tfd.RfcProvCertif = "LSO1306189R5"
	tfd.SelloCFD = "EJEMPLO-SELLO-CFD"
	tfd.NoCertificadoSAT = "EJEMPLO-CERT-SAT"
	tfd.SelloSAT = "EJEMPLO-SELLO-SAT"
	return &tfd, nil
}

type TimbreFiscalDigital struct {
	UUID             string
	FechaTimbrado    string
	RfcProvCertif    string
	SelloCFD         string
	NoCertificadoSAT string
	SelloSAT         string
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
