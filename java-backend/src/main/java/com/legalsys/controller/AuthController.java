package com.legalsys.controller;

import com.legalsys.dto.*;
import com.legalsys.entity.User;
import com.legalsys.service.UserService;
import com.legalsys.util.JwtUtil;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/v1/auth")
@RequiredArgsConstructor
public class AuthController {

    private final UserService userService;
    private final JwtUtil jwtUtil;

    @PostMapping("/login")
    public ResponseEntity<ApiResponse<LoginResponse>> login(@Valid @RequestBody LoginRequest request) {
        return userService.findByEmployeeId(request.getEmployeeId())
                .filter(user -> user.getPassword().equals(request.getPassword()))
                .filter(user -> "active".equals(user.getStatus()))
                .map(user -> {
                    String token = jwtUtil.generateToken(user.getId(), user.getEmployeeId(), user.getRole());
                    LoginResponse response = new LoginResponse(token, new LoginResponse.UserDTO(
                            user.getId(),
                            user.getEmployeeId(),
                            user.getName(),
                            user.getRole(),
                            user.getDepartmentId()
                    ));
                    return ResponseEntity.ok(ApiResponse.success(response));
                })
                .orElse(ResponseEntity.status(401).body(ApiResponse.error(401, "工号或密码错误")));
    }

    @PostMapping("/logout")
    public ResponseEntity<ApiResponse<Void>> logout() {
        return ResponseEntity.ok(ApiResponse.success());
    }

    @GetMapping("/me")
    public ResponseEntity<ApiResponse<User>> getCurrentUser(@RequestHeader("Authorization") String authHeader) {
        String token = authHeader.replace("Bearer ", "");
        String userId = jwtUtil.getUserIdFromToken(token);

        return userService.findById(userId)
                .map(user -> ResponseEntity.ok(ApiResponse.success(user)))
                .orElse(ResponseEntity.status(404).body(ApiResponse.error(404, "用户不存在")));
    }
}
