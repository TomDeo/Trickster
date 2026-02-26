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

// ================================
// PATRONES CONOCIDOS DE CONTRASEÑAS
// Basados en análisis de RockYou, HaveIBeenPwned y otros leaks reales
// ================================

// sufijos numéricos más comunes encontrados en leaks (orden por frecuencia real)
var commonNumericSuffixes = []string{
	"1", "12", "123", "1234", "12345", "123456",
	"2", "21", "0", "01", "007", "69", "99", "00",
	"11", "22", "33", "44", "55", "66", "77", "88",
	"111", "222", "333", "777", "999", "000",
	"1111", "2222", "3333", "1212", "2121",
	"12345678", "123123", "321", "4321",
}

// prefijos numéricos comunes antes de la palabra
var commonNumericPrefixes = []string{
	"1", "12", "123", "0", "01", "007", "69", "99",
}

// sufijos de símbolos más comunes en contraseñas filtradas
var commonSymbolSuffixes = []string{
	"!", "!!", "!!!", ".", "*", "@", "#", "$", "%", "&",
	"!1", "!12", "!123", "@1", "#1", "!@", "!@#", "!@#$",
	".", "..", "...",
}

// años comunes usados en contraseñas (nacimientos + años recientes)
var commonYears = []string{
	"1970", "1971", "1972", "1973", "1974", "1975",
	"1976", "1977", "1978", "1979", "1980", "1981",
	"1982", "1983", "1984", "1985", "1986", "1987",
	"1988", "1989", "1990", "1991", "1992", "1993",
	"1994", "1995", "1996", "1997", "1998", "1999",
	"2000", "2001", "2002", "2003", "2004", "2005",
	"2020", "2021", "2022", "2023", "2024", "2025",
}

var commonYearsShort = []string{
	"70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99",
	"00", "01", "02", "03", "04", "05",
	"20", "21", "22", "23", "24", "25",
}

// palabras clave comunes que la gente agrega a su nombre/apellido
var commonKeywords = []string{
	"pass", "password", "passwd", "clave", "key",
	"love", "amor", "mi", "el", "la",
	"admin", "root", "user",
	"web", "net", "mail",
	"junior", "senior", "jr",
}

// patrones de teclado muy frecuentes en contraseñas filtradas
var keyboardPatterns = []string{
	"123", "1234", "12345", "123456", "1234567", "12345678",
	"qwerty", "qwerty123", "asdf", "asdfgh", "zxcvbn",
	"abc", "abcd", "abcde", "abc123",
	"111", "1111", "11111",
	"000", "0000", "00000",
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

	result := GenerateFromProfile(p)

	fmt.Printf("\n\033[32m[+] Total generado: %d palabras\033[0m\n", len(result))

	outputPath := utils.AskStringRequired("Ruta de salida (ej: /home/user/perfil.txt)")

	if err := output.WriteWordlist(result, outputPath); err != nil {
		utils.Error("Error al guardar: " + err.Error())
		return
	}

	output.PrintStats(outputPath, len(result))
}

