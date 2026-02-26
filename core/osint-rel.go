package core

import (
	"fmt"
	"strings"
	"trickster/transforms"
	"trickster/utils"
)

// ================================================================
// MÓDULO: OSINT RELATIVES
//
// Captura nombres de personas/mascotas del entorno del objetivo y los
// combina con años relevantes.
//
// Patrones cubiertos (fuente: análisis de contraseñas filtradas de usuarios
// latinoamericanos + CUPP "special relationships" module):
//
//   [NombreHijo][AñoNacimientoHijo]   → "valentina2015", "Valentina2015!"
//   [NombreHijo][AñoActual]            → "valentina2024"
//   [NombreMascota][Número]            → "firulais123", "Firulais1!"
//   [NombreHijo]+[NombrePareja]        → "valentinajuan", "ValentinaJuan"
//   [NombreObjetivo]+[NombreHijo]      → "carlosvalentina"
//   [NombreMascota]+[AñoAdopción]      → "firulais2019"
//
// El módulo también acepta otros familiares: padres, hermanos, pareja.
// ================================================================

// Relative representa a una persona o mascota del entorno del objetivo
type Relative struct {
	Nombre    string
	TipoVinc  string // "hijo", "pareja", "mascota", "padre", "madre", "hermano", etc.
	AnioNac   string // año de nacimiento o adopción (opcional)
}

// RelativesProfile agrupa todos los familiares/mascotas capturados por OSINT
type RelativesProfile struct {
	Parientes []Relative
}

// AskRelatives guía al usuario para ingresar familiares y mascotas del objetivo.
// Devuelve un RelativesProfile listo para usar en generación.
func AskRelatives() RelativesProfile {
	fmt.Println()
	utils.Info("Ingresá personas/mascotas del entorno del objetivo (OSINT: Instagram, Facebook, etc.)")
	utils.Info("Dejá en blanco el nombre para terminar.")
	fmt.Println()

	rp := RelativesProfile{}

	tiposComunes := []string{"hijo/a", "pareja", "mascota", "padre", "madre", "hermano/a", "amigo/a"}
	_ = tiposComunes // referencia para el prompt

	for i := 1; i <= 10; i++ {
		nombre := utils.AskOptional(fmt.Sprintf("  Nombre %d (persona/mascota)", i))
		if strings.TrimSpace(nombre) == "" {
			break
		}
		tipo := utils.AskOptional(fmt.Sprintf("  Vínculo (ej: hijo, mascota, pareja)"))
		anio := utils.AskOptional(fmt.Sprintf("  Año nacimiento/adopción (opcional, ej: 2015)"))

		rp.Parientes = append(rp.Parientes, Relative{
			Nombre:   strings.TrimSpace(nombre),
			TipoVinc: strings.TrimSpace(tipo),
			AnioNac:  strings.TrimSpace(anio),
		})
	}

	return rp
}

