package terse

import (
	"regexp"
	"strings"
)

func strippedPrefix(code, prefix string) bool {
	tc := strings.TrimSpace(code)
	return strings.HasPrefix(tc, prefix) && len(tc) > len(prefix)
}

func strippedBegin(code, prefix string) bool {
	return strings.HasPrefix(strings.TrimSpace(code), prefix)
}

func strippedLine(code string) string {
	return strings.TrimSpace(code)
}

var firstTextRegex = regexp.MustCompile("^[a-z0-9_-]+")

func firstTextToken(code string) string {
	return firstTextRegex.FindString(code)
}

func doctypeCode(code string) bool {
	return strippedBegin(code, "!!")
}

func execCode(code string) bool {
	return strippedPrefix(code, "=")
}

func execContCode(code string) bool {
	return strippedPrefix(code, "/=")
}

func verbatimCode(code string) bool {
	return strippedPrefix(code, "/")
}

func commentCode(code string) bool {
	return strippedBegin(code, "//")
}

func tagCode(code string) bool {
	_, ok := ValidElements[firstTextToken(code)]
	return ok || strippedPrefix(code, "%")
}

func filterCode(code string) bool {
	line := strippedLine(code)
	if line[0] == ':' {
		_, ok := Filters[firstTextToken(line[1:])]
		return ok
	}
	return false
}

func blockCode(code string) bool {
	line := strippedLine(code)
	if line[0] == '[' {
		ftt := firstTextToken(line[1:])
		return line[len(ftt)+1] == ']'
	}
	return false
}

func defineBlockCode(code string) bool {
	line := strippedLine(code)
	if line[0] == '^' {
		ftt := firstTextToken(line[1:])
		return line[len(ftt)+1] == ']'
	}
	return false
}

func execBlockCode(code string) bool {
	line := strippedLine(code)
	if line[0] == '$' {
		ftt := firstTextToken(line[1:])
		return line[len(ftt)+1] == ']'
	}
	return false
}

func extendCode(code string) bool {
	return strippedPrefix(code, "@@")
}

func yieldCode(code string) bool {
	return strippedLine(code)[0] == '@'
}

func ifCode(code string) bool {
	return strippedLine(code)[0] == '?'
}

func elseCode(code string) bool {
	return strippedLine(code) == "!?"
}

func rangeCode(code string) bool {
	return strippedLine(code)[0] == '&'
}

func rangeElseCode(code string) bool {
	return strippedLine(code) == "!&"
}

func withCode(code string) bool {
	return strippedLine(code)[0] == '>'
}

func withElseCode(code string) bool {
	return strippedLine(code) == "!>"
}

func idClassCode(code string) bool {
	c := strippedLine(code)[0]
	return c == '.' || c == '#'
}
