package ascii

import (
	// "fmt"
	"os"
	"strings"
)

// 	Reand banner files
// 	func Render(text string, banner string) {
// 	bannerLines := readBanner(banner)

// 	// convert banner file to map
// 	asciiMap := buildAsciiMap(bannerLines)

// 	// print ascii art
// 	printAscii(text, asciiMap)

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

	lines := strings.Split(text, "\\n")

	// subStr := "kit"

	colorCode := Colors[colorArg]
	reset := Colors["reset"]

	for i, line := range lines {

		subIndex := strings.Index(line, subStr)
		// fmt.Println("subIndex:", subIndex, line, len(line), len(subStr), 8 == subIndex+len(subStr))

		if line == "" {
			if i != 0 {
				result.WriteString("\n")
			}
			continue
		}

		for row := 0; row < 8; row++ {

			for i, char := range line {
				// fmt.Println(subIndex, i, line[i:])
				if i == subIndex {
					// color from colors map
					result.WriteString(colorCode)
				}
				result.WriteString(asciiMap[char][row])
				if i == subIndex+len(subStr)-1 || i == len(line) - 1 {
					result.WriteString(reset)
					// fmt.Println("resetting color:", subIndex, i, line[i:])
					subIndex = strings.Index(line[i:], subStr) + i
					// fmt.Println("after resetting color:", subIndex, i, line[i:])
				}
			}
			subIndex = strings.Index(line, subStr)

			result.WriteString("\n")
		}
	}
	return result.String()
}
