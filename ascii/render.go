package ascii

import (
	"os"
	"strings"
)

/*
ReadBanner reads the banner file and returns all lines.

Example banner file structure:
Each character is represented by 8 lines,
with a blank line separating characters.
*/
func ReadBanner(file string) ([]string, error) {

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Split file into lines
	lines := strings.Split(string(data), "\n")

	return lines, nil
}

/*
BuildAsciiMap converts banner lines into a map.

Example mapping:

'A' → [
	line1,
	line2,
	...
	line8
]
*/
func BuildAsciiMap(lines []string) map[rune][]string {

	asciiMap := make(map[rune][]string)

	char := 32 // ASCII code for space

	/*
		Each character block in banner file = 9 lines
		8 lines of drawing + 1 empty separator
	*/

	for i := 1; i < len(lines); i += 9 {

		// Store the 8 ASCII lines for this character
		asciiMap[rune(char)] = lines[i : i+8]

		char++
	}

	return asciiMap
}

/*
ANSI color escape codes used for terminal coloring
*/
var colors = map[string]string{
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

/*
PrintAscii converts text into ASCII art.

Parameters:
text      → input string
asciiMap  → mapping of characters to ASCII art
subStr    → substring that should be colored
colorName → color name from the colors map
*/
func PrintAscii(text string, asciiMap map[rune][]string, subStr string, colorName string) string {

	var result strings.Builder

	// Split input text on literal "\n"
	lines := strings.Split(text, "\\n")

	// Get ANSI color codes
	color := colors[colorName]
	reset := colors["reset"]

	for _, line := range lines {

		// Find first occurrence of substring
		subIndex := strings.Index(line, subStr)

		// Handle empty lines
		if line == "" {
			result.WriteString("\n")
			continue
		}

		/*
			Each ASCII character has height = 8 lines
			So we print row-by-row
		*/
		for row := 0; row < 8; row++ {

			for i, char := range line {

				// Start coloring when substring begins
				if i == subIndex {
					result.WriteString(color)
				}

				// Write ASCII art segment
				result.WriteString(asciiMap[char][row])

				// Stop coloring when substring ends
				if i == subIndex+len(subStr)-1 {

					result.WriteString(reset)

					/*
						Look for the next occurrence of substring
						after the current position
					*/
					next := strings.Index(line[i+1:], subStr)

					if next != -1 {
						subIndex = next + i + 1
					} else {
						subIndex = -1
					}
				}
			}

			// Move to next ASCII row
			result.WriteString("\n")
		}
	}

	return result.String()
}