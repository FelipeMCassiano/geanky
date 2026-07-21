
> **@RestController**
# 📄 Technical Specification: `BoatSyncController`

> **Package:** sync
> **Dependencies (Imports):**
> - java.io.IOException
> - java.security.Principal
> - java.util.List
> - java.util.Map
> - java.util.UUID
> - org.slf4j.LoggerFactory
> - org.springframework.http.HttpStatus
> - org.springframework.http.MediaType
> - org.springframework.http.ResponseEntity
> - org.springframework.http.codec.ServerSentEvent
> - org.springframework.web.bind.annotation.GetMapping
> - org.springframework.web.bind.annotation.PathVariable
> - org.springframework.web.bind.annotation.PostMapping
> - org.springframework.web.bind.annotation.RequestBody
> - org.springframework.web.bind.annotation.RequestParam
> - org.springframework.web.bind.annotation.RestController
> - com.rfidbrasil.core.dto.request.SyncRequest
> - com.rfidbrasil.core.dto.response.SyncLastestResponse
> - com.rfidbrasil.core.service.DownloadService
> - com.rfidbrasil.core.service.SyncService
> - com.rfidbrasil.core.service.sync.BoatSseSyncService
> - [com.rfidbrasil.core.service.sync.CloudSyncUploadService](CloudSyncUploadService.md) 🔗
> - com.rfidbrasil.core.service.sync.MachineSync
> - com.rfidbrasil.core.utils.response.ContentLengthResponseBuilder
> - reactor.core.publisher.Flux
> **Automatically generated documentation** by the Geanky tool.

---

## 1. Quick Summary (API & State)
A high-level overview of the class, its internal state, and available methods.

**Internal State & Dependencies:**


- `private static final ` **log** (`Logger`)


- `private final ` **service** (`SyncService`)


- `private final ` **syncService** ([CloudSyncUploadService](CloudSyncUploadService.md)) 🔗


- `private final ` **sseSyncService** (`BoatSseSyncService`)


- `private final ` **downloadService** (`DownloadService`)


- `private final ` **machineSync** (`MachineSync`)


**Available Methods:**
- **subscEmitter(UUID uuid)** ➞ returns `Flux<ServerSentEvent<String>>`
- **portalReadingsSyncPaginated(Long timestamp, Long id, Principal principal)** ➞ returns `ResponseEntity<?>` (throws IOException)
- **latestSync(Principal principal)** ➞ returns `ResponseEntity<?>` (throws IOException)
- **triggerCloudSync(SyncRequest request, Principal principal)** ➞ returns `ResponseEntity<?>`


---

## 2. Class Dependencies & State
Visual representation of the internal state and external dependencies this class maintains.

```mermaid
flowchart LR
    classDef classNode fill:#2b3137,stroke:#fff,stroke-width:2px,color:#fff;
    classDef stateNode fill:#f4f6f8,stroke:#d0d7de,color:#24292f;
    classDef extNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    ThisClass["BoatSyncController"]:::classNode

    
    
    ThisClass -- "Depends on" ---> Dep_log["Logger"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_service["SyncService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_syncService["CloudSyncUploadService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_sseSyncService["BoatSseSyncService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_downloadService["DownloadService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_machineSync["MachineSync"]:::extNode
    
    
```

---

## 3. Deep Dive (Constructors & Methods)


### 🛠️ Constructors

<details>
<summary><b>BoatSyncController</b>(<i>SyncService</i> service, <i>CloudSyncUploadService</i> syncService, <i>BoatSseSyncService</i> sseSyncService, <i>DownloadService</i> downloadService, <i>MachineSync</i> machineSync) (Click to expand)</summary>

> **Signature:**
> `public BoatSyncController(SyncService service, CloudSyncUploadService syncService, BoatSseSyncService sseSyncService, DownloadService downloadService, MachineSync machineSync)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: BoatSyncController(service, syncService, sseSyncService, downloadService, ma...)

```

**Step-by-Step Logic:**


1. Set 'this.service' to 'service'

1. Set 'this.syncService' to 'syncService'

1. Set 'this.sseSyncService' to 'sseSyncService'

1. Set 'this.downloadService' to 'downloadService'

1. Set 'this.machineSync' to 'machineSync'


</details>




### ⚙️ Methods

