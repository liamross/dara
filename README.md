# Dara language

This is a purely for-fun dynamically-typed language that is a 🚧WORK IN PROGRESS 🚧.

I'm using it to learn how to write interpreters + compilers for fun. It will use
elements from both Lox (from the book
[Crafting Interpreters](https://craftinginterpreters.com/) and Monkey (from the
book [Writing an Interpreter in Go](https://interpreterbook.com/)), since I am
using these resources to learn.

However as I develop it I hope to diverge from both in order to achieve a
simpler syntax and no semicolons (probably using the same smart placement as Go
uses).

## Usage

If you clone the repo, you can run the _repl_ by compiling to binary, or running
go run main.go. More to come.

## Key features (compared to Monkey and Lox)

- has full number support (like Lox, Monkey only supports integers)
- has `>=` and `<=` (like Lox, Monkey does not include tokens for those )
- supports `//` and `/**/` style comments (Monkey has no comments, Lox only had `//`)
- supports `%` (modulo - not supported by either language)
- no brackets around if conditions (not supported by either language)
- if is a statement rather than an expression in order to support `else if`
  (Monkey has if as an expression which did not easily support `else if`)

> Dara is second in Irish, since this language follows in the footsteps of Lox
> and Monkey.

## Example of valid Dara (currently only the lexer is implemented though, will update with changed plans)

```rust
let five = 5;
let num = 1.234;

let add = fn(a, b) {
    return a + b;
}

if add(five, num) > 2 {
    num = 1;
} else if five > num {
    num = 2;
}

/* The following is just to get all the possible characters out: */
// And both types of comments!

if (1 >= 2 <= 3 > 4 < -(5 % 6) == 7 && true || false) {
    callSomeFunction();
}
```

## Goal syntax of Dara (fingers crossed)

```go
// Using := to indicate initialization, and = to indicate assignment.
five := 5
num := 1.234

num = 1.2345

// No semicolons!
add := fn(a, b) {
    return a + b
}
```
