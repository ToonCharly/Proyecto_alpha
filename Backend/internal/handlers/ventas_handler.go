package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// Constantes para reducir duplicación de literales
const (
	MetodoNoPermitido  = "Método no permitido"
	ErrorProcesarDatos = "Error al procesar los datos"
	ContentTypeHeader  = "Content-Type"
	ApplicationJSON    = "application/json"
	SerieMinLength     = 30
	ErrorBuscarVentas  = "Error al buscar ventas"
	ErrorBuscarPedido  = "Error al buscar pedido"
	ErrorBaseDatos     = "Error de base de datos"
)

// VentasHandler maneja las peticiones de consulta de ventas
func VentasHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, MetodoNoPermitido, http.StatusMethodNotAllowed)
			return
		}

		serie := r.URL.Query().Get("serie")
		if len(serie) < SerieMinLength {
			http.Error(w, "La serie debe tener al menos 30 caracteres", http.StatusBadRequest)
			return
		} // CONSULTA CON GROUP BY PARA EVITAR DUPLICADOS Y OBTENER CLAVES SAT REALES
		query := `
            SELECT 
                p.id_pedido, 
                p.clave_pedido, 
                d.descripcion AS producto, 
                SUM(d.cantidad) AS cantidad_total, 
                AVG(d.precio) AS precio_promedio, 
                AVG(d.iva) AS iva_pedido,                
                SUM(d.descuento) AS descuento_total,
                d.idproducto AS codigo_producto,
                COALESCE(pr.sat_clave, '01010101') AS sat_clave_real,
                COALESCE(pr.sat_medida, 'H87') AS sat_medida_real
            FROM optimus.crm_pedidos p 
            JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
            LEFT JOIN optimus.crm_productos pr ON d.idproducto = pr.idproducto
            WHERE p.clave_pedido = ?
            GROUP BY p.id_pedido, p.clave_pedido, d.idproducto, d.descripcion, pr.sat_clave, pr.sat_medida
            ORDER BY d.idproducto`

		rows, err := db.Query(query, serie)
		if err != nil {
			log.Printf("Error al buscar ventas: %v", err)
			http.Error(w, "Error al buscar ventas", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var ventas []map[string]interface{}
		// Usar mapa para deduplicar por codigo_producto
		productosUnicos := make(map[int]map[string]interface{})

		for rows.Next() {
			var idPedido int
			var clavePedido, producto string
			var codigoProducto int
			var cantidad, precio, ivaPedido, descuento float64
			var satClaveReal, satMedidaReal string

			if err := rows.Scan(
				&idPedido,
				&clavePedido,
				&producto,
				&cantidad,
				&precio,
				&ivaPedido,
				&descuento,
				&codigoProducto,
				&satClaveReal,
				&satMedidaReal,
			); err != nil {
				log.Printf("Error al escanear los resultados: %v", err)
				http.Error(w, ErrorProcesarDatos, http.StatusInternalServerError)
				return
			}

			// Solo procesar si no hemos visto este codigo_producto antes
			if _, existe := productosUnicos[codigoProducto]; !existe {
				// Calcular el total básico
				subtotal := (cantidad * precio) - descuento
				ivaTotal := subtotal * (ivaPedido / 100.0)
				totalConImpuestos := subtotal + ivaTotal

				productosUnicos[codigoProducto] = map[string]interface{}{
					"idPedido":        idPedido,
					"clavePedido":     clavePedido,
					"producto":        producto,
					"cantidad":        cantidad,
					"precio":          precio,
					"iva":             ivaPedido,
					"descuento":       descuento,
					"subtotal":        subtotal,
					"total":           totalConImpuestos,
					"codigo_producto": codigoProducto,
					"sat_clave":       satClaveReal,  // Usar valor real de la BD
					"sat_medida":      satMedidaReal, // Usar valor real de la BD
				}

				log.Printf("✅ PRODUCTO ÚNICO: ID=%d, Desc='%s', Cant=%.2f, Precio=%.2f, SAT_Clave='%s', SAT_Medida='%s'",
					codigoProducto, producto, cantidad, precio, satClaveReal, satMedidaReal)
			} else {
				log.Printf("⚠️ PRODUCTO DUPLICADO IGNORADO: ID=%d, Desc='%s'",
					codigoProducto, producto)
			}
		}

		// Convertir mapa a slice
		for _, producto := range productosUnicos {
			ventas = append(ventas, producto)
		}

		log.Printf("🔍 RESUMEN: Total filas SQL=%d, Productos únicos=%d",
			len(productosUnicos), len(ventas))
		w.Header().Set(ContentTypeHeader, ApplicationJSON)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ventas": ventas,
			"total":  len(ventas),
		})
	}
}

