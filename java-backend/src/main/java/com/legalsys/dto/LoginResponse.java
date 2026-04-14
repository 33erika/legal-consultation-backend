package com.legalsys.dto;

import lombok.Data;
import lombok.AllArgsConstructor;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class LoginResponse {
    private String token;
    private UserDTO user;

    @Data
    @AllArgsConstructor
    @NoArgsConstructor
    public static class UserDTO {
        private String id;
        private String employeeId;
        private String name;
        private String role;
        private String department;
    }
}
