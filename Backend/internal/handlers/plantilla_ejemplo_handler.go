package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// PlantillaEjemploHandler - Versión robusta con validación completa
func PlantillaEjemploHandler(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Manejar preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Solo permitir GET
	if r.Method != http.MethodGet {
		log.Printf("❌ Método no permitido: %s", r.Method)
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Construir ruta del archivo
	dir, _ := os.Getwd()
	filePath := filepath.Join(dir, "public", "assets", "plantilla_ejemplo_factura.docx")

	log.Printf("🔍 [PlantillaEjemplo] Solicitando archivo desde: %s", filePath)
	log.Printf("🌐 [PlantillaEjemplo] User-Agent: %s", r.Header.Get("User-Agent"))

	// Verificar que el archivo existe
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("❌ [PlantillaEjemplo] Error: Archivo no encontrado en %s: %v", filePath, err)
		http.Error(w, "Archivo no encontrado", http.StatusNotFound)
		return
	}

	log.Printf("✅ [PlantillaEjemplo] Archivo encontrado - Tamaño: %d bytes, Modificado: %v",
		fileInfo.Size(), fileInfo.ModTime())

	// Verificar tamaño esperado (debe ser alrededor de 13380 bytes)
	if fileInfo.Size() < 10000 || fileInfo.Size() > 50000 {
		log.Printf("❌ [PlantillaEjemplo] Tamaño de archivo sospechoso: %d bytes (esperado: ~13380)", fileInfo.Size())
		http.Error(w, "Archivo corrupto", http.StatusInternalServerError)
		return
	}

	// Leer el archivo completo en memoria para verificar integridad
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("❌ [PlantillaEjemplo] Error al leer archivo: %v", err)
		http.Error(w, "Error al leer archivo", http.StatusInternalServerError)
		return
	}

	// Verificar que los datos leídos coinciden con el tamaño del archivo
	if int64(len(fileData)) != fileInfo.Size() {
		log.Printf("❌ [PlantillaEjemplo] Discrepancia en tamaño: archivo=%d bytes, leído=%d bytes",
			fileInfo.Size(), len(fileData))
		http.Error(w, "Error de integridad del archivo", http.StatusInternalServerError)
		return
	}

	// Verificar que es un archivo ZIP válido (los .docx son ZIP)
	if len(fileData) < 4 || fileData[0] != 0x50 || fileData[1] != 0x4B {
		log.Printf("❌ [PlantillaEjemplo] No es un archivo ZIP válido. Header: %x %x",
			fileData[0], fileData[1])
		http.Error(w, "Archivo Word corrupto", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [PlantillaEjemplo] Archivo verificado: %d bytes, header ZIP válido", len(fileData))

	// Configurar headers HTTP para descarga
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	w.Header().Set("Content-Disposition", "attachment; filename=\"plantilla_ejemplo_factura.docx\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(fileData)))

	// Headers de cache
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("ETag", fmt.Sprintf("\"%d-%d\"", fileInfo.ModTime().Unix(), fileInfo.Size()))

	log.Printf("🌐 [PlantillaEjemplo] Headers configurados para descarga")

	// Escribir el archivo completo al response
	bytesWritten, err := w.Write(fileData)
	if err != nil {
		log.Printf("❌ [PlantillaEjemplo] Error al enviar archivo: %v", err)
		return
	}

	if bytesWritten != len(fileData) {
		log.Printf("⚠️ [PlantillaEjemplo] ADVERTENCIA: Bytes enviados (%d) != bytes del archivo (%d)",
			bytesWritten, len(fileData))
	} else {
		log.Printf("✅ [PlantillaEjemplo] Archivo enviado exitosamente: %d bytes", bytesWritten)
	}
}

// verificarArchivoWord - Función simplificada para verificar validez del archivo
func verificarArchivoWord(filePath string) error {
	// Verificar que existe
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("archivo no encontrado: %v", err)
	}

	// Verificar que no está vacío
	if fileInfo.Size() == 0 {
		return fmt.Errorf("archivo vacío")
	}

	// Verificar que tiene un tamaño mínimo razonable para un .docx
	if fileInfo.Size() < 1000 {
		return fmt.Errorf("archivo muy pequeño (%d bytes), posiblemente corrupto", fileInfo.Size())
	}

	// Abrir y verificar los primeros bytes
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("no se puede abrir: %v", err)
	}
	defer file.Close()

	// Leer header (los .docx son archivos ZIP)
	header := make([]byte, 4)
	n, err := file.Read(header)
	if err != nil || n < 4 {
		return fmt.Errorf("no se puede leer header: %v", err)
	}

	// Verificar signatura ZIP (PK)
	if header[0] != 0x50 || header[1] != 0x4B {
		return fmt.Errorf("no es un archivo ZIP válido (header: %x %x)", header[0], header[1])
	}

	return nil
}

