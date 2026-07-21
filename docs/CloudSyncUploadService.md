
> **@Slf4j**
> **@Service**
# 📄 Technical Specification: `CloudSyncUploadService`

> **Package:** sync
> **Dependencies (Imports):**
> - java.net.http.HttpHeaders
> - java.time.Instant
> - java.util.ArrayList
> - java.util.List
> - java.util.Optional
> - java.util.UUID
> - java.util.concurrent.CompletableFuture
> - java.util.concurrent.CompletionException
> - java.util.concurrent.Executor
> - java.util.stream.Collectors
> - org.springframework.beans.factory.annotation.Qualifier
> - org.springframework.beans.factory.annotation.Value
> - org.springframework.http.HttpEntity
> - org.springframework.http.HttpMethod
> - org.springframework.http.HttpStatus
> - org.springframework.http.MediaType
> - org.springframework.http.ResponseEntity
> - org.springframework.http.client.SimpleClientHttpRequestFactory
> - org.springframework.stereotype.Service
> - org.springframework.web.client.HttpClientErrorException
> - org.springframework.web.client.RestTemplate
> - com.fasterxml.jackson.databind.JsonNode
> - com.fasterxml.jackson.databind.ObjectMapper
> - com.rfidbrasil.core.dto.ItemModificationStatusDTO
> - com.rfidbrasil.core.dto.ItemStatusSyncDTO
> - com.rfidbrasil.core.dto.MusteringSituationSyncDTO
> - com.rfidbrasil.core.dto.MusteringSyncDTO
> - com.rfidbrasil.core.dto.projection.ItemStatusTransitionProjection
> - com.rfidbrasil.core.dto.request.SyncRequest
> - com.rfidbrasil.core.dto.request.UpdateSyncStatusRequest
> - com.rfidbrasil.core.dto.response.PingSyncResponse
> - com.rfidbrasil.core.dto.response.SyncEvent
> - com.rfidbrasil.core.dto.response.SyncLastestResponse
> - com.rfidbrasil.core.enums.EPacketStatus
> - com.rfidbrasil.core.enums.ESyncStatus
> - com.rfidbrasil.core.exception.throwable.AppException
> - com.rfidbrasil.core.model.sync.SyncControlModel
> - com.rfidbrasil.core.model.sync.SyncPacketChunk
> - com.rfidbrasil.core.model.sync.SyncPacketMaster
> - com.rfidbrasil.core.repository.InventoryRepository
> - com.rfidbrasil.core.repository.InventorySituationRepository
> - com.rfidbrasil.core.repository.ItemRepository
> - com.rfidbrasil.core.repository.LocationRepository
> - com.rfidbrasil.core.repository.SyncControlRepository
> - com.rfidbrasil.core.repository.SyncPacketChunkRepository
> - com.rfidbrasil.core.repository.SyncPacketMasterRepository
> - com.rfidbrasil.core.service.DownloadService
> - com.rfidbrasil.core.service.SyncService
> - com.rfidbrasil.core.service.sync.engine.PacketChunkerEngine
> - com.rfidbrasil.core.service.sync.proto.ItemModificationSync
> - com.rfidbrasil.core.service.sync.proto.ItemStatusSync
> - com.rfidbrasil.core.service.sync.proto.MusteringSituationSync
> - com.rfidbrasil.core.service.sync.proto.MusteringSync
> - com.rfidbrasil.core.service.sync.proto.SyncPackage
> - lombok.extern.slf4j.Slf4j
> **Automatically generated documentation** by the Geanky tool.

---

## 1. Quick Summary (API & State)
A high-level overview of the class, its internal state, and available methods.

**Internal State & Dependencies:**


- `private final ` **machineClient** (`MachineClient`)


- `private final ` **situationRepository** (`InventorySituationRepository`)


- `private final ` **inventoryRepository** (`InventoryRepository`)


- `private final ` **musteringBuilderService** (`SyncService`)


- `private final ` **httpClient** (`CloudSyncIntegrationClient`)


- `private final ` **syncControlRepository** (`SyncControlRepository`)


- `private final ` **sseSyncService** (`BoatSseSyncService`)


- `private final ` **downloadService** (`DownloadService`)


- `private final ` **parallelSyncExecutor** (`Executor`)


- `private final ` **itemRepository** (`ItemRepository`)


- `private final ` **masterRepository** (`SyncPacketMasterRepository`)


- `private final ` **chunkRepository** (`SyncPacketChunkRepository`)


- `private final ` **chunkerEngine** (`PacketChunkerEngine`)


- `private final ` **restTemplate** (`RestTemplate`)


- `@Value("${cloud.api.sync.url:http://cloud:8889/api/v2/sync}")` `private ` **cloudUrl** (`String`)


- `@Value("${cloud.api.auth.url:http://cloud:8889/api/v2/auth}")` `private ` **authUrl** (`String`)


- `private ` **destination** (`String`)


