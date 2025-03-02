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

type ReturnStmt struct {
	Expr Node
}

func (expr *ReturnStmt) accept(visitor Visitor) {
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

type VarDecl struct {
	Identifier Token
	Value      Node
}

func (expr *VarDecl) accept(visitor Visitor) {
}

type WhileStmt struct {
	Condition Node
	Then      Node
}

func (expr *WhileStmt) accept(visitor Visitor) {
}

type FuncDecl struct {
	Name       Token
	Body       Node
	Parameters []Node
}

func (expr *FuncDecl) accept(visitor Visitor) {
}

type ClassDecl struct {
	Name   Token
	Parent *Token
	Body   []Node
}

func (expr *ClassDecl) accept(visitor Visitor) {
}
