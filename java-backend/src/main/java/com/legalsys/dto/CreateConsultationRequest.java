package com.legalsys.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class CreateConsultationRequest {
    @NotBlank(message = "标题不能为空")
    @Size(max = 100, message = "标题不能超过100字")
    private String title;

    @NotBlank(message = "描述不能为空")
    @Size(max = 5000, message = "描述不能超过5000字")
    private String description;

    private String category = "other";
    private String urgency = "normal";
    private String extraData;
}
