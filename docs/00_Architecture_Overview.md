
# 🌍 Global Architecture Diagram

> Visão geral de alto nível mostrando as dependências entre todas as classes analisadas.

```mermaid
flowchart LR
    %% Styling
    classDef classNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff,cursor:pointer;

    %% Nodes Creation
    
    UsuarioController["UsuarioController"]:::classNode
    click UsuarioController "UsuarioController.md" "Acessar UsuarioController"
    
    UsuarioService["UsuarioService"]:::classNode
    click UsuarioService "UsuarioService.md" "Acessar UsuarioService"
    

    %% Relationships / Dependencies
    
    
     <!-- CHAMA A NOSSA FUNÇÃO GO! -->
    
    
    
    
    
        
        
        
        
            UsuarioController -->|"Calls:<br><b>validarEAtivarUsuario()<br>registrarLog()</b>"| UsuarioService
        

    
    
    
    
     <!-- CHAMA A NOSSA FUNÇÃO GO! -->
    
    
    
    
    
    
    
```
