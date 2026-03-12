package main

import (
	"ascii-art-color/ascii"
	"fmt"
	"os"
	"strings"
)

func main() {

	/*
		Valid usages:

		go run . --color=<color> "text"
		go run . --color=<color> <substring> "text"
	*/

	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Println("Usage: go run . [OPTION] [STRING]")
		fmt.Println()
		fmt.Println("EX: go run . --color=<color> <substring to be colored> \"something\"")
		return
	}

	colorArg := os.Args[1]

	if !strings.HasPrefix(colorArg, "--color=") {
		fmt.Println("Usage: go run . [OPTION] [STRING]")
		fmt.Println()
		fmt.Println("EX: go run . --color=<color> <substring to be colored> \"something\"")
		return
	}

	// extract color name
	colorArgs := strings.TrimPrefix(colorArg, "--color=")

	// validate color exists in the Colors map
	if _, ok := ascii.Colors[colorArgs]; !ok {
		fmt.Printf("error: '%s' is not a supported color\n", colorArgs)
		fmt.Println("Supported colors: black, red, green, yellow, blue, magenta, purple, cyan, white, orange")
		return
	}

	var subStr string
	var input string

	// case 1: no substring (color entire text)
	if len(os.Args) == 3 {
		input = os.Args[2]
		subStr = ""
	}

	// case 2: substring provided
	if len(os.Args) == 4 {
		subStr = os.Args[2]
		input = os.Args[3]
	}

	// read banner
	bannerLines, err := ascii.ReadBanner("standard.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	asciiMap := ascii.BuildAsciiMap(bannerLines)

	result := ascii.PrintAscii(input, asciiMap, subStr, colorArgs)

	fmt.Print(result)
}
