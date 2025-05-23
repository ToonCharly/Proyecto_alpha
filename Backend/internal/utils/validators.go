package utils

import "strings"

// IfEmpty devuelve un valor predeterminado si el valor proporcionado está vacío.
func IfEmpty(value, defaultValue string) string {
    if strings.TrimSpace(value) == "" {
        return defaultValue
    }
    return value
}