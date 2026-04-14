package com.legalsys.entity;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "notification_configs")
public class NotificationConfig {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @Column(name = "config_key", unique = true, nullable = false)
    private String configKey;

    @Column(columnDefinition = "TEXT")
    private String configValue;

    private String description;
}
