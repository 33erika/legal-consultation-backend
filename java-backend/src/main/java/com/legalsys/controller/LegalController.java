package com.legalsys.controller;

import com.legalsys.dto.ApiResponse;
import com.legalsys.dto.PageResponse;
import com.legalsys.entity.Consultation;
import com.legalsys.entity.User;
import com.legalsys.service.ConsultationService;
import com.legalsys.service.UserService;
import com.legalsys.util.JwtUtil;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/legal")
@RequiredArgsConstructor
public class LegalController {

    private final ConsultationService consultationService;
    private final UserService userService;
    private final JwtUtil jwtUtil;

    @GetMapping("/dashboard")
    public ResponseEntity<ApiResponse<Map<String, Object>>> dashboard(
            @RequestHeader("Authorization") String authHeader) {
        Map<String, Object> stats = consultationService.getStats();
        return ResponseEntity.ok(ApiResponse.success(stats));
    }

    @GetMapping("/consultation-pool")
    public ResponseEntity<ApiResponse<PageResponse<Consultation>>> consultationPool(
            @RequestParam(required = false) String urgency,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "20") int pageSize) {
        PageRequest pageRequest = PageRequest.of(page - 1, pageSize, Sort.by(
                Sort.Order.desc("urgency"),
                Sort.Order.desc("submittedAt")
        ));
        Page<Consultation> pageResult = consultationService.getConsultationPool(urgency, pageRequest);

        PageResponse<Consultation> response = PageResponse.of(
                pageResult.getContent(),
                pageResult.getTotalElements(),
                page,
                pageSize
        );
        return ResponseEntity.ok(ApiResponse.success(response));
    }

    @GetMapping("/my-tasks")
    public ResponseEntity<ApiResponse<PageResponse<Consultation>>> myTasks(
            @RequestHeader("Authorization") String authHeader,
            @RequestParam(required = false) String status,
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "20") int pageSize) {
        String token = authHeader.replace("Bearer ", "");
        String userId = jwtUtil.getUserIdFromToken(token);

        PageRequest pageRequest = PageRequest.of(page - 1, pageSize, Sort.by(Sort.Direction.DESC, "submittedAt"));
        Page<Consultation> pageResult = consultationService.getMyTasks(userId, status, pageRequest);

        PageResponse<Consultation> response = PageResponse.of(
                pageResult.getContent(),
                pageResult.getTotalElements(),
                page,
                pageSize
        );
        return ResponseEntity.ok(ApiResponse.success(response));
    }

    @GetMapping("/staff-list")
    public ResponseEntity<ApiResponse<List<User>>> staffList() {
        List<User> staff = userService.findLegalStaff();
        return ResponseEntity.ok(ApiResponse.success(staff));
    }
}
