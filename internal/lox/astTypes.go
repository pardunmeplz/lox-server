package lox

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
	visitFuncDecl(*FuncDecl)
	visitClassDecl(*ClassDecl)
}

type Primary struct {
	Value   any
	ValType string
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
	Expression Node
}

func (expr *Group) Accept(visitor Visitor) {
	visitor.visitGroup(expr)
}

type Variable struct {
	Identifier Token
	Definition Token
}

func (expr *Variable) Accept(visitor Visitor) {
	visitor.visitVariable(expr)
}

type This struct {
	Identifier Token
}

func (expr *This) Accept(visitor Visitor) {
	visitor.visitThis(expr)
}

type Super struct {
	Identifier Token
	Property   Token
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
	Property Token
}

func (expr *GetExpr) Accept(visitor Visitor) {
	visitor.visitGetExpr(expr)
}
