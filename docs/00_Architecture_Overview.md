
# 🌍 Global Architecture Diagram

> Visão geral de alto nível mostrando as dependências entre todas as classes analisadas e seus respectivos pacotes.

```mermaid
flowchart LR
    %% Styling
    classDef classNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    %% Nodes Creation Grouped by Package
    
    subgraph Default Package
        
        UsuarioController["UsuarioController"]:::classNode
        
        UsuarioService["UsuarioService"]:::classNode
        
    end
    

    %% Relationships / Dependencies
    
    
    
    
    
    
    
    
        
        
        
        
            UsuarioController -->|"Calls:<br><b>validarEAtivarUsuario()<br>registrarLog()</b>"| UsuarioService
        

    
    
    
    
    
    
    
    
    
    
    
    
```
