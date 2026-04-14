package com.legalsys.controller;

import com.legalsys.dto.*;
import com.legalsys.entity.Consultation;
import com.legalsys.entity.ConsultationReply;
import com.legalsys.service.ConsultationService;
import com.legalsys.util.JwtUtil;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/consultations")
@RequiredArgsConstructor
public class ConsultationController {

    private final ConsultationService consultationService;
    private final JwtUtil jwtUtil;

    @PostMapping
    public ResponseEntity<ApiResponse<Consultation>> create(
            @RequestHeader("Authorization") String authHeader,
            @Valid @RequestBody CreateConsultationRequest request) {
        String token = authHeader.replace("Bearer ", "");
        String userId = jwtUtil.getUserIdFromToken(token);

        Consultation consultation = new Consultation();
        consultation.setSubmitterId(userId);
        consultation.setTitle(request.getTitle());
        consultation.setDescription(request.getDescription());
        consultation.setCategory(request.getCategory());
        consultation.setUrgency(request.getUrgency());
        consultation.setExtraData(request.getExtraData());

        Consultation created = consultationService.createConsultation(consultation);
        return ResponseEntity.ok(ApiResponse.success(created));
    }

    @GetMapping
    public ResponseEntity<ApiResponse<PageResponse<Consultation>>> list(
            @RequestParam(required = false) String status,
            @RequestParam(required = false) String category,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "20") int pageSize) {
        PageRequest pageRequest = PageRequest.of(page - 1, pageSize, Sort.by(Sort.Direction.DESC, "submittedAt"));
        Page<Consultation> pageResult = consultationService.getConsultations(null, status, null, category, pageRequest);

        PageResponse<Consultation> response = PageResponse.of(
                pageResult.getContent(),
                pageResult.getTotalElements(),
                page,
                pageSize
        );
        return ResponseEntity.ok(ApiResponse.success(response));
    }

    @GetMapping("/{id}")
    public ResponseEntity<ApiResponse<Consultation>> get(@PathVariable String id) {
        return consultationService.getById(id)
                .map(c -> ResponseEntity.ok(ApiResponse.success(c)))
                .orElse(ResponseEntity.status(404).body(ApiResponse.error(404, "咨询不存在")));
    }

    @GetMapping("/{id}/replies")
    public ResponseEntity<ApiResponse<List<ConsultationReply>>> getReplies(@PathVariable String id) {
        List<ConsultationReply> replies = consultationService.getReplies(id);
        return ResponseEntity.ok(ApiResponse.success(replies));
    }

    @PostMapping("/{id}/accept")
    public ResponseEntity<ApiResponse<Consultation>> accept(
            @RequestHeader("Authorization") String authHeader,
            @PathVariable String id,
            @RequestBody AcceptConsultationRequest request) {
        String token = authHeader.replace("Bearer ", "");
        String userId = jwtUtil.getUserIdFromToken(token);

        Consultation result = consultationService.accept(id, userId,
                request.getInternalCategory(), request.getComplexSubCategory());
        if (result == null) {
            return ResponseEntity.status(404).body(ApiResponse.error(404, "咨询不存在"));
        }
        return ResponseEntity.ok(ApiResponse.success(result));
    }

    @PostMapping("/{id}/reply")
    public ResponseEntity<ApiResponse<ConsultationReply>> reply(
            @RequestHeader("Authorization") String authHeader,
            @PathVariable String id,
            @Valid @RequestBody ReplyRequest request) {
        String token = authHeader.replace("Bearer ", "");
        String userId = jwtUtil.getUserIdFromToken(token);

        ConsultationReply reply = consultationService.reply(id, userId, request.getContent());
        if (reply == null) {
            return ResponseEntity.status(404).body(ApiResponse.error(404, "咨询不存在"));
        }
        return ResponseEntity.ok(ApiResponse.success(reply));
    }

    @PostMapping("/{id}/request-supplement")
    public ResponseEntity<ApiResponse<Void>> requestSupplement(
            @PathVariable String id,
            @RequestBody Map<String, String> body) {
        consultationService.requestSupplement(id, body.get("content"));
        return ResponseEntity.ok(ApiResponse.success());
    }

    @PostMapping("/{id}/close")
    public ResponseEntity<ApiResponse<Consultation>> close(
            @PathVariable String id,
            @RequestBody CloseConsultationRequest request) {
        Consultation result = consultationService.close(id,
                request.getInternalCategory(), request.getComplexSubCategory());
        if (result == null) {
            return ResponseEntity.status(404).body(ApiResponse.error(404, "咨询不存在"));
        }
        return ResponseEntity.ok(ApiResponse.success(result));
    }

    @PostMapping("/{id}/transfer")
    public ResponseEntity<ApiResponse<Consultation>> transfer(
            @PathVariable String id,
            @RequestBody TransferRequest request) {
        Consultation result = consultationService.transfer(id, request.getNewHandlerId(), request.getReason());
        if (result == null) {
            return ResponseEntity.status(404).body(ApiResponse.error(404, "咨询不存在"));
        }
        return ResponseEntity.ok(ApiResponse.success(result));
    }

    @PostMapping("/{id}/rate")
    public ResponseEntity<ApiResponse<Void>> rate(
            @PathVariable String id,
            @RequestBody Map<String, Integer> body) {
        consultationService.rate(id, body.get("rating"));
        return ResponseEntity.ok(ApiResponse.success());
    }

    @GetMapping("/{id}/similar")
    public ResponseEntity<ApiResponse<List<Consultation>>> similar(
            @PathVariable String id,
            @RequestParam(defaultValue = "5") int limit) {
        return consultationService.getById(id)
                .map(c -> {
                    List<Consultation> similar = consultationService.getSimilar(id, c.getTitle(), limit);
                    return ResponseEntity.ok(ApiResponse.success(similar));
                })
                .orElse(ResponseEntity.status(404).body(ApiResponse.error(404, "咨询不存在")));
    }

    @GetMapping("/search")
    public ResponseEntity<ApiResponse<PageResponse<Consultation>>> search(
            @RequestParam String keyword,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "20") int pageSize) {
        PageRequest pageRequest = PageRequest.of(page - 1, pageSize, Sort.by(Sort.Direction.DESC, "submittedAt"));
        Page<Consultation> pageResult = consultationService.search(keyword, pageRequest);

        PageResponse<Consultation> response = PageResponse.of(
                pageResult.getContent(),
                pageResult.getTotalElements(),
                page,
                pageSize
        );
        return ResponseEntity.ok(ApiResponse.success(response));
    }
}
