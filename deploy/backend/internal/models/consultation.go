package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// ============ 常量定义 ============

// 咨询状态
const (
	ConsultationStatusPending    = "pending"    // 待处理
	ConsultationStatusProcessing  = "processing" // 处理中
	ConsultationStatusReplied     = "replied"    // 已回复
	ConsultationStatusClosed      = "closed"     // 已结案
)

// 紧急程度
const (
	UrgencyNormal     = "normal"      // 一般
	UrgencyUrgent     = "urgent"      // 紧急
	UrgencyVeryUrgent = "very_urgent" // 非常紧急
)

// 咨询类型（员工选择）
const (
	ConsultationTypeComplaint = "complaint" // 客诉
	ConsultationTypeContract  = "contract"  // 合同
	ConsultationTypeLabor     = "labor"      // 劳动
	ConsultationTypeIP        = "ip"         // 知识产权
	ConsultationTypeDispute   = "dispute"    // 纠纷咨询
	ConsultationTypeOther     = "other"      // 其他
)

// 内部分类（法务内部使用）
const (
	InternalCategorySimple  = "simple"  // 简单问题
	InternalCategoryComplex = "complex"  // 复杂问题
)

// 复杂问题子分类
const (
	ComplexSubCategoryDispute   = "dispute"   // 纠纷咨询
	ComplexSubCategoryContract  = "contract"  // 合同
	ComplexSubCategoryLabor     = "labor"     // 劳动
	ComplexSubCategoryComplaint = "complaint" // 客诉
	ComplexSubCategoryIP        = "ip"         // 知识产权
	ComplexSubCategoryOther     = "other"     // 其他
)

// ============ JSON 类型 ============

type JSONType map[string]interface{}

func (j JSONType) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONType) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// ============ 咨询主体 ============

type Consultation struct {
	ID                 string      `gorm:"primaryKey;size:36" json:"id"`
	TicketNo           string      `gorm:"uniqueIndex;size:50;not null" json:"ticket_no"` // CONS-YYYYMMDD-序号
	Title              string      `gorm:"size:100;not null" json:"title"`
	Description        string      `gorm:"type:text;not null" json:"description"`
	Urgency            string      `gorm:"size:20;not null" json:"urgency"`
	Status             string      `gorm:"size:20;not null;default:pending" json:"status"`

	// 员工选择类型
	ConsultationType string `gorm:"size:50" json:"consultation_type"`

	// 员工填写的扩展信息（JSON存储）
	ExtensionData JSONType `gorm:"type:jsonb" json:"extension_data"`

	// 法务内部分类
	InternalCategory   string `gorm:"size:20" json:"internal_category"`    // simple/complex
	ComplexSubCategory string `gorm:"size:50" json:"complex_sub_category"`  // 复杂问题子分类

	// 关联信息
	SubmitterID string  `gorm:"size:36;not null" json:"submitter_id"`
	HandlerID   *string  `gorm:"size:36" json:"handler_id"`

	// 评价
	Rating *int `json:"rating"` // 1-5星

	// 时间戳
	SubmittedAt    time.Time  `json:"submitted_at"`
	AcceptedAt     *time.Time `json:"accepted_at"`
	FirstRepliedAt *time.Time `json:"first_replied_at"`
	ClosedAt       *time.Time `json:"closed_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// 关联
	Submitter     *User                  `gorm:"foreignKey:SubmitterID" json:"submitter,omitempty"`
	Handler       *User                  `gorm:"foreignKey:HandlerID" json:"handler,omitempty"`
	Replies       []ConsultationReply    `gorm:"foreignKey:ConsultationID" json:"replies,omitempty"`
	Attachments   []ConsultationAttachment `gorm:"foreignKey:ConsultationID" json:"attachments,omitempty"`
	OperationLogs []OperationLog         `gorm:"foreignKey:EntityID;constraint:OnDelete:CASCADE" json:"operation_logs,omitempty"`
}

func (Consultation) TableName() string {
	return "consultations"
}

// ============ 回复记录 ============

// 回复类型
const (
	ReplyTypeStaff   = "staff"   // 员工补充
	ReplyTypeLegal   = "legal"   // 法务回复
	ReplyTypeSystem  = "system"  // 系统消息
)

type ConsultationReply struct {
	ID              string      `gorm:"primaryKey;size:36" json:"id"`
	ConsultationID  string      `gorm:"size:36;not null;index" json:"consultation_id"`
	ReplyType       string      `gorm:"size:20;not null" json:"reply_type"` // staff/legal/system
	Content         string      `gorm:"type:text" json:"content"`
	AuthorID        string      `gorm:"size:36" json:"author_id"`
	CreatedAt       time.Time   `json:"created_at"`

	Author      *User        `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Attachments []Attachment `gorm:"foreignKey:EntityID;constraint:OnDelete:CASCADE" json:"attachments,omitempty"`
}

func (ConsultationReply) TableName() string {
	return "consultation_replies"
}
