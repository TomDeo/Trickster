package output

import (
	"bufio"
	"fmt"
	"os"
)

// WriteWordlist escribe un slice de strings a un archivo .txt, una palabra por línea.
// Usa bufio para escritura eficiente (no carga todo en memoria de una vez).
func WriteWordlist(words []string, filepath string) error {
	// Crear o sobreescribir el archivo
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("no se pudo crear el archivo: %w", err)
	}
	defer file.Close()

	// bufio.Writer agrupa escrituras pequeñas en bloques grandes → mucho más rápido
	writer := bufio.NewWriterSize(file, 1024*64) // buffer de 64KB

	for _, word := range words {
		// Escribir cada palabra seguida de salto de línea
		if _, err := fmt.Fprintln(writer, word); err != nil {
			return fmt.Errorf("error al escribir palabra: %w", err)
		}
	}

	// Importante: flush vuelca el buffer al disco al final
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error al finalizar escritura: %w", err)
	}

	return nil
}

// WriteWordlistStream es igual pero recibe un canal (channel) de strings.
// Útil cuando los datos se generan en tiempo real (para wordlists muy grandes).
func WriteWordlistStream(words <-chan string, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("no se pudo crear el archivo: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 1024*64)

	// Lee del canal hasta que se cierre
	for word := range words {
		if _, err := fmt.Fprintln(writer, word); err != nil {
			return err
		}
	}

	return writer.Flush()
}

// PrintStats imprime un resumen al terminar
func PrintStats(filepath string, count int) {
	fmt.Printf("\n\033[32m[+] Wordlist guardada en: %s\033[0m\n", filepath)
	fmt.Printf("\033[32m[+] Total de palabras generadas: %d\033[0m\n\n", count)
}
