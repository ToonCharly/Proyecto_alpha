package db

import (
    "database/sql"
    "fmt"

    _ "github.com/go-sql-driver/mysql" 
)

// Connect establece la conexión a la base de datos y la devuelve.
func Connect() (*sql.DB, error) {
    dsn := "alpha_junior:GHtLop23_P54@tcp(199.89.55.249:3306)/optimus"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("no se pudo conectar a la base de datos: %w", err)
    }

    fmt.Println("Conexión exitosa a la base de datos")
    return db, nil
}

func ConnectToOptimus() (*sql.DB, error) {
    return Connect() // Ya conecta a optimus por defecto
}