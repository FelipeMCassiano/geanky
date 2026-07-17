package java

import (
	"fmt"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Package struct {
	Scope string `json:"scope"`
	Name  string `json:"name"`
}

type CaptureHandler func(node *tree_sitter.Node, content []byte, classData *ClassJava) error

type Variable struct {
	Modifiers  []Modifier `json:"modifiers"`
	TypeName   string     `json:"typeName"`
	Declarator string     `json:"declarator"`
}

type ClassJava struct {
	Package      Package      `json:"package"`
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
	Body          Block      `json:"body"`
	ReturnType    string     `json:"returnType"`
	isConstructor bool
}

type Modifier struct {
	Modifier string `json:"modifier"`
}

type Block struct {
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
	"package":     parsePackage,
}

func parsePackage(node *tree_sitter.Node, content []byte, classData *ClassJava) error {
	scopeNode := node.ChildByFieldName("scope")
	nameNode := node.ChildByFieldName("name")

	if nameNode == nil {
		return fmt.Errorf("Name nao pode ser nulo")
	}
	if scopeNode == nil {
		return fmt.Errorf("Scope nao pode ser nulo")
	}
	classData.Package = Package{
		Name:  nameNode.Utf8Text(content),
		Scope: scopeNode.Utf8Text(content),
	}

	return nil
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
	newContructor.Body = parseBlock(constructorBody, content)

	classData.Constructors = append(classData.Constructors, newContructor)
	return nil
}
func parseField(node *tree_sitter.Node, content []byte, classData *ClassJava) error {
	typeNode := node.ChildByFieldName("type")
	declaratorNode := node.ChildByFieldName("declarator")

	if typeNode == nil || declaratorNode == nil {
		return nil
	}

	nameParamNode := declaratorNode.ChildByFieldName("name")
	if nameParamNode == nil {
		return nil
	}

	newField := Variable{
		Modifiers:  extractModifiers(node, content), // AGORA extraímos os modificadores!
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

	bodyNode := node.ChildByFieldName("body")

	newMethod.Body = parseBlock(bodyNode, content)

	classData.Methods = append(classData.Methods, newMethod)
	return nil
}

func parseBlock(node *tree_sitter.Node, content []byte) Block {
	childCount := node.ChildCount()
	var newBlock Block

	if node == nil {
		return newBlock
	}

	for i := range childCount {
		child := node.Child(i)
		if handler, exists := statementsHandlers[child.Kind()]; exists {
			newStatement, err := handler(child, content)
			if err == nil {
				newBlock.Statements = append(newBlock.Statements, newStatement)
			}
		}
	}
	return newBlock
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
