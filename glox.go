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
	interpreter := visitor.NewAstInterpreter()
	node := parser.Parse()
	node.Accept(interpreter)
}
