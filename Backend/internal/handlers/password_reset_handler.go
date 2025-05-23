package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

// Almacenamiento temporal de tokens y control de intentos
var (
	passwordResetTokens   = make(map[string]PasswordResetInfo)
	passwordResetAttempts = make(map[string][]time.Time)
	resetMutex            = &sync.Mutex{}
)

// PasswordResetInfo almacena información sobre un token de recuperación
type PasswordResetInfo struct {
	UserID    int
	Username  string
	Email     string
	Token     string
	ExpiresAt time.Time
}

// ResetPasswordRequestHandler maneja las solicitudes de recuperación de contraseña
func ResetPasswordRequestHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Solo permitir método POST
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Decodificar la solicitud
		var request struct {
			Email string `json:"email"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
			return
		}

		// Validación básica de email
		if request.Email == "" || !strings.Contains(request.Email, "@") {
			http.Error(w, "Correo electrónico inválido", http.StatusBadRequest)
			return
		}

		// Verificar límites de intentos
		resetMutex.Lock()
		now := time.Now()
		attempts := passwordResetAttempts[request.Email]

		// Limpiar intentos antiguos (más de 24 horas)
		var recentAttempts []time.Time
		for _, t := range attempts {
			if now.Sub(t) < 24*time.Hour {
				recentAttempts = append(recentAttempts, t)
			}
		}

		// Verificar si ha habido más de 5 intentos en las últimas 24 horas
		if len(recentAttempts) >= 5 {
			// Verificar si el último intento fue hace menos de 30 minutos
			if len(recentAttempts) > 0 && now.Sub(recentAttempts[len(recentAttempts)-1]) < 30*time.Minute {
				resetMutex.Unlock()
				http.Error(w, "Has excedido el límite de intentos. Intenta de nuevo más tarde.",
					http.StatusTooManyRequests)
				return
			}
		}

		// Registrar este intento
		passwordResetAttempts[request.Email] = append(recentAttempts, now)
		resetMutex.Unlock()

		// Buscar usuario por email
		var userID int
		var username string
		query := "SELECT id, username FROM usuarios WHERE email = ?"
		err = db.QueryRow(query, request.Email).Scan(&userID, &username)

		// Si el usuario no existe, enviar respuesta genérica por seguridad
		if err != nil {
			if err == sql.ErrNoRows {
				// Por seguridad, no revelar si el email existe o no
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Si el correo existe, recibirás instrucciones para restablecer tu contraseña",
				})
				return
			}

			log.Printf("Error al buscar usuario: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		// Generar token aleatorio
		token, err := generateRandomToken()
		if err != nil {
			log.Printf("Error al generar token: %v", err)
			http.Error(w, "Error al procesar la solicitud", http.StatusInternalServerError)
			return
		}

		// Guardar información de recuperación
		passwordResetTokens[token] = PasswordResetInfo{
			UserID:    userID,
			Username:  username,
			Email:     request.Email,
			Token:     token,
			ExpiresAt: time.Now().Add(30 * time.Minute), // Token válido por 30 minutos
		}

		// Enviar correo (en background)
		go sendPasswordResetEmail(request.Email, username, token)

		// Responder con éxito
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Si el correo existe, recibirás instrucciones para restablecer tu contraseña",
		})
	}
}

// ResetPasswordHandler maneja el restablecimiento de contraseña
func ResetPasswordHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Solo permitir método POST
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Decodificar la solicitud
		var request struct {
			Token       string `json:"token"`
			NewPassword string `json:"newPassword"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
			return
		}

		// Validación básica
		if request.Token == "" || request.NewPassword == "" {
			http.Error(w, "Datos incompletos", http.StatusBadRequest)
			return
		}

		// Verificar token
		resetInfo, exists := passwordResetTokens[request.Token]
		if !exists {
			http.Error(w, "Token inválido o expirado", http.StatusBadRequest)
			return
		}

		// Verificar si el token ha expirado
		if time.Now().After(resetInfo.ExpiresAt) {
			delete(passwordResetTokens, request.Token) // Eliminar token expirado
			http.Error(w, "Token expirado", http.StatusBadRequest)
			return
		}

		// Actualizar contraseña en la base de datos
		_, err = db.Exec("UPDATE usuarios SET password = ? WHERE id = ?", request.NewPassword, resetInfo.UserID)
		if err != nil {
			log.Printf("Error al actualizar contraseña: %v", err)
			http.Error(w, "Error al actualizar contraseña", http.StatusInternalServerError)
			return
		}

		// Eliminar token utilizado
		delete(passwordResetTokens, request.Token)

		// Responder con éxito
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Contraseña actualizada correctamente",
		})
	}
}

// Función para generar un token aleatorio
func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Función para enviar correo de recuperación
func sendPasswordResetEmail(email, username, token string) {
	// Configuración para diferentes proveedores
	var smtpConfig struct {
		Host     string
		Port     string
		Username string
		Password string
	}

	// Define qué proveedor usar (1=Gmail, 2=Dominio propio, 3=Outlook/Hotmail)
	proveedorCorreo := 1 // CAMBIA ESTO según el proveedor que quieras usar

	// Configuraciones disponibles
	switch proveedorCorreo {
	case 1: // Gmail
		smtpConfig.Host = "smtp.gmail.com"
		smtpConfig.Port = "587"
		smtpConfig.Username = "verificadorcontra00@gmail.com" // Correo configurado
		smtpConfig.Password = "acyhcauwbiezbiql"              // Contraseña (sin espacios)
	case 2: // Dominio propio
		smtpConfig.Host = "smtp.tudominio.com"     // CAMBIA ESTO
		smtpConfig.Port = "587"                    // Verifica el puerto con tu proveedor
		smtpConfig.Username = "info@tudominio.com" // CAMBIA ESTO
		smtpConfig.Password = "tu_contraseña"      // CAMBIA ESTO
	case 3: // Outlook/Hotmail
		smtpConfig.Host = "smtp-mail.outlook.com"
		smtpConfig.Port = "587"
		smtpConfig.Username = "tu_correo@outlook.com" // CAMBIA ESTO (o @hotmail.com)
		smtpConfig.Password = "tu_contraseña"         // CAMBIA ESTO
	default:
		log.Printf("Error: Proveedor de correo no configurado correctamente")
		return
	}

	// URL de tu aplicación
	appURL := "http://localhost:5173" // URL correcta donde corre tu React
	resetLink := appURL + "/restablecer-password?token=" + token

	// Contenido del correo
	to := []string{email}
	subject := "Recuperación de contraseña"
	body := "Hola " + username + ",\n\n" +
		"Has solicitado restablecer tu contraseña. Para continuar, haz clic en el siguiente enlace:\n\n" +
		resetLink + "\n\n" +
		"Este enlace expirará en 30 minutos.\n\n" +
		"Si no solicitaste este cambio, puedes ignorar este correo.\n\n" +
		"Saludos,\n" +
		""

	message := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body)

	// Autenticación
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	// Envío del correo
	err := smtp.SendMail(smtpConfig.Host+":"+smtpConfig.Port, auth, smtpConfig.Username, to, message)
	if err != nil {
		log.Printf("Error al enviar correo de recuperación: %v", err)
	} else {
		log.Printf("Correo de recuperación enviado a: %s", email)
	}
}

// Función para limpiar tokens expirados (llamar periódicamente)
func LimpiarTokensExpirados() {
	ahora := time.Now()
	resetMutex.Lock()
	defer resetMutex.Unlock()

	for token, info := range passwordResetTokens {
		if ahora.After(info.ExpiresAt) {
			delete(passwordResetTokens, token)
		}
	}
}
