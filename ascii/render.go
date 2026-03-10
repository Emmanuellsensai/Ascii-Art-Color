package ascii

import (
	"fmt"
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

func PrintAscii(text string, asciiMap map[rune][]string) {
	lines := strings.Split(text, "\\n")
	for i, line := range lines {
		if line == "" {
			if i != 0 {
				fmt.Println()
			}
			continue
		}
		for row := 0; row < 8; row++ {
			for _, char := range line {
				fmt.Print(asciiMap[char][row])
			}
			fmt.Println()
		}
	}
}

// Colors maps color names to ANSI escape codes
var Colors = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"orange":  "\033[38;5;214m",
	"reset":   "\033[0m",
}

// PrintAsciiColor prints ASCII art with color applied to all occurrences
// of substr in the text. If substr is empty, the whole output is colored.
func PrintAsciiColor(text, color, substr string, asciiMap map[rune][]string) {
	colorCode := Colors[color]
	reset := Colors["reset"]

	lines := strings.Split(text, "\\n")
	for _, line := range lines {
		if line == "" {
			fmt.Println()
			continue
		}

		// Build a set of rune positions that should be colored.
		// We convert the line to a rune slice so positions match
		// how range iterates over characters.
		runes := []rune(line)
		subRunes := []rune(substr)
		colored := make([]bool, len(runes))

		if substr == "" {
			// Color everything
			for i := range colored {
				colored[i] = true
			}
		} else {
			// Mark every position that falls inside an occurrence of substr
			for i := 0; i <= len(runes)-len(subRunes); i++ {
				match := true
				for j := 0; j < len(subRunes); j++ {
					if runes[i+j] != subRunes[j] {
						match = false
						break
					}
				}
				if match {
					for j := 0; j < len(subRunes); j++ {
						colored[i+j] = true
					}
				}
			}
		}

		for row := 0; row < 8; row++ {
			for i, char := range runes {
				artRow := asciiMap[char][row]
				if colored[i] {
					fmt.Print(colorCode + artRow + reset)
				} else {
					fmt.Print(artRow)
				}
			}
			fmt.Println()
		}
	}
}
