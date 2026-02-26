package core

import (
	"strconv"
	"strings"
	"trickster/transforms"
)

// ================================================================
// MÓDULO: PATRONES LOCALES ARGENTINA
//
// Complementa profiler.go con patrones específicos del contexto
// argentino/rioplatense no cubiertos por las listas genéricas de CUPP.
//
// Fuentes:
//   - Análisis de leaks de foros argentinos (foros.biz, etc.)
//   - Comportamiento típico en uso de redes sociales locales
//   - Vocabulario y jerga argentina en contraseñas
// ================================================================

// argSuffixPatterns: patrones de repetición locales muy comunes.
// En Argentina es frecuente duplicar años o agregar el mismo año repetido.
// Ej: "pedro20242024", "perla2024.2024", "carla9090"
//
// Esta función los genera para un token + año dado.
func ArgRepeatPatterns(token, anio, anioCorto string) []string {
	if token == "" {
		return nil
	}

	var result []string
	seen := make(map[string]bool)

	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	// Patrones de año duplicado (muy comunes en Argentina según análisis de leaks locales)
	if anio != "" {
		add(token + anio + anio)             // juan19901990
		add(token + anio + "." + anio)       // juan1990.1990
		add(token + anio + "_" + anio)       // juan1990_1990
	}
	if anioCorto != "" {
		add(token + anioCorto + anioCorto)         // juan9090
		add(token + anioCorto + "." + anioCorto)   // juan90.90
	}
	if anio != "" && anioCorto != "" {
		add(token + anio + anioCorto)        // juan199090
		add(token + anioCorto + anio)        // juan901990
	}

	// Patrón "token + año actual repetido" (muy visto en 2023-2025)
	for _, yr := range []string{"2023", "2024", "2025"} {
		add(token + yr + yr)
		add(token + yr[2:] + yr[2:]) // token + 2324
	}

	return result
}

// ArgLeetLocal: leet speak extendido con variaciones rioplatenses.
// Más allá del leet estándar, en Argentina se ven estos patrones:
//   - "q" en lugar de "k" (qiero, kiero)
//   - "x" en lugar de "ch" o "ks" (mucho → mucho, extra → extra)
//   - Uso de "ph" → no tan común, pero sí:
//   - Duplicación de letras al final: "vickyyy", "carlosss"
//   - Puntuación al final triplicada: "vicky!!!", "carlos..."
func ArgLetterDuplication(token string) []string {
	if token == "" {
		return nil
	}

	runes := []rune(token)
	last := string(runes[len(runes)-1])

	var result []string
	seen := make(map[string]bool)
	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	// Duplicar/triplicar la última letra
	add(token + last)            // carross
	add(token + last + last)     // carloss

	// Sufijos de puntuación triplicados (muy comunes en contraseñas arg)
	add(token + "...")
	add(token + "!!!")
	add(token + "???")

	// "q" reemplazando "c" o "k" al inicio (kiero → qiero)
	if len(runes) > 0 && (runes[0] == 'c' || runes[0] == 'k') {
		withQ := "q" + string(runes[1:])
		add(withQ)
	}

	return result
}

// ArgCommonPhrases: frases/palabras muy usadas en contraseñas argentinas.
// Detectadas en análisis de leaks de servicios con base en Argentina.
var ArgCommonPhrases = []string{
	// Expresiones locales
	"pelotudo", "boludo", "chabón", "chabon", "pibe", "mina",
	"groso", "grosso", "capo", "fenomeno",
	"fenomenal", "genial",

	// Términos afectivos muy usados en contraseñas de parejas
	"amor", "mi amor", "miamor", "corazon", "mi vida", "mivida",
	"bebe", "bb", "baby", "nena", "nene", "gordo", "gorda",
	"flaco", "flaca", "cielo", "vita",

	// Jerga de gaming/internet argentina
	"crack", "master", "pro", "noob",
	"negro", "negra", "blanquito",

	// Palabras del fútbol (muy presentes en contraseñas arg)
	"boca", "river", "racing", "independiente", "sanlorenzo",
	"huracan", "velez", "lanus", "belgrano", "talleres",
	"estudiantes", "gimnasia", "newells", "rosariocentral",
	"bocajuniors", "riverplate",

	// Números de camiseta típicos
	"numero10", "numero9", "eldiez", "elnueve",

	// Frases de uso común
	"password", "contrasena", "contraseña", "miclave", "clave",
	"micon", "miconta", "micuenta",

	// Sufijos "argentinos" de complejidad (para políticas de contraseña)
	// La gente resuelve el "necesitás mayúscula+número+símbolo" así:
	// NombreCapitalized + año + !
}

