package main

import (
	"fmt"
	"os"

	"ascii-art/ascii"
)

func main() {

	/*
		This function reads the user input from os.Args and sends it to the Render
		function which handles the ASCII art generation.
	*/
	if len(os.Args) != 2 {
		return
	}

	input := os.Args[1]

	bannerLines, err := ascii.ReadBanner("standard.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	asciiMap := ascii.BuildAsciiMap(bannerLines)

	result := ascii.PrintAscii(input, asciiMap)

	fmt.Print(result)
}
