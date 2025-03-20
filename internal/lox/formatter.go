package lox

import (
	"fmt"
	"strings"
)

type Formatter struct {
	code         strings.Builder
	scope        int
	stopNewLines bool
	queueNewLine bool
	lastWrite    string
}

func (formatter *Formatter) write(code string) {
	if formatter.queueNewLine && !strings.HasPrefix(code, "\n") && !strings.HasSuffix(formatter.lastWrite, "\n") {
		formatter.code.WriteString("\n")
	}
	formatter.queueNewLine = false
	formatter.code.WriteString(code)
	formatter.lastWrite = code
}

func (formatter *Formatter) visitComment(comment *Comment) {
	value, ok := comment.Comment.Value.(string)
	if !ok {
		return
	}
	formatter.write(fmt.Sprintf("//%s", value))
	formatter.addNewLine()
}

func (formatter *Formatter) visitNewLine(*NewLine) {
	formatter.addIndentation()
	formatter.write(fmt.Sprintf("\n"))
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
		formatter.write("nil")
	case "boolean":
		if primary.Value == false {
			formatter.write("false")
		} else if primary.Value == true {
			formatter.write("true")
		}
	case "number":
		switch primary.Value.(type) {
		case float64:
			value := primary.Value.(float32)
			formatter.write(fmt.Sprintf("%f", value))
		case int:
			value := primary.Value.(int)
			formatter.write(fmt.Sprintf("%d", value))

		}
	case "string":
		value := primary.Value.(string)
		formatter.write(fmt.Sprintf("\"%s\"", value))
	}

}

func (formatter *Formatter) visitBinary(binary *Binary) {
	binary.Left.Accept(formatter)
	switch binary.Operation {
	case STAR:
		formatter.write(" * ")
	case SLASH:
		formatter.write(" / ")
	case PLUS:
		formatter.write(" + ")
	case MINUS:
		formatter.write(" - ")
	case GREATER:
		formatter.write(" > ")
	case GREATEREQUAL:
		formatter.write(" >= ")
	case LESS:
		formatter.write(" < ")
	case LESSEQUAL:
		formatter.write(" <= ")
	case EQUALEQUAL:
		formatter.write(" == ")
	case BANGEQUAL:
		formatter.write(" != ")
	case AND:
		formatter.write(" && ")
	case OR:
		formatter.write(" || ")
	}
	binary.Right.Accept(formatter)
}

func (formatter *Formatter) visitUnary(unary *Unary) {
	switch unary.Operation {
	case MINUS:
		formatter.write("-")
	case BANG:
		formatter.write("!")
	}
	unary.Expression.Accept(formatter)
}

func (formatter *Formatter) visitGroup(group *Group) {
	formatter.write("(")

	group.Expression.Accept(formatter)

	formatter.write(")")
}
func (formatter *Formatter) visitVariable(variable *Variable) {
	name, ok := variable.Identifier.Value.(string)
	if !ok {
		return
	}
	formatter.write(name)
}

func (formatter *Formatter) visitThis(*This) {
	formatter.write("this")
}

func (formatter *Formatter) visitSuper(super *Super) {
	property, ok := super.Property.Value.(string)
	if !ok {
		return
	}
	formatter.write(fmt.Sprintf("super.%s", property))
}

func (formatter *Formatter) visitAssignment(assignment *Assignment) {
	assignment.Identifier.Accept(formatter)
	formatter.write(" = ")
	assignment.Value.Accept(formatter)
}

func (formatter *Formatter) visitCall(call *Call) {
	call.Callee.Accept(formatter)
	formatter.write("(")

	for i, argument := range call.Argument {
		if i != 0 {
			formatter.write(",")
		}
		argument.Accept(formatter)
	}

	formatter.write(")")
}

func (formatter *Formatter) visitGetExpr(getExpr *GetExpr) {
	getExpr.Object.Accept(formatter)
	name, ok := getExpr.Property.Value.(string)
	if !ok {
		return
	}
	formatter.write(fmt.Sprintf(".%s", name))
}

func (formatter *Formatter) visitExprStmt(exprStmt *ExpressionStmt) {
	formatter.addIndentation()
	exprStmt.Expr.Accept(formatter)
	formatter.write(";")
	formatter.addNewLine()
}

