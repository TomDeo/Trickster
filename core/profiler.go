package core

import (
	"fmt"
	"strings"
	"trickster/output"
	"trickster/transforms"
	"trickster/utils"
)

// Profile contiene todos los datos personales del objetivo
type Profile struct {
	Nombre          string
	Apellido        string
	DNI             string
	FechaNacimiento string // formato: DDMMYYYY o DDMMYY
	Dia             string // extraído de la fecha
	Mes             string
	Anio            string
	AnioCorto       string // últimos 2 dígitos
	EquipoFutbol    string
	Edad            string
	Ciudad          string
	OldPass1        string
	OldPass2        string
	OldPass3        string
}

// RunProfiler es el punto de entrada del Módulo 3.
// Recolecta datos personales y genera una wordlist muy robusta.
func RunProfiler() {
	fmt.Println("\n\033[1m[ MÓDULO 3 - PERFIL AVANZADO ]\033[0m\n")
	utils.Info("Ingresa los datos del objetivo. Los campos opcionales pueden dejarse en blanco.")
	utils.Warn("Omitir campos reduce la cantidad de combinaciones generadas.")
	fmt.Println()

	// ---- Recolección de datos ----
	p := Profile{}

	p.Nombre = utils.AskOptional("Nombre")
	p.Apellido = utils.AskOptional("Apellido")
	p.DNI = utils.AskOptional("DNI / Cédula / ID")

	fechaRaw := utils.AskOptional("Fecha de nacimiento (DDMMAAAA, ej: 15031990)")
	p.FechaNacimiento = strings.TrimSpace(fechaRaw)
	// Si la fecha tiene 8 dígitos, extraemos partes
	if len(p.FechaNacimiento) == 8 {
		p.Dia = p.FechaNacimiento[0:2]
		p.Mes = p.FechaNacimiento[2:4]
		p.Anio = p.FechaNacimiento[4:8]
		p.AnioCorto = p.FechaNacimiento[6:8]
	}

	p.EquipoFutbol = utils.AskOptional("Equipo de fútbol favorito")
	p.Edad = utils.AskOptional("Edad actual")
	p.Ciudad = utils.AskOptional("Ciudad")
	p.OldPass1 = utils.AskOptional("Contraseña antigua 1")
	p.OldPass2 = utils.AskOptional("Contraseña antigua 2")
	p.OldPass3 = utils.AskOptional("Contraseña antigua 3")

	fmt.Println()
	utils.Info("Procesando perfil y generando wordlist...")

	// ---- Generación de wordlist ----
	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		s = strings.TrimSpace(s)
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	// Recolectar todos los campos como tokens base
	// Solo los que no están vacíos
	tokens := collectTokens(p)

	// 1. Variantes completas de cada token individual
	for _, token := range tokens {
		for _, v := range transforms.AllVariants(token) {
			add(v)
		}
	}

	// 2. Combinaciones de 2 tokens
	for i, a := range tokens {
		for j, b := range tokens {
			if i == j {
				continue
			}
			// Combinación directa
			add(transforms.Combine(a, b))
			add(transforms.Combine(transforms.Capitalize(a), b))
			add(transforms.Combine(a, transforms.Capitalize(b)))
			add(transforms.Combine(transforms.Capitalize(a), transforms.Capitalize(b)))

			// Con separadores comunes
			for _, sep := range []string{"_", ".", "-", "@"} {
				add(transforms.CombineWith(a, sep, b))
				add(transforms.CombineWith(transforms.Capitalize(a), sep, b))
			}
		}
	}

	// 3. Combinaciones nombre + fecha (muy comunes en contraseñas reales)
	if p.Nombre != "" && p.Anio != "" {
		nombre := strings.ToLower(p.Nombre)
		add(nombre + p.Anio)
		add(nombre + p.AnioCorto)
		add(transforms.Capitalize(nombre) + p.Anio)
		add(transforms.Capitalize(nombre) + p.AnioCorto)
		add(p.Anio + nombre)
		add(nombre + p.Dia + p.Mes)
		add(nombre + p.Mes + p.Dia)
	}

	// 4. Nombre repetido (ej: carloscarlos)
	if p.Nombre != "" {
		n := strings.ToLower(p.Nombre)
		add(n + n)
		add(transforms.Capitalize(n) + n)
	}

	// 5. Año repetido (ej: 19901990)
	if p.Anio != "" {
		add(p.Anio + p.Anio)
		add(p.AnioCorto + p.AnioCorto)
	}

	// 6. Fecha completa en diferentes formatos
	if p.Dia != "" && p.Mes != "" && p.Anio != "" {
		add(p.Dia + p.Mes + p.Anio)
		add(p.Anio + p.Mes + p.Dia)
		add(p.Dia + "/" + p.Mes + "/" + p.Anio)
		add(p.Dia + "-" + p.Mes + "-" + p.Anio)
		if p.Nombre != "" {
			n := strings.ToLower(p.Nombre)
			add(n + p.Dia + p.Mes + p.Anio)
			add(transforms.Capitalize(n) + p.Dia + p.Mes + p.Anio)
		}
	}

	// 7. Variantes de contraseñas antiguas (muy alta probabilidad de reutilización con mutaciones)
	for _, oldPass := range []string{p.OldPass1, p.OldPass2, p.OldPass3} {
		if oldPass == "" {
			continue
		}
		for _, v := range transforms.AllVariants(oldPass) {
			add(v)
		}
		// Contraseña vieja + año actual (patrón frecuente de "actualización")
		if p.Anio != "" {
			add(oldPass + p.Anio)
			add(oldPass + p.AnioCorto)
		}
		// Contraseña vieja + sufijos comunes
		for _, suf := range transforms.CommonSuffixes {
			add(oldPass + suf)
		}
	}

	// 8. DNI con variantes
	if p.DNI != "" {
		add(p.DNI)
		if p.Nombre != "" {
			add(strings.ToLower(p.Nombre) + p.DNI)
			add(transforms.Capitalize(p.Nombre) + p.DNI)
		}
	}

	// ---- Guardar resultado ----
	fmt.Printf("\n\033[32m[+] Total antes de deduplicar: %d palabras\033[0m\n", len(result))

	outputPath := utils.AskStringRequired("Ruta de salida (ej: /home/user/perfil.txt)")

	if err := output.WriteWordlist(result, outputPath); err != nil {
		utils.Error("Error al guardar: " + err.Error())
		return
	}

	output.PrintStats(outputPath, len(result))
}

// collectTokens recolecta todos los campos del perfil como tokens base (no vacíos)
func collectTokens(p Profile) []string {
	var tokens []string

	fields := []string{
		p.Nombre, p.Apellido, p.DNI,
		p.FechaNacimiento, p.Dia, p.Mes, p.Anio, p.AnioCorto,
		p.EquipoFutbol, p.Edad, p.Ciudad,
	}

	seen := make(map[string]bool)
	for _, f := range fields {
		f = strings.ToLower(strings.TrimSpace(f))
		if f != "" && !seen[f] {
			seen[f] = true
			tokens = append(tokens, f)
		}
	}

	return tokens
}
