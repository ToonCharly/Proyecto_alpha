package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// VentasHandler maneja las peticiones de consulta de ventas
func VentasHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		serie := r.URL.Query().Get("serie")
		if len(serie) < 30 {
			http.Error(w, "La serie debe tener al menos 30 caracteres", http.StatusBadRequest)
			return
		}

		// Actualizar la consulta SQL para incluir precio_o e iva
		query := `
            SELECT 
                p.id_pedido, 
                p.clave_pedido, 
                d.descripcion AS producto, 
                d.cantidad, 
                d.precio, 
                d.precio_o,           -- Nuevo campo
                d.iva,                -- Nuevo campo
                d.descuento,
                d.idproducto AS codigo_producto, 
                pr.clave AS categoria_producto,
                pr.sat_clave,           
                pr.sat_medida           
            FROM optimus.crm_pedidos p 
            JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
            LEFT JOIN optimus.crm_productos pr ON d.idproducto = pr.idproducto 
            WHERE p.clave_pedido = ?`

		rows, err := db.Query(query, serie)
		if err != nil {
			log.Printf("Error al buscar ventas: %v", err)
			http.Error(w, "Error al buscar ventas", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var ventas []map[string]interface{}
		for rows.Next() {
			var idPedido int
			var clavePedido, producto string
			var codigoProducto int
			var categoriaProducto sql.NullString
			var satClave sql.NullString
			var satMedida sql.NullString
			var cantidad, precio, precio_o, iva, descuento float64 // Añadido precio_o e iva

			if err := rows.Scan(
				&idPedido,
				&clavePedido,
				&producto,
				&cantidad,
				&precio,
				&precio_o, // Nuevo campo
				&iva,      // Nuevo campo
				&descuento,
				&codigoProducto,
				&categoriaProducto,
				&satClave,
				&satMedida,
			); err != nil {
				log.Printf("Error al escanear los resultados: %v", err)
				http.Error(w, "Error al procesar los datos", http.StatusInternalServerError)
				return
			}

			// Usar valores por defecto cuando los campos son NULL
			categoriaProductoVal := "N/A"
			if categoriaProducto.Valid {
				categoriaProductoVal = categoriaProducto.String
			}

			satClaveVal := ""
			if satClave.Valid {
				satClaveVal = satClave.String
			}

			satMedidaVal := ""
			if satMedida.Valid {
				satMedidaVal = satMedida.String
			}

			// Calcular el total incluyendo el IVA
			subtotal := (cantidad * precio) - descuento
			totalConIva := subtotal + (subtotal * (iva / 100.0))

			ventas = append(ventas, map[string]interface{}{
				"idPedido":           idPedido,
				"clavePedido":        clavePedido,
				"producto":           producto,
				"cantidad":           cantidad,
				"precio":             precio,
				"precio_o":           precio_o, // Nuevo campo en la respuesta
				"iva":                iva,      // Nuevo campo en la respuesta
				"descuento":          descuento,
				"subtotal":           subtotal,
				"total":              totalConIva, // Ahora incluye IVA
				"codigo_producto":    codigoProducto,
				"categoria_producto": categoriaProductoVal,
				"sat_clave":          satClaveVal,
				"sat_medida":         satMedidaVal,
			})
		}

		w.Header().Set("Content-Type", "application/json")
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
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		serie := r.URL.Query().Get("serie")
		if len(serie) < 30 {
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
                pr.sat_medida,
                CASE 
                    WHEN pr.idproducto IS NULL THEN 'Producto no existe en crm_productos'
                    WHEN pr.clave IS NULL OR pr.clave = '' THEN 'Clave vacía o nula'
                    ELSE 'Clave disponible'
                END AS diagnostico
            FROM optimus.crm_pedidos p 
            JOIN optimus.crm_pedidos_det d ON p.id_pedido = d.id_pedido 
            LEFT JOIN optimus.crm_productos pr ON d.idproducto = pr.idproducto 
            WHERE p.clave_pedido = ?`

		rows, err := db.Query(query, serie)
		if err != nil {
			log.Printf("Error al buscar ventas: %v", err)
			http.Error(w, "Error al buscar ventas", http.StatusInternalServerError)
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
			var cantidad, precio, precio_o, iva, descuento float64 // Añadido precio_o e iva

			if err := rows.Scan(
				&idPedido,
				&clavePedido,
				&idProductoDet,
				&producto,
				&cantidad,
				&precio,
				&precio_o, // Nuevo campo
				&iva,      // Nuevo campo
				&descuento,
				&idProductoCrm,
				&claveProducto,
				&satClave,
				&satMedida,
				&estadoDiagnostico,
			); err != nil {
				log.Printf("Error al escanear resultados diagnóstico: %v", err)
				http.Error(w, "Error al procesar los datos", http.StatusInternalServerError)
				return
			}

			// Procesar valores nulos
			var idProductoCrmValue interface{}
			if idProductoCrm.Valid {
				idProductoCrmValue = idProductoCrm.Int64
			} else {
				idProductoCrmValue = nil
			}

			// Procesar los nuevos campos
			satClaveVal := "N/A"
			if satClave.Valid {
				satClaveVal = satClave.String
			}

			satMedidaVal := "N/A"
			if satMedida.Valid {
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
				"precio_o":          precio_o, // Nuevo campo en la respuesta
				"iva":               iva,      // Nuevo campo en la respuesta
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

		w.Header().Set("Content-Type", "application/json")
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
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		clave := r.URL.Query().Get("clave")
		if len(clave) < 30 {
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
			http.Error(w, "Error al buscar pedido", http.StatusInternalServerError)
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pedido": pedidoInfo,
		})
	}
}
