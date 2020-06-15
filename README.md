# Dara language

A super simple dynamic language with a focus on simplicity and writing speed.

## Why

This is a purely for-fun language that I am building to learn how to write a
fully-featured interpreter and compiler. The core architecture was written while
learning from three sources:

- [Writing an Interpreter in Go](https://interpreterbook.com/)
- [Crafting Interpreters](https://craftinginterpreters.com/)
- [Go Source Code](https://github.com/golang/go/tree/master/src/go)

However, it differs from all these languages in syntax and features (see below).

## Completed

- [x] Lexer
- [x] Error reporting with line-level accuracy
- [x] Most of the parser
- [x] Basic REPL (only prints parser results at the moment, so RPPL?)

## TODO

- [ ] Add line numbers to evaluator error reporting (and then columns eventually
  - see below)
- [ ] Implement parsing objects and arrays
- [ ] **Remove all semicolons**
- [ ] Build an evaluator
- [ ] Complete REPL
- [ ] Build a compiler (stretch goal)
- [ ] Improve error messaging, and have column-level accuracy (stretch goal)

## Usage

If you clone the repo, you can run the _repl_ by compiling to binary, or running
`go run main.go`. More to come.

## Current Valid Dara (subject to change wildly)

```go
// Declare values with `:=` (no declaration keyword).
five := 5;
num := 1.234;

// Dara uses `nil` to indicate the absence of a value.
other := nil;

// Can declare an identifier without assigning a value. Value will be `nil`.
add;

// Assign values to existing identifiers using `=`. Functions are values.
add = fn(a, b) {
    return a + b;
}

// No brackets around the logic in if statements.
if 1 > 2 {
    num = 1;
} else if five > 2 {
    num = 2;
}

/* Dara also allows multi-line comments using c-style syntax. */

// Available types:
noValue := nil;
string := "string";
number := 1.234;
function := fn(a, b) { return a + b };

// Available logical operators:
// < > ! == != >= <= && || (work on strings: < > == != >= <=)

// Available arithmetic operators:
//  + - * / % (work on strings: +)
```

## Other Language Quirks

- no truthy or falsy values (must use explicit booleans)
