package visitor

import (
	"time"

	"github.com/jameslahm/glox/ast"
)

type Clock struct {
}

func (c *Clock) Call(v ast.Visitor, arguments []interface{}) interface{} {
	return time.Now().Unix()
}

func (c *Clock) Arity() int {
	return 0
}
