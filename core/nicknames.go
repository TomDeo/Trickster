package core

import (
	"strings"
	"trickster/transforms"
)

// ================================
// DICCIONARIO DE APODOS
// Más de 100 nombres comunes en español con sus apodos reales
// ================================

var nicknameDict = map[string][]string{
	// ---- Masculinos ----
	"alejandro": {"ale", "alex", "alejo", "alejan", "xander", "jandro"},
	"alberto":   {"beto", "alber", "al", "bert"},
	"alfredo":   {"fredo", "fred", "alfre", "alfi"},
	"andres":    {"andy", "andre", "andresito", "andes"},
	"antonio":   {"tony", "toni", "anto", "toño", "antoñito"},
	"ariel":     {"ari", "arielito"},
	"augusto":   {"gus", "augus", "tito"},
	"benjamin":  {"benja", "benji", "ben"},
	"carlos":    {"cali", "carl", "carli", "cha", "charly", "carlitos", "carly"},
	"christian": {"chris", "cris", "chri", "christianito"},
	"claudio":   {"clau", "claudito"},
	"cristian":  {"cris", "chris", "cristi", "cristiano"},
	"daniel":    {"dani", "dan", "danny", "danielito"},
	"david":     {"dav", "davit", "davidcho", "dave"},
	"diego":     {"die", "diegui", "dieghito", "diegote"},
	"eduardo":   {"edu", "eduar", "eddie", "lalo", "edy"},
	"emilio":    {"emi", "emilito", "mil"},
	"enrique":   {"quique", "kike", "enri", "henry"},
	"ernesto":   {"nesto", "ernes", "nestor", "ernie"},
	"esteban":   {"este", "esteba", "teba", "steve"},
	"ezequiel":  {"eze", "ezequielito", "zeke"},
	"facundo":   {"facu", "faku", "facun"},
	"federico":  {"fede", "fed", "freddy", "fedo"},
	"felipe":    {"feli", "pipe", "pippo", "pipo"},
	"fernando":  {"fer", "fercho", "nando", "fernan"},
	"francisco": {"fran", "pancho", "cisco", "franci", "kiko", "paco", "francho"},
	"gabriel":   {"gabi", "gabo", "gabri", "gabriel"},
	"gonzalo":   {"gonza", "gonzi", "gonzo", "gon"},
	"guillermo": {"guille", "willy", "will", "guillerme", "memo"},
	"gustavo":   {"gus", "gusti", "tavo", "gusta"},
	"hernando":  {"hernan", "hernie", "nando"},
	"horacio":   {"hora", "horacito", "hor"},
	"hugo":      {"hugito", "huguito", "hug"},
	"ignacio":   {"nacho", "igna", "iñaki", "nachito"},
	"ivan":      {"iva", "ivancho", "ivanito"},
	"javier":    {"javi", "xavi", "javiercho", "jabier"},
	"jesus":     {"jesusito", "chucho", "chuy", "jesu"},
	"jorge":     {"jorgito", "jorgi", "george", "jor"},
	"jose":      {"pepe", "josecito", "chepe", "joselito", "josepe"},
	"juan":      {"juancho", "juancito", "juani", "johnny", "juanito"},
	"julian":    {"juli", "juliancho", "juliancito"},
	"leandro":   {"lea", "lean", "leandrito", "lechuga"},
	"leonardo":  {"leo", "leon", "leonardito", "lenny"},
	"lorenzo":   {"loren", "lorencito", "renzo"},
	"lucas":     {"luca", "luquitas", "luqui"},
	"luis":      {"lucho", "luisito", "luisi"},
	"manuel":    {"manu", "manolo", "manuelito", "man"},
	"marcelo":   {"marce", "marchelo", "marcelito"},
	"marcos":    {"marki", "marquitos", "marc"},
	"martin":    {"marti", "martincho", "tito", "martn"},
	"mateo":     {"mate", "mateito", "teo"},
	"matias":    {"mati", "maticho", "tias"},
	"mauricio":  {"mauri", "mau", "mauricito"},
	"maximiliano": {"maxi", "max", "maxito"},
	"miguel":    {"miguelito", "migue", "mike", "mikel"},
	"nicolas":   {"nico", "nikolas", "nicol", "nicolasito"},
	"oscar":     {"osquitar", "osca", "ozzy"},
	"pablo":     {"pabli", "pablito", "pabs"},
	"patricio":  {"patri", "pato", "patrizio"},
	"paulo":     {"pau", "paulito"},
	"pedro":     {"pedrito", "pete", "pedrolo", "piero"},
	"rafael":    {"rafa", "rafita", "rafo"},
	"raul":      {"raulito", "rau", "ralito"},
	"ricardo":   {"ricky", "rico", "ricar", "richardito"},
	"roberto":   {"rober", "beto", "bob", "robertito"},
	"rodrigo":   {"rodri", "rod", "rodrigo", "rodriguito"},
	"romulo":    {"romi", "romu", "roms"},
	"ruben":     {"rubi", "rubencito", "rube"},
	"salvador":  {"salva", "chava", "salvadorito"},
	"santiago":  {"santi", "sandy", "santito", "tiago"},
	"sebastian": {"seba", "sebas", "sebi", "sebita"},
	"sergio":    {"sergi", "sergy", "sergito"},
	"tomas":     {"tomi", "tommy", "tomasito", "tom"},
	"victor":    {"vict", "victo", "vic", "victorito"},
	"walter":    {"wally", "walterito", "walt"},

	// ---- Femeninos ----
	"adriana":   {"adri", "adry", "adrianita"},
	"agustina":  {"agus", "tina", "agustinita"},
	"alejandra": {"ale", "alex", "aleja", "alejandrita"},
	"andrea":    {"andy", "andre", "andreita"},
	"analia":    {"ana", "anali", "analita"},
	"angelica":  {"angel", "angie", "angeliquita"},
	"barbara":   {"barbi", "barbie", "barbarita"},
	"beatriz":   {"bea", "beti", "beatricita"},
	"belen":     {"belu", "belit", "belencita"},
	"brenda":    {"bren", "brendita"},
	"camila":    {"cami", "camilita", "mila"},
	"carla":     {"carlita", "carly", "car"},
	"carolina":  {"caro", "carol", "carolinita", "lina"},
	"catalina":  {"cata", "cati", "catita", "lina"},
	"cecilia":   {"ceci", "cecilita", "cily"},
	"celeste":   {"cele", "celestita"},
	"claudia":   {"clau", "claudiita"},
	"constanza": {"coni", "consti", "constanzita"},
	"daniela":   {"dani", "danny", "danielita"},
	"diana":     {"dianita", "diany"},
	"elena":     {"ele", "elenita", "lena"},
	"emilia":    {"emi", "emilita", "mili"},
	"estefania": {"este", "stefi", "estefanita", "fany"},
	"eugenia":   {"euge", "eugenita"},
	"fernanda":  {"fer", "ferchu", "fernandita", "nanda"},
	"florencia": {"flor", "florita", "florencita"},
	"gabriela":  {"gabi", "gabo", "gabrielita"},
	"gisela":    {"gise", "giselita"},
	"graciela":  {"graci", "grace", "gracielita"},
	"guadalupe": {"guada", "lupe", "lupita", "guadalupita"},
	"jimena":    {"jime", "jimenita"},
	"josefina":  {"jofi", "jose", "josefinita", "fina"},
	"julieta":   {"juli", "julie", "julietita"},
	"karina":    {"kari", "karinita"},
	"laura":     {"lau", "laurita", "lauri"},
	"leticia":   {"leti", "leticita"},
	"lorena":    {"lore", "lorenita"},
	"lucia":     {"lu", "luci", "lucita", "luciana"},
	"luciana":   {"luci", "lu", "lucianita", "lucy"},
	"luisa":     {"lui", "luisita"},
	"magdalena": {"magda", "made", "magdalenita"},
	"marcela":   {"marce", "marcelita"},
	"maria":     {"mari", "mary", "maruja", "mariita"},
	"mariana":   {"mari", "mary", "marianita", "ana"},
	"marina":    {"mari", "marinita"},
	"martina":   {"marti", "tinita", "martinita"},
	"mercedes":  {"meche", "merce", "merche"},
	"micaela":   {"mica", "micaelita", "micky"},
	"monica":    {"moni", "monicaita"},
	"natalia":   {"nati", "natalita", "nats"},
	"noelia":    {"noe", "noelita"},
	"paola":     {"pao", "paolita"},
	"patricia":  {"patri", "pato", "patricita"},
	"paula":     {"pau", "paulita"},
	"pilar":     {"pili", "pilarcha", "pilarita"},
	"romina":    {"romi", "rominita"},
	"rosa":      {"rosita", "rosi"},
	"sabrina":   {"sabri", "sabrinita"},
	"silvana":   {"silvi", "silvarita"},
	"silvia":    {"silvi", "silvita"},
	"sofia":     {"sofi", "sofita"},
	"soledad":   {"sole", "soledadita"},
	"valeria":   {"vale", "vali", "valeriacita"},
	"valentina": {"vale", "valen", "valentinita", "tina"},
	"vanesa":    {"vane", "vanesita"},
	"veronica":  {"vero", "veronicaita", "veri"},
	"victoria":  {"vicky", "vic", "victorita", "viki"},
	"viviana":   {"vivi", "vivianita"},
	"yamila":    {"yami", "yamilita"},
	"yesica":    {"yesi", "yesicita"},
}

