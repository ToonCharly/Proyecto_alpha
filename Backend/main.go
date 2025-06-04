package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"carlos/Facts/Backend/internal/db"
	"carlos/Facts/Backend/internal/handlers"
	"carlos/Facts/Backend/internal/models"
	"carlos/Facts/Backend/internal/utils"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Inicializar la conexión a la base de datos
	db.InitDB()

	// Obtener conexión a optimus para los handlers que la necesitan
	optimusDB, err := db.ConnectToOptimus()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos optimus: %v", err)
	}

	// Crear directorios necesarios
	directorios := []string{"./templates", "./templates/facturas"}
	for _, dir := range directorios {
		if err := utils.CreateDirectory(dir); err != nil {
			log.Fatalf("Error al crear directorio %s: %v", dir, err)
		}
	}

	// Definir endpoints
	http.Handle("/api/factura", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		criterio := strings.TrimSpace(r.URL.Query().Get("criterio"))
		if criterio == "" {
			utils.RespondWithError(w, "Criterio no proporcionado")
			return
		}
		handlers.BuscarFactura(db.GetDB(), w, criterio)
	})))

	http.Handle("/api/generar_factura", utils.EnableCors(http.HandlerFunc(handlers.GenerarFacturaHandler)))

	http.Handle("/api/plantillas", utils.EnableCors(http.HandlerFunc(handlers.BuscarPlantillasHandler)))

	// Usar la conexión a optimus para los handlers que la necesitan
	http.Handle("/api/ventas", utils.EnableCors(http.HandlerFunc(handlers.VentasHandler(optimusDB))))

	http.Handle("/api/historial_facturas", utils.EnableCors(http.HandlerFunc(handlers.HistorialFacturasHandler(db.GetDB()))))

	// Endpoint para obtener regímenes fiscales
	http.Handle("/api/regimenes-fiscales", utils.EnableCors(http.HandlerFunc(handlers.GetRegimenesFiscales)))

	// Usar la conexión a optimus para el diagnóstico
	http.Handle("/api/optimus/diagnostico", utils.EnableCors(http.HandlerFunc(handlers.DiagnosticoVentasHandler(optimusDB))))

	// Endpoint para registrar usuarios
	http.Handle("/api/registrar_usuario", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var registerData struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&registerData); err != nil {
			log.Printf("Error decodificando JSON: %v", err)
			utils.RespondWithError(w, "Error al leer los datos: "+err.Error())
			return
		}

		if registerData.Username == "" || registerData.Email == "" || registerData.Phone == "" || registerData.Password == "" {
			utils.RespondWithError(w, "Todos los campos son obligatorios")
			return
		}

		userData, err := db.RegisterUser(registerData.Username, registerData.Email, registerData.Phone, registerData.Password)
		if err != nil {
			log.Printf("Error al registrar usuario: %v", err)
			utils.RespondWithError(w, fmt.Sprintf("Error al registrar el usuario: %v", err))
			return
		}

		// Responder con los datos del usuario, incluyendo el rol
		utils.RespondWithJSON(w, http.StatusCreated, userData)
	})))

	// Endpoint para actualizar información de usuario
	http.Handle("/api/actualizar_usuario", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var userData struct {
			Email        string `json:"email"`
			Username     string `json:"username"`
			Phone        string `json:"phone"`
			Direccion    string `json:"direccion"`
			CodigoPostal string `json:"codigo_postal"`
			Ciudad       string `json:"ciudad"`
			Estado       string `json:"estado"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&userData); err != nil {
			log.Printf("Error decodificando JSON de actualización: %v", err)
			utils.RespondWithError(w, "Error al leer los datos de actualización")
			return
		}

		if userData.Email == "" || userData.Username == "" {
			utils.RespondWithError(w, "Email y nombre de usuario son obligatorios")
			return
		}

		err := db.UpdateUser(db.GetDB(), userData.Email, userData.Username,
			userData.Phone, userData.Direccion,
			userData.CodigoPostal, userData.Ciudad, userData.Estado)

		if err != nil {
			log.Printf("Error al actualizar usuario: %v", err)
			utils.RespondWithError(w, fmt.Sprintf("Error al actualizar datos: %v", err))
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Datos actualizados correctamente",
		})
	})))

	// Endpoint para login de usuarios
	http.Handle("/api/login", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var loginData struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&loginData); err != nil {
			log.Printf("Error decodificando JSON de login: %v", err)
			utils.RespondWithError(w, "Error al leer los datos de login")
			return
		}

		if (loginData.Email == "" && loginData.Username == "") || loginData.Password == "" {
			utils.RespondWithError(w, "Se requiere email/username y contraseña")
			return
		}

		userData, err := db.LoginUser(loginData.Email, loginData.Username, loginData.Password)
		if err != nil {
			log.Printf("Error en login: %v", err)
			http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, userData)
	})))

	// Endpoint para manejar empresas
	http.Handle("/api/empresas", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var empresa models.Empresa
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&empresa); err != nil {
				log.Printf("Error decodificando JSON de empresa: %v", err)
				utils.RespondWithError(w, "Error al leer los datos de la empresa")
				return
			}

			if empresa.RFC == "" || empresa.RazonSocial == "" || empresa.IdUsuario == 0 {
				utils.RespondWithError(w, "RFC, Razón Social e ID de Usuario son obligatorios")
				return
			}

			id, err := models.InsertarEmpresa(empresa)
			if err != nil {
				log.Printf("Error al insertar empresa: %v", err)
				utils.RespondWithError(w, fmt.Sprintf("Error al registrar la empresa: %v", err))
				return
			}

			empresa.ID = int(id)
			utils.RespondWithJSON(w, http.StatusOK, empresa)

		case http.MethodGet:
			idUsuarioStr := r.URL.Query().Get("id_usuario")
			if idUsuarioStr == "" {
				http.Error(w, "El parámetro id_usuario es requerido", http.StatusBadRequest)
				return
			}

			idUsuario, err := strconv.Atoi(idUsuarioStr)
			if err != nil {
				http.Error(w, "El parámetro id_usuario debe ser un número entero", http.StatusBadRequest)
				return
			}

			empresas, err := models.ObtenerEmpresasPorUsuario(idUsuario)
			if err != nil {
				http.Error(w, "Error al obtener las empresas", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(empresas)

		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})))

	// NUEVO: Endpoint para manejar operaciones en empresas específicas por ID
	http.Handle("/api/empresas/", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el ID de la empresa de la URL
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "URL inválida. Se esperaba /api/empresas/{id}", http.StatusBadRequest)
			return
		}

		// Convertir el ID de string a entero
		empresaID, err := strconv.Atoi(pathParts[3])
		if err != nil {
			http.Error(w, "ID de empresa inválido. Debe ser un número entero.", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodDelete:
			// Eliminar la empresa
			err := models.EliminarEmpresa(empresaID)
			if err != nil {
				log.Printf("Error al eliminar empresa con ID %d: %v", empresaID, err)
				http.Error(w, fmt.Sprintf("Error al eliminar empresa: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": fmt.Sprintf("Empresa con ID %d eliminada exitosamente", empresaID),
			})

		case http.MethodGet:
			// Obtener una empresa específica por ID
			empresa, err := models.ObtenerEmpresaPorID(empresaID)
			if err != nil {
				log.Printf("Error al obtener empresa con ID %d: %v", empresaID, err)
				http.Error(w, fmt.Sprintf("Error al obtener empresa: %v", err), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(empresa)

		case http.MethodPut:
			// Actualizar empresa existente
			var empresa models.Empresa
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&empresa); err != nil {
				log.Printf("Error decodificando JSON de empresa: %v", err)
				http.Error(w, "Error al leer los datos de la empresa", http.StatusBadRequest)
				return
			}

			// Asegurarse de que el ID en la URL sea el mismo que en el cuerpo
			empresa.ID = empresaID

			err := models.ActualizarEmpresa(empresa)
			if err != nil {
				log.Printf("Error al actualizar empresa: %v", err)
				http.Error(w, fmt.Sprintf("Error al actualizar empresa: %v", err), http.StatusInternalServerError)
				return
			}

			// Obtener la empresa actualizada
			empresaActualizada, err := models.ObtenerEmpresaPorID(empresaID)
			if err != nil {
				log.Printf("Empresa actualizada pero error al recuperarla: %v", err)
				http.Error(w, "Empresa actualizada pero error al recuperar los datos", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(empresaActualizada)

		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})))

	// NUEVO: Endpoint para verificar conexión del servidor
	http.Handle("/api/ping", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Servidor funcionando correctamente",
		})
	})))

	// Endpoint para ventas en la base de datos optimus
	http.Handle("/api/optimus/ventas", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		serie := r.URL.Query().Get("serie")
		if serie == "" {
			utils.RespondWithError(w, "Serie no proporcionada")
			return
		}

		// Crear una conexión específica a la base de datos optimus
		optimusDB, err := db.ConnectToOptimus()
		if err != nil {
			log.Printf("Error al conectar a la base de datos optimus: %v", err)
			utils.RespondWithError(w, "Error al conectar a la base de datos")
			return
		}
		defer optimusDB.Close()

		// Consultar las ventas en la base de datos optimus uniendo las tablas crm_pedidos y crm_pedidos_det
		rows, err := optimusDB.Query(`
        SELECT p.id_pedido, pd.descripcion AS producto, 
               pd.cantidad, pd.precio_o AS precio, pd.descuento, 
               (pd.precio_o * pd.cantidad) - pd.descuento AS total 
        FROM crm_pedidos p
        JOIN crm_pedidos_det pd ON p.id_pedido = pd.id_pedido
        WHERE p.clave_pedido LIKE ?`, "%"+serie+"%")
		if err != nil {
			log.Printf("Error al consultar ventas: %v", err)
			utils.RespondWithError(w, "Error al buscar las ventas")
			return
		}
		defer rows.Close()

		var ventas []map[string]interface{}
		for rows.Next() {
			var idPedido, producto string
			var cantidad int
			var precio, descuento, total float64

			if err := rows.Scan(&idPedido, &producto, &cantidad, &precio, &descuento, &total); err != nil {
				log.Printf("Error al escanear fila: %v", err)
				continue
			}

			venta := map[string]interface{}{
				"idPedido":  idPedido,
				"producto":  producto,
				"cantidad":  cantidad,
				"precio":    precio,
				"descuento": descuento,
				"total":     total,
			}

			ventas = append(ventas, venta)
		}

		if err = rows.Err(); err != nil {
			log.Printf("Error durante la iteración de filas: %v", err)
			utils.RespondWithError(w, "Error al procesar los resultados")
			return
		}

		// Devolver el resultado
		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"ventas": ventas,
		})
	})))

	// NUEVO: Endpoint para historial de facturas (formato con guión)
	http.Handle("/api/historial-facturas", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Obtener historial de facturas por usuario
			idUsuarioStr := r.URL.Query().Get("id_usuario")
			if idUsuarioStr == "" {
				http.Error(w, "El parámetro id_usuario es requerido", http.StatusBadRequest)
				return
			}

			idUsuario, err := strconv.Atoi(idUsuarioStr)
			if err != nil {
				http.Error(w, "El parámetro id_usuario debe ser un número entero", http.StatusBadRequest)
				return
			}

			facturas, err := models.ObtenerHistorialFacturasPorUsuario(idUsuario)
			if err != nil {
				log.Printf("Error al obtener historial de facturas: %v", err)
				http.Error(w, "Error al obtener el historial de facturas", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(facturas)

		case http.MethodPost:
			// Registrar una nueva factura en el historial
			var historialData struct {
				IDUsuario           int     `json:"id_usuario"`
				RFCReceptor         string  `json:"rfc_receptor"`
				RazonSocialReceptor string  `json:"razon_social_receptor"`
				ClaveTicket         string  `json:"clave_ticket"`
				Total               float64 `json:"total"`
				UsoCFDI             string  `json:"uso_cfdi"`
				Observaciones       string  `json:"observaciones"`
			}

			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&historialData); err != nil {
				log.Printf("Error decodificando JSON de historial: %v", err)
				utils.RespondWithError(w, "Error al leer los datos del historial")
				return
			}

			// Validar datos obligatorios
			if historialData.IDUsuario == 0 || historialData.RFCReceptor == "" ||
				historialData.RazonSocialReceptor == "" || historialData.ClaveTicket == "" ||
				historialData.Total <= 0 || historialData.UsoCFDI == "" {
				utils.RespondWithError(w, "Faltan campos obligatorios para el registro de historial")
				return
			}

			// Insertar en base de datos
			id, err := models.InsertarHistorialFactura(
				historialData.IDUsuario,
				historialData.RFCReceptor,
				historialData.RazonSocialReceptor,
				historialData.ClaveTicket,
				historialData.Total,
				historialData.UsoCFDI,
				historialData.Observaciones,
			)

			if err != nil {
				log.Printf("Error al registrar historial: %v", err)
				utils.RespondWithError(w, fmt.Sprintf("Error al registrar en el historial: %v", err))
				return
			}

			utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
				"id":      id,
				"message": "Factura registrada correctamente en el historial",
			})

		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})))

	// Endpoint para descargar una factura del historial
	http.Handle("/api/descargar-factura/", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Obtener el ID de la factura de la URL
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "URL inválida. Se esperaba /api/descargar-factura/{id}", http.StatusBadRequest)
			return
		}

		// Convertir el ID de string a entero
		facturaID, err := strconv.Atoi(pathParts[3])
		if err != nil {
			http.Error(w, "ID de factura inválido. Debe ser un número entero.", http.StatusBadRequest)
			return
		}
		handlers.DescargarFacturaHandler(w, r, facturaID)
	})))

	// Añadir nuevos endpoints (corregido)
	http.Handle("/api/reset-password-request", utils.EnableCors(http.HandlerFunc(handlers.ResetPasswordRequestHandler(db.GetDB()))))
	http.Handle("/api/reset-password", utils.EnableCors(http.HandlerFunc(handlers.ResetPasswordHandler(db.GetDB()))))

	// Endpoint para obtener todos los usuarios (solo para admins)
	http.Handle("/api/usuarios", utils.EnableCors(http.HandlerFunc(handlers.GetAllUsersHandler)))

	// Endpoint para actualizar el rol de un usuario
	http.Handle("/api/usuarios/", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		// Verificar que la ruta sea /api/usuarios/{id}/rol
		if len(pathParts) >= 5 && pathParts[4] == "rol" {
			handlers.UpdateUserRoleHandler(w, r)
			return
		}
		http.Error(w, "Endpoint no encontrado", http.StatusNotFound)
	})))

	// Endpoint para obtener datos fiscales
	http.Handle("/api/datos-fiscales", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Obtener ID de usuario de la solicitud
		idUsuarioStr := r.URL.Query().Get("id_usuario")
		if idUsuarioStr == "" {
			http.Error(w, "El parámetro id_usuario es requerido", http.StatusBadRequest)
			return
		}

		// Convertir a entero
		idUsuario, err := strconv.Atoi(idUsuarioStr)
		if err != nil {
			http.Error(w, "El parámetro id_usuario debe ser un número entero", http.StatusBadRequest)
			return
		}

		// Verificar si el usuario es administrador
		var isAdmin bool
		err = db.GetDB().QueryRow("SELECT role = 'admin' FROM usuarios WHERE id = ?", idUsuario).Scan(&isAdmin)
		if err != nil {
			log.Printf("Error al verificar rol de usuario: %v", err)
			http.Error(w, "Error al verificar permisos de usuario", http.StatusInternalServerError)
			return
		}

		// Obtener datos fiscales (siempre será un solo registro)
		query := `SELECT id, rfc, razon_social, direccion_fiscal, codigo_postal, 
              ruta_csd_key, ruta_csd_cer, clave_csd, regimen_fiscal 
              FROM datos_fiscales LIMIT 1`

		var datosFiscales struct {
			ID              int    `json:"id"`
			RFC             string `json:"rfcEmisor"`
			RazonSocial     string `json:"razonSocial"`
			DireccionFiscal string `json:"direccionFiscal"`
			CodigoPostal    string `json:"cp"`
			RutaCsdKey      string `json:"rutaCsdKey"`
			RutaCsdCer      string `json:"rutaCsdCer"`
			ClaveCsd        string `json:"claveArchivoCSD"`
			RegimenFiscal   string `json:"regimenFiscal"`
		}

		err = db.GetDB().QueryRow(query).Scan(
			&datosFiscales.ID,
			&datosFiscales.RFC,
			&datosFiscales.RazonSocial,
			&datosFiscales.DireccionFiscal,
			&datosFiscales.CodigoPostal,
			&datosFiscales.RutaCsdKey,
			&datosFiscales.RutaCsdCer,
			&datosFiscales.ClaveCsd,
			&datosFiscales.RegimenFiscal,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				// Si no hay datos, devolver objeto vacío (HTTP 200)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{})
				return
			}
			// Otro error
			log.Printf("Error al consultar datos fiscales: %v", err)
			http.Error(w, "Error al consultar datos fiscales", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(datosFiscales)
	})))

	// Endpoint para actualizar datos fiscales
	http.Handle("/api/actualizar-datos-fiscales", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Procesar el formulario multipart para los archivos
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			log.Printf("Error al procesar formulario: %v", err)
			http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
			return
		}

		// Obtener datos del formulario
		rfc := r.FormValue("rfc")
		razonSocial := r.FormValue("razon_social")
		direccionFiscal := r.FormValue("direccion_fiscal")
		codigoPostal := r.FormValue("codigo_postal")
		claveCsd := r.FormValue("clave_csd")
		regimenFiscal := r.FormValue("regimen_fiscal")

		// Validar campos obligatorios
		if rfc == "" || razonSocial == "" || codigoPostal == "" || regimenFiscal == "" {
			log.Printf("Faltan campos obligatorios: RFC=%s, RazonSocial=%s, CP=%s, RegimenFiscal=%s",
				rfc, razonSocial, codigoPostal, regimenFiscal)
			http.Error(w, "RFC, Razón Social, Código Postal y Régimen Fiscal son obligatorios", http.StatusBadRequest)
			return
		}

		// Crear directorio para certificados si no existe
		certDir := "./certificados"
		if err := utils.CreateDirectory(certDir); err != nil {
			log.Printf("Error al crear directorio para certificados: %v", err)
			http.Error(w, "Error al procesar los archivos", http.StatusInternalServerError)
			return
		}

		// Variables para rutas de archivos
		var rutaCsdKey, rutaCsdCer string

		// Manejar archivo CSD Key
		csdKeyFile, _, err := r.FormFile("csdKey")
		if err == nil {
			defer csdKeyFile.Close()

			// Crear ruta única para el archivo
			rutaCsdKey = fmt.Sprintf("%s/%s_key.key", certDir, rfc)

			// Guardar archivo
			keyFile, err := os.Create(rutaCsdKey)
			if err != nil {
				log.Printf("Error al crear archivo key: %v", err)
				http.Error(w, "Error al guardar archivo CSD Key", http.StatusInternalServerError)
				return
			}
			defer keyFile.Close()

			_, err = io.Copy(keyFile, csdKeyFile)
			if err != nil {
				log.Printf("Error al copiar archivo key: %v", err)
				http.Error(w, "Error al guardar archivo CSD Key", http.StatusInternalServerError)
				return
			}
		}

		// Manejar archivo CSD Cer
		csdCerFile, _, err := r.FormFile("csdCer")
		if err == nil {
			defer csdCerFile.Close()

			// Crear ruta única para el archivo
			rutaCsdCer = fmt.Sprintf("%s/%s_cer.cer", certDir, rfc)

			// Guardar archivo
			cerFile, err := os.Create(rutaCsdCer)
			if err != nil {
				log.Printf("Error al crear archivo cer: %v", err)
				http.Error(w, "Error al guardar archivo CSD Cer", http.StatusInternalServerError)
				return
			}
			defer cerFile.Close()

			_, err = io.Copy(cerFile, csdCerFile)
			if err != nil {
				log.Printf("Error al copiar archivo cer: %v", err)
				http.Error(w, "Error al guardar archivo CSD Cer", http.StatusInternalServerError)
				return
			}
		}

		// Conectar a la base de datos
		dbConn, err := db.ConnectUserDB()
		if err != nil {
			log.Printf("Error al conectar a la base de datos: %v", err)
			http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
			return
		}
		defer dbConn.Close()

		// Verificar si ya existen datos fiscales
		var count int
		err = dbConn.QueryRow("SELECT COUNT(*) FROM datos_fiscales").Scan(&count)
		if err != nil {
			log.Printf("Error al verificar datos existentes: %v", err)
			http.Error(w, "Error al procesar la solicitud", http.StatusInternalServerError)
			return
		}

		var query string
		var args []interface{}

		if count == 0 {
			// Insertar nuevos datos
			query = `INSERT INTO datos_fiscales 
                (rfc, razon_social, direccion_fiscal, codigo_postal, 
                ruta_csd_key, ruta_csd_cer, clave_csd, regimen_fiscal) 
                VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
			args = []interface{}{
				rfc, razonSocial, direccionFiscal, codigoPostal,
				rutaCsdKey, rutaCsdCer, claveCsd, regimenFiscal,
			}
		} else {
			// Actualizar datos existentes
			query = `UPDATE datos_fiscales SET 
                rfc = ?, razon_social = ?, direccion_fiscal = ?, codigo_postal = ?,
                regimen_fiscal = ?`
			args = []interface{}{
				rfc, razonSocial, direccionFiscal, codigoPostal, regimenFiscal,
			}

			// Añadir campos opcionales solo si se proporcionaron
			if claveCsd != "" {
				query += ", clave_csd = ?"
				args = append(args, claveCsd)
			}

			if rutaCsdKey != "" {
				query += ", ruta_csd_key = ?"
				args = append(args, rutaCsdKey)
			}

			if rutaCsdCer != "" {
				query += ", ruta_csd_cer = ?"
				args = append(args, rutaCsdCer)
			}

			query += " WHERE id = 1"
		}

		// Ejecutar consulta
		_, err = dbConn.Exec(query, args...)
		if err != nil {
			log.Printf("Error al guardar datos fiscales: %v", err)
			http.Error(w, "Error al guardar los datos fiscales", http.StatusInternalServerError)
			return
		}

		// Responder éxito
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Datos fiscales actualizados correctamente",
		})
	})))

	// Iniciar limpieza programada de tokens de recuperación
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				handlers.LimpiarTokensExpirados()
			}
		}
	}()

	// Endpoint para obtener detalles de usuario por identificador (email o username)
	http.Handle("/api/usuario", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		identifier := r.URL.Query().Get("identifier")
		if identifier == "" {
			http.Error(w, "Parámetro identifier requerido", http.StatusBadRequest)
			return
		}

		// Conectar a la base de datos
		userDB, err := db.ConnectUserDB()
		if err != nil {
			log.Printf("Error al conectar a la base de datos: %v", err)
			http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
			return
		}
		defer userDB.Close()

		// Buscar usuario por email o username
		query := `SELECT id, username, email, phone, role FROM usuarios WHERE email = ? OR username = ?`
		var user struct {
			ID       int
			Username string
			Email    string
			Phone    string
			Role     string
		}

		err = userDB.QueryRow(query, identifier, identifier).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Usuario no encontrado", http.StatusNotFound)
				return
			}
			log.Printf("Error al consultar usuario: %v", err)
			http.Error(w, "Error al consultar usuario", http.StatusInternalServerError)
			return
		}

		// Preparar respuesta
		userData := map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
			"role":     user.Role,
		}

		// Añadir detalles adicionales si existen
		var direccion, codigoPostal, ciudad, estado sql.NullString
		queryDetalles := `SELECT direccion, codigo_postal, ciudad, estado 
                      FROM usuario_detalles WHERE usuario_id = ?`

		err = userDB.QueryRow(queryDetalles, user.ID).Scan(&direccion, &codigoPostal, &ciudad, &estado)
		if err == nil {
			if direccion.Valid {
				userData["direccion"] = direccion.String
			}
			if codigoPostal.Valid {
				userData["codigoPostal"] = codigoPostal.String
			}
			if ciudad.Valid {
				userData["ciudad"] = ciudad.String
			}
			if estado.Valid {
				userData["estado"] = estado.String
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userData)
	})))

	// Iniciar el servidor
	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
