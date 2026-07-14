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
	Accessed Access `json:"accessed"`
}

type Identifier struct {
	Name string `json:"name"`
}

type Access struct {
	Object     string     `json:"object"`
	Identifier Identifier `json:"identifier"`
}

type ExpressionHandler func(node *tree_sitter.Node, content []byte) (Expression, error)

var expressionHandlers = map[string]ExpressionHandler{
	"assignment_expression": parseAssignment,
}

func (b Binary) isExpression() {}

func (a Access) isExpression() {}

func (a Assignment) isExpression() {}

func (i Identifier) isExpression() {}

func parseBinary(node *tree_sitter.Node, content []byte) (Expression, error) {
	return Binary{}, nil
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
