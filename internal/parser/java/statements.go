package java

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type StatementHandler func(node *tree_sitter.Node, content []byte) (Statement, error)

var statementsHandlers map[string]StatementHandler

func init() {
	statementsHandlers = map[string]StatementHandler{
		"expression_statement": parseExpressionStatement,
		"if_statement":         parseIfStatement,
		"return_statement":     parseReturnStatement,
	}
}
func parseReturnStatement(node *tree_sitter.Node, content []byte) (Statement, error) {
	var returnValue Expression
	for i := range node.ChildCount() {
		child := node.Child(i)
		if child.IsNamed() {
			returnValue = routeExpression(child, content)
			break
		}
	}

	return Statement{Expressions: []Expression{ReturnNode{returnValue}}}, nil
}

func parseIfStatement(node *tree_sitter.Node, content []byte) (Statement, error) {
	conditionNode := node.ChildByFieldName("condition")

	conditionExpr := routeExpression(conditionNode, content)
	consequenceNode := node.ChildByFieldName("consequence")
	consequence := parseBlock(consequenceNode, content)

	return Statement{Expressions: []Expression{IfNode{Condition: conditionExpr, Consequence: consequence}}}, nil
}

func parseExpressionStatement(node *tree_sitter.Node, content []byte) (Statement, error) {
	var expressions []Expression
	for i := range node.ChildCount() {
		child := node.Child(i)

		if handler, exists := expressionHandlers[child.Kind()]; exists {
			newExpression, err := handler(child, content)
			if err == nil {
				expressions = append(expressions, newExpression)
			}
		}
	}
	return Statement{expressions}, nil
}
