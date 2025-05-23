package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "carlos/Facts/Backend/internal/db"
)

// RegimenFiscal representa un régimen fiscal.
type RegimenFiscal struct {
    ID          int    `json:"id"`          // idregimenfiscal
    Codigo      string `json:"codigo"`      // c_regimenfiscal
    Descripcion string `json:"descripcion"` // descripcion
}

// GetRegimenesFiscales maneja la solicitud para obtener los regímenes fiscales.
func GetRegimenesFiscales(w http.ResponseWriter, r *http.Request) {
    log.Println("Handler GetRegimenesFiscales ejecutado")

    if r.Method != http.MethodGet {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    dbConn, err := db.Connect()
    if err != nil {
        log.Printf("Error al conectar a la base de datos: %v", err)
        http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
        return
    }
    defer dbConn.Close()

    // Actualiza la consulta SQL con los nombres correctos de las columnas
    query := "SELECT idregimenfiscal, c_regimenfiscal, descripcion FROM efac_regimenfiscal"
    log.Printf("Ejecutando consulta: %s", query)

    rows, err := dbConn.Query(query)
    if err != nil {
        log.Printf("Error al ejecutar la consulta: %v", err)
        http.Error(w, "Error al consultar la base de datos", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var regimenes []RegimenFiscal
    for rows.Next() {
        var regimen RegimenFiscal
        if err := rows.Scan(&regimen.ID, &regimen.Codigo, &regimen.Descripcion); err != nil {
            log.Printf("Error al procesar los datos: %v", err)
            http.Error(w, "Error al procesar los datos", http.StatusInternalServerError)
            return
        }
        regimenes = append(regimenes, regimen)
    }

    if len(regimenes) == 0 {
        log.Println("No se encontraron regímenes fiscales")
        http.Error(w, "No se encontraron regímenes fiscales", http.StatusNotFound)
        return
    }

    log.Printf("Regímenes fiscales obtenidos: %+v", regimenes)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(regimenes)
}