package java

func JavaQueries() string {
	return `
	(package_declaration) @package
	(import_declaration) @import
	(class_declaration) @class
	(field_declaration) @field
	(constructor_declaration) @constructor
	(method_declaration) @method
	`
}
