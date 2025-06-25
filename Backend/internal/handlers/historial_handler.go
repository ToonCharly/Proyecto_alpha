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

			facturas, err := models.ObtenerHistorialFacturasPorUsuario(idUsuario)
			if err != nil {
				log.Printf("Error al obtener facturas: %v", err)
				http.Error(w, "Error al obtener facturas", http.StatusInternalServerError)
				return
			}

			// Importante: Devolver directamente el array, no un objeto que lo contenga
			json.NewEncoder(w).Encode(facturas)
			return
		}

		// Manejar solicitudes POST (guardar nueva factura)
		if r.Method == http.MethodPost {
			var factura struct {
				IDUsuario           int     `json:"id_usuario"`
				RFCReceptor         string  `json:"rfc_receptor"`
				RazonSocialReceptor string  `json:"razon_social_receptor"`
				ClaveTicket         string  `json:"clave_ticket"`
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
