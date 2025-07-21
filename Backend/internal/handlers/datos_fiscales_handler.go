package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"carlos/Facts/Backend/internal/db"
)

// GetDatosFiscalesHandler maneja la solicitud GET para obtener datos fiscales
func GetDatosFiscalesHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID de usuario del token JWT (implementado en el middleware de autenticación)
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	// Verificar si hay un parámetro id_usuario en la consulta
	idParam := r.URL.Query().Get("id_usuario")
	if idParam != "" {
		// Si el usuario actual no es admin, no permitir consultar otros usuarios
		isAdmin, _ := r.Context().Value("isAdmin").(bool)
		if !isAdmin {
			http.Error(w, "No tiene permisos para consultar estos datos", http.StatusForbidden)
			return
		}

		// Convertir el ID del parámetro a entero
		var err error
		userID, err = strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
			return
		}
	}

	// Obtener datos fiscales de la base de datos
	datos, err := db.ObtenerDatosFiscales(userID)
	if err != nil {
		// Verificar si el error es porque no se encontraron datos
		if err.Error() == "no se encontraron datos fiscales para el usuario" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			// Log del error para depuración
			fmt.Printf("Error al obtener datos fiscales: %v\n", err)
			http.Error(w, "Error al obtener datos fiscales", http.StatusInternalServerError)
		}
		return
	}

	// Responder con JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datos)
}

// UpdateDatosFiscalesHandler maneja la petición para actualizar datos fiscales
func UpdateDatosFiscalesHandler(w http.ResponseWriter, r *http.Request) {
	// IMPORTANTE: Parsear el formulario multipart primero
	r.ParseMultipartForm(32 << 20) // 32MB máximo

	// Depuración - imprimir todos los valores del formulario
	fmt.Println("==== DATOS RECIBIDOS EN EL FORMULARIO ====")
	for key, values := range r.Form {
		fmt.Printf("%s: %v\n", key, values)
	}

	// Obtener el ID de usuario del formulario de manera explícita
	idUsuarioStr := r.FormValue("id_usuario")
	fmt.Printf("ID recibido del formulario: '%s'\n", idUsuarioStr)

	// Convertir a entero de manera segura
	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil || idUsuario <= 0 {
		http.Error(w, "ID de usuario inválido o no proporcionado", http.StatusBadRequest)
		fmt.Printf("Error al procesar ID de usuario: %v\n", err)
		return
	}

	fmt.Printf("ID de usuario procesado correctamente: %d\n", idUsuario)

	// Obtener ID del usuario desde el contexto
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	// Añadir log para verificar que userID es válido
	fmt.Printf("Actualizando datos fiscales para usuario ID: %d\n", userID)

	// Verificar que el usuario sea admin
	isAdmin, ok := r.Context().Value("isAdmin").(bool)
	if !ok || !isAdmin {
		http.Error(w, "Solo los administradores pueden actualizar datos fiscales", http.StatusForbidden)
		return
	}

	// Parsear formulario multipart
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	// Obtener datos del formulario
	rfc := r.FormValue("rfc")
	razonSocial := r.FormValue("razon_social")
	direccionFiscal := r.FormValue("direccion_fiscal")
	codigoPostal := r.FormValue("codigo_postal")
	claveCSD := r.FormValue("clave_csd")
	regimenFiscal := r.FormValue("regimen_fiscal")
	serieDf := r.FormValue("serie_df") // AGREGAR: Obtener serie_df del formulario

	// LOG DETALLADO para depuración
	fmt.Printf("🔍 HANDLER - Datos recibidos del formulario:\n")
	fmt.Printf("   rfc: '%s'\n", rfc)
	fmt.Printf("   razon_social: '%s'\n", razonSocial)
	fmt.Printf("   serie_df: '%s' (longitud: %d)\n", serieDf, len(serieDf))
	fmt.Printf("   direccion_fiscal: '%s'\n", direccionFiscal)
	fmt.Printf("   codigo_postal: '%s'\n", codigoPostal)

	// Verificar el userID y loguearlo
	fmt.Printf("Valor inicial de userID: %d\n", userID)

	// Obtener el ID de usuario del formulario (opcional, para permitir a admins actualizar datos de otros usuarios)
	idUsuarioParam := r.FormValue("id_usuario")
	if idUsuarioParam != "" && isAdmin {
		idUsuarioInt, err := strconv.Atoi(idUsuarioParam)
		if err == nil && idUsuarioInt > 0 {
			userID = idUsuarioInt
			fmt.Printf("Admin actualizando datos para otro usuario: %d\n", userID)
		} else {
			fmt.Printf("Error al convertir id_usuario del formulario: %v\n", err)
		}
	}

	// VALIDACIÓN CRÍTICA: Asegurar que userID es válido
	if userID <= 0 {
		http.Error(w, "ID de usuario inválido o no proporcionado", http.StatusBadRequest)
		fmt.Printf("ERROR: ID de usuario inválido: %d\n", userID)
		return
	}

	// Validar campos requeridos
	if rfc == "" || razonSocial == "" || codigoPostal == "" || regimenFiscal == "" {
		http.Error(w, "Faltan campos requeridos", http.StatusBadRequest)
		return
	}

	// Variables para archivos binarios y nombres
	var archivoCSDKey, archivoCSDCer []byte
	var nombreArchivoKey, nombreArchivoCer string

	// Procesar archivo CSD KEY si se proporciona
	fileCSDKey, fileCSDKeyHeader, err := r.FormFile("csdKey")
	if err == nil {
		defer fileCSDKey.Close()
		archivoCSDKey, err = ioutil.ReadAll(fileCSDKey)
		nombreArchivoKey = fileCSDKeyHeader.Filename
		if err != nil {
			http.Error(w, "Error al leer archivo CSD KEY", http.StatusInternalServerError)
			return
		}
	}

	// Procesar archivo CSD CER si se proporciona
	fileCSDCer, fileCSDCerHeader, err := r.FormFile("csdCer")
	if err == nil {
		defer fileCSDCer.Close()
		archivoCSDCer, err = ioutil.ReadAll(fileCSDCer)
		nombreArchivoCer = fileCSDCerHeader.Filename
		if err != nil {
			http.Error(w, "Error al leer archivo CSD CER", http.StatusInternalServerError)
			return
		}
	}

	// Obtener keyPath del formulario
	keyPath := r.FormValue("key_path")

	// Guardar datos fiscales y obtener el idDatosFiscales
	idDatosFiscales, err := db.GuardarDatosFiscales(
		rfc, razonSocial, direccionFiscal, codigoPostal,
		archivoCSDKey, archivoCSDCer,
		nombreArchivoKey, nombreArchivoCer,
		claveCSD, regimenFiscal, serieDf,
		userID,
		keyPath, // NUEVO argumento para la ruta .key
	)
	if err != nil {
		fmt.Printf("Error al guardar datos fiscales para usuario %d: %v\n", userID, err)
		http.Error(w, "Error al guardar datos fiscales", http.StatusInternalServerError)
		return
	}

	// Responder con éxito y el idDatosFiscales
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":            "success",
		"message":           "Datos fiscales actualizados correctamente",
		"id_datos_fiscales": idDatosFiscales,
	})
}

