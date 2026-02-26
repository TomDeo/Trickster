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
	Mascota         string
	Pareja          string
	OldPass1        string
	OldPass2        string
	OldPass3        string
}

// ================================================================
// TABLAS DE MUTACIÓN
// Basadas en análisis de RockYou, HIBP, CUPP config y hashcat best64/d3ad0ne
// ================================================================

// numSuffixes: sufijos numéricos ordenados por frecuencia real en leaks.
// CUPP usa 0-100 por defecto y cubre el ~80% de los casos numéricos.
// Añadimos años 4 dígitos y patrones de teclado frecuentes.
var numSuffixes = buildNumSuffixes()

func buildNumSuffixes() []string {
	seen := make(map[string]bool)
	var s []string
	addUniq := func(v string) {
		if !seen[v] {
			seen[v] = true
			s = append(s, v)
		}
	}

	// 0–100 completo (CUPP default)
	for i := 0; i <= 100; i++ {
		addUniq(fmt.Sprintf("%d", i))
	}
	// Años de 2 dígitos
	for i := 60; i <= 99; i++ {
		addUniq(fmt.Sprintf("%d", i))
	}
	// Años de 4 dígitos
	for y := 1960; y <= 2025; y++ {
		addUniq(fmt.Sprintf("%d", y))
	}
	// Patrones de teclado numérico
	for _, p := range []string{
		"123", "1234", "12345", "123456", "1234567", "12345678", "123456789",
		"111", "222", "333", "444", "555", "666", "777", "888", "999",
		"000", "1111", "2222", "3333", "4444", "5555",
		"1212", "2121", "3131", "1122", "2211",
		"321", "4321", "54321", "112", "121", "211",
		"007", "069", "420", "666", "101", "404",
	} {
		addUniq(p)
	}
	return s
}

// specialSuffixes: sufijos de símbolos más frecuentes en leaks reales
// Fuente: CUPP specialchars config + hashcat best64 analysis + NotSoSecure study
var specialSuffixes = []string{
	"!", "!!", "!!!", ".", "..", "...",
	"@", "#", "$", "%", "&", "*",
	"!1", "!12", "!123", "!1234",
	"!@", "!@#", "!@#$", "!@#$%",
	"@1", "@123", "@1234",
	"#1", "#123",
	"$1", "$123",
	"1!", "1!!", "12!", "123!", "1234!",
	".*", "*", "**",
	"_1", "_12", "_123",
	"-1", "-12", "-123",
}

// numSymbolSuffixes: combinaciones número+símbolo muy comunes en políticas
// de contraseñas que exigen complejidad (ej: Carlos1990!)
var numSymbolSuffixes = []string{
	"1!", "1!!", "2!", "12!", "123!", "1234!",
	"1@", "12@", "123@",
	"1#", "12#", "123#",
	"01!", "0!", "00!",
	"1.", "12.", "123.",
}

// numPrefixes: prefijos numéricos comunes antes de la palabra
var numPrefixes = []string{
	"1", "12", "123", "0", "00", "01", "007",
	"69", "99", "11", "22", "33",
}

// specialPrefixes: prefijos de símbolos comunes
var specialPrefixes = []string{"!", "@", "#", "$"}

// leetTable: tabla de sustituciones leet.
// Fuente: CUPP config (a=4,i=1,e=3,t=7,o=0,s=5,g=9,z=2) + extensiones comunes
var leetTable = map[rune][]string{
	'a': {"4", "@"},
	'e': {"3"},
	'i': {"1", "!"},
	'o': {"0"},
	's': {"5", "$"},
	't': {"7"},
	'g': {"9"},
	'z': {"2"},
	'l': {"1"},
	'b': {"8"},
}

// passwordKeywords: palabras que la gente mezcla con su nombre en contraseñas
// Fuente: estudios de HIBP + análisis de behavior de usuarios
var passwordKeywords = []string{
	"pass", "password", "passwd", "clave", "key", "secret",
	"amor", "love", "mi", "baby", "bebe", "bb",
	"admin", "root", "user", "web", "mail", "net",
	"123", "1234", "12345",
	"forever", "always", "lucky",
}

