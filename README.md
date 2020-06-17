# Dara language

A super simple dynamic language with a focus on conciseness and readability.

## Why

Firstly, to learn how to write interpreters and compilers! The core architecture
is being written while learning (primarily) from three sources:

- [Writing an Interpreter in Go](https://interpreterbook.com/)
- [Crafting Interpreters](https://craftinginterpreters.com/)
- [Go Source Code](https://github.com/golang/go/tree/master/src/go)

Secondly, to build a dynamic language syntax and utility that I would personally
like to use. It draws heavily on Go, JavaScript and Rust syntax (in that order)
and aims to reduce the amount of code needed to achieve things, while
maintaining high levels of readability.

## Status

### Completed

- [x] Lexer
- [x] Error reporting with line-level accuracy
- [x] Most of the parser
- [x] Complete _REPL_ (only works on values implemented in evaluator though)

### In progress

- [ ] Implement evaluator
- [ ] Implement parsing objects and arrays

### To Do (roughly in order)

- [ ] Add line numbers to evaluator error reporting
- [ ] **Remove all semicolons**
- [ ] Build a compiler (stretch goal)
- [ ] Improve all error messaging, and have column-level accuracy (stretch goal)

## Usage

If you clone the repo, you can run the _REPL_ by compiling to binary, or running
`go run main.go`. More to come.

## Current Valid Dara (subject to change wildly)

```go
// Declare values with `:=` (no declaration keyword).
five := 5;
num := 1.234;

// Dara uses `nil` to indicate the absence of a value. If you want to declare a
// variable without assigning a value, use `:= nil`.
// other;     // (not allowed)
other := nil; // allowed

// Assign values to existing identifiers using `=`. Functions are values.
add = fn(a, b) {
    return a + b;
}

// No brackets around the logic in if statements. No truthy or falsy values,
// must use booleans in if statements.
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
function := fn(a, b) { return a + b; };

// Available logical operators:
// < > ! == != >= <= && ||   (work on strings: < > == != >= <=)

// Available arithmetic operators:
//  + - * / %                (work on strings: +)
```
