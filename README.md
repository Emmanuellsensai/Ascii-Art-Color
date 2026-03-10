# ASCII Art — Complete Project Guide

A Go program that converts any string into large ASCII art using pre-made banner font files, rendered character by character directly in the terminal.

---

## Table of Contents

1. [What This Project Does](#what-this-project-does)
2. [How It Works — The Core Concept](#how-it-works)
3. [Prerequisites](#prerequisites)
4. [Project Structure](#project-structure)
5. [Setting Up](#setting-up)
6. [The Banner Files Explained](#the-banner-files-explained)
7. [The Code — File by File, Line by Line](#the-code)
   - [go.mod](#gomod)
   - [ascii/render.go — ReadBanner](#readbanner)
   - [ascii/render.go — BuildAsciiMap](#buildasciimap)
   - [ascii/render.go — PrintAscii](#printascii)
   - [main.go](#maingo)
   - [main_test.go](#main_testgo)
8. [Running the Program](#running-the-program)
9. [Expected Outputs](#expected-outputs)
10. [Common Bugs and Fixes](#common-bugs-and-fixes)
11. [Extensions](#extensions)
12. [Submission Checklist](#submission-checklist)
13. [Resources](#resources)

---

## What This Project Does

You build a command-line tool in Go that takes a string argument and prints it as large ASCII art in the terminal. Each character in the input is looked up in a banner font file and rendered as an 8-line-tall graphic. All characters on the same line are printed side by side.

```
$ go run . "Hi"

 _    _   _  
| |  | | (_) 
| |__| |  _  
|  __  | | | 
| |  | | | | 
|_|  |_| |_| 
             
             
```

---

## How It Works

The program has three clear steps:

**Step 1 — Read the banner file**
`ReadBanner` opens a `.txt` file (e.g. `standard.txt`) and loads every line into a slice of strings.

**Step 2 — Build the character map**
`BuildAsciiMap` takes those lines and organizes them into a lookup table where every character (`A`–`Z`, `a`–`z`, `0`–`9`, symbols, space) maps to its 8-line art representation.

**Step 3 — Print the art**
`PrintAscii` takes the user's input string, looks up each character in the map, and prints the art row by row so all characters appear side by side.

---

## Prerequisites

- Go installed (version 1.18 or higher recommended)
- A terminal or VS Code with integrated terminal
- The three banner files: `standard.txt`, `shadow.txt`, `thinkertoy.txt`

Check your Go version:
```bash
go version
```

---

## Project Structure

```
ascii-art/
├── main.go              ← entry point, argument handling
├── go.mod               ← Go module definition
├── ascii/
│   └── render.go        ← all core logic: read, build map, print
├── standard.txt         ← standard banner font
├── shadow.txt           ← shadow banner font
└── thinkertoy.txt       ← thinkertoy banner font
```

> The `.txt` banner files must sit in the root of your project — the same level as `main.go`, NOT inside the `ascii/` folder. This is because `os.ReadFile` looks relative to where you run `go run .`

---

## Setting Up

**1. Create your project folder**
```bash
mkdir ascii-art
cd ascii-art
```

**2. Initialise the Go module**
```bash
go mod init ascii-art
```
This creates `go.mod`. The module name `ascii-art` is critical — it's what `main.go` uses when importing `"ascii-art/ascii"`.

**3. Create the ascii package folder and file**
```bash
mkdir ascii
touch ascii/render.go
```

**4. Create your main file**
```bash
touch main.go
```

**5. Copy the banner files into the root folder**
Copy `standard.txt`, `shadow.txt`, and `thinkertoy.txt` into your `ascii-art/` root folder.

---

## The Banner Files Explained

This is the most important concept to understand before reading any code.

Open `standard.txt`. You will see:

- The file **starts with a blank line**
- Every character is **exactly 8 lines tall**
- Characters are **separated by one blank line** — making each character slot **9 lines total**
- Characters start at **ASCII 32** (the space character) and go up to ASCII 126 (`~`)

So the structure looks like this (dots represent spaces for visibility):

```
[blank line]          ← the file starts here (line 0)
........              ← line 1  \
........              ← line 2   |
........              ← line 3   |  space character (ASCII 32)
........              ← line 4   |  8 lines of art
........              ← line 5   |
........              ← line 6   |
........              ← line 7   |
........              ← line 8  /
[blank line]          ← line 9  (separator between characters)
 _                    ← line 10 \
| |                   ← line 11  |
| |                   ← line 12  |  ! character (ASCII 33)
| |                   ← line 13  |  8 lines of art
|_|                   ← line 14  |
(_)                   ← line 15  |
                      ← line 16  |
                      ← line 17 /
[blank line]          ← line 18 (separator)
...and so on for all 95 printable characters
```

**Key insight — why the loop steps by 9:**

Each character takes up exactly 9 lines in the file (8 art lines + 1 blank separator). When `BuildAsciiMap` loops with `i += 9`, it jumps from one character's block to the next.

- `i = 0` → space character (ASCII 32)
- `i = 9` → `!` character (ASCII 33)
- `i = 18` → `"` character (ASCII 34)
- `i = 27` → `#` character (ASCII 35)

> In this implementation, the blank line at the very top of the file is treated as line 0 of the space character's 8-line block. This works fine because all 8 lines of the space character are blank anyway.

---

## The Code — File by File, Line by Line

### go.mod

Created automatically by `go mod init ascii-art`. It looks like:

```
module ascii-art

go 1.21
```

The module name `ascii-art` is what makes the import path `"ascii-art/ascii"` work in `main.go`. Replace `1.21` with whatever `go version` shows on your machine.

---

### ascii/render.go

This single file lives in the `ascii/` subfolder and contains all three core functions. Because it's in a subfolder with its own package name, `main.go` must explicitly import it.

**The full file:**

```go
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

    for i := 0; i < len(lines); i += 9 {
        asciiMap[rune(char)] = lines[i : i+8]
        char++
    }

    return asciiMap
}

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

**`package ascii`**

This declares the package name. Every file inside the `ascii/` folder must have this at the top. It is different from `package main` — this is a reusable package, not an entry point. Functions in this package are accessed from outside as `ascii.FunctionName`.

---

**`func ReadBanner(file string) ([]string, error)`**

Defines a function that takes one input (the filename) and returns two things simultaneously — Go allows multiple return values:
- `[]string` — a slice of strings (every line of the file)
- `error` — either `nil` (success) or a description of what went wrong

Returning an error alongside the result is Go's standard way of handling operations that can fail, like reading a file.

---

**`data, err := os.ReadFile(file)`**

`os.ReadFile` opens the file, reads every byte, closes the file, and returns two values at once using `:=` (short variable declaration — creates both variables and infers their types automatically):
- `data` — a `[]byte` (the raw contents of the file as bytes)
- `err` — `nil` if successful, or an error value if something went wrong

---

**`if err != nil { return nil, err }`**

`nil` in Go means "nothing" or "no value". If `err` is not nil, the file couldn't be read — maybe it doesn't exist or can't be opened. We return early: `nil` for the slice (we have nothing useful to give back) and pass the error up to `main.go` to deal with. This pattern — check for error immediately, return early — is used throughout Go code.

---

**`lines := strings.Split(string(data), "\n")`**

Two operations in one line:
- `string(data)` converts the raw `[]byte` into a readable Go string
- `strings.Split(..., "\n")` cuts that string into pieces at every newline character, returning a `[]string`

After this, `lines[0]` is the first line of the file, `lines[1]` is the second, and so on. The entire banner file is now in memory as individually accessible strings.

---

**`return lines, nil`**

Return the completed slice and `nil` for the error — meaning everything worked fine.

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

Takes the slice of lines as input and returns a `map[rune][]string`.

A **map** is a lookup table — like a dictionary. You give it a key, it gives you a value back instantly.
- Key type: `rune` — Go's type for a single character. A rune is just an `int32` that stores the character's Unicode/ASCII number. `'A'` is stored as `65`, `' '` as `32`, `'!'` as `33`. You can do math with runes: `rune(65)` gives you `'A'`.
- Value type: `[]string` — a slice of 8 strings (the 8 art lines for that character)

So `asciiMap['A']` gives you back the 8 lines that draw the letter A. `asciiMap['!']` gives you the 8 lines for an exclamation mark.

---

**`asciiMap := make(map[rune][]string)`**

`make` creates an empty map ready to use. You must use `make` — just writing `var asciiMap map[rune][]string` would create a `nil` map that crashes when you try to add anything to it. Think of `make` as building the empty filing cabinet before you start filing things in it.

---

**`char := 32`**

Start at ASCII code 32, which is the space character. This is exactly where the banner file begins. `char` is a plain integer that tracks which ASCII character we're currently processing. It increments by 1 each loop iteration to move through the full ASCII range.

---

**`for i := 0; i < len(lines); i += 9`**

The key loop. `i` starts at 0 and steps forward by 9 on each iteration. Why 9? Because each character occupies exactly 9 lines in the banner file — 8 art lines plus 1 blank separator line.

So the loop visits:
- `i = 0` → start of space (ASCII 32)
- `i = 9` → start of `!` (ASCII 33)
- `i = 18` → start of `"` (ASCII 34)
- `i = 27` → start of `#` (ASCII 35)
- ...continuing until all 95 characters are processed

---

**`asciiMap[rune(char)] = lines[i : i+8]`**

This is the most important line in `BuildAsciiMap`. Two things:

`rune(char)` converts the integer `char` (e.g. `65`) into a rune — the actual character `'A'`. This becomes the map key.

`lines[i : i+8]` is a **slice expression**. It extracts elements from index `i` up to but not including `i+8`, giving exactly 8 elements — the 8 art lines for this character. The 9th line (the blank separator) is automatically skipped because the loop jumps to `i+9` next iteration.

Put together: *"store the 8 art lines starting at position i, filed under the key for this character"*.

---

**`char++`**

Move to the next ASCII character. After processing space (32), move to `!` (33). After `!` move to `"` (34). After all 95 characters, the loop ends because `i >= len(lines)`.

---

**`return asciiMap`**

Return the completed map. After this function runs, `asciiMap['H']` gives 8 strings, `asciiMap['e']` gives 8 strings, etc. — ready for `PrintAscii` to use.

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

**What it does:** Takes the user's input and the character map, then prints the final ASCII art to the terminal. This is where the actual rendering happens.

---

**The core problem — why you can't print one character at a time:**

Your first instinct might be: loop through each character, print its 8 lines, then move to the next. But that stacks them vertically:

```
 _    _          ← all 8 rows of H printed
| |  | |
| |__| |
|  __  |
| |  | |
|_|  |_|


 _               ← all 8 rows of i printed
(_)
 _
| |
| |
|_|
```

That's wrong. You need side by side:

```
 _    _   _      ← row 0 of H  +  row 0 of i  printed on same line
| |  | | (_)     ← row 1 of H  +  row 1 of i
| |__| |  _      ← row 2 of H  +  row 2 of i
|  __  | | |     ← row 3 of H  +  row 3 of i
| |  | | | |     ← row 4 of H  +  row 4 of i
|_|  |_| |_|     ← row 5 of H  +  row 5 of i
                 ← row 6 of H  +  row 6 of i
                 ← row 7 of H  +  row 7 of i
```

To achieve this, you must think in rows first, characters second. For row 0, print row 0 of every character. For row 1, print row 1 of every character. That is the double loop pattern at the heart of this function.

---

**`func PrintAscii(text string, asciiMap map[rune][]string)`**

Takes the input string and the lookup map. Returns nothing — output goes straight to the terminal via `fmt.Print`.

---

**`lines := strings.Split(text, "\\n")`**

Splits the input on the literal two-character sequence: backslash then `n`.

This is subtle and important. When a user types:
```bash
go run . "Hello\nThere"
```

Go receives `Hello\nThere` as a 12-character string — a backslash followed by `n`. It is NOT a real newline. The shell passes it literally without interpreting it.

In a Go string literal, `"\\n"` means: a backslash character followed by `n` (two characters). This correctly matches and splits the user's literal `\n` sequence.

Compare the two:
```go
strings.Split(text, "\n")   // splits on a REAL newline (1 character) — WRONG here
strings.Split(text, "\\n")  // splits on backslash+n (2 characters) — CORRECT
```

Result examples:
```
"Hello\nThere"    → ["Hello", "There"]
"Hello\n\nThere"  → ["Hello", "", "There"]   ← middle empty string = blank line
"\n"              → ["", ""]
```

---

**`for _, line := range lines`**

`range` on a slice gives two values per iteration: the index and the value. The index is discarded with `_` because we don't need it. `line` holds each segment of the input.

For `"Hello\nThere"` split into `["Hello", "There"]`:
- Iteration 1: `line = "Hello"`
- Iteration 2: `line = "There"`

---

**`if line == "" { fmt.Println(); continue }`**

An empty `line` means the user typed `\n\n` — a double newline. After splitting, this produces an empty string between the two segments: `["Hello", "", "There"]`.

`fmt.Println()` with no arguments prints a single blank line, which is the correct output for an empty segment (a gap between two art blocks).

`continue` skips all the rendering code below and jumps to the next iteration. There's nothing to draw for an empty segment.

---

**`for row := 0; row < 8; row++`**

The outer rendering loop — first half of the double loop. Runs exactly 8 times, once for each row of the ASCII art height. `row` goes from 0 to 7.

For the current input segment, we need to produce 8 output lines. This loop controls which row we're currently building.

---

**`for _, char := range line`**

The inner loop — second half of the double loop. `range` on a string gives the byte index (discarded with `_`) and the actual character as a `rune`.

For `line = "Hi"`:
- Iteration 1: `char = 'H'` (rune 72)
- Iteration 2: `char = 'i'` (rune 105)

For each character in the input, we grab that character's art for the current row.

---

**`fmt.Print(asciiMap[char][row])`**

The single most important line in the whole program. Three things chained:

- `asciiMap[char]` — looks up this character in the map, getting back its `[]string` of 8 art lines
- `[row]` — indexes into that slice to get just the one string for the current row (0 through 7)
- `fmt.Print(...)` — prints it WITHOUT a newline

Using `Print` instead of `Println` is critical. Because there's no newline added, the next character's row gets printed immediately after on the same terminal line — this is what produces the side-by-side effect.

**Visualizing one full pass for `"Hi"` at `row = 0`:**
```
char = 'H'  →  asciiMap['H'][0]  =  " _    _  "  →  prints " _    _  "
char = 'i'  →  asciiMap['i'][0]  =  " _  "        →  prints " _  "

Terminal line so far: " _    _   _  "   (no newline yet)
```

---

**`fmt.Println()`**

After all characters have contributed their piece of row `row`, this prints a newline. The cursor moves to the next terminal line, ready for `row = 1`.

After all 8 rows complete, the full block of side-by-side ASCII art for this input segment is printed to the terminal.

---

### main.go

```go
package main

import (
    "fmt"
    "os"

    "ascii-art/ascii"
)

func main() {

    if len(os.Args) != 2 {
        return
    }

    input := os.Args[1]

    bannerLines, err := ascii.ReadBanner("thinkertoy.txt")
    if err != nil {
        fmt.Println(err)
        return
    }

    asciiMap := ascii.BuildAsciiMap(bannerLines)

    ascii.PrintAscii(input, asciiMap)
}
```

---

**`package main`**

Every Go program needs exactly one `main` package with exactly one `main()` function. This is where execution begins when you run `go run .`

---

**`"ascii-art/ascii"`**

Imports the `ascii` package from the subfolder. The import path format is `<module-name>/<folder-name>`. Since `go.mod` defines the module as `ascii-art` and the code is in a folder called `ascii`, the import is `ascii-art/ascii`. After importing, you call its exported functions with the prefix `ascii.` — e.g. `ascii.ReadBanner(...)`.

Note: in Go, only functions starting with a capital letter are exported (accessible from other packages). That's why the functions are named `ReadBanner`, `BuildAsciiMap`, and `PrintAscii` — not `readBanner` etc.

---

**`if len(os.Args) != 2 { return }`**

`os.Args` is a slice of strings containing everything typed on the command line:
- `os.Args[0]` = the program itself
- `os.Args[1]` = the first real argument (the string to render)

`len(os.Args) != 2` means the user didn't provide exactly one argument. We `return` to exit silently. A better version would print a usage message — this is improved in the `ascii-art-fs` extension.

---

**`input := os.Args[1]`**

Store the user's input string. This is the string that will be rendered into ASCII art.

---

**`bannerLines, err := ascii.ReadBanner("thinkertoy.txt")`**

Call `ReadBanner` from the imported `ascii` package. Currently hardcoded to `"thinkertoy.txt"` — this gets updated in the `ascii-art-fs` extension to accept any banner by name.

---

**`if err != nil { fmt.Println(err); return }`**

If the file couldn't be read, print the error and exit. This is the same error-checking pattern used in `render.go`.

---

**`asciiMap := ascii.BuildAsciiMap(bannerLines)`**

Build the character lookup map from the raw lines. After this, `asciiMap` contains all 95 characters mapped to their art.

---

**`ascii.PrintAscii(input, asciiMap)`**

Render and print the ASCII art. No return value — the output goes directly to the terminal.

---

### main_test.go

Tests are required by the audit. Create this file in the root of your project alongside `main.go`.

```go
package main

import (
    "testing"

    "ascii-art/ascii"
)

// TestReadBanner tests that the banner file loads without error
func TestReadBanner(t *testing.T) {
    lines, err := ascii.ReadBanner("standard.txt")
    if err != nil {
        t.Fatal("Could not read banner file:", err)
    }
    if len(lines) == 0 {
        t.Error("Expected lines from banner file, got empty slice")
    }
}

// TestBuildAsciiMap tests the map contains the expected number of characters
func TestBuildAsciiMap(t *testing.T) {
    lines, err := ascii.ReadBanner("standard.txt")
    if err != nil {
        t.Fatal(err)
    }

    asciiMap := ascii.BuildAsciiMap(lines)

    if len(asciiMap) != 95 {
        t.Errorf("Expected 95 characters in map, got %d", len(asciiMap))
    }

    // Check that 'A' exists and has exactly 8 lines
    artLines, ok := asciiMap['A']
    if !ok {
        t.Error("Expected 'A' to exist in map")
    }
    if len(artLines) != 8 {
        t.Errorf("Expected 8 lines for 'A', got %d", len(artLines))
    }
}

// TestAllBanners tests that all three banner files load and produce valid maps
func TestAllBanners(t *testing.T) {
    banners := []string{"standard.txt", "shadow.txt", "thinkertoy.txt"}
    for _, b := range banners {
        lines, err := ascii.ReadBanner(b)
        if err != nil {
            t.Errorf("Failed to load %s: %v", b, err)
            continue
        }
        asciiMap := ascii.BuildAsciiMap(lines)
        if len(asciiMap) == 0 {
            t.Errorf("Empty map for banner %s", b)
        }
    }
}

// TestSpaceCharacterExists tests that the space character is in the map
func TestSpaceCharacterExists(t *testing.T) {
    lines, _ := ascii.ReadBanner("standard.txt")
    asciiMap := ascii.BuildAsciiMap(lines)

    _, ok := asciiMap[' ']
    if !ok {
        t.Error("Expected space character in map")
    }
}

// TestUpperAndLowerAreDifferent tests that 'A' and 'a' have different art
func TestUpperAndLowerAreDifferent(t *testing.T) {
    lines, _ := ascii.ReadBanner("standard.txt")
    asciiMap := ascii.BuildAsciiMap(lines)

    upper := asciiMap['A']
    lower := asciiMap['a']

    if len(upper) == 0 || len(lower) == 0 {
        t.Error("Expected both A and a to exist in map")
    }

    if upper[0] == lower[0] {
        t.Error("Expected A and a to have different art on row 0")
    }
}
```

**Why each test exists:**

- `TestReadBanner` — confirms the file read works without crashing
- `TestBuildAsciiMap` — confirms the map has all 95 characters and each has 8 lines
- `TestAllBanners` — confirms all three font files work, not just standard
- `TestSpaceCharacterExists` — space (ASCII 32) is the first character and easiest to miss
- `TestUpperAndLowerAreDifferent` — catches bugs where the loop step is wrong

Run all tests:
```bash
go test ./... -v
```

---

## Running the Program

```bash
# Basic usage (currently uses thinkertoy.txt)
go run . "hello"

# With pipe to see exact line endings ($ marks end of each line)
go run . "hello" | cat -e

# Newline between words
go run . "Hello\nThere" | cat -e

# Double newline (blank line between art blocks)
go run . "Hello\n\nThere" | cat -e

# Numbers and special characters
go run . "1Hello 2There" | cat -e

# Run all tests
go test ./... -v
```

---

## Expected Outputs

### `go run . "hello"` — standard banner
```
 _              _   _          
| |            | | | |         
| |__     ___  | | | |   ___   
|  _ \   / _ \ | | | |  / _ \  
| | | | |  __/ | | | | | (_) | 
|_| |_|  \___| |_| |_|  \___/  
                               
                               
```

### `go run . "hello"` — shadow banner
```
                                          
_|            _| _|                       
_|_|_|   _|_| _| _|   _|_|               
_|    _|_|_|_|_| _| _|    _|             
_|    _|_|    _| _| _|    _|             
_|    _|  _|_|_|_|   _|_|               
                                         
                                         
```

### `go run . "Hello\nThere"` — standard banner
```
 _    _          _   _          
| |  | |        | | | |         
| |__| |   ___  | | | |   ___   
|  __  |  / _ \ | | | |  / _ \  
| |  | | |  __/ | | | | | (_) | 
|_|  |_|  \___| |_| |_|  \___/  
                                
                                
 _______   _                           
|__   __| | |                          
   | |    | |__     ___   _ __    ___  
   | |    |  _ \   / _ \ | '__|  / _ \
   | |    | | | | |  __/ | |    |  __/
   |_|    |_| |_|  \___| |_|     \___|
                                      
                                      
```

---

## Common Bugs and Fixes

### Characters appear shifted — A renders as B, etc.
Your loop step is wrong. It must be `i += 9` — each character takes exactly 9 lines. If you use `i += 8`, you drift by one line per character and everything shifts.

### `\n` in input is not producing separate art blocks
Make sure you split on the two-character literal sequence, not a real newline:
```go
// WRONG — splits on actual newline character (1 char)
strings.Split(text, "\n")

// CORRECT — splits on backslash followed by n (2 chars)
strings.Split(text, "\\n")
```

### `no such file or directory` when reading banner files
The banner `.txt` files must be in the project root — same folder where you run `go run .`. They should NOT be inside the `ascii/` subfolder. `os.ReadFile` resolves paths relative to your working directory, which is wherever you run the command from.

### `undefined: ascii.ReadBanner`
Check two things:
- `ascii/render.go` starts with `package ascii` (not `package main`)
- `main.go` imports `"ascii-art/ascii"` and your `go.mod` says `module ascii-art`

### Functions not accessible from main.go
In Go, only functions starting with a capital letter are exported from a package. `ReadBanner` ✓ but `readBanner` ✗. Check that all three functions in `render.go` start with a capital letter.

### Windows `\r\n` line ending issues
If your banner files were edited on Windows, lines have a trailing `\r` invisible character that corrupts lookups. Fix in `ReadBanner` after reading the file:
```go
content := strings.ReplaceAll(string(data), "\r\n", "\n")
lines := strings.Split(content, "\n")
```

---

## Extensions

Once the base project works and passes all audit tests, implement these in order.

---

### ascii-art-fs (Banner Selection)

Allows choosing a banner as a second argument:
```bash
go run . "hello" shadow
go run . "hello" thinkertoy
```

Update `main.go`:
```go
func main() {
    if len(os.Args) < 2 || len(os.Args) > 3 {
        fmt.Println("Usage: go run . [STRING] [BANNER]")
        fmt.Println("\nEX: go run . something standard")
        os.Exit(1)
    }

    input := os.Args[1]

    banner := "standard"
    if len(os.Args) == 3 {
        banner = os.Args[2]
    }

    bannerLines, err := ascii.ReadBanner(banner + ".txt")
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    asciiMap := ascii.BuildAsciiMap(bannerLines)
    ascii.PrintAscii(input, asciiMap)
}
```

---

### ascii-art-output (Write to File)

Adds `--output=<filename>` to write art to a file:
```bash
go run . --output=banner.txt "hello" standard
```

Add a `RenderToString` function in `render.go` that returns the art as a string instead of printing it, then write it to a file with `os.WriteFile`.

---

### ascii-art-color (Colorized Output)

Adds `--color=<color>` to colorize output using ANSI escape codes:
```bash
go run . --color=red "hello"
```

ANSI codes tell the terminal to change text color:
```
\033[31m  ← start red text
\033[0m   ← reset back to default
```

Add a color map to `render.go`:
```go
var ColorCodes = map[string]string{
    "red":     "\033[31m",
    "green":   "\033[32m",
    "yellow":  "\033[33m",
    "blue":    "\033[34m",
    "cyan":    "\033[36m",
    "white":   "\033[37m",
    "reset":   "\033[0m",
}
```

---

### ascii-art-justify (Text Alignment)

Adds `--align=<left|right|center|justify>`:
```bash
go run . --align=center "hello" standard
```

You need the terminal width to calculate padding. On Linux you can get it via `syscall.TIOCGWINSZ` or fall back to a default of 80 columns.

---

### ascii-art-reverse (Decode Art Back to Text)

Reads an ASCII art file and figures out what string it represents:
```bash
go run . --reverse=file.txt
```

Approach:
1. Load the banner map normally
2. Build a reverse map: join each character's 8 lines as a single string → that string maps back to the character
3. Read the art file, split into 8-line blocks, match each block against the reverse map

---

## Submission Checklist

**Functionality**
- [ ] `go run . "hello"` matches the exact audit output
- [ ] `go run . ""` produces no output
- [ ] `go run . "\n"` produces one blank line
- [ ] `go run . "Hello\nThere"` produces two separate art blocks
- [ ] `go run . "Hello\n\nThere"` has a blank line between the blocks
- [ ] Numbers, uppercase, lowercase, and special characters all render correctly
- [ ] All three banner files are in the project root

**Code Quality**
- [ ] `ascii/render.go` starts with `package ascii`
- [ ] `main.go` starts with `package main`
- [ ] No external packages — only Go standard library
- [ ] No commented-out code or debug print statements left in
- [ ] `go build .` compiles with zero errors

**Tests**
- [ ] `main_test.go` exists in the project root
- [ ] `go test ./... -v` passes with no failures
- [ ] Tests cover: file loading, map building, all three banners, space character, upper vs lower

**File Structure**
- [ ] `main.go` and `go.mod` in root
- [ ] `ascii/render.go` in the `ascii/` subfolder
- [ ] `standard.txt`, `shadow.txt`, `thinkertoy.txt` in root

---

## Resources

| Topic | Link |
|-------|------|
| Go standard library | https://pkg.go.dev/std |
| os package (ReadFile, Args) | https://pkg.go.dev/os |
| strings package | https://pkg.go.dev/strings |
| fmt package | https://pkg.go.dev/fmt |
| Writing Go tests | https://go.dev/doc/tutorial/add-a-test |
| Go packages and imports | https://go.dev/doc/code |
| ASCII character table | https://www.asciitable.com |
| Go tour (interactive basics) | https://go.dev/tour |
| Effective Go (best practices) | https://go.dev/doc/effective_go |

---

*Only the Go standard library is allowed for this project. No third-party packages.*