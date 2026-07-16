
package internal.testadata.java.controllers;

import internal.testadata.java.services.UsuarioService;

public class UsuarioController {

    private String nomeSistema;
    private UsuarioService service; 

    public UsuarioController(String nomeSistema, UsuarioService service) {
        this.nomeSistema = nomeSistema;
        this.service = service;
    }

    public boolean processarUsuario(int idade, String status) {
        if (this.service.validarEAtivarUsuario(idade, status)) {
            this.service.registrarLog("Processo concluido no sistema " + this.nomeSistema);
            return true;
        }

        return false;
    }
}