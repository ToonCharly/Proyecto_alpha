package handlers

import (
	"net/http"
)

// Registra los endpoints principales en el router
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/timbrar", TimbrarCFDIHandler)
}
