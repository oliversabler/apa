package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/oliversabler/apa/evaluator"
	"github.com/oliversabler/apa/lexer"
	"github.com/oliversabler/apa/object"
	"github.com/oliversabler/apa/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}
