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
	case MethodInvocation:
		fmt.Printf("%s;", formatExpression(e))
	case IfNode:
		fmt.Printf("%s", formatExpression(e))
	default:
		fmt.Printf("%#v;", expr)
	}
}
func formatExpression(expr Expression) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case Assignment:
		return fmt.Sprintf("%s = %s", formatExpression(e.Left), formatExpression(e.Right))
	case Binary:
		return fmt.Sprintf("%s %s %s", formatExpression(e.Left), e.Operator, formatExpression(e.Right))
	case IfNode:
		condStr := formatExpression(e.Condition)

		// 1. Se o bloco do IF estiver vazio
		if len(e.Consequence.Statements) == 0 {
			return fmt.Sprintf("if (%s) { }", condStr)
		}

		// 2. Se tiver código, abre a chave e quebra a linha
		result := fmt.Sprintf("if (%s) {\n", condStr)

		// 3. Itera sobre o novo Block
		for _, stmt := range e.Consequence.Statements {
			// Adiciona espaços para identar o que está dentro do IF
			result += "             "
			for _, expr := range stmt.Expressions {
				// Mágica da recursão: vai formatar variáveis, invocações, atribuições...
				result += formatExpression(expr) + ";"
			}
			result += "\n" // Quebra de linha após cada instrução
		}

		// 4. Fecha a chave do IF alinhando com a margem
		result += "         }"
		return result
	case MethodInvocation:
		// 1. Extraímos e formatamos todos os argumentos
		var argsStr string
		for i, arg := range e.Args {
			argsStr += formatExpression(arg) // Recursão! Funciona para literais, identificadores, etc.
			if i < len(e.Args)-1 {
				argsStr += ", " // Adiciona vírgula entre os argumentos
			}
		}

		// CORREÇÃO AQUI: != nil em vez de != "", e usar formatExpression no Object
		if e.Accessed.Object != nil {
			return fmt.Sprintf("%s.%s(%s)", formatExpression(e.Accessed.Object), e.Accessed.Identifier.Name, argsStr)
		}
		return fmt.Sprintf("%s(%s)", e.Accessed.Identifier.Name, argsStr)

	case ReturnNode:
		// Para imprimir os "return true;" que capturamos!
		if e.Value != nil {
			return fmt.Sprintf("return %s", formatExpression(e.Value))
		}
		return "return"

	case Literal:
		return e.Value
	case Access:
		// CORREÇÃO AQUI: != nil em vez de != "", e usar formatExpression no Object
		if e.Object != nil {
			return fmt.Sprintf("%s.%s", formatExpression(e.Object), e.Identifier.Name)
		}
		return e.Identifier.Name
	case Identifier:
		return e.Name
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
