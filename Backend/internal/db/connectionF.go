package db

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitDB inicializa la conexión a la base de datos
func InitDB() {
    var err error

    // Credenciales de la base de datos
    username := "root"
    password := "facts_senior"
    hostname := "localhost"
    port := "3306"
    dbname := "Usuario"

    // Construcción del DSN (Data Source Name)
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbname)

    // Conexión a la base de datos
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Error al conectar a la base de datos: %v", err)
    }

    // Verifica la conexión
    err = db.Ping()
    if err != nil {
        log.Fatalf("No se pudo conectar a la base de datos: %v", err)
    }

    fmt.Println("Conexión exitosa a la base de datos")
}

// GetDB devuelve la conexión a la base de datos
func GetDB() *sql.DB {
    if db == nil {
        log.Fatal("La conexión a la base de datos no está inicializada")
    }
    return db
}