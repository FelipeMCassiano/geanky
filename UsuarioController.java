package com.exemplo;

public class UsuarioController {
    
    private String nomeSistema;

    public UsuarioController(String nomeSistema) {
        this.nomeSistema = nomeSistema;
    }

    // Método com parâmetro e condicional para o Tree-sitter mapear
    public boolean processarUsuario(int idade, String status) {
        if (idade >= 18 && status.equals("ATIVO")) {
            System.out.println("Usuário processado no " + this.nomeSistema);
            return true;
        }
        
        return false;
    }
}
