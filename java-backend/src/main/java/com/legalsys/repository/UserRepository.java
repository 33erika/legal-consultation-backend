package com.legalsys.repository;

import com.legalsys.entity.User;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.List;
import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<User, String> {

    Optional<User> findByEmployeeId(String employeeId);

    Page<User> findByEmployeeIdContainingOrNameContaining(String employeeId, String name, Pageable pageable);

    List<User> findByRole(String role);

    List<User> findByRoleAndStatus(String role, String status);

    long countByRole(String role);
}
