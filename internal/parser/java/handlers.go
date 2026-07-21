package java

import (
	"fmt"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Package struct {
	Scope string `json:"scope"`
	Name  string `json:"name"`
}

type CaptureHandler func(node *tree_sitter.Node, content []byte, classData *ClassJava) error

type Variable struct {
	Annotations []string   `json:"annotations"`
	Modifiers   []Modifier `json:"modifiers"`
	TypeName    string     `json:"typeName"`
	Declarator  string     `json:"declarator"`
	Value       Expression `json:"value"`
}

func (v Variable) isExpression() {}

type ClassJava struct {
	Annotations  []string     `json:"annotations"`
	Package      Package      `json:"package"`
	Imports      []string     `json:"imports"`
	Modifiers    []Modifier   `json:"modifiers"`
	Name         string       `json:"name"`
	Constructors []Executable `json:"constructors"`
	Fields       []Variable   `json:"fields"`
	Methods      []Executable `json:"methods"`
}

type Executable struct {
	Annotations   []string   `json:"annotations"`
	Modifiers     []Modifier `json:"modifiers"`
	Name          string     `json:"name"`
	Parameters    []Variable `json:"parameters"`
	Body          Block      `json:"body"`
	ReturnType    string     `json:"returnType"`
	IsConstructor bool       `json:"isConstructor"`
	Throws        string     `json:"throws"`
}

type Modifier struct {
	Modifier string `json:"modifier"`
}

type Block struct {
	Statements []Statement `json:"statements"`
}

type Statement struct {
	Expressions []Expression `json:"expressions"`
}

var handlers = map[string]CaptureHandler{
	"class":       parseClass,
	"method":      parseMethod,
	"field":       parseField,
	"constructor": parseConstructor,
	"package":     parsePackage,
	"import":      parseImport,
}

func parseImport(node *tree_sitter.Node, content []byte, classData *ClassJava) error {
	importNode := node.NamedChild(0)
	classData.Imports = append(classData.Imports, importNode.Utf8Text(content))
	return nil
}

func parsePackage(node *tree_sitter.Node, content []byte, classData *ClassJava) error {
	pkgNode := node.NamedChild(0)
	scopeNode := pkgNode.ChildByFieldName("scope")
	nameNode := pkgNode.ChildByFieldName("name")

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
	modifiers, annotations := extractModifiersAndAnnotations(node, content)
	classData.Annotations = annotations
	classData.Modifiers = modifiers
	classData.Name = nameNode.Utf8Text(content)

	return nil
}

func parseConstructor(node *tree_sitter.Node, content []byte, classData *ClassJava) error {

	modifiers, annotations := extractModifiersAndAnnotations(node, content)
	newContructor := Executable{
		Annotations:   annotations,
		Modifiers:     modifiers,
		IsConstructor: true,
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
	modifiers, annotations := extractModifiersAndAnnotations(node, content)

	newField := Variable{
		Annotations: annotations,
		Modifiers:   modifiers,
		Declarator:  nameParamNode.Utf8Text(content),
		TypeName:    typeNode.Utf8Text(content),
	}

	classData.Fields = append(classData.Fields, newField)
	return nil
}
func parseMethod(node *tree_sitter.Node, content []byte, classData *ClassJava) error {

	modifiers, annotations := extractModifiersAndAnnotations(node, content)
	var newMethod = Executable{
		Annotations:   annotations,
		Modifiers:     modifiers,
		ReturnType:    node.ChildByFieldName("type").Utf8Text(content),
		Name:          node.ChildByFieldName("name").Utf8Text(content),
		IsConstructor: false,
	}
	parseParameters(node, content, &newMethod)

	bodyNode := node.ChildByFieldName("body")

	newMethod.Body = parseBlock(bodyNode, content)
	for i := range node.ChildCount() {
		child := node.Child(i)
		if child.Kind() == "throws" {
			rawText := child.Utf8Text(content)

			cleanText := strings.Replace(rawText, "throws", "", 1)
			newMethod.Throws = strings.TrimSpace(cleanText)
			break
		}
	}
	fmt.Println(newMethod.Throws)
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
func extractModifiersAndAnnotations(node *tree_sitter.Node, content []byte) ([]Modifier, []string) {
	var modifiers []Modifier
	var annotations []string

	for i := range node.ChildCount() {
		child := node.Child(i)
		if child.Kind() == "modifiers" {
			for j := range child.ChildCount() {
				modChild := child.Child(j)
				kind := modChild.Kind()

				// Identifica se o modificador é uma anotação
				if kind == "annotation" || kind == "marker_annotation" {
					// Pega a string completa: ex: "@RequestMapping(value = \"/user\")"
					annotations = append(annotations, modChild.Utf8Text(content))
				} else {
					modifiers = append(modifiers, Modifier{modChild.Utf8Text(content)})
				}
			}
			break
		}
	}
	return modifiers, annotations
}
