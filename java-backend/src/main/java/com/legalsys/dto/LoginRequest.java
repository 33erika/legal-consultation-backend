package com.legalsys.dto;

import jakarta.validation.constraints.NotBlank;
import lombok.Data;

@Data
public class LoginRequest {
    @NotBlank(message = "工号不能为空")
    private String employeeId;

    @NotBlank(message = "密码不能为空")
    private String password;
}
