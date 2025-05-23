package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectUserDB establece conexión con la base de datos de usuarios (en Docker)
func ConnectUserDB() (*sql.DB, error) {
	username := "root"
	password := "facts_senior"
	hostname := "localhost"
	port := "3306"
	dbname := "Usuario"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos de usuarios: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error conectando a la base de datos de usuarios: %w", err)
	}

	fmt.Println("Conexión exitosa a la base de datos de usuarios")
	return db, nil
}

// RegisterUser inserta un nuevo usuario en la base de datos de usuarios.
func RegisterUser(username, email, phone, password string) (map[string]interface{}, error) {
	db, err := ConnectUserDB()
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos de usuarios: %w", err)
	}
	defer db.Close()

	// Verificar duplicados (como ya lo tienes)
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuarios WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error al verificar duplicados: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("El correo ya está registrado. Por favor, use otro.")
	}

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM usuarios WHERE phone = ?)", phone).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error al verificar duplicados: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("El número de teléfono ya está registrado. Por favor, use otro.")
	}

	// Verificar si es el primer usuario
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM usuarios").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("error al verificar usuarios existentes: %w", err)
	}

	// Determinar el rol
	role := "user"
	if count == 0 {
		role = "admin"
	}

	// Insertar usuario con rol
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queryInsert := `INSERT INTO usuarios (username, email, phone, password, role) VALUES (?, ?, ?, ?, ?)`
	result, err := tx.Exec(queryInsert, username, email, phone, password, role)
	if err != nil {
		return nil, fmt.Errorf("error al registrar el usuario: %w", err)
	}

	// Obtener el ID del usuario recién insertado
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ID del usuario: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	// Preparar respuesta
	userData := map[string]interface{}{
		"id":       userID,
		"username": username,
		"email":    email,
		"phone":    phone,
		"role":     role,
	}

	fmt.Println("Usuario registrado exitosamente con rol:", role)
	return userData, nil
}

// LoginUser verifica las credenciales del usuario y devuelve sus datos si son correctas
func LoginUser(email, username, password string) (map[string]interface{}, error) {
	db, err := ConnectUserDB()
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos de usuarios: %w", err)
	}
	defer db.Close()

	var query string
	var args []interface{}

	if email != "" {
		query = `SELECT id, username, email, phone, role FROM usuarios WHERE email = ? AND password = ?`
		args = []interface{}{email, password}
	} else if username != "" {
		query = `SELECT id, username, email, phone, role FROM usuarios WHERE username = ? AND password = ?`
		args = []interface{}{username, password}
	} else {
		return nil, fmt.Errorf("debe proporcionar email o nombre de usuario")
	}

	row := db.QueryRow(query, args...)

	var user struct {
		ID       int
		Username string
		Email    string
		Phone    string
		Role     string
	}

	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("credenciales inválidas")
		}
		return nil, fmt.Errorf("error al consultar usuario: %w", err)
	}

	userData := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"phone":    user.Phone,
		"role":     user.Role,
	}

	// El resto de tu código para obtener detalles adicionales se mantiene igual
	var direccion, codigoPostal, ciudad, estado sql.NullString
	queryDetalles := `SELECT direccion, codigo_postal, ciudad, estado 
                      FROM usuario_detalles 
                      WHERE usuario_id = ?`

	err = db.QueryRow(queryDetalles, user.ID).Scan(&direccion, &codigoPostal, &ciudad, &estado)

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

	return userData, nil
}

// UpdateUser actualiza la información de un usuario solo en la tabla usuarios
func UpdateUser(db *sql.DB, email, username, phone, direccion, codigoPostal, ciudad, estado string) error {
	conn, err := ConnectUserDB()
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos de usuarios: %w", err)
	}
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Verificar que el usuario existe
	var userID int
	err = tx.QueryRow("SELECT id FROM usuarios WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no existe un usuario con el email proporcionado")
		}
		return fmt.Errorf("error al buscar usuario: %w", err)
	}

	// Actualizar toda la información en la tabla usuarios
	_, err = tx.Exec(
		"UPDATE usuarios SET username=?, phone=? WHERE email=?",
		username, phone, email,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar datos del usuario: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	fmt.Println("Información de usuario actualizada exitosamente")
	return nil
}
