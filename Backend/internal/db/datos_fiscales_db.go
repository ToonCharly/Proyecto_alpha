package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// GuardarDatosFiscales guarda o actualiza los datos fiscales incluyendo archivos binarios
func GuardarDatosFiscales(
	rfc, razonSocial, direccionFiscal, codigoPostal string,
	archivoCSDKey, archivoCSDCer []byte,
	claveCSD, regimenFiscal string,
	usuarioID int) error {

	// Verificar que el ID de usuario sea válido
	if usuarioID <= 0 {
		return fmt.Errorf("ID de usuario inválido: %d", usuarioID)
	}

	// Más log para depuración
	log.Printf("GuardarDatosFiscales: iniciando para usuario ID %d", usuarioID)
	log.Printf("Datos: rfc=%s, razon_social=%s, id_usuario=%d", rfc, razonSocial, usuarioID)

	// Conectar a la base de datos
	db, err := ConnectUserDB()
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer db.Close()

	// Crear directorio para certificados si no existe
	if err := os.MkdirAll("uploads/certificados", 0755); err != nil {
		return fmt.Errorf("error al crear directorio para certificados: %w", err)
	}

	// Guardar archivos en disco si fueron proporcionados
	var rutaCSDKey, rutaCSDCer string
	if len(archivoCSDKey) > 0 {
		rutaCSDKey = fmt.Sprintf("uploads/certificados/%d_key.key", usuarioID)
		if err := ioutil.WriteFile(rutaCSDKey, archivoCSDKey, 0644); err != nil {
			return fmt.Errorf("error al guardar archivo KEY: %w", err)
		}
		log.Printf("Archivo KEY guardado en: %s", rutaCSDKey)
	}

	if len(archivoCSDCer) > 0 {
		rutaCSDCer = fmt.Sprintf("uploads/certificados/%d_cer.cer", usuarioID)
		if err := ioutil.WriteFile(rutaCSDCer, archivoCSDCer, 0644); err != nil {
			return fmt.Errorf("error al guardar archivo CER: %w", err)
		}
		log.Printf("Archivo CER guardado en: %s", rutaCSDCer)
	}

	// ENFOQUE DIRECTO: Insertar registro completo o actualizarlo

	// Iniciar transacción
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}

	// Si algo falla, hacer rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Verificar si ya existe un registro para este usuario
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM datos_fiscales WHERE id_usuario = ?)", usuarioID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error al verificar datos existentes: %w", err)
	}

	// Preparar valores para fechas
	now := time.Now().Format("2006-01-02 15:04:05")

	if exists {
		// Actualizar registro existente
		log.Printf("Actualizando registro existente para usuario ID %d", usuarioID)

		query := `
            UPDATE datos_fiscales 
            SET rfc = ?, razon_social = ?, direccion_fiscal = ?, 
                codigo_postal = ?, regimen_fiscal = ?, clave_csd = ?,
                fecha_actualizacion = ?
            WHERE id_usuario = ?
        `
		_, err = tx.Exec(query, rfc, razonSocial, direccionFiscal, codigoPostal,
			regimenFiscal, claveCSD, now, usuarioID)

		if err != nil {
			return fmt.Errorf("error al actualizar datos básicos: %w", err)
		}

		// Actualizar rutas de archivos solo si se proporcionaron nuevos archivos
		if len(archivoCSDKey) > 0 {
			_, err = tx.Exec("UPDATE datos_fiscales SET ruta_csd_key = ? WHERE id_usuario = ?",
				rutaCSDKey, usuarioID)
			if err != nil {
				return fmt.Errorf("error al actualizar ruta KEY: %w", err)
			}
		}

		if len(archivoCSDCer) > 0 {
			_, err = tx.Exec("UPDATE datos_fiscales SET ruta_csd_cer = ? WHERE id_usuario = ?",
				rutaCSDCer, usuarioID)
			if err != nil {
				return fmt.Errorf("error al actualizar ruta CER: %w", err)
			}
		}
	} else {
		// Insertar nuevo registro
		log.Printf("Insertando nuevo registro para usuario ID %d", usuarioID)

		// IMPORTANTE: USAR INSERCIÓN EXPLÍCITA CON id_usuario PRIMERO
		query := `
            INSERT INTO datos_fiscales (
                id_usuario, rfc, razon_social, direccion_fiscal, 
                codigo_postal, regimen_fiscal, clave_csd, 
                ruta_csd_key, ruta_csd_cer, 
                fecha_creacion, fecha_actualizacion
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `
		// Log de la consulta y valores para depuración
		log.Printf("SQL: %s", query)
		log.Printf("Valores: id_usuario=%d, rfc=%s, razon_social=%s", usuarioID, rfc, razonSocial)

		// CRUCIAL: Usar la variable now en lugar de NOW()
		_, err = tx.Exec(query,
			usuarioID, rfc, razonSocial, direccionFiscal,
			codigoPostal, regimenFiscal, claveCSD,
			rutaCSDKey, rutaCSDCer, now, now)

		if err != nil {
			return fmt.Errorf("error al insertar datos fiscales: %w", err)
		}
	}

	// Confirmar transacción
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	log.Printf("Datos fiscales guardados exitosamente para el usuario %d", usuarioID)
	return nil
}

