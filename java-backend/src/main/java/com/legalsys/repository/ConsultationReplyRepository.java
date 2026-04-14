package com.legalsys.repository;

import com.legalsys.entity.ConsultationReply;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.List;

@Repository
public interface ConsultationReplyRepository extends JpaRepository<ConsultationReply, String> {

    List<ConsultationReply> findByConsultationIdOrderByCreatedAtAsc(String consultationId);

    long countByConsultationId(String consultationId);
}
