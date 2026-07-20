package com.rfidbrasil.core.service.sync;

import java.net.http.HttpHeaders;
import java.time.Instant;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionException;
import java.util.concurrent.Executor;
import java.util.stream.Collectors;

import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.http.client.SimpleClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.RestTemplate;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.rfidbrasil.core.dto.ItemModificationStatusDTO;
import com.rfidbrasil.core.dto.ItemStatusSyncDTO;
import com.rfidbrasil.core.dto.MusteringSituationSyncDTO;
import com.rfidbrasil.core.dto.MusteringSyncDTO;
import com.rfidbrasil.core.dto.projection.ItemStatusTransitionProjection;
import com.rfidbrasil.core.dto.request.SyncRequest;
import com.rfidbrasil.core.dto.request.UpdateSyncStatusRequest;
import com.rfidbrasil.core.dto.response.PingSyncResponse;
import com.rfidbrasil.core.dto.response.SyncEvent;
import com.rfidbrasil.core.dto.response.SyncLastestResponse;
import com.rfidbrasil.core.enums.EPacketStatus;
import com.rfidbrasil.core.enums.ESyncStatus;
import com.rfidbrasil.core.exception.throwable.AppException;
import com.rfidbrasil.core.model.sync.SyncControlModel;
import com.rfidbrasil.core.model.sync.SyncPacketChunk;
import com.rfidbrasil.core.model.sync.SyncPacketMaster;
import com.rfidbrasil.core.repository.InventoryRepository;
import com.rfidbrasil.core.repository.InventorySituationRepository;
import com.rfidbrasil.core.repository.ItemRepository;
import com.rfidbrasil.core.repository.LocationRepository;
import com.rfidbrasil.core.repository.SyncControlRepository;
import com.rfidbrasil.core.repository.SyncPacketChunkRepository;
import com.rfidbrasil.core.repository.SyncPacketMasterRepository;
import com.rfidbrasil.core.service.DownloadService;
import com.rfidbrasil.core.service.SyncService;
import com.rfidbrasil.core.service.sync.engine.PacketChunkerEngine;
import com.rfidbrasil.core.service.sync.proto.ItemModificationSync;
import com.rfidbrasil.core.service.sync.proto.ItemStatusSync;
import com.rfidbrasil.core.service.sync.proto.MusteringSituationSync;
import com.rfidbrasil.core.service.sync.proto.MusteringSync;
import com.rfidbrasil.core.service.sync.proto.SyncPackage;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@Service
public class CloudSyncUploadService {

    private final MachineClient machineClient;
    private final InventorySituationRepository situationRepository;
    private final InventoryRepository inventoryRepository;
    private final SyncService musteringBuilderService;
    private final CloudSyncIntegrationClient httpClient;
    private final SyncControlRepository syncControlRepository;
    private final BoatSseSyncService sseSyncService;
    private final DownloadService downloadService;
    private final Executor parallelSyncExecutor;
    private final ItemRepository itemRepository;
    private final SyncPacketMasterRepository masterRepository;
    private final SyncPacketChunkRepository chunkRepository;
    private final PacketChunkerEngine chunkerEngine;

    private final RestTemplate restTemplate;

    @Value("${cloud.api.sync.url:http://cloud:8889/api/v2/sync}")
    private String cloudUrl;
    @Value("${cloud.api.auth.url:http://cloud:8889/api/v2/auth}")
    private String authUrl;

    private String destination;

    public CloudSyncUploadService(
            InventorySituationRepository situationRepository,
            LocationRepository locationRepository,
            InventoryRepository inventoryRepository,
            SyncService musteringBuilderService,
            CloudSyncIntegrationClient httpClient,
            SyncControlRepository syncControlRepository,
            BoatSseSyncService sseSyncService,
            RestTemplate restTemplate,
            DownloadService downloadService,
            ItemRepository itemRepository,
            PacketChunkerEngine chunkerEngine,
            SyncPacketChunkRepository syncPacketChunkRepository,
            SyncPacketMasterRepository syncPacketMasterRepository,
            @Qualifier("parallelSyncExecutor") Executor parallelSyncExecutor, MachineClient machineClient) {
        this.situationRepository = situationRepository;
        this.inventoryRepository = inventoryRepository;
        this.musteringBuilderService = musteringBuilderService;
        this.httpClient = httpClient;
        this.syncControlRepository = syncControlRepository;
        this.sseSyncService = sseSyncService;
        this.parallelSyncExecutor = parallelSyncExecutor;
        this.restTemplate = restTemplate;
        this.downloadService = downloadService;
        this.itemRepository = itemRepository;
        this.machineClient = machineClient;
        this.chunkerEngine = chunkerEngine;
        this.chunkRepository = syncPacketChunkRepository;
        this.masterRepository = syncPacketMasterRepository;
    }

