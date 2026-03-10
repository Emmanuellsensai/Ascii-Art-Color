# main_test.go — Test File Guide

This file tests the two core functions of the ASCII Art project: `ReadBanner` and `BuildAsciiMap`. It does not test `PrintAscii` directly because that function prints straight to the terminal and returns nothing — there is nothing to check against.

---

## How to Run the Tests

From inside your project folder:

```bash
# Run all tests
go test ./...

# Run all tests with details (shows each test name passing or failing)
go test ./... -v

# Run one specific test by name
go test -run TestReadBanner -v
```

A passing run looks like this:
```
--- PASS: TestReadBanner (0.00s)
--- PASS: TestReadBannerMissingFile (0.00s)
--- PASS: TestBuildAsciiMap (0.00s)
--- PASS: TestMapContainsExpectedChars (0.00s)
--- PASS: TestUpperLowerDifferent (0.00s)
PASS
ok  	ascii-art	0.004s
```

---

## The Package and Imports

```go
package main
```
The test file belongs to `package main` — the same package as `main.go`. This is required so Go knows these tests belong to this project.

```go
import (
    "strings"
    "testing"

    "ascii-art/ascii"
)
```

Three imports:
- `"testing"` — Go's built-in test package. Provides the `*testing.T` type used in every test function to report failures.
- `"strings"` — used in the last test to join slices of strings together for comparison.
- `"ascii-art/ascii"` — your own `ascii` package from the `ascii/` subfolder. This gives the test file access to `ReadBanner` and `BuildAsciiMap`.

---

## How Go Tests Work

Every test function must follow this exact pattern:
```go
func TestSomeName(t *testing.T) {
    // test code here
}
```

- The function name **must start with `Test`** — Go's test runner only picks up functions with this prefix.
- `t *testing.T` is the test helper passed in automatically. You use it to report failures:
  - `t.Error(...)` — marks the test as failed but continues running the rest of the function
  - `t.Errorf(...)` — same as `t.Error` but lets you format a message with values (like `%s`, `%d`, `%v`)
  - `t.Fatal(...)` — marks the test as failed and **stops immediately** — used when there's no point continuing

---

## Test 1 — TestReadBanner

```go
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
```

**What it tests:** All three banner files (`standard.txt`, `shadow.txt`, `thinkertoy.txt`) can be read successfully.

**Step by step:**

`banners := []string{"standard.txt", "shadow.txt", "thinkertoy.txt"}` — creates a slice with all three filenames so we can test them all in one loop instead of writing three separate tests.

`for _, b := range banners` — loops through each filename. The `_` discards the index since we only need the filename value `b`.

`lines, err := ascii.ReadBanner(b)` — calls your `ReadBanner` function with each filename. It returns the file's lines and an error value.

`if err != nil` — if the error is not nil, something went wrong reading the file. `t.Errorf` reports the failure with the filename and the specific error message. The `%s` and `%v` are format verbs — `%s` prints a string, `%v` prints any value in its default format.

`if len(lines) == 0` — even if there was no error, we double-check the returned slice actually has content. An empty slice would mean the file was read but had nothing in it.

**Why it matters:** If any banner file is missing or unreadable, the whole program breaks. This test catches that immediately.

---

## Test 2 — TestReadBannerMissingFile

```go
func TestReadBannerMissingFile(t *testing.T) {
    _, err := ascii.ReadBanner("missing.txt")
    if err == nil {
        t.Error("Expected an error for missing file, got nil")
    }
}
```

**What it tests:** When you pass a filename that does not exist, `ReadBanner` returns an error instead of silently returning nothing.

**Step by step:**

`_, err := ascii.ReadBanner("missing.txt")` — calls `ReadBanner` with a file that definitely does not exist. The `_` discards the returned lines slice because we don't care about it here — we only care about the error.

`if err == nil` — this condition is the **opposite** of what you normally check. Here we WANT an error. If `err` is nil it means `ReadBanner` claimed everything was fine when it shouldn't have been — that's the failure.

`t.Error(...)` — reports the failure. No formatting needed here since the message is plain text.

**Why it matters:** This tests that your error handling actually works. A function that never returns errors even when things go wrong is dangerous — your program would continue running with bad data instead of stopping cleanly.

---

## Test 3 — TestBuildAsciiMap

```go
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
```

