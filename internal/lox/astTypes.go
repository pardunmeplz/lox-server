package lox

type Node interface {
	accept(Visitor)
}

type Visitor interface {
}

type Primary struct {
	Value   any
	ValType string
}

func (expr *Primary) accept(visitor Visitor) {
}

type Binary struct {
	Left      Node
	Right     Node
	Operation int
}

func (expr *Binary) accept(visitor Visitor) {
}

type Unary struct {
	Expression Node
	Operation  int
}

func (expr *Unary) accept(visitor Visitor) {
}

type Group struct {
	Expression Node
}

func (expr *Group) accept(visitor Visitor) {
}

type Variable struct {
	Identifier Token
}

func (expr *Variable) accept(visitor Visitor) {
}

type This struct {
	Identifier Token
}

func (expr *This) accept(visitor Visitor) {
}

type Super struct {
	Identifier Token
}

func (expr *Super) accept(visitor Visitor) {
}

type Assignment struct {
	Value      Node
	Identifier Node
}

func (expr *Assignment) accept(visitor Visitor) {
}

type Call struct {
	Callee   Node
	Argument []Node
}

func (expr *Call) accept(visitor Visitor) {
}
