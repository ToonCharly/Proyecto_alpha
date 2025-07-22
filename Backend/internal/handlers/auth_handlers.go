package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"Facts/internal/db"
)

// LoginHandler maneja la autenticación de usuarios con enfoque simplificado
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var loginData struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginData); err != nil {
		log.Printf("Error decodificando JSON de login: %v", err)
		http.Error(w, "Error al leer los datos de login", http.StatusBadRequest)
		return
	}

	if (loginData.Email == "" && loginData.Username == "") || loginData.Password == "" {
		http.Error(w, "Se requiere email/username y contraseña", http.StatusBadRequest)
		return
	}

	userData, err := db.LoginUser(loginData.Email, loginData.Username, loginData.Password)
	if err != nil {
		log.Printf("Error en login: %v", err)
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	// ENFOQUE EXTREMADAMENTE SIMPLIFICADO
	// Generar un token simple pero reconocible
	tokenString := fmt.Sprintf("TOKEN_%d_%s", time.Now().Unix(), userData["email"])

	// Crear una respuesta con JSON crudo predefinido
	hardcodedResponse := fmt.Sprintf(`{
		"id": %v,
		"username": "%s",
		"email": "%s",
		"phone": "%s",
		"role": "%s",
		"token": "%s"
	}`,
		userData["id"],
		userData["username"],
		userData["email"],
		userData["phone"],
		userData["role"],
		tokenString)

	// Logs para confirmar que el token está incluido
	log.Printf("RESPUESTA ENVIADA (longitud: %d):", len(hardcodedResponse))
	log.Printf("%s", hardcodedResponse)

	// Enviar respuesta directamente sin serialización adicional
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hardcodedResponse))
}
