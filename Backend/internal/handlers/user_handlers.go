package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"Facts/internal/db"
	"Facts/internal/utils"

	"github.com/golang-jwt/jwt/v4"
)

// Función para verificar token JWT
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return utils.JWTKey, nil // Usar la clave centralizada
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token inválido")
}

// GetAllUsersHandler obtiene todos los usuarios registrados en el sistema
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Manejar preflight OPTIONS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Verificar método
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar autorización (token)
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		// Si no hay token, simplemente continuar sin verificación para desarrollo
		// En producción, deberías descomentar la siguiente línea:
		// http.Error(w, "No autorizado", http.StatusUnauthorized)
		// return
		log.Println("Advertencia: Solicitud sin token de autorización")
	} else {
		// Extraer token del formato Bearer
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[7:]
		}

		// Verificar token (opcional durante desarrollo)
		claims, err := VerifyToken(tokenString)
		if err != nil {
			log.Printf("Token inválido: %v", err)
			// En desarrollo, continuar sin token válido
			// En producción, descomentar:
			// http.Error(w, "Token inválido", http.StatusUnauthorized)
			// return
		} else {
			// Verificar si es admin
			role, ok := claims["role"].(string)
			if !ok || role != "admin" {
				log.Printf("Usuario no es administrador: %v", claims["email"])
				// En desarrollo, continuar aunque no sea admin
				// En producción, descomentar:
				// http.Error(w, "Acceso denegado. Se requiere rol de administrador", http.StatusForbidden)
				// return
			}
		}
	}

	// Obtener todos los usuarios
	users, err := db.GetAllUsers()
	if err != nil {
		log.Printf("Error al obtener usuarios: %v", err)
		utils.RespondWithError(w, "Error al obtener la lista de usuarios")
		return
	}

	// Responder con la lista de usuarios
	utils.RespondWithJSON(w, http.StatusOK, users)
}

// UpdateUserRoleHandler actualiza el rol de un usuario (admin o no)
func UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Manejar preflight OPTIONS
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Verificar método
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar autorización (token)
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Extraer token del formato Bearer
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = tokenString[7:]
	}

	// Aquí podrías verificar si el usuario tiene permisos de administrador
	// por ahora, simplemente continuamos

	// Obtener ID del usuario de la URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}

	userIDStr := pathParts[3]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo de la solicitud
	var requestBody struct {
		IsAdmin bool `json:"isAdmin"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		log.Printf("Error decodificando cuerpo: %v", err)
		http.Error(w, "Error al leer datos de la solicitud", http.StatusBadRequest)
		return
	}

	// Actualizar el rol del usuario
	if err := db.UpdateUserRole(userID, requestBody.IsAdmin); err != nil {
		log.Printf("Error al actualizar rol: %v", err)
		utils.RespondWithError(w, fmt.Sprintf("Error al actualizar rol: %v", err))
		return
	}

	// Responder con éxito
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"mensaje": fmt.Sprintf("Rol de usuario actualizado exitosamente a %s",
			map[bool]string{true: "administrador", false: "usuario"}[requestBody.IsAdmin]),
	})
}
