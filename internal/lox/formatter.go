package lox

import (
	"fmt"
	"strings"
)

type Formatter struct {
	code  strings.Builder
	scope int
}

func (formatter *Formatter) Format(ast []Node) string {
	formatter.code.Reset()
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
		switch primary.Value.(type) {
		case float64:
			value := primary.Value.(float32)
			formatter.code.WriteString(fmt.Sprintf("%f", value))
		case int:
			value := primary.Value.(int)
			formatter.code.WriteString(fmt.Sprintf("%d", value))

		}
	case "string":
		value := primary.Value.(string)
		formatter.code.WriteString(value)
	}

}

func (formatter *Formatter) visitBinary(binary *Binary) {
	binary.Left.Accept(formatter)
	switch binary.Operation {
	case STAR:
		formatter.code.WriteString(" * ")
	case SLASH:
		formatter.code.WriteString(" / ")
	case PLUS:
		formatter.code.WriteString(" + ")
	case MINUS:
		formatter.code.WriteString(" - ")
	case GREATER:
		formatter.code.WriteString(" > ")
	case GREATEREQUAL:
		formatter.code.WriteString(" >= ")
	case LESS:
		formatter.code.WriteString(" < ")
	case LESSEQUAL:
		formatter.code.WriteString(" <= ")
	case EQUALEQUAL:
		formatter.code.WriteString(" == ")
	case BANGEQUAL:
		formatter.code.WriteString(" != ")
	case AND:
		formatter.code.WriteString(" && ")
	case OR:
		formatter.code.WriteString(" || ")
	}
	binary.Right.Accept(formatter)
}

func (formatter *Formatter) visitUnary(unary *Unary) {
	switch unary.Operation {
	case MINUS:
		formatter.code.WriteString("- ")
	case BANG:
		formatter.code.WriteString("! ")
	}
	unary.Expression.Accept(formatter)
}

func (formatter *Formatter) visitGroup(group *Group) {
	formatter.code.WriteString("(")

	group.Expression.Accept(formatter)

	formatter.code.WriteString(")")
}
func (formatter *Formatter) visitVariable(variable *Variable) {
	name, ok := variable.Identifier.Value.(string)
	if !ok {
		return
	}
	formatter.code.WriteString(name)
}

func (formatter *Formatter) visitThis(*This) {
	formatter.code.WriteString("this")
}

func (formatter *Formatter) visitSuper(super *Super) {
	property, ok := super.Property.Value.(string)
	if !ok {
		return
	}
	formatter.code.WriteString(fmt.Sprintf("super.%s", property))
}

func (formatter *Formatter) visitAssignment(assignment *Assignment) {
	assignment.Identifier.Accept(formatter)
	formatter.code.WriteString(" = ")
	assignment.Value.Accept(formatter)
}

func (formatter *Formatter) visitCall(call *Call) {
	call.Callee.Accept(formatter)
	formatter.code.WriteString("(")

	for i, argument := range call.Argument {
		if i != 0 {
			formatter.code.WriteString(",")
		}
		argument.Accept(formatter)
	}

	formatter.code.WriteString(")")
}

func (formatter *Formatter) visitGetExpr(getExpr *GetExpr) {
	getExpr.Object.Accept(formatter)
	name, ok := getExpr.Property.Value.(string)
	if !ok {
		return
	}
	formatter.code.WriteString(fmt.Sprintf(".%s", name))
}

func (formatter *Formatter) visitExprStmt(exprStmt *ExpressionStmt) {
	formatter.addIndentation()
	exprStmt.Expr.Accept(formatter)
	formatter.code.WriteString(";\n")
}

func (formatter *Formatter) visitPrint(printstmt *PrintStmt) {
	formatter.addIndentation()
	formatter.code.WriteString("print ")
	printstmt.Expr.Accept(formatter)
	formatter.code.WriteString(";\n")
}

func (formatter *Formatter) visitReturn(returnStmt *ReturnStmt) {
	formatter.addIndentation()
	formatter.code.WriteString("returnStmt ")
	returnStmt.Accept(formatter)
	formatter.code.WriteString(";\n")
}
func (formatter *Formatter) visitBlock(block *BlockStmt) {
	formatter.addIndentation()
	formatter.code.WriteString("{\n")
	formatter.scope++

	for _, stmt := range block.Body {
		stmt.Accept(formatter)
	}

	formatter.scope--
	formatter.addIndentation()
	formatter.code.WriteString("}\n")

}
func (formatter *Formatter) visitIf(ifStmt *IfStmt) {
	formatter.addIndentation()
	formatter.code.WriteString("if(")
	ifStmt.Condition.Accept(formatter)
	formatter.code.WriteString(")")
	ifStmt.Then.Accept(formatter)

	if ifStmt.Else != nil {
		formatter.code.WriteString(" else ")
		ifStmt.Else.Accept(formatter)
	}

}

func (formatter *Formatter) visitVarDecl(varDecl *VarDecl) {

	formatter.addIndentation()
	name, ok := varDecl.Identifier.Value.(string)
	if !ok {
		return
	}

	formatter.code.WriteString(fmt.Sprintf("var %s", name))

	if varDecl.Value != nil {
		formatter.code.WriteString(" = ")
		varDecl.Value.Accept(formatter)
	}

	formatter.code.WriteString(";\n")

}

func (formatter *Formatter) visitWhile(while *WhileStmt) {
	formatter.addIndentation()
	formatter.code.WriteString("if(")
	while.Condition.Accept(formatter)
	formatter.code.WriteString(")")
	while.Then.Accept(formatter)
}

func (formatter *Formatter) visitFuncDecl(function *FuncDecl) {
	formatter.addIndentation()
	name, ok := function.Name.Value.(string)
	if !ok {
		return
	}
	formatter.code.WriteString(fmt.Sprintf("func %s(", name))

	for i, param := range function.Parameters {
		if i != 0 {
			formatter.code.WriteString(",")
		}
		param.Accept(formatter)
	}
	formatter.code.WriteString(")")
	function.Body.Accept(formatter)
}

func (formatter *Formatter) visitClassDecl(class *ClassDecl) {
	formatter.addIndentation()
	name, ok := class.Name.Value.(string)
	if !ok {
		return
	}
	formatter.code.WriteString(fmt.Sprintf("class %s ", name))

	if class.Parent != nil {
		name, ok := class.Parent.Value.(string)
		if !ok {
			return
		}
		formatter.code.WriteString(fmt.Sprintf("< %s {\n", name))
	} else {
		formatter.code.WriteString(fmt.Sprintf("{\n"))
	}
	formatter.scope++

	for _, method := range class.Body {
		method.Accept(formatter)
	}
	formatter.scope--
	formatter.addIndentation()
	formatter.code.WriteString(fmt.Sprintf("}\n"))

}

func (formatter *Formatter) addIndentation() {
	for range formatter.scope {
		formatter.code.WriteString("    ")
	}
}
