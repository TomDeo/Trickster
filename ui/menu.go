package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"trickster/core"
)

// Colores ANSI para la terminal
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// banner imprime el logo de Trickster al inicio
func banner() {
	fmt.Println(colorRed + colorBold)
	fmt.Println(`
 ████████╗██████╗ ██╗ ██████╗██╗  ██╗███████╗████████╗███████╗██████╗ 
 ╚══██╔══╝██╔══██╗██║██╔════╝██║ ██╔╝██╔════╝╚══██╔══╝██╔════╝██╔══██╗
    ██║   ██████╔╝██║██║     █████╔╝ ███████╗   ██║   █████╗  ██████╔╝
    ██║   ██╔══██╗██║██║     ██╔═██╗ ╚════██║   ██║   ██╔══╝  ██╔══██╗
    ██║   ██║  ██║██║╚██████╗██║  ██╗███████║   ██║   ███████╗██║  ██║
    ╚═╝   ╚═╝  ╚═╝╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝
	`)
	fmt.Println(colorReset)
	fmt.Println(colorCyan + "  [ Wordlist Generator Tool ]" + colorReset)
	fmt.Println(colorYellow + "  v0.1 - by you\n" + colorReset)
}

// menu imprime las opciones disponibles
func menu() {
	fmt.Println(colorBold + "  MENU PRINCIPAL" + colorReset)
	fmt.Println(colorGreen + "  [1]" + colorReset + " Crear máscaras desde wordlist")
	fmt.Println(colorGreen + "  [2]" + colorReset + " Crear variantes guiadas")
	fmt.Println(colorGreen + "  [3]" + colorReset + " Perfil avanzado (modo completo)")
	fmt.Println(colorRed + "  [0]" + colorReset + " Salir")
	fmt.Println()
}

// prompt muestra el símbolo de entrada estilo Metasploit
func prompt() string {
	fmt.Print(colorRed + "trickster" + colorReset + " > ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Run es la función principal que lanza el loop del menú
func Run() {
	banner()

	for {
		menu()
		opcion := prompt()

		switch opcion {
		case "1":
			// Módulo 1: máscaras básicas desde wordlist existente
			core.RunMasks()

		case "2":
			// Módulo 2: variantes guiadas por preguntas al usuario
			core.RunVariants()

		case "3":
			// Módulo 3: perfil completo con datos personales
			// Este módulo incluye todo lo que hacen el 1 y el 2
			core.RunProfiler()

		case "0":
			fmt.Println(colorYellow + "\n[*] Saliendo de Trickster. Hasta luego.\n" + colorReset)
			os.Exit(0)

		default:
			fmt.Println(colorRed + "\n[!] Opción no válida. Intenta de nuevo.\n" + colorReset)
		}
	}
}
