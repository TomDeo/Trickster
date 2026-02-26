package transforms

import (
	"fmt"
	"strings"
	"unicode"
)

// ================================
// TRANSFORMACIONES BÁSICAS
// ================================

func ToUpper(s string) string { return strings.ToUpper(s) }
func ToLower(s string) string { return strings.ToLower(s) }

// Capitalize pone la primera letra en mayúscula: "carlos" → "Carlos"
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// CapitalizeAll para nombres compuestos: "juan pablo" → "Juan Pablo"
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
// CAPITALIZACIÓN AVANZADA
// ================================

// ToggleCase alterna mayúsculas: "carlos" → "cArLoS"
func ToggleCase(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if i%2 == 0 {
			runes[i] = unicode.ToLower(r)
		} else {
			runes[i] = unicode.ToUpper(r)
		}
	}
	return string(runes)
}

// ToggleCaseInverse empieza en mayúscula: "carlos" → "CaRlOs"
func ToggleCaseInverse(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if i%2 == 0 {
			runes[i] = unicode.ToUpper(r)
		} else {
			runes[i] = unicode.ToLower(r)
		}
	}
	return string(runes)
}

// CapitalizeLast pone en mayúscula la última letra: "carlos" → "carloS"
func CapitalizeLast(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[len(runes)-1] = unicode.ToUpper(runes[len(runes)-1])
	return string(runes)
}

// CapitalizeFirstLast primera y última en mayúscula: "carlos" → "CarloS"
func CapitalizeFirstLast(s string) string {
	return CapitalizeLast(Capitalize(s))
}

// Duplicate repite la palabra: "carlos" → "carloscarlos"
func Duplicate(s string) string { return s + s }

// DuplicateCapitalized: "carlos" → "carlosCARLOS"
func DuplicateCapitalized(s string) string {
	return ToLower(s) + ToUpper(s)
}

// ================================
// LEET SPEAK
// ================================

var leetMapBasic = map[rune]string{
	'a': "4", 'A': "4",
	'e': "3", 'E': "3",
	'i': "1", 'I': "1",
	'o': "0", 'O': "0",
	's': "5", 'S': "5",
	't': "7", 'T': "7",
	'b': "8", 'B': "8",
	'g': "9", 'G': "9",
}

var leetMapFull = map[rune]string{
	'a': "4", 'A': "4",
	'e': "3", 'E': "3",
	'i': "1", 'I': "1",
	'o': "0", 'O': "0",
	's': "5", 'S': "5",
	't': "7", 'T': "7",
	'b': "8", 'B': "8",
	'g': "9", 'G': "9",
	'l': "1", 'L': "1",
	'z': "2", 'Z': "2",
	'q': "9", 'Q': "9",
}

func applyLeetMap(s string, m map[rune]string) string {
	var result strings.Builder
	for _, r := range s {
		if replacement, ok := m[r]; ok {
			result.WriteString(replacement)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Leet aplica leet speak básico: "carlos" → "c4rl0s"
func Leet(s string) string { return applyLeetMap(s, leetMapBasic) }

// LeetFull aplica leet speak extendido
func LeetFull(s string) string { return applyLeetMap(s, leetMapFull) }

// LeetPartial aplica leet solo a vocales
func LeetPartial(s string) string {
	vocales := map[rune]string{
		'a': "4", 'A': "4",
		'e': "3", 'E': "3",
		'i': "1", 'I': "1",
		'o': "0", 'O': "0",
		'u': "v", 'U': "V",
	}
	return applyLeetMap(s, vocales)
}

// ================================
// TRUNCADO
// ================================

// TruncateLeft toma los primeros N caracteres: "carlos" n=3 → "car"
func TruncateLeft(s string, n int) string {
	runes := []rune(s)
	if n >= len(runes) {
		return s
	}
	return string(runes[:n])
}

// TruncateRight toma los últimos N caracteres: "carlos" n=3 → "los"
func TruncateRight(s string, n int) string {
	runes := []rune(s)
	if n >= len(runes) {
		return s
	}
	return string(runes[len(runes)-n:])
}

// RemoveVowels elimina vocales: "carlos" → "crls"
func RemoveVowels(s string) string {
	vowels := "aeiouáéíóúAEIOUÁÉÍÓÚ"
	var result strings.Builder
	for _, r := range s {
		if !strings.ContainsRune(vowels, r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ================================
// SUFIJOS / PREFIJOS
// ================================

// YearRange genera años (y sus versiones cortas de 2 dígitos) en un rango
func YearRange(from, to int) []string {
	var years []string
	for y := from; y <= to; y++ {
		years = append(years, fmt.Sprintf("%d", y))
		years = append(years, fmt.Sprintf("%02d", y%100))
	}
	return years
}

var CommonSuffixes = []string{
	// Secuencias numéricas
	"1", "12", "123", "1234", "12345", "123456",
	"0", "01", "02", "007", "09", "10",
	// Símbolos
	"!", "!!", "!123", "#", "@", ".*", ".", "*", "?",
	"!@#", "@123", "#123",
	// Años
	"2024", "2023", "2022", "2021", "2020", "2019", "2018", "2000", "1999",
	// Números con símbolos
	"1!", "123!", "1234!", "12!", "99", "00",
	// Dobles
	"11", "22", "33", "44", "55", "66", "77", "88", "99",
}

var CommonPrefixes = []string{
	"!", "@", "#", "123", "000", "1", "el", "la", "los",
}

func AppendSuffix(base, suffix string) string  { return base + suffix }
func PrependPrefix(prefix, base string) string { return prefix + base }
func Combine(a, b string) string               { return a + b }
func CombineWith(a, sep, b string) string      { return a + sep + b }

// ================================
// GENERADOR CENTRAL DE VARIANTES
// ================================

// AllVariants genera todas las variantes posibles de un string base.
// Es el motor central que usan los 3 módulos.
func AllVariants(base string) []string {
	if base == "" {
		return nil
	}

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
	leetFull := LeetFull(lower)
	leetCap := Leet(cap)
	rev := Reverse(lower)

	// --- Variantes base ---
	add(lower)
	add(upper)
	add(cap)
	add(leet)
	add(leetFull)
	add(leetCap)
	add(rev)
	add(ToggleCase(lower))
	add(ToggleCaseInverse(lower))
	add(RemoveVowels(lower))
	add(CapitalizeLast(lower))
	add(CapitalizeFirstLast(lower))
	add(Duplicate(lower))
	add(DuplicateCapitalized(lower))
	add(LeetPartial(lower))

	// --- Con sufijos comunes ---
	for _, suf := range CommonSuffixes {
		add(lower + suf)
		add(cap + suf)
		add(upper + suf)
		add(leet + suf)
		add(leetCap + suf)
	}

	// --- Con prefijos comunes ---
	for _, pre := range CommonPrefixes {
		add(pre + lower)
		add(pre + cap)
	}

	// --- Con rango de años 2000-2025 ---
	for _, year := range YearRange(2000, 2025) {
		add(lower + year)
		add(cap + year)
		add(leet + year)
	}

	// --- Truncados con sufijos: car123, carl! etc ---
	for _, n := range []int{3, 4, 5} {
		trunc := TruncateLeft(lower, n)
		if trunc != lower {
			add(trunc)
			for _, suf := range CommonSuffixes {
				add(trunc + suf)
			}
		}
	}

	return results
}