<details>
<summary><b>subscEmitter</b>(<i>UUID</i> uuid) ➞ `Flux<ServerSentEvent<String>>` (Click to expand)</summary>

> **Signature:**
> `@GetMapping(value = "/events/{uuid}", produces = MediaType.TEXT_EVENT_STREAM_VALUE)`
> `public Flux<ServerSentEvent<String>> subscEmitter(UUID uuid)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: subscEmitter(uuid)
    ThisClass-->>Caller: return sseSyncService.subscribe(uuid)

```

**Step-by-Step Logic:**


1. Return the result of: Invoke 'sseSyncService.subscribe' with parameters: 'uuid'


</details>

<details>
<summary><b>portalReadingsSyncPaginated</b>(<i>Long</i> timestamp, <i>Long</i> id, <i>Principal</i> principal) ➞ `ResponseEntity<?>` (Click to expand)</summary>

> **Signature:**
> `@GetMapping("/portal-readings")`
> `public ResponseEntity<?> portalReadingsSyncPaginated(Long timestamp, Long id, Principal principal) throws IOException`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: portalReadingsSyncPaginated(timestamp, id, principal)
    participant service
    ThisClass->>service: syncPortalReadingsPaginated(timestamp, id, principal)
    alt readings.size() < SyncService.MAX_ITEMS + 1
    ThisClass-->>Caller: return ContentLengthResponseBuilder.createResponse(readings, Htt...
    else
    participant readings
    ThisClass->>readings: remove(SyncService.MAX_ITEMS)
    ThisClass-->>Caller: return ContentLengthResponseBuilder.createResponse(readings, Htt...
    end

```

**Step-by-Step Logic:**


1. Declare variable 'readings' of type 'List<Map<String, Object>>' and initialize it with 'Invoke 'service.syncPortalReadingsPaginated' with parameters: 'timestamp', 'id', 'principal''

1. If Invoke 'readings.size' (no parameters) < SyncService.MAX_ITEMS + 1
   then:
      - Return the result of: Invoke 'ContentLengthResponseBuilder.createResponse' with parameters: 'readings', 'HttpStatus.OK', 'principal'
   else:
      - Declare variable 'nextOffset' of type 'Long' and initialize it with '(Long) readings.get(SyncService.MAX_ITEMS).get("id")'
      - Invoke 'readings.remove' with parameters: 'SyncService.MAX_ITEMS'
      - Return the result of: Invoke 'ContentLengthResponseBuilder.createResponse' with parameters: 'readings', 'HttpStatus.PARTIAL_CONTENT', 'principal', '"v2/sync/portal-readings/page?timestamp=0&offset=" + nextOffset'


</details>

<details>
<summary><b>latestSync</b>(<i>Principal</i> principal) ➞ `ResponseEntity<?>` (Click to expand)</summary>

> **Signature:**
> `@GetMapping("/latest")`
> `public ResponseEntity<?> latestSync(Principal principal) throws IOException`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: latestSync(principal)
    participant syncService
    ThisClass->>syncService: latestSync()
    ThisClass-->>Caller: return ContentLengthResponseBuilder.ok(response, principal)

```

**Step-by-Step Logic:**


1. Declare variable 'response' of type 'SyncLastestResponse' and initialize it with 'Invoke 'syncService.latestSync' (no parameters)'

1. Return the result of: Invoke 'ContentLengthResponseBuilder.ok' with parameters: 'response', 'principal'


</details>

<details>
<summary><b>triggerCloudSync</b>(<i>SyncRequest</i> request, <i>Principal</i> principal) ➞ `ResponseEntity<?>` (Click to expand)</summary>

> **Signature:**
> `@PostMapping("/cloud")`
> `public ResponseEntity<?> triggerCloudSync(SyncRequest request, Principal principal)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: triggerCloudSync(request, principal)
    alt try
    ThisClass-->>Caller: return ContentLengthResponseBuilder.ok(machineSync.executeSync(t...
    else catch 
    participant log
    ThisClass->>log: error('[REST] Erro ao tentar acionar o sync manual: {}', e.getM...)
    participant e
    ThisClass->>e: getMessage()
    ThisClass-->>Caller: return ResponseEntity.internalServerError().body('Falha ao inici...
    end

```

**Step-by-Step Logic:**


1. Execute a safe block (try) catching potential exceptions


</details>


