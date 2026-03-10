package main

import (
	"fmt"
	"os"
	"strings"

	"ascii-art/ascii"
)

func main() {
	args := os.Args[1:]

	// Extract --color=VALUE flag if present
	color := ""
	remaining := []string{}
	for _, a := range args {
		if strings.HasPrefix(a, "--color=") {
			color = strings.TrimPrefix(a, "--color=")
		} else {
			remaining = append(remaining, a)
		}
	}

	// Validate color
	if color != "" {
		if _, ok := ascii.Colors[color]; !ok {
			fmt.Println("Usage: go run . [--color=COLOR] [SUBSTRING] STRING [BANNER]")
			return
		}
	}

	// Parse positional args: STRING, optional BANNER, optional SUBSTRING
	// Valid combinations after flag removal:
	//   1 arg:  STRING
	//   2 args: STRING BANNER  OR  SUBSTRING STRING  (only if --color given)
	//   3 args: SUBSTRING STRING BANNER  (only if --color given)
	knownBanners := map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}

	input, substr, banner := "", "", "standard"

	switch len(remaining) {
	case 1:
		input = remaining[0]
	case 2:
		if knownBanners[remaining[1]] {
			input, banner = remaining[0], remaining[1]
		} else if color != "" {
			substr, input = remaining[0], remaining[1]
		} else {
			fmt.Println("Usage: go run . [--color=COLOR] [SUBSTRING] STRING [BANNER]")
			return
		}
	case 3:
		if color != "" && knownBanners[remaining[2]] {
			substr, input, banner = remaining[0], remaining[1], remaining[2]
		} else {
			fmt.Println("Usage: go run . [--color=COLOR] [SUBSTRING] STRING [BANNER]")
			return
		}
	default:
		fmt.Println("Usage: go run . [--color=COLOR] [SUBSTRING] STRING [BANNER]")
		return
	}

	// Load banner and build map
	bannerLines, err := ascii.ReadBanner(banner + ".txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	asciiMap := ascii.BuildAsciiMap(bannerLines)

	// Render
	if color != "" {
		ascii.PrintAsciiColor(input, color, substr, asciiMap)
	} else {
		ascii.PrintAscii(input, asciiMap)
	}
}