package com.legalsys.repository;

import com.legalsys.entity.ConsultationTypeConfig;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.List;
import java.util.Optional;

@Repository
public interface ConsultationTypeConfigRepository extends JpaRepository<ConsultationTypeConfig, String> {

    Optional<ConsultationTypeConfig> findByType(String type);

    List<ConsultationTypeConfig> findByEnabledTrueOrderBySortOrderAsc();
}
