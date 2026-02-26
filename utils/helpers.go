package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AskString hace una pregunta al usuario y devuelve la respuesta como string.
// Es la función base para todos los módulos interactivos.
func AskString(question string) string {
	fmt.Print("\033[33m[?]\033[0m " + question + ": ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	return strings.TrimSpace(answer)
}

// AskStringRequired igual que AskString pero repite hasta que el usuario ingrese algo
func AskStringRequired(question string) string {
	for {
		val := AskString(question)
		if val != "" {
			return val
		}
		fmt.Println("\033[31m[!] Este campo es requerido.\033[0m")
	}
}

// AskOptional igual que AskString pero indica que es opcional
func AskOptional(question string) string {
	return AskString(question + " (opcional, Enter para omitir)")
}

// ReadWordlistFile lee un archivo de texto línea por línea y devuelve un slice de strings.
// Ignora líneas vacías y espacios extra.
func ReadWordlistFile(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			words = append(words, line)
		}
	}

	return words, scanner.Err()
}

// DeduplicateAndSort elimina duplicados de un slice sin ordenar
// Usa un map como set - O(n) en tiempo
func Deduplicate(words []string) []string {
	seen := make(map[string]bool, len(words))
	result := make([]string, 0, len(words))

	for _, w := range words {
		if !seen[w] {
			seen[w] = true
			result = append(result, w)
		}
	}
	return result
}

// Info imprime un mensaje informativo en cyan
func Info(msg string) {
	fmt.Println("\033[36m[*]\033[0m " + msg)
}

// Success imprime un mensaje de éxito en verde
func Success(msg string) {
	fmt.Println("\033[32m[+]\033[0m " + msg)
}

// Warn imprime una advertencia en amarillo
func Warn(msg string) {
	fmt.Println("\033[33m[!]\033[0m " + msg)
}

// Error imprime un error en rojo
func Error(msg string) {
	fmt.Println("\033[31m[!]\033[0m " + msg)
}
