package com.legalsys.repository;

import com.legalsys.entity.Attachment;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.List;

@Repository
public interface AttachmentRepository extends JpaRepository<Attachment, String> {

    List<Attachment> findByEntityTypeAndEntityId(String entityType, String entityId);

    List<Attachment> findByUploaderId(String uploaderId);
}
