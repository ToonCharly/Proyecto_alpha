package handlers

import (
	"log"
	"net/http"

	"carlos/Facts/Backend/internal/db"
	"carlos/Facts/Backend/internal/utils"
)

// DebugVentasHandler muestra información de la tabla ventas_det
func DebugVentasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	database := db.GetDB()

	// Obtener las últimas 10 series de ventas
	query := `
		SELECT DISTINCT serie, COUNT(*) as productos
		FROM ventas_det 
		GROUP BY serie 
		ORDER BY serie DESC 
		LIMIT 10
	`

	rows, err := database.Query(query)
	if err != nil {
		log.Printf("Error al consultar ventas_det: %v", err)
		http.Error(w, "Error al consultar ventas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type SerieInfo struct {
		Serie     string `json:"serie"`
		Productos int    `json:"productos"`
	}

	var series []SerieInfo
	for rows.Next() {
		var info SerieInfo
		err := rows.Scan(&info.Serie, &info.Productos)
		if err != nil {
			log.Printf("Error al escanear serie: %v", err)
			continue
		}
		series = append(series, info)
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Series encontradas en ventas_det",
		"series":  series,
		"total":   len(series),
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}
