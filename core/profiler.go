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
	FechaNacimiento string
	Dia             string
	Mes             string
	Anio            string
	AnioCorto       string
	EquipoFutbol    string
	Edad            string
	Ciudad          string
	OldPass1        string
	OldPass2        string
	OldPass3        string
}

func RunProfiler() {
	fmt.Println("\n\033[1m[ MÓDULO 3 - PERFIL AVANZADO ]\033[0m\n")
	utils.Info("Ingresa los datos del objetivo. Los campos opcionales pueden dejarse en blanco.")
	utils.Warn("Omitir campos reduce la cantidad de combinaciones generadas.")
	fmt.Println()

	p := Profile{}

	p.Nombre = utils.AskOptional("Nombre")
	p.Apellido = utils.AskOptional("Apellido")
	p.DNI = utils.AskOptional("DNI / Cédula / ID")

	fechaRaw := utils.AskOptional("Fecha de nacimiento (DDMMAAAA, ej: 15031990)")
	p.FechaNacimiento = strings.TrimSpace(fechaRaw)
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

	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		s = strings.TrimSpace(s)
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	tokens := collectTokens(p)

	// 1. Variantes completas de cada token individual
	for _, token := range tokens {
		for _, v := range transforms.AllVariants(token) {
			add(v)
		}
	}

	// 2. Combinaciones de 2 tokens con separadores
	for i, a := range tokens {
		for j, b := range tokens {
			if i == j {
				continue
			}
			add(transforms.Combine(a, b))
			add(transforms.Combine(transforms.Capitalize(a), b))
			add(transforms.Combine(a, transforms.Capitalize(b)))
			add(transforms.Combine(transforms.Capitalize(a), transforms.Capitalize(b)))
			for _, sep := range []string{"_", ".", "-", "@"} {
				add(transforms.CombineWith(a, sep, b))
				add(transforms.CombineWith(transforms.Capitalize(a), sep, b))
			}
		}
	}

	// 3. Nombre repetido
	if p.Nombre != "" {
		n := strings.ToLower(p.Nombre)
		add(n + n)
		add(transforms.Capitalize(n) + n)
		add(n + transforms.Capitalize(n))
	}

	// 4. Año repetido y variantes
	if p.Anio != "" {
		add(p.Anio + p.Anio)           // 19901990
		add(p.AnioCorto + p.AnioCorto) // 9090
		add(p.Anio + p.AnioCorto)      // 199090
		add(p.AnioCorto + p.Anio)      // 901990
		for _, suf := range transforms.CommonSuffixes {
			add(p.Anio + suf)      // 1990!
			add(p.AnioCorto + suf) // 90!
		}
	}

	// 5. Fecha completa en múltiples formatos
	if p.Dia != "" && p.Mes != "" && p.Anio != "" {
		add(p.Dia + p.Mes + p.Anio)                 // 15031990
		add(p.Anio + p.Mes + p.Dia)                 // 19900315
		add(p.Dia + "/" + p.Mes + "/" + p.Anio)     // 15/03/1990
		add(p.Dia + "-" + p.Mes + "-" + p.Anio)     // 15-03-1990
		add(p.Dia + "." + p.Mes + "." + p.Anio)     // 15.03.1990
		add(p.Dia + p.Mes + p.AnioCorto)             // 150390
		add(p.Anio + p.Dia + p.Mes)                  // 19901503
		add(p.Dia + p.Mes)                            // 1503
		add(p.Mes + p.Dia)                            // 0315
		for _, suf := range transforms.CommonSuffixes {
			add(p.Dia + p.Mes + p.Anio + suf)        // 15031990!
		}
	}

	// 6. Patrones nombre + año (los más comunes en contraseñas reales)
	if p.Nombre != "" && p.Anio != "" {
		n := strings.ToLower(p.Nombre)
		nc := transforms.Capitalize(n)
		nu := transforms.ToUpper(n)

		// Básicos
		add(n + p.Anio)          // carlos1990
		add(nc + p.Anio)         // Carlos1990
		add(nu + p.Anio)         // CARLOS1990
		add(p.Anio + n)          // 1990carlos
		add(n + p.AnioCorto)     // carlos90
		add(nc + p.AnioCorto)    // Carlos90
		add(p.AnioCorto + n)     // 90carlos
		add(p.Anio + nu)         // 1990CARLOS

		// Año repetido con nombre
		add(n + p.Anio + p.Anio)         // carlos19901990
		add(n + p.AnioCorto + p.AnioCorto) // carlos9090

		// Sándwich año-nombre-año
		add(p.Anio + n + p.Anio)         // 1990carlos1990
		add(p.AnioCorto + n + p.AnioCorto) // 90carlos90

		// nombre + año + nombre
		add(n + p.Anio + n)              // carlos1990carlos
		add(nc + p.Anio + n)             // Carlos1990carlos

		// Nombre + separador + año
		for _, sep := range []string{".", "_", "-", "@"} {
			add(n + sep + p.Anio)        // carlos.1990
			add(nc + sep + p.Anio)       // Carlos.1990
			add(n + sep + p.AnioCorto)   // carlos.90
			add(p.Anio + sep + n)        // 1990.carlos
		}

		// Nombre + año + símbolo
		for _, suf := range []string{"!", "@", "#", ".", "*", "!!"} {
			add(n + p.Anio + suf)        // carlos1990!
			add(nc + p.Anio + suf)       // Carlos1990!
			add(n + p.AnioCorto + suf)   // carlos90!
		}
	}

	// 7. Nombre + fecha completa
	if p.Nombre != "" && p.Dia != "" && p.Mes != "" && p.Anio != "" {
		n := strings.ToLower(p.Nombre)
		nc := transforms.Capitalize(n)
		add(n + p.Dia + p.Mes + p.Anio)             // carlos15031990
		add(nc + p.Dia + p.Mes + p.Anio)            // Carlos15031990
		add(p.Dia + p.Mes + p.Anio + n)             // 15031990carlos
		add(n + p.Dia + p.Mes)                       // carlos1503
		add(n + p.Mes + p.Anio)                      // carlos031990
	}

	// 8. Patrones nombre + apellido
	if p.Nombre != "" && p.Apellido != "" {
		n := strings.ToLower(p.Nombre)
		a := strings.ToLower(p.Apellido)
		nc := transforms.Capitalize(n)
		ac := transforms.Capitalize(a)

		add(n + a)                                       // carlosgomez
		add(nc + ac)                                     // CarlosGomez
		add(nc + a)                                      // Carlosgomez
		add(n + "_" + a)                                 // carlos_gomez
		add(n + "." + a)                                 // carlos.gomez
		add(transforms.TruncateLeft(n, 1) + a)          // cgomez
		add(transforms.TruncateLeft(n, 1) + "_" + a)    // c_gomez
		add(transforms.TruncateLeft(n, 1) + "." + a)    // c.gomez
		add(a + n)                                       // gomezcarlos
		add(ac + nc)                                     // GomezCarlos

		if p.Anio != "" {
			add(n + a + p.Anio)                          // carlosgomez1990
			add(nc + ac + p.Anio)                        // CarlosGomez1990
			add(transforms.TruncateLeft(n, 1) + a + p.Anio) // cgomez1990
		}
	}

	// 9. Contraseñas antiguas con mutaciones
	for _, oldPass := range []string{p.OldPass1, p.OldPass2, p.OldPass3} {
		if oldPass == "" {
			continue
		}
		for _, v := range transforms.AllVariants(oldPass) {
			add(v)
		}
		if p.Anio != "" {
			add(oldPass + p.Anio)
			add(oldPass + p.AnioCorto)
			add(p.Anio + oldPass)
		}
		for _, suf := range transforms.CommonSuffixes {
			add(oldPass + suf)
		}
	}

	// 10. DNI con variantes
	if p.DNI != "" {
		add(p.DNI)
		if p.Nombre != "" {
			n := strings.ToLower(p.Nombre)
			add(n + p.DNI)
			add(transforms.Capitalize(n) + p.DNI)
			add(p.DNI + n)
		}
	}

	// 11. Apodos + todas sus variantes y combinaciones
	if p.Nombre != "" {
		for _, nickVariant := range NicknameVariants(p.Nombre) {
			add(nickVariant)
		}
		for _, nick := range GetNicknames(p.Nombre) {
			if p.Anio != "" {
				add(nick + p.Anio)
				add(nick + p.AnioCorto)
				add(transforms.Capitalize(nick) + p.Anio)
				add(p.Anio + nick)
				add(nick + p.Anio + p.Anio)              // apodo + año repetido
				add(p.Anio + nick + p.Anio)              // sándwich con apodo
			}
			if p.Apellido != "" {
				add(nick + strings.ToLower(p.Apellido))
				add(transforms.Capitalize(nick) + transforms.Capitalize(p.Apellido))
			}
			for _, suf := range []string{"!", "123", "1234", "@", "."} {
				add(nick + suf)
				add(transforms.Capitalize(nick) + suf)
			}
		}
	}

	// ---- Guardar resultado ----
	fmt.Printf("\n\033[32m[+] Total generado: %d palabras\033[0m\n", len(result))

	outputPath := utils.AskStringRequired("Ruta de salida (ej: /home/user/perfil.txt)")

	if err := output.WriteWordlist(result, outputPath); err != nil {
		utils.Error("Error al guardar: " + err.Error())
		return
	}

	output.PrintStats(outputPath, len(result))
}

// collectTokens recolecta todos los campos del perfil como tokens base
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