// standaloneKeyboard: patrones de teclado autónomos muy frecuentes en leaks
var standaloneKeyboard = []string{
	"qwerty", "qwerty123", "qwertyuiop",
	"asdf", "asdfgh", "asdfghjkl",
	"zxcvbn",
	"abc123", "abc", "abcd",
	"password", "pass", "passwd",
	"admin", "root", "user", "login",
	"letmein", "welcome",
	"iloveyou",
	"dragon", "monkey", "shadow",
	"sunshine", "princess", "football",
	"superman", "batman",
}

// ================================================================
// FUNCIÓN PRINCIPAL
// ================================================================

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
	p.Mascota = utils.AskOptional("Nombre de mascota")
	p.Pareja = utils.AskOptional("Nombre de pareja / familiar cercano")
	p.Edad = utils.AskOptional("Edad actual")
	p.Ciudad = utils.AskOptional("Ciudad")
	p.OldPass1 = utils.AskOptional("Contraseña antigua 1")
	p.OldPass2 = utils.AskOptional("Contraseña antigua 2")
	p.OldPass3 = utils.AskOptional("Contraseña antigua 3")

	fmt.Println()
	utils.Info("Procesando perfil y generando wordlist...")

	result := GenerateFromProfile(p)

	fmt.Printf("\n\033[32m[+] Total generado: %d palabras\033[0m\n", len(result))

	outputPath := utils.AskStringRequired("Ruta de salida (ej: /home/user/perfil.txt)")

	if err := output.WriteWordlist(result, outputPath); err != nil {
		utils.Error("Error al guardar: " + err.Error())
		return
	}

	output.PrintStats(outputPath, len(result))
}

// ================================================================
// GENERADOR PRINCIPAL — separado del I/O para facilitar testing
// ================================================================

