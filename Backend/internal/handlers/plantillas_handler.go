package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

// BuscarPlantillasHandler maneja la búsqueda de plantillas disponibles.
func BuscarPlantillasHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    plantillasDir := "./templates/facturas"

    if _, err := os.Stat(plantillasDir); os.IsNotExist(err) {
        if err := os.MkdirAll(plantillasDir, 0755); err != nil {
            log.Printf("Error al crear directorio de plantillas: %v", err)
            http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
            return
        }
    }

    files, err := os.ReadDir(plantillasDir)
    if err != nil {
        log.Printf("Error al leer directorio de plantillas: %v", err)
        http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        return
    }

    type PlantillaInfo struct {
        Nombre string `json:"nombre"`
        Ruta   string `json:"ruta"`
    }

    var plantillas []PlantillaInfo
    for _, file := range files {
        if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".docx") {
            plantillas = append(plantillas, PlantillaInfo{
                Nombre: file.Name(),
                Ruta:   filepath.Join(plantillasDir, file.Name()),
            })
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "plantillas": plantillas,
        "total":      len(plantillas),
    })
}