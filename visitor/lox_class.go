package visitor

type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{
		Name:    name,
		Methods: make(map[string]*LoxFunction),
	}
}

func (c *LoxClass) Call(v *AstInterpreter, arguments []interface{}) interface{} {
	instance := &LoxInstance{
		Class: c,
	}
	if initializer, ok := c.Methods["init"]; ok {
		initializer.Bind(instance).Call(v, arguments)
	}
	return instance
}

func (c *LoxClass) Arity() int {
	if initialzer, ok := c.Methods["init"]; ok {
		return initialzer.Arity()
	}
	return 0
}
