package java

import tree_sitter "github.com/tree-sitter/go-tree-sitter"

type StatementHandler func(node *tree_sitter.Node, content []byte) (Statement, error)

var statementsHandlers = map[string]StatementHandler{
	"expression_statement": parseExpressionStatement,
	"if_statement":         parseIfStatement,
}

func parseIfStatement(node *tree_sitter.Node, content []byte) (Statement, error) {
	// conditionNode := node.ChildByFieldName("condition")
	return Statement{}, nil
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