**Available Methods:**
- **setDestination(String destination)** ➞ returns `void`
- **latestSync()** ➞ returns `SyncLastestResponse`
- **syncToCloud(SyncRequest request)** ➞ returns `SyncEvent` (throws AppException)
- **authorizeSync()** ➞ returns `String`
- **backgroundProcess(UUID eventKey, SyncRequest request, String token, Long syncTimestamp, Long locationId)** ➞ returns `void`
- **processDownload(UUID eventKey, SyncRequest request, String token, Long locationId)** ➞ returns `void` (throws AppException)
- **getControlWithRetry(UUID eventKey)** ➞ returns `SyncControlModel`
- **processCloudSyncInBackground(UUID eventKey, String token, Long syncTimestamp, Long locationId)** ➞ returns `void`
- **saveMasterAndChunksToDatabase(List<byte[]> slices, String checksum)** ➞ returns `SyncPacketMaster`
- **cleanupOldUploadPackages()** ➞ returns `void`
- **buildProtobufPackage(List<ItemStatusSyncDTO> status, List<MusteringSyncDTO> mustering, List<ItemModificationStatusDTO> modifications)** ➞ returns `SyncPackage`
- **sendStatusToApi(UUID eventKey, ESyncStatus status, String message, String token, Long timestamp, Long locationId)** ➞ returns `boolean`
- **fetchItemStatusTransitions(List<ItemStatusTransitionProjection> transitions)** ➞ returns `List<ItemStatusSyncDTO>`
- **markDataAsSynced(List<ItemStatusTransitionProjection> transitions, List<MusteringSyncDTO> musterings)** ➞ returns `void`


---

## 2. Class Dependencies & State
Visual representation of the internal state and external dependencies this class maintains.

```mermaid
flowchart LR
    classDef classNode fill:#2b3137,stroke:#fff,stroke-width:2px,color:#fff;
    classDef stateNode fill:#f4f6f8,stroke:#d0d7de,color:#24292f;
    classDef extNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    ThisClass["CloudSyncUploadService"]:::classNode

    
    
    ThisClass -- "Depends on" ---> Dep_machineClient["MachineClient"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_situationRepository["InventorySituationRepository"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_inventoryRepository["InventoryRepository"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_musteringBuilderService["SyncService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_httpClient["CloudSyncIntegrationClient"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_syncControlRepository["SyncControlRepository"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_sseSyncService["BoatSseSyncService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_downloadService["DownloadService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_parallelSyncExecutor["Executor"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_itemRepository["ItemRepository"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_masterRepository["SyncPacketMasterRepository"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_chunkRepository["SyncPacketChunkRepository"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_chunkerEngine["PacketChunkerEngine"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_restTemplate["RestTemplate"]:::extNode
    
    
    
    ThisClass -- "Maintains State" --- State_cloudUrl(["String<br>cloudUrl"]):::stateNode
    
    
    
    ThisClass -- "Maintains State" --- State_authUrl(["String<br>authUrl"]):::stateNode
    
    
    
    ThisClass -- "Maintains State" --- State_destination(["String<br>destination"]):::stateNode
    
    
```

---

## 3. Deep Dive (Constructors & Methods)


### 🛠️ Constructors

<details>
<summary><b>CloudSyncUploadService</b>(<i>InventorySituationRepository</i> situationRepository, <i>LocationRepository</i> locationRepository, <i>InventoryRepository</i> inventoryRepository, <i>SyncService</i> musteringBuilderService, <i>CloudSyncIntegrationClient</i> httpClient, <i>SyncControlRepository</i> syncControlRepository, <i>BoatSseSyncService</i> sseSyncService, <i>RestTemplate</i> restTemplate, <i>DownloadService</i> downloadService, <i>ItemRepository</i> itemRepository, <i>PacketChunkerEngine</i> chunkerEngine, <i>SyncPacketChunkRepository</i> syncPacketChunkRepository, <i>SyncPacketMasterRepository</i> syncPacketMasterRepository, <i>Executor</i> parallelSyncExecutor, <i>MachineClient</i> machineClient) (Click to expand)</summary>

> **Signature:**
> `public CloudSyncUploadService(InventorySituationRepository situationRepository, LocationRepository locationRepository, InventoryRepository inventoryRepository, SyncService musteringBuilderService, CloudSyncIntegrationClient httpClient, SyncControlRepository syncControlRepository, BoatSseSyncService sseSyncService, RestTemplate restTemplate, DownloadService downloadService, ItemRepository itemRepository, PacketChunkerEngine chunkerEngine, SyncPacketChunkRepository syncPacketChunkRepository, SyncPacketMasterRepository syncPacketMasterRepository, Executor parallelSyncExecutor, MachineClient machineClient)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: CloudSyncUploadService(situationRepository, locationRepository, inventoryReposit...)

```

**Step-by-Step Logic:**


1. Set 'this.situationRepository' to 'situationRepository'

1. Set 'this.inventoryRepository' to 'inventoryRepository'

1. Set 'this.musteringBuilderService' to 'musteringBuilderService'

1. Set 'this.httpClient' to 'httpClient'

1. Set 'this.syncControlRepository' to 'syncControlRepository'

1. Set 'this.sseSyncService' to 'sseSyncService'

1. Set 'this.parallelSyncExecutor' to 'parallelSyncExecutor'

1. Set 'this.restTemplate' to 'restTemplate'

1. Set 'this.downloadService' to 'downloadService'

1. Set 'this.itemRepository' to 'itemRepository'

1. Set 'this.machineClient' to 'machineClient'

1. Set 'this.chunkerEngine' to 'chunkerEngine'

1. Set 'this.chunkRepository' to 'syncPacketChunkRepository'

1. Set 'this.masterRepository' to 'syncPacketMasterRepository'


</details>




### ⚙️ Methods

<details>
<summary><b>setDestination</b>(<i>String</i> destination) ➞ `void` (Click to expand)</summary>

> **Signature:**
> `public void setDestination(String destination)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: setDestination(destination)

```

**Step-by-Step Logic:**


1. Set 'this.destination' to 'destination'


</details>

<details>
<summary><b>latestSync</b>() ➞ `SyncLastestResponse` (Click to expand)</summary>

