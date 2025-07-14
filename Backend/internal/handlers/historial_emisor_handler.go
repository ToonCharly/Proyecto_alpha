package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"carlos/Facts/Backend/internal/db"
)

type FacturaEmisorHistorial struct {
	ID         int     `json:"id"`
	Fecha      string  `json:"fecha_emision"`
	Folio      string  `json:"folio"`
	Total      float64 `json:"total"`
	Cliente    string  `json:"razon_social_receptor"`
	RFCCliente string  `json:"rfc_receptor"`
	Usuario    struct {
		ID     int    `json:"id"`
		Nombre string `json:"nombre"`
		Email  string `json:"email"`
	} `json:"usuario"`
}

func HistorialEmisorHandler(w http.ResponseWriter, r *http.Request) {
	idEmpresaStr := r.URL.Query().Get("id_empresa")
	if idEmpresaStr == "" {
		http.Error(w, "Falta parámetro id_empresa", http.StatusBadRequest)
		return
	}
	idEmpresa, err := strconv.Atoi(idEmpresaStr)
	if err != nil {
		http.Error(w, "id_empresa inválido", http.StatusBadRequest)
		return
	}
	dbconn := db.GetDB()
	query := `
		SELECT f.id, f.fecha_emision, f.folio, f.total, f.razon_social_receptor, f.rfc_receptor,
		       u.id, u.nombre, u.email
		FROM facturas f
		JOIN usuarios u ON f.id_usuario = u.id
		WHERE f.id_empresa = ?
		ORDER BY f.fecha_emision DESC`
	rows, err := dbconn.Query(query, idEmpresa)
	if err != nil {
		http.Error(w, "Error al consultar historial", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var historial []FacturaEmisorHistorial
	for rows.Next() {
		var f FacturaEmisorHistorial
		var userID int
		var userNombre, userEmail string
		if err := rows.Scan(&f.ID, &f.Fecha, &f.Folio, &f.Total, &f.Cliente, &f.RFCCliente, &userID, &userNombre, &userEmail); err != nil {
			continue
		}
		f.Usuario.ID = userID
		f.Usuario.Nombre = userNombre
		f.Usuario.Email = userEmail
		historial = append(historial, f)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(historial)
}