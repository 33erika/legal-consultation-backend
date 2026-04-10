package models

import (
	"time"
)

// 附件关联类型
const (
	AttachmentEntityConsultation     = "consultation"
	AttachmentEntityReply            = "reply"
	AttachmentEntityTemplate         = "template"
	AttachmentEntityTemplateRequest  = "template_request"
)

type Attachment struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	EntityType  string    `gorm:"size:50;not null" json:"entity_type"` // consultation/reply/template
	EntityID    string    `gorm:"size:36;not null;index" json:"entity_id"`
	FileName    string    `gorm:"size:255;not null" json:"file_name"`
	FilePath    string    `gorm:"size:500;not null" json:"file_path"`
	FileSize    int64     `gorm:"not null" json:"file_size"`
	ContentType string    `gorm:"size:100" json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Attachment) TableName() string {
	return "attachments"
}

// 咨询附件关联表
type ConsultationAttachment struct {
	ID             string `gorm:"primaryKey;size:36" json:"id"`
	ConsultationID string `gorm:"size:36;not null;index" json:"consultation_id"`
	AttachmentID   string `gorm:"size:36;not null" json:"attachment_id"`

	Attachment *Attachment `gorm:"foreignKey:AttachmentID" json:"attachment,omitempty"`
}

func (ConsultationAttachment) TableName() string {
	return "consultation_attachments"
}
