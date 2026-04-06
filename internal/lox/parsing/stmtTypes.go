package parsing

import lox_sc "lox-server/internal/lox/lexical_analysis"

type ExpressionStmt struct {
	Expr Node
}

func (expr *ExpressionStmt) Accept(visitor Visitor) {
	visitor.visitExprStmt(expr)
}

type PrintStmt struct {
	Expr Node
}

func (expr *PrintStmt) Accept(visitor Visitor) {
	visitor.visitPrint(expr)
}

type ReturnStmt struct {
	Expr         Node
	ReturnsValue bool
}

func (expr *ReturnStmt) Accept(visitor Visitor) {
	visitor.visitReturn(expr)
}

type BlockStmt struct {
	Body         []Node
	BlockContext int
}

func (expr *BlockStmt) Accept(visitor Visitor) {
	visitor.visitBlock(expr)
}

type IfStmt struct {
	Condition Node
	Then      Node
	Else      Node
}

func (expr *IfStmt) Accept(visitor Visitor) {
	visitor.visitIf(expr)
}

type VarDecl struct {
	Identifier  lox_sc.Token
	Value       Node
	Initialized bool
}

func (expr *VarDecl) Accept(visitor Visitor) {
	visitor.visitVarDecl(expr)
}

type ForStmt struct {
	Initializer Node
	Condition   Node
	Assignment  Node
	Body        Node
}

func (expr *ForStmt) Accept(visitor Visitor) {
	visitor.visitFor(expr)
}

type WhileStmt struct {
	Condition Node
	Then      Node
}

func (expr *WhileStmt) Accept(visitor Visitor) {
	visitor.visitWhile(expr)
}

type FuncDecl struct {
	Name         lox_sc.Token
	Body         Node
	Parameters   []Node
	FunctionType int
}

func (expr *FuncDecl) Accept(visitor Visitor) {
	visitor.visitFuncDecl(expr)
}

type ClassDecl struct {
	Name   lox_sc.Token
	Parent *lox_sc.Token
	Body   []Node
}

func (expr *ClassDecl) Accept(visitor Visitor) {
	visitor.visitClassDecl(expr)
}

type NewLine struct {
	Token lox_sc.Token
}

func (expr *NewLine) Accept(visitor Visitor) {
	visitor.visitNewLine(expr)
}