// GenerateFromRelatives genera candidatos de contraseña a partir de
// los familiares/mascotas del objetivo combinados con el perfil principal.
func GenerateFromRelatives(rp RelativesProfile, p Profile) []string {
	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		s = strings.TrimSpace(s)
		l := len([]rune(s))
		if l < 4 || l > 28 || seen[s] {
			return
		}
		seen[s] = true
		result = append(result, s)
	}

	// Años relevantes para combinaciones con hijos/mascotas:
	// Cubrimos 2005-2025 (años en que la mayoría tiene hijos o mascotas)
	childYears := []string{
		"2005", "2006", "2007", "2008", "2009",
		"2010", "2011", "2012", "2013", "2014",
		"2015", "2016", "2017", "2018", "2019",
		"2020", "2021", "2022", "2023", "2024", "2025",
	}
	childYearsShort := []string{
		"05", "06", "07", "08", "09",
		"10", "11", "12", "13", "14",
		"15", "16", "17", "18", "19",
		"20", "21", "22", "23", "24", "25",
	}

	nombreObjetivo := strings.ToLower(strings.TrimSpace(p.Nombre))
	apellidoObjetivo := strings.ToLower(strings.TrimSpace(p.Apellido))

	for _, rel := range rp.Parientes {
		rn := strings.ToLower(strings.TrimSpace(rel.Nombre))
		if rn == "" {
			continue
		}
		rnc := transforms.Capitalize(rn)
		rnu := transforms.ToUpper(rn)
		rnLeet := leetSimple(rn)

		// ── Formas base del familiar ─────────────────────────
		add(rn)
		add(rnc)
		add(rnu)
		add(rnLeet)
		add(transforms.Reverse(rn))

		// ── Familiar + sufijos numéricos comunes ─────────────
		// (patrones más frecuentes en contraseñas con nombres propios)
		for _, num := range []string{
			"1", "2", "3", "12", "21", "123", "1234", "12345",
			"0", "00", "01", "007",
			"111", "222", "333", "777", "999",
		} {
			add(rn + num)
			add(rnc + num)
		}

		// ── Familiar + sufijos de símbolo ─────────────────────
		for _, sp := range []string{"!", "!!", ".", "@", "#", "1!", "123!", "!1"} {
			add(rn + sp)
			add(rnc + sp)
		}

		// ── Familiar + año conocido (si se ingresó) ───────────
		if rel.AnioNac != "" {
			ay := rel.AnioNac
			ayShort := ""
			if len(ay) == 4 {
				ayShort = ay[2:]
			}

			add(rn + ay)
			add(rnc + ay)
			add(rnu + ay)
			add(rnLeet + ay)
			add(ay + rn)
			add(ay + rnc)

			if ayShort != "" {
				add(rn + ayShort)
				add(rnc + ayShort)
				add(ayShort + rn)
			}

			for _, sp := range []string{"!", "@", "#", ".", "1", "123"} {
				add(rn + ay + sp)
				add(rnc + ay + sp)
				if ayShort != "" {
					add(rn + ayShort + sp)
					add(rnc + ayShort + sp)
				}
			}

			// Sándwich: año-familiar-año
			add(ay + rn + ay)
			if ayShort != "" {
				add(ayShort + rn + ayShort)
			}
		}

		// ── Familiar × años relevantes (aunque no se conozca el año) ─
		// Para hijos/mascotas cubrimos 2005-2025
		isMascotaOHijo := strings.Contains(rel.TipoVinc, "hijo") ||
			strings.Contains(rel.TipoVinc, "mascota") ||
			strings.Contains(rel.TipoVinc, "hija") ||
			rel.TipoVinc == ""

		if isMascotaOHijo {
			for _, y := range childYears {
				add(rn + y)
				add(rnc + y)
				add(y + rn)
			}
			for _, y := range childYearsShort {
				add(rn + y)
				add(rnc + y)
			}
		}

		// ── Familiar + nombre/apellido del objetivo ───────────
		if nombreObjetivo != "" {
			no := nombreObjetivo
			noc := transforms.Capitalize(no)

			add(no + rn)
			add(rn + no)
			add(noc + rnc)
			add(rnc + noc)
			add(no + "_" + rn)
			add(rn + "_" + no)
			add(no + "." + rn)

			if p.Anio != "" {
				add(no + rn + p.Anio)
				add(rn + no + p.Anio)
				add(noc + rnc + p.Anio)
			}
		}

		if apellidoObjetivo != "" {
			ao := apellidoObjetivo
			aoc := transforms.Capitalize(ao)
			add(rn + ao)
			add(rnc + aoc)
			add(ao + rn)
			if p.Anio != "" {
				add(rn + ao + p.Anio)
				add(rnc + aoc + p.Anio)
			}
		}

		// ── Combinaciones entre parientes ─────────────────────
		for _, rel2 := range rp.Parientes {
			rn2 := strings.ToLower(strings.TrimSpace(rel2.Nombre))
			if rn2 == "" || rn2 == rn {
				continue
			}
			rn2c := transforms.Capitalize(rn2)

			add(rn + rn2)
			add(rnc + rn2c)
			add(rn + "_" + rn2)

			if p.Anio != "" {
				add(rn + rn2 + p.Anio)
				add(rnc + rn2c + p.Anio)
			}
		}

		// ── Apodos del familiar ───────────────────────────────
		nicks := GetNicknames(rel.Nombre)
		for _, nick := range nicks {
			nc := transforms.Capitalize(nick)
			add(nick)
			add(nc)

			for _, num := range []string{"1", "12", "123", "0", "00"} {
				add(nick + num)
				add(nc + num)
			}
			for _, sp := range []string{"!", "!!", "@", "."} {
				add(nick + sp)
				add(nc + sp)
			}

			if rel.AnioNac != "" {
				add(nick + rel.AnioNac)
				add(nc + rel.AnioNac)
			}
			if isMascotaOHijo {
				for _, y := range childYears {
					add(nick + y)
					add(nc + y)
				}
			}
			if nombreObjetivo != "" {
				add(nombreObjetivo + nick)
				add(nick + nombreObjetivo)
			}
		}

		// ── Leet de todos los años del familiar ───────────────
		if rel.AnioNac != "" {
			for _, lv := range leetAllVariants(rn) {
				add(lv + rel.AnioNac)
				add(transforms.Capitalize(lv) + rel.AnioNac)
			}
		}
	}

	return result
}
