package goname

import (
	"strings"
)

func UpperCamelCase(value string) string {
	splited := splitCases(value)
	for i, str := range splited {
		splited[i] = uppserCamelOneWord(str)
	}
	return strings.Join(splited, "")
}

func LowerCamelCase(value string) string {
	splited := splitCases(value)
	for i, str := range splited {
		if i == 0 {
			splited[i] = strings.ToLower(str)
		} else {
			splited[i] = uppserCamelOneWord(str)
		}
	}
	return strings.Join(splited, "")
}

func LowerSnakeCase(value string) string {
	splited := splitCases(value)
	for i, str := range splited {
		splited[i] = strings.ToLower(str)
	}
	return strings.Join(splited, "_")
}

func splitCases(str string) []string {
	splited := []string{}
	runeStr := []rune(str)

	for 0 < len(runeStr) {
		i := 0
		size := len(runeStr)
		plower := false
		pupper := false
		palfabet := false
		upperCases := false
		pnumber := false
		for {
			r := runeStr[i]
			lower := 'a' <= r && r <= 'z'
			upper := 'A' <= r && r <= 'Z'
			alfabet := lower || upper
			number := '0' <= r && r <= '9'
			other := !(alfabet || number)

			if other {
				if i != 0 {
					splited = append(splited, string(runeStr[:i]))
				}
				i++
				break
			}

			if i > 0 {
				if plower && upper {
					splited = append(splited, string(runeStr[:i]))
					break
				}
				if pnumber && !number {
					splited = append(splited, string(runeStr[:i]))
					break
				}
				if palfabet && !alfabet {
					splited = append(splited, string(runeStr[:i]))
					break
				}
				if upperCases && lower {
					i -= 1
					splited = append(splited, string(runeStr[:i]))
					break
				}
				if pupper && upper {
					upperCases = true
				}
			}
			plower = lower
			pupper = upper
			palfabet = alfabet
			pnumber = number

			i++

			if i >= size {
				if len(runeStr) > 0 {
					splited = append(splited, string(runeStr))
				}
				break
			}
		}
		runeStr = runeStr[i:]
	}

	return concatInitial(splited)
}

// https://github.com/golang/lint/blob/master/lint.go#L770
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

func concatInitial(splited []string) []string {
	newSplited := []string{}
	buffer := []string{}
	for _, word := range splited {
		buffer = append(buffer, word)
		concatWord := strings.Join(buffer, "")
		concatUpperWord := strings.ToUpper(concatWord)

		if commonInitialisms[concatUpperWord] {
			buffer = []string{}
			newSplited = append(newSplited, concatUpperWord)
			continue
		}

		hasPrefix := false
		for initial := range commonInitialisms {
			hasPrefix = strings.HasPrefix(initial, concatUpperWord)
			if hasPrefix {
				break
			}
		}

		if !hasPrefix {
			newSplited = append(newSplited, buffer...)
			buffer = []string{}
		}
	}

	return newSplited
}

func uppserCamelOneWord(word string) string {
	if word == "" {
		return ""
	}

	upper := strings.ToUpper(word)
	if commonInitialisms[upper] {
		return upper
	}

	runeWord := []rune(word)
	return strings.ToUpper(string(runeWord[:1])) + strings.ToLower(string(runeWord[1:]))
}
