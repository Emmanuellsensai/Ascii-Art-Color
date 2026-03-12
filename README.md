# ascii-art-color

A command-line tool written in Go that converts text into ASCII art and lets you color the entire output — or just a specific word — using ANSI terminal colors.

---

## Table of Contents

1. [What Does This Program Do?](#1-what-does-this-program-do)
2. [How to Run It](#2-how-to-run-it)
3. [Project Structure](#3-project-structure)
4. [How the Font Files Work](#4-how-the-font-files-work)
5. [Understanding render.go — Line by Line](#5-understanding-rendergo--line-by-line)
   - [Package and Imports](#51-package-and-imports)
   - [ReadBanner](#52-readbanner)
   - [BuildAsciiMap](#53-buildasciimap)
   - [The Colors Map](#54-the-colors-map)
   - [PrintAscii](#55-printascii)
6. [Understanding main.go — Line by Line](#6-understanding-maingo--line-by-line)
7. [The Coloring Logic — A Full Walkthrough](#7-the-coloring-logic--a-full-walkthrough)
8. [Example Runs](#8-example-runs)

---

## 1. What Does This Program Do?

When you type text into this program it draws it as giant ASCII art characters in your terminal. Each letter is made from 8 rows of plain text characters placed side by side. On top of that you can:

- Color the **entire** output one color
- Color only a **specific word or substring** while leaving the rest uncolored

```bash
go run . --color=red kit "kitten"
```

The letters `k`, `i`, `t` appear in red. The letters `t`, `e`, `n` appear in the default terminal color.

---

## 2. How to Run It

```bash
# Color the entire text
go run . --color=<color> "your text"

# Color only a specific part of the text
go run . --color=<color> <substring> "your text"

# Multi-line art (type a literal backslash then n)
go run . --color=<color> "Hello\nWorld"
```

**Available colors:** `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `purple`, `cyan`, `white`, `orange`

---

## 3. Project Structure

```
ascii-art-color/
├── main.go          ← Entry point: reads arguments, validates them, calls render
├── ascii/
│   └── render.go    ← Core engine: reads font, builds lookup table, draws colored art
├── standard.txt     ← Font file: the "standard" ASCII art style
├── shadow.txt       ← Font file: the "shadow" ASCII art style
└── thinkertoy.txt   ← Font file: the "thinkertoy" ASCII art style
```

> **Important:** The `.txt` font files must stay in the project root — the same folder where you run `go run .`. Go looks for files relative to wherever you run the command from.

---

## 4. How the Font Files Work

Before reading any code you need to understand the font files, because the entire program is built around their structure.

Open `standard.txt`. Here is what you find:

```
              ← one blank line at the very top of the file
              ← 8 lines for the SPACE character (all blank)
              ←
              ←
              ←
              ←
              ←
              ←
              ←
              ← one blank separator line
 _            ← 8 lines for ! (ASCII code 33)
| |
| |
| |
|_|
(_)

              ← one blank separator line
...and so on for every printable character up to ~
```

**The rules:**
- Every character is exactly **8 lines tall**
- After those 8 lines there is always **1 blank separator line**
- So each character occupies **9 lines total** in the file
- The file starts at ASCII code **32** (space) and ends at code **126** (`~`)
- There is **one extra blank line at the very top** of the file before anything begins

This is why the code steps by `9` and starts at line index `1` — to skip that leading blank line.

---

## 5. Understanding render.go — Line by Line

This file is the engine of the program. It does three things: reads the font file from disk, builds a character lookup table, and draws the final colored art.

---

### 5.1 Package and Imports

```go
package ascii
```

This declares that everything in this file belongs to the `ascii` package. In Go, files in the same folder share a package name. Because this says `ascii` and not `main`, it cannot run on its own — it gets imported and called by `main.go`.

**Important Go rule:** only names starting with a **capital letter** are visible outside their package. That is why `ReadBanner`, `BuildAsciiMap`, `PrintAscii`, and `Colors` all start with capitals — so `main.go` can use them.

```go
import (
    // "fmt"
    "os"
    "strings"
)
```

`"fmt"` is **commented out** — it was used during development for debug printing and is no longer needed. The `//` makes Go ignore that line entirely.

| Package | What it gives us |
|---|---|
| `"os"` | `os.ReadFile` — reads an entire file from disk |
| `"strings"` | `strings.Split`, `strings.Index`, `strings.Builder` — text tools |

Below the imports there is a block of commented-out code:

```go
// 	Reand banner files
// 	func Render(text string, banner string) {
// 	bannerLines := readBanner(banner)
// ...
```

This is an old version of the code that was replaced. The `//` at the start of each line makes Go ignore it. It is left in as a record of how the code used to look.

---

### 5.2 ReadBanner

```go
func ReadBanner(file string) ([]string, error) {
```

Declares a function named `ReadBanner`.

- **Input:** `file string` — the name of the font file to open, e.g. `"standard.txt"`
- **Output:** two values at once:
  - `[]string` — a **slice** (ordered list) of strings, one per line of the file
  - `error` — `nil` if everything worked, or a description of what went wrong

Returning two values at once is completely normal in Go. The pattern is always: *"here is the result, and here is whether anything went wrong."*

```go
    data, err := os.ReadFile(file)
```

`os.ReadFile` opens the file, reads every byte, and closes it — all in one call. The `:=` is Go's shorthand for "create these variables and figure out their types automatically." We get back:

- `data` — the raw contents of the file as bytes (`[]byte`)
- `err` — `nil` if it worked, or an error if something failed

```go
    if err != nil {
        return nil, err
    }
```

`nil` means "nothing" in Go — it is the zero value for errors. If `err` is not `nil`, something went wrong (the file is missing, the path is wrong, etc.). We immediately stop and send the error back to whoever called us — `main.go` — so it can decide what to do. We return `nil` for the lines because there is nothing useful to give back.

This **check immediately and return early** style is the standard Go way of handling errors.

```go
    lines := strings.Split(string(data), "\n")
    return lines, nil
```

`string(data)` converts the raw bytes into a readable Go string. `strings.Split(..., "\n")` cuts that string at every real newline character, giving back a `[]string` where each element is one line of the file.

Then we return the list and `nil` for the error — meaning success, nothing went wrong.

---

### 5.3 BuildAsciiMap

```go
func BuildAsciiMap(lines []string) map[rune][]string {
```

Takes the raw list of lines from `ReadBanner` and organizes them into a **lookup table** — give it a character, get back that character's 8 art lines instantly.

The return type `map[rune][]string`:
- **`map`** — a lookup table, like a dictionary
- **`rune`** — Go's type for a single character. Internally just a number: `'A'` is 65, space is 32, `'!'` is 33
- **`[]string`** — a slice of strings, always exactly 8 here (one per art row)

So `asciiMap['H']` gives back the 8 lines that draw the letter H.

```go
    asciiMap := make(map[rune][]string)
```

`make` creates an empty, ready-to-use map. You must use `make` in Go — declaring the variable alone gives you a `nil` map that crashes the moment you try to write to it.

```go
    char := 32
```

A plain integer tracking which ASCII character we are currently on. We start at `32` because that is the space character — the first character stored in the font files. It will count up: 32 = space, 33 = `!`, 34 = `"` … all the way to 126 = `~`.

```go
    for i := 1; i < len(lines); i += 9 {
```

Walks through the font file lines. Three things to understand:

- **`i := 1`** — skip line index 0, the extra blank line at the top of every font file
- **`i += 9`** — step by 9 each time: 8 art lines + 1 blank separator = 9 lines per character
- **`i < len(lines)`** — stop when we run out of lines

Positions visited: `1` (space), `10` (!), `19` ("), `28` (#) … and so on.

```go
        asciiMap[rune(char)] = lines[i : i+8]
        char++
```

`rune(char)` converts the integer (e.g. `65`) into the character `'A'` — this becomes the map key.

`lines[i : i+8]` is a **slice expression** — it extracts exactly 8 elements starting at index `i`, up to but not including `i+8`. Those are the 8 art rows for this character. The 9th line (the blank separator) is skipped automatically because the loop jumps by 9.

`char++` moves to the next ASCII character and the loop continues.

```go
    return asciiMap
```

Returns the completed lookup table. After this, any character can be looked up instantly for its 8 art rows.

---

### 5.4 The Colors Map

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

`var` outside any function creates a **package-level variable** — it lives for the entire run of the program. Because `Colors` starts with a capital letter, `main.go` can also access it as `ascii.Colors`.

The values are **ANSI escape codes** — special character sequences that terminals interpret as formatting commands rather than text to display.

Breaking down `"\033[31m"` (red):
- `\033` — the escape character (decimal 27). Every ANSI sequence starts with this
- `[` — opens the command
- `31` — the color number (31 = red, 32 = green, 33 = yellow, 34 = blue, 35 = magenta, 36 = cyan, 37 = white, 30 = black)
- `m` — closes the command

So writing `\033[31mhello\033[0m` to the terminal displays `hello` in red. The `\033[0m` (`"reset"`) cancels the color — without it, everything printed afterward would also be red.

**`"orange"` uses `\033[38;5;214m`** — orange is not one of the 8 original ANSI colors. The format `38;5;N` selects from an extended 256-color palette. Number 214 is a warm orange. Most modern terminals support it.

**`"magenta"` and `"purple"` share the same code** — ANSI color 35 looks somewhere between the two depending on your terminal, so both names are mapped to it.

**`"reset"` is stored in the same map** — so cancelling color is always `Colors["reset"]`, keeping everything in one place.

---

### 5.5 PrintAscii

This is the most important function. It takes the user's text, looks up each character's art, applies color, and returns the complete finished art as a string.

```go
func PrintAscii(text string, asciiMap map[rune][]string, subStr, colorArg string) string {
```

Four inputs, one output:

| Parameter | Type | What it is |
|---|---|---|
| `text` | `string` | The text to render, e.g. `"kitten"` |
| `asciiMap` | `map[rune][]string` | The lookup table from `BuildAsciiMap` |
| `subStr` | `string` | The word to color, or `""` to color everything |
| `colorArg` | `string` | The color name, e.g. `"red"` |

Note: `subStr, colorArg string` on one line with one `string` at the end is Go shorthand for two parameters of the same type.

The function returns a `string` — the full finished art with all ANSI color codes embedded.

```go
    var result strings.Builder
```

`strings.Builder` is a growing text buffer. Every `result.WriteString(...)` call adds more text to it. At the very end `result.String()` converts everything into one big string. This is much faster than repeatedly concatenating strings with `+`, which creates a brand new string in memory every time.

```go
    lines := strings.Split(text, "\\n")
```

Splits the input on the two-character sequence backslash + letter n — not a real newline. When a user types:

```bash
go run . --color=red "Hello\nWorld"
```

The shell passes Go the string `Hello\nWorld` — literally a `\` and an `n` sitting next to each other, not a real newline character. In a Go string literal `"\\n"` represents exactly those two characters. So this split breaks `"Hello\nWorld"` into `["Hello", "World"]` correctly.

Examples:
```
"Hello\nWorld"    → ["Hello", "World"]       two art blocks
"Hello\n\nWorld"  → ["Hello", "", "World"]   two blocks with a blank gap
"Hello"           → ["Hello"]                single block, no split
```

```go
    // subStr := "kit"
```

A commented-out line left over from early development when `subStr` was hardcoded for testing. The `//` makes Go ignore it completely.

```go
    colorCode := Colors[colorArg]
    reset := Colors["reset"]
```

Look up the two ANSI codes needed throughout the render. If `colorArg` is `""`, then `Colors[""]` returns `""` (Go returns the zero value for missing map keys, not a crash). Writing an empty string to the output does nothing — so the no-color case handles itself.

```go
    for i, line := range lines {
```

Loop through each text segment. `range` on a slice gives the index `i` and the value `line` each iteration. We need `i` for the empty-line check below.

```go
        subIndex := strings.Index(line, subStr)
        // fmt.Println("subIndex:", subIndex, line, len(line), len(subStr), 8 == subIndex+len(subStr))
```

`strings.Index` searches `line` for `subStr` and returns the **byte position** where it first appears. If not found, returns `-1`. If `subStr` is `""`, returns `0`.

The commented-out `fmt.Println` was a debug line used during development to print the state of `subIndex` at each step. It is left in but ignored.

`subIndex` will be compared against the current character position throughout the inner loops to know exactly when to start and stop the color.

```go
        if line == "" {
            if i != 0 {
                result.WriteString("\n")
            }
            continue
        }
```

An empty `line` means the user typed `\n\n` (double newline). After splitting, that creates an empty string between two segments. We write one real newline to create a visual blank gap in the output — but only if this is not the very first segment (`i != 0`), to avoid a stray blank line before anything has been drawn. `continue` skips the rest of the loop and moves to the next segment.

---

#### The Double Loop — How Characters Appear Side by Side

You cannot draw all 8 rows of one character and then all 8 rows of the next — that would stack them vertically:

```
 _
| |
...     ← all of 'h' stacked here
(_)
 _
(_)     ← all of 'i' below it
```

To get them **side by side**, you draw **one row across all characters**, then the next row across all characters. That is what the two nested loops below do.

```go
        for row := 0; row < 8; row++ {
```

The **outer loop** — runs exactly 8 times. Each iteration draws one horizontal stripe across all characters in the current line segment.

```go
            for i, char := range line {
                // fmt.Println(subIndex, i, line[i:])
```

The **inner loop** — runs once per character in `line`. `range` on a string gives:
- `i` — the **byte position** of this character in the string
- `char` — the character itself as a `rune`

The commented `fmt.Println` was another debug line that printed `subIndex`, the current position `i`, and the remaining portion of the line at each step.

> **Note:** `i` here is a new variable that shadows the outer `i` from the segments loop. They are separate — this one belongs only to the character loop. Also note that `i` is a **byte** position. For standard ASCII characters (letters, numbers, punctuation) byte position equals character position. For characters outside ASCII like `é` or `€`, a single character can take 2 or 3 bytes, so `i` can jump by more than 1.

```go
                if i == subIndex {
                    // color from colors map
                    result.WriteString(colorCode)
                }
```

When the current byte position `i` matches `subIndex` (the start of the colored word), we write the color start escape code **before** writing the art for this character. From this point the terminal displays output in the chosen color.

For example: `subStr = "kit"` starts at `subIndex = 2`. When `i` reaches `2`, we write `"\033[31m"` (red), then draw the art for `k` — which appears in red.

```go
                result.WriteString(asciiMap[char][row])
```

Look up this character in the map and write its art for the current row:
- `asciiMap[char]` — returns the `[]string` of 8 art lines for this character
- `[row]` — picks the line for the current row (0 through 7)

```go
                if i == subIndex+len(subStr)-1 || i == len(line)-1 {
                    result.WriteString(reset)
                    // fmt.Println("resetting color:", subIndex, i, line[i:])
                    subIndex = strings.Index(line[i:], subStr) + i
                    // fmt.Println("after resetting color:", subIndex, i, line[i:])
                }
```

This is the condition for **stopping the color** and resetting back to the terminal default.

Two commented debug lines are present — they printed the state before and after the re-seek during development.

The condition has two parts separated by `||` (meaning "or" — if either is true, the block runs):

**Part 1 — `i == subIndex+len(subStr)-1`**

We just drew the last character of `subStr`. For example `"kit"` has length 3. If it starts at position 2, the last character is at `2 + 3 - 1 = 4`. When `i` reaches 4, the word is complete — write the reset code.

**Part 2 — `i == len(line)-1`**

A safety guard. `len(line)-1` is the index of the very last character in the line. This ensures the color is always cancelled at the end of every line, regardless of whether Part 1 fired. Without this, if `subStr` is not found in the line at all, the color could bleed into the next line of output.

After resetting, immediately look for the **next occurrence** of `subStr` in the remaining portion of the line:

```go
                    subIndex = strings.Index(line[i:], subStr) + i
```

- `line[i:]` — the portion of the line from the current position forward
- `strings.Index(line[i:], subStr)` — finds the next match, returns `-1` if none
- `+ i` — adjusts the result back to the position in the full original line

This re-seek means the loop will color again when it reaches the next occurrence.

```go
            }
            subIndex = strings.Index(line, subStr)
```

After finishing all characters for this row, reset `subIndex` back to the **first** occurrence of `subStr` in the full line. This is needed because the next row starts from the beginning of `line` again.

```go
            result.WriteString("\n")
        }
    }
    return result.String()
```

After every row's characters are written, add a real newline to move to the next terminal line. Once all segments and all rows are done, `result.String()` assembles everything in the buffer into one complete string and returns it to `main.go`.

---

## 6. Understanding main.go — Line by Line

`main.go` is where Go starts executing. Its job is to read what the user typed, validate it, and call the rendering functions.

```go
package main
```

Every runnable Go program needs exactly one `package main` with one `func main()`. This is it.

```go
import (
    "ascii-art-color/ascii"
    "fmt"
    "os"
    "strings"
)
```

| Import | Purpose |
|---|---|
| `"ascii-art-color/ascii"` | Our `render.go` functions and variables |
| `"fmt"` | `fmt.Println`, `fmt.Printf`, `fmt.Print` — printing to the terminal |
| `"os"` | `os.Args` — reading command-line arguments |
| `"strings"` | `strings.HasPrefix`, `strings.TrimPrefix` — inspecting the flag |

The path `"ascii-art-color/ascii"` combines the module name from `go.mod` with the subfolder name. After this import, everything exported from `render.go` is accessed with the `ascii.` prefix.

```go
func main() {

    /*
        Valid usages:

        go run . --color=<color> "text"
        go run . --color=<color> <substring> "text"
    */
```

`/* ... */` is a multi-line comment — Go ignores it entirely. It documents the two valid ways to call this program.

---

#### Checking the Argument Count

```go
    if len(os.Args) != 3 && len(os.Args) != 4 {
        fmt.Println("Usage: go run . [OPTION] [STRING]")
        fmt.Println()
        fmt.Println("EX: go run . --color=<color> <substring to be colored> \"something\"")
        return
    }
```

`os.Args` is a slice of everything typed on the command line. **Index 0 is always the program name itself**, so user arguments start at index 1.

For `go run . --color=red "hello"`:
```
os.Args[0] = "./ascii-art-color"   ← always the program
os.Args[1] = "--color=red"
os.Args[2] = "hello"
len(os.Args) = 3
```

For `go run . --color=red kit "a kitten"`:
```
os.Args[0] = "./ascii-art-color"
os.Args[1] = "--color=red"
os.Args[2] = "kit"
os.Args[3] = "a kitten"
len(os.Args) = 4
```

If the count is anything other than 3 or 4, the input is wrong. We print usage and `return` — which exits `main()` and therefore the whole program.

---

#### Validating the Flag Format

```go
    colorArg := os.Args[1]

    if !strings.HasPrefix(colorArg, "--color=") {
        fmt.Println("Usage: go run . [OPTION] [STRING]")
        fmt.Println()
        fmt.Println("EX: go run . --color=<color> <substring to be colored> \"something\"")
        return
    }
```

Store the first user argument and check it starts with exactly `"--color="`. `strings.HasPrefix` returns `true` if it does, `false` if it doesn't. The `!` flips that — if it does **not** start with `"--color="`, print usage and stop.

This catches: `go run . red "hello"` (missing prefix), `go run . --colour=red "hello"` (typo), and `go run . --color red "hello"` (space instead of `=`).

---

#### Extracting the Color Name

```go
    // extract color name
    colorArgs := strings.TrimPrefix(colorArg, "--color=")
```

`strings.TrimPrefix` removes the `"--color="` part from the front and returns what remains. So `"--color=red"` becomes `"red"`. This is the value passed through the rest of the program.

---

#### Validating the Color Exists

```go
    // validate color exists in the Colors map
    if _, ok := ascii.Colors[colorArgs]; !ok {
        fmt.Printf("error: '%s' is not a supported color\n", colorArgs)
        fmt.Println("Supported colors: black, red, green, yellow, blue, magenta, purple, cyan, white, orange")
        return
    }
```

This is the **comma ok idiom** — Go's standard way to check if a key exists in a map.

When you look up a map key in Go you can receive two values: the value stored at that key, and a boolean `ok` that is `true` if the key exists and `false` if it doesn't. We use `_` to discard the actual ANSI string (we don't need it here — just checking existence).

If `ok` is `false`, the color the user typed is not in the `Colors` map. `!ok` means "key was NOT found." We print a clear error message showing which color failed and list the valid options, then stop.

`fmt.Printf` lets you embed values into a string using `%s` as a placeholder for a string. So `fmt.Printf("error: '%s' is not a supported color\n", colorArgs)` would print something like:

```
error: 'purpel' is not a supported color
```

---

#### Parsing the Positional Arguments

```go
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
```

`var` gives both variables the zero value for strings: `""` (empty string).

**Three arguments total** → `go run . --color=red "hello"` → `os.Args[2]` is the text to render, `subStr` stays empty meaning color the whole output.

**Four arguments total** → `go run . --color=red kit "a kitten"` → `os.Args[2]` is the substring to color (`"kit"`), `os.Args[3]` is the full text to render (`"a kitten"`).

---

#### Loading the Font and Rendering

```go
    // read banner
    bannerLines, err := ascii.ReadBanner("standard.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
```

Open and read the font file. If it fails (missing file, wrong directory), print the error and stop.

```go
    asciiMap := ascii.BuildAsciiMap(bannerLines)
```

Turn the raw lines into the character lookup table. Fast, runs once per program call.

```go
    result := ascii.PrintAscii(input, asciiMap, subStr, colorArgs)

    fmt.Print(result)
```

`PrintAscii` does all the rendering and returns the full art as one string. We print it with `fmt.Print` — not `fmt.Println` — because `PrintAscii` already manages all its own newlines internally. Using `Println` would add an unwanted extra blank line at the bottom.

---

## 7. The Coloring Logic — A Full Walkthrough

Let's trace through exactly what happens for a concrete example.

**Command:** `go run . --color=red kit "kitten"`

After `main.go` parses everything:
- `subStr = "kit"`, `input = "kitten"`, `colorArg = "red"`

Inside `PrintAscii`:
```
lines     = ["kitten"]     (no \n in input, so one segment only)
colorCode = "\033[31m"     (ANSI red)
reset     = "\033[0m"
```

For `line = "kitten"`:
```
subIndex = strings.Index("kitten", "kit") = 0
```

Found at position 0 — `k` is the very first character.

Now the double loop. Tracing **row 0** character by character:

```
i=0, char='k':
  i == subIndex? (0 == 0) → YES → write "\033[31m"        ← RED STARTS
  write asciiMap['k'][0]
  i == 0+3-1=2? NO.  i == len("kitten")-1=5? NO.

i=1, char='i':
  i == 0? NO
  write asciiMap['i'][0]
  i == 2? NO.  i == 5? NO.

i=2, char='t':
  i == 0? NO
  write asciiMap['t'][0]
  i == 2? YES → write "\033[0m"                            ← RED STOPS
  re-seek: strings.Index("kitten"[2:], "kit") + 2
         = strings.Index("tten", "kit") + 2
         = -1 + 2 = 1    (not found in remaining text)

i=3, char='t':
  i == 1? NO
  write asciiMap['t'][0]
  i == 5? NO.

i=4, char='e':
  write asciiMap['e'][0]
  i == 5? NO.

i=5, char='n':
  write asciiMap['n'][0]
  i == 5? YES → write "\033[0m"    (safety reset — end of line)

write "\n"

After row 0: subIndex = strings.Index("kitten", "kit") = 0
             (reset back to first match, ready for row 1)
```

Row 0 result in the output buffer:
```
\033[31m[k-row0][i-row0][t-row0]\033[0m[t-row0][e-row0][n-row0]\n
```

The terminal sees `\033[31m`, switches to red, draws `k`, `i`, `t` in red, then sees `\033[0m` and switches back — `kit` appears red, `ten` appears in the default color.

This exact process repeats for rows 1 through 7, with `subIndex` reset to `0` at the start of each row. All 8 rows of `kit` are red, all 8 rows of `ten` are plain — producing correctly side-by-side colored art.

---

## 8. Example Runs

```bash
# Color the entire word
go run . --color=cyan "Hello"

# Color only "ell" inside "Hello"
go run . --color=yellow ell "Hello"

# Color "kit" everywhere it appears
go run . --color=red kit "a kitten has a kit"

# Multi-line art
go run . --color=green "Hello\nWorld"

# Double newline creates a gap between blocks
go run . --color=blue "Hello\n\nWorld"

# Wrong color — prints error message
go run . --color=purpel "hello"

# Too few arguments — prints usage
go run . --color=red

# Missing flag entirely — prints usage
go run . "hello"

# Run all tests
go test ./... -v
```

**Available colors:**

| Name | ANSI Code | Note |
|---|---|---|
| `black` | `\033[30m` | Standard |
| `red` | `\033[31m` | Standard |
| `green` | `\033[32m` | Standard |
| `yellow` | `\033[33m` | Standard |
| `blue` | `\033[34m` | Standard |
| `magenta` | `\033[35m` | Standard |
| `purple` | `\033[35m` | Same code as magenta |
| `cyan` | `\033[36m` | Standard |
| `white` | `\033[37m` | Standard |
| `orange` | `\033[38;5;214m` | 256-color palette — works on most modern terminals |

---

> Built with Go — no external dependencies, just the standard library.