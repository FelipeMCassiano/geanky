package java

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type StatementHandler func(node *tree_sitter.Node, content []byte) (Statement, error)

var statementsHandlers map[string]StatementHandler

func init() {
	statementsHandlers = map[string]StatementHandler{
		"expression_statement":       parseExpressionStatement,
		"if_statement":               parseIfStatement,
		"return_statement":           parseReturnStatement,
		"local_variable_declaration": parseLocalVariableDeclaration,
		"try_statement":              parseTryStatement,
		"for_statement":              parseForStatement,
		"enhanced_for_statement":     parseEnhancedForStatement,
		"while_statement":            parseWhileStatement,
		"throw_statement":            parseThrowStatement,
		"break_statement":            parseBreakStatement,
	}
}

func parseForStatement(node *tree_sitter.Node, content []byte) (Statement, error) {
	initNode := node.ChildByFieldName("init")
	condNode := node.ChildByFieldName("condition")
	updateNode := node.ChildByFieldName("update")
	bodyNode := node.ChildByFieldName("body")
	var initExpr, condExpr, updateExpr Expression

	if initNode != nil {
		if initNode.Kind() == "local_variable_declaration" {
			stmt, _ := parseLocalVariableDeclaration(initNode, content)
			if len(stmt.Expressions) > 0 {
				initExpr = stmt.Expressions[0]
			}
		} else {
			initExpr = routeExpression(initNode, content)
		}
	}
	if condNode != nil {
		condExpr = routeExpression(condNode, content)
	}
	if updateNode != nil {
		updateExpr = routeExpression(updateNode, content)
	}

	body := parseBlock(bodyNode, content)

	return Statement{Expressions: []Expression{ForNode{
		Init:      initExpr,
		Condition: condExpr,
		Update:    updateExpr,
		Body:      body,
	}}}, nil
}
func parseEnhancedForStatement(node *tree_sitter.Node, content []byte) (Statement, error) {
	typeNode := node.ChildByFieldName("type")
	nameNode := node.ChildByFieldName("name")
	valueNode := node.ChildByFieldName("value")
	bodyNode := node.ChildByFieldName("body")
	typeStr, nameStr := "", ""
	if typeNode != nil {
		typeStr = typeNode.Utf8Text(content)
	}
	if nameNode != nil {
		nameStr = nameNode.Utf8Text(content)
	}

	valueExpr := routeExpression(valueNode, content)
	body := parseBlock(bodyNode, content)

	return Statement{[]Expression{
		EnhancedForNode{
			Type:  typeStr,
			Name:  nameStr,
			Value: valueExpr,
			Body:  body,
		},
	}}, nil
}
func parseThrowStatement(node *tree_sitter.Node, content []byte) (Statement, error) {
	for i := range node.ChildCount() {
		child := node.Child(i)
		if child.IsNamed() {
			return Statement{Expressions: []Expression{ThrowNode{
				Value: routeExpression(child, content),
			}}}, nil
		}
	}
	return Statement{}, nil
}

func parseBreakStatement(node *tree_sitter.Node, content []byte) (Statement, error) {
	return Statement{Expressions: []Expression{BreakNode{}}}, nil
}
func parseWhileStatement(node *tree_sitter.Node, content []byte) (Statement, error) {
	condNode := node.ChildByFieldName("condition")
	bodyNode := node.ChildByFieldName("body")

	condExpr := routeExpression(condNode, content)
	body := parseBlock(bodyNode, content)

	return Statement{Expressions: []Expression{WhileNode{
		Condition: condExpr,
		Body:      body,
	}}}, nil
}

func parseTryStatement(node *tree_sitter.Node, content []byte) (Statement, error) {
	bodyNode := node.ChildByFieldName("body")
	tryBody := parseBlock(bodyNode, content)

	catches := []CatchNode{}

	for i := range node.ChildCount() {
		child := node.Child(i)

		if child.Kind() == "catch_clause" {
			catchBodyNode := child.ChildByFieldName("body")
			paramNode := child.ChildByFieldName("parameter")
			paramStr := ""
			if paramNode != nil {
				paramStr = paramNode.Utf8Text(content)
			}
			catches = append(catches, CatchNode{
				Parameter: paramStr,
				Body:      parseBlock(catchBodyNode, content),
			})
		}
	}
	return Statement{Expressions: []Expression{
		TryNode{
			Body:    tryBody,
			Catches: catches,
		},
	}}, nil
}
func parseLocalVariableDeclaration(node *tree_sitter.Node, content []byte) (Statement, error) {

	typeNode := node.ChildByFieldName("type")
	declaratorNode := node.ChildByFieldName("declarator")
	nameNode := declaratorNode.ChildByFieldName("name")
	valueNode := declaratorNode.ChildByFieldName("value")
	expr := routeExpression(valueNode, content)

	return Statement{Expressions: []Expression{
		Variable{
			TypeName:   typeNode.Utf8Text(content),
			Declarator: nameNode.Utf8Text(content),
			Value:      expr,
		},
	}}, nil
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
	var alternativeBlock *Block
	alternativeNode := node.ChildByFieldName("alternative")
	if alternativeNode != nil {
		alt := parseBlock(alternativeNode, content)
		alternativeBlock = &alt
	}

	return Statement{Expressions: []Expression{IfNode{Condition: conditionExpr, Consequence: consequence, Alternative: alternativeBlock}}}, nil
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
