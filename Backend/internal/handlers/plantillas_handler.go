package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"Facts/internal/db"
	"Facts/internal/utils"
)

// SubirPlantillaHandler maneja la subida de plantillas
func SubirPlantillaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de usuario del token o parámetro
	idUsuarioStr := r.URL.Query().Get("id_usuario")
	if idUsuarioStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Procesar el formulario multipart
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Obtener archivo
	file, handler, err := r.FormFile("plantilla")
	if err != nil {
		http.Error(w, "Error al obtener el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validar extensión
	fileName := handler.Filename
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".doc" && ext != ".docx" {
		http.Error(w, "Tipo de archivo no permitido. Use .doc o .docx", http.StatusBadRequest)
		return
	}

	// Descripción opcional
	descripcion := r.FormValue("descripcion")

	// Crear directorio si no existe
	plantillasDir := "./public/plantillas"
	if err := utils.CreateDirectory(plantillasDir); err != nil {
		http.Error(w, "Error al crear directorio", http.StatusInternalServerError)
		return
	}

	// Generar nombre único para el archivo
	timestamp := time.Now().UnixNano()
	uniqueFileName := fmt.Sprintf("%d_%s", timestamp, fileName)
	filePath := filepath.Join(plantillasDir, uniqueFileName)

	// Guardar archivo físicamente
	dest, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error al crear archivo", http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	_, err = io.Copy(dest, file)
	if err != nil {
		http.Error(w, "Error al guardar archivo", http.StatusInternalServerError)
		return
	}

	// Conectar a la base de datos
	dbConn := db.GetDB()

	// Iniciar transacción
	tx, err := dbConn.Begin()
	if err != nil {
		http.Error(w, "Error de base de datos", http.StatusInternalServerError)
		return
	}

	// Desactivar plantilla activa actual si se marca como activa
	activar := r.FormValue("activar") == "true"
	if activar {
		_, err = tx.Exec("UPDATE plantillas_factura SET activa = false WHERE id_usuario = ?", idUsuario)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error al actualizar plantillas", http.StatusInternalServerError)
			return
		}
	}

	// Insertar registro en base de datos
	result, err := tx.Exec(
		"INSERT INTO plantillas_factura (id_usuario, nombre, ruta_archivo, descripcion, activa) VALUES (?, ?, ?, ?, ?)",
		idUsuario, fileName, filePath, descripcion, activar,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al guardar plantilla en base de datos", http.StatusInternalServerError)
		return
	}

	// Obtener ID generado
	plantillaID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al obtener ID de plantilla", http.StatusInternalServerError)
		return
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	// Responder éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Plantilla guardada correctamente",
		"id":      plantillaID,
		"nombre":  fileName,
		"activa":  activar,
	})
}

// ListarPlantillasHandler obtiene las plantillas de un usuario
func ListarPlantillasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de usuario
	idUsuarioStr := r.URL.Query().Get("id_usuario")
	if idUsuarioStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Conectar a la base de datos
	dbConn := db.GetDB()

	// Consultar plantillas
	rows, err := dbConn.Query(
		"SELECT id, nombre, descripcion, activa, fecha_creacion FROM plantillas_factura WHERE id_usuario = ? ORDER BY fecha_creacion DESC",
		idUsuario,
	)
	if err != nil {
		http.Error(w, "Error al consultar plantillas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Preparar resultado
	var plantillas []map[string]interface{}

	// Procesar resultados
	for rows.Next() {
		var id int
		var nombre, descripcion string
		var activa bool
		var fechaCreacion string

		if err := rows.Scan(&id, &nombre, &descripcion, &activa, &fechaCreacion); err != nil {
			log.Printf("Error al escanear plantilla: %v", err)
			continue
		}

		plantillas = append(plantillas, map[string]interface{}{
			"id":             id,
			"nombre":         nombre,
			"descripcion":    descripcion,
			"activa":         activa,
			"fecha_creacion": fechaCreacion,
		})
	}

	// Verificar errores de iteración
	if err := rows.Err(); err != nil {
		http.Error(w, "Error al procesar plantillas", http.StatusInternalServerError)
		return
	}

	// Responder con plantillas
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"plantillas": plantillas,
	})
}

