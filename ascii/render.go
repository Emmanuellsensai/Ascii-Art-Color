package ascii

import (
	"os"
	"strings"
)

func ReadBanner(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func BuildAsciiMap(lines []string) map[rune][]string {
	asciiMap := make(map[rune][]string)
	char := 32
	for i := 1; i < len(lines); i += 9 {
		asciiMap[rune(char)] = lines[i : i+8]
		char++
	}
	return asciiMap
}

var Colors = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"purple":  "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"orange":  "\033[38;5;214m",
	"reset":   "\033[0m",
}

func PrintAscii(text string, asciiMap map[rune][]string, subStr, colorArg string) string {
	var result strings.Builder

	colorCode := Colors[colorArg]
	reset := Colors["reset"]

	lines := strings.Split(text, "\\n")

	for i, line := range lines {
		if line == "" {
			if i != 0 {
				result.WriteString("\n")
			}
			continue
		}

		for row := 0; row < 8; row++ {
			if subStr == "" {
				// No substring — color every character individually.
				// We must write reset after EACH character so the color
				// does not bleed into characters we did not intend to color.
				// (Without reset, the terminal keeps the color active until
				// it receives another escape code.)
				for _, char := range line {
					result.WriteString(colorCode + asciiMap[char][row] + reset)
				}
			} else {
				// Substring provided — color only the matching parts.
				// strings.Split removes every occurrence of subStr from the
				// line and gives back the plain pieces that sat between them.
				// We render each plain piece as-is, then re-insert subStr in
				// color between every pair of pieces.
				parts := strings.Split(line, subStr)
				for p, part := range parts {
					for _, char := range part {
						result.WriteString(asciiMap[char][row])
					}
					// After every piece except the last, there was one
					// occurrence of subStr — render it in color.
					if p < len(parts)-1 {
						for _, char := range subStr {
							result.WriteString(colorCode + asciiMap[char][row] + reset)
						}
					}
				}
			}
			result.WriteString("\n")
		}
	}
	return result.String()
}