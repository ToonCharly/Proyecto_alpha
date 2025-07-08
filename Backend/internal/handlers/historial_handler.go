package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"carlos/Facts/Backend/internal/models"
)

// HistorialFactura representa los datos para crear un registro en el historial
type HistorialFactura struct {
	IDUsuario           int     `json:"id_usuario"`
	RFCReceptor         string  `json:"rfc_receptor"`
	RazonSocialReceptor string  `json:"razon_social_receptor"`
	ClaveTicket         string  `json:"clave_ticket"`
	Folio               string  `json:"folio"`
	Total               float64 `json:"total"`
	UsoCFDI             string  `json:"uso_cfdi"`
	Observaciones       string  `json:"observaciones"`
	Estado              string  `json:"estado"`
}

// HistorialFacturasHandler maneja las peticiones relacionadas con el historial de facturas
func HistorialFacturasHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Manejar solicitudes GET (obtener facturas)
		if r.Method == http.MethodGet {
			idUsuarioStr := r.URL.Query().Get("id_usuario")
			if idUsuarioStr == "" {
				http.Error(w, "Se requiere id_usuario", http.StatusBadRequest)
				return
			}

			idUsuario, err := strconv.Atoi(idUsuarioStr)
			if err != nil {
				http.Error(w, "id_usuario debe ser un número", http.StatusBadRequest)
				return
			}

			// Obtener parámetros de paginación
			paginaStr := r.URL.Query().Get("pagina")
			limiteStr := r.URL.Query().Get("limite")

			// Valores por defecto
			pagina := 1
			limite := 10

			// Si se proporciona página, usarla
			if paginaStr != "" {
				if p, err := strconv.Atoi(paginaStr); err == nil && p > 0 {
					pagina = p
				}
			}

			// Si se proporciona límite, usarlo (máximo 50)
			if limiteStr != "" {
				if l, err := strconv.Atoi(limiteStr); err == nil && l > 0 && l <= 50 {
					limite = l
				}
			}

			log.Printf("Obteniendo historial de facturas - Usuario: %d, Página: %d, Límite: %d", idUsuario, pagina, limite)

			facturas, totalFacturas, err := models.ObtenerHistorialFacturasPorUsuarioConPaginacion(idUsuario, pagina, limite)
			if err != nil {
				log.Printf("Error al obtener facturas: %v", err)
				http.Error(w, "Error al obtener facturas", http.StatusInternalServerError)
				return
			}

			// Calcular información de paginación
			totalPaginas := (totalFacturas + limite - 1) / limite
			tieneSiguiente := pagina < totalPaginas
			tieneAnterior := pagina > 1

			// Devolver respuesta con información de paginación
			response := map[string]interface{}{
				"facturas": facturas,
				"paginacion": map[string]interface{}{
					"pagina_actual":   pagina,
					"limite":          limite,
					"total_facturas":  totalFacturas,
					"total_paginas":   totalPaginas,
					"tiene_siguiente": tieneSiguiente,
					"tiene_anterior":  tieneAnterior,
				},
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		// Manejar solicitudes POST (guardar nueva factura)
		if r.Method == http.MethodPost {
			var factura struct {
				IDUsuario           int     `json:"id_usuario"`
				RFCReceptor         string  `json:"rfc_receptor"`
				RazonSocialReceptor string  `json:"razon_social_receptor"`
				ClaveTicket         string  `json:"clave_ticket"`
				Folio               string  `json:"folio"`
				Total               float64 `json:"total"`
				UsoCFDI             string  `json:"uso_cfdi"`
				Observaciones       string  `json:"observaciones"`
			}

			if err := json.NewDecoder(r.Body).Decode(&factura); err != nil {
				http.Error(w, "Error al decodificar datos", http.StatusBadRequest)
				return
			}

			id, err := models.InsertarHistorialFactura(
				factura.IDUsuario,
				factura.RFCReceptor,
				factura.RazonSocialReceptor,
				factura.ClaveTicket,
				factura.Folio,
				factura.Total,
				factura.UsoCFDI,
				factura.Observaciones,
			)

			if err != nil {
				log.Printf("Error al insertar en historial: %v", err)
				http.Error(w, "Error al guardar factura en historial", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":      id,
				"mensaje": "Factura guardada correctamente en el historial",
			})
			return
		}

		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// BuscarHistorialFacturasHandler maneja las búsquedas en el historial de facturas
func BuscarHistorialFacturasHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Obtener parámetros de la URL
		idUsuarioStr := r.URL.Query().Get("id_usuario")
		if idUsuarioStr == "" {
			log.Printf("Error: Se requiere id_usuario")
			http.Error(w, "Se requiere id_usuario", http.StatusBadRequest)
			return
		}

		idUsuario, err := strconv.Atoi(idUsuarioStr)
		if err != nil {
			log.Printf("Error: id_usuario debe ser un número: %v", err)
			http.Error(w, "id_usuario debe ser un número", http.StatusBadRequest)
			return
		}

		// Obtener criterios de búsqueda opcionales
		folio := r.URL.Query().Get("folio")
		rfcReceptor := r.URL.Query().Get("rfc_receptor")
		razonSocial := r.URL.Query().Get("razon_social_receptor")

		// Obtener parámetros de paginación
		paginaStr := r.URL.Query().Get("pagina")
		limiteStr := r.URL.Query().Get("limite")

		// Valores por defecto para paginación
		pagina := 1
		limite := 10

		// Si se proporciona página, usarla
		if paginaStr != "" {
			if p, err := strconv.Atoi(paginaStr); err == nil && p > 0 {
				pagina = p
			}
		}

		// Si se proporciona límite, usarlo (máximo 50)
		if limiteStr != "" {
			if l, err := strconv.Atoi(limiteStr); err == nil && l > 0 && l <= 50 {
				limite = l
			}
		}

		log.Printf("Búsqueda de facturas - Usuario: %d, Folio: '%s', RFC: '%s', Razón Social: '%s', Página: %d, Límite: %d",
			idUsuario, folio, rfcReceptor, razonSocial, pagina, limite)

		// Si no hay criterios de búsqueda, devolver todas las facturas del usuario con paginación
		if folio == "" && rfcReceptor == "" && razonSocial == "" {
			log.Printf("Sin criterios de búsqueda, obteniendo todas las facturas del usuario %d con paginación", idUsuario)
			facturas, totalFacturas, err := models.ObtenerHistorialFacturasPorUsuarioConPaginacion(idUsuario, pagina, limite)
			if err != nil {
				log.Printf("Error al obtener facturas: %v", err)
				http.Error(w, "Error al obtener facturas", http.StatusInternalServerError)
				return
			}

			// Calcular información de paginación
			totalPaginas := (totalFacturas + limite - 1) / limite
			tieneSiguiente := pagina < totalPaginas
			tieneAnterior := pagina > 1

			response := map[string]interface{}{
				"facturas": facturas,
				"paginacion": map[string]interface{}{
					"pagina_actual":   pagina,
					"limite":          limite,
					"total_facturas":  totalFacturas,
					"total_paginas":   totalPaginas,
					"tiene_siguiente": tieneSiguiente,
					"tiene_anterior":  tieneAnterior,
				},
			}

			log.Printf("Se encontraron %d facturas de %d totales", len(facturas), totalFacturas)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Buscar facturas con criterios específicos y paginación
		log.Printf("Buscando facturas con criterios específicos y paginación")
		facturas, totalFacturas, err := models.BuscarHistorialFacturasConPaginacion(idUsuario, folio, rfcReceptor, razonSocial, pagina, limite)
		if err != nil {
			log.Printf("Error al buscar facturas: %v", err)
			http.Error(w, "Error al buscar facturas", http.StatusInternalServerError)
			return
		}

		// Calcular información de paginación
		totalPaginas := (totalFacturas + limite - 1) / limite
		tieneSiguiente := pagina < totalPaginas
		tieneAnterior := pagina > 1

		response := map[string]interface{}{
			"facturas": facturas,
			"paginacion": map[string]interface{}{
				"pagina_actual":   pagina,
				"limite":          limite,
				"total_facturas":  totalFacturas,
				"total_paginas":   totalPaginas,
				"tiene_siguiente": tieneSiguiente,
				"tiene_anterior":  tieneAnterior,
			},
		}

		log.Printf("Búsqueda completada, se encontraron %d facturas de %d totales", len(facturas), totalFacturas)
		json.NewEncoder(w).Encode(response)
	}
}
