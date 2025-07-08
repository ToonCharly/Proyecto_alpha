package handlers

import (
	"fmt"
	"log"

	"carlos/Facts/Backend/internal/db"
	"carlos/Facts/Backend/internal/models"
)

// LlenarDatosEmisor llena automáticamente los datos del emisor desde los datos fiscales activos del usuario
func LlenarDatosEmisor(factura *models.Factura, userID int) error {
	log.Printf("DEBUG - LlenarDatosEmisor iniciando para usuario ID: %d", userID)

	// Intentar obtener datos fiscales del usuario especificado
	datosFiscales, err := obtenerDatosFiscalesUsuario(userID)
	if err != nil {
		log.Printf("INFO - No se encontraron datos fiscales para usuario %d: %v", userID, err)

		// Si no es el usuario administrador, intentar usar datos del administrador (usuario 1) como respaldo
		if userID != 1 {
			log.Printf("INFO - Intentando usar datos del administrador como respaldo...")
			datosFiscales, err = obtenerDatosFiscalesUsuario(1)
			if err != nil {
				log.Printf("WARNING - No se pudieron cargar datos del administrador: %v", err)
				return fmt.Errorf("no hay datos fiscales disponibles")
			}
			log.Printf("SUCCESS - Usando datos del administrador como emisor")
		} else {
			return fmt.Errorf("no hay datos fiscales activos configurados para el usuario administrador")
		}
	}

	log.Printf("DEBUG - Datos fiscales obtenidos exitosamente")

	// Mapear datos fiscales a campos del emisor en la factura
	mapearDatosFiscales(factura, datosFiscales)

	log.Printf("Datos del emisor llenados automáticamente para usuario ID %d", userID)
	log.Printf("Emisor: RFC=%s, RazonSocial=%s, CP=%s, RegimenFiscal=%s",
		factura.EmisorRFC, factura.EmisorRazonSocial, factura.EmisorCodigoPostal, factura.EmisorRegimenFiscal)

	return nil
}

// obtenerDatosFiscalesUsuario obtiene los datos fiscales de un usuario específico desde la base de datos
func obtenerDatosFiscalesUsuario(userID int) (map[string]interface{}, error) {
	log.Printf("DEBUG - Obteniendo datos fiscales desde BD para usuario ID: %d", userID)

	// Usar la función de base de datos existente
	datosFiscales, err := db.ObtenerDatosFiscales(userID)
	if err != nil {
		log.Printf("ERROR - No se pudieron obtener datos fiscales de BD para usuario %d: %v", userID, err)
		return nil, fmt.Errorf("no se encontraron datos fiscales para el usuario %d: %v", userID, err)
	}
	log.Printf("DEBUG - Datos fiscales obtenidos de BD exitosamente")
	log.Printf("DEBUG - Serie obtenida de BD: '%v'", datosFiscales["serie_df"])
	log.Printf("DEBUG - Tipo de serie_df: %T", datosFiscales["serie_df"])
	
	return datosFiscales, nil
}

// mapearDatosFiscales mapea los datos fiscales a los campos del emisor en la factura
func mapearDatosFiscales(factura *models.Factura, datosFiscales map[string]interface{}) {
	if rfc, ok := datosFiscales["rfc"].(string); ok {
		factura.EmisorRFC = rfc
		log.Printf("DEBUG - EmisorRFC asignado: %s", rfc)
	}

	if razonSocial, ok := datosFiscales["razon_social"].(string); ok {
		factura.EmisorRazonSocial = razonSocial
		log.Printf("DEBUG - EmisorRazonSocial asignado: %s", razonSocial)
	}

	if nombreComercial, ok := datosFiscales["nombre_comercial"].(string); ok {
		factura.EmisorNombreComercial = nombreComercial
	}

	if direccionFiscal, ok := datosFiscales["direccion_fiscal"].(string); ok {
		factura.EmisorDireccionFiscal = direccionFiscal
	}

	if direccion, ok := datosFiscales["direccion"].(string); ok {
		factura.EmisorDireccion = direccion
	}

	if colonia, ok := datosFiscales["colonia"].(string); ok {
		factura.EmisorColonia = colonia
	}

	if codigoPostal, ok := datosFiscales["codigo_postal"].(string); ok {
		factura.EmisorCodigoPostal = codigoPostal
	}

	if ciudad, ok := datosFiscales["ciudad"].(string); ok {
		factura.EmisorCiudad = ciudad
	}

	if estado, ok := datosFiscales["estado"].(string); ok {
		factura.EmisorEstado = estado
	}

	if regimenFiscal, ok := datosFiscales["regimen_fiscal"].(string); ok {
		factura.EmisorRegimenFiscal = regimenFiscal
	}

	if metodoPago, ok := datosFiscales["metodo_pago"].(string); ok {
		factura.EmisorMetodoPago = metodoPago
	}

	if tipoPago, ok := datosFiscales["tipo_pago"].(string); ok {
		factura.EmisorTipoPago = tipoPago
	}

	if condicionPago, ok := datosFiscales["condicion_pago"].(string); ok {
		factura.EmisorCondicionPago = condicionPago
	}

	// Agregar serie de datos fiscales
	if serie, ok := datosFiscales["serie_df"].(string); ok && serie != "" {
		factura.Serie = serie
		log.Printf("DEBUG - Serie asignada a factura: %s", serie)
	} else {
		log.Printf("DEBUG - No se encontró serie o está vacía en datos fiscales")
	}
}
