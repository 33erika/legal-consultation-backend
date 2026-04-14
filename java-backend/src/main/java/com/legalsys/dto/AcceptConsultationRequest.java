package com.legalsys.dto;

import lombok.Data;

@Data
public class AcceptConsultationRequest {
    private String internalCategory = "simple";
    private String complexSubCategory;
}
