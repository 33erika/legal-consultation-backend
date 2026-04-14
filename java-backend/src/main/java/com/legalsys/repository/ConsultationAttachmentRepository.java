package com.legalsys.repository;

import com.legalsys.entity.ConsultationAttachment;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;
import java.util.List;

@Repository
public interface ConsultationAttachmentRepository extends JpaRepository<ConsultationAttachment, String> {

    List<ConsultationAttachment> findByConsultationId(String consultationId);

    void deleteByConsultationId(String consultationId);
}
