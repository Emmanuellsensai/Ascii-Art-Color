package main

import (
	"strings"
	"testing"

	"ascii-art-color/ascii"
)

// Test that all three banner files load without error
func TestReadBanner(t *testing.T) {
	banners := []string{"standard.txt", "shadow.txt", "thinkertoy.txt"}
	for _, b := range banners {
		lines, err := ascii.ReadBanner(b)
		if err != nil {
			t.Errorf("Could not read %s: %v", b, err)
		}
		if len(lines) == 0 {
			t.Errorf("%s loaded but returned no lines", b)
		}
	}
}

// Test that a missing file returns an error
func TestReadBannerMissingFile(t *testing.T) {
	_, err := ascii.ReadBanner("missing.txt")
	if err == nil {
		t.Error("Expected an error for missing file, got nil")
	}
}

// Test that the map has exactly 95 characters and each has 8 lines
func TestBuildAsciiMap(t *testing.T) {
	lines, err := ascii.ReadBanner("standard.txt")
	if err != nil {
		t.Fatal(err)
	}

	asciiMap := ascii.BuildAsciiMap(lines)

	if len(asciiMap) != 95 {
		t.Errorf("Expected 95 characters in map, got %d", len(asciiMap))
	}

	for char, art := range asciiMap {
		if len(art) != 8 {
			t.Errorf("Character %q has %d lines, expected 8", char, len(art))
		}
	}
}

// Test that key characters exist in the map
func TestMapContainsExpectedChars(t *testing.T) {
	lines, _ := ascii.ReadBanner("standard.txt")
	asciiMap := ascii.BuildAsciiMap(lines)

	testChars := "ABCabc123 !@#"
	for _, char := range testChars {
		if _, ok := asciiMap[char]; !ok {
			t.Errorf("Character %q not found in map", char)
		}
	}
}

// Test that uppercase and lowercase have different art
func TestUpperLowerDifferent(t *testing.T) {
	lines, _ := ascii.ReadBanner("standard.txt")
	asciiMap := ascii.BuildAsciiMap(lines)

	upper := strings.Join(asciiMap['A'], "")
	lower := strings.Join(asciiMap['a'], "")

	if upper == lower {
		t.Error("Expected uppercase A and lowercase a to have different art")
	}
}