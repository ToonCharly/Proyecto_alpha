package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	// Iniciar el servidor
	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
