package parser

func ClassAndMethodNames() string {

	queryPattern := `
	(class_declaration
		name: (identifier) @class_name)
	`
	return queryPattern
}
