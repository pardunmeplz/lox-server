package lox

type ExpressionStmt struct {
	Expr Node
}

func (expr *ExpressionStmt) accept(visitor Visitor) {
}

type PrintStmt struct {
	Expr Node
}

func (expr *PrintStmt) accept(visitor Visitor) {
}

type BlockStmt struct {
	Body []Node
}

func (expr *BlockStmt) accept(visitor Visitor) {
}

type IfStmt struct {
	Condition Node
	Then      Node
	Else      Node
}

func (expr *IfStmt) accept(visitor Visitor) {
}