// DescargarCertificadoHandler permite descargar los certificados CSD
func DescargarCertificadoHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener ID del usuario desde el contexto
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Usuario no autenticado", http.StatusUnauthorized)
		return
	}

	// Obtener el tipo de certificado a descargar (key o cer)
	tipoArchivo := r.URL.Query().Get("tipo")
	if tipoArchivo != "key" && tipoArchivo != "cer" {
		http.Error(w, "Tipo de archivo inválido", http.StatusBadRequest)
		return
	}

	// Obtener los archivos binarios
	archivoKey, archivoCer, err := db.ObtenerCertificadosCSD(userID)
	if err != nil {
		http.Error(w, "Error al obtener certificados", http.StatusInternalServerError)
		return
	}

	// Obtener datos fiscales para los nombres de archivos
	datos, err := db.ObtenerDatosFiscales(userID)
	if err != nil {
		http.Error(w, "Error al obtener datos fiscales", http.StatusInternalServerError)
		return
	}

	var nombreArchivo string
	var contenidoArchivo []byte

	if tipoArchivo == "key" {
		// Usar el nombre del archivo desde datos si está disponible
		if datos["archivo_csd_key"] != nil {
			nombreArchivo = "certificado.key" // O usar un nombre más específico
		} else {
			nombreArchivo = "certificado.key"
		}
		contenidoArchivo = archivoKey
	} else {
		// Usar el nombre del archivo desde datos si está disponible
		if datos["archivo_csd_cer"] != nil {
			nombreArchivo = "certificado.cer" // O usar un nombre más específico
		} else {
			nombreArchivo = "certificado.cer"
		}
		contenidoArchivo = archivoCer
	}

	// Configurar los headers para la descarga
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", nombreArchivo))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(contenidoArchivo)))

	// Escribir el contenido del archivo
	w.Write(contenidoArchivo)
}
