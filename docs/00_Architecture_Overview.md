
# 🌍 Global Architecture Diagram

> Visão geral de alto nível mostrando as dependências entre todas as classes analisadas e seus respectivos pacotes.

```mermaid
flowchart LR
    %% Styling
    classDef classNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    %% Nodes Creation Grouped by Package
    
    subgraph controllers
        
        UsuarioController["UsuarioController"]:::classNode
        
    end
    
    subgraph models
        
        UserModel["UserModel"]:::classNode
        
    end
    
    subgraph services
        
        UsuarioService["UsuarioService"]:::classNode
        
    end
    
    subgraph sync
        
        BoatSyncController["BoatSyncController"]:::classNode
        
    end
    

    %% Relationships / Dependencies
    
    
    
    
    
    
    
    
        
        
        
        
            UsuarioController -->|"Calls:<br><b>validarEAtivarUsuario(Invoke &#39;userModel.getIdade&#39; (no parameters), status)&lt;br&gt;registrarLog(&#39;Processo concluido no sistema &#39; plus this.nomeSistema)</b>"| UsuarioService
        

    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
        
        
        
        
            BoatSyncController -->|"Depends on"| Logger
        

    
    
    
        
        
        
        
            BoatSyncController -->|"Depends on"| SyncService
        

    
    
    
        
        
        
        
            BoatSyncController -->|"Depends on"| CloudSyncUploadService
        

    
    
    
        
        
        
        
            BoatSyncController -->|"Calls:<br><b>subscribe(uuid)</b>"| BoatSseSyncService
        

    
    
    
        
        
        
        
            BoatSyncController -->|"Depends on"| DownloadService
        

    
    
    
        
        
        
        
            BoatSyncController -->|"Depends on"| MachineSync
        

    
    
    
```
