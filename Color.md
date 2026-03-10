# ASCII Art Color — Complete Project Guide

A Go command-line tool that converts any string into large ASCII art using pre-made banner font files, with full ANSI color support for coloring the entire output or only specific substrings.

---

## Table of Contents

1. [What This Project Does](#what-this-project-does)
2. [How It Works — The Core Concept](#how-it-works)
3. [Project Structure](#project-structure)
4. [How to Run](#how-to-run)
5. [The Code — File by File, Line by Line](#the-code)
   - [go.mod](#gomod)
   - [ascii/render.go — ReadBanner](#readbanner)
   - [ascii/render.go — BuildAsciiMap](#buildasciimap)
   - [ascii/render.go — PrintAscii](#printascii)
   - [ascii/render.go — Colors Map](#colors-map)
   - [ascii/render.go — PrintAsciiColor](#printasciicolor)
   - [main.go](#maingo)
   - [main_test.go](#main_testgo)
6. [The Color Logic Explained](#the-color-logic-explained)
7. [All Supported Commands](#all-supported-commands)
8. [Audit Test Cases](#audit-test-cases)
9. [Common Bugs and Fixes](#common-bugs-and-fixes)
10. [Submission Checklist](#submission-checklist)

---

## What This Project Does

The base project converts a string into large ASCII art printed in the terminal. The color extension builds on top of that and adds the ability to colorize the entire output or only specific substrings using ANSI terminal escape codes — the same mechanism your terminal uses to display colored text in git output, compiler errors, and shell prompts.

```bash
go run . "Hi"
```

```
 _    _   _
| |  | | (_)
| |__| |  _
|  __  | | |
| |  | | | |
|_|  |_| |_|


```

```bash
go run . --color=orange GuYs "HeY GuYs"
```

Output: `HeY ` in the terminal's default color, `GuYs` in orange — across all 8 rows of the ASCII art simultaneously.

---

## How It Works

The program has four steps:

**Step 1 — Read the banner file**
`ReadBanner` opens a `.txt` file (e.g. `standard.txt`) and loads every line into a slice of strings.

**Step 2 — Build the character map**
`BuildAsciiMap` takes those lines and organizes them into a lookup table: every character (`A`–`Z`, `a`–`z`, `0`–`9`, symbols, space) maps to its 8-line art representation.

**Step 3 — Parse the command-line arguments**
`main.go` reads the flags and positional arguments to determine: which color was requested (if any), which substring to color (if any), which banner font to use, and what string to render.

**Step 4 — Print the art**
If a color was given, `PrintAsciiColor` renders each character's art and wraps matching positions in ANSI escape codes. Otherwise the plain `PrintAscii` is used.

---

## Project Structure

```
ascii-art/
├── main.go              ← entry point and argument parsing
├── go.mod               ← module: ascii-art, go 1.22.2
├── main_test.go         ← 5 unit tests
├── ascii/
│   └── render.go        ← ReadBanner, BuildAsciiMap, PrintAscii, Colors, PrintAsciiColor
├── standard.txt         ← standard banner font
├── shadow.txt           ← shadow banner font
└── thinkertoy.txt       ← thinkertoy banner font
```

> The `.txt` banner files must be in the project root — same level as `main.go`. They must NOT be inside `ascii/`. This is because `os.ReadFile` resolves paths relative to where you run `go run .`

---

## How to Run

```bash
# Basic — whole string, standard banner
go run . "hello"

# Choose a banner
go run . "hello" shadow
go run . "hello" thinkertoy

# Color the whole string
go run . --color=red "hello"

# Color only a substring
go run . --color=orange GuYs "HeY GuYs"

# Color + choose banner
go run . --color=cyan "hello" shadow

# Color substring + choose banner
go run . --color=blue kit "a king kitten have kit" standard

# Newline between art blocks
go run . "Hello\nThere"

# Run all tests
go test ./... -v
```

---

## The Code — File by File, Line by Line

### go.mod

Created automatically by `go mod init ascii-art`. It looks like:

```
module ascii-art

go 1.22.2
```

The module name `ascii-art` is what makes the import path `"ascii-art/ascii"` work in `main.go`. The module name in this file and the import path in every file that imports the package must match exactly.

---

### ascii/render.go

This single file lives inside the `ascii/` subfolder and contains all core logic. The full file is:

```go
package ascii

import (
    "fmt"
    "os"
    "strings"
)

func ReadBanner(file string) ([]string, error) { ... }
func BuildAsciiMap(lines []string) map[rune][]string { ... }
func PrintAscii(text string, asciiMap map[rune][]string) { ... }
var Colors = map[string]string{ ... }
func PrintAsciiColor(text, color, substr string, asciiMap map[rune][]string) { ... }
```

Because it lives in a subfolder with its own package name `ascii`, `main.go` must explicitly import it with `import "ascii-art/ascii"` and call all its functions as `ascii.FunctionName(...)`.

---

### ReadBanner

```go
func ReadBanner(file string) ([]string, error) {
    data, err := os.ReadFile(file)
    if err != nil {
        return nil, err
    }
    lines := strings.Split(string(data), "\n")
    return lines, nil
}
```

**What it does:** Opens a banner `.txt` file and returns every line as a slice of strings.

---

**`func ReadBanner(file string) ([]string, error)`**

Defines a function that takes one input (the filename as a string) and returns two values simultaneously — Go allows multiple return values:

- `[]string` — a slice of strings (every line of the file)
- `error` — either `nil` (no problem) or a description of what went wrong

Returning an error alongside the result is Go's standard pattern for operations that can fail, like reading a file from disk.

---

**`data, err := os.ReadFile(file)`**

`os.ReadFile` opens the file, reads every byte, closes it, and returns two values at once using `:=` (short variable declaration — creates both variables and infers their types automatically):

- `data` — a `[]byte` containing the raw contents of the file
- `err` — `nil` on success, or an error describing what went wrong

---

**`if err != nil { return nil, err }`**

`nil` in Go means "nothing" or "no value". If `err` is not nil, the file couldn't be read. We return early with `nil` for the lines slice (nothing useful to give back) and pass the error up so the caller — `main.go` — can decide what to do. This early-return pattern is idiomatic Go.

---

**`lines := strings.Split(string(data), "\n")`**

Two operations chained:

- `string(data)` converts the raw `[]byte` into a readable Go string
- `strings.Split(..., "\n")` cuts the string into a slice at every newline character

After this, `lines[0]` is the first line of the file, `lines[1]` is the second, and so on. The entire banner file is now accessible line by line.

---

**`return lines, nil`**

Return the completed slice and `nil` for the error — everything worked.

---

### BuildAsciiMap

```go
func BuildAsciiMap(lines []string) map[rune][]string {
    asciiMap := make(map[rune][]string)
    char := 32
    for i := 0; i < len(lines); i += 9 {
        asciiMap[rune(char)] = lines[i : i+8]
        char++
    }
    return asciiMap
}
```

**What it does:** Takes the raw lines from `ReadBanner` and organizes them into a lookup table — given any character, instantly get its 8 art lines back.

---

**`func BuildAsciiMap(lines []string) map[rune][]string`**

Returns a `map[rune][]string`.

A **map** is a lookup table — like a dictionary. You give it a key, it gives you a value instantly.

- Key type: `rune` — Go's type for a single Unicode character. A rune is an `int32` storing the character's code point. `'A'` is rune `65`, `' '` is rune `32`, `'!'` is rune `33`.
- Value type: `[]string` — a slice of 8 strings (the 8 art rows for that character)

So `asciiMap['A']` returns the 8 strings that together draw the letter A.

---

**`asciiMap := make(map[rune][]string)`**

`make` creates an empty but usable map. Writing `var asciiMap map[rune][]string` would create a `nil` map that panics on first write. Think of `make` as building the empty filing cabinet before you start filing things.

---

**`char := 32`**

Start at ASCII code 32, the space character. This matches where the banner files begin. `char` is a plain integer that increments by 1 per loop iteration, tracking which character we're currently assigning art to.

---

**`for i := 0; i < len(lines); i += 9`**

The critical loop. `i` starts at 0 and advances by 9 each iteration. Why 9? Because each character in the banner file occupies exactly 9 lines: 8 lines of actual art followed by 1 blank separator line.

- `i = 0` → space (ASCII 32)
- `i = 9` → `!` (ASCII 33)
- `i = 18` → `"` (ASCII 34)
- ... and so on for all 95 printable characters

---

**`asciiMap[rune(char)] = lines[i : i+8]`**

The most important line in `BuildAsciiMap`. Two things happening:

`rune(char)` converts the integer `char` (e.g. `65`) to a rune (the actual character `'A'`). This becomes the map key.

`lines[i : i+8]` is a slice expression. It extracts elements from index `i` up to (but not including) `i+8` — exactly 8 elements, the 8 art rows for this character. The 9th line (the blank separator) is automatically skipped because the loop steps to `i+9` next.

Together: _"Store the 8 art lines starting at position i, filed under this character's key."_

---

**`char++`**

Move to the next ASCII character. After space (32) comes `!` (33), then `"` (34), and so on. This keeps `char` in sync with which character's block `i` has reached in the file.

---

### PrintAscii

```go
func PrintAscii(text string, asciiMap map[rune][]string) {
    lines := strings.Split(text, "\\n")
    for _, line := range lines {
        if line == "" {
            fmt.Println()
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
```

**What it does:** Takes the user's input string, looks up each character in the map, and prints the ASCII art side by side in the terminal. Used when no color flag is provided.

---

**The core rendering problem — why you can't print one character at a time:**

Your first instinct might be to loop through each character, print all 8 of its rows, then move to the next character. But that stacks them vertically:

```
 _    _
| |  | |     ← H's 8 rows printed, then...
...
 _
(_)          ← i's 8 rows printed below
```

That's wrong. The correct output is side by side:

```
 _    _   _      ← row 0 of H, then row 0 of i, on the same terminal line
| |  | | (_)     ← row 1 of H, then row 1 of i
...
```

To achieve this you must think in **rows first, characters second**. For each row (0 through 7), print that row's art for every character in the input, then print a newline. This is the double-loop pattern at the heart of `PrintAscii` — and `PrintAsciiColor`.

---

**`lines := strings.Split(text, "\\n")`**

Splits the input on the literal two-character sequence backslash-then-n.

When a user types `go run . "Hello\nThere"`, Go receives the 12-character string `Hello\nThere` — a literal backslash followed by `n`. It is NOT a real newline character. The shell passes it as-is.

In a Go string literal, `"\\n"` means: a backslash character followed by `n`. This correctly matches the user's `\n` and splits on it.

```go
strings.Split(text, "\n")   // WRONG — splits on actual newline (1 byte)
strings.Split(text, "\\n")  // CORRECT — splits on backslash + n (2 bytes)
```

Result: `"Hello\nThere"` becomes `["Hello", "There"]`. `"Hello\n\nThere"` becomes `["Hello", "", "There"]`.

---

**`if line == "" { fmt.Println(); continue }`**

An empty `line` means the user typed `\n\n` — a double newline, producing a gap between two art blocks. `fmt.Println()` with no arguments prints one blank line. `continue` skips all rendering for this iteration and jumps to the next segment.

---

**`for row := 0; row < 8; row++`**

The outer rendering loop. Runs exactly 8 times — once per row of ASCII art height. This drives which row we are currently building across all characters.

---

**`for _, char := range line`**

The inner loop. `range` on a string yields the byte index (discarded with `_`) and each character as a `rune`. For `"Hi"`: first `char = 'H'`, then `char = 'i'`.

---

**`fmt.Print(asciiMap[char][row])`**

Three things chained:

- `asciiMap[char]` — looks up this character, gets back its 8-element `[]string`
- `[row]` — indexes into that slice to get the single string for the current row
- `fmt.Print(...)` — prints it **without** a newline

Using `Print` instead of `Println` is critical. Because no newline is added, the next character's row is printed immediately after on the same terminal line — this is what produces the side-by-side effect.

---

**`fmt.Println()`** (after the inner loop)

After every character has contributed its piece of row `row`, this moves the cursor to the next terminal line. After all 8 rows are done, the full art block for this segment is complete.

---

### Colors Map

```go
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
```

**What it is:** A package-level variable (declared with `var`, not inside any function) that maps human-readable color names to their ANSI terminal escape codes.

---

**`var Colors = map[string]string{...}`**

`var` at package level means `Colors` is accessible from anywhere inside the `ascii` package, and — because it starts with a capital letter — also from outside the package. `main.go` reads `ascii.Colors` directly to validate the user's color argument.

A `map[string]string` maps string keys to string values. Each key is a color name the user types on the command line. Each value is the ANSI escape code sequence for that color.

---

**How ANSI escape codes work:**

An ANSI escape code is a special character sequence that terminals interpret as a formatting instruction instead of text to display.

```
\033[31m   →   ESC [ 31 m   →   start red text
\033[0m    →   ESC [ 0  m   →   reset to terminal default
```

`\033` is the escape character (octal 33, decimal 27). It is followed by `[`, then a number or series of numbers, then `m` (the "Select Graphic Rendition" command). The number selects the color or effect.

Wrapping text between a color code and a reset code:

```
\033[31mhello\033[0m
```

causes the terminal to display `hello` in red. Without the reset, all subsequent text would also be red.

---

**`"30"`–`"37"` — standard ANSI colors:**

The numbers 30 through 37 are reserved by the ANSI standard for the 8 basic foreground colors. These work in virtually every terminal on every platform. The sequence `\033[30m` through `\033[37m` maps to:

- 30 → black
- 31 → red
- 32 → green
- 33 → yellow
- 34 → blue
- 35 → magenta
- 36 → cyan
- 37 → white

---

**`"38;5;214"` — 256-color mode for orange:**

There is no standard ANSI code for orange — the original 8 colors jump from yellow (33) straight to blue (34). To get orange, we use the 256-color extension: `38;5;N` where `N` is a number from 0 to 255 selecting a specific color from the extended palette. Color 214 is a bright orange in this palette. Not all terminals support 256-color mode, but most modern terminals (iTerm2, GNOME Terminal, Windows Terminal) do.

---

**`"reset": "\033[0m"`**

Color 0 in the ANSI standard is a special value meaning "reset all attributes to their defaults." After printing colored art, this code immediately restores the terminal to its normal state so the next character (if not colored) appears in the user's default terminal color.

---

**Why `Colors` is exported (capital C):**

Go's visibility rule: identifiers starting with a capital letter are exported from their package and accessible from other packages. `Colors` is exported because `main.go` needs to reach into the `ascii` package to validate user input:

```go
if _, ok := ascii.Colors[color]; !ok {
    fmt.Println("Usage: go run . [--color=COLOR] [SUBSTRING] STRING [BANNER]")
    return
}
```

If the user types `--color=purple` and `"purple"` is not a key in the map, this check catches it and prints the usage message before any rendering happens.

---

### PrintAsciiColor

```go
func PrintAsciiColor(text, color, substr string, asciiMap map[rune][]string) {
    colorCode := Colors[color]
    reset := Colors["reset"]

    lines := strings.Split(text, "\\n")
    for _, line := range lines {
        if line == "" {
            fmt.Println()
            continue
        }

        runes := []rune(line)
        subRunes := []rune(substr)
        colored := make([]bool, len(runes))

        if substr == "" {
            for i := range colored {
                colored[i] = true
            }
        } else {
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
```

**What it does:** Extends `PrintAscii` with two new responsibilities — deciding which character positions in the input should be colored, then wrapping those positions' art rows in ANSI escape codes when rendering.

---

**`func PrintAsciiColor(text, color, substr string, asciiMap map[rune][]string)`**

Four parameters:

- `text` — the full string to render (e.g. `"HeY GuYs"`)
- `color` — the color name chosen by the user (e.g. `"orange"`)
- `substr` — the substring to colorize (e.g. `"GuYs"`), or `""` to color everything
- `asciiMap` — the lookup table built by `BuildAsciiMap`

When three string parameters share a type, Go lets you write them together: `text, color, substr string` instead of `text string, color string, substr string`. The result is identical.

---

**`colorCode := Colors[color]`**
**`reset := Colors["reset"]`**

Look up the two ANSI codes we will use for the rest of this function:

- `colorCode` — the escape sequence to start the chosen color (e.g. `"\033[38;5;214m"` for orange)
- `reset` — `"\033[0m"`, the sequence to cancel the color after each colored art row

Fetching these once into local variables means we're not doing a map lookup inside every iteration of the innermost loop.

---

**`lines := strings.Split(text, "\\n")`**

Same as in `PrintAscii`. Splits on the literal backslash-n sequence to support multi-line art blocks. See the `PrintAscii` section above for a full explanation.

---

**`for _, line := range lines`**
**`if line == "" { fmt.Println(); continue }`**

Same as `PrintAscii`. An empty segment means a `\n\n` in the input — print a blank line and skip to the next segment.

---

**`runes := []rune(line)`**

Converts the string `line` into a slice of runes — `[]rune`.

In Go, a `string` is a sequence of bytes (UTF-8 encoded). When you `range` over a string, Go automatically decodes each multi-byte character into a `rune`. But if you need to index a string by character position (not byte position), you must first convert it to `[]rune`.

This matters for the `colored` array: `colored[i]` must align perfectly with the `i`-th character in the rendering loop. If we indexed into the string as bytes, multi-byte characters would cause `i` to be off. Converting to `[]rune` first guarantees that index `i` in `colored` corresponds to the exact same `i`-th character in the render loop below.

---

**`subRunes := []rune(substr)`**

Same conversion for the substring. We need to compare it character by character against windows of `runes`, so it also needs to be a `[]rune`.

---

**`colored := make([]bool, len(runes))`**

Creates a boolean slice with one entry per character in the input line. Every entry starts as `false` (Go's zero value for `bool`). An entry at position `i` being `true` means: _"when rendering character `i`, wrap its art row in the color code."_

This is the key data structure of the function. Rather than making a coloring decision mid-render, we precompute the entire map of which positions are colored before touching `asciiMap` at all. The rendering loop below simply reads from this array.

---

**The two cases — empty substr vs. specific substr:**

**Case 1: `substr == ""`**

```go
if substr == "" {
    for i := range colored {
        colored[i] = true
    }
}
```

When the user provides no substring (e.g. `go run . --color=red "hello"`), the intent is to color the entire output. Setting every entry in `colored` to `true` achieves that: every character will be colored.

---

**Case 2: substr is a specific string**

```go
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
```

This is a **sliding window** substring search. Rather than using `strings.Contains`, we scan position by position through `runes` and check whether `subRunes` starts at each position. If it does, every position in that window is marked as `true` in `colored`.

**Why not `strings.Contains` or `strings.ContainsRune`?**

An earlier approach checked:

```go
strings.ContainsRune(substr, char)
```

This asks: _"does the character appear anywhere inside the substr string?"_ — not _"is the character at this position part of an occurrence of substr in the input?"_

For `--color=orange GuYs "HeY GuYs"`, the character `Y` appears inside `"GuYs"`. So `strings.ContainsRune("GuYs", 'Y')` returns `true`. That means the `Y` in `HeY` would also get colored — which is wrong. We want to color only the characters that are actually part of a `"GuYs"` match in the input, not every character that happens to appear in the substr string somewhere.

The sliding window fixes this by tracking _positions_, not individual characters.

---

**Walking through the sliding window for `--color=orange GuYs "HeY GuYs"`:**

```
runes    = ['H','e','Y',' ','G','u','Y','s']   (len = 8)
subRunes = ['G','u','Y','s']                   (len = 4)
colored  = [F, F, F, F, F, F, F, F]           (all false initially)

Loop range: i = 0 to len(runes)-len(subRunes) = 8-4 = 4, so i goes 0,1,2,3,4

i=0: runes[0]='H' vs subRunes[0]='G' → no match, inner loop breaks immediately
i=1: runes[1]='e' vs subRunes[0]='G' → no match
i=2: runes[2]='Y' vs subRunes[0]='G' → no match
i=3: runes[3]=' ' vs subRunes[0]='G' → no match
i=4: runes[4]='G' vs subRunes[0]='G' → continue
     runes[5]='u' vs subRunes[1]='u' → continue
     runes[6]='Y' vs subRunes[2]='Y' → continue
     runes[7]='s' vs subRunes[3]='s' → match = true!
     → mark colored[4], colored[5], colored[6], colored[7] = true

Final: colored = [F, F, F, F, T, T, T, T]
```

Result: `H`, `e`, `Y`, ` ` render in default color. `G`, `u`, `Y`, `s` render in orange — even though `Y` appears in both words, only the `Y` at position 6 (inside `GuYs`) gets colored.

**The loop upper bound `i <= len(runes)-len(subRunes)`:**

This prevents the inner loop from reading past the end of `runes`. If `runes` has 8 characters and `subRunes` has 4, the last possible match start is position 4 (positions 4,5,6,7 cover the last 4 characters). Checking position 5 would need positions 5,6,7,8 — but index 8 is out of bounds. The bound `len(runes)-len(subRunes)` computes exactly where to stop.

---

**The rendering loop:**

```go
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
```

This is identical in structure to `PrintAscii`'s double loop — rows first, characters second. The only difference is in the inner body:

**`for i, char := range runes`**

Notice `i, char` instead of `_, char`. Here we keep the index `i` because we need to check `colored[i]` for each character. Since we converted `line` to `[]rune` before the loop, this `i` is a character-position index, perfectly aligned with the `colored` array.

**`artRow := asciiMap[char][row]`**

Look up this character's art string for the current row. Storing it in a local variable `artRow` avoids the double map lookup in the if/else below.

**`if colored[i] { fmt.Print(colorCode + artRow + reset) } else { fmt.Print(artRow) }`**

If position `i` is marked for coloring, wrap the art row between the ANSI color code and the reset code. Otherwise, print the art row as plain text.

String concatenation with `+` is used here. Each of the three strings — `colorCode`, `artRow`, `reset` — is small, so this is perfectly efficient. The result is a single string like `"\033[38;5;214m _    _ \033[0m"` which the terminal interprets as orange `_    _` followed by a color reset.

**Why reset after every art row?**

The reset is applied immediately after each colored art string — not once at the end of a line. This means uncolored characters printed on the same terminal line (those where `colored[i]` is `false`) are completely unaffected. Without resetting after each colored segment, the color would "bleed" into the characters that follow it on the same row.

**`fmt.Println()`**

After all characters for this row have been printed, move to the next terminal line. Same as in `PrintAscii`.

---

### main.go

```go
package main

import (
    "fmt"
    "os"
    "strings"

    "ascii-art/ascii"
)

func main() {
    args := os.Args[1:]

    color := ""
    remaining := []string{}
    for _, a := range args {
        if strings.HasPrefix(a, "--color=") {
            color = strings.TrimPrefix(a, "--color=")
        } else {
            remaining = append(remaining, a)
        }
    }

    if color != "" {
        if _, ok := ascii.Colors[color]; !ok {
            fmt.Println("Usage: go run . [--color=COLOR] [SUBSTRING] STRING [BANNER]")
            return
        }
    }

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

    bannerLines, err := ascii.ReadBanner(banner + ".txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    asciiMap := ascii.BuildAsciiMap(bannerLines)

    if color != "" {
        ascii.PrintAsciiColor(input, color, substr, asciiMap)
    } else {
        ascii.PrintAscii(input, asciiMap)
    }
}
```

**What it does:** Parses all command-line arguments, validates the color, decides which combination of arguments was given, loads the correct banner file, and calls the appropriate render function.

---

**`package main`**

Every Go program needs exactly one `main` package with one `main()` function. This is where execution begins when you run `go run .`

---

**`import "strings"`**

The `strings` package is new in `main.go` for the color extension. It provides `strings.HasPrefix` (to detect `--color=...`) and `strings.TrimPrefix` (to extract the value after `--color=`).

---

**`import "ascii-art/ascii"`**

Imports the `ascii` package from the `ascii/` subfolder. The path is `<module-name>/<folder-name>`. Functions and variables in the package are accessed as `ascii.ReadBanner`, `ascii.Colors`, etc. Only identifiers starting with a capital letter are exported — which is why all functions and the `Colors` map are capitalized.

---

**`args := os.Args[1:]`**

`os.Args` is a slice containing everything typed on the command line:

- `os.Args[0]` → the program name (e.g. `"./ascii-art"`)
- `os.Args[1]` onward → the actual arguments the user typed

`os.Args[1:]` is a slice expression that gives us everything from index 1 to the end — all user-provided arguments, with the program name removed. Working with `args` instead of `os.Args` directly makes the rest of the parsing cleaner.

---

**Flag extraction loop:**

```go
color := ""
remaining := []string{}
for _, a := range args {
    if strings.HasPrefix(a, "--color=") {
        color = strings.TrimPrefix(a, "--color=")
    } else {
        remaining = append(remaining, a)
    }
}
```

This loop separates the `--color=VALUE` flag from the positional arguments in one pass.

`strings.HasPrefix(a, "--color=")` returns `true` if the argument starts with `"--color="`. If it does, `strings.TrimPrefix(a, "--color=")` removes that prefix and leaves just the color value — e.g. `"--color=orange"` becomes `"orange"`.

Arguments that are not the color flag are collected into `remaining`. After the loop:

- `color` holds the user's color choice, or `""` if none was given
- `remaining` holds only the positional arguments (substring, string, banner)

This separation is important because it makes the `switch len(remaining)` below predictable — it counts only positional arguments, regardless of how many flags were given.

---

**Color validation:**

```go
if color != "" {
    if _, ok := ascii.Colors[color]; !ok {
        fmt.Println("Usage: go run . [--color=COLOR] [SUBSTRING] STRING [BANNER]")
        return
    }
}
```

If the user specified a color, check that it exists as a key in the `Colors` map. The two-value map lookup `_, ok := ascii.Colors[color]` is Go's idiomatic existence check: `ok` is `true` if the key was found, `false` if not. We don't need the value itself (just confirming existence), so it's discarded with `_`.

If the color is unknown — `"purple"`, `"pink"`, a typo — the program prints the usage message and exits before doing any work. This also catches the common mistake of writing `--color red` without the `=` sign: Go's shell tokenization splits that into two arguments `"--color"` and `"red"`. `"--color"` does not have the prefix `"--color="` so `color` stays `""` and `"red"` ends up in `remaining`, which then causes a usage error in the switch below.

---

**`knownBanners := map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}`**

A local map used to test whether an argument is a recognized banner name. This is needed to disambiguate the 2-argument case: `go run . "hello" shadow` has the banner as the second argument, while `go run . --color=red GuYs "HeY GuYs"` has a substring as the first and the string as the second.

Looking up a key in this map returns `true` if the banner is known, `false` (the zero value for `bool`) if it is not — so `knownBanners["shadow"]` is `true` and `knownBanners["hello"]` is `false`.

---

**`input, substr, banner := "", "", "standard"`**

Declare the three variables that will hold the parsed positional arguments, with their defaults. `banner` defaults to `"standard"` — if the user doesn't specify a banner, we use `standard.txt`.

---

**`switch len(remaining)`**

Handles every valid combination of positional arguments after the flag has been stripped:

**`case 1:`** One positional argument — it must be the string to render.

```
go run . "hello"
go run . --color=red "hello"
```

**`case 2:`** Two positional arguments — could be:

- `STRING BANNER` (always valid)
- `SUBSTRING STRING` (only valid when `--color` was given, because substrings only make sense with a color)

```go
if knownBanners[remaining[1]] {
    input, banner = remaining[0], remaining[1]
} else if color != "" {
    substr, input = remaining[0], remaining[1]
} else {
    // usage error
}
```

We check the second argument first: if it is a known banner name, treat the first as the string and second as the banner. Otherwise, if a color flag was given, treat the first as the substring and the second as the string to render. If neither condition holds (no color flag, and second arg is not a banner), print a usage message.

**`case 3:`** Three positional arguments — `SUBSTRING STRING BANNER`. This is only valid when `--color` was given, and the third argument must be a known banner.

```
go run . --color=blue kit "a king kitten" standard
```

**`default:`** Zero arguments, or four or more — always invalid. Print usage.

---

**`bannerLines, err := ascii.ReadBanner(banner + ".txt")`**

Constructs the filename by appending `".txt"` to the banner name, then calls `ReadBanner`. If the file can't be read (e.g. a missing `.txt` file), the error is printed and the program exits.

---

**`if color != "" { ascii.PrintAsciiColor(...) } else { ascii.PrintAscii(...) }`**

The final decision. If a color was provided, use the color renderer. Otherwise, use the plain renderer. This keeps the base behavior — no color, same output as before — completely unchanged for commands that don't use `--color`.

---

### main_test.go

```go
package main

import (
    "strings"
    "testing"
    "ascii-art/ascii"
)

func TestReadBanner(t *testing.T) { ... }
func TestReadBannerMissingFile(t *testing.T) { ... }
func TestBuildAsciiMap(t *testing.T) { ... }
func TestMapContainsExpectedChars(t *testing.T) { ... }
func TestUpperLowerDifferent(t *testing.T) { ... }
```

Tests live in `package main` — the same package as `main.go` — and are run with:

```bash
go test ./... -v
```

The tests do not test `PrintAscii` or `PrintAsciiColor` directly because both functions write straight to the terminal and return nothing. Testing terminal output requires capturing stdout, which adds complexity. Instead, the tests focus on the two functions that return values: `ReadBanner` and `BuildAsciiMap`.

---

**TestReadBanner** — Confirms all three banner files load successfully and produce non-empty line slices. Fails immediately if any file is missing or unreadable.

**TestReadBannerMissingFile** — Confirms that `ReadBanner` returns a non-nil error when the file does not exist. This tests the error-handling path rather than the success path.

**TestBuildAsciiMap** — Confirms the map has exactly 95 entries (all printable ASCII characters from space to `~`) and that every entry has exactly 8 art lines. This directly catches the most common bug: a wrong loop step in `BuildAsciiMap`.

**TestMapContainsExpectedChars** — Spot-checks that specific characters (`A`, `a`, `1`, `!`, space, etc.) are present and accessible by their rune values.

**TestUpperLowerDifferent** — Confirms `asciiMap['A']` and `asciiMap['a']` contain different art. Catches off-by-one errors in the map building loop where two characters end up sharing the same art block.

---

## The Color Logic Explained

### How ANSI codes work in the terminal

When your terminal receives text to display, it reads the bytes one by one. Most bytes are simply displayed. But when it encounters the escape character (`\033`, decimal 27), it switches into "command mode": the following bytes up to and including `m` are interpreted as a formatting command, not as text.

```
\033[31m   → "start displaying in red"
\033[0m    → "reset to default"
```

In `PrintAsciiColor`, each art row for a colored character is individually wrapped:

```go
fmt.Print(colorCode + artRow + reset)
```

This produces output like:

```
\033[31m _    _ \033[0m\033[31m| |  | |\033[0m ...
```

Which the terminal renders as: the art rows for colored characters in red, and the art rows for uncolored characters in the terminal's default color — even though all of this appears on the same physical line in the output.

---

### Why position-based tracking is correct

The `colored []bool` array marks positions, not character identities. Consider the string `"HeY GuYs"` with substr `"GuYs"`.

Both the `Y` in `HeY` and the `Y` in `GuYs` are the same rune (`'Y'`). A naive approach that checks "is this character in the substr?" would color both `Y`s. The position-based sliding window approach finds that only the `Y` at position 6 is part of a `"GuYs"` match, so only that `Y` gets colored.

This distinction between _"character identity"_ and _"character position in a match"_ is the entire reason `PrintAsciiColor` is more complex than a simple `strings.ContainsRune` check would be.

---

### Multiple occurrences of the substring

The sliding window marks **all** occurrences of `substr`, not just the first. For `--color=cyan kit "a king kitten have kit"`:

```
runes = ['a',' ','k','i','n','g',' ','k','i','t','t','e','n',' ','h','a','v','e',' ','k','i','t']

Match at i=2:  'k','i' — but next char is 'n', not 't' — no match
Match at i=7:  'k','i','t' — match! mark positions 7,8,9
Match at i=19: 'k','i','t' — match! mark positions 19,20,21
```

Both `"kit"` occurrences (in `"king"` there's no `kit`, in `"kitten"` there is, and in `"kit"` at the end there is) end up colored. The word `"king"` is not affected because `'k','i','n'` does not match `'k','i','t'` — the comparison fails at the third character.

---

## All Supported Commands

```bash
go run . "hello"                                        # standard banner, no color
go run . "hello" shadow                                 # shadow banner
go run . "hello" thinkertoy                             # thinkertoy banner
go run . --color=red "hello"                            # whole string in red
go run . --color=blue "hello" shadow                    # whole string in blue, shadow banner
go run . --color=orange GuYs "HeY GuYs"                 # only GuYs in orange
go run . --color=cyan kit "a king kitten have kit"      # all occurrences of kit in cyan
go run . --color=green B "RGB()"                        # only B characters in green
go run . "Hello\nThere"                                 # two separate art blocks
go run . "Hello\n\nThere"                               # two blocks with blank line between
go run . --color=magenta GuYs "HeY GuYs" shadow         # substring colored, shadow font
go run . --color=yellow kit "a king kitten" thinkertoy  # substring colored, thinkertoy font
```

**Available colors:** `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`, `orange`

---

## Audit Test Cases

| Command                                     | Expected result                                |
| ------------------------------------------- | ---------------------------------------------- |
| `go run . --color red "hello"`              | Missing `=` → usage message                    |
| `go run . --color=red "hello world"`        | Entire output in red                           |
| `go run . --color=green "1 + 1 = 2"`        | Entire output in green including special chars |
| `go run . --color=orange GuYs "HeY GuYs"`   | Only `GuYs` in orange, `HeY ` in default       |
| `go run . --color=blue B "RGB()"`           | Only `B` characters in blue                    |
| `go run . "hello" shadow`                   | Rendered in shadow banner font                 |
| `go run . --color=cyan kit "a king kitten"` | All occurrences of `kit` in cyan               |
| `go run . ""`                               | No output                                      |

---

## Common Bugs and Fixes

**Wrong color applied — characters outside the substring get colored**

Caused by `strings.ContainsRune(substr, char)` which checks whether the character appears anywhere in the substr string, not whether this character's position in the input is inside an actual occurrence of the substr. Fix: use the sliding window position-marking approach in `PrintAsciiColor`.

**`--color red "hello"` not caught as wrong format**

When `=` is missing, the shell gives Go two separate arguments: `"--color"` and `"red"`. `"--color"` does not match `strings.HasPrefix(a, "--color=")`, so `color` stays `""`. Then `"red"` ends up in `remaining` as an unexpected positional argument, which triggers the `default` usage error in the switch. The format error is caught correctly — just not with a specific error message about the `=` sign.

**Characters shifted — `A` renders as `B`**

The loop step in `BuildAsciiMap` must be `i += 9`. Each character occupies exactly 9 lines (8 art lines + 1 blank separator). Using `i += 8` causes a 1-line drift per character, shifting every subsequent character's art by one line.

**`\n` in input not splitting into separate art blocks**

The split must target the two-character literal sequence `\` then `n`, not a real newline:

```go
strings.Split(text, "\\n")   // correct
strings.Split(text, "\n")    // wrong — matches a real newline character
```

**Banner file not found**

The `.txt` files must be in the project root alongside `main.go`, not inside `ascii/`. `os.ReadFile` resolves paths relative to the working directory where you run `go run .`

**Color bleeds into uncolored characters**

Make sure the reset code `Colors["reset"]` is appended immediately after every colored art row: `colorCode + artRow + reset`. Without the reset, the terminal continues rendering in the selected color until it receives another escape sequence.

**Panic: index out of range in sliding window**

The upper bound of the outer loop must be `i <= len(runes)-len(subRunes)`, not `i < len(runes)`. Without this, when `i` is near the end of `runes`, the inner loop `runes[i+j]` reads past the end of the slice.

---

## Submission Checklist

**Functionality**

- [ ] `go run . "hello"` renders correctly in standard font
- [ ] `go run . "hello" shadow` uses shadow banner
- [ ] `go run . "hello" thinkertoy` uses thinkertoy banner
- [ ] `go run . --color=red "hello"` colors the whole output red
- [ ] `go run . --color=orange GuYs "HeY GuYs"` colors only `GuYs`
- [ ] `go run . --color=cyan kit "a king kitten"` colors all `kit` occurrences
- [ ] `go run . --color red "hello"` prints usage message (missing `=`)
- [ ] `go run . --color=purple "hello"` prints usage message (unknown color)
- [ ] `go run . "Hello\nThere"` produces two separate art blocks
- [ ] `go run . "Hello\n\nThere"` has a blank line between the blocks
- [ ] `go run . ""` produces no output
- [ ] Numbers, symbols, uppercase, lowercase all render correctly

**Code**

- [ ] `ascii/render.go` starts with `package ascii`
- [ ] `main.go` starts with `package main`
- [ ] No external packages — standard library only
- [ ] `go build .` compiles with zero errors

**Tests**

- [ ] `go test ./... -v` passes all 5 tests
- [ ] Tests cover: all 3 banner files, missing file error, 95 chars in map, 8 lines each, key characters, upper vs lower

**Files**

- [ ] `main.go`, `go.mod`, `main_test.go` in root
- [ ] `ascii/render.go` in `ascii/` subfolder
- [ ] `standard.txt`, `shadow.txt`, `thinkertoy.txt` in root

---

_Only the Go standard library is used. No third-party packages._
