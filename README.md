# ascii-art-color 🎨

A command-line tool written in Go that converts text into ASCII art and lets you color the entire output — or just specific words — using ANSI terminal colors.

---

## Table of Contents

1. [What Does This Program Do?](#what-does-this-program-do)
2. [How to Run It](#how-to-run-it)
3. [Project Structure](#project-structure)
4. [Understanding `render.go` — Line by Line](#understanding-rendergo--line-by-line)
   - [Package and Imports](#package-and-imports)
   - [ReadBanner](#readbanner)
   - [BuildAsciiMap](#buildasciimap)
   - [The Colors Map](#the-colors-map)
   - [PrintAscii](#printascii)
5. [Understanding `main.go` — Line by Line](#understanding-maingo--line-by-line)
   - [Parsing the --color Flag](#parsing-the---color-flag)
   - [Validating the Color](#validating-the-color)
   - [Figuring Out the Arguments](#figuring-out-the-arguments)
   - [Loading the Banner and Printing](#loading-the-banner-and-printing)
6. [How the Banner Files Work](#how-the-banner-files-work)
7. [The Coloring Logic Explained](#the-coloring-logic-explained)
8. [Example Runs](#example-runs)

---

## What Does This Program Do?

When you type something into this program, it draws it in giant ASCII art characters using one of three font styles (`standard`, `shadow`, or `thinkertoy`). On top of that, you can tell it:

- Color the **entire** text one color.
- Color only a **specific substring** (a word or part of a word) while leaving the rest uncolored.

Example — coloring just `kit` in the word `kitten`:

```
go run . --color=red kit "kitten"
```

The letters `k`, `i`, `t` will be drawn in red. The letters `t`, `e`, `n` will be drawn in the default terminal color.

---

## How to Run It

Make sure you have Go installed, then:

```bash
# Color the entire text
go run . --color=<color> "your text"

# Color only a specific part of the text
go run . --color=<color> <substring> "your text"
```

**Available colors:** `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `purple`, `cyan`, `white`, `orange`

**Available banner styles:** `standard` (default), `shadow`, `thinkertoy`

```bash
go run . --color=blue "Hello"
```

---

## Project Structure

```
ascii-art-color/
├── main.go          ← Entry point: reads your arguments, calls render functions
├── ascii/
│   └── render.go    ← Core logic: reads fonts, builds the art, handles colors
├── standard.txt     ← Font file: the "standard" ASCII art style
├── shadow.txt       ← Font file: the "shadow" ASCII art style
└── thinkertoy.txt   ← Font file: the "thinkertoy" ASCII art style
```

The font `.txt` files are pre-made files where every printable character (from space to `~`) is drawn using 8 lines of plain text characters. The program reads them and looks up the right character's art when building the output.

---

## Understanding `render.go` — Line by Line

This file lives inside the `ascii` folder and is the brain of the program. It does three things: reads the font file, builds a lookup table, and draws the final colored output.

---

### Package and Imports

```go
package ascii
```

This line says: "Everything in this file belongs to the `ascii` package." In Go, files in the same folder share a package name. Because this package is named `ascii` (not `main`), it can't run on its own — it gets imported and used by `main.go`.

```go
import (
    "os"
    "strings"
)
```

`os` gives us tools to read files from the disk.  
`strings` gives us tools to work with text — splitting, searching, building strings efficiently.

---

### ReadBanner

```go
func ReadBanner(file string) ([]string, error) {
```

This declares a function named `ReadBanner`. It takes one input: `file`, which is the name of the font file (like `"standard.txt"`). It returns two things: a slice of strings (`[]string` — think of this as a list of lines) and an `error` (which will be `nil` if everything worked, or contain the problem description if it didn't).

Returning errors like this is Go's standard way of handling failure — no exceptions, just an explicit check: "did it work?"

```go
    data, err := os.ReadFile(file)
```

`os.ReadFile` opens the file and reads its entire contents as raw bytes. The `:=` is Go's shorthand for "declare this variable and assign to it at the same time." We get two things back: `data` (the file contents) and `err` (any error that occurred).

```go
    if err != nil {
        return nil, err
    }
```

If something went wrong (file doesn't exist, wrong path, etc.), `err` won't be `nil`. We immediately return `nil` (empty list) and the error so the caller knows what happened. This is the Go pattern: check the error right after every operation that could fail.

```go
    lines := strings.Split(string(data), "\n")
    return lines, nil
```

`string(data)` converts the raw bytes into a readable string. `strings.Split(..., "\n")` then cuts that string at every newline character, giving us a list where each item is one line of the file. We return that list and `nil` for the error (meaning success).

---

### BuildAsciiMap

```go
func BuildAsciiMap(lines []string) map[rune][]string {
```

This function takes the list of lines from the font file and turns it into a **map** — a lookup table where you give it a character and it gives you that character's 8 art lines back.

A `rune` in Go is a single character (it handles Unicode properly, unlike a plain `byte`). So this map is: character → list of 8 art lines.

```go
    asciiMap := make(map[rune][]string)
```

`make` is how you create a map in Go. This gives us an empty map ready to be filled.

```go
    char := 32
```

ASCII character 32 is the **space** character. The font files store characters starting from space and going through all printable characters up to `~` (ASCII 126). So we start our counter at 32.

```go
    for i := 1; i < len(lines); i += 9 {
```

This loop walks through the lines of the font file. Why start at `1`? Because font files have one blank line at the very top (line index 0), so the first real character starts at line 1. Why step by `9`? Each character takes up exactly 8 lines of art, plus 1 blank separator line between characters — so every 9 lines is one new character.

```go
        asciiMap[rune(char)] = lines[i : i+8]
```

`rune(char)` converts our integer counter into a character. `lines[i : i+8]` is a **slice** — it takes lines from index `i` up to (but not including) `i+8`, giving us exactly 8 lines. We store those 8 lines as the art for this character.

```go
        char++
    }
    return asciiMap
}
```

We increment `char` to move to the next ASCII character and loop again. After the loop, every printable character from space to `~` has its 8-line art stored in the map.

---

### The Colors Map

```go
var Colors = map[string]string{
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
```

`var` declares a package-level variable — it exists as long as the program runs, and because it starts with a capital letter (`Colors`), it is **exported** (visible and usable from `main.go`).

The values like `"\033[31m"` are **ANSI escape codes** — special sequences that terminals understand as instructions, not visible text. Breaking one down:

- `\033` — the Escape character (also written `\x1b`)
- `[` — starts the control sequence
- `31` — the color code (31 = red)
- `m` — ends the sequence, meaning "apply this color"

When your terminal sees `\033[31m`, it switches the text color to red. When it sees `\033[0m` (the `"reset"` entry), it switches back to the default color. Everything printed between those two codes appears colored.

`orange` uses a different format (`38;5;214`) because orange isn't one of the 8 standard terminal colors — it uses the 256-color extended palette, where color 214 happens to be orange.

---

### PrintAscii

This is the most important function. It takes your text, looks up each character's art, applies colors, and returns the full ASCII art as one big string.

```go
func PrintAscii(text string, asciiMap map[rune][]string, subStr, colorArg string) string {
```

Four inputs:

- `text` — what the user typed (e.g., `"kitten"`)
- `asciiMap` — the lookup table built by `BuildAsciiMap`
- `subStr` — the part to color (e.g., `"kit"`), or empty string if coloring everything
- `colorArg` — the color name (e.g., `"red"`), or empty string if no color

It returns a `string` — the complete, ready-to-print ASCII art.

```go
    var result strings.Builder
```

A `strings.Builder` is an efficient way to build a large string piece by piece. Instead of doing `bigString = bigString + newPiece` (which creates a brand new string every time and wastes memory), a Builder collects all the pieces and assembles them once at the end.

```go
    lines := strings.Split(text, "\\n")
```

This splits the input on `\n` — but notice it's `"\\n"` (a backslash followed by the letter n), not a real newline character. That's because users type `\n` literally in their terminal argument to mean "new line here." For example: `go run . "Hello\nWorld"` should draw `Hello` and `World` on separate art lines. The `\\n` in Go source code represents the two-character sequence the user actually typed.

```go
    colorCode := Colors[colorArg]
    reset := Colors["reset"]
```

Look up the ANSI code for the requested color once upfront. If `colorArg` is `""` (empty), `Colors[""]` returns `""` (empty string), which is fine — writing an empty string doesn't change anything.

```go
    for i, line := range lines {
```

Loop through each line segment of the text. `i` is the index (0, 1, 2...) and `line` is the actual content. We need `i` below to handle the very first segment specially.

```go
        if line == "" {
            if i != 0 {
                result.WriteString("\n")
            }
            continue
        }
```

If this segment is empty (the user typed `\n` at the start or used back-to-back `\n`), we write a single newline (to create a blank line in the output) and `continue` — skip the rest of the loop body and go to the next segment. The `i != 0` check prevents adding a stray blank line before anything has been printed yet.

```go
        for row := 0; row < 8; row++ {
```

This is the **core rendering loop**. Each character's art is 8 rows tall. To make characters appear side by side (like normal text), we render **one row of every character across the whole line**, then move to the next row. If we did it the other way — all 8 rows of one character, then all 8 rows of the next — the letters would stack vertically instead of sitting beside each other.

Think of it like printing a spreadsheet row by row instead of column by column.

```go
            if colorArg != "" && subStr != "" {
```

This branch handles the **substring coloring** case — where only certain parts of the text should be colored.

```go
                parts := strings.Split(line, subStr)
```

`strings.Split` cuts the line at every occurrence of `subStr` and gives back the pieces in between. For example:

```
line   = "a kitten has a kit"
subStr = "kit"
parts  = ["a ", "ten has a ", ""]
```

The `subStr` itself is removed — it sits _between_ the pieces. The number of colored substrings to insert equals `len(parts) - 1` (always one fewer than the number of pieces).

```go
                for p, part := range parts {
                    for _, ch := range part {
                        result.WriteString(asciiMap[ch][row])
                    }
```

For each uncolored piece, write the art for each of its characters at the current row. No color codes here — this is the plain text between the colored substrings.

```go
                    if p < len(parts)-1 {
                        for _, ch := range subStr {
                            result.WriteString(colorCode + asciiMap[ch][row] + reset)
                        }
                    }
```

After each piece (except the last one), insert the colored version of `subStr`. `p < len(parts)-1` is the guard that prevents inserting an extra colored substring after the final piece. We wrap each character's art individually with `colorCode + art + reset` to prevent color bleeding into adjacent characters.

```go
            } else {
                for _, ch := range line {
                    if colorArg != "" {
                        result.WriteString(colorCode + asciiMap[ch][row] + reset)
                    } else {
                        result.WriteString(asciiMap[ch][row])
                    }
                }
            }
```

This is the simpler branch. If we're coloring everything (or nothing at all), just loop through every character in the line and either wrap it in color codes or write it plain.

```go
            result.WriteString("\n")
        }
    }
    return result.String()
}
```

After all 8 rows for a line of text are written, add a newline so the next block of art starts fresh on a new terminal line. Once all lines are processed, `result.String()` assembles everything the Builder collected into one final string and returns it.

---

## Understanding `main.go` — Line by Line

`main.go` is the entry point — the file Go runs first. Its only job is to read what the user typed, validate it, and call the functions in `render.go`.

```go
package main
```

Every runnable Go program needs exactly one `main` package with exactly one `main()` function. This is it.

```go
import (
    "ascii-art-color/ascii"
    "fmt"
    "os"
    "strings"
)
```

- `ascii-art-color/ascii` — imports our own `ascii` package (the `render.go` file). The path matches the module name in `go.mod` plus the subfolder name.
- `fmt` — for printing to the terminal (`fmt.Println`, `fmt.Print`).
- `os` — for reading command-line arguments via `os.Args`.
- `strings` — for manipulating argument strings.

---

### Parsing the `--color` Flag

```go
colorArg := os.Args[1]
```

`os.Args` is a list of everything the user typed to run the program. `os.Args[0]` is always the program name itself. So `os.Args[1]` is the first argument the user provided — which we expect to be something like `--color=red`.

```go
if !strings.HasPrefix(colorArg, "--color=") {
    fmt.Println("Usage: go run . [OPTION] [STRING]")
    fmt.Println()
    fmt.Println("EX: go run . --color=<color> <substring to be colored> \"something\"")
    return
}
```

`strings.HasPrefix` checks if the string starts with `"--color="`. If it doesn't — meaning the user typed something unexpected — we print the usage instructions and `return` to exit `main()` early. Without a valid `--color=` flag, there's nothing meaningful to do.

```go
colorArgs := strings.TrimPrefix(colorArg, "--color=")
```

`strings.TrimPrefix` removes the `"--color="` part from the front, leaving just the color name. So `"--color=red"` becomes `"red"`. This is what gets passed to `PrintAscii`.

---

### Validating the Color

A safe check to add would be:

```go
if colorArgs != "" {
    if _, ok := ascii.Colors[colorArgs]; !ok {
        fmt.Println("Unknown color:", colorArgs)
        return
    }
}
```

The `_, ok` pattern is Go's way of checking if a key exists in a map. The map lookup returns two values: the value (which we discard using `_`) and a boolean `ok` that is `true` if the key was found. If `ok` is `false`, the user typed an unknown color and we should tell them rather than silently produce uncolored output.

---

### Figuring Out the Arguments

```go
var subStr string
var input string
```

`var` declares variables with their **zero value** — for strings, that's `""` (an empty string). This means both variables start out empty, and we only fill them in if the user provided them.

```go
if len(os.Args) == 3 {
    input = os.Args[2]
    subStr = ""
}

if len(os.Args) == 4 {
    subStr = os.Args[2]
    input = os.Args[3]
}
```

`len(os.Args)` counts the total items in the args list, including the program name at index 0:

- **3 total** → `program --color=red "text"` → last item is the text, no substring.
- **4 total** → `program --color=red kit "kitten"` → second item is the substring, third is the text.

---

### Loading the Banner and Printing

```go
bannerLines, err := ascii.ReadBanner("standard.txt")
if err != nil {
    fmt.Println(err)
    return
}
```

Call our `ReadBanner` function with the font file name. If it returns an error (file not found, etc.), print the error and exit. The font file must be in the same folder you run the program from.

```go
asciiMap := ascii.BuildAsciiMap(bannerLines)
```

Turn the raw lines into our character lookup table. This is fast and only needs to happen once per run.

```go
result := ascii.PrintAscii(input, asciiMap, subStr, colorArgs)
fmt.Print(result)
```

Call `PrintAscii` with everything it needs. It returns the full finished art as a string. `fmt.Print` (without `ln`) prints it exactly as-is — no extra newline added, since `PrintAscii` already manages newlines internally.

---

## How the Banner Files Work

Each `.txt` file stores every printable ASCII character from space (32) to tilde `~` (126), in order.

Here's what the structure looks like for the letter `A` in `standard.txt`:

```
    /\        ← row 0
   /  \       ← row 1
  / /\ \      ← row 2
 / ____ \     ← row 3
/_/    \_\    ← row 4
              ← row 5 (blank — still part of A's 8 lines)
              ← row 6
              ← row 7
              ← row 8: blank separator line between characters
```

Every character is exactly 8 rows tall, followed by 1 blank separator line = **9 lines per character**.

`BuildAsciiMap` reads these in order, starting at line 1 (skipping the leading blank line at the top of the file), stepping 9 lines at a time. Since space is ASCII 32 and `char` starts at 32 and increments by 1 for every 9-line block, after the whole file every character maps to exactly the right art.

---

## The Coloring Logic Explained

The key insight that makes coloring work correctly is `strings.Split`.

Say the user runs: `go run . --color=red kit "a kitten has a kit"`

```
line   = "a kitten has a kit"
subStr = "kit"
```

`strings.Split("a kitten has a kit", "kit")` returns:

```
["a ", "ten has a ", ""]
```

Now we process these pieces one by one. Between each pair of adjacent pieces, we insert the colored `kit`. Visually:

```
"a "          → plain art
"kit"         → RED colored art   ← inserted between piece 0 and piece 1
"ten has a "  → plain art
"kit"         → RED colored art   ← inserted between piece 1 and piece 2
""            → (empty, nothing to write)
```

This works for **any number of occurrences** automatically — however many times `kit` appears, `strings.Split` creates that many gaps to fill with colored art.

The guard `p < len(parts)-1` ensures we don't add a colored `kit` after the final piece. And if `kit` doesn't appear at all, `strings.Split` returns a single piece (the whole string unchanged) — the guard fires zero times, and nothing gets colored. Correct behavior in every case.

---

## Example Runs

```bash
# Color the whole word
go run . --color=cyan "Hello"

# Color only "ell" inside "Hello"
go run . --color=yellow ell "Hello"

# Multiple occurrences — both "kit" segments colored
go run . --color=red kit "a kitten has a kit"

# Multiline input using \n
go run . --color=green "Hello\nWorld"

# Nothing should be colored — "xyz" doesn't appear in "hello"
go run . --color=red xyz "hello world"
```

---

> Built with Go — no external dependencies, just the standard library.
