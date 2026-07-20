
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
        
        CloudSyncUploadService["CloudSyncUploadService"]:::classNode
        
    end
    

    %% Relationships / Dependencies
    
    
    
    
    
    
    
    
        
        
        
        
            UsuarioController -->|"Calls:<br><b>validarEAtivarUsuario(..., status)&lt;br&gt;registrarLog(...)</b>"| UsuarioService
        

    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
        
        
        
        
            BoatSyncController -->|"Depends on"| Logger
        

    
    
    
        
        
        
        
            BoatSyncController -->|"Depends on"| SyncService
        

    
    
    
        
        
        
        
            BoatSyncController -->|"Depends on"| CloudSyncUploadService
        

    
    
    
        
        
        
        
            BoatSyncController -->|"Calls:<br><b>subscribe(uuid)</b>"| BoatSseSyncService
        

    
    
    
        
        
        
        
            BoatSyncController -->|"Depends on"| DownloadService
        

    
    
    
        
        
        
        
            BoatSyncController -->|"Depends on"| MachineSync
        

    
    
    
    
    
    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Depends on"| MachineClient
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Calls:<br><b>markAsSynced(situationIds)</b>"| InventorySituationRepository
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Calls:<br><b>markAsSyncedByNames(inventoryNames)</b>"| InventoryRepository
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Depends on"| SyncService
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Depends on"| CloudSyncIntegrationClient
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Calls:<br><b>saveAndFlush(newControl)</b>"| SyncControlRepository
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Calls:<br><b>sendSyncStatusOnChange(..., ..., ...)&lt;br&gt;sendSyncStatusOnChange(eventKey, ..., ...)</b>"| BoatSseSyncService
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Depends on"| DownloadService
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Depends on"| Executor
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Depends on"| ItemRepository
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Calls:<br><b>save(master)</b>"| SyncPacketMasterRepository
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Depends on"| SyncPacketChunkRepository
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Depends on"| PacketChunkerEngine
        

    
    
    
        
        
        
        
            CloudSyncUploadService -->|"Depends on"| RestTemplate
        

    
    
    
    
    
    
    
    
    
```
