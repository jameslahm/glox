package glox

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jameslahm/glox/ast"
	"github.com/jameslahm/glox/lexer"
	"github.com/jameslahm/glox/visitor"
)

type Glox struct {
}

func (g *Glox) RunFile(path string) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("ReadFile error: %v\n", err)
	}
	g.Run(string(buf))
}

func (g *Glox) RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")
	for scanner.Scan() {
		line := scanner.Text()
		g.Run(line)
		fmt.Print(">> ")
	}
}

func (g *Glox) Run(script string) {
	lex := lexer.NewLexer(script)
	lex.Lex()
	parser := ast.NewParser(lex.Tokens)
	node := parser.Parse()

	if len(parser.Errors) != 0 {
		fmt.Println("Error: parse failed")
		return
	}

	resolver := visitor.NewResolver()
	node.Accept(resolver)

	if len(resolver.Errors) != 0 {
		fmt.Println("Error: resolve failed")
		return
	}

	interpreter := visitor.NewAstInterpreter(resolver.VariableBindingDistances)

	node.Accept(interpreter)
}
