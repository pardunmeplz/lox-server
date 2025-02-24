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

type ExpressionStmt struct {
	Expr Node
}

func (expr *ExpressionStmt) accept(visitor Visitor) {

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
