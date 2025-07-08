package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"carlos/Facts/Backend/internal/db"
	"carlos/Facts/Backend/internal/models"
	"carlos/Facts/Backend/internal/services"
	"carlos/Facts/Backend/internal/utils"
)

// LogoHandler maneja las rutas relacionadas con logos
type LogoHandler struct {
	logoService *services.LogoService
}

// NewLogoHandler crea un nuevo handler de logos
func NewLogoHandler() *LogoHandler {
	return &LogoHandler{logoService: services.NewLogoService()}
}

// SubirLogo maneja la subida de logos
func (h *LogoHandler) SubirLogo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var logoReq models.LogoRequest
	if err := json.NewDecoder(r.Body).Decode(&logoReq); err != nil {
		utils.RespondWithError(w, "Error al decodificar JSON: "+err.Error())
		return
	}

	// Validaciones
	if logoReq.IdUsuario == 0 {
		utils.RespondWithError(w, "ID de usuario es requerido")
		return
	}

	if logoReq.NombreLogo == "" {
		utils.RespondWithError(w, "Nombre del logo es requerido")
		return
	}

	if logoReq.ImagenBase64 == "" {
		utils.RespondWithError(w, "Imagen es requerida")
		return
	}

	if !services.ValidarTipoImagen(logoReq.TipoMime) {
		utils.RespondWithError(w, "Tipo de imagen no válido")
		return
	}

	// Guardar logo
	logo, err := h.logoService.GuardarLogo(logoReq)
	if err != nil {
		utils.RespondWithError(w, "Error al guardar logo: "+err.Error())
		return
	}

	// Respuesta
	response := models.LogoResponse{
		ID:            logo.ID,
		IdUsuario:     logo.IdUsuario,
		IdEmpresa:     logo.IdEmpresa,
		NombreLogo:    logo.NombreLogo,
		TipoMime:      logo.TipoMime,
		TamañoArchivo: logo.TamañoArchivo,
		EstadoLogo:    logo.EstadoLogo,
		FechaCreacion: logo.FechaCreacion.Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logo subido exitosamente",
		"data":    response,
	})
}

// ObtenerLogosUsuario obtiene todos los logos de un usuario
func (h *LogoHandler) ObtenerLogosUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idUsuarioStr := r.URL.Query().Get("id_usuario")
	if idUsuarioStr == "" {
		http.Error(w, "ID de usuario es requerido", http.StatusBadRequest)
		return
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		http.Error(w, "ID de usuario no válido", http.StatusBadRequest)
		return
	}

	logos, err := h.logoService.ListarLogosPorUsuario(idUsuario)
	if err != nil {
		http.Error(w, "Error al obtener logos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    logos,
	})
}

// ObtenerLogoImagen obtiene la imagen de un logo específico
func (h *LogoHandler) ObtenerLogoImagen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idLogoStr := r.URL.Query().Get("id")
	if idLogoStr == "" {
		http.Error(w, "ID de logo es requerido", http.StatusBadRequest)
		return
	}

	idLogo, err := strconv.Atoi(idLogoStr)
	if err != nil {
		http.Error(w, "ID de logo no válido", http.StatusBadRequest)
		return
	}

	logo, err := h.logoService.ObtenerLogoPorID(idLogo)
	if err != nil {
		http.Error(w, "Logo no encontrado: "+err.Error(), http.StatusNotFound)
		return
	}

	// Configurar headers para la imagen
	w.Header().Set("Content-Type", logo.TipoMime)
	w.Header().Set("Content-Length", strconv.Itoa(len(logo.ImagenLogo)))
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache por 24 horas

	// Escribir la imagen
	w.Write(logo.ImagenLogo)
}

// ObtenerLogoActivoUsuario obtiene el logo activo de un usuario
func (h *LogoHandler) ObtenerLogoActivoUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, "Método no permitido")
		return
	}

	idUsuarioStr := r.URL.Query().Get("id_usuario")
	if idUsuarioStr == "" {
		utils.RespondWithError(w, "ID de usuario es requerido")
		return
	}

	// Limpiar el formato "default:1" a solo "1"
	if len(idUsuarioStr) > 8 && idUsuarioStr[:8] == "default:" {
		idUsuarioStr = idUsuarioStr[8:]
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		utils.RespondWithError(w, "ID de usuario no válido: "+err.Error())
		return
	}

	logo, err := h.logoService.ObtenerLogoActivoPorUsuario(idUsuario)
	if err != nil {
		utils.RespondWithError(w, "Logo activo no encontrado: "+err.Error())
		return
	}

	// Configurar headers para la imagen
	w.Header().Set("Content-Type", logo.TipoMime)
	w.Header().Set("Content-Length", strconv.Itoa(len(logo.ImagenLogo)))
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache por 24 horas

	// Escribir la imagen
	w.Write(logo.ImagenLogo)
}

