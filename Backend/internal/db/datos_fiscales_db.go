package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"Facts/internal/utils"
)

func GuardarDatosFiscales(
	rfc, razonSocial, direccionFiscal, codigoPostal string,
	archivoCSDKey, archivoCSDCer []byte,
	nombreArchivoKey, nombreArchivoCer,
	claveCSD, regimenFiscal, serieDf string,
	usuarioID int,
	keyPath string, // NUEVO: ruta al archivo .key
) (int, error) {

	if usuarioID <= 0 {
		return 0, fmt.Errorf("ID de usuario inválido: %d", usuarioID)
	}

	log.Printf("GuardarDatosFiscales: iniciando para usuario ID %d", usuarioID)

	db, err := ConnectUserDB()
	if err != nil {
		return 0, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM datos_fiscales WHERE id_usuario = ?)", usuarioID).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("error al verificar datos existentes: %w", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	var idDatosFiscales int64

	// Guardar archivos en disco y obtener rutas
	var rutaKey, rutaCer string
	baseDir := fmt.Sprintf("./certificados/%d", usuarioID)
	utils.CreateDirectory(baseDir)
	if len(archivoCSDKey) > 0 {
		rutaKey = fmt.Sprintf("%s/%s", baseDir, nombreArchivoKey)
		err := utils.SaveFile(rutaKey, archivoCSDKey)
		if err != nil {
			return 0, fmt.Errorf("error al guardar archivo KEY en disco: %w", err)
		}
	}
	if len(archivoCSDCer) > 0 {
		rutaCer = fmt.Sprintf("%s/%s", baseDir, nombreArchivoCer)
		log.Printf("[GUARDAR_CER] Ruta final del archivo .cer: %s", rutaCer)
		err := utils.SaveFile(rutaCer, archivoCSDCer)
		if err != nil {
			return 0, fmt.Errorf("error al guardar archivo CER en disco: %w", err)
		}
	}

	if exists {
		// Obtener el id actual de datos_fiscales
		err = db.QueryRow("SELECT id FROM datos_fiscales WHERE id_usuario = ?", usuarioID).Scan(&idDatosFiscales)
		if err != nil {
			return 0, fmt.Errorf("error al obtener id de datos fiscales: %w", err)
		}

		// Actualizar registro existente
		query := `
			UPDATE datos_fiscales 
			SET rfc = ?, razon_social = ?, direccion_fiscal = ?, 
				codigo_postal = ?, regimen_fiscal = ?, clave_csd = ?,
				serie_df = ?, fecha_actualizacion = ?, ruta_archivo_key = ?, ruta_archivo_cer = ?
			WHERE id_usuario = ?
		`
		_, err = tx.Exec(query, rfc, razonSocial, direccionFiscal, codigoPostal,
			regimenFiscal, claveCSD, serieDf, now, rutaKey, rutaCer, usuarioID)
		if err != nil {
			return 0, fmt.Errorf("error al actualizar datos básicos: %w", err)
		}
	} else {
		// Insertar nuevo registro
		query := `
			INSERT INTO datos_fiscales (
				id_usuario, rfc, razon_social, direccion_fiscal, 
				codigo_postal, regimen_fiscal, clave_csd, serie_df,
				ruta_archivo_key, ruta_archivo_cer,
				fecha_creacion, fecha_actualizacion
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		result, err := tx.Exec(query,
			usuarioID, rfc, razonSocial, direccionFiscal,
			codigoPostal, regimenFiscal, claveCSD, serieDf,
			rutaKey, rutaCer,
			now, now)
		if err != nil {
			return 0, fmt.Errorf("error al insertar datos fiscales: %w", err)
		}
		idDatosFiscales, err = result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("error al obtener id insertado: %w", err)
		}
	}

	// Actualiza el campo id_datos_fiscales en la tabla usuarios
	_, err = tx.Exec("UPDATE usuarios SET id_datos_fiscales = ? WHERE id = ?", idDatosFiscales, usuarioID)
	if err != nil {
		return 0, fmt.Errorf("error al actualizar id_datos_fiscales en usuarios: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	log.Printf("Datos fiscales guardados exitosamente para el usuario %d (empresa %d)", usuarioID, idDatosFiscales)
	return int(idDatosFiscales), nil
}

// ObtenerDatosFiscales obtiene los datos fiscales de un usuario
func ObtenerDatosFiscales(userID int) (map[string]interface{}, error) {
	db, err := ConnectUserDB()
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer db.Close()

	query := `
	SELECT archivo_cer, archivo_key, archivo_cer_pem, archivo_key_pem, id, rfc, razon_social, direccion_fiscal, direccion, colonia, 
		   codigo_postal, ciudad, estado, regimen_fiscal, clave_csd, 
		   serie_df, ruta_archivo_key, ruta_archivo_cer
	FROM datos_fiscales
	WHERE id_usuario = ?
`

	var archivoCer, archivoKey []byte
	var archivoCerPEM, archivoKeyPEM sql.NullString
	var id int
	var rfc, razonSocial string
	var direccionFiscal, direccion, colonia, codigoPostal, ciudad, estado, regimenFiscal sql.NullString
	var claveCSD, serieDf sql.NullString
	var rutaArchivoKey, rutaArchivoCer sql.NullString

	err = db.QueryRow(query, userID).Scan(
		&archivoCer, &archivoKey, &archivoCerPEM, &archivoKeyPEM, &id, &rfc, &razonSocial, &direccionFiscal, &direccion, &colonia,
		&codigoPostal, &ciudad, &estado, &regimenFiscal, &claveCSD, &serieDf,
		&rutaArchivoKey, &rutaArchivoCer,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no se encontraron datos fiscales para el usuario")
		}
		return nil, fmt.Errorf("error al obtener datos fiscales: %w", err)
	}

	datos := map[string]interface{}{
		"id":              id,
		"rfc":             rfc,
		"razon_social":    razonSocial,
		"archivo_cer":     archivoCer,
		"archivo_key":     archivoKey,
		"archivo_cer_pem": "",
		"archivo_key_pem": "",
	}
	if archivoCerPEM.Valid {
		datos["archivo_cer_pem"] = archivoCerPEM.String
	}
	if archivoKeyPEM.Valid {
		datos["archivo_key_pem"] = archivoKeyPEM.String
	}

	if direccionFiscal.Valid {
		datos["direccion_fiscal"] = direccionFiscal.String
	} else {
		datos["direccion_fiscal"] = ""
	}
	if direccion.Valid {
		datos["direccion"] = direccion.String
	} else {
		datos["direccion"] = ""
	}
	if colonia.Valid {
		datos["colonia"] = colonia.String
	} else {
		datos["colonia"] = ""
	}
	if codigoPostal.Valid {
		datos["codigo_postal"] = codigoPostal.String
	} else {
		datos["codigo_postal"] = ""
	}
	if ciudad.Valid {
		datos["ciudad"] = ciudad.String
	} else {
		datos["ciudad"] = ""
	}
	if estado.Valid {
		datos["estado"] = estado.String
	} else {
		datos["estado"] = ""
	}
	if regimenFiscal.Valid {
		datos["regimen_fiscal"] = regimenFiscal.String
	} else {
		datos["regimen_fiscal"] = ""
	}
	if claveCSD.Valid {
		datos["clave_csd"] = claveCSD.String
	}
	if serieDf.Valid {
		datos["serie_df"] = serieDf.String
	} else {
		datos["serie_df"] = ""
	}

	// Para saber si hay archivos, revisa la longitud de los []byte
	datos["tiene_key"] = len(archivoKey) > 0
	datos["tiene_cer"] = len(archivoCer) > 0

	// Ya no se asignan los nombres de archivo, solo las rutas locales
	// Agregar las rutas locales de los archivos .key y .cer
	if rutaArchivoKey.Valid {
		datos["key_path"] = rutaArchivoKey.String
	} else {
		datos["key_path"] = ""
	}
	if rutaArchivoCer.Valid {
		datos["cer_path"] = rutaArchivoCer.String
	} else {
		datos["cer_path"] = ""
	}

	return datos, nil
}

// ObtenerCertificadosCSD obtiene los archivos binarios de certificados directamente de la base de datos
func ObtenerCertificadosCSD(userID int) ([]byte, []byte, error) {
	db, err := ConnectUserDB()
	if err != nil {
		return nil, nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer db.Close()

	var keyData, cerData []byte
	err = db.QueryRow(
		"SELECT archivo_key, archivo_cer FROM datos_fiscales WHERE id_usuario = ?",
		userID,
	).Scan(&keyData, &cerData)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("no se encontraron certificados para el usuario")
		}
		return nil, nil, fmt.Errorf("error al obtener certificados: %w", err)
	}

	return keyData, cerData, nil
}
