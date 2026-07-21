package com.example.parser;

import java.util.ArrayList;
import java.util.List;
import java.util.function.Supplier;

// 1. record_declaration
// public record UserDTO(int id, String name) {
// }

public class ParserTestCase {

    // null_literal, true, false (literais booleanos e nulos)
    private boolean isActive = true;
    private String status = null;

    public ParserTestCase() {
        this.isActive = false;
    }

    public UserDTO processData(List<String> names) throws Exception {

        // 2. local_variable_declaration e 3. object_creation_expression
        List<UserDTO> users = new ArrayList<>();

        // 4. enhanced_for_statement
        for (String name : names) {

            // 5. unary_expression (!) e if_statement com alternative (else)
            if (!name.isEmpty()) {

                // 6. try_statement e catch_clause
                try {
                    // 7. ternary_expression
                    String validName = (name != null) ? name : "Unknown";
                    users.add(new UserDTO(1, validName));
                } catch (RuntimeException e) {
                    // 8. throw_statement
                    throw new Exception("Erro interno", e);
                }
            } else {
                // 9. break_statement
                break;
            }
        }

        // 10. for_statement e 11. update_expression (i++)
        for (int i = 0; i < 5; i++) {
            System.out.println(i);
        }

        // 12. while_statement
        int counter = 0;
        while (counter < 3) {
            counter++;
        }

        // 13. lambda_expression
        Supplier<String> defaultName = () -> "Default";

        // 14. method_reference (String::toUpperCase)
        names.forEach(String::toUpperCase);

        // return_statement final
        return new UserDTO(0, defaultName.get());
    }
}