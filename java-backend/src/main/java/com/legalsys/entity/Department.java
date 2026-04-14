package com.legalsys.entity;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "departments")
public class Department {

    @Id
    private String id;

    @Column(nullable = false)
    private String name;

    @Column(name = "parent_id")
    private String parentId;
}