**What it tests:** `BuildAsciiMap` correctly parses the banner file into a map with exactly 95 characters, each having exactly 8 art lines.

**Step by step:**

`lines, err := ascii.ReadBanner("standard.txt")` — load the banner file first. We need real lines to pass to `BuildAsciiMap`.

`if err != nil { t.Fatal(err) }` — if the file can't be read, use `t.Fatal` to stop immediately. There is no point running the rest of this test without valid data — it would just produce confusing errors.

`asciiMap := ascii.BuildAsciiMap(lines)` — build the map from the loaded lines.

`if len(asciiMap) != 95` — there are exactly 95 printable ASCII characters (from space at code 32 to `~` at code 126). The map must have all of them. `%d` is the format verb for integers.

`for char, art := range asciiMap` — loop through every entry in the map to check each one individually. `char` is the rune (character key), `art` is the `[]string` of 8 lines.

`if len(art) != 8` — every character must have exactly 8 lines. If any has fewer or more, the banner file was parsed incorrectly. `%q` prints a character with quotes around it (e.g. `'A'`), which makes error messages easier to read.

**Why it matters:** This is the most thorough test. If the `i += 9` loop step in `BuildAsciiMap` is wrong, characters will be shifted and their art lines will be wrong — this test catches that directly.

---

## Test 4 — TestMapContainsExpectedChars

```go
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
```

**What it tests:** A representative sample of characters — uppercase, lowercase, numbers, space, and symbols — all exist as keys in the map.

**Step by step:**

`lines, _ := ascii.ReadBanner("standard.txt")` — the `_` discards the error this time. We know this file exists from Test 1, so we keep this test short and focused.

`testChars := "ABCabc123 !@#"` — a string covering the main character categories: uppercase letters, lowercase letters, digits, a space, and special characters. Ranging over a string in Go gives you each character as a rune automatically.

`for _, char := range testChars` — loop through each character in the test string.

`if _, ok := asciiMap[char]; !ok` — this is Go's map existence check. When you look up a key in a map, Go returns two values: the value and a boolean `ok`. If `ok` is `false`, the key is not in the map. We discard the value with `_` since we only care whether the key exists.

**Why it matters:** Even if the map has 95 entries, they might be mapped to the wrong keys. This spot-checks that specific well-known characters are findable by their actual character value.

---

## Test 5 — TestUpperLowerDifferent

```go
func TestUpperLowerDifferent(t *testing.T) {
    lines, _ := ascii.ReadBanner("standard.txt")
    asciiMap := ascii.BuildAsciiMap(lines)

    upper := strings.Join(asciiMap['A'], "")
    lower := strings.Join(asciiMap['a'], "")

    if upper == lower {
        t.Error("Expected uppercase A and lowercase a to have different art")
    }
}
```

**What it tests:** The letter `A` and the letter `a` have genuinely different ASCII art — confirming that uppercase and lowercase are stored and retrieved separately.

**Step by step:**

`asciiMap['A']` and `asciiMap['a']` — look up both letters directly using their rune values as keys.

`strings.Join(asciiMap['A'], "")` — `asciiMap['A']` returns a `[]string` of 8 lines. `strings.Join` concatenates all 8 into one single string with no separator. This gives us the complete art for `A` as one string we can compare.

We join all 8 rows instead of comparing just one row because the first row of many characters is blank — both `A` and `a` have an empty string as row 0. Comparing only row 0 would incorrectly say they are the same. Joining everything together means the full shape of the letter is compared.

`if upper == lower` — if both joined strings are identical, the map has stored the same art for both cases — which would be a parsing bug.

**Why it matters:** This catches a specific off-by-one bug in `BuildAsciiMap`. If the loop step or starting position is slightly wrong, uppercase and lowercase letters can end up pointing to the same art block.

---

## Summary Table

| Test | Function tested | What it checks |
|------|----------------|----------------|
| `TestReadBanner` | `ReadBanner` | All 3 banner files load with content |
| `TestReadBannerMissingFile` | `ReadBanner` | Missing file returns an error |
| `TestBuildAsciiMap` | `BuildAsciiMap` | Map has 95 chars, each with 8 lines |
| `TestMapContainsExpectedChars` | `BuildAsciiMap` | Key characters are findable in the map |
| `TestUpperLowerDifferent` | `BuildAsciiMap` | `A` and `a` have different art |