> **Signature:**
> `public SyncLastestResponse latestSync()`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant syncControlRepository
    participant syncModelOpt

    Caller->>ThisClass: latestSync()
    ThisClass->>syncControlRepository: findFirstByOrderByIdDesc()
    alt syncModelOpt.isEmpty()
    ThisClass-->>Caller: return new SyncLastestResponse('', ESyncStatus.NONE, 0L, '')
    end
    ThisClass->>syncModelOpt: get()
    ThisClass-->>Caller: return new SyncLastestResponse(eventKey, syncControlModel.getSta...

```

**Step-by-Step Logic:**


1. Declare variable 'syncModelOpt' of type 'Optional<SyncControlModel>' and initialize it with 'Invoke 'syncControlRepository.findFirstByOrderByIdDesc' (no parameters)'

1. If Invoke 'syncModelOpt.isEmpty' (no parameters)
   then:
      - Return the result of: new SyncLastestResponse("", ESyncStatus.NONE, 0L, "")

1. Declare variable 'syncControlModel' of type 'SyncControlModel' and initialize it with 'Invoke 'syncModelOpt.get' (no parameters)'

1. Declare variable 'timestamp' of type 'Long' and initialize it with 'syncControlModel.getLastSyncAt() != null
                ? syncControlModel.getLastSyncAt().toEpochMilli()
                : 0L'

1. Declare variable 'eventKey' of type 'String' and initialize it with 'syncControlModel.getEventKey() != null
                ? syncControlModel.getEventKey().toString()
                : ""'

1. Return the result of: new SyncLastestResponse(eventKey, syncControlModel.getStatus(),
                timestamp, syncControlModel.getErrorMessage())


</details>

<details>
<summary><b>syncToCloud</b>(<i>SyncRequest</i> request) ➞ `SyncEvent` (Click to expand)</summary>

> **Signature:**
> `public SyncEvent syncToCloud(SyncRequest request) throws AppException`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant log
    participant sseSyncService
    participant control
    participant UUID
    participant request
    participant newControl
    participant syncControlRepository

    Caller->>ThisClass: syncToCloud(request)
    ThisClass->>log: info('sync iniciado')
    ThisClass->>ThisClass: authorizeSync()
    ThisClass->>ThisClass: orElseGet(() -> { SyncControlModel newControl = new SyncControlMode...)
    alt ESyncStatus.UPLOAD_IN_PROGRESS.equals(control.getStatus()...
    ThisClass->>sseSyncService: sendSyncStatusOnChange(control.getEventKey(), 'Sync já em andamento...', contro...)
    ThisClass->>control: getEventKey()
    ThisClass->>control: getStatus()
    ThisClass-->>Caller: return new SyncEvent(control.getEventKey().toString(), control.g...
    end
    ThisClass->>control: getLastSyncAt()
    alt request.getEventKey() != null && !request.getEventKey().t...
    ThisClass->>UUID: fromString(request.getEventKey())
    ThisClass->>request: getEventKey()
    else
    ThisClass->>UUID: randomUUID()
    end
    ThisClass->>ThisClass: orElseGet(() -> new SyncControlModel())
    ThisClass->>newControl: setEventKey(targetEventKey)
    ThisClass->>newControl: setStatus(ESyncStatus.UPLOAD_IN_PROGRESS)
    ThisClass->>newControl: setLastSyncAt(lastSuccessfulSyncDate)
    ThisClass->>newControl: setErrorMessage('')
    alt newControl.getCreatedAt() == null
    ThisClass->>newControl: setCreatedAt(new java.util.Date())
    end
    ThisClass->>syncControlRepository: saveAndFlush(newControl)
    ThisClass->>newControl: getEventKey()
    ThisClass->>sseSyncService: sendSyncStatusOnChange(eventKey, '', ESyncStatus.IN_PROGRESS)
    ThisClass->>ThisClass: backgroundProcess(eventKey, request, token, syncTimestamp, request.getLocat...)
    ThisClass->>request: getLocationId()
    ThisClass->>log: info('[SYNC] Retornando eventKey {} para o frontend conectar n...)
    ThisClass-->>Caller: return new SyncEvent(eventKey.toString(), ESyncStatus.IN_PROGRES...

```

**Step-by-Step Logic:**


1. Invoke 'log.info' with parameters: '"sync iniciado"'

1. Declare variable 'token' of type 'String' and initialize it with 'Invoke 'authorizeSync' (no parameters)'

1. Declare variable 'control' of type 'SyncControlModel' and initialize it with 'Invoke 'Invoke 'syncControlRepository.findFirstByOrderByIdDesc' (no parameters).orElseGet' with parameters: '() -> {
                    SyncControlModel newControl = new SyncControlModel();
                    newControl.setEventKey(UUID.randomUUID());
                    return newControl;
                }''

1. If Invoke 'ESyncStatus.UPLOAD_IN_PROGRESS.equals' with parameters: 'Invoke 'control.getStatus' (no parameters)' || Invoke 'ESyncStatus.DOWNLOAD_IN_PROGRESS.equals' with parameters: 'Invoke 'control.getStatus' (no parameters)' || Invoke 'ESyncStatus.IN_PROGRESS.equals' with parameters: 'Invoke 'control.getStatus' (no parameters)'
   then:
      - Declare variable 'lastSync' of type 'Long' and initialize it with 'control.getLastSyncAt() != null ? control.getLastSyncAt().toEpochMilli() : null'
      - Invoke 'sseSyncService.sendSyncStatusOnChange' with parameters: 'Invoke 'control.getEventKey' (no parameters)', '"Sync já em andamento..."', 'Invoke 'control.getStatus' (no parameters)'
      - Return the result of: new SyncEvent(control.getEventKey().toString(), control.getStatus(), lastSync)

1. Declare variable 'lastSuccessfulSyncDate' of type 'Instant' and initialize it with 'Invoke 'control.getLastSyncAt' (no parameters)'

1. Declare variable 'targetEventKey' of type 'UUID'

1. If Invoke 'request.getEventKey' (no parameters) != null && !request.getEventKey().trim().isEmpty()
   then:
      - Set 'targetEventKey' to 'Invoke 'UUID.fromString' with parameters: 'Invoke 'request.getEventKey' (no parameters)''
   else:
      - Set 'targetEventKey' to 'Invoke 'UUID.randomUUID' (no parameters)'

1. Declare variable 'newControl' of type 'SyncControlModel' and initialize it with 'Invoke 'Invoke 'syncControlRepository.findFirstByEventKeyOrderByIdDesc' with parameters: 'targetEventKey'.orElseGet' with parameters: '() -> new SyncControlModel()''

1. Invoke 'newControl.setEventKey' with parameters: 'targetEventKey'

1. Invoke 'newControl.setStatus' with parameters: 'ESyncStatus.UPLOAD_IN_PROGRESS'

1. Invoke 'newControl.setLastSyncAt' with parameters: 'lastSuccessfulSyncDate'

1. Invoke 'newControl.setErrorMessage' with parameters: '""'

1. If Invoke 'newControl.getCreatedAt' (no parameters) == null
   then:
      - Invoke 'newControl.setCreatedAt' with parameters: 'new java.util.Date()'

1. Invoke 'syncControlRepository.saveAndFlush' with parameters: 'newControl'

1. Declare variable 'eventKey' of type 'UUID' and initialize it with 'Invoke 'newControl.getEventKey' (no parameters)'

1. Declare variable 'syncTimestamp' of type 'Long' and initialize it with 'lastSuccessfulSyncDate != null ? lastSuccessfulSyncDate.toEpochMilli()
                : request.getTimestamp()'

1. Invoke 'sseSyncService.sendSyncStatusOnChange' with parameters: 'eventKey', '""', 'ESyncStatus.IN_PROGRESS'

1. Invoke 'backgroundProcess' with parameters: 'eventKey', 'request', 'token', 'syncTimestamp', 'Invoke 'request.getLocationId' (no parameters)'

1. Invoke 'log.info' with parameters: '"[SYNC] Retornando eventKey {} para o frontend conectar no SSE..."', 'eventKey'

1. Return the result of: new SyncEvent(eventKey.toString(), ESyncStatus.IN_PROGRESS, syncTimestamp)


</details>

<details>
<summary><b>authorizeSync</b>() ➞ `String` (Click to expand)</summary>

> **Signature:**
> `public String authorizeSync()`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant headers
    participant restTemplate
    participant mapper
    participant responseEntity

    Caller->>ThisClass: authorizeSync()
    ThisClass->>headers: setContentType(MediaType.APPLICATION_JSON)
    alt try
    ThisClass->>restTemplate: exchange(url, HttpMethod.POST, requestEntity, String.class)
    ThisClass->>mapper: readTree(responseEntity.getBody())
    ThisClass->>responseEntity: getBody()
    ThisClass->>ThisClass: asText()
    ThisClass-->>Caller: return token
    else catch 
    alt e.getStatusCode() == HttpStatus.UNAUTHORIZED || e.getStat...
    ThisClass-->>Caller: throw new RuntimeException('Sync unauthorized')
    end
    ThisClass-->>Caller: throw e
    else catch 
    ThisClass-->>Caller: throw new RuntimeException('Erro ao buscar token: ' + e.getMess...
    end

```

**Step-by-Step Logic:**


1. Declare variable 'req' of type 'AuthRequest' and initialize it with 'new AuthRequest("suporte@cbo.com", "RFIDBrasil")'

1. Declare variable 'url' of type 'String' and initialize it with 'this.destination == null ? authUrl : this.destination + "/auth"'

1. Declare variable 'headers' of type 'HttpHeaders' and initialize it with 'new HttpHeaders()'

1. Invoke 'headers.setContentType' with parameters: 'MediaType.APPLICATION_JSON'

1. Declare variable 'requestEntity' of type 'HttpEntity<?>' and initialize it with 'new HttpEntity<>(req, headers)'

1. Execute a safe block (try) catching potential exceptions


</details>

<details>
<summary><b>backgroundProcess</b>(<i>UUID</i> eventKey, <i>SyncRequest</i> request, <i>String</i> token, <i>Long</i> syncTimestamp, <i>Long</i> locationId) ➞ `void` (Click to expand)</summary>

> **Signature:**
> `private void backgroundProcess(UUID eventKey, SyncRequest request, String token, Long syncTimestamp, Long locationId)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: backgroundProcess(eventKey, request, token, syncTimestamp, locationId)
    ThisClass->>ThisClass: exceptionally(ex -> { Throwable rootCause = ex.getCause() != null ? ex....)

```

**Step-by-Step Logic:**


1. Invoke 'Invoke 'Invoke 'Invoke 'CompletableFuture.runAsync' with parameters: '() -> {
                    try {
                        processCloudSyncInBackground(eventKey, token, syncTimestamp, locationId);
                    } catch (Exception e) {
                        throw new CompletionException(e);
                    }
                }', 'parallelSyncExecutor'.thenRunAsync' with parameters: '() -> {
                    try {
                        SyncControlModel control = getControlWithRetry(eventKey);
                        control.setStatus(ESyncStatus.DOWNLOAD_IN_PROGRESS);
                        syncControlRepository.save(control);
                        sseSyncService.sendSyncStatusOnChange(eventKey,
                                "Upload concluído. Verificando novidades na Nuvem...",
                                ESyncStatus.DOWNLOAD_IN_PROGRESS);

                        log.info("[SYNC-BACKGROUND] Iniciando Fase 2: DOWNLOAD. Verificando nuvem...");
                        PingSyncResponse ping = httpClient.pingCloud(token, this.destination, syncTimestamp);

                        if (ping != null && ping.isDownloadReady()) {
                            log.info("[SYNC-BACKGROUND] Nuvem sinalizou que há dados. Baixando...");
                            request.setTimestamp(ping.getTimestampToDownload());
                            processDownload(eventKey, request, token, locationId);
                        } else {
                            log.info("[SYNC-BACKGROUND] Nenhuma novidade na nuvem. O barco já está atualizado!");
                        }

                    } catch (Exception e) {
                        throw new CompletionException("Falha na etapa de Download: " + e.getMessage(), e);
                    }
                }', 'parallelSyncExecutor'.thenRun' with parameters: '() -> {
                    SyncControlModel control = getControlWithRetry(eventKey);

                    Instant now = Instant.now();
                    control.setStatus(ESyncStatus.SUCCESS_PENDING);
                    control.setErrorMessage("");
                    control.setLastSyncAt(now);
                    control = syncControlRepository.saveAndFlush(control);
                    sseSyncService.sendSyncStatusOnChange(eventKey, "", ESyncStatus.SUCCESS);

                    boolean sendApiSuccess = sendStatusToApi(eventKey, ESyncStatus.SUCCESS, "", token,
                            now.toEpochMilli(),
                            locationId);
                    boolean sendMachineSuccess = machineClient.sendSuccessSync(
                            this.destination == null ? cloudUrl : this.destination,
                            eventKey.toString(),
                            now.toEpochMilli());
                    if (sendApiSuccess && sendMachineSuccess) {
                        control.setStatus(ESyncStatus.SUCCESS);
                        syncControlRepository.save(control);
                    } else {
                        log.error(
                                "Error ao avisar sucesso para api ou machine, na proxima iteracao do machine havera outra tentativa");
                    }

                }'.exceptionally' with parameters: 'ex -> {
                    Throwable rootCause = ex.getCause() != null ? ex.getCause() : ex;
                    log.error("[SYNC-BACKGROUND] Erro durante a sincronização: {}", rootCause.getMessage(), rootCause);

                    try {
                        SyncControlModel control = getControlWithRetry(eventKey);
                        control.setStatus(ESyncStatus.FAILED_PENDING);
                        String errorMsg = rootCause.getMessage() != null ? rootCause.getMessage() : "Erro desconhecido";
                        control.setErrorMessage(errorMsg.length() > 1000 ? errorMsg.substring(0, 1000) : errorMsg);
                        control = syncControlRepository.saveAndFlush(control);
                        sseSyncService.sendSyncStatusOnChange(eventKey, rootCause.getMessage(), ESyncStatus.FAILED);
                        boolean sendApiStatusSuccessfully = sendStatusToApi(eventKey, ESyncStatus.FAILED,
                                rootCause.getMessage(), token,
                                0L, locationId);

                        if (sendApiStatusSuccessfully) {
                            control.setStatus(ESyncStatus.FAILED);
                            syncControlRepository.save(control);
                        } else {
                            log.error(
                                    "Error ao avisar falha para api, na proxima iteracao do machine havera outra tentativa");
                        }

                    } catch (Exception dbEx) {
                        log.error("[SYNC-BACKGROUND] Falha crítica ao salvar status FAILED: {}", dbEx.getMessage());
                    }
                    return null;
                }'


</details>

<details>
<summary><b>processDownload</b>(<i>UUID</i> eventKey, <i>SyncRequest</i> request, <i>String</i> token, <i>Long</i> locationId) ➞ `void` (Click to expand)</summary>

> **Signature:**
> `private void processDownload(UUID eventKey, SyncRequest request, String token, Long locationId) throws AppException`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant downloadService

    Caller->>ThisClass: processDownload(eventKey, request, token, locationId)
    alt try
    ThisClass->>downloadService: download(eventKey, request, token, locationId, this.destination)
    else catch 
    ThisClass-->>Caller: throw new AppException(500, e.getMessage())
    end

```

**Step-by-Step Logic:**


1. Execute a safe block (try) catching potential exceptions


</details>

<details>
<summary><b>getControlWithRetry</b>(<i>UUID</i> eventKey) ➞ `SyncControlModel` (Click to expand)</summary>

> **Signature:**
> `private SyncControlModel getControlWithRetry(UUID eventKey)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant syncControlRepository
    participant Thread

    Caller->>ThisClass: getControlWithRetry(eventKey)
    loop for i < 15
    ThisClass->>syncControlRepository: findFirstByEventKeyOrderByIdDesc(eventKey)
    alt opt.isPresent()
    ThisClass-->>Caller: return opt.get()
    end
    alt try
    ThisClass->>Thread: sleep(200)
    else catch 
    ThisClass->>ThisClass: interrupt()
    Note right of ThisClass: break loop
    end
    end
    ThisClass-->>Caller: throw new RuntimeException('Controle de sync não encontrado no...

```

**Step-by-Step Logic:**


1. Start loop (for) initializing 'Declare variable 'i' of type 'int' and initialize it with '0'', continuing while 'i < 15' is true, and updating 'i++'

1. Throw exception: new RuntimeException("Controle de sync não encontrado no banco para a key: " + eventKey)


</details>

<details>
<summary><b>processCloudSyncInBackground</b>(<i>UUID</i> eventKey, <i>String</i> token, <i>Long</i> syncTimestamp, <i>Long</i> locationId) ➞ `void` (Click to expand)</summary>

> **Signature:**
> `private void processCloudSyncInBackground(UUID eventKey, String token, Long syncTimestamp, Long locationId)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant log
    participant factory
    participant headers
    participant localRestTemplate
    participant masterRepository
    participant pendingMaster
    participant master
    participant situationRepository
    participant musteringBuilderService
    participant chunkerEngine
    participant slices
    participant chunkRepository
    participant unsentChunks
    participant httpClient
    participant chunk
    participant String
    participant sseSyncService
    participant Thread

    Caller->>ThisClass: processCloudSyncInBackground(eventKey, token, syncTimestamp, locationId)
    ThisClass->>log: info('[SYNC-BACKGROUND] Iniciando processamento em background ...)
    alt try
    ThisClass->>factory: setConnectTimeout(5000)
    ThisClass->>factory: setReadTimeout(15000)
    ThisClass->>headers: setContentType(MediaType.APPLICATION_JSON)
    ThisClass->>headers: setBearerAuth(token)
    ThisClass->>localRestTemplate: exchange(url + '/events/trigger', HttpMethod.PUT, requestEntity, S...)
    alt !response.getStatusCode().is2xxSuccessful()
    ThisClass-->>Caller: throw new RuntimeException('Erro na resposta da nuvem: ' + resp...
    end
    ThisClass->>masterRepository: findFirstPending()
    alt pendingMaster.isPresent()
    ThisClass->>pendingMaster: get()
    ThisClass->>log: info('[SYNC-BACKGROUND] Retomando envio de pacote pendente ID:...)
    ThisClass->>master: getId()
    else
    ThisClass->>situationRepository: findStatusTransitions()
    ThisClass->>ThisClass: fetchItemStatusTransitions(transitions)
    ThisClass->>musteringBuilderService: buildPendingMusteringsPayload()
    ThisClass->>ThisClass: toList()
    alt itemStatusPayload.isEmpty() && musteringPayload.isEmpty()...
    ThisClass->>log: info('[SYNC-BACKGROUND] nenhum status, mustering ou alteraçã...)
    ThisClass-->>Caller: return 
    end
    ThisClass->>ThisClass: buildProtobufPackage(itemStatusPayload, musteringPayload, itemModificationStat...)
    ThisClass->>chunkerEngine: compressToGzip(pbPackage)
    ThisClass->>chunkerEngine: sliceIntoChunks(gzipPayload)
    ThisClass->>chunkerEngine: calculateSha256(gzipPayload)
    ThisClass->>ThisClass: saveMasterAndChunksToDatabase(slices, checksum)
    ThisClass->>ThisClass: markDataAsSynced(transitions, musteringPayload)
    ThisClass->>log: info('[SYNC-BACKGROUND] Empacotamento concluído. Gerados {} c...)
    ThisClass->>slices: size()
    end
    ThisClass->>chunkRepository: findUnsyncedByMasterId(master.getId())
    ThisClass->>master: getId()
    ThisClass->>log: info('[SYNC-BACKGROUND] Enviando {} chunks restantes...', unse...)
    ThisClass->>unsentChunks: size()
    loop for each chunk in unsentChunks
    loop while !sucesso
    ThisClass->>httpClient: sendChunk(chunk, master, currentToken, this.destination)
    alt !sucesso
    alt networkFailures % 3 == 0
    alt try
    ThisClass->>ThisClass: authorizeSync()
    else catch 
    ThisClass->>log: warn('[SYNC-BACKGROUND] Tentativa de renovar token offline fal...)
    end
    end
    ThisClass->>log: warn('[SYNC-BACKGROUND] Rede falhou no fragmento {}. Retentand...)
    ThisClass->>chunk: getSequenceNumber()
    ThisClass->>String: format('Conexão instável. Reconectando e enviando pacote %d de...)
    ThisClass->>chunk: getSequenceNumber()
    ThisClass->>master: getTotalChunks()
    ThisClass->>sseSyncService: sendSyncStatusOnChange(eventKey, uiMessage, ESyncStatus.IN_PROGRESS)
    alt try
    ThisClass->>Thread: sleep(delaySeconds * 1000)
    else catch 
    ThisClass->>ThisClass: interrupt()
    ThisClass-->>Caller: throw new RuntimeException('A thread de sincronização foi mor...
    end
    end
    end
    ThisClass->>chunk: setSynced(true)
    ThisClass->>chunkRepository: save(chunk)
    ThisClass->>master: setProcessedChunks(master.getProcessedChunks() + 1)
    ThisClass->>master: getProcessedChunks()
    ThisClass->>masterRepository: save(master)
    end
    ThisClass->>master: setStatus(EPacketStatus.DONE)
    ThisClass->>masterRepository: save(master)
    ThisClass->>ThisClass: cleanupOldUploadPackages()
    ThisClass->>log: info('[SYNC-BACKGROUND] Upload de chunks concluído com sucesso!')
    else catch 
    ThisClass-->>Caller: throw new RuntimeException('Falha no upload: ' + e.getMessage()...
    end

```

**Step-by-Step Logic:**


1. Invoke 'log.info' with parameters: '"[SYNC-BACKGROUND] Iniciando processamento em background para o eventKey: {}"', 'eventKey'

1. Execute a safe block (try) catching potential exceptions


</details>

<details>
<summary><b>saveMasterAndChunksToDatabase</b>(<i>List<byte[]></i> slices, <i>String</i> checksum) ➞ `SyncPacketMaster` (Click to expand)</summary>

> **Signature:**
> `private SyncPacketMaster saveMasterAndChunksToDatabase(List<byte[]> slices, String checksum)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant master
    participant UUID
    participant slices
    participant Instant
    participant chunk
    participant chunks

    Caller->>ThisClass: saveMasterAndChunksToDatabase(slices, checksum)
    ThisClass->>master: setId(UUID.randomUUID())
    ThisClass->>UUID: randomUUID()
    ThisClass->>master: setStatus(EPacketStatus.PENDING)
    ThisClass->>master: setTotalChunks(slices.size())
    ThisClass->>slices: size()
    ThisClass->>master: setProcessedChunks(0)
    ThisClass->>master: setChecksum(checksum)
    ThisClass->>master: setNextRetryAt(Instant.now())
    ThisClass->>Instant: now()
    loop for i < slices.size()
    ThisClass->>chunk: setId(UUID.randomUUID())
    ThisClass->>UUID: randomUUID()
    ThisClass->>chunk: setMaster(master)
    ThisClass->>chunk: setSequenceNumber(i + 1)
    ThisClass->>chunk: setChunkData(slices.get(i))
    ThisClass->>slices: get(i)
    ThisClass->>chunk: setSynced(false)
    ThisClass->>chunks: add(chunk)
    end
    ThisClass->>master: setChunks(chunks)
    ThisClass-->>Caller: return masterRepository.save(master)

```

**Step-by-Step Logic:**


1. Declare variable 'master' of type 'SyncPacketMaster' and initialize it with 'new SyncPacketMaster()'

1. Invoke 'master.setId' with parameters: 'Invoke 'UUID.randomUUID' (no parameters)'

1. Invoke 'master.setStatus' with parameters: 'EPacketStatus.PENDING'

1. Invoke 'master.setTotalChunks' with parameters: 'Invoke 'slices.size' (no parameters)'

1. Invoke 'master.setProcessedChunks' with parameters: '0'

1. Invoke 'master.setChecksum' with parameters: 'checksum'

1. Invoke 'master.setNextRetryAt' with parameters: 'Invoke 'Instant.now' (no parameters)'

1. Declare variable 'chunks' of type 'List<SyncPacketChunk>' and initialize it with 'new ArrayList<>()'

1. Start loop (for) initializing 'Declare variable 'i' of type 'int' and initialize it with '0'', continuing while 'i < Invoke 'slices.size' (no parameters)' is true, and updating 'i++'

1. Invoke 'master.setChunks' with parameters: 'chunks'

1. Return the result of: Invoke 'masterRepository.save' with parameters: 'master'


</details>

<details>
<summary><b>cleanupOldUploadPackages</b>() ➞ `void` (Click to expand)</summary>

> **Signature:**
> `private void cleanupOldUploadPackages()`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant masterRepository
    participant log
    participant oldMasters
    participant e

    Caller->>ThisClass: cleanupOldUploadPackages()
    alt try
    ThisClass->>masterRepository: findOldMastersToCleanup()
    alt oldMasters != null && !oldMasters.isEmpty()
    ThisClass->>log: info('[UPLOAD-CLEANUP] Deletando {} pacotes antigos (Master e ...)
    ThisClass->>oldMasters: size()
    ThisClass->>masterRepository: deleteAll(oldMasters)
    end
    else catch 
    ThisClass->>log: warn('[UPLOAD-CLEANUP] Falha ao tentar limpar pacotes antigos:...)
    ThisClass->>e: getMessage()
    end

```

**Step-by-Step Logic:**


1. Execute a safe block (try) catching potential exceptions


</details>

<details>
<summary><b>buildProtobufPackage</b>(<i>List<ItemStatusSyncDTO></i> status, <i>List<MusteringSyncDTO></i> mustering, <i>List<ItemModificationStatusDTO></i> modifications) ➞ `SyncPackage` (Click to expand)</summary>

> **Signature:**
> `private SyncPackage buildProtobufPackage(List<ItemStatusSyncDTO> status, List<MusteringSyncDTO> mustering, List<ItemModificationStatusDTO> modifications)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant SyncPackage
    participant pb
    participant itemBuilder
    participant musteringBuilder
    participant sitBuilder
    participant modBuilder

    Caller->>ThisClass: buildProtobufPackage(status, mustering, modifications)
    ThisClass->>SyncPackage: newBuilder()
    ThisClass->>pb: setTimestamp(Instant.now().toEpochMilli())
    ThisClass->>ThisClass: toEpochMilli()
    alt status != null
    loop for each s in status
    ThisClass->>ThisClass: setItemId(s.getItemId() != null ? s.getItemId() : 0L)
    ThisClass->>pb: addItemStatus(itemBuilder.build())
    ThisClass->>itemBuilder: build()
    end
    end
    alt mustering != null
    loop for each m in mustering
    ThisClass->>ThisClass: setInventoryType(m.getType() != null ? m.getType().name() : '')
    alt m.getSituations() != null
    loop for each sit in m.getSituations()
    ThisClass->>ThisClass: setReadingDate(sit.getReadingDate() != null ? sit.getReadingDate().toEpo...)
    ThisClass->>musteringBuilder: addSituations(sitBuilder.build())
    ThisClass->>sitBuilder: build()
    end
    end
    ThisClass->>pb: addMustering(musteringBuilder.build())
    ThisClass->>musteringBuilder: build()
    end
    end
    alt modifications != null
    loop for each mod in modifications
    ThisClass->>ThisClass: setLastModifiedDate(mod.getLastMovedDate().toEpochMilli())
    ThisClass->>ThisClass: toEpochMilli()
    ThisClass->>pb: addItemModifications(modBuilder.build())
    ThisClass->>modBuilder: build()
    end
    end
    ThisClass-->>Caller: return pb.build()

```

**Step-by-Step Logic:**


1. Declare variable 'pb' of type 'SyncPackage.Builder' and initialize it with 'Invoke 'SyncPackage.newBuilder' (no parameters)'

1. Invoke 'pb.setTimestamp' with parameters: 'Invoke 'Invoke 'Instant.now' (no parameters).toEpochMilli' (no parameters)'

1. If status != null
   then:
      - Loop through each 's' in the collection 'status'

1. If mustering != null
   then:
      - Loop through each 'm' in the collection 'mustering'

1. If modifications != null
   then:
      - Loop through each 'mod' in the collection 'modifications'

1. Return the result of: Invoke 'pb.build' (no parameters)


</details>

<details>
<summary><b>sendStatusToApi</b>(<i>UUID</i> eventKey, <i>ESyncStatus</i> status, <i>String</i> message, <i>String</i> token, <i>Long</i> timestamp, <i>Long</i> locationId) ➞ `boolean` (Click to expand)</summary>

> **Signature:**
> `public boolean sendStatusToApi(UUID eventKey, ESyncStatus status, String message, String token, Long timestamp, Long locationId)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant headers
    participant restTemplate
    participant log
    participant e

    Caller->>ThisClass: sendStatusToApi(eventKey, status, message, token, timestamp, locationId)
    alt try
    ThisClass->>headers: setContentType(MediaType.APPLICATION_JSON)
    ThisClass->>headers: setBearerAuth(token)
    ThisClass->>restTemplate: exchange(url, HttpMethod.PUT, requestEntity, String.class)
    alt response.getStatusCode().is2xxSuccessful()
    ThisClass->>log: info('[SYNC-BACKGROUND] Status {} avisado para a API com suces...)
    ThisClass-->>Caller: return true
    end
    else catch 
    ThisClass->>log: error('[SYNC-BACKGROUND] Falha ao avisar a API sobre o status {...)
    ThisClass->>e: getMessage()
    ThisClass-->>Caller: return false
    end
    ThisClass-->>Caller: return false

```

**Step-by-Step Logic:**


1. Execute a safe block (try) catching potential exceptions

1. Return the result of: false


</details>

<details>
<summary><b>fetchItemStatusTransitions</b>(<i>List<ItemStatusTransitionProjection></i> transitions) ➞ `List<ItemStatusSyncDTO>` (Click to expand)</summary>

> **Signature:**
> `private List<ItemStatusSyncDTO> fetchItemStatusTransitions(List<ItemStatusTransitionProjection> transitions)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: fetchItemStatusTransitions(transitions)
    ThisClass-->>Caller: return transitions.stream().map(t -> { ItemStatusSyncDTO dto = n...

```

**Step-by-Step Logic:**


1. Return the result of: Invoke 'Invoke 'Invoke 'transitions.stream' (no parameters).map' with parameters: 't -> {
            ItemStatusSyncDTO dto = new ItemStatusSyncDTO();
            dto.setEpc(t.getEpc());
            dto.setStatus(t.getSituation());
            dto.setReadingDate(t.getReadingDate());
            dto.setAntennaNumber(t.getAntennaNumber());
            dto.setPortalMac(t.getPortalMac());
            dto.setInventoryId(t.getInventoryId());
            dto.setItemId(t.getItemId());
            return dto;
        }'.collect' with parameters: 'Invoke 'Collectors.toList' (no parameters)'


</details>

<details>
<summary><b>markDataAsSynced</b>(<i>List<ItemStatusTransitionProjection></i> transitions, <i>List<MusteringSyncDTO></i> musterings) ➞ `void` (Click to expand)</summary>

> **Signature:**
> `private void markDataAsSynced(List<ItemStatusTransitionProjection> transitions, List<MusteringSyncDTO> musterings)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass
    participant Collectors
    participant situationRepository
    participant inventoryRepository

    Caller->>ThisClass: markDataAsSynced(transitions, musterings)
    alt !transitions.isEmpty()
    ThisClass->>ThisClass: collect(Collectors.toList())
    ThisClass->>Collectors: toList()
    ThisClass->>situationRepository: markAsSynced(situationIds)
    end
    alt !musterings.isEmpty()
    ThisClass->>ThisClass: collect(Collectors.toList())
    ThisClass->>Collectors: toList()
    ThisClass->>inventoryRepository: markAsSyncedByNames(inventoryNames)
    end

```

**Step-by-Step Logic:**


1. If !transitions.isEmpty()
   then:
      - Declare variable 'situationIds' of type 'List<Long>' and initialize it with 'Invoke 'Invoke 'Invoke 'transitions.stream' (no parameters).map' with parameters: 'ItemStatusTransitionProjection::getId'.collect' with parameters: 'Invoke 'Collectors.toList' (no parameters)''
      - Invoke 'situationRepository.markAsSynced' with parameters: 'situationIds'

1. If !musterings.isEmpty()
   then:
      - Declare variable 'inventoryNames' of type 'List<String>' and initialize it with 'Invoke 'Invoke 'Invoke 'musterings.stream' (no parameters).map' with parameters: 'MusteringSyncDTO::getInventoryName'.collect' with parameters: 'Invoke 'Collectors.toList' (no parameters)''
      - Invoke 'inventoryRepository.markAsSyncedByNames' with parameters: 'inventoryNames'


</details>


