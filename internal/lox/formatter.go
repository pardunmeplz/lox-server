package lox

import (
	"fmt"
	"strings"
)

type Formatter struct {
	code strings.Builder
}

func (formatter *Formatter) Format(ast []Node) string {
	for _, node := range ast {
		node.Accept(formatter)
	}
	return formatter.code.String()

}

func (formatter *Formatter) visitPrimary(primary *Primary) {
	switch primary.ValType {
	case "nil":
		formatter.code.WriteString("nil")
	case "boolean":
		if primary.Value == false {
			formatter.code.WriteString("false")
		} else if primary.Value == true {
			formatter.code.WriteString("true")
		}
	case "number":
		value := primary.Value.(float64)
		formatter.code.WriteString(fmt.Sprintf("%f ", value))
	case "string":
		value := primary.Value.(string)
		formatter.code.WriteString(value)
	}

}

func (formatter *Formatter) visitBinary(*Binary) {

}
func (formatter *Formatter) visitUnary(*Unary) {

}
func (formatter *Formatter) visitGroup(*Group) {

}
func (formatter *Formatter) visitVariable(*Variable) {

}
func (formatter *Formatter) visitThis(*This) {

}
func (formatter *Formatter) visitSuper(*Super) {

}
func (formatter *Formatter) visitAssignment(*Assignment) {

}
func (formatter *Formatter) visitCall(*Call) {

}
func (formatter *Formatter) visitGetExpr(*GetExpr) {

}
func (formatter *Formatter) visitExprStmt(*ExpressionStmt) {

}
func (formatter *Formatter) visitPrint(*PrintStmt) {

}
func (formatter *Formatter) visitReturn(*ReturnStmt) {

}
func (formatter *Formatter) visitBlock(*BlockStmt) {

}
func (formatter *Formatter) visitIf(*IfStmt) {

}
func (formatter *Formatter) visitVarDecl(*VarDecl) {

}
func (formatter *Formatter) visitWhile(*WhileStmt) {

}
func (formatter *Formatter) visitFuncDecl(*FuncDecl) {

}
func (formatter *Formatter) visitClassDecl(*ClassDecl) {

}
