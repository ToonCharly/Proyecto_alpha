package handlers

import (
    "net/http"
    "io/ioutil"
    "carlos/Facts/Backend/internal/utils"
    "database/sql"
)

// AltaDatosFiscalesHandler maneja el alta de datos fiscales y guarda binario y PEM
func AltaDatosFiscalesHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Parsear multipart/form-data
        err := r.ParseMultipartForm(10 << 20) // 10MB
        if err != nil {
            http.Error(w, "Error al parsear formulario: "+err.Error(), http.StatusBadRequest)
            return
        }

        archivoCer, _, err := r.FormFile("archivo_cer")
        if err != nil {
            http.Error(w, "Error leyendo archivo .cer: "+err.Error(), http.StatusBadRequest)
            return
        }
        defer archivoCer.Close()
        cerBytes, _ := ioutil.ReadAll(archivoCer)

        archivoKey, _, err := r.FormFile("archivo_key")
        if err != nil {
            http.Error(w, "Error leyendo archivo .key: "+err.Error(), http.StatusBadRequest)
            return
        }
        defer archivoKey.Close()
        keyBytes, _ := ioutil.ReadAll(archivoKey)

        // Convertir a PEM
        cerPEM, err := utils.SavePEMFromDER(cerBytes, "CERTIFICATE", "")
        if err != nil {
            http.Error(w, "Error convirtiendo .cer a PEM: "+err.Error(), http.StatusInternalServerError)
            return
        }
        keyPEM, err := utils.SavePEMFromDER(keyBytes, "PRIVATE KEY", "")
        if err != nil {
            http.Error(w, "Error convirtiendo .key a PEM: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Ejemplo de SQL para guardar ambos formatos
        // Asume que recibes también los demás campos requeridos (rfc, razon_social, etc.)
        query := `INSERT INTO datos_fiscales (rfc, razon_social, archivo_cer, archivo_key, archivo_cer_pem, archivo_key_pem) VALUES (?, ?, ?, ?, ?, ?)`
        _, err = db.Exec(query, r.FormValue("rfc"), r.FormValue("razon_social"), cerBytes, keyBytes, cerPEM, keyPEM)
        if err != nil {
            http.Error(w, "Error guardando en base de datos: "+err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        w.Write([]byte("Datos fiscales guardados correctamente con PEM"))
    }
}
