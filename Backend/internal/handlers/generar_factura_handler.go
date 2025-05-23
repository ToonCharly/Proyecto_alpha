package handlers

import (
    "bytes"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "strings"

    "carlos/Facts/Backend/internal/models"
    "carlos/Facts/Backend/internal/services"
)

func GenerarFacturaHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    var factura models.Factura
    var plantillaBytes []byte

    // Detectar el tipo de contenido
    contentType := r.Header.Get("Content-Type")
    if contentType == "" {
        http.Error(w, "Content-Type no especificado", http.StatusBadRequest)
        return
    }

    if strings.Contains(contentType, "multipart/form-data") {
        // Procesar como multipart/form-data
        if err := r.ParseMultipartForm(10 << 20); err != nil {
            log.Printf("Error al parsear multipart form: %v", err)
            http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
            return
        }

        // Leer el campo "datos" que contiene el JSON de la factura
        facturaData := r.FormValue("datos")
        if facturaData == "" {
            http.Error(w, "No se encontraron datos de factura", http.StatusBadRequest)
            return
        }

        if err := json.Unmarshal([]byte(facturaData), &factura); err != nil {
            log.Printf("Error al decodificar JSON en multipart: %v", err)
            http.Error(w, "Error al procesar los datos: "+err.Error(), http.StatusBadRequest)
            return
        }

        // Leer el archivo de la plantilla
        plantillaFile, _, err := r.FormFile("plantilla")
        if err == nil {
            defer plantillaFile.Close()
            plantillaBytes, err = io.ReadAll(plantillaFile)
            if err != nil {
                log.Printf("Error al leer la plantilla: %v", err)
                http.Error(w, "Error al leer la plantilla", http.StatusInternalServerError)
                return
            }
        }
    } else if strings.Contains(contentType, "application/json") {
        // Procesar como JSON simple
        body, err := io.ReadAll(r.Body)
        if err != nil {
            log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
            http.Error(w, "Error al leer la solicitud", http.StatusBadRequest)
            return
        }

        if err := json.Unmarshal(body, &factura); err != nil {
            log.Printf("Error al decodificar JSON: %v", err)
            http.Error(w, "Error al procesar los datos: "+err.Error(), http.StatusBadRequest)
            return
        }
    }

    // Procesar la plantilla si se proporcionó
    var pdfBuffer *bytes.Buffer
    var err error
    if len(plantillaBytes) > 0 {
        pdfBuffer, err = services.ProcesarPlantilla(factura, plantillaBytes)
        if err != nil {
            log.Printf("Error al procesar plantilla: %v", err)
            http.Error(w, "Error al procesar la plantilla", http.StatusInternalServerError)
            return
        }
    } else {
        // Generar el PDF sin plantilla
        pdfBuffer, err = services.GenerarPDF(factura, nil)
        if err != nil {
            log.Printf("Error al generar PDF: %v", err)
            http.Error(w, "Error al generar la factura", http.StatusInternalServerError)
            return
        }
    }

    // Generar XML para la factura
    xmlBytes, err := services.GenerarXML(factura)
    if err != nil {
        log.Printf("Error al generar XML: %v", err)
        http.Error(w, "Error al generar XML de la factura", http.StatusInternalServerError)
        return
    }

    // Usar la función CrearZIP para empaquetar PDF y XML
    zipBuffer, err := services.CrearZIP(pdfBuffer.Bytes(), xmlBytes)
    if err != nil {
        log.Printf("Error al crear ZIP: %v", err)
        http.Error(w, "Error al crear archivo ZIP", http.StatusInternalServerError)
        return
    }

    // Enviar el archivo ZIP como respuesta
    w.Header().Set("Content-Type", "application/zip")
    w.Header().Set("Content-Disposition", "attachment; filename=factura.zip")
    _, err = w.Write(zipBuffer.Bytes())
    if err != nil {
        log.Printf("Error al enviar archivo ZIP: %v", err)
    }
}