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
	Condition   Expression `json:"condition"`
	Consequence Block      `json:"consequence"`
	Alternative *Block     `json:"alternative"`
}

type MethodInvocation struct {
	Accessed Access       `json:"accessed"`
	Args     []Expression `json:"args"`
}

type Identifier struct {
	Name string `json:"name"`
}

type Access struct {
	Object     Expression `json:"object"`
	Identifier Identifier `json:"identifier"`
}
type Literal struct {
	Value string `json:"value"`
}

type ReturnNode struct {
	Value Expression `json:"value"`
}
type TryNode struct {
	Body    Block       `json:"body"`
	Catches []CatchNode `json:"catches"`
}

type CatchNode struct {
	Parameter string `json:"parameter"`
	Body      Block  `json:"body"`
}
type ForNode struct {
	Init      Expression `json:"init"`
	Condition Expression `json:"condition"`
	Update    Expression `json:"update"`
	Body      Block      `json:"body"`
}

type EnhancedForNode struct {
	Modifiers []Modifier `json:"modifiers"`
	Type      string     `json:"type"`
	Name      string     `json:"name"`
	Value     Expression `json:"value"`
	Body      Block      `json:"body"`
}

type WhileNode struct {
	Condition Expression `json:"condition"`
	Body      Block      `json:"body"`
}

type ThrowNode struct {
	Value Expression `json:"value"`
}

type BreakNode struct{}

type ExpressionHandler func(node *tree_sitter.Node, content []byte) (Expression, error)

var expressionHandlers map[string]ExpressionHandler

func init() {
	expressionHandlers = map[string]ExpressionHandler{
		"assignment_expression":    parseAssignment,
		"binary_expression":        parseBinary,
		"parenthesized_expression": parseParenthesized,
		"method_invocation":        parseMethodInvocation,
		"string_literal":           parseLiteral,
		"decimal_integer_literal":  parseLiteral,
		"identifier":               parseIdentifier,
		"return_statement":         parseReturnNode,
		"field_access":             parseAccess,
	}
}

func (t TryNode) isExpression()          {}
func (b Binary) isExpression()           {}
func (m MethodInvocation) isExpression() {}
func (l Literal) isExpression()          {}
func (i IfNode) isExpression()           {}
func (r ReturnNode) isExpression()       {}
func (f ForNode) isExpression()          {}
func (e EnhancedForNode) isExpression()  {}
func (w WhileNode) isExpression()        {}
func (t ThrowNode) isExpression()        {}
func (b BreakNode) isExpression()        {}
func (a Access) isExpression()           {}
func (a Assignment) isExpression()       {}
func (i Identifier) isExpression()       {}

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

func parseAccess(node *tree_sitter.Node, content []byte) (Expression, error) {
	objectNode := node.ChildByFieldName("object")
	fieldNode := node.ChildByFieldName("field")

	fieldName := ""
	if fieldNode != nil {
		fieldName = fieldNode.Utf8Text(content)
	}

	return Access{
		Object:     routeExpression(objectNode, content),
		Identifier: Identifier{Name: fieldName},
	}, nil
}

func parseReturnNode(node *tree_sitter.Node, content []byte) (Expression, error) {
	for i := range node.ChildCount() {
		child := node.Child(i)
		if child.IsNamed() {
			return ReturnNode{Value: routeExpression(child, content)}, nil
		}
	}
	return ReturnNode{Value: nil}, nil
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
	methodName := ""
	if nameNode != nil {
		methodName = nameNode.Utf8Text(content)
	}

	var args []Expression
	argsNode := node.ChildByFieldName("arguments")

	for i := range argsNode.ChildCount() {
		child := argsNode.Child(i)
		if child.IsNamed() {
			expr := routeExpression(child, content)
			args = append(args, expr)
		}

	}

	return MethodInvocation{
		Accessed: Access{
			Object:     routeExpression(objNode, content),
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

	if leftNode == nil {
		return Assignment{}, fmt.Errorf("Tem que ter variavel de destino")
	}

	rightNode := node.ChildByFieldName("right")
	if rightNode == nil {
		return Assignment{}, fmt.Errorf("Tem que ter valor para o destino")
	}

	return Assignment{
		Left:  routeExpression(leftNode, content),
		Right: routeExpression(rightNode, content),
	}, nil
}
