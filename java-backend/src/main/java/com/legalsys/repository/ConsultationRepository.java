package com.legalsys.repository;

import com.legalsys.entity.Consultation;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;

@Repository
public interface ConsultationRepository extends JpaRepository<Consultation, String> {

    Optional<Consultation> findByTicketNo(String ticketNo);

    Page<Consultation> findBySubmitterId(String submitterId, Pageable pageable);

    Page<Consultation> findBySubmitterIdAndStatus(String submitterId, String status, Pageable pageable);

    Page<Consultation> findByStatus(String status, Pageable pageable);

    Page<Consultation> findByHandlerId(String handlerId, Pageable pageable);

    Page<Consultation> findByHandlerIdAndStatus(String handlerId, String status, Pageable pageable);

    @Query("SELECT c FROM Consultation c WHERE c.status = :status AND c.urgency = :urgency")
    Page<Consultation> findByStatusAndUrgency(@Param("status") String status, @Param("urgency") String urgency, Pageable pageable);

    // 统计相关
    long countByStatus(String status);

    long countByStatusAndSubmittedAtBetween(String status, LocalDateTime start, LocalDateTime end);

    long countByHandlerIdAndStatus(String handlerId, String status);

    long countByHandlerIdAndClosedAtBetween(String handlerId, LocalDateTime start, LocalDateTime end);

    long countByInternalCategory(String internalCategory);

    @Query("SELECT COUNT(c) FROM Consultation c WHERE c.submittedAt >= :start AND c.submittedAt < :end")
    long countSubmittedBetween(@Param("start") LocalDateTime start, @Param("end") LocalDateTime end);

    @Query("SELECT COUNT(c) FROM Consultation c WHERE c.closedAt >= :start AND c.closedAt < :end")
    long countClosedBetween(@Param("start") LocalDateTime start, @Param("end") LocalDateTime end);

    // 搜索
    @Query("SELECT c FROM Consultation c WHERE " +
           "LOWER(c.title) LIKE LOWER(CONCAT('%', :keyword, '%')) OR " +
           "LOWER(c.description) LIKE LOWER(CONCAT('%', :keyword, '%'))")
    Page<Consultation> search(@Param("keyword") String keyword, Pageable pageable);

    // 相似问题
    @Query("SELECT c FROM Consultation c WHERE " +
           "LOWER(c.title) LIKE LOWER(CONCAT('%', :keyword, '%')) AND c.id != :excludeId")
    List<Consultation> findSimilar(@Param("keyword") String keyword, @Param("excludeId") String excludeId, Pageable pageable);
}
