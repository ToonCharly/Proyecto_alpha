package utils

import "net/http"

// EnableCors es un middleware que agrega los encabezados CORS apropiados
func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el origen de la solicitud
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Valor por defecto para desarrollo local
			origin = "http://localhost:5173"
		}

		// Usar origen específico en lugar de comodín *
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Añadir este header cuando se usan credenciales
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Manejar preflight OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
