package com.legalsys.entity;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "consultation_type_configs")
public class ConsultationTypeConfig {

    @Id
    private String id;

    @Column(unique = true, nullable = false)
    private String type;

    @Column(nullable = false)
    private String name;

    @Column(columnDefinition = "TEXT")
    private String keywords;  // JSON

    @Column(columnDefinition = "TEXT")
    private String fields;  // JSON

    @Column(name = "sort_order")
    private Integer sortOrder = 0;

    @Column(nullable = false)
    private Boolean enabled = true;
}