// GenerateArgPatterns genera candidatos específicos del contexto argentino.
// Se llama desde RunProfiler después de GenerateFromProfile para agregar
// el vocabulario local sin duplicar la lógica base.
func GenerateArgPatterns(p Profile) []string {
	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		s = trimAndCheck(s)
		if s == "" || seen[s] {
			return
		}
		seen[s] = true
		result = append(result, s)
	}

	n := ""
	nc := ""
	if p.Nombre != "" {
		n = lowerTrim(p.Nombre)
		nc = capFirst(n)
	}

	// ── 1. Frases locales combinadas con el nombre ────────────────
	for _, phrase := range ArgCommonPhrases {
		add(phrase)
		if n != "" {
			add(n + phrase)
			add(nc + phrase)
			add(phrase + n)
			add(phrase + nc)
		}
		if p.Anio != "" {
			add(phrase + p.Anio)
			add(phrase + p.AnioCorto)
		}
	}

	// ── 2. Clubes de fútbol de Argentina × año ────────────────────
	clubes := []string{
		"boca", "bocajuniors", "river", "riverplate",
		"racing", "independiente", "sanlorenzo",
		"huracan", "velez", "lanus", "belgrano",
		"talleres", "estudiantes", "gimnasia",
		"newells", "rosariocentral", "banfield",
		"platense", "sarmiento", "tigre",
	}

	// Si el objetivo ingresó equipo, ya está cubierto en profiler.go.
	// Acá cubrimos los más frecuentes sin importar el equipo declarado.
	for _, club := range clubes {
		add(club)
		add(capFirst(club))
		if p.Anio != "" {
			add(club + p.Anio)
			add(capFirst(club) + p.Anio)
			add(club + p.AnioCorto)
		}
		if n != "" {
			add(n + club)
			add(club + n)
			add(nc + capFirst(club))
		}
		for _, num := range []string{"1", "10", "9", "11", "123"} {
			add(club + num)
		}
	}

	// ── 3. Patrones de año duplicado (comportamiento local) ───────
	if n != "" && p.Anio != "" {
		for _, pat := range ArgRepeatPatterns(n, p.Anio, p.AnioCorto) {
			add(pat)
		}
		for _, pat := range ArgRepeatPatterns(nc, p.Anio, p.AnioCorto) {
			add(pat)
		}
	}
	if p.Apellido != "" {
		aLow := lowerTrim(p.Apellido)
		aCap := capFirst(aLow)
		if p.Anio != "" {
			for _, pat := range ArgRepeatPatterns(aLow, p.Anio, p.AnioCorto) {
				add(pat)
			}
			for _, pat := range ArgRepeatPatterns(aCap, p.Anio, p.AnioCorto) {
				add(pat)
			}
		}
	}

	// ── 4. Duplicación de letras y puntuación local ───────────────
	for _, base := range []string{n, nc} {
		if base == "" {
			continue
		}
		for _, pat := range ArgLetterDuplication(base) {
			add(pat)
		}
	}
	if p.Apellido != "" {
		for _, base := range []string{lowerTrim(p.Apellido), capFirst(lowerTrim(p.Apellido))} {
			for _, pat := range ArgLetterDuplication(base) {
				add(pat)
			}
		}
	}

	// ── 5. CUIL: relacionado con el DNI, muy usado como contraseña ─
	// Formato CUIL: 20-XXXXXXXD-N (el DNI va en el medio)
	// La gente a veces usa su CUIL completo o parcial como contraseña.
	if p.DNI != "" {
		dniClean := p.DNI
		for _, r := range []string{".", "-", " "} {
			dniClean = strings.ReplaceAll(dniClean, r, "")
		}
		// Prefijos CUIL más comunes para ciudadanos argentinos
		for _, prefix := range []string{"20", "23", "24", "27"} {
			// No conocemos el dígito verificador, generamos los posibles (0-9)
			for d := 0; d <= 9; d++ {
				cuil := prefix + dniClean + strconv.Itoa(d)
				add(cuil)
				add(prefix + "-" + dniClean + "-" + strconv.Itoa(d))
			}
		}
	}

	return result
}

// ── helpers locales ───────────────────────────────────────────────

func lowerTrim(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func capFirst(s string) string {
	return transforms.Capitalize(s)
}

func trimAndCheck(s string) string {
	s = strings.TrimSpace(s)
	l := len([]rune(s))
	if l < 4 || l > 28 {
		return ""
	}
	return s
}
