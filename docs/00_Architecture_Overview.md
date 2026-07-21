
# 🌍 Global Architecture Diagram

> Visão geral de alto nível mostrando as dependências entre todas as classes analisadas e seus respectivos pacotes.

```mermaid
flowchart LR
    classDef classNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    
    subgraph com.exampleparser
        
        ParserTestCase["ParserTestCase"]:::classNode
        
    end
    
    subgraph com.rfidbrasil.core.controllersync
        
        BoatSyncController["BoatSyncController"]:::classNode
        
    end
    
    subgraph com.rfidbrasil.core.servicesync
        
        CloudSyncUploadService["CloudSyncUploadService"]:::classNode
        
    end
    
    subgraph internal.testadata.javacontrollers
        
        UsuarioController["UsuarioController"]:::classNode
        
    end
    
    subgraph internal.testadata.javamodels
        
        UserModel["UserModel"]:::classNode
        
    end
    
    subgraph internal.testadata.javaservices
        
        UsuarioService["UsuarioService"]:::classNode
        
    end
    

    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
        
        
            UsuarioController -->|"Calls:<br><b>validarEAtivarUsuario(..., status)<br>registrarLog(...)</b>"| UsuarioService
        
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
        
        
            BoatSyncController -->|"Calls:<br><b>error(..., ...)</b>"| Logger
        
    
    
    
        
        
            BoatSyncController -->|"Calls:<br><b>syncPortalReadingsPaginated(timestamp, id, principal)</b>"| SyncService
        
    
    
    
        
        
            BoatSyncController -->|"Calls:<br><b>latestSync()</b>"| CloudSyncUploadService
        
    
    
    
        
        
            BoatSyncController -->|"Calls:<br><b>subscribe(uuid)</b>"| BoatSseSyncService
        
    
    
    
        
        
            BoatSyncController -->|"Depends on"| DownloadService
        
    
    
    
        
        
            BoatSyncController -->|"Calls:<br><b>executeSync(true, ...)</b>"| MachineSync
        
    
    
    
    
    
    
    
    
        
        
            CloudSyncUploadService -->|"Depends on"| MachineClient
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>findStatusTransitions()<br>markAsSynced(situationIds)</b>"| InventorySituationRepository
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>markAsSyncedByNames(inventoryNames)</b>"| InventoryRepository
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>buildPendingMusteringsPayload()</b>"| SyncService
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>sendChunk(chunk, master, currentToken, ...)</b>"| CloudSyncIntegrationClient
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>findFirstByOrderByIdDesc()<br>saveAndFlush(newControl)<br>findFirstByEventKeyOrderByIdDesc(eventKey)</b>"| SyncControlRepository
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>sendSyncStatusOnChange(eventKey, uiMessage, ...)<br>sendSyncStatusOnChange(..., ..., ...)<br>sendSyncStatusOnChange(eventKey, ..., ...)</b>"| BoatSseSyncService
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>download(eventKey, request, token, locationId, ...)</b>"| DownloadService
        
    
    
    
        
        
            CloudSyncUploadService -->|"Depends on"| Executor
        
    
    
    
        
        
            CloudSyncUploadService -->|"Depends on"| ItemRepository
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>findFirstPending()<br>save(master)<br>findOldMastersToCleanup()<br>deleteAll(oldMasters)</b>"| SyncPacketMasterRepository
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>findUnsyncedByMasterId(...)<br>save(chunk)</b>"| SyncPacketChunkRepository
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>sliceIntoChunks(gzipPayload)<br>calculateSha256(gzipPayload)<br>compressToGzip(pbPackage)</b>"| PacketChunkerEngine
        
    
    
    
        
        
            CloudSyncUploadService -->|"Calls:<br><b>exchange(url, ..., requestEntity, String.class)</b>"| RestTemplate
        
    
    
    
    
    
    
    
    
    
```
