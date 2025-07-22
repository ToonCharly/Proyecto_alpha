package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Facts/internal/db"
	"Facts/internal/handlers"
	"Facts/internal/models"
	"Facts/internal/utils"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	http.Handle("/api/facturas-empresa-activa", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		idUsuarioStr := r.URL.Query().Get("id_usuario")
		if idUsuarioStr == "" {
			log.Printf("[facturas-empresa-activa] Falta parámetro id_usuario")
			http.Error(w, "El parámetro id_usuario es requerido", http.StatusBadRequest)
			return
		}
		idUsuario, err := strconv.Atoi(idUsuarioStr)
		if err != nil {
			log.Printf("[facturas-empresa-activa] id_usuario no es entero: %v", idUsuarioStr)
			http.Error(w, "El parámetro id_usuario debe ser un número entero", http.StatusBadRequest)
			return
		}

		log.Printf("[facturas-empresa-activa] Buscando empresa activa para usuario ID: %d", idUsuario)
		var idDatosFiscales int
		err = db.GetDB().QueryRow("SELECT id FROM datos_fiscales WHERE id_usuario = ? AND activo = 1 LIMIT 1", idUsuario).Scan(&idDatosFiscales)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("[facturas-empresa-activa] No hay empresa activa para usuario %d", idUsuario)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode([]interface{}{})
				return
			}
			log.Printf("[facturas-empresa-activa] Error SQL al obtener idDatosFiscales para usuario %d: %v", idUsuario, err)
			http.Error(w, "Error al buscar la empresa activa: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("[facturas-empresa-activa] idDatosFiscales encontrado: %d", idDatosFiscales)
		rows, err := db.GetDB().Query("SELECT * FROM facturas WHERE id_datos_fiscales = ?", idDatosFiscales)
		if err != nil {
			log.Printf("[facturas-empresa-activa] Error SQL al consultar facturas para id_datos_fiscales %d: %v", idDatosFiscales, err)
			http.Error(w, "Error al buscar facturas: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			log.Printf("[facturas-empresa-activa] Error al obtener columnas: %v", err)
			http.Error(w, "Error al procesar resultados: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var facturas []map[string]interface{}
		for rows.Next() {
			values := make([]interface{}, len(cols))
			valuePtrs := make([]interface{}, len(cols))
			for i := range values {
				valuePtrs[i] = &values[i]
			}
			if err := rows.Scan(valuePtrs...); err != nil {
				log.Printf("[facturas-empresa-activa] Error al escanear fila: %v", err)
				continue
			}
			factura := make(map[string]interface{})
			for i, col := range cols {
				var v interface{}
				val := values[i]
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
				factura[col] = v
			}
			facturas = append(facturas, factura)
		}

		log.Printf("[facturas-empresa-activa] Facturas encontradas: %d", len(facturas))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(facturas)
	})))
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

	// Configurar FileServer para archivos estáticos
	fs := http.FileServer(http.Dir("./public/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Los logos ahora se manejan únicamente desde la base de datos
	// Eliminado el servicio de archivos estáticos para logos

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

	// Endpoint robusto único para generar factura
	http.Handle("/api/generar-factura", utils.EnableCors(http.HandlerFunc(handlers.GenerarFacturaHandler)))

	http.HandleFunc("/api/historial-emisor", handlers.HistorialEmisorHandler)

	// Endpoint que devuelve información sobre la factura generada
	http.Handle("/api/generar-factura-info", utils.EnableCors(http.HandlerFunc(handlers.GenerarFacturaConInfoHandler)))

	// Test endpoints removed - handlers not implemented

	// Endpoint para debug - ver qué series existen en ventas_det
	http.Handle("/api/debug-ventas", utils.EnableCors(http.HandlerFunc(handlers.DebugVentasHandler)))

	// Endpoint para debug - ver estructura y datos de ventas_det
	http.Handle("/api/debug-ventas-det", utils.EnableCors(http.HandlerFunc(handlers.DebugVentasDetHandler)))

	http.Handle("/api/plantillas/subir", utils.EnableCors(http.HandlerFunc(handlers.SubirPlantillaHandler)))
	http.Handle("/api/plantillas/listar", utils.EnableCors(http.HandlerFunc(handlers.ListarPlantillasHandler)))
	http.Handle("/api/plantillas/activar", utils.EnableCors(http.HandlerFunc(handlers.ActivarPlantillaHandler)))
	http.Handle("/api/plantillas/eliminar", utils.EnableCors(http.HandlerFunc(handlers.EliminarPlantillaHandler)))
	http.Handle("/api/plantillas", utils.EnableCors(http.HandlerFunc(handlers.BuscarPlantillasHandler)))

	// Endpoints para manejo de logos en base de datos
	logoDBHandler := handlers.NewLogoHandler()
	http.Handle("/api/logos/subir", utils.EnableCors(http.HandlerFunc(logoDBHandler.SubirLogo)))
	http.Handle("/api/logos/listar", utils.EnableCors(http.HandlerFunc(logoDBHandler.ObtenerLogosUsuario)))
	http.Handle("/api/logos/obtener-activo", utils.EnableCors(http.HandlerFunc(logoDBHandler.ObtenerLogoActivoUsuario)))
	http.Handle("/api/logos/obtener-activo-json", utils.EnableCors(http.HandlerFunc(logoDBHandler.ObtenerLogoActivoJSON)))
	http.Handle("/api/logos/descargar", utils.EnableCors(http.HandlerFunc(logoDBHandler.ObtenerLogoImagen)))
	http.Handle("/api/logos/activar", utils.EnableCors(http.HandlerFunc(logoDBHandler.ActivarLogo)))
	http.Handle("/api/logos/eliminar", utils.EnableCors(http.HandlerFunc(logoDBHandler.EliminarLogo)))
	http.Handle("/api/logos/debug", utils.EnableCors(http.HandlerFunc(handlers.DebugLogosHandler)))

	// Usar la conexión a optimus para los handlers que la necesitan
	http.Handle("/api/ventas", utils.EnableCors(http.HandlerFunc(handlers.VentasHandler(optimusDB))))

	// Endpoint para guardar ventas en la tabla ventas_det
	http.Handle("/api/ventas/guardar", utils.EnableCors(http.HandlerFunc(handlers.GuardarVentasHandler(db.GetDB()))))

	http.Handle("/api/historial_facturas", utils.EnableCors(http.HandlerFunc(handlers.HistorialFacturasHandler(db.GetDB()))))

	// Endpoint para búsqueda en historial de facturas
	http.Handle("/api/buscar-facturas", utils.EnableCors(http.HandlerFunc(handlers.BuscarHistorialFacturasHandler(db.GetDB()))))

	// Endpoint para obtener regímenes fiscales
	http.Handle("/api/regimenes-fiscales", utils.EnableCors(http.HandlerFunc(handlers.GetRegimenesFiscales)))
	// Usar la conexión a optimus para el diagnóstico
	http.Handle("/api/optimus/diagnostico", utils.EnableCors(http.HandlerFunc(handlers.DiagnosticoVentasHandler(optimusDB))))

	// Endpoints para impuestos
	http.Handle("/api/impuestos", utils.EnableCors(http.HandlerFunc(handlers.ImpuestosHandler(optimusDB))))
	http.Handle("/api/productos-con-impuestos", utils.EnableCors(http.HandlerFunc(handlers.ProductosConImpuestosHandler(optimusDB))))

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
	http.Handle("/api/historial-facturas", utils.EnableCors(http.HandlerFunc(handlers.HistorialFacturasHandler(db.GetDB()))))

	// Endpoint para descargar una factura del historial
	// ...existing code...

	// Añadir nuevos endpoints (corregido)
	http.Handle("/api/reset-password-request", utils.EnableCors(http.HandlerFunc(handlers.ResetPasswordRequestHandler(db.GetDB()))))
	http.Handle("/api/reset-password", utils.EnableCors(http.HandlerFunc(handlers.ResetPasswordHandler(db.GetDB()))))

	// Endpoint para obtener todos los usuarios (solo para admins)
	http.Handle("/api/usuarios", utils.EnableCors(http.HandlerFunc(handlers.GetAllUsersHandler)))

	// Endpoint temporal para verificar tablas disponibles
	http.Handle("/api/verificar-tablas", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Conectar a la base de datos
		database, err := db.ConnectToAlpha()
		if err != nil {
			log.Printf("Error al conectar a la base de datos: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}
		defer database.Close()

		// Consultar las tablas disponibles
		rows, err := database.Query("SHOW TABLES")
		if err != nil {
			log.Printf("Error al consultar tablas: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tablas []string
		for rows.Next() {
			var tabla string
			if err := rows.Scan(&tabla); err != nil {
				log.Printf("Error al escanear tabla: %v", err)
				continue
			}
			tablas = append(tablas, tabla)
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"tablas": tablas,
		})
	})))

	// Endpoint para buscar empresa por RFC en adm_empresas_rfc
	http.Handle("/api/buscar-empresa-rfc", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		rfc := r.URL.Query().Get("rfc")
		if rfc == "" {
			http.Error(w, "El parámetro RFC es requerido", http.StatusBadRequest)
			return
		}

		// Buscar en la tabla adm_empresas_rfc con JOIN a adm_metodopago, efac_regimenfiscal, adm_tipopagos y adm_condicionpago
		var result struct {
			IDEmpresa          int    `json:"idempresa"`
			IDRFC              int    `json:"idrfc"`
			IDMetodo           int    `json:"idmetodo"`
			MetodoPago         string `json:"metodo_pago"`
			IDRegimenFiscal    int    `json:"idregimenfiscal"`
			CRegimenFiscal     string `json:"c_regimenfiscal"`
			DescripcionRegimen string `json:"descripcion_regimen"`
			IDTipoPago         int    `json:"idtipopago"`
			TipoPago           string `json:"tipo_pago"`
			IDCondicion        int    `json:"idcondicion"`
			CondicionPago      string `json:"condicion_pago"`
		}

		// Conectar a la base de datos Alpha
		alphaDB, err := db.ConnectToAlpha()
		if err != nil {
			log.Printf("Error al conectar a la base de datos Alpha: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}
		defer alphaDB.Close()

		query := `SELECT er.idempresa, er.idrfc, er.idmetodo, 
				  COALESCE(mp.metodo, 'No disponible') AS metodo_pago,
				  er.idregimenfiscal,
				  COALESCE(rf.c_regimenfiscal, 'No disponible') AS c_regimenfiscal,
				  COALESCE(rf.descripcion, 'No disponible') AS descripcion_regimen,
				  COALESCE(er.idtipopago, 0) AS idtipopago,
				  COALESCE(tp.tipo, 'No disponible') AS tipo_pago,
				  COALESCE(er.idcondicion, 0) AS idcondicion,
				  COALESCE(cp.condicion, 'No disponible') AS condicion_pago
				  FROM adm_empresas_rfc er
				  LEFT JOIN adm_metodopago mp ON er.idmetodo = mp.idmetodo
				  LEFT JOIN efac_regimenfiscal rf ON er.idregimenfiscal = rf.idregimenfiscal
				  LEFT JOIN adm_tipopagos tp ON er.idtipopago = tp.idtipo
				  LEFT JOIN adm_condicionpago cp ON er.idcondicion = cp.idcondicion
				  WHERE er.rfc = ?`
		err = alphaDB.QueryRow(query, rfc).Scan(&result.IDEmpresa, &result.IDRFC, &result.IDMetodo, &result.MetodoPago, &result.IDRegimenFiscal, &result.CRegimenFiscal, &result.DescripcionRegimen, &result.IDTipoPago, &result.TipoPago, &result.IDCondicion, &result.CondicionPago)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "RFC no encontrado", http.StatusNotFound)
				return
			}
			log.Printf("Error al buscar RFC: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, result)
	})))

	// Endpoint para timbrar factura CFDI 4.0
	http.Handle("/api/timbrar-factura", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Leer datos de la factura desde el body
		var facturaData models.FacturaCFDI
		if err := json.NewDecoder(r.Body).Decode(&facturaData); err != nil {
			log.Printf("Error decodificando JSON de factura: %v", err)
			utils.RespondWithError(w, "Error al leer los datos de la factura")
			return
		}

		// Lógica de timbrado (debe estar implementada en tu módulo PAC o handlers)
		resultadoTimbrado, err := handlers.TimbrarFacturaHandler(facturaData)
		if err != nil {
			log.Printf("Error al timbrar factura: %v", err)
			utils.RespondWithError(w, fmt.Sprintf("Error al timbrar la factura: %v", err))
			return
		}

		// Guardar el resultado en la base de datos
		err = db.GuardarFacturaTimbrada(db.GetDB(), resultadoTimbrado)
		if err != nil {
			log.Printf("Error al guardar factura timbrada: %v", err)
			utils.RespondWithError(w, "Error al guardar la factura timbrada")
			return
		}

		// Responder con el XML/PDF/timbre fiscal digital
		utils.RespondWithJSON(w, http.StatusOK, resultadoTimbrado)
	})))

	// Endpoint para obtener detalles de empresa por ID desde adm_empresa
	http.Handle("/api/empresa-detalle", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		idEmpresaStr := r.URL.Query().Get("idempresa")
		if idEmpresaStr == "" {
			http.Error(w, "El parámetro idempresa es requerido", http.StatusBadRequest)
			return
		}

		idEmpresa, err := strconv.Atoi(idEmpresaStr)
		if err != nil {
			http.Error(w, "El parámetro idempresa debe ser un número", http.StatusBadRequest)
			return
		}

		// Conectar a la base de datos Alpha
		alphaDB, err := db.ConnectToAlpha()
		if err != nil {
			log.Printf("Error al conectar a la base de datos Alpha: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}
		defer alphaDB.Close()

		var empresa struct {
			IDEmpresa       int    `json:"idempresa"`
			NombreComercial string `json:"nombre_comercial"`
			RazonSocial     string `json:"razon_social"`
			RFC             string `json:"rfc"`
			Direccion1      string `json:"direccion1"`
			Colonia         string `json:"colonia"`
			CP              string `json:"cp"`
			Ciudad          string `json:"ciudad"`
			Estado          string `json:"estado"`
		}

		query := `SELECT e.idempresa, 
			  COALESCE(e.nombre_comercial, 'No disponible') AS nombre_comercial,
			  COALESCE(e.razon_social, 'No disponible') AS razon_social,
			  COALESCE(e.rfc, 'No disponible') AS rfc,
			  COALESCE(e.direccion1, 'No disponible') AS direccion1,
			  COALESCE(e.colonia, 'No disponible') AS colonia,
			  COALESCE(e.cp, 'No disponible') AS cp,
			  COALESCE(e.ciudad, 'No disponible') AS ciudad,
			  COALESCE(est.estado, 'No disponible') AS estado
			  FROM adm_empresas e
			  LEFT JOIN adm_estados_mex est ON e.estado = est.idestado
			  WHERE e.idempresa = ?`

		err = alphaDB.QueryRow(query, idEmpresa).Scan(
			&empresa.IDEmpresa,
			&empresa.NombreComercial,
			&empresa.RazonSocial,
			&empresa.RFC,
			&empresa.Direccion1,
			&empresa.Colonia,
			&empresa.CP,
			&empresa.Ciudad,
			&empresa.Estado,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Empresa no encontrada", http.StatusNotFound)
				return
			}
			log.Printf("Error al obtener detalles de empresa: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, empresa)
	})))

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

		// Obtener datos fiscales para el usuario específico
		query := `SELECT id, rfc, razon_social, direccion_fiscal, direccion, 
		  codigo_postal, clave_csd, regimen_fiscal,
		  nombre_comercial, colonia, ciudad, estado, serie_df
		  FROM datos_fiscales WHERE id_usuario = ?`

		var datosFiscales struct {
			ID              int    `json:"id"`
			RFC             string `json:"rfcEmisor"`
			RazonSocial     string `json:"razonSocial"`
			DireccionFiscal string `json:"direccionFiscal"`
			Direccion       string `json:"direccion"`
			CodigoPostal    string `json:"cp"`
			// RutaCsdKey and RutaCsdCer removed (now handled as binary in DB)
			ClaveCsd        string `json:"claveArchivoCSD"`
			RegimenFiscal   string `json:"regimenFiscal"`
			NombreComercial string `json:"nombreComercial"`
			Colonia         string `json:"colonia"`
			Ciudad          string `json:"ciudad"`
			Estado          string `json:"estado"`
			SerieDf         string `json:"serieDf"`
		}

		err = db.GetDB().QueryRow(query, idUsuario).Scan(
			&datosFiscales.ID,
			&datosFiscales.RFC,
			&datosFiscales.RazonSocial,
			&datosFiscales.DireccionFiscal,
			&datosFiscales.Direccion,
			&datosFiscales.CodigoPostal,
			&datosFiscales.ClaveCsd,
			&datosFiscales.RegimenFiscal,
			&datosFiscales.NombreComercial,
			&datosFiscales.Colonia,
			&datosFiscales.Ciudad,
			&datosFiscales.Estado,
			&datosFiscales.SerieDf,
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
		nombreComercial := r.FormValue("nombre_comercial")
		direccionFiscal := r.FormValue("direccion_fiscal")
		direccion := r.FormValue("direccion")
		colonia := r.FormValue("colonia")
		ciudad := r.FormValue("ciudad")
		estado := r.FormValue("estado")
		codigoPostal := r.FormValue("codigo_postal")
		claveCsd := r.FormValue("clave_csd")
		regimenFiscal := r.FormValue("regimen_fiscal")
		serieDf := r.FormValue("serie_df")
		idUsuarioStr := r.FormValue("id_usuario")

		// LOG para depuración
		log.Printf("🔍 MAIN.GO - Datos recibidos:")
		log.Printf("   rfc: '%s'", rfc)
		log.Printf("   nombre_comercial: '%s'", nombreComercial)
		log.Printf("   direccion: '%s'", direccion)
		log.Printf("   colonia: '%s'", colonia)
		log.Printf("   ciudad: '%s'", ciudad)
		log.Printf("   estado: '%s'", estado)
		log.Printf("   serie_df: '%s' (longitud: %d)", serieDf, len(serieDf))

		// Validar campos obligatorios
		if rfc == "" || razonSocial == "" || codigoPostal == "" || regimenFiscal == "" || idUsuarioStr == "" {
			log.Printf("Faltan campos obligatorios: RFC=%s, RazonSocial=%s, CP=%s, RegimenFiscal=%s, IDUsuario=%s",
				rfc, razonSocial, codigoPostal, regimenFiscal, idUsuarioStr)
			http.Error(w, "RFC, Razón Social, Código Postal, Régimen Fiscal e ID de Usuario son obligatorios", http.StatusBadRequest)
			return
		}

		// Convertir id_usuario a entero
		idUsuario, err := strconv.Atoi(idUsuarioStr)
		if err != nil {
			log.Printf("Error al convertir id_usuario: %v", err)
			http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
			return
		}

		// Leer archivos del form (si existen)
		var archivoCSDKey, archivoCSDCer []byte
		var nombreArchivoKey, nombreArchivoCer string

		csdKeyFile, fileCSDKeyHeader, err := r.FormFile("csdKey")
		if err == nil && csdKeyFile != nil {
			defer csdKeyFile.Close()
			archivoCSDKey, _ = io.ReadAll(csdKeyFile)
			if fileCSDKeyHeader != nil {
				nombreArchivoKey = fileCSDKeyHeader.Filename
			}
		}

		csdCerFile, fileCSDCerHeader, err := r.FormFile("csdCer")
		if err == nil && csdCerFile != nil {
			defer csdCerFile.Close()
			archivoCSDCer, _ = io.ReadAll(csdCerFile)
			if fileCSDCerHeader != nil {
				nombreArchivoCer = fileCSDCerHeader.Filename
			}
		}

		// Llama a la función de tu capa db que guarda todo correctamente (incluyendo archivos)
		_, err = db.GuardarDatosFiscales(
			rfc, razonSocial, direccionFiscal, codigoPostal,
			archivoCSDKey, archivoCSDCer,
			nombreArchivoKey, nombreArchivoCer,
			claveCsd, regimenFiscal, serieDf,
			idUsuario,
			"", // keyPath (ruta al archivo .key), se pasa vacío por defecto
		)
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

	// Añadir este nuevo endpoint para diagnóstico del historial

	// Endpoint para diagnosticar problemas con el historial
	http.Handle("/api/diagnostico-historial", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Verificar estructura de la tabla
		var tablaInfo struct {
			Columnas []string                 `json:"columnas"`
			Filas    int                      `json:"filas"`
			Muestra  []map[string]interface{} `json:"muestra"`
		}

		// Obtener columnas
		columnas, err := db.GetDB().Query("SHOW COLUMNS FROM historial_facturas")
		if err != nil {
			log.Printf("Error al obtener estructura: %v", err)
			http.Error(w, "Error al verificar estructura de tabla", http.StatusInternalServerError)
			return
		}
		defer columnas.Close()

		for columnas.Next() {
			var field, tipo, nulo, key, default_val, extra string
			if err := columnas.Scan(&field, &tipo, &nulo, &key, &default_val, &extra); err != nil {
				continue
			}
			tablaInfo.Columnas = append(tablaInfo.Columnas, field)
		}

		// Contar filas
		var count int
		err = db.GetDB().QueryRow("SELECT COUNT(*) FROM historial_facturas").Scan(&count)
		if err != nil {
			log.Printf("Error al contar filas: %v", err)
			tablaInfo.Filas = -1
		} else {
			tablaInfo.Filas = count
		}

		// Obtener muestra de hasta 5 filas
		if count > 0 {
			rows, err := db.GetDB().Query("SELECT * FROM historial_facturas ORDER BY fecha_emision DESC LIMIT 5")
			if err == nil {
				defer rows.Close()

				// Obtener nombres de columnas
				cols, err := rows.Columns()
				if err == nil {
					for rows.Next() {
						// Crear slice para almacenar valores de columnas
						values := make([]interface{}, len(cols))
						valuePtrs := make([]interface{}, len(cols))
						for i := range values {
							valuePtrs[i] = &values[i]
						}

						if err := rows.Scan(valuePtrs...); err != nil {
							continue
						}

						// Crear mapa para esta fila
						fila := make(map[string]interface{})
						for i, col := range cols {
							var v interface{}
							val := values[i]

							b, ok := val.([]byte)
							if ok {
								v = string(b)
							} else {
								v = val
							}
							fila[col] = v
						}
						tablaInfo.Muestra = append(tablaInfo.Muestra, fila)
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tablaInfo)
	})))

	// Endpoint para restablecer datos fiscales
	http.Handle("/api/restablecer-datos-fiscales", utils.EnableCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodDelete {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Verificar que sea un administrador (opcional - la validación también está en frontend)
		var body struct {
			IDUsuario int `json:"id_usuario"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
			return
		}

		// Verificar si el usuario es administrador
		var isAdmin bool
		err := db.GetDB().QueryRow("SELECT role = 'admin' FROM usuarios WHERE id = ?", body.IDUsuario).Scan(&isAdmin)
		if err != nil {
			log.Printf("Error al verificar rol de usuario: %v", err)
			http.Error(w, "Error al verificar permisos de usuario", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			http.Error(w, "Solo los administradores pueden restablecer datos fiscales", http.StatusForbidden)
			return
		}

		// Conectar a la base de datos
		dbConn, err := db.ConnectUserDB()
		if err != nil {
			log.Printf("Error al conectar a la base de datos: %v", err)
			http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
			return
		}
		defer dbConn.Close()

		// Eliminar los datos fiscales del usuario específico
		_, err = dbConn.Exec("DELETE FROM datos_fiscales WHERE id_usuario = ?", body.IDUsuario)
		if err != nil {
			log.Printf("Error al eliminar datos fiscales: %v", err)
			http.Error(w, "Error al restablecer los datos fiscales", http.StatusInternalServerError)
			return
		}

		log.Printf("Datos fiscales restablecidos por usuario administrador ID: %d", body.IDUsuario)

		// Responder éxito
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Datos fiscales restablecidos correctamente",
		})
	})))

	// En main.go asegúrate de tener esta ruta
	http.HandleFunc("/api/plantillas/ejemplo", handlers.PlantillaEjemploHandler)

	// Iniciar el servidor
	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
