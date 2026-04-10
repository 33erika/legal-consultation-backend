package models

import (
	"time"
)

// 通知配置
type NotificationConfig struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	ConfigKey   string    `gorm:"uniqueIndex;size:100;not null" json:"config_key"`
	ConfigValue string    `gorm:"type:text" json:"config_value"`
	Description string    `gorm:"size:255" json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (NotificationConfig) TableName() string {
	return "notification_configs"
}

// 操作日志
type OperationLog struct {
	ID         string    `gorm:"primaryKey;size:36" json:"id"`
	EntityType string    `gorm:"size:50;not null" json:"entity_type"` // consultation/template_request
	EntityID   string    `gorm:"size:36;not null;index" json:"entity_id"`
	Action     string    `gorm:"size:50;not null" json:"action"` // create/update/delete/status_change
	OperatorID string    `gorm:"size:36;not null" json:"operator_id"`
	OldValue   JSONType  `gorm:"type:jsonb" json:"old_value"`
	NewValue   JSONType  `gorm:"type:jsonb" json:"new_value"`
	IPAddress  string    `gorm:"size:50" json:"ip_address"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`

	Operator *User `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

// 案例收藏
type CaseCollection struct {
	ID             string    `gorm:"primaryKey;size:36" json:"id"`
	ConsultationID string    `gorm:"size:36;not null;uniqueIndex" json:"consultation_id"`
	CollectorID    string    `gorm:"size:36;not null" json:"collector_id"`
	Tags           JSONType  `gorm:"type:jsonb" json:"tags"` // ["劳动纠纷", "合同审核"]
	CreatedAt      time.Time `json:"created_at"`

	Consultation *Consultation `gorm:"foreignKey:ConsultationID" json:"consultation,omitempty"`
	Collector    *User         `gorm:"foreignKey:CollectorID" json:"collector,omitempty"`
}

func (CaseCollection) TableName() string {
	return "case_collections"
}

// 咨询类型配置（引导配置）
type ConsultationTypeConfig struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	Type      string    `gorm:"uniqueIndex;size:50;not null" json:"type"` // complaint/contract/labor/ip/dispute
	Name      string    `gorm:"size:100;not null" json:"name"`
	Keywords  JSONType  `gorm:"type:jsonb" json:"keywords"`  // 触发关键词列表
	Fields    JSONType  `gorm:"type:jsonb" json:"fields"`    // 引导字段配置
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ConsultationTypeConfig) TableName() string {
	return "consultation_type_configs"
}
