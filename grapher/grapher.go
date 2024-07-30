package grapher

import (
	"bytes"
	"fmt"
	"log"
	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/google/uuid"
)

type Grapher struct {
	Parser *parser.Parser
}

func New(input string) *Grapher {
	l := lexer.New(input)
	p := parser.New(l)
	grapher := &Grapher{Parser: p}
	return grapher
}

func (g *Grapher) GetDot() string {
	program := g.Parser.ParseProgram()
	graphviz := graphviz.New()
	graph, err := graphviz.Graph()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		graphviz.Close()
	}()

	root, err := graph.CreateNode("program\n" + program.String())
	if err != nil {
		log.Fatal("Error creating graph node " + err.Error())
	}
	evalGraph(graph, program, root, "")

	var buf bytes.Buffer
	if err := graphviz.Render(graph, "dot", &buf); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(buf.String())
}

func evalGraph(graph *cgraph.Graph, ast_node ast.Node, parent *cgraph.Node, edgeLabel string) {
	var graph_node *cgraph.Node

	switch ast_node := ast_node.(type) {
	case *ast.Program:
		for _, stmt := range ast_node.Statements {
			evalGraph(graph, stmt, parent, "statement")
		}
		return

	case *ast.LetStatement:
		n, err := graph.CreateNode("LET_STATEMENT\n" + ast_node.String())
		graph_node = n
		if err != nil {
			fmt.Printf("Error creating graph node " + err.Error())
			return
		}
		evalGraph(graph, ast_node.Name, graph_node, "Name")
		evalGraph(graph, ast_node.Value, graph_node, "Value")

	case *ast.FunctionLiteralExpression:
		n, err := graph.CreateNode("FUNCTION_LITERAL\n" + ast_node.String())
		graph_node = n
		if err != nil {
			fmt.Printf("Error creating graph node " + err.Error())
			return
		}
		for _, param := range ast_node.Parameters {
			evalGraph(graph, param, graph_node, "Parameter")
		}
		evalGraph(graph, ast_node.Body, graph_node, "Body")

	case *ast.Identifier:
		n, err := graph.CreateNode("IDENTIFIER\n" + ast_node.String())
		graph_node = n
		if err != nil {
			fmt.Printf("Error creating graph node " + err.Error())
			return
		}

	case *ast.IntegerLiteral:
		n, err := graph.CreateNode("INTEGER_LITERAL\n" + ast_node.String())
		graph_node = n
		if err != nil {
			fmt.Printf("Error creating graph node " + err.Error())
			return
		}

	case *ast.BlockStatement:
		n, err := graph.CreateNode("BLOCK_STATEMENT\n" + ast_node.String())
		graph_node = n
		if err != nil {
			fmt.Printf("Error creating graph node " + err.Error())
			return
		}
		for _, stmt := range ast_node.Statements {
			evalGraph(graph, stmt, graph_node, "statement")
		}

	case *ast.ExpressionStatement:
		n, err := graph.CreateNode("EXPRESSION_STATEMENT\n" + ast_node.String())
		graph_node = n
		if err != nil {
			fmt.Printf("Error creating graph node " + err.Error())
			return
		}
		evalGraph(graph, ast_node.Expression, graph_node, "Expression")

	case *ast.FunctionCallExpression:
		n, err := graph.CreateNode("FUNCTION_CALL\n" + ast_node.String())
		graph_node = n
		if err != nil {
			fmt.Printf("Error creating graph node " + err.Error())
			return
		}
		for _, param := range ast_node.Parameters {
			evalGraph(graph, param, graph_node, "Parameter")
		}
		evalGraph(graph, ast_node.Function, graph_node, "Function")

	case *ast.InfixExpression:
		n, err := graph.CreateNode(fmt.Sprintf("INFIX_EXPRESSION\nOperator: %s\n%s", ast_node.Operator, ast_node.String()))
		graph_node = n
		if err != nil {
			fmt.Printf("Error creating graph node " + err.Error())
			return
		}
		evalGraph(graph, ast_node.Left, graph_node, "Left")
		evalGraph(graph, ast_node.Right, graph_node, "Right")

	default:
		n, err := graph.CreateNode(fmt.Sprintf("%T\n%s", ast_node, ast_node.String()))
		graph_node = n
		if err != nil {
			fmt.Printf("Error creating graph node " + err.Error())
			return
		}
	}

	e, err := graph.CreateEdge(uuid.New().String(), parent, graph_node)
	if err != nil {
		fmt.Printf("Error creating graph edge " + err.Error())
		return
	}
	e.SetLabel(edgeLabel)
}
