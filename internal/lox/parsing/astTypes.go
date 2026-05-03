package parsing

import (
	// lox_dg "lox-server/internal/lox/diagnostics"
	lox_sc "lox-server/internal/lox/lexical_analysis"
)

type Node interface {
	Accept(Visitor)
}

type Visitor interface {
	visitPrimary(*Primary)
	visitBinary(*Binary)
	visitUnary(*Unary)
	visitGroup(*Group)
	visitVariable(*Variable)
	visitThis(*This)
	visitSuper(*Super)
	visitAssignment(*Assignment)
	visitCall(*Call)
	visitGetExpr(*GetExpr)
	visitExprStmt(*ExpressionStmt)
	visitPrint(*PrintStmt)
	visitReturn(*ReturnStmt)
	visitBlock(*BlockStmt)
	visitIf(*IfStmt)
	visitVarDecl(*VarDecl)
	visitWhile(*WhileStmt)
	visitFor(*ForStmt)
	visitFuncDecl(*FuncDecl)
	visitClassDecl(*ClassDecl)
	visitNewLine(*NewLine)
	visitComment(*Comment)
}

type Comment struct {
	Token  *lox_sc.Token
	Inline bool
}

func (expr *Comment) Accept(visitor Visitor) {
	visitor.visitComment(expr)
}

type Primary struct {
	Value any
	// ValType string
	Token *lox_sc.Token
}

func (expr *Primary) Accept(visitor Visitor) {
	visitor.visitPrimary(expr)
}

type Binary struct {
	Left      Node
	Right     Node
	Operation int
}

func (expr *Binary) Accept(visitor Visitor) {
	visitor.visitBinary(expr)
}

type Unary struct {
	Expression Node
	Operation  int
}

func (expr *Unary) Accept(visitor Visitor) {
	visitor.visitUnary(expr)
}

type Group struct {
	OpenParan  *lox_sc.Token
	CloseParan *lox_sc.Token
	Expression Node
}

func (expr *Group) Accept(visitor Visitor) {
	visitor.visitGroup(expr)
}

type Variable struct {
	Identifier *lox_sc.Token
	Definition *lox_sc.Token
}

func (expr *Variable) Accept(visitor Visitor) {
	visitor.visitVariable(expr)
}

type This struct {
	Identifier *lox_sc.Token
}

func (expr *This) Accept(visitor Visitor) {
	visitor.visitThis(expr)
}

type Super struct {
	Identifier *lox_sc.Token
	Super      *lox_sc.Token
}

func (expr *Super) Accept(visitor Visitor) {
	visitor.visitSuper(expr)
}

type Assignment struct {
	Value      Node
	Identifier Node
}

func (expr *Assignment) Accept(visitor Visitor) {
	visitor.visitAssignment(expr)
}

type Call struct {
	Callee   Node
	Argument []Node
}

func (expr *Call) Accept(visitor Visitor) {
	visitor.visitCall(expr)
}

type GetExpr struct {
	Object   Node
	Property *lox_sc.Token
}

func (expr *GetExpr) Accept(visitor Visitor) {
	visitor.visitGetExpr(expr)
}