// ObtenerDatosFiscales obtiene los datos fiscales de un usuario
func ObtenerDatosFiscales(userID int) (map[string]interface{}, error) {
	db, err := ConnectUserDB()
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer db.Close()

	query := `
        SELECT id, rfc, razon_social, direccion_fiscal, codigo_postal, 
               regimen_fiscal, clave_csd, ruta_csd_key, ruta_csd_cer,
               fecha_creacion, fecha_actualizacion
        FROM datos_fiscales
        WHERE id_usuario = ?
    `

	var id int
	var rfc, razonSocial string
	var direccionFiscal, codigoPostal, regimenFiscal sql.NullString
	var claveCSD, rutaCSDKey, rutaCSDCer sql.NullString
	var fechaCreacion, fechaActualizacion sql.NullTime

	err = db.QueryRow(query, userID).Scan(
		&id, &rfc, &razonSocial, &direccionFiscal, &codigoPostal,
		&regimenFiscal, &claveCSD, &rutaCSDKey, &rutaCSDCer,
		&fechaCreacion, &fechaActualizacion,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no se encontraron datos fiscales para el usuario")
		}
		return nil, fmt.Errorf("error al obtener datos fiscales: %w", err)
	}

	datos := map[string]interface{}{
		"id":           id,
		"rfc":          rfc,
		"razon_social": razonSocial,
	}

	// Añadir campos que pueden ser NULL solo si tienen valor
	if direccionFiscal.Valid {
		datos["direccion_fiscal"] = direccionFiscal.String
	} else {
		datos["direccion_fiscal"] = ""
	}

	if codigoPostal.Valid {
		datos["codigo_postal"] = codigoPostal.String
	} else {
		datos["codigo_postal"] = ""
	}

	if regimenFiscal.Valid {
		datos["regimen_fiscal"] = regimenFiscal.String
	} else {
		datos["regimen_fiscal"] = ""
	}

	if claveCSD.Valid {
		datos["clave_csd"] = claveCSD.String
	}

	if rutaCSDKey.Valid && rutaCSDKey.String != "" {
		datos["tiene_key"] = true
	} else {
		datos["tiene_key"] = false
	}

	if rutaCSDCer.Valid && rutaCSDCer.String != "" {
		datos["tiene_cer"] = true
	} else {
		datos["tiene_cer"] = false
	}

	if fechaCreacion.Valid {
		datos["fecha_creacion"] = fechaCreacion.Time
	}

	if fechaActualizacion.Valid {
		datos["fecha_actualizacion"] = fechaActualizacion.Time
	}

	return datos, nil
}

// ObtenerCertificadosCSD obtiene los archivos binarios de certificados
func ObtenerCertificadosCSD(userID int) ([]byte, []byte, error) {
	db, err := ConnectUserDB()
	if err != nil {
		return nil, nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer db.Close()

	var rutaKey, rutaCer sql.NullString
	err = db.QueryRow(
		"SELECT ruta_csd_key, ruta_csd_cer FROM datos_fiscales WHERE id_usuario = ?",
		userID,
	).Scan(&rutaKey, &rutaCer)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("no se encontraron certificados para el usuario")
		}
		return nil, nil, fmt.Errorf("error al obtener rutas de certificados: %w", err)
	}

	var keyData, cerData []byte

	if rutaKey.Valid && rutaKey.String != "" {
		keyData, err = ioutil.ReadFile(rutaKey.String)
		if err != nil && !os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("error al leer archivo KEY: %w", err)
		}
	}

	if rutaCer.Valid && rutaCer.String != "" {
		cerData, err = ioutil.ReadFile(rutaCer.String)
		if err != nil && !os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("error al leer archivo CER: %w", err)
		}
	}

	return keyData, cerData, nil
}
