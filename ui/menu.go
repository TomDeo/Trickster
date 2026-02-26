package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"trickster/core"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

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

func menu() {
	fmt.Println(colorBold + "  MENU PRINCIPAL" + colorReset)
	fmt.Println(colorGreen + "  [1]" + colorReset + " Generar contraseñas")
	fmt.Println(colorRed + "  [0]" + colorReset + " Salir")
	fmt.Println()
}

func menuPasswords() {
	fmt.Println(colorBold + "  GENERAR CONTRASEÑAS" + colorReset)
	fmt.Println(colorGreen + "  [1]" + colorReset + " Crear máscaras desde wordlist")
	fmt.Println(colorGreen + "  [2]" + colorReset + " Crear variantes guiadas")
	fmt.Println(colorGreen + "  [3]" + colorReset + " Perfil avanzado (modo completo)")
	fmt.Println(colorYellow + "  [0]" + colorReset + " Volver")
	fmt.Println()
}

func prompt() string {
	fmt.Print(colorRed + "trickster" + colorReset + " > ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func runPasswordsMenu() {
	for {
		menuPasswords()
		opcion := prompt()
		switch opcion {
		case "1":
			core.RunMasks()
		case "2":
			core.RunVariants()
		case "3":
			core.RunProfiler()
		case "0":
			return // vuelve al menú principal
		default:
			fmt.Println(colorRed + "\n[!] Opción no válida. Intenta de nuevo.\n" + colorReset)
		}
	}
}

func Run() {
	banner()
	for {
		menu()
		opcion := prompt()
		switch opcion {
		case "1":
			runPasswordsMenu()
		case "0":
			fmt.Println(colorYellow + "\n[*] Saliendo de Trickster. Hasta luego.\n" + colorReset)
			os.Exit(0)
		default:
			fmt.Println(colorRed + "\n[!] Opción no válida. Intenta de nuevo.\n" + colorReset)
		}
	}
}