func GenerateFromProfile(p Profile) []string {
	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		s = strings.TrimSpace(s)
		l := len([]rune(s))
		if l < 4 || l > 28 {
			return
		}
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	// ── PASO 1: Construir átomos base (tokens personales) ─────────
	atoms := buildAtoms(p)

	// ── PASO 2: Expandir cada átomo textual a sus formas de casing ─
	type expandedWord struct {
		lower string
		cap   string
		upper string
		leet  string // leet simple de lower
		leetC string // leet de Cap
	}

	var textForms []expandedWord
	var numForms []string // átomos numéricos van directo

	for _, a := range atoms {
		if a.isNumber {
			numForms = append(numForms, a.val)
			add(a.val)
			continue
		}
		e := expandedWord{
			lower: a.val,
			cap:   transforms.Capitalize(a.val),
			upper: transforms.ToUpper(a.val),
			leet:  leetSimple(a.val),
			leetC: leetSimple(transforms.Capitalize(a.val)),
		}
		textForms = append(textForms, e)

		// Formas base directas
		add(e.lower)
		add(e.cap)
		add(e.upper)
		add(e.leet)
		add(e.leetC)
		add(transforms.Reverse(e.lower))
		add(transforms.Reverse(e.cap))

		// Leet con todas las variantes (solo para tokens cortos ≤8 chars)
		for _, lv := range leetAllVariants(e.lower) {
			add(lv)
			add(transforms.Capitalize(lv))
		}
	}

	// ── PASO 3: Núcleo — cada forma × todos los sufijos/prefijos ──
	// Esta es la operación que multiplica de ~12k a >200k candidatos.
	// CUPP aplica 0-100 + años + chars especiales a CADA forma.
	for _, e := range textForms {
		// Sufijos numéricos (0-100 + años + patrones de teclado)
		for _, num := range numSuffixes {
			add(e.lower + num)
			add(e.cap + num)
			if e.leet != e.lower {
				add(e.leet + num)
			}
		}

		// Sufijos de símbolos especiales
		for _, sp := range specialSuffixes {
			add(e.lower + sp)
			add(e.cap + sp)
		}

		// Sufijos número+símbolo (cumple políticas de complejidad)
		for _, ns := range numSymbolSuffixes {
			add(e.lower + ns)
			add(e.cap + ns)
		}

		// Prefijos numéricos
		for _, pre := range numPrefixes {
			add(pre + e.lower)
			add(pre + e.cap)
		}

		// Prefijos de símbolos
		for _, pre := range specialPrefixes {
			add(pre + e.lower)
			add(pre + e.cap)
		}

		// Número + palabra + número (patrón tipo 1carlos1, 123carlos123)
		for _, num := range []string{"1", "12", "123", "0", "01", "00", "007"} {
			add(num + e.lower + num)
			add(num + e.cap + num)
		}

		// Palabra duplicada (carlos → carloscarlos)
		add(e.lower + e.lower)
		add(e.cap + e.lower)
		add(e.cap + e.cap)
	}

	// ── PASO 4: Año real del objetivo × todas las formas ──────────
	// El año personal es el multiplicador más efectivo en contraseñas reales.
	if p.Anio != "" {
		for _, e := range textForms {
			add(e.lower + p.Anio)
			add(e.cap + p.Anio)
			add(e.upper + p.Anio)
			add(e.leet + p.Anio)
			add(e.leetC + p.Anio)
			add(p.Anio + e.lower)
			add(p.Anio + e.cap)
			add(e.lower + p.AnioCorto)
			add(e.cap + p.AnioCorto)
			add(p.AnioCorto + e.lower)
			add(p.AnioCorto + e.cap)

			// forma + año + símbolo
			for _, sp := range specialSuffixes {
				add(e.lower + p.Anio + sp)
				add(e.cap + p.Anio + sp)
				add(e.lower + p.AnioCorto + sp)
				add(e.cap + p.AnioCorto + sp)
			}
			// forma + año + número simple
			for _, num := range []string{"1", "2", "3", "12", "123"} {
				add(e.lower + p.Anio + num)
				add(e.cap + p.Anio + num)
			}

			// Sándwich: año-forma-año
			add(p.Anio + e.lower + p.Anio)
			add(p.AnioCorto + e.lower + p.AnioCorto)
			add(p.Anio + e.cap + p.Anio)

			// forma + año repetido
			add(e.lower + p.Anio + p.Anio)
			add(e.cap + p.Anio + p.Anio)
			add(e.lower + p.AnioCorto + p.AnioCorto)

			// Separadores entre forma y año
			for _, sep := range []string{".", "_", "-", "@"} {
				add(e.lower + sep + p.Anio)
				add(e.cap + sep + p.Anio)
				add(e.lower + sep + p.AnioCorto)
				add(e.cap + sep + p.AnioCorto)
				add(p.Anio + sep + e.lower)
				add(p.Anio + sep + e.cap)
			}
		}

		// Variantes del año solo
		add(p.Anio + p.Anio)
		add(p.AnioCorto + p.AnioCorto)
		add(p.Anio + p.AnioCorto)
		for _, sp := range specialSuffixes {
			add(p.Anio + sp)
			add(p.AnioCorto + sp)
		}
	}

	// ── PASO 5: Combinaciones de 2 formas entre sí ────────────────
	type sf struct{ lower, cap string }
	var simpleForms []sf
	for _, e := range textForms {
		simpleForms = append(simpleForms, sf{e.lower, e.cap})
	}

	for i, a := range simpleForms {
		for j, b := range simpleForms {
			if i == j {
				continue
			}
			add(a.lower + b.lower)
			add(a.cap + b.cap)
			add(a.cap + b.lower)
			add(a.lower + b.cap)

			for _, sep := range []string{".", "_", "-", "@"} {
				add(a.lower + sep + b.lower)
				add(a.cap + sep + b.cap)
				add(a.cap + sep + b.lower)
			}

			if p.Anio != "" {
				add(a.lower + b.lower + p.Anio)
				add(a.cap + b.cap + p.Anio)
				add(a.cap + b.cap + p.AnioCorto)
				for _, sep := range []string{".", "_", "-"} {
					add(a.cap + sep + b.cap + sep + p.Anio)
				}
			}

			// Sufijos comunes sobre la combinación
			for _, sp := range []string{"!", "1", "12", "123", "1234", "@", "#", "1!"} {
				add(a.lower + b.lower + sp)
				add(a.cap + b.cap + sp)
			}
		}
	}

	// ── PASO 6: Fechas en múltiples formatos ─────────────────────
	if p.Dia != "" && p.Mes != "" && p.Anio != "" {
		fechas := []string{
			p.Dia + p.Mes + p.Anio,
			p.Anio + p.Mes + p.Dia,
			p.Dia + p.Mes + p.AnioCorto,
			p.Dia + p.Mes,
			p.Mes + p.Anio,
			p.Mes + p.Dia,
			p.Anio + p.Dia + p.Mes,
			p.Dia + "-" + p.Mes + "-" + p.Anio,
			p.Dia + "/" + p.Mes + "/" + p.Anio,
			p.Dia + "." + p.Mes + "." + p.Anio,
			p.Anio + "-" + p.Mes + "-" + p.Dia,
		}

		for _, fecha := range fechas {
			add(fecha)
			for _, sp := range specialSuffixes {
				add(fecha + sp)
			}
		}

		// Nombre/Apellido + cada formato de fecha
		for _, e := range textForms {
			for _, fecha := range fechas {
				add(e.lower + fecha)
				add(e.cap + fecha)
				add(fecha + e.lower)
				add(fecha + e.cap)
				for _, sp := range []string{"!", "@", "#", "1", "123"} {
					add(e.lower + fecha + sp)
					add(e.cap + fecha + sp)
				}
			}
		}
	}

	// ── PASO 7: Apodos × todos los sufijos ────────────────────────
	if p.Nombre != "" {
		nicks := GetNicknames(p.Nombre)
		for _, nick := range nicks {
			nc := transforms.Capitalize(nick)
			nl := leetSimple(nick)

			add(nick)
			add(nc)
			if nl != nick {
				add(nl)
			}

			// Sufijos completos sobre cada apodo
			for _, num := range numSuffixes {
				add(nick + num)
				add(nc + num)
			}
			for _, sp := range specialSuffixes {
				add(nick + sp)
				add(nc + sp)
			}
			for _, ns := range numSymbolSuffixes {
				add(nick + ns)
				add(nc + ns)
			}

			if p.Anio != "" {
				add(nick + p.Anio)
				add(nc + p.Anio)
				add(nick + p.AnioCorto)
				add(nc + p.AnioCorto)
				add(p.Anio + nick)
				add(p.Anio + nc)
				add(nick + p.Anio + p.Anio)
				add(p.Anio + nick + p.Anio)
				for _, sp := range specialSuffixes {
					add(nick + p.Anio + sp)
					add(nc + p.Anio + sp)
				}
			}

			if p.Apellido != "" {
				a := strings.ToLower(p.Apellido)
				ac := transforms.Capitalize(a)
				add(nick + a)
				add(nc + ac)
				add(nick + "_" + a)
				add(nick + "." + a)
				if p.Anio != "" {
					add(nick + a + p.Anio)
					add(nc + ac + p.Anio)
				}
				for _, sp := range []string{"!", "1", "123", "@"} {
					add(nick + a + sp)
					add(nc + ac + sp)
				}
			}
		}
	}

	// ── PASO 8: Inicial del nombre + apellido ─────────────────────
	if p.Nombre != "" && p.Apellido != "" {
		n := strings.ToLower(p.Nombre)
		a := strings.ToLower(p.Apellido)
		ini := string([]rune(n)[0])

		add(ini + a)
		add(ini + "." + a)
		add(ini + "_" + a)
		add(strings.ToUpper(ini) + transforms.Capitalize(a))

		for _, num := range numSuffixes[:50] {
			add(ini + a + num)
		}
		for _, sp := range specialSuffixes {
			add(ini + a + sp)
		}
		if p.Anio != "" {
			add(ini + a + p.Anio)
			add(ini + a + p.AnioCorto)
		}
	}

	// ── PASO 9: DNI con variantes ─────────────────────────────────
	if p.DNI != "" {
		add(p.DNI)
		for _, sp := range specialSuffixes {
			add(p.DNI + sp)
		}
		for _, e := range textForms {
			add(e.lower + p.DNI)
			add(e.cap + p.DNI)
			add(p.DNI + e.lower)
			add(p.DNI + e.cap)
			for _, sp := range specialSuffixes {
				add(e.lower + p.DNI + sp)
				add(e.cap + p.DNI + sp)
			}
		}
	}

	// ── PASO 10: Contraseñas antiguas con mutación profunda ───────
	// La gente suele mutar su contraseña anterior añadiendo sufijos,
	// cambiando el año o haciendo pequeñas variaciones. Este es el patrón
	// más efectivo cuando se conocen contraseñas previas.
	for _, oldPass := range []string{p.OldPass1, p.OldPass2, p.OldPass3} {
		if oldPass == "" {
			continue
		}
		add(oldPass)
		add(transforms.Capitalize(oldPass))
		add(transforms.ToUpper(oldPass))
		add(leetSimple(oldPass))

		for _, num := range numSuffixes {
			add(oldPass + num)
		}
		for _, sp := range specialSuffixes {
			add(oldPass + sp)
		}
		for _, ns := range numSymbolSuffixes {
			add(oldPass + ns)
		}
		for _, pre := range numPrefixes {
			add(pre + oldPass)
		}

		if p.Anio != "" {
			add(oldPass + p.Anio)
			add(oldPass + p.AnioCorto)
			add(p.Anio + oldPass)
			for _, sp := range specialSuffixes {
				add(oldPass + p.Anio + sp)
			}
		}

		for _, e := range textForms {
			add(oldPass + e.lower)
			add(e.lower + oldPass)
		}
	}

	// ── PASO 11: Palabras clave × bases personales ────────────────
	for _, e := range textForms {
		for _, kw := range passwordKeywords {
			add(e.lower + kw)
			add(kw + e.lower)
			add(e.cap + kw)
			add(kw + e.cap)
			add(e.lower + "_" + kw)
			add(kw + "_" + e.lower)
		}
	}

	// ── PASO 12: Patrones de teclado autónomos ───────────────────
	for _, kp := range standaloneKeyboard {
		add(kp)
		if p.Nombre != "" {
			n := strings.ToLower(p.Nombre)
			nc := transforms.Capitalize(n)
			add(n + kp)
			add(nc + kp)
			add(kp + n)
			add(kp + nc)
		}
	}

	// ── PASO 13: Leet recursivo sobre combinaciones clave ─────────
	// Solo sobre las combinaciones más probables para no explotar
	if p.Nombre != "" && p.Anio != "" {
		for _, lv := range leetAllVariants(strings.ToLower(p.Nombre)) {
			add(lv + p.Anio)
			add(transforms.Capitalize(lv) + p.Anio)
		}
	}
	if p.Nombre != "" && p.Apellido != "" {
		combined := strings.ToLower(p.Nombre) + strings.ToLower(p.Apellido)
		if len([]rune(combined)) <= 10 {
			for _, lv := range leetAllVariants(combined) {
				add(lv)
				add(transforms.Capitalize(lv))
				if p.Anio != "" {
					add(lv + p.Anio)
				}
			}
		}
	}

	return result
}

