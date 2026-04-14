package com.legalsys.entity;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "consultation_attachments")
public class ConsultationAttachment {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @Column(name = "consultation_id", nullable = false)
    private String consultationId;

    @Column(name = "attachment_id", nullable = false)
    private String attachmentId;

    @Transient
    private Attachment attachment;
}
