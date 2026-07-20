
# 🌍 Global Architecture Diagram

> Visão geral de alto nível mostrando as dependências entre todas as classes analisadas e seus respectivos pacotes.

```mermaid
flowchart LR
    %% Styling
    classDef classNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    %% Nodes Creation Grouped by Package
    
    subgraph internal.testadata.java.controllers
        
        UsuarioController["UsuarioController"]:::classNode
        
    end
    
    subgraph internal.testadata.java.models
        
        UserModel["UserModel"]:::classNode
        
    end
    
    subgraph internal.testadata.java.services
        
        UsuarioService["UsuarioService"]:::classNode
        
    end
    

    %% Relationships / Dependencies
    
    
    
    
    
    
    
    
        
        
        
        
            UsuarioController -->|"Calls:<br><b>validarEAtivarUsuario()<br>registrarLog()</b>"| UsuarioService
        

    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
```