// ================================================================
// HELPERS INTERNOS
// ================================================================

// leetSimple aplica sustitución leet simple (primera opción por letra).
func leetSimple(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if subs, ok := leetTable[r]; ok {
			b.WriteString(subs[0])
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// leetAllVariants genera todas las combinaciones posibles de sustitución leet.
// Limitado a tokens de ≤10 chars para evitar explosión combinatoria.
func leetAllVariants(s string) []string {
	runes := []rune(strings.ToLower(s))
	if len(runes) > 10 {
		return []string{leetSimple(s)}
	}

	results := []string{""}
	for _, r := range runes {
		subs, ok := leetTable[r]
		if !ok {
			for i := range results {
				results[i] += string(r)
			}
			continue
		}
		options := append([]string{string(r)}, subs...)
		var newResults []string
		for _, prev := range results {
			for _, opt := range options {
				newResults = append(newResults, prev+opt)
			}
		}
		results = newResults
	}

	// Filtrar el original (ya fue agregado como átomo base)
	original := strings.ToLower(s)
	var filtered []string
	for _, v := range results {
		if v != original {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// atom representa un token base con metadatos
type atom struct {
	val      string
	isNumber bool
}

// buildAtoms recolecta todos los campos del perfil como átomos,
// distinguiendo tokens de texto de los numéricos.
func buildAtoms(p Profile) []atom {
	var atoms []atom
	seen := make(map[string]bool)

	add := func(val string, isNum bool) {
		val = strings.ToLower(strings.TrimSpace(val))
		if val == "" || len([]rune(val)) < 2 || seen[val] {
			return
		}
		seen[val] = true
		atoms = append(atoms, atom{val, isNum})
	}

	add(p.Nombre, false)
	add(p.Apellido, false)
	add(p.EquipoFutbol, false)
	add(p.Ciudad, false)
	add(p.Mascota, false)
	add(p.Pareja, false)

	add(p.DNI, true)
	add(p.Anio, true)
	add(p.AnioCorto, true)
	add(p.Dia, true)
	add(p.Mes, true)
	add(p.FechaNacimiento, true)
	add(p.Edad, true)

	return atoms
}

// collectTokens mantiene compatibilidad con el resto del código
func collectTokens(p Profile) []string {
	atoms := buildAtoms(p)
	tokens := make([]string, 0, len(atoms))
	for _, a := range atoms {
		tokens = append(tokens, a.val)
	}
	return tokens
}