    public void setDestination(String destination) {
        this.destination = destination;
    }

    public SyncLastestResponse latestSync() {
        Optional<SyncControlModel> syncModelOpt = syncControlRepository.findFirstByOrderByIdDesc();
        if (syncModelOpt.isEmpty()) {
            return new SyncLastestResponse("", ESyncStatus.NONE, 0L, "");
        }

        SyncControlModel syncControlModel = syncModelOpt.get();

        Long timestamp = syncControlModel.getLastSyncAt() != null
                ? syncControlModel.getLastSyncAt().toEpochMilli()
                : 0L;

        String eventKey = syncControlModel.getEventKey() != null
                ? syncControlModel.getEventKey().toString()
                : "";

        return new SyncLastestResponse(eventKey, syncControlModel.getStatus(),
                timestamp, syncControlModel.getErrorMessage());
    }

    public SyncEvent syncToCloud(SyncRequest request) throws AppException {
        log.info("sync iniciado");
        String token = authorizeSync();

        SyncControlModel control = syncControlRepository.findFirstByOrderByIdDesc()
                .orElseGet(() -> {
                    SyncControlModel newControl = new SyncControlModel();
                    newControl.setEventKey(UUID.randomUUID());
                    return newControl;
                });

        if (ESyncStatus.UPLOAD_IN_PROGRESS.equals(control.getStatus()) ||
                ESyncStatus.DOWNLOAD_IN_PROGRESS.equals(control.getStatus()) ||
                ESyncStatus.IN_PROGRESS.equals(control.getStatus())) {

            Long lastSync = control.getLastSyncAt() != null ? control.getLastSyncAt().toEpochMilli() : null;
            sseSyncService.sendSyncStatusOnChange(control.getEventKey(), "Sync já em andamento...",
                    control.getStatus());
            return new SyncEvent(control.getEventKey().toString(), control.getStatus(), lastSync);
        }

        Instant lastSuccessfulSyncDate = control.getLastSyncAt();

        UUID targetEventKey;
        if (request.getEventKey() != null && !request.getEventKey().trim().isEmpty()) {
            targetEventKey = UUID.fromString(request.getEventKey());
        } else {
            targetEventKey = UUID.randomUUID();
        }

        SyncControlModel newControl = syncControlRepository.findFirstByEventKeyOrderByIdDesc(targetEventKey)
                .orElseGet(() -> new SyncControlModel());

        // 3. Atualiza os dados do controle (reaproveitado ou novo)
        newControl.setEventKey(targetEventKey);
        newControl.setStatus(ESyncStatus.UPLOAD_IN_PROGRESS);
        newControl.setLastSyncAt(lastSuccessfulSyncDate);
        newControl.setErrorMessage("");

        if (newControl.getCreatedAt() == null) {
            newControl.setCreatedAt(new java.util.Date());
        }

        // Salva e garante que está no banco antes de disparar o processo
        syncControlRepository.saveAndFlush(newControl);

        final UUID eventKey = newControl.getEventKey();

        final Long syncTimestamp = lastSuccessfulSyncDate != null ? lastSuccessfulSyncDate.toEpochMilli()
                : request.getTimestamp();

        sseSyncService.sendSyncStatusOnChange(eventKey, "", ESyncStatus.IN_PROGRESS);

        // Dispara o trabalho pesado em background
        backgroundProcess(eventKey, request, token, syncTimestamp, request.getLocationId());

        log.info("[SYNC] Retornando eventKey {} para o frontend conectar no SSE...", eventKey);
        return new SyncEvent(eventKey.toString(), ESyncStatus.IN_PROGRESS, syncTimestamp);
    }

    private record AuthRequest(String email, String pass) {
    }

    public String authorizeSync() {
        AuthRequest req = new AuthRequest("suporte@cbo.com", "RFIDBrasil");

        String url = this.destination == null ? authUrl : this.destination + "/auth";
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<?> requestEntity = new HttpEntity<>(req, headers);

        try {
            ResponseEntity<String> responseEntity = restTemplate.exchange(
                    url,
                    HttpMethod.POST,
                    requestEntity,
                    String.class);

            ObjectMapper mapper = new ObjectMapper();
            JsonNode root = mapper.readTree(responseEntity.getBody());

            String token = root.path("data").path("token").asText();

            return token;

        } catch (HttpClientErrorException e) {
            if (e.getStatusCode() == HttpStatus.UNAUTHORIZED || e.getStatusCode() == HttpStatus.FORBIDDEN) {
                throw new RuntimeException("Sync unauthorized");
            }
            throw e;
        } catch (Exception e) {
            throw new RuntimeException("Erro ao buscar token: " + e.getMessage());
        }
    }

