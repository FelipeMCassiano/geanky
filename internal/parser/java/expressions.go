package java

import (
	"fmt"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Expression interface {
	isExpression()
}
type Assignment struct {
	Left  Expression `json:"left"`
	Right Expression `json:"right"`
}
type Binary struct {
	Left     Expression `json:"left"`
	Operator string     `json:"operator"`
	Right    Expression `json:"right"`
}
type IfNode struct {
	Condition Expression `json:"condition"`
}

type MethodInvocation struct {
	Accessed Access       `json:"accessed"`
	Args     []Expression `json:"args"`
}

type Identifier struct {
	Name string `json:"name"`
}

type Access struct {
	Object     string     `json:"object"`
	Identifier Identifier `json:"identifier"`
}
type Literal struct {
	Value string `json:"value"`
}

type ExpressionHandler func(node *tree_sitter.Node, content []byte) (Expression, error)

var expressionHandlers map[string]ExpressionHandler

// A função init() é executada automaticamente pelo Go antes do main()
func init() {
	expressionHandlers = map[string]ExpressionHandler{
		"assignment_expression":    parseAssignment,
		"binary_expression":        parseBinary,
		"parenthesized_expression": parseParenthesized,
		"method_invocation":        parseMethodInvocation,
		"string_literal":           parseLiteral,
		"decimal_integer_literal":  parseLiteral,
		"identifier":               parseIdentifier,
	}
}

func (b Binary) isExpression()           {}
func (m MethodInvocation) isExpression() {}
func (l Literal) isExpression()          {}
func (i IfNode) isExpression()           {}

func (a Access) isExpression() {}

func (a Assignment) isExpression() {}

func (i Identifier) isExpression() {}

func routeExpression(node *tree_sitter.Node, content []byte) Expression {
	if node == nil {
		return nil
	}

	if handler, exists := expressionHandlers[node.Kind()]; exists {
		expr, err := handler(node, content)
		if err == nil {
			return expr
		}
	}
	return Identifier{Name: node.Utf8Text(content)}
}

func parseParenthesized(node *tree_sitter.Node, content []byte) (Expression, error) {
	for i := range node.ChildCount() {
		child := node.Child(i)
		if child.IsNamed() {
			return routeExpression(child, content), nil
		}
	}
	return nil, fmt.Errorf("expressão entre parênteses vazia")

}
func parseMethodInvocation(node *tree_sitter.Node, content []byte) (Expression, error) {
	objNode := node.ChildByFieldName("object")
	nameNode := node.ChildByFieldName("name")
	obj := ""
	if objNode != nil {
		obj = objNode.Utf8Text(content)
	}
	methodName := ""
	if nameNode != nil {
		methodName = nameNode.Utf8Text(content)
	}

	var args []Expression
	argsNode := node.ChildByFieldName("arguments")

	for i := range argsNode.ChildCount() {
		child := argsNode.Child(i)
		expr := routeExpression(child, content)
		args = append(args, expr)

	}

	return MethodInvocation{
		Accessed: Access{
			Object:     obj,
			Identifier: Identifier{Name: methodName},
		},
		Args: args,
	}, nil

}

func parseLiteral(node *tree_sitter.Node, content []byte) (Expression, error) {
	return Literal{node.Utf8Text(content)}, nil
}
func parseIdentifier(node *tree_sitter.Node, content []byte) (Expression, error) {

	return Identifier{node.Utf8Text(content)}, nil

}

func parseBinary(node *tree_sitter.Node, content []byte) (Expression, error) {
	leftNode := node.ChildByFieldName("left")
	rightNode := node.ChildByFieldName("right")
	operatorNode := node.ChildByFieldName("operator")
	op := ""
	if operatorNode != nil {
		op = operatorNode.Utf8Text(content)
	} else if node.ChildCount() > 1 {
		op = node.Child(1).Utf8Text(content)
	}

	return Binary{
		Left:     routeExpression(leftNode, content),
		Operator: op,
		Right:    routeExpression(rightNode, content),
	}, nil
}

func parseAssignment(node *tree_sitter.Node, content []byte) (Expression, error) {
	leftNode := node.ChildByFieldName("left")
	var leftExpression Expression

	if leftNode == nil {
		return Assignment{}, fmt.Errorf("Tem que ter variavel de destino")
	}
	if leftNode.Kind() == "field_access" {
		objectNode := leftNode.ChildByFieldName("object")
		fieldNode := leftNode.ChildByFieldName("field")
		leftExpression = Access{
			Object:     objectNode.Utf8Text(content),
			Identifier: Identifier{Name: fieldNode.Utf8Text(content)},
		}
	}
	if leftNode.Kind() == "identifier" {
		leftExpression = Identifier{Name: leftNode.Utf8Text(content)}
	}

	rightNode := node.ChildByFieldName("right")
	if rightNode == nil {
		return Assignment{}, fmt.Errorf("Tem que ter valor para o destino")
	}

	return Assignment{
		Left:  leftExpression,
		Right: Identifier{rightNode.Utf8Text(content)},
	}, nil
}
