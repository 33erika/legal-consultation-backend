package com.legalsys.entity;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;
import java.time.LocalDateTime;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "consultations")
public class Consultation {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @Column(name = "ticket_no", unique = true, nullable = false)
    private String ticketNo;  // 工单号 CONS-20260413-001

    @Column(nullable = false)
    private String title;

    @Column(columnDefinition = "TEXT")
    private String description;

    @Column(nullable = false)
    private String category;  // complaint, contract, labor, ip, dispute, other

    @Column(nullable = false)
    private String urgency;  // normal, high, urgent, very_urgent

    @Column(nullable = false)
    private String status;  // pending, accepted, in_progress, waiting_supplement, completed, closed

    // 内部分类
    @Column(name = "internal_category")
    private String internalCategory;  // simple, complex

    @Column(name = "complex_sub_category")
    private String complexSubCategory;  // dispute, contract, labor, complaint, ip, other

    // 提交人信息
    @Column(name = "submitter_id", nullable = false)
    private String submitterId;

    @Transient
    private User submitter;

    // 处理人信息
    @Column(name = "handler_id")
    private String handlerId;

    @Transient
    private User handler;

    // 时间戳
    @Column(name = "submitted_at", nullable = false)
    private LocalDateTime submittedAt = LocalDateTime.now();

    @Column(name = "accepted_at")
    private LocalDateTime acceptedAt;

    @Column(name = "first_replied_at")
    private LocalDateTime firstRepliedAt;

    @Column(name = "closed_at")
    private LocalDateTime closedAt;

    @Column(name = "updated_at")
    private LocalDateTime updatedAt = LocalDateTime.now();

    // 评分
    private Integer rating;

    // 扩展字段 (JSON格式存储)
    @Column(columnDefinition = "TEXT")
    private String extraData;
}
