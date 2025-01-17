package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// program
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// let
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	} else {
		out.WriteString("nil")
	}
	out.WriteString(";")

	return out.String()
}

// return
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	} else {
		out.WriteString("nil")
	}
	out.WriteString(";")
	return out.String()
}

// expression statement
type ExpressionStatement struct {
	Token      token.Token // the first token in the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return "nil"
}

// identifier
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// integer literal
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

// prefix expression
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// infix expression
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// infix expression
type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (be *BooleanExpression) expressionNode()      {}
func (be *BooleanExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BooleanExpression) String() string       { return be.Token.Literal }

// Block statement
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// if expression
type IfExpression struct {
	Token       token.Token // the IF token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// Function literal
type FunctionLiteralExpression struct {
	Token      token.Token // the IF token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteralExpression) expressionNode()      {}
func (fl *FunctionLiteralExpression) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteralExpression) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, param := range fl.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// Function call
type FunctionCallExpression struct {
	Token      token.Token // the IF token
	Function   Expression  // identifier or function literal
	Parameters []Expression
}

func (fc *FunctionCallExpression) expressionNode()      {}
func (fc *FunctionCallExpression) TokenLiteral() string { return fc.Token.Literal }
func (fc *FunctionCallExpression) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, param := range fc.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(fc.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")

	return out.String()
}

// String
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// Array
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, element := range al.Elements {
		elements = append(elements, element.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ","))
	out.WriteString("]")

	return out.String()
}

// Index expression
type IndexingExpression struct {
	Token  token.Token
	Index  Expression
	Target Expression
}

func (ie *IndexingExpression) expressionNode()      {}
func (ie *IndexingExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexingExpression) String() string {
	return fmt.Sprintf("%s[%s]", ie.Target.String(), ie.Index.String())
}

// Hash
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for k, v := range hl.Pairs {
		pairs = append(pairs, fmt.Sprintf(`%s: %s`, k.String(), v.String()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ","))
	out.WriteString("}")

	return out.String()
}
