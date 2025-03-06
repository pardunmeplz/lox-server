package lox

type Node interface {
	accept(Visitor)
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

func (expr *Primary) accept(visitor Visitor) {
	visitor.visitPrimary(expr)
}

type Binary struct {
	Left      Node
	Right     Node
	Operation int
}

func (expr *Binary) accept(visitor Visitor) {
	visitor.visitBinary(expr)
}

type Unary struct {
	Expression Node
	Operation  int
}

func (expr *Unary) accept(visitor Visitor) {
	visitor.visitUnary(expr)
}

type Group struct {
	Expression Node
}

func (expr *Group) accept(visitor Visitor) {
	visitor.visitGroup(expr)
}

type Variable struct {
	Identifier Token
	Definition Token
}

func (expr *Variable) accept(visitor Visitor) {
	visitor.visitVariable(expr)
}

type This struct {
	Identifier Token
}

func (expr *This) accept(visitor Visitor) {
	visitor.visitThis(expr)
}

type Super struct {
	Identifier Token
	Property   Token
}

func (expr *Super) accept(visitor Visitor) {
	visitor.visitSuper(expr)
}

type Assignment struct {
	Value      Node
	Identifier Node
}

func (expr *Assignment) accept(visitor Visitor) {
	visitor.visitAssignment(expr)
}

type Call struct {
	Callee   Node
	Argument []Node
}

func (expr *Call) accept(visitor Visitor) {
	visitor.visitCall(expr)
}

type GetExpr struct {
	Object   Node
	Property Token
}

func (expr *GetExpr) accept(visitor Visitor) {
	visitor.visitGetExpr(expr)
}
