
package internal.testadata.java.controllers;

import internal.testadata.java.models.UserModel;
import internal.testadata.java.services.UsuarioService;

public class UsuarioController {

    private String nomeSistema;
    private UsuarioService service;

    public UsuarioController(String nomeSistema, UsuarioService service) {
        this.nomeSistema = nomeSistema;
        this.service = service;
    }

    public boolean processarUsuario(UserModel userModel, String status) {
        if (this.service.validarEAtivarUsuario(userModel.getIdade(), status)) {
            this.service.registrarLog("Processo concluido no sistema " + this.nomeSistema);
            return true;
        }

        return false;
    }
}