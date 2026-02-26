package core

import (
	"fmt"
	"trickster/output"
	"trickster/transforms"
	"trickster/utils"
)

// RunMasks es el punto de entrada del Módulo 1.
// Toma una wordlist existente y genera variantes con máscaras alfanuméricas.
func RunMasks() {
	fmt.Println("\n\033[1m[ MÓDULO 1 - MÁSCARAS DESDE WORDLIST ]\033[0m\n")
	utils.Info("Este módulo toma tu wordlist y genera cientos de variantes automáticamente.")
	fmt.Println()

	// 1. Pedir ruta del archivo de entrada
	inputPath := utils.AskStringRequired("Ruta de tu wordlist de entrada (ej: /home/user/palabras.txt)")

	// 2. Leer el archivo
	words, err := utils.ReadWordlistFile(inputPath)
	if err != nil {
		utils.Error("No se pudo leer el archivo: " + err.Error())
		return
	}
	utils.Success(fmt.Sprintf("Cargadas %d palabras base.", len(words)))

	// 3. Generar variantes para cada palabra
	utils.Info("Generando variantes...")
	var allVariants []string

	for _, word := range words {
		variants := transforms.AllVariants(word)
		allVariants = append(allVariants, variants...)
	}

	// 4. Eliminar duplicados
	allVariants = utils.Deduplicate(allVariants)

	// 5. Pedir ruta de salida
	outputPath := utils.AskStringRequired("Ruta de salida para la wordlist generada (ej: /home/user/output.txt)")

	// 6. Escribir al archivo
	if err := output.WriteWordlist(allVariants, outputPath); err != nil {
		utils.Error("Error al guardar: " + err.Error())
		return
	}

	output.PrintStats(outputPath, len(allVariants))
}
