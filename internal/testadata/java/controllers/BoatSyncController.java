package com.rfidbrasil.core.controller.sync;

import java.io.IOException;
import java.security.Principal;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.http.codec.ServerSentEvent;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.rfidbrasil.core.dto.request.SyncRequest;
import com.rfidbrasil.core.dto.response.SyncLastestResponse;
import com.rfidbrasil.core.service.DownloadService;
import com.rfidbrasil.core.service.SyncService;
import com.rfidbrasil.core.service.sync.BoatSseSyncService;
import com.rfidbrasil.core.service.sync.CloudSyncUploadService;
import com.rfidbrasil.core.service.sync.MachineSync;
import com.rfidbrasil.core.utils.response.ContentLengthResponseBuilder;

import reactor.core.publisher.Flux;

@RestController
public class BoatSyncController extends BaseSyncController {

    private static final Logger log = LoggerFactory.getLogger(BoatSyncController.class);
    private final SyncService service;
    private final CloudSyncUploadService syncService;
    private final BoatSseSyncService sseSyncService;
    private final DownloadService downloadService;
    private final MachineSync machineSync;

    public BoatSyncController(SyncService service,
            CloudSyncUploadService syncService,
            BoatSseSyncService sseSyncService,
            DownloadService downloadService, MachineSync machineSync) {
        this.service = service;
        this.syncService = syncService;
        this.sseSyncService = sseSyncService;
        this.downloadService = downloadService;
        this.machineSync = machineSync;
    }

    @GetMapping(value = "/events/{uuid}", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<ServerSentEvent<String>> subscEmitter(@PathVariable("uuid") UUID uuid) {
        return sseSyncService.subscribe(uuid);
    }

    @GetMapping("/portal-readings")
    public ResponseEntity<?> portalReadingsSyncPaginated(
            @RequestParam(name = "timestamp", defaultValue = "0", required = false) Long timestamp,
            @RequestParam(name = "offset", defaultValue = "-1", required = false) Long id,
            Principal principal) throws IOException {

        List<Map<String, Object>> readings = service.syncPortalReadingsPaginated(timestamp, id, principal);

        if (readings.size() < (SyncService.MAX_ITEMS + 1)) {
            return ContentLengthResponseBuilder.createResponse(readings, HttpStatus.OK, principal);
        } else {
            Long nextOffset = (Long) readings.get(SyncService.MAX_ITEMS).get("id");
            readings.remove(SyncService.MAX_ITEMS);

            return ContentLengthResponseBuilder.createResponse(
                    readings,
                    HttpStatus.PARTIAL_CONTENT,
                    principal,
                    "v2/sync/portal-readings/page?timestamp=0&offset=" + nextOffset);
        }
    }

    @GetMapping("/latest")
    public ResponseEntity<?> latestSync(Principal principal) throws IOException {
        SyncLastestResponse response = syncService.latestSync();
        return ContentLengthResponseBuilder.ok(response, principal);
    }

    @PostMapping("/cloud")
    public ResponseEntity<?> triggerCloudSync(@RequestBody SyncRequest request, Principal principal) {
        try {

            return ContentLengthResponseBuilder.ok(machineSync.executeSync(true, request.getTimestamp()));

        } catch (Exception e) {
            log.error("[REST] Erro ao tentar acionar o sync manual: {}", e.getMessage());
            return ResponseEntity.internalServerError().body("Falha ao iniciar sincronização: " + e.getMessage());
        }
    }

    // @PostMapping("/download")
    // public ResponseEntity<?> triggerDownload(@RequestBody SyncRequest request,
    // @RequestParam("event") String event) {
    // downloadService.download(UUID.fromString(event), request,
    // request.getToken());
    // return ResponseEntity.ok(null);
    // }

}