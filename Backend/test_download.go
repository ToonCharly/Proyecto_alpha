// Test para verificar descarga de plantilla Word
// Ejecutar con: go run test_download.go

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func testDownload() {
	// Test directo del endpoint
	fmt.Println("🧪 Iniciando test de descarga de plantilla Word...")

	// URL del endpoint
	url := "http://localhost:8080/api/plantillas/ejemplo"

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Hacer request
	fmt.Printf("📡 Haciendo request a: %s\n", url)
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("❌ Error en request: %v", err)
	}
	defer resp.Body.Close()

	// Verificar status code
	fmt.Printf("📊 Status Code: %d %s\n", resp.StatusCode, resp.Status)
	if resp.StatusCode != 200 {
		log.Fatalf("❌ Status code incorrecto: %d", resp.StatusCode)
	}

	// Mostrar headers de respuesta
	fmt.Println("📋 Headers de respuesta:")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Crear archivo de salida
	outputFile := "plantilla_test_download.docx"
	out, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("❌ Error al crear archivo de salida: %v", err)
	}
	defer out.Close()

	// Copiar contenido
	fmt.Println("💾 Descargando archivo...")
	bytesWritten, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("❌ Error al copiar contenido: %v", err)
	}

	fmt.Printf("✅ Descarga completada: %d bytes\n", bytesWritten)

	// Verificar el archivo descargado
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		log.Fatalf("❌ Error al verificar archivo descargado: %v", err)
	}

	fmt.Printf("📁 Archivo creado: %s\n", outputFile)
	fmt.Printf("📏 Tamaño final: %d bytes\n", fileInfo.Size())
	fmt.Printf("🕒 Fecha de creación: %v\n", fileInfo.ModTime())

	// Verificar header del archivo (debe ser ZIP)
	file, err := os.Open(outputFile)
	if err != nil {
		log.Fatalf("❌ Error al abrir archivo para verificación: %v", err)
	}
	defer file.Close()

	header := make([]byte, 4)
	n, err := file.Read(header)
	if err != nil || n < 4 {
		log.Fatalf("❌ Error al leer header del archivo: %v", err)
	}

	fmt.Printf("🔍 Header del archivo: %02X %02X %02X %02X\n", header[0], header[1], header[2], header[3])

	if header[0] == 0x50 && header[1] == 0x4B {
		fmt.Println("✅ Archivo ZIP válido (formato .docx correcto)")
	} else {
		fmt.Println("❌ Archivo NO es ZIP válido")
	}

	// Comparar con archivo original
	originalFile := "public/assets/plantilla_ejemplo_factura.docx"
	originalInfo, err := os.Stat(originalFile)
	if err != nil {
		fmt.Printf("⚠️ No se puede acceder al archivo original: %v\n", err)
	} else {
		fmt.Printf("📊 Comparación:\n")
		fmt.Printf("  Archivo original: %d bytes\n", originalInfo.Size())
		fmt.Printf("  Archivo descargado: %d bytes\n", fileInfo.Size())

		if originalInfo.Size() == fileInfo.Size() {
			fmt.Println("✅ Tamaños coinciden perfectamente")
		} else {
			fmt.Printf("❌ PROBLEMA: Los tamaños NO coinciden (diferencia: %d bytes)\n",
				originalInfo.Size()-fileInfo.Size())
		}
	}

	fmt.Println("🎯 Test completado.")
}
