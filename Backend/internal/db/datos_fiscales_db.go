package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func GuardarDatosFiscales(
	rfc, razonSocial, direccionFiscal, codigoPostal string,
	archivoCSDKey, archivoCSDCer []byte,
	nombreArchivoKey, nombreArchivoCer,
	claveCSD, regimenFiscal, serieDf string,
	usuarioID int) (int, error) {

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
				serie_df = ?, fecha_actualizacion = ?
			WHERE id_usuario = ?
		`
		_, err = tx.Exec(query, rfc, razonSocial, direccionFiscal, codigoPostal,
			regimenFiscal, claveCSD, serieDf, now, usuarioID)
		if err != nil {
			return 0, fmt.Errorf("error al actualizar datos básicos: %w", err)
		}

		// Actualizar archivos solo si se proporcionan
		if len(archivoCSDKey) > 0 {
			_, err = tx.Exec(
				"UPDATE datos_fiscales SET archivo_key = ?, nombre_archivo_key = ? WHERE id_usuario = ?",
				archivoCSDKey, nombreArchivoKey, usuarioID,
			)
			if err != nil {
				return 0, fmt.Errorf("error al actualizar archivo KEY: %w", err)
			}
		}
		if len(archivoCSDCer) > 0 {
			_, err = tx.Exec(
				"UPDATE datos_fiscales SET archivo_cer = ?, nombre_archivo_cer = ? WHERE id_usuario = ?",
				archivoCSDCer, nombreArchivoCer, usuarioID,
			)
			if err != nil {
				return 0, fmt.Errorf("error al actualizar archivo CER: %w", err)
			}
		}
	} else {
		// Insertar nuevo registro
		query := `
			INSERT INTO datos_fiscales (
				id_usuario, rfc, razon_social, direccion_fiscal, 
				codigo_postal, regimen_fiscal, clave_csd, serie_df,
				archivo_key, nombre_archivo_key, archivo_cer, nombre_archivo_cer,
				fecha_creacion, fecha_actualizacion
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		result, err := tx.Exec(query,
			usuarioID, rfc, razonSocial, direccionFiscal,
			codigoPostal, regimenFiscal, claveCSD, serieDf,
			archivoCSDKey, nombreArchivoKey, archivoCSDCer, nombreArchivoCer,
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
	SELECT archivo_cer, archivo_key, id, rfc, razon_social, direccion_fiscal, direccion, colonia, 
		   codigo_postal, ciudad, estado, regimen_fiscal, clave_csd, 
		   serie_df, nombre_archivo_key, nombre_archivo_cer
	FROM datos_fiscales
	WHERE id_usuario = ?
`

	var archivoCer, archivoKey []byte
	var id int
	var rfc, razonSocial string
	var direccionFiscal, direccion, colonia, codigoPostal, ciudad, estado, regimenFiscal sql.NullString
	var claveCSD, serieDf sql.NullString
	var nombreArchivoKey, nombreArchivoCer sql.NullString

	err = db.QueryRow(query, userID).Scan(
		&archivoCer, &archivoKey, &id, &rfc, &razonSocial, &direccionFiscal, &direccion, &colonia,
		&codigoPostal, &ciudad, &estado, &regimenFiscal, &claveCSD, &serieDf,
		&nombreArchivoKey, &nombreArchivoCer,
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
		"archivo_cer":  archivoCer,
		"archivo_key":  archivoKey,
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

	if nombreArchivoKey.Valid {
		datos["nombre_archivo_key"] = nombreArchivoKey.String
	}
	if nombreArchivoCer.Valid {
		datos["nombre_archivo_cer"] = nombreArchivoCer.String
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