    private void backgroundProcess(UUID eventKey, SyncRequest request, String token, Long syncTimestamp,
            Long locationId) {
        CompletableFuture
                .runAsync(() -> {
                    try {
                        processCloudSyncInBackground(eventKey, token, syncTimestamp, locationId);
                    } catch (Exception e) {
                        throw new CompletionException(e);
                    }
                }, parallelSyncExecutor)
                .thenRunAsync(() -> {
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
                }, parallelSyncExecutor)
                .thenRun(() -> {
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

                })
                .exceptionally(ex -> {
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
                });
    }

    private void processDownload(UUID eventKey, SyncRequest request, String token, Long locationId)
            throws AppException {
        try {
            downloadService.download(eventKey, request, token, locationId, this.destination);
        } catch (Exception e) {
            throw new AppException(500, e.getMessage());
        }
    }

    private SyncControlModel getControlWithRetry(UUID eventKey) {
        for (int i = 0; i < 15; i++) {
            Optional<SyncControlModel> opt = syncControlRepository.findFirstByEventKeyOrderByIdDesc(eventKey);
            if (opt.isPresent()) {
                return opt.get();
            }
            try {
                Thread.sleep(200);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                break;
            }
        }
        throw new RuntimeException("Controle de sync não encontrado no banco para a key: " + eventKey);
    }

