package com.legalsys.service;

import com.legalsys.entity.*;
import com.legalsys.repository.*;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.DayOfWeek;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.format.DateTimeFormatter;
import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class ConsultationService {

    private final ConsultationRepository consultationRepository;
    private final ConsultationReplyRepository replyRepository;
    private final ConsultationAttachmentRepository attachmentRepository;
    private final UserRepository userRepository;
    private final NotificationService notificationService;

    private static final DateTimeFormatter TICKET_DATE_FORMAT = DateTimeFormatter.ofPattern("yyyyMMdd");

    // 创建咨询
    @Transactional
    public Consultation createConsultation(Consultation consultation) {
        // 生成工单号
        String today = LocalDate.now().format(TICKET_DATE_FORMAT);
        long count = consultationRepository.count() + 1;
        String ticketNo = String.format("CONS-%s-%03d", today, count);
        consultation.setTicketNo(ticketNo);
        consultation.setStatus("pending");
        consultation.setSubmittedAt(LocalDateTime.now());

        Consultation saved = consultationRepository.save(consultation);

        // 通知法务
        notificationService.notifyNewConsultation(saved);

        return saved;
    }

    // 获取咨询详情
    public Optional<Consultation> getById(String id) {
        return consultationRepository.findById(id).map(c -> {
            enrichConsultation(c);
            return c;
        });
    }

    // 获取咨询列表
    public Page<Consultation> getConsultations(String submitterId, String status, String handlerId,
                                                String category, Pageable pageable) {
        Page<Consultation> page;
        if (submitterId != null && !submitterId.isEmpty()) {
            page = consultationRepository.findBySubmitterId(submitterId, pageable);
        } else if (status != null && !status.isEmpty()) {
            page = consultationRepository.findByStatus(status, pageable);
        } else if (handlerId != null && !handlerId.isEmpty()) {
            page = consultationRepository.findByHandlerId(handlerId, pageable);
        } else {
            page = consultationRepository.findAll(pageable);
        }
        page.forEach(this::enrichConsultation);
        return page;
    }

    // 获取咨询池（待处理）
    public Page<Consultation> getConsultationPool(String urgency, Pageable pageable) {
        Page<Consultation> page;
        if (urgency != null && !urgency.isEmpty()) {
            page = consultationRepository.findByStatusAndUrgency("pending", urgency, pageable);
        } else {
            page = consultationRepository.findByStatus("pending", pageable);
        }
        page.forEach(this::enrichConsultation);
        return page;
    }

    // 获取我的待办
    public Page<Consultation> getMyTasks(String handlerId, String status, Pageable pageable) {
        Page<Consultation> page;
        if (status != null && !status.isEmpty()) {
            page = consultationRepository.findByHandlerIdAndStatus(handlerId, status, pageable);
        } else {
            page = consultationRepository.findByHandlerId(handlerId, pageable);
        }
        page.forEach(this::enrichConsultation);
        return page;
    }

    // 接单
    @Transactional
    public Consultation accept(String id, String handlerId, String internalCategory, String complexSubCategory) {
        return consultationRepository.findById(id).map(c -> {
            c.setHandlerId(handlerId);
            c.setAcceptedAt(LocalDateTime.now());
            c.setStatus("in_progress");
            c.setInternalCategory(internalCategory != null ? internalCategory : "simple");
            c.setComplexSubCategory(complexSubCategory);
            c.setUpdatedAt(LocalDateTime.now());
            return enrichConsultation(consultationRepository.save(c));
        }).orElse(null);
    }

    // 回复
    @Transactional
    public ConsultationReply reply(String consultationId, String userId, String content) {
        Consultation consultation = consultationRepository.findById(consultationId).orElse(null);
        if (consultation == null) return null;

        // 创建回复记录
        ConsultationReply reply = new ConsultationReply();
        reply.setConsultationId(consultationId);
        reply.setUserId(userId);
        reply.setContent(content);
        reply.setType("reply");
        reply.setCreatedAt(LocalDateTime.now());
        replyRepository.save(reply);

        // 更新咨询状态
        if (consultation.getFirstRepliedAt() == null) {
            consultation.setFirstRepliedAt(LocalDateTime.now());
        }
        consultation.setStatus("waiting_supplement");
        consultation.setUpdatedAt(LocalDateTime.now());
        enrichConsultation(consultationRepository.save(consultation));

        // 通知员工
        notificationService.notifyConsultationReplied(consultation);

        return reply;
    }

    // 要求补充资料
    @Transactional
    public void requestSupplement(String consultationId, String content) {
        consultationRepository.findById(consultationId).ifPresent(c -> {
            // 创建系统消息
            ConsultationReply reply = new ConsultationReply();
            reply.setConsultationId(consultationId);
            reply.setUserId("system");
            reply.setContent("法务要求补充以下资料：" + content);
            reply.setType("system");
            reply.setCreatedAt(LocalDateTime.now());
            replyRepository.save(reply);

            c.setStatus("waiting_supplement");
            c.setUpdatedAt(LocalDateTime.now());
            consultationRepository.save(c);

            notificationService.notifyConsultationReplied(c);
        });
    }

    // 结案
    @Transactional
    public Consultation close(String consultationId, String internalCategory, String complexSubCategory) {
        return consultationRepository.findById(consultationId).map(c -> {
            c.setStatus("closed");
            c.setClosedAt(LocalDateTime.now());
            c.setInternalCategory(internalCategory != null ? internalCategory : c.getInternalCategory());
            c.setComplexSubCategory(complexSubCategory != null ? complexSubCategory : c.getComplexSubCategory());
            c.setUpdatedAt(LocalDateTime.now());
            enrichConsultation(consultationRepository.save(c));

            notificationService.notifyConsultationClosed(c);
            return c;
        }).orElse(null);
    }

    // 转交
    @Transactional
    public Consultation transfer(String consultationId, String newHandlerId, String reason) {
        return consultationRepository.findById(consultationId).map(c -> {
            String oldHandlerId = c.getHandlerId();
            c.setHandlerId(newHandlerId);
            c.setUpdatedAt(LocalDateTime.now());
            enrichConsultation(consultationRepository.save(c));

            notificationService.notifyConsultationTransferred(c, oldHandlerId);
            return c;
        }).orElse(null);
    }

    // 评分
    @Transactional
    public void rate(String consultationId, Integer rating) {
        consultationRepository.findById(consultationId).ifPresent(c -> {
            c.setRating(rating);
            c.setUpdatedAt(LocalDateTime.now());
            consultationRepository.save(c);
        });
    }

    // 搜索
    public Page<Consultation> search(String keyword, Pageable pageable) {
        return consultationRepository.search(keyword, pageable).map(this::enrichConsultation);
    }

    // 相似问题
    public List<Consultation> getSimilar(String id, String title, int limit) {
        return consultationRepository.findSimilar(title, id, PageRequest.of(0, limit));
    }

    // 获取回复列表
    public List<ConsultationReply> getReplies(String consultationId) {
        return replyRepository.findByConsultationIdOrderByCreatedAtAsc(consultationId).stream()
                .peek(r -> userRepository.findById(r.getUserId()).ifPresent(r::setUser))
                .toList();
    }

    // 统计
    public java.util.Map<String, Object> getStats() {
        LocalDate today = LocalDate.now();
        LocalDateTime startOfDay = today.atStartOfDay();
        LocalDateTime endOfDay = today.atTime(LocalTime.MAX);
        LocalDateTime startOfWeek = today.with(DayOfWeek.MONDAY).atStartOfDay();

        java.util.Map<String, Object> stats = new java.util.HashMap<>();
        stats.put("todayNew", consultationRepository.countByStatusAndSubmittedAtBetween("pending", startOfDay, endOfDay) +
                               consultationRepository.countByStatusAndSubmittedAtBetween("accepted", startOfDay, endOfDay) +
                               consultationRepository.countByStatusAndSubmittedAtBetween("in_progress", startOfDay, endOfDay) +
                               consultationRepository.countByStatusAndSubmittedAtBetween("waiting_supplement", startOfDay, endOfDay));
        stats.put("pending", consultationRepository.countByStatus("pending"));
        stats.put("thisWeekClosed", consultationRepository.countClosedBetween(startOfWeek, endOfDay));

        return stats;
    }

    private Consultation enrichConsultation(Consultation c) {
        userRepository.findById(c.getSubmitterId()).ifPresent(c::setSubmitter);
        if (c.getHandlerId() != null) {
            userRepository.findById(c.getHandlerId()).ifPresent(c::setHandler);
        }
        return c;
    }

    // 辅助方法
    public Page<Consultation> findByHandlerIdAndStatus(String handlerId, String status, Pageable pageable) {
        return consultationRepository.findByHandlerIdAndStatus(handlerId, status, pageable);
    }
}
