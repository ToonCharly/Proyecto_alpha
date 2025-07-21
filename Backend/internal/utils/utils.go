package utils

import (
    "os"
    "encoding/json"
    "log"
    "net/http"
)

// SaveFile guarda un archivo en disco en la ruta especificada.
func SaveFile(path string, data []byte) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = f.Write(data)
    return err
}

// RespondWithJSON envía una respuesta JSON al cliente
func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Printf("Error al serializar JSON: %v", err)
        http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        return
    }

    // Agregar log para depuración
    log.Printf("Respuesta JSON (primeros 100 bytes): %s", string(jsonData[:min(len(jsonData), 100)]))

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write(jsonData)
}

// Función auxiliar min
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}