package transforms

import (
	"strings"
	"unicode"
)

// ================================
// TRANSFORMACIONES BÁSICAS
// ================================

// ToUpper convierte toda la palabra a mayúsculas: "carlos" → "CARLOS"
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToLower convierte toda la palabra a minúsculas: "CARLOS" → "carlos"
func ToLower(s string) string {
	return strings.ToLower(s)
}

// Capitalize pone la primera letra en mayúscula: "carlos" → "Carlos"
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// CapitalizeAll pone en mayúscula la primera letra de cada "palabra"
// útil para nombres compuestos: "juan pablo" → "Juan Pablo"
func CapitalizeAll(s string) string {
	return strings.Title(strings.ToLower(s)) //nolint
}

// Reverse invierte el string: "carlos" → "solrac"
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ================================
// LEET SPEAK
// ================================

// leetMap define las sustituciones clásicas de leet speak
var leetMap = map[rune]string{
	'a': "4", 'A': "4",
	'e': "3", 'E': "3",
	'i': "1", 'I': "1",
	'o': "0", 'O': "0",
	's': "5", 'S': "5",
	't': "7", 'T': "7",
	'b': "8", 'B': "8",
	'g': "9", 'G': "9",
}

// Leet aplica sustituciones de leet speak: "carlos" → "c4rl0s"
func Leet(s string) string {
	var result strings.Builder
	for _, r := range s {
		if replacement, ok := leetMap[r]; ok {
			result.WriteString(replacement)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// LeetPartial aplica leet solo a vocales: "carlos" → "c4rl0s"
// (misma que Leet en este caso, pero puede extenderse)
func LeetPartial(s string) string {
	vocales := map[rune]string{
		'a': "4", 'A': "4",
		'e': "3", 'E': "3",
		'i': "1", 'I': "1",
		'o': "0", 'O': "0",
		'u': "v", 'U': "V",
	}
	var result strings.Builder
	for _, r := range s {
		if replacement, ok := vocales[r]; ok {
			result.WriteString(replacement)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ================================
// SUFIJOS / PREFIJOS COMUNES
// ================================

// suffixes lista de sufijos numéricos y símbolos más comunes en contraseñas reales
var CommonSuffixes = []string{
	"1", "12", "123", "1234", "12345", "123456",
	"!", "!!", "!123", "#", "@", ".*",
	"01", "02", "99", "00", "2024", "2023", "2022",
}

// prefixes lista de prefijos comunes
var CommonPrefixes = []string{
	"!", "@", "#", "123", "000",
}

// AppendSuffix agrega un sufijo a la palabra: "carlos" + "123" → "carlos123"
func AppendSuffix(base, suffix string) string {
	return base + suffix
}

// PrependPrefix agrega un prefijo: "123" + "carlos" → "123carlos"
func PrependPrefix(prefix, base string) string {
	return prefix + base
}

// ================================
// COMBINACIONES DE MÚLTIPLES CAMPOS
// ================================

// Combine une dos campos directamente: "carlos" + "1990" → "carlos1990"
func Combine(a, b string) string {
	return a + b
}

// CombineWith une dos campos con separador: "carlos" + "_" + "1990" → "carlos_1990"
func CombineWith(a, sep, b string) string {
	return a + sep + b
}

// ================================
// GENERADOR DE VARIANTES DE UNA PALABRA
// ================================

// AllVariants genera todas las variantes estándar de un string base.
// Este es el motor central que usan los 3 módulos.
func AllVariants(base string) []string {
	if base == "" {
		return nil
	}

	// Usamos un map para evitar duplicados
	seen := make(map[string]bool)
	var results []string

	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			results = append(results, s)
		}
	}

	lower := ToLower(base)
	upper := ToUpper(base)
	cap := Capitalize(base)
	leet := Leet(lower)
	leetCap := Leet(cap)
	rev := Reverse(lower)

	// Variantes base
	add(lower)
	add(upper)
	add(cap)
	add(leet)
	add(leetCap)
	add(rev)

	// Con sufijos comunes
	for _, suf := range CommonSuffixes {
		add(lower + suf)
		add(cap + suf)
		add(upper + suf)
		add(leet + suf)
	}

	// Con prefijos comunes
	for _, pre := range CommonPrefixes {
		add(pre + lower)
		add(pre + cap)
	}

	return results
}
