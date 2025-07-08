package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// LogoPlantillaHandler maneja las operaciones relacionadas con logos de plantillas
type LogoPlantillaHandler struct{}

// NewLogoPlantillaHandler crea una nueva instancia del handler
func NewLogoPlantillaHandler() *LogoPlantillaHandler {
	return &LogoPlantillaHandler{}
}

// SubirLogo maneja la subida de logos para plantillas
func (h *LogoPlantillaHandler) SubirLogo(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parsear el formulario multipart
	err := r.ParseMultipartForm(2 << 20) // 2MB límite
	if err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Obtener el archivo del logo
	file, header, err := r.FormFile("logo")
	if err != nil {
		http.Error(w, "Error al obtener el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validar tipo de archivo
	allowedTypes := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/gif":     true,
		"image/svg+xml": true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		http.Error(w, "Tipo de archivo no permitido. Use JPG, PNG, GIF o SVG", http.StatusBadRequest)
		return
	}

	// Validar tamaño del archivo
	if header.Size > 2*1024*1024 { // 2MB
		http.Error(w, "Archivo demasiado grande. Máximo 2MB", http.StatusBadRequest)
		return
	}

	// Obtener ID de usuario (opcional)
	idUsuario := r.URL.Query().Get("id_usuario")
	if idUsuario == "" {
		idUsuario = "default"
	}

	// Crear directorio de logos si no existe
	logoDir := filepath.Join("public", "logos")
	if err := os.MkdirAll(logoDir, 0755); err != nil {
		http.Error(w, "Error al crear directorio de logos", http.StatusInternalServerError)
		return
	}

	// Generar nombre único para el archivo
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		// Determinar extensión basada en el tipo MIME
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/svg+xml":
			ext = ".svg"
		}
	}

	filename := fmt.Sprintf("logo_plantilla_user_%s%s", idUsuario, ext)
	filePath := filepath.Join(logoDir, filename)

	// Crear el archivo en el servidor
	destFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error al crear archivo en el servidor", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Copiar el contenido del archivo
	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, "Error al guardar archivo", http.StatusInternalServerError)
		return
	}

	// Respuesta exitosa
	response := map[string]interface{}{
		"success":  true,
		"message":  "Logo subido correctamente",
		"filename": filename,
		"path":     fmt.Sprintf("/public/logos/%s", filename),
		"url":      fmt.Sprintf("http://localhost:8080/public/logos/%s", filename),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ObtenerLogo devuelve el logo actual del usuario
func (h *LogoPlantillaHandler) ObtenerLogo(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de usuario
	idUsuario := r.URL.Query().Get("id_usuario")
	if idUsuario == "" {
		idUsuario = "default"
	}

	// Buscar logo existente
	logoDir := filepath.Join("public", "logos")
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".svg"}

	var logoPath string
	var logoExists bool

	for _, ext := range extensions {
		filename := fmt.Sprintf("logo_plantilla_user_%s%s", idUsuario, ext)
		fullPath := filepath.Join(logoDir, filename)

		if _, err := os.Stat(fullPath); err == nil {
			logoPath = fmt.Sprintf("http://localhost:8080/public/logos/%s", filename)
			logoExists = true
			break
		}
	}

	response := map[string]interface{}{
		"exists": logoExists,
		"url":    logoPath,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// EliminarLogo elimina el logo del usuario
func (h *LogoPlantillaHandler) EliminarLogo(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "DELETE" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de usuario
	idUsuario := r.URL.Query().Get("id_usuario")
	if idUsuario == "" {
		idUsuario = "default"
	}

	// Buscar y eliminar logo existente
	logoDir := filepath.Join("public", "logos")
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".svg"}

	var eliminado bool

	for _, ext := range extensions {
		filename := fmt.Sprintf("logo_plantilla_user_%s%s", idUsuario, ext)
		fullPath := filepath.Join(logoDir, filename)

		if _, err := os.Stat(fullPath); err == nil {
			if err := os.Remove(fullPath); err == nil {
				eliminado = true
			}
		}
	}

	response := map[string]interface{}{
		"success": eliminado,
		"message": "Logo eliminado correctamente",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListarLogos lista todos los logos disponibles (para administración)
func (h *LogoPlantillaHandler) ListarLogos(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	logoDir := filepath.Join("public", "logos")

	// Leer directorio de logos
	files, err := os.ReadDir(logoDir)
	if err != nil {
		// Si el directorio no existe, devolver lista vacía
		response := map[string]interface{}{
			"logos": []interface{}{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var logos []map[string]interface{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if strings.HasPrefix(filename, "logo_plantilla_user_") {
			// Extraer ID de usuario del nombre del archivo
			parts := strings.Split(filename, "_")
			if len(parts) >= 4 {
				userPart := strings.Join(parts[3:], "_")
				userId := strings.Split(userPart, ".")[0]

				info, _ := file.Info()
				logo := map[string]interface{}{
					"filename": filename,
					"userId":   userId,
					"url":      fmt.Sprintf("http://localhost:8080/public/logos/%s", filename),
					"size":     info.Size(),
					"modified": info.ModTime(),
				}
				logos = append(logos, logo)
			}
		}
	}

	response := map[string]interface{}{
		"logos": logos,
		"count": len(logos),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// obtenerIDUsuario extrae el ID de usuario de los parámetros de la URL
func obtenerIDUsuario(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id_usuario")
	if idStr == "" {
		return 1, nil // Valor por defecto
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("ID de usuario inválido: %s", idStr)
	}

	return id, nil
}
