package models

import (
	"database/sql"
	"fmt"
	"log"

	"Facts/internal/db"
)

// FolioControl representa el control de folios para facturas
type FolioControl struct {
	ID                 int    `json:"id"`
	Serie              string `json:"serie"`
	UltimoFolio        int64  `json:"ultimo_folio"`
	EmpresaID          int    `json:"empresa_id"`
	FechaCreacion      string `json:"fecha_creacion"`
	FechaActualizacion string `json:"fecha_actualizacion"`
	Activo             bool   `json:"activo"`
}

// GenerarFolio genera el siguiente número de folio para una serie y empresa específica
func GenerarFolio(db *sql.DB, serie string, empresaID int) (string, error) {
	// Usar la función almacenada para obtener el siguiente folio de manera atómica
	var nuevoFolio int64

	query := "SELECT obtener_siguiente_folio(?, ?)"
	err := db.QueryRow(query, serie, empresaID).Scan(&nuevoFolio)
	if err != nil {
		log.Printf("Error al generar folio: %v", err)
		return "", fmt.Errorf("error al generar folio: %v", err)
	}

	// Si el folio es 0, significa que no existe la serie para esa empresa
	if nuevoFolio == 0 {
		// Crear nuevo registro para la serie y empresa
		_, err = db.Exec(
			"INSERT INTO folio_control (serie, ultimo_folio, empresa_id) VALUES (?, 1, ?)",
			serie, empresaID,
		)
		if err != nil {
			return "", fmt.Errorf("error al crear serie de folios: %v", err)
		}
		nuevoFolio = 1
	}

	// Formatear el folio con la serie
	folioCompleto := fmt.Sprintf("%s%d", serie, nuevoFolio)

	log.Printf("Folio generado: %s para empresa %d", folioCompleto, empresaID)
	return folioCompleto, nil
}

// GenerarFolioConPadding genera un folio con padding de ceros
func GenerarFolioConPadding(serie string, empresaID int, padding int) (string, error) {
	// Conectar a la base de datos Usuario donde está la tabla folio_control
	userDB, err := db.ConnectUserDB()
	if err != nil {
		return "", fmt.Errorf("error al conectar a la base de datos: %v", err)
	}
	defer userDB.Close()

	// Usar la función almacenada para obtener el siguiente folio
	var nuevoFolio int64

	query := "SELECT obtener_siguiente_folio(?, ?)"
	err = userDB.QueryRow(query, serie, empresaID).Scan(&nuevoFolio)
	if err != nil {
		log.Printf("Error al generar folio: %v", err)
		return "", fmt.Errorf("error al generar folio: %v", err)
	}

	// Si el folio es 0, crear nuevo registro
	if nuevoFolio == 0 {
		_, err = userDB.Exec(
			"INSERT INTO folio_control (serie, ultimo_folio, empresa_id) VALUES (?, 1, ?)",
			serie, empresaID,
		)
		if err != nil {
			return "", fmt.Errorf("error al crear serie de folios: %v", err)
		}
		nuevoFolio = 1
	}

	// Formatear con padding de ceros
	formato := fmt.Sprintf("%%s%%0%dd", padding)
	folioCompleto := fmt.Sprintf(formato, serie, nuevoFolio)

	log.Printf("Folio generado con padding: %s para empresa %d", folioCompleto, empresaID)
	return folioCompleto, nil
}

// ObtenerUltimoFolio obtiene el último folio usado para una serie y empresa
func ObtenerUltimoFolio(db *sql.DB, serie string, empresaID int) (int64, error) {
	var ultimoFolio int64

	query := `SELECT ultimo_folio FROM folio_control 
			  WHERE serie = ? AND empresa_id = ? AND activo = TRUE`

	err := db.QueryRow(query, serie, empresaID).Scan(&ultimoFolio)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No existe la serie, empezará desde 1
		}
		return 0, fmt.Errorf("error al obtener último folio: %v", err)
	}

	return ultimoFolio, nil
}

// ValidarFolioUnico verifica que un folio no esté duplicado
func ValidarFolioUnico(folio string, empresaID int) (bool, error) {
	userDB, err := db.ConnectUserDB()
	if err != nil {
		return false, fmt.Errorf("error al conectar a la base de datos: %v", err)
	}
	defer userDB.Close()

	var count int

	query := `SELECT COUNT(*) FROM historial_facturas 
			  WHERE numero_folio = ? AND id_empresa = ?`

	err = userDB.QueryRow(query, folio, empresaID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error al validar folio único: %v", err)
	}

	return count == 0, nil
}

// CrearSerieEmpresa crea una nueva serie de folios para una empresa
func CrearSerieEmpresa(db *sql.DB, serie string, empresaID int) error {
	query := `INSERT INTO folio_control (serie, ultimo_folio, empresa_id) 
			  VALUES (?, 0, ?) 
			  ON DUPLICATE KEY UPDATE activo = TRUE`

	_, err := db.Exec(query, serie, empresaID)
	if err != nil {
		return fmt.Errorf("error al crear serie de empresa: %v", err)
	}

	log.Printf("Serie %s creada para empresa %d", serie, empresaID)
	return nil
}
