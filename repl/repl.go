package repl

import (
	"bufio"
	"dara/lexer"
	"dara/token"
	"fmt"
	"io"
)

const PROMPT = "-> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		var (
			line = scanner.Text()
			l    = lexer.New(line)
		)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%-8v (%v)\n", tok.Type, tok.Literal)
		}
	}
}
