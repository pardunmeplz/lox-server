package lox

type ExpressionStmt struct {
	Expr Node
}

func (expr *ExpressionStmt) accept(visitor Visitor) {
	visitor.visitExprStmt(expr)
}

type PrintStmt struct {
	Expr Node
}

func (expr *PrintStmt) accept(visitor Visitor) {
	visitor.visitPrint(expr)
}

type ReturnStmt struct {
	Expr Node
}

func (expr *ReturnStmt) accept(visitor Visitor) {
	visitor.visitReturn(expr)
}

type BlockStmt struct {
	Body []Node
}

func (expr *BlockStmt) accept(visitor Visitor) {
	visitor.visitBlock(expr)
}

type IfStmt struct {
	Condition Node
	Then      Node
	Else      Node
}

func (expr *IfStmt) accept(visitor Visitor) {
	visitor.visitIf(expr)
}

type VarDecl struct {
	Identifier Token
	Value      Node
}

func (expr *VarDecl) accept(visitor Visitor) {
	visitor.visitVarDecl(expr)
}

type WhileStmt struct {
	Condition Node
	Then      Node
}

func (expr *WhileStmt) accept(visitor Visitor) {
	visitor.visitWhile(expr)
}

type FuncDecl struct {
	Name       Token
	Body       Node
	Parameters []Node
}

func (expr *FuncDecl) accept(visitor Visitor) {
	visitor.visitFuncDecl(expr)
}

type ClassDecl struct {
	Name   Token
	Parent *Token
	Body   []Node
}

func (expr *ClassDecl) accept(visitor Visitor) {
	visitor.visitClassDecl(expr)
}
