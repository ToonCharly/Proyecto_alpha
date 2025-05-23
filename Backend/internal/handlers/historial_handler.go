package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

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

// HistorialFacturasHandler obtiene el historial de facturas
func HistorialFacturasHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
            return
        }

        query := `SELECT idfactura, idempresa, rfc, razon_social, subtotal, impuestos, total, fecha_emision, estatus FROM adm_efacturas ORDER BY fecha_emision DESC`
        rows, err := db.Query(query)
        if err != nil {
            log.Printf("Error al consultar el historial de facturas: %v", err)
            http.Error(w, "Error al consultar el historial de facturas", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var facturas []models.Factura
        for rows.Next() {
            var factura models.Factura
            err := rows.Scan(&factura.IdFactura, &factura.IdEmpresa, &factura.RFC, &factura.RazonSocial, &factura.Subtotal, &factura.Impuestos, &factura.Total, &factura.FechaEmision, &factura.Estatus)
            if err != nil {
                log.Printf("Error al escanear los resultados: %v", err)
                http.Error(w, "Error al procesar los datos", http.StatusInternalServerError)
                return
            }
            facturas = append(facturas, factura)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "facturas": facturas,
            "total":    len(facturas),
        })
    }
}

// CreateHistorialFacturaHandler maneja la creación de registros en el historial
func CreateHistorialFacturaHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
            return
        }

        var historial HistorialFactura
        if err := json.NewDecoder(r.Body).Decode(&historial); err != nil {
            log.Printf("Error al decodificar la petición: %v", err)
            http.Error(w, "Error al procesar la petición", http.StatusBadRequest)
            return
        }

        // Antes de ejecutar la consulta, verifica si el usuario existe:
        var exists bool
        err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuarios WHERE id = ?)", historial.IDUsuario).Scan(&exists)
        if err != nil {
            log.Printf("Error al verificar si el usuario existe: %v", err)
        }
        if !exists {
            log.Printf("ADVERTENCIA: El usuario con ID %d no existe en la base de datos", historial.IDUsuario)
            http.Error(w, "El usuario no existe", http.StatusBadRequest)
            return
        }

        // Añade logs detallados antes de la ejecución
        log.Printf("Intentando guardar historial para usuario ID: %d", historial.IDUsuario)

        // Consulta de inserción con todos los campos
        query := `INSERT INTO historial_facturas (id_usuario, rfc_receptor, razon_social_receptor, clave_ticket, total, uso_cfdi, observaciones, estado) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

        // Ejecutar la consulta en una sola línea para evitar problemas de sintaxis
        _, err = db.Exec(query,
            historial.IDUsuario,
            historial.RFCReceptor,
            historial.RazonSocialReceptor,
            historial.ClaveTicket,
            historial.Total,
            historial.UsoCFDI,
            historial.Observaciones,
            historial.Estado)

        if err != nil {
            log.Printf("Error al insertar en el historial: %v", err)
            http.Error(w, "Error al guardar en el historial", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{
            "message": "Registro guardado correctamente en el historial",
        })
    }
}