// ObtenerLogoActivoJSON obtiene el logo activo de un usuario en formato JSON para el NavBar
func (h *LogoHandler) ObtenerLogoActivoJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, "Método no permitido")
		return
	}

	idUsuarioStr := r.URL.Query().Get("id_usuario")
	if idUsuarioStr == "" {
		utils.RespondWithError(w, "ID de usuario es requerido")
		return
	}

	// Limpiar el formato "default:1" a solo "1"
	if len(idUsuarioStr) > 8 && idUsuarioStr[:8] == "default:" {
		idUsuarioStr = idUsuarioStr[8:]
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		utils.RespondWithError(w, "ID de usuario no válido: "+err.Error())
		return
	}

	logo, err := h.logoService.ObtenerLogoActivoPorUsuario(idUsuario)
	if err != nil {
		// Devolver respuesta indicando que no existe logo activo
		response := map[string]interface{}{
			"exists":  false,
			"message": "No se encontró logo activo",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convertir imagen a base64
	imagenBase64 := base64.StdEncoding.EncodeToString(logo.ImagenLogo)

	response := map[string]interface{}{
		"exists":         true,
		"id":             logo.ID,
		"nombre_logo":    logo.NombreLogo,
		"tipo":           logo.TipoMime,
		"imagen_base64":  imagenBase64,
		"fecha_creacion": logo.FechaCreacion.Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache por 1 hora
	json.NewEncoder(w).Encode(response)
}

// ActivarLogo activa un logo específico
func (h *LogoHandler) ActivarLogo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		IdLogo    int `json:"id_logo"`
		IdUsuario int `json:"id_usuario"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if request.IdLogo == 0 || request.IdUsuario == 0 {
		http.Error(w, "ID de logo e ID de usuario son requeridos", http.StatusBadRequest)
		return
	}

	err := h.logoService.ActivarLogo(request.IdLogo, request.IdUsuario)
	if err != nil {
		http.Error(w, "Error al activar logo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logo activado exitosamente",
	})
}

// EliminarLogo elimina un logo
func (h *LogoHandler) EliminarLogo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, "Método no permitido")
		return
	}

	var request struct {
		IdLogo    int `json:"id_logo"`
		IdUsuario int `json:"id_usuario"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, "Error al decodificar JSON: "+err.Error())
		return
	}

	if request.IdLogo == 0 || request.IdUsuario == 0 {
		utils.RespondWithError(w, "ID de logo e ID de usuario son requeridos")
		return
	}

	err := h.logoService.EliminarLogo(request.IdLogo, request.IdUsuario)
	if err != nil {
		utils.RespondWithError(w, "Error al eliminar logo: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logo eliminado exitosamente",
	})
}

// DebugLogosHandler - endpoint temporal para debuggear logos
func DebugLogosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idUsuarioStr := r.URL.Query().Get("id_usuario")
	if idUsuarioStr == "" {
		http.Error(w, "ID de usuario requerido", http.StatusBadRequest)
		return
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Conectar a la base de datos
	dbConn, err := db.ConnectUserDB()
	if err != nil {
		log.Printf("Error al conectar a la base de datos: %v", err)
		http.Error(w, "Error de conexión a base de datos", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Obtener todos los logos del usuario
	query := `SELECT id, id_usuario, nombre_logo, tipo_mime, tamaño_archivo, estado_logo, fecha_creacion FROM logos_empresas WHERE id_usuario = ?`

	rows, err := dbConn.Query(query, idUsuario)
	if err != nil {
		log.Printf("Error al ejecutar query: %v", err)
		http.Error(w, "Error al consultar logos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logos []map[string]interface{}

	for rows.Next() {
		var id, idUsuarioDB, tamanoArchivo, estadoLogo int
		var nombreLogo, tipoMime, fechaCreacion string

		err = rows.Scan(&id, &idUsuarioDB, &nombreLogo, &tipoMime, &tamanoArchivo, &estadoLogo, &fechaCreacion)
		if err != nil {
			log.Printf("Error al escanear fila: %v", err)
			continue
		}

		logo := map[string]interface{}{
			"id":             id,
			"id_usuario":     idUsuarioDB,
			"nombre_logo":    nombreLogo,
			"tipo_mime":      tipoMime,
			"tamaño_archivo": tamanoArchivo,
			"estado_logo":    estadoLogo,
			"fecha_creacion": fechaCreacion,
			"es_activo":      estadoLogo == 1,
		}

		logos = append(logos, logo)
	}

	response := map[string]interface{}{
		"total_logos": len(logos),
		"logos":       logos,
		"usuario_id":  idUsuario,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ConfigurarRutasLogos configura las rutas para el manejo de logos
func (h *LogoHandler) ConfigurarRutasLogos(mux *http.ServeMux) {
	mux.HandleFunc("/api/logos/subir", h.SubirLogo)
	mux.HandleFunc("/api/logos/usuario", h.ObtenerLogosUsuario)
	mux.HandleFunc("/api/logos/imagen", h.ObtenerLogoImagen)
	mux.HandleFunc("/api/logos/activo", h.ObtenerLogoActivoUsuario)
	mux.HandleFunc("/api/logos/activar", h.ActivarLogo)
	mux.HandleFunc("/api/logos/eliminar", h.EliminarLogo)
	mux.HandleFunc("/api/logos/debug", DebugLogosHandler)
}