// ================================
// REGLAS DE GENERACIÓN DE APODOS
// ================================

// generateRuleBasedNicknames genera apodos automáticos por reglas
// cuando el nombre no está en el diccionario
func generateRuleBasedNicknames(name string) []string {
	lower := strings.ToLower(name)
	var nicks []string
	seen := make(map[string]bool)

	add := func(s string) {
		s = strings.TrimSpace(s)
		if s != "" && s != lower && !seen[s] {
			seen[s] = true
			nicks = append(nicks, s)
		}
	}

	runes := []rune(lower)
	length := len(runes)

	// Truncados clásicos
	if length > 3 {
		add(string(runes[:3])) // primeras 3 letras
	}
	if length > 4 {
		add(string(runes[:4])) // primeras 4 letras
	}
	if length > 5 {
		add(string(runes[:5])) // primeras 5 letras
	}

	// Diminutivos en español
	add(lower + "ito")
	add(lower + "ita")
	add(lower + "in")
	add(lower + "i")
	add(lower + "y")

	// Si termina en vocal, truncar y agregar diminutivo
	if length > 3 {
		base := string(runes[:length-1]) // quitar última letra
		add(base + "i")
		add(base + "ito")
		add(base + "ita")
	}

	// Si termina en 'o', reemplazar por 'i' (ej: roberto → roberti)
	if length > 2 && runes[length-1] == 'o' {
		add(string(runes[:length-1]) + "i")
	}

	// Si termina en 'a', reemplazar por 'i' (ej: camila → camili)
	if length > 2 && runes[length-1] == 'a' {
		add(string(runes[:length-1]) + "i")
	}

	return nicks
}

// GetNicknames devuelve todos los apodos posibles para un nombre dado.
// Combina el diccionario curado + generación por reglas.
func GetNicknames(name string) []string {
	lower := strings.ToLower(strings.TrimSpace(name))
	if lower == "" {
		return nil
	}

	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		s = strings.ToLower(strings.TrimSpace(s))
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	// 1. Buscar en diccionario curado
	if nicks, ok := nicknameDict[lower]; ok {
		for _, n := range nicks {
			add(n)
		}
	}

	// 2. Agregar siempre los generados por reglas
	for _, n := range generateRuleBasedNicknames(lower) {
		add(n)
	}

	return result
}

// NicknameVariants devuelve todos los apodos + todas sus variantes (leet, caps, sufijos)
// Este es el método más completo para el módulo 3
func NicknameVariants(name string) []string {
	nicks := GetNicknames(name)
	if len(nicks) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var result []string

	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	for _, nick := range nicks {
		// Cada apodo pasa por AllVariants para obtener todas sus mutaciones
		for _, variant := range transforms.AllVariants(nick) {
			add(variant)
		}
	}

	return result
}
