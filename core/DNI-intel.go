package core

import (
	"fmt"
	"strings"
	"trickster/transforms"
)

// ================================================================
// MÓDULO: DNI INTELIGENTE ARGENTINA
//
// El DNI argentino es secuencial a nivel nacional desde 1968 (Ley 17.671).
// Correlación año de nacimiento → rango de DNI (fuente: RENAPER + datos públicos):
//
//   Nacidos ~1945-1959 → DNI  1.000.000 –  9.999.999
//   Nacidos ~1960-1969 → DNI 10.000.000 – 19.999.999
//   Nacidos ~1970-1979 → DNI 20.000.000 – 29.999.999
//   Nacidos ~1980-1989 → DNI 30.000.000 – 39.999.999
//   Nacidos ~1990-1999 → DNI 40.000.000 – 49.999.999  (±500k de margen)
//   Nacidos ~2000-2009 → DNI 50.000.000 – 59.999.999
//   Nacidos ~2010-2023 → DNI 70.000.000 – 79.999.999  (60M reservado para CUIL/CUIT extranjeros)
//
// NOTA: La correlación NO es exacta. Hay margen de ±1-2M por:
//   - Personas que tramitaron el DNI tardíamente
//   - Regularización de extranjeros en ciertos períodos
//   - Reemisión por pérdida/daño
//   - Aceleración demográfica post-2012 (digitalización masiva)
//
// El generador produce el rango probable ±margen, formateado como lo
// escribe la gente en contraseñas: con y sin puntos, con y sin ceros iniciales.
// ================================================================

// dniRangeForBirthYear devuelve el rango [min, max] de DNI probable
// para una persona nacida en el año dado.
// Margen amplio para cubrir tardíos y excepciones.
func dniRangeForBirthYear(birthYear int) (min, max int) {
	switch {
	case birthYear <= 1949:
		return 1_000_000, 4_999_999
	case birthYear <= 1959:
		return 4_000_000, 9_999_999
	case birthYear <= 1969:
		return 9_000_000, 19_999_999
	case birthYear <= 1979:
		return 18_000_000, 29_999_999
	case birthYear <= 1984:
		return 27_000_000, 34_999_999
	case birthYear <= 1989:
		return 33_000_000, 39_999_999
	case birthYear <= 1994:
		return 38_000_000, 45_999_999
	case birthYear <= 1999:
		return 44_000_000, 51_999_999
	case birthYear <= 2004:
		return 50_000_000, 57_999_999
	case birthYear <= 2009:
		return 55_000_000, 59_999_999
	case birthYear <= 2015:
		return 70_000_000, 74_999_999
	default:
		return 73_000_000, 79_999_999
	}
}

// DNIFormats genera todas las formas en que una persona escribe su DNI
// en una contraseña: con/sin puntos, con/sin espacios, compacto.
func DNIFormats(dni int) []string {
	raw := fmt.Sprintf("%d", dni)

	// Formato con puntos (30.123.456)
	withDots := formatDNIWithDots(raw)

	// Formato compacto sin puntos (30123456)
	compact := raw

	// Formato con guiones (30-123-456) — menos común pero existe
	withDashes := formatDNIWithDashes(raw)

	var forms []string
	seen := make(map[string]bool)
	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			forms = append(forms, s)
		}
	}

	add(compact)
	add(withDots)
	add(withDashes)

	// Sin el primer dígito (a veces omiten el grupo de millones si es 1 dígito)
	if len(raw) == 7 {
		add(raw[1:]) // 1234567 → 234567
	}

	return forms
}

func formatDNIWithDots(raw string) string {
	// 30123456 → 30.123.456
	// 7 dígitos: X.XXX.XXX
	// 8 dígitos: XX.XXX.XXX
	l := len(raw)
	if l == 7 {
		return raw[0:1] + "." + raw[1:4] + "." + raw[4:7]
	}
	if l == 8 {
		return raw[0:2] + "." + raw[2:5] + "." + raw[5:8]
	}
	return raw
}

func formatDNIWithDashes(raw string) string {
	l := len(raw)
	if l == 7 {
		return raw[0:1] + "-" + raw[1:4] + "-" + raw[4:7]
	}
	if l == 8 {
		return raw[0:2] + "-" + raw[2:5] + "-" + raw[5:8]
	}
	return raw
}

// GenerateDNICandidates genera candidatos de DNI para el año de nacimiento dado.
// Produce una wordlist de DNIs probables en todos los formatos relevantes,
// más combinaciones con el nombre si se provee.
//
// step controla la densidad: step=1000 produce ~1000 candidatos por millón de rango,
// step=5000 produce ~200 por millón (más rápido, menos exhaustivo).
func GenerateDNICandidates(birthYear int, nombre string, step int) []string {
	if step <= 0 {
		step = 1000
	}

	min, max := dniRangeForBirthYear(birthYear)
	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		s = strings.TrimSpace(s)
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	n := strings.ToLower(strings.TrimSpace(nombre))
	nc := transforms.Capitalize(n)

	for dni := min; dni <= max; dni += step {
		for _, f := range DNIFormats(dni) {
			add(f)

			// DNI solo con sufijos comunes
			add(f + "!")
			add(f + ".")

			// Nombre + DNI (patrón muy común en Argentina)
			if n != "" {
				add(n + f)
				add(nc + f)
				add(f + n)
				add(f + nc)
				add(n + "." + f)
				add(n + "_" + f)
			}
		}
	}

	return result
}

// DNIVariantsFromKnown genera variantes cuando el DNI ya se conoce exactamente.
// Más exhaustivo que GenerateDNICandidates porque el DNI real ya está dado.
func DNIVariantsFromKnown(dniStr string, nombre string, apellido string, anio string) []string {
	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			return
		}
		seen[s] = true
		result = append(result, s)
	}

	n := strings.ToLower(strings.TrimSpace(nombre))
	nc := transforms.Capitalize(n)
	a := strings.ToLower(strings.TrimSpace(apellido))
	ac := transforms.Capitalize(a)

	// Normalizar el DNI: quitar puntos y guiones
	dniClean := strings.ReplaceAll(dniStr, ".", "")
	dniClean = strings.ReplaceAll(dniClean, "-", "")
	dniClean = strings.TrimSpace(dniClean)

	dniNum := 0
	fmt.Sscanf(dniClean, "%d", &dniNum)

	var dniFormats []string
	if dniNum > 0 {
		dniFormats = DNIFormats(dniNum)
	} else {
		dniFormats = []string{dniClean}
	}

	for _, f := range dniFormats {
		add(f)

		// DNI + sufijos
		for _, suf := range []string{"!", "@", "#", ".", "1", "12", "123"} {
			add(f + suf)
		}

		// Nombre + DNI
		if n != "" {
			add(n + f)
			add(nc + f)
			add(f + n)
			add(f + nc)
			add(n + "." + f)
			add(n + "_" + f)
			add(nc + "." + f)
		}

		// Apellido + DNI
		if a != "" {
			add(a + f)
			add(ac + f)
			add(f + a)
		}

		// DNI + año
		if anio != "" {
			add(f + anio)
			add(anio + f)
			add(f + anio[2:]) // DNI + año corto
		}

		// Nombre + DNI + sufijo
		if n != "" {
			for _, suf := range []string{"!", "@", "1", "123"} {
				add(n + f + suf)
				add(nc + f + suf)
			}
		}
	}

	return result
}