// ActivarPlantillaHandler establece una plantilla como activa
func ActivarPlantillaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar solicitud
	var request struct {
		IDUsuario   int `json:"id_usuario"`
		IDPlantilla int `json:"id_plantilla"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error al decodificar solicitud", http.StatusBadRequest)
		return
	}

	// Validar datos
	if request.IDUsuario <= 0 || request.IDPlantilla <= 0 {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	// Conectar a la base de datos
	dbConn := db.GetDB()

	// Iniciar transacción
	tx, err := dbConn.Begin()
	if err != nil {
		http.Error(w, "Error de base de datos", http.StatusInternalServerError)
		return
	}

	// Desactivar todas las plantillas del usuario
	_, err = tx.Exec("UPDATE plantillas_factura SET activa = false WHERE id_usuario = ?", request.IDUsuario)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al actualizar plantillas", http.StatusInternalServerError)
		return
	}

	// Activar la plantilla seleccionada
	result, err := tx.Exec(
		"UPDATE plantillas_factura SET activa = true WHERE id = ? AND id_usuario = ?",
		request.IDPlantilla, request.IDUsuario,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al activar plantilla", http.StatusInternalServerError)
		return
	}

	// Verificar que se haya actualizado un registro
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Error al verificar actualización", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		tx.Rollback()
		http.Error(w, "Plantilla no encontrada o no pertenece al usuario", http.StatusNotFound)
		return
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		http.Error(w, "Error al confirmar transacción", http.StatusInternalServerError)
		return
	}

	// Responder éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Plantilla activada correctamente",
	})
}

// EliminarPlantillaHandler elimina una plantilla
func EliminarPlantillaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar solicitud
	var request struct {
		IDUsuario   int `json:"id_usuario"`
		IDPlantilla int `json:"id_plantilla"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Error al decodificar solicitud", http.StatusBadRequest)
		return
	}

	// Validar datos
	if request.IDUsuario <= 0 || request.IDPlantilla <= 0 {
		http.Error(w, "IDs inválidos", http.StatusBadRequest)
		return
	}

	// Conectar a la base de datos
	dbConn := db.GetDB()

	// Obtener la ruta del archivo antes de eliminar
	var rutaArchivo string
	err := dbConn.QueryRow(
		"SELECT ruta_archivo FROM plantillas_factura WHERE id = ? AND id_usuario = ?",
		request.IDPlantilla, request.IDUsuario,
	).Scan(&rutaArchivo)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Plantilla no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al consultar plantilla", http.StatusInternalServerError)
		}
		return
	}

	// Eliminar registro de la base de datos
	result, err := dbConn.Exec(
		"DELETE FROM plantillas_factura WHERE id = ? AND id_usuario = ?",
		request.IDPlantilla, request.IDUsuario,
	)
	if err != nil {
		http.Error(w, "Error al eliminar plantilla", http.StatusInternalServerError)
		return
	}

	// Verificar que se haya eliminado un registro
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error al verificar eliminación", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Plantilla no encontrada o no pertenece al usuario", http.StatusNotFound)
		return
	}

	// Eliminar archivo físico
	if rutaArchivo != "" {
		if err := os.Remove(rutaArchivo); err != nil {
			// No fallamos si el archivo no existe, solo registramos
			log.Printf("Error al eliminar archivo físico: %v", err)
		}
	}

	// Responder éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Plantilla eliminada correctamente",
	})
}

// ObtenerPlantillaActivaHandler obtiene la plantilla activa de un usuario
func ObtenerPlantillaActivaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de usuario
	idUsuarioStr := r.URL.Query().Get("id_usuario")
	if idUsuarioStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Conectar a la base de datos
	dbConn := db.GetDB()

	// Consultar plantilla activa
	var plantilla struct {
		ID          int    `json:"id"`
		Nombre      string `json:"nombre"`
		Descripcion string `json:"descripcion"`
		RutaArchivo string `json:"ruta_archivo"`
	}

	err = dbConn.QueryRow(
		"SELECT id, nombre, descripcion, ruta_archivo FROM plantillas_factura WHERE id_usuario = ? AND activa = true",
		idUsuario,
	).Scan(&plantilla.ID, &plantilla.Nombre, &plantilla.Descripcion, &plantilla.RutaArchivo)

	if err != nil {
		if err == sql.ErrNoRows {
			// No hay plantilla activa, devolver objeto vacío
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"plantilla_activa": nil,
			})
			return
		}
		http.Error(w, "Error al consultar plantilla activa", http.StatusInternalServerError)
		return
	}

	// Responder con plantilla activa
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"plantilla_activa": plantilla,
	})
}

// BuscarPlantillasHandler maneja la búsqueda de plantillas por nombre
func BuscarPlantillasHandler(w http.ResponseWriter, r *http.Request) {
	// Solo permitir método GET
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener término de búsqueda
	query := r.URL.Query().Get("q")
	if query == "" {
		// Si no hay término de búsqueda, redireccionar a listar todas
		ListarPlantillasHandler(w, r)
		return
	}

	query = strings.ToLower(query)

	// Obtener ID de usuario
	idUsuarioStr := r.URL.Query().Get("id_usuario")
	if idUsuarioStr == "" {
		http.Error(w, "ID de usuario no proporcionado", http.StatusBadRequest)
		return
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Conectar a la base de datos
	dbConn := db.GetDB()

	// Consultar plantillas que coincidan con la búsqueda
	rows, err := dbConn.Query(
		"SELECT id, nombre, descripcion, activa, fecha_creacion FROM plantillas_factura WHERE id_usuario = ? AND nombre LIKE ? ORDER BY fecha_creacion DESC",
		idUsuario, "%"+query+"%",
	)
	if err != nil {
		http.Error(w, "Error al buscar plantillas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Preparar resultado
	var plantillas []map[string]interface{}

	// Procesar resultados
	for rows.Next() {
		var id int
		var nombre, descripcion string
		var activa bool
		var fechaCreacion string

		if err := rows.Scan(&id, &nombre, &descripcion, &activa, &fechaCreacion); err != nil {
			log.Printf("Error al escanear plantilla: %v", err)
			continue
		}

		plantillas = append(plantillas, map[string]interface{}{
			"id":             id,
			"nombre":         nombre,
			"descripcion":    descripcion,
			"activa":         activa,
			"fecha_creacion": fechaCreacion,
		})
	}

	// Verificar errores de iteración
	if err := rows.Err(); err != nil {
		http.Error(w, "Error al procesar plantillas", http.StatusInternalServerError)
		return
	}

	// Responder con plantillas
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"plantillas": plantillas,
	})
}
