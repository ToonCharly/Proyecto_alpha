package models

import (
	"fmt"
	"sync"
	"time"
)

// Generador de folios thread-safe sin dependencia de base de datos
type FolioGenerator struct {
	mu      sync.Mutex
	counter int64
}

var (
	folioGen     *FolioGenerator
	folioGenOnce sync.Once
)

// GetFolioGenerator obtiene la instancia singleton del generador de folios
func GetFolioGenerator() *FolioGenerator {
	folioGenOnce.Do(func() {
		folioGen = &FolioGenerator{
			counter: 0,
		}
	})
	return folioGen
}

// GenerarFolio genera un folio único basado en timestamp y contador incremental
func (fg *FolioGenerator) GenerarFolio(serie string) (string, error) {
	fg.mu.Lock()
	defer fg.mu.Unlock()

	// Incrementar contador
	fg.counter++

	// Generar folio usando timestamp + contador para garantizar unicidad
	timestamp := time.Now().Unix() // Timestamp en segundos
	folio := fmt.Sprintf("%s%d%04d", serie, timestamp, fg.counter%10000)

	return folio, nil
}

// GenerarFolioSimple genera un folio simple incremental
func (fg *FolioGenerator) GenerarFolioSimple(serie string) (string, error) {
	fg.mu.Lock()
	defer fg.mu.Unlock()

	fg.counter++

	// Folio simple: Serie + contador con padding
	folio := fmt.Sprintf("%s%06d", serie, fg.counter)

	return folio, nil
}

// ResetCounter reinicia el contador (solo para pruebas)
func (fg *FolioGenerator) ResetCounter() {
	fg.mu.Lock()
	defer fg.mu.Unlock()
	fg.counter = 0
}

// GetCurrentCounter obtiene el valor actual del contador
func (fg *FolioGenerator) GetCurrentCounter() int64 {
	fg.mu.Lock()
	defer fg.mu.Unlock()
	return fg.counter
}
