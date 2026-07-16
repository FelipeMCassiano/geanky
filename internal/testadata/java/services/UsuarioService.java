
package internal.testadata.java.services;

public class UsuarioService {

    private String nomeBanco;
    private boolean conexaoAtiva;

    public UsuarioService(String nomeBanco) {
        this.nomeBanco = nomeBanco;
        this.conexaoAtiva = true;
    }

    public boolean validarEAtivarUsuario(int idade, String status) {
        if (idade >= 18 && status == "ativo") {
            System.out.println("Usuario validado com sucesso no banco " + this.nomeBanco);
            return true;
        }

        System.out.println("Falha na validacao");
        return false;
    }

    public void registrarLog(String acao) {
        this.conexaoAtiva = false;
        System.out.println(acao);
    }
}