// DiagnosticoVentasHandler genera un diagnóstico detallado sobre los productos
func DiagnosticoVentasHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, MetodoNoPermitido, http.StatusMethodNotAllowed)
			return
		}

		serie := r.URL.Query().Get("serie")
		if len(serie) < SerieMinLength {
			http.Error(w, "La serie debe tener al menos 30 caracteres", http.StatusBadRequest)
			return
		}

		// Actualizar la consulta detallada para diagnóstico
		query := `
            SELECT 
                p.id_pedido, 
                p.clave_pedido, 
                d.idproducto AS id_producto_det,
                d.descripcion AS producto, 
                d.cantidad, 
                d.precio, 
                d.precio_o,           -- Nuevo campo
                d.iva,                -- Nuevo campo
                d.descuento, 
                pr.idproducto AS id_producto_crm,
                pr.clave AS clave_producto,
                pr.sat_clave,
                pr.sat_medida,                CASE 
                    WHEN pr.idproducto IS NULL THEN 'Producto no existe en crm_productos'
                    WHEN pr.sat_clave IS NULL OR pr.sat_clave = '' OR pr.sat_clave = '0' THEN 'Sin clave SAT'
                    WHEN pr.sat_medida IS NULL OR pr.sat_medida = '' THEN 'Sin unidad SAT'
                    WHEN pr.clave IS NULL OR pr.clave = '' THEN 'Sin clave de producto'
                    ELSE 'Configuración completa'
                END AS diagnostico
            FROM optimus.crm_pedidos p 
            JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
            LEFT JOIN optimus.crm_productos pr ON d.idproducto = pr.idproducto 
            WHERE p.clave_pedido = ?`

		rows, err := db.Query(query, serie)
		if err != nil {
			log.Printf("Error al buscar ventas: %v", err)
			http.Error(w, ErrorBuscarVentas, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var diagnostico []map[string]interface{}
		for rows.Next() {
			var idPedido int
			var idProductoDet int
			var idProductoCrm sql.NullInt64
			var clavePedido, producto, claveProducto, estadoDiagnostico string
			var satClave sql.NullString
			var satMedida sql.NullString
			var cantidad, precio, precioOriginal, iva, descuento float64

			if err := rows.Scan(
				&idPedido,
				&clavePedido,
				&idProductoDet,
				&producto,
				&cantidad,
				&precio,
				&precioOriginal, // Campo precio_o renombrado
				&iva,            // Nuevo campo
				&descuento,
				&idProductoCrm,
				&claveProducto,
				&satClave,
				&satMedida,
				&estadoDiagnostico,
			); err != nil {
				log.Printf("Error al escanear resultados diagnóstico: %v", err)
				http.Error(w, ErrorProcesarDatos, http.StatusInternalServerError)
				return
			}

			// Procesar valores nulos
			var idProductoCrmValue interface{}
			if idProductoCrm.Valid {
				idProductoCrmValue = idProductoCrm.Int64
			} else {
				idProductoCrmValue = nil
			} // Procesar los nuevos campos con lógica mejorada
			satClaveVal := "No configurado"
			if satClave.Valid && satClave.String != "" && satClave.String != "0" {
				satClaveVal = satClave.String
			}

			satMedidaVal := "No disponible"
			if satMedida.Valid && satMedida.String != "" {
				satMedidaVal = satMedida.String
			}

			// Calcular el total incluyendo el IVA
			subtotal := (cantidad * precio) - descuento
			totalConIva := subtotal + (subtotal * (iva / 100.0))

			diagnostico = append(diagnostico, map[string]interface{}{
				"idPedido":          idPedido,
				"clavePedido":       clavePedido,
				"idProductoDet":     idProductoDet,
				"producto":          producto,
				"cantidad":          cantidad,
				"precio":            precio,
				"precio_o":          precioOriginal, // Campo precio_o renombrado
				"iva":               iva,            // Nuevo campo en la respuesta
				"descuento":         descuento,
				"idProductoCrm":     idProductoCrmValue,
				"claveProducto":     claveProducto,
				"sat_clave":         satClaveVal,
				"sat_medida":        satMedidaVal,
				"estadoDiagnostico": estadoDiagnostico,
				"subtotal":          subtotal,
				"total":             totalConIva, // Ahora incluye IVA
			})
		}
		w.Header().Set(ContentTypeHeader, ApplicationJSON)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"diagnostico": diagnostico,
			"total":       len(diagnostico),
		})
	}
}