// GenerateFromProfile contiene toda la lógica de generación separada del I/O.
// Esto permite testear la generación sin simular input de usuario.
func GenerateFromProfile(p Profile) []string {
	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		s = strings.TrimSpace(s)
		if len(s) < 3 || len(s) > 32 {
			return // filtrar passwords demasiado cortas o largas
		}
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	tokens := collectTokens(p)

	// ============================================================
	// BLOQUE 1: Variantes individuales de cada token
	// ============================================================
	for _, token := range tokens {
		for _, v := range transforms.AllVariants(token) {
			add(v)
		}
	}

	// ============================================================
	// BLOQUE 2: Combinaciones de 2 tokens (nombre+apellido, etc.)
	// ============================================================
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

	// ============================================================
	// BLOQUE 3: Nombre repetido
	// ============================================================
	if p.Nombre != "" {
		n := strings.ToLower(p.Nombre)
		nc := transforms.Capitalize(n)
		add(n + n)
		add(nc + n)
		add(n + nc)
		add(nc + nc)
	}

	// ============================================================
	// BLOQUE 4: Año repetido y variantes
	// ============================================================
	if p.Anio != "" {
		add(p.Anio + p.Anio)
		add(p.AnioCorto + p.AnioCorto)
		add(p.Anio + p.AnioCorto)
		add(p.AnioCorto + p.Anio)
		for _, suf := range transforms.CommonSuffixes {
			add(p.Anio + suf)
			add(p.AnioCorto + suf)
		}
	}

	// ============================================================
	// BLOQUE 5: Fecha completa en múltiples formatos
	// ============================================================
	if p.Dia != "" && p.Mes != "" && p.Anio != "" {
		add(p.Dia + p.Mes + p.Anio)
		add(p.Anio + p.Mes + p.Dia)
		add(p.Dia + "/" + p.Mes + "/" + p.Anio)
		add(p.Dia + "-" + p.Mes + "-" + p.Anio)
		add(p.Dia + "." + p.Mes + "." + p.Anio)
		add(p.Dia + p.Mes + p.AnioCorto)
		add(p.Anio + p.Dia + p.Mes)
		add(p.Dia + p.Mes)
		add(p.Mes + p.Dia)
		add(p.Mes + p.Anio)
		add(p.Dia + p.Anio)
		// Formato ISO y americano
		add(p.Anio + "-" + p.Mes + "-" + p.Dia)
		add(p.Mes + "/" + p.Dia + "/" + p.Anio)
		for _, suf := range transforms.CommonSuffixes {
			add(p.Dia + p.Mes + p.Anio + suf)
			add(p.Dia + p.Mes + p.AnioCorto + suf)
		}
	}

	// ============================================================
	// BLOQUE 6: Nombre + año (patrones más comunes en leaks)
	// ============================================================
	if p.Nombre != "" && p.Anio != "" {
		n := strings.ToLower(p.Nombre)
		nc := transforms.Capitalize(n)
		nu := transforms.ToUpper(n)

		add(n + p.Anio)
		add(nc + p.Anio)
		add(nu + p.Anio)
		add(p.Anio + n)
		add(p.Anio + nc)
		add(n + p.AnioCorto)
		add(nc + p.AnioCorto)
		add(p.AnioCorto + n)
		add(p.Anio + nu)

		// Año repetido con nombre
		add(n + p.Anio + p.Anio)
		add(n + p.AnioCorto + p.AnioCorto)
		add(nc + p.Anio + p.Anio)

		// Sándwich año-nombre-año
		add(p.Anio + n + p.Anio)
		add(p.AnioCorto + n + p.AnioCorto)
		add(p.Anio + nc + p.Anio)

		// nombre + año + nombre
		add(n + p.Anio + n)
		add(nc + p.Anio + n)

		// Separadores
		for _, sep := range []string{".", "_", "-", "@"} {
			add(n + sep + p.Anio)
			add(nc + sep + p.Anio)
			add(n + sep + p.AnioCorto)
			add(nc + sep + p.AnioCorto)
			add(p.Anio + sep + n)
			add(p.Anio + sep + nc)
		}

		// Nombre + año + símbolo
		for _, suf := range commonSymbolSuffixes {
			add(n + p.Anio + suf)
			add(nc + p.Anio + suf)
			add(n + p.AnioCorto + suf)
			add(nc + p.AnioCorto + suf)
		}

		// Nombre + sufijos numéricos comunes (además del año real)
		for _, num := range commonNumericSuffixes {
			add(n + num)
			add(nc + num)
		}

		// Prefijos numéricos + nombre
		for _, pre := range commonNumericPrefixes {
			add(pre + n)
			add(pre + nc)
		}
	}

	// ============================================================
	// BLOQUE 7: Nombre + fecha completa
	// ============================================================
	if p.Nombre != "" && p.Dia != "" && p.Mes != "" && p.Anio != "" {
		n := strings.ToLower(p.Nombre)
		nc := transforms.Capitalize(n)
		add(n + p.Dia + p.Mes + p.Anio)
		add(nc + p.Dia + p.Mes + p.Anio)
		add(p.Dia + p.Mes + p.Anio + n)
		add(n + p.Dia + p.Mes)
		add(n + p.Mes + p.Anio)
		add(nc + p.Dia + p.Mes)
		add(nc + p.Mes + p.Anio)
		// Formatos inversos
		add(n + p.Anio + p.Mes + p.Dia)
		add(nc + p.Anio + p.Mes + p.Dia)
	}

	// ============================================================
	// BLOQUE 8: Apellido + año (patrón muy frecuente, faltaba antes)
	// ============================================================
	if p.Apellido != "" && p.Anio != "" {
		a := strings.ToLower(p.Apellido)
		ac := transforms.Capitalize(a)
		au := transforms.ToUpper(a)

		add(a + p.Anio)
		add(ac + p.Anio)
		add(au + p.Anio)
		add(p.Anio + a)
		add(a + p.AnioCorto)
		add(ac + p.AnioCorto)
		add(p.AnioCorto + a)

		for _, sep := range []string{".", "_", "-", "@"} {
			add(a + sep + p.Anio)
			add(ac + sep + p.Anio)
			add(a + sep + p.AnioCorto)
		}
		for _, suf := range commonSymbolSuffixes {
			add(a + p.Anio + suf)
			add(ac + p.Anio + suf)
			add(a + p.AnioCorto + suf)
		}
		for _, num := range commonNumericSuffixes {
			add(a + num)
			add(ac + num)
		}
	}

	// ============================================================
	// BLOQUE 9: Nombre + apellido y combinaciones
	// ============================================================
	if p.Nombre != "" && p.Apellido != "" {
		n := strings.ToLower(p.Nombre)
		a := strings.ToLower(p.Apellido)
		nc := transforms.Capitalize(n)
		ac := transforms.Capitalize(a)

		add(n + a)
		add(nc + ac)
		add(nc + a)
		add(n + "_" + a)
		add(n + "." + a)
		add(n + "-" + a)
		add(transforms.TruncateLeft(n, 1) + a)
		add(transforms.TruncateLeft(n, 1) + "_" + a)
		add(transforms.TruncateLeft(n, 1) + "." + a)
		add(a + n)
		add(ac + nc)
		add(a + "_" + n)
		add(a + "." + n)

		// Iniciales
		ini := transforms.TruncateLeft(n, 1) + transforms.TruncateLeft(a, 1)
		add(ini)

		if p.Anio != "" {
			add(n + a + p.Anio)
			add(nc + ac + p.Anio)
			add(transforms.TruncateLeft(n, 1) + a + p.Anio)
			add(a + n + p.Anio)
			add(ac + nc + p.Anio)
			add(n + a + p.AnioCorto)
			add(nc + ac + p.AnioCorto)
			// Con separadores
			for _, sep := range []string{".", "_", "-"} {
				add(n + sep + a + sep + p.Anio)
				add(nc + sep + ac + sep + p.Anio)
			}
		}

		if p.Dia != "" && p.Mes != "" && p.Anio != "" {
			add(n + a + p.Dia + p.Mes + p.Anio)
			add(nc + ac + p.Dia + p.Mes + p.Anio)
		}

		// Nombre+apellido con sufijos numéricos comunes
		for _, num := range commonNumericSuffixes {
			add(n + a + num)
			add(nc + ac + num)
		}
		for _, suf := range commonSymbolSuffixes {
			add(n + a + suf)
			add(nc + ac + suf)
		}
	}

	// ============================================================
	// BLOQUE 10: Equipo de fútbol (muy común en Latam)
	// ============================================================
	if p.EquipoFutbol != "" {
		eq := strings.ToLower(p.EquipoFutbol)
		eqc := transforms.Capitalize(eq)

		for _, v := range transforms.AllVariants(eq) {
			add(v)
		}
		if p.Anio != "" {
			add(eq + p.Anio)
			add(eqc + p.Anio)
			add(eq + p.AnioCorto)
			add(eqc + p.AnioCorto)
			add(p.Anio + eq)
		}
		if p.Nombre != "" {
			n := strings.ToLower(p.Nombre)
			add(n + eq)
			add(eq + n)
			add(transforms.Capitalize(n) + eqc)
		}
		for _, num := range commonNumericSuffixes {
			add(eq + num)
			add(eqc + num)
		}
		for _, suf := range commonSymbolSuffixes {
			add(eq + suf)
			add(eqc + suf)
		}
	}

	// ============================================================
	// BLOQUE 11: Ciudad
	// ============================================================
	if p.Ciudad != "" {
		c := strings.ToLower(p.Ciudad)
		cc := transforms.Capitalize(c)

		for _, v := range transforms.AllVariants(c) {
			add(v)
		}
		if p.Anio != "" {
			add(c + p.Anio)
			add(cc + p.Anio)
			add(c + p.AnioCorto)
		}
		if p.Nombre != "" {
			n := strings.ToLower(p.Nombre)
			add(n + c)
			add(c + n)
		}
		for _, num := range commonNumericSuffixes {
			add(c + num)
			add(cc + num)
		}
	}

	// ============================================================
	// BLOQUE 12: Nombre + año de otros años comunes (no solo el de nacimiento)
	// Cubre el caso donde la persona usa su nombre + año actual o reciente
	// ============================================================
	if p.Nombre != "" {
		n := strings.ToLower(p.Nombre)
		nc := transforms.Capitalize(n)
		for _, yr := range commonYears {
			add(n + yr)
			add(nc + yr)
		}
		for _, yr := range commonYearsShort {
			add(n + yr)
			add(nc + yr)
		}
	}

	// ============================================================
	// BLOQUE 13: Apellido + año de años comunes
	// ============================================================
	if p.Apellido != "" {
		a := strings.ToLower(p.Apellido)
		ac := transforms.Capitalize(a)
		for _, yr := range commonYears {
			add(a + yr)
			add(ac + yr)
		}
		for _, yr := range commonYearsShort {
			add(a + yr)
			add(ac + yr)
		}
	}

	// ============================================================
	// BLOQUE 14: Contraseñas antiguas con mutaciones profundas
	// ============================================================
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
			add(p.AnioCorto + oldPass)
		}
		for _, suf := range transforms.CommonSuffixes {
			add(oldPass + suf)
		}
		for _, suf := range commonSymbolSuffixes {
			add(oldPass + suf)
		}
		for _, num := range commonNumericSuffixes {
			add(oldPass + num)
		}
		// Patrones de mutación incremental (pass1 → pass2, pass! → pass!!)
		add(oldPass + "1")
		add(oldPass + "2")
		add(oldPass + "!")
		add(oldPass + "!!")
		add("1" + oldPass)
		if p.Nombre != "" {
			n := strings.ToLower(p.Nombre)
			add(oldPass + n)
			add(n + oldPass)
		}
	}

	// ============================================================
	// BLOQUE 15: DNI con variantes
	// ============================================================
	if p.DNI != "" {
		add(p.DNI)
		if p.Nombre != "" {
			n := strings.ToLower(p.Nombre)
			add(n + p.DNI)
			add(transforms.Capitalize(n) + p.DNI)
			add(p.DNI + n)
		}
		if p.Apellido != "" {
			a := strings.ToLower(p.Apellido)
			add(a + p.DNI)
			add(p.DNI + a)
		}
		for _, suf := range commonSymbolSuffixes {
			add(p.DNI + suf)
		}
	}

	// ============================================================
	// BLOQUE 16: Apodos completos con todas sus variantes
	// ============================================================
	if p.Nombre != "" {
		for _, nickVariant := range NicknameVariants(p.Nombre) {
			add(nickVariant)
		}
		for _, nick := range GetNicknames(p.Nombre) {
			if p.Anio != "" {
				add(nick + p.Anio)
				add(nick + p.AnioCorto)
				add(transforms.Capitalize(nick) + p.Anio)
				add(transforms.Capitalize(nick) + p.AnioCorto)
				add(p.Anio + nick)
				add(p.AnioCorto + nick)
				add(nick + p.Anio + p.Anio)
				add(p.Anio + nick + p.Anio)
			}
			if p.Apellido != "" {
				a := strings.ToLower(p.Apellido)
				add(nick + a)
				add(transforms.Capitalize(nick) + transforms.Capitalize(a))
				add(nick + "_" + a)
				add(nick + "." + a)
			}
			for _, suf := range commonSymbolSuffixes {
				add(nick + suf)
				add(transforms.Capitalize(nick) + suf)
			}
			for _, num := range commonNumericSuffixes {
				add(nick + num)
				add(transforms.Capitalize(nick) + num)
			}
			// Apodo + años comunes (no solo el de nacimiento)
			for _, yr := range commonYears {
				add(nick + yr)
				add(transforms.Capitalize(nick) + yr)
			}
		}
	}

	// ============================================================
	// BLOQUE 17: Nombre + palabras clave comunes
	// Patrón: carlos_admin, carlospass, micarlos, etc.
	// ============================================================
	if p.Nombre != "" {
		n := strings.ToLower(p.Nombre)
		nc := transforms.Capitalize(n)
		for _, kw := range commonKeywords {
			add(n + kw)
			add(kw + n)
			add(nc + kw)
			add(kw + nc)
			add(n + "_" + kw)
			add(kw + "_" + n)
		}
	}

	// ============================================================
	// BLOQUE 18: Patrones de teclado populares combinados con nombre
	// ============================================================
	if p.Nombre != "" {
		n := strings.ToLower(p.Nombre)
		nc := transforms.Capitalize(n)
		for _, kp := range keyboardPatterns {
			add(n + kp)
			add(nc + kp)
			add(kp + n)
		}
	}

	// ============================================================
	// BLOQUE 19: Edad con variantes
	// ============================================================
	if p.Edad != "" {
		e := p.Edad
		if p.Nombre != "" {
			n := strings.ToLower(p.Nombre)
			add(n + e)
			add(transforms.Capitalize(n) + e)
			add(e + n)
		}
		if p.Apellido != "" {
			a := strings.ToLower(p.Apellido)
			add(a + e)
			add(e + a)
		}
		for _, suf := range commonSymbolSuffixes {
			add(e + suf)
		}
	}

	// ============================================================
	// BLOQUE 20: Solo nombre/apellido con sufijos numéricos comunes
	// (Cubre cuando el año real no es conocido)
	// ============================================================
	bases := []string{}
	if p.Nombre != "" {
		bases = append(bases, strings.ToLower(p.Nombre))
	}
	if p.Apellido != "" {
		bases = append(bases, strings.ToLower(p.Apellido))
	}
	for _, base := range bases {
		bc := transforms.Capitalize(base)
		for _, suf := range commonSymbolSuffixes {
			add(base + suf)
			add(bc + suf)
		}
		// Duplicados: carloscarlos, gomezgomez
		add(base + base)
		add(bc + bc)
	}

	return result
}

// collectTokens recolecta todos los campos del perfil como tokens base.
// Filtra tokens vacíos, duplicados y demasiado cortos.
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
		if f != "" && len([]rune(f)) >= 2 && !seen[f] {
			seen[f] = true
			tokens = append(tokens, f)
		}
	}

	return tokens
}