func (formatter *Formatter) visitPrint(printstmt *PrintStmt) {
	formatter.addIndentation()
	formatter.write("print ")
	printstmt.Expr.Accept(formatter)
	formatter.write(";")
	formatter.addNewLine()
}

func (formatter *Formatter) visitReturn(returnStmt *ReturnStmt) {
	formatter.addIndentation()
	if returnStmt.ReturnsValue {
		formatter.write("return ")
		returnStmt.Expr.Accept(formatter)
		formatter.write(";")
		formatter.addNewLine()
	} else {
		formatter.write("return;")
		formatter.addNewLine()
	}
}

func (formatter *Formatter) visitBlock(block *BlockStmt) {
	if block.BlockContext == BLOCK_CONTEXT {
		formatter.addIndentation()
	}
	formatter.write("{")
	formatter.addNewLine()
	formatter.scope++

	for _, stmt := range block.Body {
		stmt.Accept(formatter)
	}

	formatter.scope--
	formatter.addIndentation()
	formatter.write("}")
	formatter.addNewLine()

}
func (formatter *Formatter) visitIf(ifStmt *IfStmt) {
	formatter.addIndentation()
	formatter.write("if (")
	ifStmt.Condition.Accept(formatter)
	formatter.write(") ")
	ifStmt.Then.Accept(formatter)

	if ifStmt.Else != nil {
		formatter.write(" else ")
		ifStmt.Else.Accept(formatter)
	}

}

func (formatter *Formatter) visitVarDecl(varDecl *VarDecl) {

	formatter.addIndentation()
	name, ok := varDecl.Identifier.Value.(string)
	if !ok {
		return
	}

	formatter.write(fmt.Sprintf("var %s", name))

	if varDecl.Initialized {
		formatter.write(" = ")
		varDecl.Value.Accept(formatter)
	}

	formatter.write(";")
	formatter.addNewLine()

}

func (formatter *Formatter) visitWhile(while *WhileStmt) {
	formatter.addIndentation()
	formatter.write("while (")
	while.Condition.Accept(formatter)
	formatter.write(") ")
	while.Then.Accept(formatter)
}

func (formatter *Formatter) visitFor(forStmt *ForStmt) {
	formatter.addIndentation()
	formatter.write("for (")
	formatter.stopNewLines = true
	if forStmt.Initializer != nil {
		forStmt.Initializer.Accept(formatter)
		formatter.write(" ")
	} else {
		formatter.write("; ")
	}
	if forStmt.Condition != nil {
		forStmt.Condition.Accept(formatter)
	}
	formatter.write("; ")
	if forStmt.Assignment != nil {
		forStmt.Assignment.Accept(formatter)
	}
	formatter.write(") ")
	formatter.stopNewLines = false
	forStmt.Body.Accept(formatter)
}

func (formatter *Formatter) visitFuncDecl(function *FuncDecl) {
	formatter.addIndentation()
	name, ok := function.Name.Value.(string)
	if !ok {
		return
	}
	if function.FunctionType == METHOD_CONTEXT {
		formatter.write(fmt.Sprintf("%s(", name))
	} else {
		formatter.write(fmt.Sprintf("fun %s(", name))
	}
	for i, param := range function.Parameters {
		if i != 0 {
			formatter.write(",")
		}
		param.Accept(formatter)
	}
	formatter.write(") ")
	function.Body.Accept(formatter)
}

func (formatter *Formatter) visitClassDecl(class *ClassDecl) {
	formatter.addIndentation()
	name, ok := class.Name.Value.(string)
	if !ok {
		return
	}
	formatter.write(fmt.Sprintf("class %s ", name))

	if class.Parent != nil {
		name, ok := class.Parent.Value.(string)
		if !ok {
			return
		}
		formatter.write(fmt.Sprintf("< %s {", name))
	} else {
		formatter.write(fmt.Sprintf("{"))
	}
	formatter.addNewLine()
	formatter.scope++

	for _, method := range class.Body {
		method.Accept(formatter)
	}
	formatter.scope--
	formatter.addIndentation()
	formatter.write(fmt.Sprintf("}"))
	formatter.addNewLine()

}

func (formatter *Formatter) addIndentation() {
	for range formatter.scope {
		formatter.code.WriteString("    ")
	}
}

func (formatter *Formatter) addNewLine() {
	if formatter.stopNewLines {
		return
	}
	formatter.queueNewLine = true
}
