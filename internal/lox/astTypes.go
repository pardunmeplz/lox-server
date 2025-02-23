package lox

type Expr interface {
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

type Expression struct {
}

func (expr *Expression) accept(visitor Visitor) {

}

type Binary struct {
	Left      Expr
	Right     Expr
	Operation rune
}

func (expr *Binary) accept(visitor Visitor) {

}

type Unary struct {
	Expression Expr
	Operation  rune
}

func (expr *Unary) accept(visitor Visitor) {

}
