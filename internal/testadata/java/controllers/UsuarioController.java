
package internal.testadata.java.controllers;

import internal.testadata.java.services.UsuarioService;

public class UsuarioController {

    private String nomeSistema;

    public UsuarioController(String nomeSistema, UsuarioService service) {
        this.nomeSistema = nomeSistema;
        this.service = service;
    }

    public boolean processarUsuario(int idade, String status) {

        
        return false;
    }
}