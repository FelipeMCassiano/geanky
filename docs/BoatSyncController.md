
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
- **portalReadingsSyncPaginated(Long timestamp, Long id, Principal principal)** ➞ returns `ResponseEntity<?>`
- **latestSync(Principal principal)** ➞ returns `ResponseEntity<?>`
- **triggerCloudSync(SyncRequest request, Principal principal)** ➞ returns `ResponseEntity<?>`


---

## 2. Class Dependencies & State
Visual representation of the internal state and external dependencies this class maintains.

```mermaid
flowchart LR
    %% Styling
    classDef classNode fill:#2b3137,stroke:#fff,stroke-width:2px,color:#fff;
    classDef stateNode fill:#f4f6f8,stroke:#d0d7de,color:#24292f;
    classDef extNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    ThisClass["BoatSyncController"]:::classNode

    %% State vs External Dependencies
    
    
    ThisClass -- "Depends on" ---> Dep_log["Logger"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_service["SyncService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_syncService["CloudSyncUploadService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_sseSyncService["BoatSseSyncService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_downloadService["DownloadService"]:::extNode
    
    
    
    ThisClass -- "Depends on" ---> Dep_machineSync["MachineSync"]:::extNode
    
    
```

---

## 3. Deep Dive (Constructors & Methods)
Expand the sections below to read the exact pseudo-code and business rules.


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

**Parameters:**

- **service** (`SyncService`)

- **syncService** (`CloudSyncUploadService`)

- **sseSyncService** (`BoatSseSyncService`)

- **downloadService** (`DownloadService`)

- **machineSync** (`MachineSync`)


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

**Parameters:**

- **uuid** (`UUID`)


**Step-by-Step Logic:**



1. Return the result of: Invoke 'sseSyncService.subscribe' with parameters: 'uuid'



</details>

<details>
<summary><b>portalReadingsSyncPaginated</b>(<i>Long</i> timestamp, <i>Long</i> id, <i>Principal</i> principal) ➞ `ResponseEntity<?>` (Click to expand)</summary>

> **Signature:**
> `@GetMapping("/portal-readings")`
> `public ResponseEntity<?> portalReadingsSyncPaginated(Long timestamp, Long id, Principal principal)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: portalReadingsSyncPaginated(timestamp, id, principal)
    alt readings.size() < SyncService.MAX_ITEMS + ...
    ThisClass-->>Caller: return ContentLengthResponseBuilder.createResponse(readings, Htt...
    end

```

**Parameters:**

- **timestamp** (`Long`)

- **id** (`Long`)

- **principal** (`Principal`)


**Step-by-Step Logic:**



1. If Invoke 'readings.size' (no parameters) is less than SyncService.MAX_ITEMS plus 1
   then:
      - Return the result of: Invoke 'ContentLengthResponseBuilder.createResponse' with parameters: 'readings', 'HttpStatus.OK', 'principal'



</details>

<details>
<summary><b>latestSync</b>(<i>Principal</i> principal) ➞ `ResponseEntity<?>` (Click to expand)</summary>

> **Signature:**
> `@GetMapping("/latest")`
> `public ResponseEntity<?> latestSync(Principal principal)`

**Sequence Diagram:**
```mermaid
sequenceDiagram
    actor Caller
    participant ThisClass

    Caller->>ThisClass: latestSync(principal)
    ThisClass-->>Caller: return ContentLengthResponseBuilder.ok(response, principal)

```

**Parameters:**

- **principal** (`Principal`)


**Step-by-Step Logic:**



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

```

**Parameters:**

- **request** (`SyncRequest`)

- **principal** (`Principal`)


**Step-by-Step Logic:**
> *Empty body.*

</details>


