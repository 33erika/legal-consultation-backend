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
@Table(name = "users")
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    @Column(name = "employee_id", unique = true, nullable = false)
    private String employeeId;  // 工号

    @Column(nullable = false)
    private String password;

    @Column(nullable = false)
    private String name;

    private String email;

    @Column(nullable = false)
    private String role;  // employee, legal_staff, legal_head, supervisor, admin

    @Column(name = "department_id")
    private String departmentId;

    @Column(nullable = false)
    private String status = "active";  // active, inactive

    @Column(name = "created_at")
    private LocalDateTime createdAt = LocalDateTime.now();

    @Column(name = "updated_at")
    private LocalDateTime updatedAt = LocalDateTime.now();
}
