package com.legalsys.service;

import com.legalsys.entity.User;
import com.legalsys.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class UserService {

    private final UserRepository userRepository;

    public Optional<User> findByEmployeeId(String employeeId) {
        return userRepository.findByEmployeeId(employeeId);
    }

    public Optional<User> findById(String id) {
        return userRepository.findById(id);
    }

    public List<User> findByRole(String role) {
        return userRepository.findByRole(role);
    }

    public List<User> findLegalStaff() {
        return userRepository.findByRoleAndStatus("legal_staff", "active");
    }

    @Transactional
    public User createUser(User user) {
        return userRepository.save(user);
    }

    @Transactional
    public User updateUser(User user) {
        user.setUpdatedAt(LocalDateTime.now());
        return userRepository.save(user);
    }

    @Transactional
    public void toggleStatus(String id, String status) {
        userRepository.findById(id).ifPresent(user -> {
            user.setStatus(status);
            user.setUpdatedAt(LocalDateTime.now());
            userRepository.save(user);
        });
    }
}
