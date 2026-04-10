package models

import (
	"time"
)

// ============ 合同模板申请 ============

// 申请状态
const (
	TemplateRequestStatusPendingApproval  = "pending_approval"  // 待上级审批
	TemplateRequestStatusSupplemented      = "supplemented"      // 补充资料中
	TemplateRequestStatusRejected          = "rejected"          // 已拒绝
	TemplateRequestStatusPendingDraft      = "pending_draft"     // 待法务拟写
	TemplateRequestStatusPendingReview     = "pending_review"    // 待法务负责人审核
	TemplateRequestStatusPublished         = "published"         // 已发布
	TemplateRequestStatusDisabled          = "disabled"          // 已禁用
)

// 申请类型
const (
	TemplateRequestTypeNew     = "new"
	TemplateRequestTypeUpdate  = "update"
)

// 模板申请
type TemplateRequest struct {
	ID              string      `gorm:"primaryKey;size:36" json:"id"`
	RequestNo       string      `gorm:"uniqueIndex;size:50;not null" json:"request_no"` // TMPL-YYYYMMDD-序号
	RequestType     string      `gorm:"size:20;not null" json:"request_type"`           // new/update

	// 关联的现有模板（更新时）
	ExistingTemplateID *string `gorm:"size:36" json:"existing_template_id"`

	ContractType    string `gorm:"size:50;not null" json:"contract_type"`
	Title           string `gorm:"size:100;not null" json:"title"`

	// 需求描述
	BusinessScenario string `gorm:"type:text" json:"business_scenario"` // 核心业务场景
	BusinessFlow     string `gorm:"type:text" json:"business_flow"`     // 业务关键流程
	KeyClauses       string `gorm:"type:text" json:"key_clauses"`       // 核心关注条款
	DiffFromExisting string `gorm:"type:text" json:"diff_from_existing"` // 与现有模板区别

	// 参考资料
	ReferenceFiles JSONType `gorm:"type:jsonb" json:"reference_files"` // [{filename, filepath}]

	// 期望完成时间
	ExpectedDate *time.Time `json:"expected_date"`

	// 当前状态
	Status string `gorm:"size:30;not null;default:pending_approval" json:"status"`

	// 当前审批节点
	CurrentStep int `gorm:"default:1" json:"current_step"` // 1:L1审批, 2:法务拟写, 3:法务负责人审核

	// 提交人
	SubmitterID string `gorm:"size:36;not null" json:"submitter_id"`

	// L1审批人（业务主管）
	L1ApproverID *string    `gorm:"size:36" json:"l1_approver_id"`
	L1ApprovedAt *time.Time `json:"l1_approved_at"`
	L1Comment    string     `gorm:"type:text" json:"l1_comment"`

	// 法务拟写人
	DrafterID *string    `gorm:"size:36" json:"drafter_id"`
	DraftedAt *time.Time `json:"drafted_at"`

	// 法务负责人审核
	ReviewerID    *string    `gorm:"size:36" json:"reviewer_id"`
	ReviewedAt    *time.Time `json:"reviewed_at"`
	ReviewComment string     `gorm:"type:text" json:"review_comment"`

	// 创建时间
	SubmittedAt time.Time `json:"submitted_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Submitter        *User                   `gorm:"foreignKey:SubmitterID" json:"submitter,omitempty"`
	L1Approver       *User                   `gorm:"foreignKey:L1ApproverID" json:"l1_approver,omitempty"`
	Drafter          *User                   `gorm:"foreignKey:DrafterID" json:"drafter,omitempty"`
	Reviewer         *User                   `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	ExistingTemplate *Template               `gorm:"foreignKey:ExistingTemplateID" json:"existing_template,omitempty"`
	Attachments      []TemplateRequestAttachment `gorm:"foreignKey:TemplateRequestID" json:"attachments,omitempty"`
	ApprovalLogs     []TemplateApprovalLog    `gorm:"foreignKey:TemplateRequestID" json:"approval_logs,omitempty"`
}

func (TemplateRequest) TableName() string {
	return "template_requests"
}

// 模板申请附件
type TemplateRequestAttachment struct {
	ID               string `gorm:"primaryKey;size:36" json:"id"`
	TemplateRequestID string `gorm:"size:36;not null;index" json:"template_request_id"`
	AttachmentID     string `gorm:"size:36;not null" json:"attachment_id"`

	Attachment *Attachment `gorm:"foreignKey:AttachmentID" json:"attachment,omitempty"`
}

func (TemplateRequestAttachment) TableName() string {
	return "template_request_attachments"
}

// 模板申请审批日志
type TemplateApprovalLog struct {
	ID               string    `gorm:"primaryKey;size:36" json:"id"`
	TemplateRequestID string    `gorm:"size:36;not null;index" json:"template_request_id"`
	ApproverID       string    `gorm:"size:36;not null" json:"approver_id"`
	Action           string    `gorm:"size:50;not null" json:"action"` // approve/reject/return_for_supplement
	Comment          string    `gorm:"type:text" json:"comment"`
	CreatedAt        time.Time `json:"created_at"`

	Approver *User `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}

func (TemplateApprovalLog) TableName() string {
	return "template_approval_logs"
}

// ============ 合同模板 ============

// 模板状态
const (
	TemplateStatusDraft     = "draft"
	TemplateStatusPublished = "published"
	TemplateStatusDisabled   = "disabled"
)

type Template struct {
	ID              string    `gorm:"primaryKey;size:36" json:"id"`
	Name            string    `gorm:"size:100;not null" json:"name"`
	ContractType    string    `gorm:"size:50;not null" json:"contract_type"`
	Version         string    `gorm:"size:20;not null" json:"version"` // v1.0, v1.1, v2.0
	Description     string    `gorm:"type:text" json:"description"`    // 模板说明
	FilePath        string    `gorm:"size:500;not null" json:"file_path"`
	EditableClauses string    `gorm:"type:text" json:"editable_clauses"` // 可编辑条款

	Status string `gorm:"size:20;not null;default:draft" json:"status"` // draft/published/disabled

	// 关联的申请
	TemplateRequestID *string  `gorm:"size:36" json:"template_request_id"`
	PublishedByID     *string  `gorm:"size:36" json:"published_by_id"`
	PublishedAt       *time.Time `json:"published_at"`

	DownloadCount int       `gorm:"default:0" json:"download_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	PublishedBy *User      `gorm:"foreignKey:PublishedByID" json:"published_by,omitempty"`
}

func (Template) TableName() string {
	return "templates"
}

type TemplateVersion struct {
	ID            string    `gorm:"primaryKey;size:36" json:"id"`
	TemplateID    string    `gorm:"size:36;not null" json:"template_id"`
	Version       string    `gorm:"size:20;not null" json:"version"`
	FilePath      string    `gorm:"size:500;not null" json:"file_path"`
	ChangeNotes   string    `gorm:"type:text" json:"change_notes"`
	PublishedByID string    `gorm:"size:36" json:"published_by_id"`
	PublishedAt   time.Time `json:"published_at"`
	CreatedAt     time.Time `json:"created_at"`
}

func (TemplateVersion) TableName() string {
	return "template_versions"
}
