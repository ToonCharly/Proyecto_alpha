package utils

import (
    "encoding/json"
    "net/http"
)

// RespondWithJSON envía una respuesta JSON al cliente.
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
    response, err := json.Marshal(payload)
    if err != nil {
        http.Error(w, "Error al generar la respuesta JSON", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(response)
}