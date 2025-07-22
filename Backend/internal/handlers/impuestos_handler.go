package handlers

import (
	"Facts/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// ImpuestosHandler maneja las operaciones CRUD de impuestos
func ImpuestosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetImpuestos(db, w, r)
		case http.MethodPost:
			handleCreateImpuesto(db, w, r)
		case http.MethodPut:
			handleUpdateImpuesto(db, w, r)
		case http.MethodDelete:
			handleDeleteImpuesto(db, w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	}
}

// handleGetImpuestos obtiene impuestos por empresa
func handleGetImpuestos(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idEmpresaStr := r.URL.Query().Get("idempresa")
	if idEmpresaStr == "" {
		http.Error(w, "ID de empresa requerido", http.StatusBadRequest)
		return
	}

	idEmpresa, err := strconv.Atoi(idEmpresaStr)
	if err != nil {
		http.Error(w, "ID de empresa inválido", http.StatusBadRequest)
		return
	}

	impuestos, err := models.GetImpuestosByEmpresa(db, idEmpresa)
	if err != nil {
		log.Printf("Error al obtener impuestos: %v", err)
		http.Error(w, "Error al obtener impuestos", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":   true,
		"impuestos": impuestos,
		"total":     len(impuestos),
	}

	json.NewEncoder(w).Encode(response)
}

// handleCreateImpuesto crea un nuevo impuesto
func handleCreateImpuesto(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var impuesto models.Impuesto
	if err := json.NewDecoder(r.Body).Decode(&impuesto); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	// Validaciones básicas
	if impuesto.IDEmpresa <= 0 {
		http.Error(w, "ID de empresa requerido", http.StatusBadRequest)
		return
	}

	// Validar que al menos tenga una descripción
	if impuesto.Descripcion == "" {
		http.Error(w, "Descripción del impuesto requerida", http.StatusBadRequest)
		return
	}

	// Validar rangos de impuestos
	if impuesto.IVA < 0 || impuesto.IVA > 100 {
		http.Error(w, "IVA debe estar entre 0 y 100", http.StatusBadRequest)
		return
	}

	if impuesto.IEPS1 < 0 || impuesto.IEPS1 > 100 {
		http.Error(w, "IEPS1 debe estar entre 0 y 100", http.StatusBadRequest)
		return
	}

	if impuesto.IEPS2 < 0 || impuesto.IEPS2 > 100 {
		http.Error(w, "IEPS2 debe estar entre 0 y 100", http.StatusBadRequest)
		return
	}

	if impuesto.IEPS3 < 0 || impuesto.IEPS3 > 100 {
		http.Error(w, "IEPS3 debe estar entre 0 y 100", http.StatusBadRequest)
		return
	}

	err := models.CreateImpuesto(db, &impuesto)
	if err != nil {
		log.Printf("Error al crear impuesto: %v", err)
		http.Error(w, "Error al crear impuesto", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"message":  "Impuesto creado exitosamente",
		"impuesto": impuesto,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// handleUpdateImpuesto actualiza un impuesto existente
func handleUpdateImpuesto(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var impuesto models.Impuesto
	if err := json.NewDecoder(r.Body).Decode(&impuesto); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	// Validaciones básicas
	if impuesto.IDIVA <= 0 {
		http.Error(w, "ID de impuesto requerido", http.StatusBadRequest)
		return
	}

	if impuesto.IDEmpresa <= 0 {
		http.Error(w, "ID de empresa requerido", http.StatusBadRequest)
		return
	}

	// Validar que al menos tenga una descripción
	if impuesto.Descripcion == "" {
		http.Error(w, "Descripción del impuesto requerida", http.StatusBadRequest)
		return
	}

	// Validar rangos de impuestos
	if impuesto.IVA < 0 || impuesto.IVA > 100 {
		http.Error(w, "IVA debe estar entre 0 y 100", http.StatusBadRequest)
		return
	}

	if impuesto.IEPS1 < 0 || impuesto.IEPS1 > 100 {
		http.Error(w, "IEPS1 debe estar entre 0 y 100", http.StatusBadRequest)
		return
	}

	if impuesto.IEPS2 < 0 || impuesto.IEPS2 > 100 {
		http.Error(w, "IEPS2 debe estar entre 0 y 100", http.StatusBadRequest)
		return
	}

	if impuesto.IEPS3 < 0 || impuesto.IEPS3 > 100 {
		http.Error(w, "IEPS3 debe estar entre 0 y 100", http.StatusBadRequest)
		return
	}

	err := models.UpdateImpuesto(db, &impuesto)
	if err != nil {
		log.Printf("Error al actualizar impuesto: %v", err)
		http.Error(w, "Error al actualizar impuesto", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Impuesto actualizado exitosamente",
	}

	json.NewEncoder(w).Encode(response)
}

// handleDeleteImpuesto elimina un impuesto
func handleDeleteImpuesto(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idIVAStr := r.URL.Query().Get("idiva")
	idEmpresaStr := r.URL.Query().Get("idempresa")

	if idIVAStr == "" || idEmpresaStr == "" {
		http.Error(w, "ID de impuesto e ID de empresa requeridos", http.StatusBadRequest)
		return
	}

	idIVA, err := strconv.Atoi(idIVAStr)
	if err != nil {
		http.Error(w, "ID de impuesto inválido", http.StatusBadRequest)
		return
	}

	idEmpresa, err := strconv.Atoi(idEmpresaStr)
	if err != nil {
		http.Error(w, "ID de empresa inválido", http.StatusBadRequest)
		return
	}

	err = models.DeleteImpuesto(db, idIVA, idEmpresa)
	if err != nil {
		log.Printf("Error al eliminar impuesto: %v", err)
		http.Error(w, "Error al eliminar impuesto", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Impuesto eliminado exitosamente",
	}

	json.NewEncoder(w).Encode(response)
}

// ProductosConImpuestosHandler obtiene productos con información de impuestos
func ProductosConImpuestosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		idEmpresaStr := r.URL.Query().Get("idempresa")
		if idEmpresaStr == "" {
			http.Error(w, "ID de empresa requerido", http.StatusBadRequest)
			return
		}

		idEmpresa, err := strconv.Atoi(idEmpresaStr)
		if err != nil {
			http.Error(w, "ID de empresa inválido", http.StatusBadRequest)
			return
		}

		productos, err := models.GetProductosConImpuestos(db, idEmpresa)
		if err != nil {
			log.Printf("Error al obtener productos con impuestos: %v", err)
			http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success":   true,
			"productos": productos,
			"total":     len(productos),
		}

		json.NewEncoder(w).Encode(response)
	}
}
