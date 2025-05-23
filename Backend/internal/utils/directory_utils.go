package utils

import (
    "log"
    "os"
)

// CreateDirectory crea un directorio si no existe.
func CreateDirectory(dir string) error {
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        if err := os.MkdirAll(dir, 0755); err != nil {
            log.Printf("Error al crear el directorio %s: %v", dir, err)
            return err
        }
        log.Printf("Directorio creado: %s", dir)
    }
    return nil
}