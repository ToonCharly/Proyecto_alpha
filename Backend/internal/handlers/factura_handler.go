package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "carlos/Facts/Backend/internal/models"
)

func BuscarFactura(db *sql.DB, w http.ResponseWriter, criterio string) {
    log.Println("Criterio recibido:", criterio)
    query := `SELECT idfactura, idempresa, rfc, razon_social, subtotal, impuestos, estatus, pagado, fecha_pago FROM adm_efacturas WHERE rfc LIKE ? OR razon_social LIKE ? LIMIT 1`
    likeCriterio := "%" + criterio + "%"
    row := db.QueryRow(query, likeCriterio, likeCriterio)

    var f models.Factura
    err := row.Scan(&f.IdFactura, &f.IdEmpresa, &f.RFC, &f.RazonSocial, &f.Subtotal, &f.Impuestos, &f.Estatus, &f.Pagado, &f.FechaPago)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        if err == sql.ErrNoRows {
            json.NewEncoder(w).Encode(map[string]string{"error": "Empresa no encontrada"})
        } else {
            json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar la base de datos"})
        }
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(f)
}