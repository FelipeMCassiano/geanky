package com.exemplo;

public class UsuarioController {

    private String nomeSistema;
    private int ano;

    public UsuarioController(String nomeSistema, int ano) {
        this.nomeSistema = nomeSistema;
        this.ano = ano;
    }

    // Método com parâmetro e condicional para o Tree-sitter mapear
    public boolean processarUsuario(int idade, String status) {
        if (idade >= 18 && status.equals("ATIVO")) {
            System.out.println("Usuário processado no " + this.nomeSistema + "ano: " + ano);
            return true;
        }

        return false;
    }
}
