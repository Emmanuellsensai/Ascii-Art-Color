package main

import (
	"fmt"
	"os"
	"strings"

	"ascii-art/ascii"
)

func main() {

	/*
		Expected command format:

		go run . --color=<color> <substring> "<text>"

		Example:
		go run . --color=red kit "a king kitten have kit"
	*/

	// Program name + 3 arguments
	if len(os.Args) != 4 {
		return
	}

	// First argument: --color=<color>
	colorArg := os.Args[1]

	// Second argument: substring to color
	subStr := os.Args[2]

	// Third argument: full text to convert to ASCII art
	input := os.Args[3]

	// Ensure the first argument starts with "--color="
	if !strings.HasPrefix(colorArg, "--color=") {
		return
	}

	// Extract the color name from "--color=<color>"
	color := strings.TrimPrefix(colorArg, "--color=")

	// Read the banner file (ASCII template)
	bannerLines, err := ascii.ReadBanner("standard.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Convert banner lines into a map:
	// rune → [8 lines of ASCII art]
	asciiMap := ascii.BuildAsciiMap(bannerLines)

	// Generate ASCII output with coloring
	result := ascii.PrintAscii(input, asciiMap, subStr, color)

	// Print final ASCII art
	fmt.Print(result)
}