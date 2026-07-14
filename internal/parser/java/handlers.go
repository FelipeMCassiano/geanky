package java

import (
	"fmt"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type CaptureHandler func(node *tree_sitter.Node, content []byte, classData *ClassJava) error

type Variable struct {
	Modifiers  []Modifier `json:"modifiers"`
	TypeName   string     `json:"typeName"`
	Declarator string     `json:"declarator"`
}

type ClassJava struct {
	Modifiers    []Modifier   `json:"modifiers"`
	Name         string       `json:"name"`
	Constructors []Executable `json:"constructors"`
	Fields       []Variable   `json:"fields"`
	Methods      []Executable `json:"methods"`
}

type Executable struct {
	Modifiers     []Modifier `json:"modifiers"`
	Name          string     `json:"name"`
	Parameters    []Variable `json:"parameters"`
	Body          Body       `json:"body"`
	ReturnType    string     `json:"returnType"`
	isConstructor bool
}

type Modifier struct {
	Modifier string `json:"modifier"`
}
type Body struct {
	Statements []Statement
}
type Statement struct {
	Expressions []Expression
}

var handlers = map[string]CaptureHandler{
	"class":       parseClass,
	"method":      parseMethod,
	"field":       parseField,
	"constructor": parseConstructor,
}

func parseClass(node *tree_sitter.Node, content []byte, classData *ClassJava) error {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return fmt.Errorf("Classe nao pode ter nome nulo")
	}
	classData.Modifiers = extractModifiers(node, content)
	classData.Name = nameNode.Utf8Text(content)

	return nil
}

func parseConstructor(node *tree_sitter.Node, content []byte, classData *ClassJava) error {
	newContructor := Executable{
		Modifiers:     extractModifiers(node, content),
		isConstructor: true,
	}

	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return fmt.Errorf("Constructor nao pode ter nome nulo")
	}
	newContructor.Name = nameNode.Utf8Text(content)

	parseParameters(node, content, &newContructor)
	constructorBody := node.ChildByFieldName("body")
	parseBody(constructorBody, content, &newContructor)

	classData.Constructors = append(classData.Constructors, newContructor)
	return nil
}
func parseField(node *tree_sitter.Node, content []byte, classData *ClassJava) error {

	typeNode := node.ChildByFieldName("type")

	if typeNode == nil {
		return fmt.Errorf("Nao pode haver variavel sem tipo")
	}
	declaratorNode := node.ChildByFieldName("declarator")
	if declaratorNode == nil {
		return fmt.Errorf("Nao foi encontrado um declarador para o campo")
	}

	nameParamNode := declaratorNode.ChildByFieldName("name")
	if nameParamNode == nil {
		return fmt.Errorf("Nao pode haver tipo sem variavel")
	}

	newField := Variable{
		Declarator: nameParamNode.Utf8Text(content),
		TypeName:   typeNode.Utf8Text(content),
	}

	classData.Fields = append(classData.Fields, newField)

	return nil
}
func parseMethod(node *tree_sitter.Node, content []byte, classData *ClassJava) error {
	var newMethod = Executable{
		Modifiers:     extractModifiers(node, content),
		ReturnType:    node.ChildByFieldName("type").Utf8Text(content),
		Name:          node.ChildByFieldName("name").Utf8Text(content),
		isConstructor: false,
	}
	parseParameters(node, content, &newMethod)

	parseBody(node, content, &newMethod)

	classData.Methods = append(classData.Methods, newMethod)
	return nil
}

func parseBody(node *tree_sitter.Node, content []byte, executableData *Executable) error {
	childCount := node.ChildCount()
	var newBody Body
	for i := range childCount {
		child := node.Child(i)
		if handler, exists := statementsHandlers[child.Kind()]; exists {
			newStatement, err := handler(child, content)
			if err == nil {
				newBody.Statements = append(newBody.Statements, newStatement)
			}
		}
	}
	executableData.Body = newBody
	return nil
}
func parseParameters(node *tree_sitter.Node, content []byte, executableData *Executable) error {
	paramsGroupNode := node.ChildByFieldName("parameters")

	if paramsGroupNode == nil {
		return nil
	}

	childCount := paramsGroupNode.ChildCount()
	for i := range childCount {
		child := paramsGroupNode.Child(i)

		if child.Kind() == "formal_parameter" {
			typeNode := child.ChildByFieldName("type")
			nameParamNode := child.ChildByFieldName("name")

			newParameter := Variable{}

			if typeNode == nil {
				return fmt.Errorf("Nao pode haver variavel sem tipo")
			}
			newParameter.TypeName = typeNode.Utf8Text(content)
			if nameParamNode == nil {
				return fmt.Errorf("Nao pode haver tipo sem variavel")
			}
			newParameter.Declarator = nameParamNode.Utf8Text(content)

			executableData.Parameters = append(executableData.Parameters, newParameter)
		}
	}
	return nil
}
func extractModifiers(node *tree_sitter.Node, content []byte) []Modifier {
	var modifiers []Modifier

	for i := range node.ChildCount() {
		child := node.Child(i)
		if child.Kind() == "modifiers" {

			for j := range child.ChildCount() {
				modChild := child.Child(j)
				modifiers = append(modifiers, Modifier{modChild.Utf8Text(content)})
			}
			break
		}
	}
	return modifiers
}
