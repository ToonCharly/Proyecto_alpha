package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"carlos/Facts/Backend/internal/db"
)

// RegimenFiscal representa un régimen fiscal.
type RegimenFiscal struct {
	ID          int    `json:"id"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
}

// ObtenerCodigoRegimenFiscal obtiene el código del SAT a partir del ID interno
func ObtenerCodigoRegimenFiscal(id string) (string, error) {
	if id == "" {
		return "", nil
	}

	dbConn, err := db.Connect()
	if err != nil {
		return "", err
	}
	defer dbConn.Close()

	query := "SELECT c_regimenfiscal FROM efac_regimenfiscal WHERE idregimenfiscal = ?"
	var codigo string
	err = dbConn.QueryRow(query, id).Scan(&codigo)
	if err != nil {
		log.Printf("Error al obtener código de régimen fiscal para ID %s: %v", id, err)
		return id, nil // Devolver el ID original si no se encuentra
	}

	log.Printf("DEBUG - Régimen fiscal convertido: ID '%s' -> Código '%s'", id, codigo)
	return codigo, nil
}

// ObtenerNombreEstado obtiene el nombre del estado a partir del ID
func ObtenerNombreEstado(id string) (string, error) {
	if id == "" || id == "0" {
		return "", nil
	}

	dbConn, err := db.Connect()
	if err != nil {
		return "", err
	}
	defer dbConn.Close()

	query := "SELECT estado FROM adm_estados_mex WHERE idestado = ?"
	var nombre string
	err = dbConn.QueryRow(query, id).Scan(&nombre)
	if err != nil {
		log.Printf("Error al obtener nombre de estado para ID %s: %v", id, err)
		return id, nil // Devolver el ID original si no se encuentra
	}

	log.Printf("DEBUG - Estado convertido: ID '%s' -> Nombre '%s'", id, nombre)
	return nombre, nil
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
