package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"carlos/Facts/Backend/internal/db"
	"carlos/Facts/Backend/internal/models"
)

func BuscarFactura(db *sql.DB, w http.ResponseWriter, criterio string) {
	log.Println("Criterio recibido:", criterio)
	query := `SELECT id as idfactura, id_usuario as idempresa, rfc_receptor as rfc, razon_social_receptor as razon_social, 
total as subtotal, 0 as impuestos, estado as estatus, '' as pagado, fecha_generacion as fecha_pago 
FROM historial_facturas WHERE rfc_receptor LIKE ? OR razon_social_receptor LIKE ? OR folio LIKE ? LIMIT 1`
	likeCriterio := "%" + criterio + "%"
	row := db.QueryRow(query, likeCriterio, likeCriterio, likeCriterio)

	var f models.Factura
	err := row.Scan(&f.IdFactura, &f.IdEmpresa, &f.RFC, &f.RazonSocial, &f.Subtotal, &f.Impuestos, &f.Estatus, &f.Pagado, &f.FechaPago)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err == sql.ErrNoRows {
			log.Printf("No se encontró factura con criterio: %s", criterio)
			json.NewEncoder(w).Encode(map[string]string{"error": "Empresa no encontrada"})
		} else {
			log.Printf("Error al buscar factura con criterio '%s': %v", criterio, err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar la base de datos"})
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f)
}

func OtraFuncion(db *sql.DB, w http.ResponseWriter, id int, empresaID int, estado int, pagado int, fechaPago string) {
	sqlQuery := fmt.Sprintf("WHERE f.ID = %d AND f.EmpresaID = %d AND f.Estado = %d AND f.Pagado = %d AND f.FechaPago = '%s'",
		id, empresaID, estado, pagado, fechaPago)

	fullQuery := "SELECT * FROM facturas " + sqlQuery
	rows, err := db.Query(fullQuery)
	if err != nil {
		log.Println("Error en la consulta:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error en la consulta"})
		return
	}
	defer rows.Close()

	var facturas []models.Factura
	for rows.Next() {
		var f models.Factura
		if err := rows.Scan(&f.IdFactura, &f.IdEmpresa, &f.RFC, &f.RazonSocial, &f.Subtotal, &f.Impuestos, &f.Estatus, &f.Pagado, &f.FechaPago); err != nil {
			log.Println("Error al escanear fila:", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error al procesar resultados"})
			return
		}
		facturas = append(facturas, f)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(facturas)
}

// Nueva función para manejar la generación de facturas
func GenerarFactura(w http.ResponseWriter, r *http.Request) {
	// Aquí iría la lógica para generar la factura
	// Por ahora, solo devolvemos una respuesta de éxito simulada

	// Simulamos un ID de usuario, esto debería ser parte de la lógica de tu aplicación
	facturaRequest := struct {
		IDUsuario int
	}{IDUsuario: 123}

	nombreArchivo := fmt.Sprintf("factura_%d_%s.pdf", facturaRequest.IDUsuario, time.Now().Format("20060102150405"))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Factura generada correctamente",
		"url":     "/facturas/descargar/" + nombreArchivo,
	})
}

// GenerarFacturaHandler maneja la generación de facturas usando plantillas
func GenerarFacturaDesdeDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar datos de la factura
	var facturaRequest struct {
		IDUsuario int    `json:"id_usuario"`
		Empresa   string `json:"empresa"`
		RFC       string `json:"rfc"`
		Conceptos []struct {
			Descripcion string  `json:"descripcion"`
			Cantidad    int     `json:"cantidad"`
			PrecioUni   float64 `json:"precio_unitario"`
		} `json:"conceptos"`
	}

	err := json.NewDecoder(r.Body).Decode(&facturaRequest)
	if err != nil {
		log.Printf("Error al decodificar solicitud: %v", err)
		http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
		return
	}

	// Crear directorio para facturas si no existe
	facturasDir := "./public/facturas"
	if _, err := os.Stat(facturasDir); os.IsNotExist(err) {
		if err := os.MkdirAll(facturasDir, 0755); err != nil {
			log.Printf("Error al crear directorio para facturas: %v", err)
			http.Error(w, "Error al generar factura", http.StatusInternalServerError)
			return
		}
	}

	// Obtener la plantilla activa para este usuario desde la base de datos
	dbConn := db.GetDB()
	var plantillaRuta string
	var plantillaNombre string

	err = dbConn.QueryRow(
		"SELECT ruta_archivo, nombre FROM plantillas_factura WHERE id_usuario = ? AND activa = true",
		facturaRequest.IDUsuario,
	).Scan(&plantillaRuta, &plantillaNombre)

	// Si no hay plantilla activa, usar una plantilla por defecto
	if err == sql.ErrNoRows {
		plantillaRuta = "./templates/facturas/plantilla_default.docx"
		plantillaNombre = "Plantilla por defecto"
		log.Printf("No se encontró plantilla activa para el usuario %d, usando plantilla por defecto",
			facturaRequest.IDUsuario)
	} else if err != nil {
		log.Printf("Error al obtener plantilla activa: %v", err)
		http.Error(w, "Error al generar la factura", http.StatusInternalServerError)
		return
	}

	// Verificar que la plantilla existe
	if _, err := os.Stat(plantillaRuta); os.IsNotExist(err) {
		log.Printf("Plantilla no encontrada: %s", plantillaRuta)
		http.Error(w, "Plantilla no encontrada", http.StatusInternalServerError)
		return
	}

	// Generar nombre único para la factura
	timestamp := time.Now().Format("20060102150405")
	nombreArchivo := fmt.Sprintf("factura_%d_%s.pdf", facturaRequest.IDUsuario, timestamp)
	rutaDestino := filepath.Join(facturasDir, nombreArchivo)

	// Aquí iría el código que genera la factura a partir de la plantilla
	// ...

	// Simulamos la generación para este ejemplo
	log.Printf("Generando factura para %s con RFC %s usando plantilla %s",
		facturaRequest.Empresa, facturaRequest.RFC, plantillaNombre)

	// Registrar la factura en la base de datos
	var total float64 = 0
	var subtotal float64 = 0

	for _, concepto := range facturaRequest.Conceptos {
		subtotal += float64(concepto.Cantidad) * concepto.PrecioUni
	}

	// Calcular impuestos (16% IVA por ejemplo)
	impuestos := subtotal * 0.16
	total = subtotal + impuestos

	// Insertar en la base de datos
	_, err = dbConn.Exec(
		"INSERT INTO facturas (id_usuario, empresa, rfc, subtotal, impuestos, total, fecha_emision, ruta_pdf) VALUES (?, ?, ?, ?, ?, ?, NOW(), ?)",
		facturaRequest.IDUsuario, facturaRequest.Empresa, facturaRequest.RFC, subtotal, impuestos, total, rutaDestino,
	)

	if err != nil {
		log.Printf("Error al registrar factura en base de datos: %v", err)
		// Continuamos a pesar del error para que el usuario reciba su factura
	}

	// Al final, devolver la URL donde se puede descargar la factura
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Factura generada correctamente",
		"factura": map[string]interface{}{
			"empresa":         facturaRequest.Empresa,
			"rfc":             facturaRequest.RFC,
			"subtotal":        subtotal,
			"impuestos":       impuestos,
			"total":           total,
			"fecha":           time.Now().Format("2006-01-02 15:04:05"),
			"url":             "/facturas/descargar/" + nombreArchivo,
			"plantilla_usada": plantillaNombre,
		},
	})
}
