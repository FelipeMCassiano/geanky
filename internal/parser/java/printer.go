package java

import "fmt"

func PrintClassSummary(classData ClassJava) {

	classMods := formatModifiers(classData.Modifiers)

	fmt.Println("==================================================")
	fmt.Printf("📦 Classe: %s%s\n", classMods, classData.Name)
	fmt.Println("==================================================")

	// --- PROPRIEDADES ---
	fmt.Printf("\n🔹 Propriedades (%d):\n", len(classData.Fields))
	if len(classData.Fields) == 0 {
		fmt.Println("   (Nenhuma propriedade encontrada)")
	} else {
		for _, f := range classData.Fields {
			fmt.Printf("   - %s%s %s;\n", formatModifiers(f.Modifiers), f.TypeName, f.Declarator)
		}
	}

	// --- CONSTRUTORES ---
	fmt.Printf("\n🛠️ Construtores (%d):\n", len(classData.Constructors))
	if len(classData.Constructors) == 0 {
		fmt.Println("   (Nenhum construtor encontrado)")
	} else {
		for _, c := range classData.Constructors {
			fmt.Printf("   - %s%s(", formatModifiers(c.Modifiers), c.Name)
			for i, p := range c.Parameters {
				fmt.Printf("%s %s", p.TypeName, p.Declarator)
				if i < len(c.Parameters)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Println(")")

			fmt.Println("     [Body]:")
			if len(c.Body.Statements) == 0 {
				fmt.Println("       (Corpo vazio ou não processado)")
			} else {
				for stmtIdx, stmt := range c.Body.Statements {
					fmt.Printf("       %d. ", stmtIdx+1)
					for _, expr := range stmt.Expressions {
						printExpressionPretty(expr)
					}
					fmt.Println()
				}
			}
		}
	}

	// --- MÉTODOS ---
	fmt.Printf("\n⚙️ Métodos (%d):\n", len(classData.Methods))
	if len(classData.Methods) == 0 {
		fmt.Println("   (Nenhum método encontrado)")
	} else {
		for _, m := range classData.Methods {

			// Diferença aqui: Incluímos o m.ReturnType antes do nome do método!
			fmt.Printf("   - %s%s %s(", formatModifiers(m.Modifiers), m.ReturnType, m.Name)

			// 1. Imprime os Parâmetros do Método
			for i, p := range m.Parameters {
				fmt.Printf("%s %s", p.TypeName, p.Declarator)
				if i < len(m.Parameters)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Println(")")

			// 2. Imprime o Body do Método
			fmt.Println("     [Body]:")
			if len(m.Body.Statements) == 0 {
				fmt.Println("       (Corpo vazio ou não processado)")
			} else {
				for stmtIdx, stmt := range m.Body.Statements {
					fmt.Printf("       %d. ", stmtIdx+1)
					for _, expr := range stmt.Expressions {
						printExpressionPretty(expr)
					}
					fmt.Println() // Quebra de linha da instrução
				}
			}
		}
	}
	fmt.Println("\n==================================================")
}
func printExpressionPretty(expr Expression) {
	switch e := expr.(type) {
	case Assignment:
		leftStr := formatExpression(e.Left)
		rightStr := formatExpression(e.Right)
		fmt.Printf("%s = %s;", leftStr, rightStr)
	case Identifier:
		fmt.Printf("%s;", e.Name)
	default:
		fmt.Printf("%#v;", expr)
	}
}

func formatExpression(expr Expression) string {
	switch e := expr.(type) {
	case Identifier:
		return e.Name
	case Access:
		if e.Object != "" {
			return fmt.Sprintf("%s.%s", e.Object, e.Identifier.Name)
		}
		return e.Identifier.Name
	default:
		return fmt.Sprintf("%v", expr)
	}
}
func formatModifiers(mods []Modifier) string {
	if len(mods) == 0 {
		return ""
	}

	var result string
	for i, m := range mods {
		result += m.Modifier
		if i < len(mods)-1 {
			result += " "
		}
	}
	return result + " "
}
