package utils

import (
    "log"
    "net/http"
)

func RespondWithError(w http.ResponseWriter, message string) {
    log.Println("Error:", message)
    http.Error(w, message, http.StatusBadRequest)
}