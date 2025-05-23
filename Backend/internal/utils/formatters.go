package utils

import (
    "fmt" 
    "time"
)

// FormatearFecha convierte una fecha en formato ISO a un formato más legible.
func FormatearFecha(fechaISO string) string {
    t, err := time.Parse("2006-01-02T15:04:05", fechaISO)
    if err != nil {
        return fechaISO // Si hay error, retornar la fecha original
    }
    return t.Format("02/01/2006 15:04:05")
}

// FormatearMoneda convierte un número flotante en un formato de moneda.
func FormatearMoneda(valor float64) string {
    return fmt.Sprintf("$%.2f", valor)
}