// InfoPedidoHandler obtiene la información básica del pedido (encabezado)
func InfoPedidoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, MetodoNoPermitido, http.StatusMethodNotAllowed)
			return
		}

		clave := r.URL.Query().Get("clave")
		if len(clave) < SerieMinLength {
			http.Error(w, "La clave debe tener al menos 30 caracteres", http.StatusBadRequest)
			return
		}

		query := `
            SELECT 
                id_pedido,
                estatus,
                tipo_docto,
                facturar
            FROM optimus.crm_pedidos
            WHERE clave_pedido = ?`

		var idPedido int
		var estatus, tipoDocto, facturar sql.NullString

		err := db.QueryRow(query, clave).Scan(&idPedido, &estatus, &tipoDocto, &facturar)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Pedido no encontrado", http.StatusNotFound)
				return
			}
			log.Printf("Error al buscar pedido: %v", err)
			http.Error(w, ErrorBuscarPedido, http.StatusInternalServerError)
			return
		}

		// Convertir valores nulos a strings vacíos
		estatusStr := ""
		if estatus.Valid {
			estatusStr = estatus.String
		}

		tipoDoctoStr := ""
		if tipoDocto.Valid {
			tipoDoctoStr = tipoDocto.String
		}

		facturarStr := ""
		if facturar.Valid {
			facturarStr = facturar.String
		}

		pedidoInfo := map[string]interface{}{
			"id_pedido":  idPedido,
			"estatus":    estatusStr,
			"tipo_docto": tipoDoctoStr,
			"facturar":   facturarStr,
		}
		w.Header().Set(ContentTypeHeader, ApplicationJSON)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pedido": pedidoInfo,
		})
	}
}

// GuardarVentasHandler almacena las ventas en la tabla ventas_det
func GuardarVentasHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Configurar CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", ContentTypeHeader)

		// Manejar preflight OPTIONS
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, MetodoNoPermitido, http.StatusMethodNotAllowed)
			return
		}

		// Estructura para recibir los datos de ventas según la nueva estructura
		var ventasData struct {
			Serie  string `json:"serie"`
			Ventas []struct {
				ClaveProducto  string  `json:"clave_producto"`
				Descripcion    string  `json:"descripcion"`
				ClaveSat       string  `json:"clave_sat"`
				UnidadSat      string  `json:"unidad_sat"`
				Cantidad       int     `json:"cantidad"`
				PrecioUnitario float64 `json:"precio_unitario"`
				Descuento      float64 `json:"descuento"`
				Total          float64 `json:"total"`
				IVA            float64 `json:"iva"`
			} `json:"ventas"`
		}

		// Decodificar JSON
		if err := json.NewDecoder(r.Body).Decode(&ventasData); err != nil {
			log.Printf("Error al decodificar JSON: %v", err)
			http.Error(w, ErrorProcesarDatos, http.StatusBadRequest)
			return
		}

		// Validar que hay datos
		if len(ventasData.Ventas) == 0 {
			http.Error(w, "No hay datos de ventas para guardar", http.StatusBadRequest)
			return
		}

		// Iniciar transacción
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Error al iniciar transacción: %v", err)
			http.Error(w, ErrorBaseDatos, http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Preparar statement para insertar en ventas_det con la nueva estructura
		stmt, err := tx.Prepare(`
			INSERT INTO ventas_det (
				serie, clave_producto, descripcion, clave_sat, unidad_sat, 
				cantidad, precio_unitario, descuento, total, iva, fecha_venta
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		`)
		if err != nil {
			log.Printf("Error al preparar statement: %v", err)
			http.Error(w, ErrorBaseDatos, http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		// Insertar cada venta
		insertados := 0
		for _, venta := range ventasData.Ventas {
			// Si no se proporciona IVA, usar el valor por defecto de 16%
			iva := venta.IVA
			if iva == 0 {
				iva = 16.0
			}

			_, err := stmt.Exec(
				ventasData.Serie,    // Serie como primer parámetro
				venta.ClaveProducto, // Agregar clave_producto
				venta.Descripcion,
				venta.ClaveSat,
				venta.UnidadSat,
				venta.Cantidad,
				venta.PrecioUnitario,
				venta.Descuento,
				venta.Total,
				iva, // Agregar IVA
			)
			if err != nil {
				log.Printf("Error al insertar venta: %v", err)
				http.Error(w, "Error al guardar venta", http.StatusInternalServerError)
				return
			}
			insertados++
		}

		// Confirmar transacción
		if err := tx.Commit(); err != nil {
			log.Printf("Error al confirmar transacción: %v", err)
			http.Error(w, "Error al confirmar guardado", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Guardadas %d ventas en ventas_det para serie: %s", insertados, ventasData.Serie)
		// Responder con éxito
		w.Header().Set(ContentTypeHeader, ApplicationJSON)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "success",
			"message":    "Ventas guardadas correctamente",
			"insertados": insertados,
			"serie":      ventasData.Serie,
		})
	}
}
