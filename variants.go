package core

import (
	"fmt"
	"strings"
	"trickster/output"
	"trickster/transforms"
	"trickster/utils"
)

// RunVariants es el punto de entrada del Módulo 2.
// Guía al usuario con preguntas para construir variantes personalizadas.
func RunVariants() {
	fmt.Println("\n\033[1m[ MÓDULO 2 - VARIANTES GUIADAS ]\033[0m\n")
	utils.Info("Responde las preguntas para personalizar las variantes generadas.")
	fmt.Println()

	// 1. Palabras base que el usuario quiere usar
	rawInput := utils.AskStringRequired("Ingresa las palabras base separadas por comas (ej: carlos,perro,boca)")
	baseParts := strings.Split(rawInput, ",")

	// Limpiar espacios de cada parte
	var bases []string
	for _, p := range baseParts {
		p = strings.TrimSpace(p)
		if p != "" {
			bases = append(bases, p)
		}
	}

	// 2. Preguntar qué transformaciones aplicar
	fmt.Println()
	utils.Info("¿Qué transformaciones deseas aplicar?")
	applyUpper := askYesNo("¿Agregar versión en MAYÚSCULAS?")
	applyCap := askYesNo("¿Agregar versión Capitalizada (Primera letra mayúscula)?")
	applyLeet := askYesNo("¿Agregar versión l33tspeak (ej: carlos → c4rl0s)?")
	applyReverse := askYesNo("¿Agregar versión al revés (ej: carlos → solrac)?")

	// 3. Sufijos personalizados
	fmt.Println()
	customSuffix := utils.AskOptional("¿Agregar sufijo personalizado? (ej: 2024, @empresa)")
	useDefaultSuffix := askYesNo("¿Agregar sufijos numéricos comunes (1, 123, 1234, !)?")

	// 4. Generar
	utils.Info("Generando variantes...")
	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	for _, base := range bases {
		lower := transforms.ToLower(base)
		add(lower)

		if applyUpper {
			add(transforms.ToUpper(base))
		}
		if applyCap {
			add(transforms.Capitalize(base))
		}
		if applyLeet {
			add(transforms.Leet(lower))
			if applyCap {
				add(transforms.Leet(transforms.Capitalize(base)))
			}
		}
		if applyReverse {
			add(transforms.Reverse(lower))
		}

		// Sufijo personalizado
		if customSuffix != "" {
			add(lower + customSuffix)
			if applyCap {
				add(transforms.Capitalize(base) + customSuffix)
			}
		}

		// Sufijos comunes
		if useDefaultSuffix {
			for _, suf := range transforms.CommonSuffixes {
				add(lower + suf)
				add(transforms.Capitalize(base) + suf)
			}
		}
	}

	// 5. Guardar
	outputPath := utils.AskStringRequired("Ruta de salida (ej: /home/user/variantes.txt)")

	if err := output.WriteWordlist(result, outputPath); err != nil {
		utils.Error("Error al guardar: " + err.Error())
		return
	}

	output.PrintStats(outputPath, len(result))
}

// askYesNo hace una pregunta sí/no y devuelve true si el usuario responde s/si/yes/y
func askYesNo(question string) bool {
	answer := utils.AskString(question + " [s/n]")
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "s" || answer == "si" || answer == "sí" || answer == "y" || answer == "yes"
}
