package java

func JavaQueries() string {
	return `
	(package_declaration) @package
	(class_declaration) @class
	(field_declaration) @field
	(constructor_declaration) @constructor
	(method_declaration) @method
	`
}