    private void processCloudSyncInBackground(UUID eventKey, String token, Long syncTimestamp, Long locationId) {
        log.info("[SYNC-BACKGROUND] Iniciando processamento em background para o eventKey: {}", eventKey);

        try {
            UpdateSyncStatusRequest syncPayload = new UpdateSyncStatusRequest(
                    ESyncStatus.IN_PROGRESS, "",
                    syncTimestamp, eventKey.toString());

            SimpleClientHttpRequestFactory factory = new SimpleClientHttpRequestFactory();
            factory.setConnectTimeout(5000);
            factory.setReadTimeout(15000);

            RestTemplate localRestTemplate = new RestTemplate(factory);

            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);
            headers.setBearerAuth(token);
            HttpEntity<UpdateSyncStatusRequest> requestEntity = new HttpEntity<>(syncPayload, headers);

            String url = this.destination == null ? cloudUrl : this.destination + "/sync";

            ResponseEntity<String> response = localRestTemplate.exchange(url + "/events/trigger",
                    HttpMethod.PUT, requestEntity, String.class);
            if (!response.getStatusCode().is2xxSuccessful()) {
                throw new RuntimeException("Erro na resposta da nuvem: " + response.getBody());
            }

            Optional<SyncPacketMaster> pendingMaster = masterRepository.findFirstPending();
            SyncPacketMaster master;

            if (pendingMaster.isPresent()) {
                master = pendingMaster.get();
                log.info("[SYNC-BACKGROUND] Retomando envio de pacote pendente ID: {}", master.getId());
            } else {
                List<ItemStatusTransitionProjection> transitions = situationRepository.findStatusTransitions();
                List<ItemStatusSyncDTO> itemStatusPayload = fetchItemStatusTransitions(transitions);
                List<MusteringSyncDTO> musteringPayload = musteringBuilderService.buildPendingMusteringsPayload();
                List<ItemModificationStatusDTO> itemModificationStatusPayload = itemRepository
                        .findItemModelsStatusForSync(syncTimestamp).stream()
                        .map(i -> ItemModificationStatusDTO.fromItem(i))
                        .filter(dto -> dto != null)
                        .toList();

                if (itemStatusPayload.isEmpty() && musteringPayload.isEmpty()
                        && itemModificationStatusPayload.isEmpty()) {
                    log.info("[SYNC-BACKGROUND] nenhum status, mustering ou alteração de item para enviar");
                    return;
                }

                SyncPackage pbPackage = buildProtobufPackage(itemStatusPayload, musteringPayload,
                        itemModificationStatusPayload);
                byte[] gzipPayload = chunkerEngine.compressToGzip(pbPackage);
                List<byte[]> slices = chunkerEngine.sliceIntoChunks(gzipPayload);
                String checksum = chunkerEngine.calculateSha256(gzipPayload);

                master = saveMasterAndChunksToDatabase(slices, checksum);

                markDataAsSynced(transitions, musteringPayload);
                log.info("[SYNC-BACKGROUND] Empacotamento concluído. Gerados {} chunks de 30 bytes.", slices.size());
            }

            List<SyncPacketChunk> unsentChunks = chunkRepository.findUnsyncedByMasterId(master.getId());
            log.info("[SYNC-BACKGROUND] Enviando {} chunks restantes...", unsentChunks.size());

            String currentToken = token;
            for (SyncPacketChunk chunk : unsentChunks) {
                boolean sucesso = false;
                int networkFailures = 0;

                while (!sucesso) {
                    sucesso = httpClient.sendChunk(chunk, master, currentToken, this.destination);

                    if (!sucesso) {
                        networkFailures++;

                        if (networkFailures % 3 == 0) {
                            try {
                                currentToken = authorizeSync();
                            } catch (Exception ignored) {
                                log.warn("[SYNC-BACKGROUND] Tentativa de renovar token offline falhou.");
                            }
                        }

                        long delaySeconds = (long) Math.min(Math.pow(2, networkFailures) * 5, 60);

                        log.warn("[SYNC-BACKGROUND] Rede falhou no fragmento {}. Retentando em {} segundos...",
                                chunk.getSequenceNumber(), delaySeconds);

                        String uiMessage = String.format("Conexão instável. Reconectando e enviando pacote %d de %d...",
                                chunk.getSequenceNumber(), master.getTotalChunks());
                        sseSyncService.sendSyncStatusOnChange(eventKey, uiMessage, ESyncStatus.IN_PROGRESS);

                        try {
                            Thread.sleep(delaySeconds * 1000);
                        } catch (InterruptedException e) {
                            Thread.currentThread().interrupt();
                            throw new RuntimeException("A thread de sincronização foi morta pela aplicação.");
                        }
                    }
                }

                //
                chunk.setSynced(true);
                chunkRepository.save(chunk);
                master.setProcessedChunks(master.getProcessedChunks() + 1);
                masterRepository.save(master);

            }

            master.setStatus(EPacketStatus.DONE);
            masterRepository.save(master);
            cleanupOldUploadPackages();

            log.info("[SYNC-BACKGROUND] Upload de chunks concluído com sucesso!");

        } catch (Exception e) {
            throw new RuntimeException("Falha no upload: " + e.getMessage(), e);
        }
    }

    private SyncPacketMaster saveMasterAndChunksToDatabase(List<byte[]> slices, String checksum) {
        SyncPacketMaster master = new SyncPacketMaster();
        master.setId(UUID.randomUUID());
        master.setStatus(EPacketStatus.PENDING);
        master.setTotalChunks(slices.size());
        master.setProcessedChunks(0);
        master.setChecksum(checksum);
        master.setNextRetryAt(Instant.now());

        List<SyncPacketChunk> chunks = new ArrayList<>();
        for (int i = 0; i < slices.size(); i++) {
            SyncPacketChunk chunk = new SyncPacketChunk();
            chunk.setId(UUID.randomUUID());
            chunk.setMaster(master);
            chunk.setSequenceNumber(i + 1);
            chunk.setChunkData(slices.get(i));
            chunk.setSynced(false);
            chunks.add(chunk);
        }
        master.setChunks(chunks);
        return masterRepository.save(master);
    }

    private void cleanupOldUploadPackages() {
        try {
            List<SyncPacketMaster> oldMasters = masterRepository.findOldMastersToCleanup();
            if (oldMasters != null && !oldMasters.isEmpty()) {
                log.info("[UPLOAD-CLEANUP] Deletando {} pacotes antigos (Master e Chunks)...", oldMasters.size());

                masterRepository.deleteAll(oldMasters);
            }
        } catch (Exception e) {
            log.warn("[UPLOAD-CLEANUP] Falha ao tentar limpar pacotes antigos: {}", e.getMessage());
        }
    }

    private SyncPackage buildProtobufPackage(List<ItemStatusSyncDTO> status, List<MusteringSyncDTO> mustering,
            List<ItemModificationStatusDTO> modifications) {
        SyncPackage.Builder pb = SyncPackage.newBuilder();
        pb.setTimestamp(Instant.now().toEpochMilli());

        if (status != null) {
            for (ItemStatusSyncDTO s : status) {
                ItemStatusSync.Builder itemBuilder = ItemStatusSync.newBuilder()
                        .setEpc(s.getEpc() != null ? s.getEpc() : "")
                        .setStatus(s.getStatus() != null ? s.getStatus() : "")
                        .setReadingDate(s.getReadingDate() != null ? s.getReadingDate().toEpochMilli() : 0L)
                        .setAntennaNumber(s.getAntennaNumber() != null ? s.getAntennaNumber() : 0)
                        .setPortalMac(s.getPortalMac() != null ? s.getPortalMac() : "")
                        .setInventoryId(s.getInventoryId() != null ? s.getInventoryId() : 0L)
                        .setItemId(s.getItemId() != null ? s.getItemId() : 0L);

                pb.addItemStatus(itemBuilder.build());
            }
        }

        if (mustering != null) {
            for (MusteringSyncDTO m : mustering) {
                MusteringSync.Builder musteringBuilder = MusteringSync.newBuilder()
                        .setInventoryName(m.getInventoryName() != null ? m.getInventoryName() : "")
                        .setCreatedAt(m.getCreatedAt() != null ? m.getCreatedAt().toEpochMilli() : 0L)
                        .setCompanyId(m.getCompanyId() != null ? m.getCompanyId() : 0L)
                        .setInventoryType(m.getType() != null ? m.getType().name() : "");

                if (m.getSituations() != null) {
                    for (MusteringSituationSyncDTO sit : m.getSituations()) {
                        MusteringSituationSync.Builder sitBuilder = MusteringSituationSync.newBuilder()
                                .setEpc(sit.getEpc() != null ? sit.getEpc() : "")
                                .setStatus(sit.getStatus() != null ? sit.getStatus() : "")
                                .setReadingDate(
                                        sit.getReadingDate() != null ? sit.getReadingDate().toEpochMilli() : 0L);

                        musteringBuilder.addSituations(sitBuilder.build());
                    }
                }
                pb.addMustering(musteringBuilder.build());
            }
        }
        if (modifications != null) {
            for (ItemModificationStatusDTO mod : modifications) {
                var modBuilder = ItemModificationSync.newBuilder()
                        .setEpc(mod.getEpc() != null ? mod.getEpc() : "")
                        .setStatus(mod.getStatus() != null ? mod.getStatus() : "")
                        .setLastModifiedDate((mod.getLastModifiedDate().toEpochMilli()))
                        .setLastModifiedDate(mod.getLastMovedDate().toEpochMilli());

                pb.addItemModifications(modBuilder.build());
            }
        }

        return pb.build();
    }

    public boolean sendStatusToApi(UUID eventKey, ESyncStatus status, String message, String token, Long timestamp,
            Long locationId) {
        try {

            String baseUrl = this.destination == null ? cloudUrl : this.destination + "/sync";
            String url = baseUrl + "/events/trigger";

            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);
            headers.setBearerAuth(token);

            UpdateSyncStatusRequest payload = new UpdateSyncStatusRequest(status, message, timestamp,
                    eventKey.toString());
            HttpEntity<UpdateSyncStatusRequest> requestEntity = new HttpEntity<>(payload, headers);

            ResponseEntity<String> response = restTemplate.exchange(url, HttpMethod.PUT, requestEntity, String.class);

            if (response.getStatusCode().is2xxSuccessful()) {
                log.info("[SYNC-BACKGROUND] Status {} avisado para a API com sucesso (eventKey: {})", status, eventKey);
                return true;
            }
        } catch (Exception e) {
            log.error("[SYNC-BACKGROUND] Falha ao avisar a API sobre o status {}: {}", status, e.getMessage());
            return false;
        }
        return false;
    }

    private List<ItemStatusSyncDTO> fetchItemStatusTransitions(List<ItemStatusTransitionProjection> transitions) {
        return transitions.stream().map(t -> {
            ItemStatusSyncDTO dto = new ItemStatusSyncDTO();
            dto.setEpc(t.getEpc());
            dto.setStatus(t.getSituation());
            dto.setReadingDate(t.getReadingDate());
            dto.setAntennaNumber(t.getAntennaNumber());
            dto.setPortalMac(t.getPortalMac());
            dto.setInventoryId(t.getInventoryId());
            dto.setItemId(t.getItemId());
            return dto;
        }).collect(Collectors.toList());
    }

    private void markDataAsSynced(List<ItemStatusTransitionProjection> transitions, List<MusteringSyncDTO> musterings) {

        if (!transitions.isEmpty()) {
            List<Long> situationIds = transitions.stream()
                    .map(ItemStatusTransitionProjection::getId)
                    .collect(Collectors.toList());

            situationRepository.markAsSynced(situationIds);
        }

        if (!musterings.isEmpty()) {
            List<String> inventoryNames = musterings.stream()
                    .map(MusteringSyncDTO::getInventoryName)
                    .collect(Collectors.toList());

            inventoryRepository.markAsSyncedByNames(inventoryNames);
        }
    }
}