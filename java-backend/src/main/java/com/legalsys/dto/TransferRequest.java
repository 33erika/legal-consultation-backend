package com.legalsys.dto;

import lombok.Data;

@Data
public class TransferRequest {
    private String newHandlerId;
    private String reason;
}
