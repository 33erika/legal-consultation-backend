---
AIGC:
    ContentProducer: Minimax Agent AI
    ContentPropagator: Minimax Agent AI
    Label: AIGC
    ProduceID: de8a816ec119fb45e78827e2b88ffd4c
    PropagateID: de8a816ec119fb45e78827e2b88ffd4c
    ReservedCode1: 3045022100a8c42a9b7ee0083b4c49f6065d005af72fd12436807f16645a97499a2ca4a5ee02202151a83b06d26c85777c929c637e96ca2de860145a96941d4294b91e23758ddf
    ReservedCode2: 30460221008f78504b8aef525abd1e10d29e40fe34b758914cde01e62abd528d4afecb2881022100feac1e5e1e6c507af54843cf8278b04971cffc22dcaccee19fe5362f2bd8b423
---

# Design: 法务部法律咨询系统

## 需求映射

| Story | 实现方式 |
|-------|---------|
| Story 1: 用户认证与权限管理 | API: 认证相关接口, 组件: LoginPage |
| Story 2: 法律咨询 - 发起咨询（员工） | API: /api/v1/consultations, 组件: ConsultationForm, ConsultationList, ConsultationDetail |
| Story 3: 法律咨询 - 处理咨询（法务专员） | API: /api/v1/consultations/{id}/accept 等, 组件: LegalWorkbench, ConsultationPool |
| Story 4-5: 合同模板管理 - 发起申请/审批 | API: /api/v1/template-requests, 组件: TemplateRequestForm, ApprovalPages |
| Story 6-7: 合同模板管理 - 拟写/审核 | API: /api/v1/templates, 组件: TemplateEditor, TemplateAudit |
| Story 8: 合同模板库管理 | API: /api/v1/templates, /api/v1/template-versions, 组件: TemplateLibrary |
| Story 9: 钉钉通知集成 | Service: DingTalkNotifier |
| Story 10: 数据统计与报表 | API: /api/v1/statistics, 组件: StatisticsDashboard |
| Story 11: 搜索与历史案例复用 | API: /api/v1/consultations/search, 组件: SearchPage, CaseLibrary |
| Story 12: 系统管理 | API: /api/v1/admin/*, 组件: AdminPages |

---

## 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (Next.js)                     │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │  Auth   │ │Consult- │ │Template │ │  Admin  │           │
│  │  Pages  │ │  ation  │ │  Pages  │ │  Pages  │           │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
└───────────────────────────┬─────────────────────────────────┘
                            │ REST API
┌───────────────────────────▼─────────────────────────────────┐
│                      Backend (Go)                            │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │   Auth  │ │Consult- │ │Template │ │  Admin  │           │
│  │ Handler │ │ation API│ │   API   │ │  API    │           │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘           │
│       │           │           │           │                 │
│  ┌────▼───────────▼───────────▼───────────▼────┐           │
│  │              Service Layer                    │           │
│  │  ConsultationService | TemplateService | ...   │           │
│  └─────────────────────┬─────────────────────────┘           │
│                        │                                     │
│  ┌─────────────────────▼─────────────────────────┐          │
│  │            Repository Layer (GORM)              │          │
│  │  ConsultationRepo | TemplateRepo | UserRepo     │          │
│  └─────────────────────┬─────────────────────────┘          │
│                        │                                     │
│  ┌─────────────────────▼─────────────────────────┐          │
│  │           Database (PostgreSQL)                │          │
│  │  users | consultations | consultation_replies  │          │
│  │  consultation_attachments | template_requests  │          │
│  │  templates | template_versions | departments   │          │
│  │  notification_configs | operation_logs         │          │
│  └───────────────────────────────────────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

---

## 数据模型

### 1. 用户与认证

```go
// backend/internal/models/user.go

package models

import (
    "time"
)

type User struct {
    ID           string      `gorm:"primaryKey;size:36" json:"id"`
    EmployeeID   string      `gorm:"uniqueIndex;size:32;not null" json:"employee_id"`  // 工号
    Password     string      `gorm:"size:255;not null" json:"-"`                       // bcrypt加密
    Name         string      `gorm:"size:100;not null" json:"name"`
    Email        string      `gorm:"size:255" json:"email"`
    Role         string      `gorm:"size:50;not null" json:"role"`                     // employee/supervisor/legal_staff/legal_head/admin
    DepartmentID *string      `gorm:"size:36" json:"department_id"`
    Status       string      `gorm:"size:20;default:active" json:"status"`             // active/inactive
    CreatedAt    time.Time   `json:"created_at"`
    UpdatedAt    time.Time   `json:"updated_at"`

    Department   *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

type Department struct {
    ID           string       `gorm:"primaryKey;size:36" json:"id"`
    Name         string       `gorm:"size:100;not null" json:"name"`
    ParentID     *string      `gorm:"size:36" json:"parent_id"`
    CreatedAt    time.Time    `json:"created_at"`
    UpdatedAt    time.Time    `json:"updated_at"`

    Parent       *Department  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Children     []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
    Users        []User       `gorm:"foreignKey:DepartmentID" json:"users,omitempty"`
}
```

### 2. 法律咨询

```go
// backend/internal/models/consultation.go

package models

import (
    "time"
)

// 咨询状态
const (
    ConsultationStatusPending    = "pending"    // 待处理
    ConsultationStatusProcessing = "processing" // 处理中
    ConsultationStatusReplied    = "replied"    // 已回复
    ConsultationStatusClosed     = "closed"    // 已结案
)

// 紧急程度
const (
    UrgencyNormal      = "normal"       // 一般
    UrgencyUrgent      = "urgent"       // 紧急
    UrgencyVeryUrgent  = "very_urgent"  // 非常紧急
)

// 咨询类型（员工选择）
const (
    ConsultationTypeComplaint    = "complaint"    // 客诉
    ConsultationTypeContract     = "contract"     // 合同
    ConsultationTypeLabor        = "labor"         // 劳动
    ConsultationTypeIP           = "ip"            // 知识产权
    ConsultationTypeDispute       = "dispute"       // 纠纷咨询
    ConsultationTypeOther        = "other"         // 其他
)

// 内部分类（法务内部使用）
const (
    InternalCategorySimple       = "simple"        // 简单问题
    InternalCategoryComplex      = "complex"        // 复杂问题
)

// 复杂问题子分类
const (
    ComplexSubCategoryDispute    = "dispute"       // 纠纷咨询
    ComplexSubCategoryContract   = "contract"      // 合同
    ComplexSubCategoryLabor      = "labor"         // 劳动
    ComplexSubCategoryComplaint   = "complaint"    // 客诉
    ComplexSubCategoryIP          = "ip"            // 知识产权
    ComplexSubCategoryOther      = "other"         // 其他
)

type Consultation struct {
    ID                  string      `gorm:"primaryKey;size:36" json:"id"`
    TicketNo            string      `gorm:"uniqueIndex;size:50;not null" json:"ticket_no"`  // CONS-YYYYMMDD-序号
    Title               string      `gorm:"size:100;not null" json:"title"`
    Description         string      `gorm:"type:text;not null" json:"description"`
    Urgency             string      `gorm:"size:20;not null" json:"urgency"`
    Status              string      `gorm:"size:20;not null;default:pending" json:"status"`

    // 员工选择类型
    ConsultationType   string      `gorm:"size:50" json:"consultation_type"`

    // 员工填写的扩展信息（JSON存储）
    ExtensionData       JSONType    `gorm:"type:jsonb" json:"extension_data"`

    // 法务内部分类
    InternalCategory    string      `gorm:"size:20" json:"internal_category"`    // simple/complex
    ComplexSubCategory  string      `gorm:"size:50" json:"complex_sub_category"` // 复杂问题子分类

    // 关联信息
    SubmitterID         string      `gorm:"size:36;not null" json:"submitter_id"`
    HandlerID           *string     `gorm:"size:36" json:"handler_id"`

    // 评价
    Rating              *int        `json:"rating"`   // 1-5星

    // 时间戳
    SubmittedAt         time.Time   `json:"submitted_at"`
    AcceptedAt          *time.Time  `json:"accepted_at"`
    FirstRepliedAt      *time.Time  `json:"first_replied_at"`
    ClosedAt            *time.Time  `json:"closed_at"`
    CreatedAt           time.Time   `json:"created_at"`
    UpdatedAt           time.Time   `json:"updated_at"`

    // 关联
    Submitter           *User       `gorm:"foreignKey:SubmitterID" json:"submitter,omitempty"`
    Handler             *User       `gorm:"foreignKey:HandlerID" json:"handler,omitempty"`
    Replies             []ConsultationReply `gorm:"foreignKey:ConsultationID" json:"replies,omitempty"`
    Attachments         []ConsultationAttachment `gorm:"foreignKey:ConsultationID" json:"attachments,omitempty"`
    OperationLogs       []OperationLog `gorm:"foreignKey:EntityID;constraint:OnDelete:CASCADE" json:"operation_logs,omitempty"`
}

// JSONType 自定义类型用于JSON存储
type JSONType map[string]interface{}
```

### 3. 咨询回复

```go
// backend/internal/models/consultation_reply.go

package models

// 回复类型
const (
    ReplyTypeStaff    = "staff"    // 员工补充
    ReplyTypeLegal    = "legal"    // 法务回复
    ReplyTypeSystem    = "system"   // 系统消息（如要求补充资料）
)

type ConsultationReply struct {
    ID             string      `gorm:"primaryKey;size:36" json:"id"`
    ConsultationID string      `gorm:"size:36;not null;index" json:"consultation_id"`
    ReplyType      string      `gorm:"size:20;not null" json:"reply_type"`
    Content        string      `gorm:"type:text" json:"content"`
    AuthorID       string      `gorm:"size:36" json:"author_id"`
    CreatedAt      time.Time   `json:"created_at"`

    Author         *User       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
    Attachments     []Attachment `gorm:"foreignKey:EntityID;constraint:OnDelete:CASCADE" json:"attachments,omitempty"`
}
```

### 4. 附件

```go
// backend/internal/models/attachment.go

package models

// 附件关联类型
const (
    AttachmentEntityConsultation = "consultation"
    AttachmentEntityReply        = "reply"
    AttachmentEntityTemplate     = "template"
    AttachmentEntityTemplateRequest = "template_request"
)

type Attachment struct {
    ID           string    `gorm:"primaryKey;size:36" json:"id"`
    EntityType   string    `gorm:"size:50;not null" json:"entity_type"`
    EntityID     string    `gorm:"size:36;not null;index" json:"entity_id"`
    FileName     string    `gorm:"size:255;not null" json:"file_name"`
    FilePath     string    `gorm:"size:500;not null" json:"file_path"`
    FileSize     int64     `gorm:"not null" json:"file_size"`
    ContentType  string    `gorm:"size:100" json:"content_type"`
    CreatedAt    time.Time `json:"created_at"`
}

type ConsultationAttachment struct {
    ID             string    `gorm:"primaryKey;size:36" json:"id"`
    ConsultationID string    `gorm:"size:36;not null;index" json:"consultation_id"`
    AttachmentID   string    `gorm:"size:36;not null" json:"attachment_id"`

    Attachment     *Attachment `gorm:"foreignKey:AttachmentID" json:"attachment,omitempty"`
}
```

### 5. 合同模板申请

```go
// backend/internal/models/template_request.go

package models

// 申请状态
const (
    TemplateRequestStatusPendingApproval = "pending_approval"  // 待上级审批
    TemplateRequestStatusSupplemented    = "supplemented"      // 补充资料中
    TemplateRequestStatusRejected        = "rejected"          // 已拒绝
    TemplateRequestStatusPendingDraft    = "pending_draft"     // 待法务拟写
    TemplateRequestStatusPendingReview   = "pending_review"    // 待法务负责人审核
    TemplateRequestStatusPublished       = "published"        // 已发布
    TemplateRequestStatusDisabled        = "disabled"          // 已禁用
)

// 申请类型
const (
    TemplateRequestTypeNew     = "new"
    TemplateRequestTypeUpdate = "update"
)

type TemplateRequest struct {
    ID              string      `gorm:"primaryKey;size:36" json:"id"`
    RequestNo       string      `gorm:"uniqueIndex;size:50;not null" json:"request_no"`  // TMPL-YYYYMMDD-序号
    RequestType     string      `gorm:"size:20;not null" json:"request_type"`  // new/update

    // 关联的现有模板（更新时）
    ExistingTemplateID *string  `gorm:"size:36" json:"existing_template_id"`

    ContractType     string    `gorm:"size:50;not null" json:"contract_type"`
    Title             string    `gorm:"size:100;not null" json:"title"`

    // 需求描述
    BusinessScenario  string    `gorm:"type:text" json:"business_scenario"`  // 核心业务场景
    BusinessFlow      string    `gorm:"type:text" json:"business_flow"`      // 业务关键流程
    KeyClauses        string    `gorm:"type:text" json:"key_clauses"`         // 核心关注条款
    DiffFromExisting  string    `gorm:"type:text" json:"diff_from_existing"`  // 与现有模板区别

    // 参考资料
    ReferenceFiles    JSONType  `gorm:"type:jsonb" json:"reference_files"`  // [{filename, filepath}]

    // 期望完成时间
    ExpectedDate      *time.Time `json:"expected_date"`

    // 当前状态
    Status            string    `gorm:"size:30;not null;default:pending_approval" json:"status"`

    // 当前审批节点
    CurrentStep       int       `gorm:"default:1" json:"current_step"`  // 1:L1审批, 2:法务拟写, 3:法务负责人审核

    // 提交人
    SubmitterID       string    `gorm:"size:36;not null" json:"submitter_id"`

    // L1审批人（业务主管）
    L1ApproverID      *string   `gorm:"size:36" json:"l1_approver_id"`
    L1ApprovedAt      *time.Time `json:"l1_approved_at"`
    L1Comment         string    `gorm:"type:text" json:"l1_comment"`

    // 法务拟写人
    DrafterID         *string   `gorm:"size:36" json:"drafter_id"`
    DraftedAt         *time.Time `json:"drafted_at"`

    // 法务负责人审核
    ReviewerID        *string   `gorm:"size:36" json:"reviewer_id"`
    ReviewedAt        *time.Time `json:"reviewed_at"`
    ReviewComment     string    `gorm:"type:text" json:"review_comment"`

    // 创建时间
    SubmittedAt       time.Time `json:"submitted_at"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`

    // 关联
    Submitter         *User      `gorm:"foreignKey:SubmitterID" json:"submitter,omitempty"`
    L1Approver        *User      `gorm:"foreignKey:L1ApproverID" json:"l1_approver,omitempty"`
    Drafter           *User      `gorm:"foreignKey:DrafterID" json:"drafter,omitempty"`
    Reviewer          *User      `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
    ExistingTemplate  *Template  `gorm:"foreignKey:ExistingTemplateID" json:"existing_template,omitempty"`
    Attachments       []TemplateRequestAttachment `gorm:"foreignKey:TemplateRequestID" json:"attachments,omitempty"`
    ApprovalLogs      []TemplateApprovalLog       `gorm:"foreignKey:TemplateRequestID" json:"approval_logs,omitempty"`
}
```

### 6. 合同模板

```go
// backend/internal/models/template.go

package models

type Template struct {
    ID              string      `gorm:"primaryKey;size:36" json:"id"`
    Name            string      `gorm:"size:100;not null" json:"name"`
    ContractType    string      `gorm:"size:50;not null" json:"contract_type"`
    Version         string      `gorm:"size:20;not null" json:"version"`  // v1.0, v1.1, v2.0
    Description     string      `gorm:"type:text" json:"description"`  // 模板说明
    FilePath        string      `gorm:"size:500;not null" json:"file_path"`
    EditableClauses string      `gorm:"type:text" json:"editable_clauses"`  // 可编辑条款

    Status          string      `gorm:"size:20;not null;default:draft" json:"status"`  // draft/published/disabled

    // 关联的申请
    TemplateRequestID *string   `gorm:"size:36" json:"template_request_id"`
    PublishedByID    *string    `gorm:"size:36" json:"published_by_id"`
    PublishedAt      *time.Time `json:"published_at"`

    DownloadCount    int        `gorm:"default:0" json:"download_count"`

    CreatedAt        time.Time  `json:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"`

    PublishedBy      *User      `gorm:"foreignKey:PublishedByID" json:"published_by,omitempty"`
    Versions         []Template `gorm:"foreignKey:Name;references:Name" json:"versions,omitempty"`
}

type TemplateVersion struct {
    ID           string      `gorm:"primaryKey;size:36" json:"id"`
    TemplateID   string      `gorm:"size:36;not null" json:"template_id"`
    Version      string      `gorm:"size:20;not null" json:"version"`
    FilePath     string      `gorm:"size:500;not null" json:"file_path"`
    ChangeNotes  string      `gorm:"type:text" json:"change_notes"`
    PublishedByID string     `gorm:"size:36" json:"published_by_id"`
    PublishedAt  time.Time   `json:"published_at"`
    CreatedAt    time.Time   `json:"created_at"`
}
```

### 7. 系统配置与日志

```go
// backend/internal/models/config.go

package models

type NotificationConfig struct {
    ID           string    `gorm:"primaryKey;size:36" json:"id"`
    ConfigKey    string    `gorm:"uniqueIndex;size:100;not null" json:"config_key"`
    ConfigValue  string    `gorm:"type:text" json:"config_value"`
    Description  string    `gorm:"size:255" json:"description"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type OperationLog struct {
    ID           string    `gorm:"primaryKey;size:36" json:"id"`
    EntityType   string    `gorm:"size:50;not null" json:"entity_type"`  // consultation/template_request
    EntityID     string    `gorm:"size:36;not null;index" json:"entity_id"`
    Action       string    `gorm:"size:50;not null" json:"action"`  // create/update/delete/status_change
    OperatorID   string    `gorm:"size:36;not null" json:"operator_id"`
    OldValue     JSONType  `gorm:"type:jsonb" json:"old_value"`
    NewValue     JSONType  `gorm:"type:jsonb" json:"new_value"`
    IPAddress    string    `gorm:"size:50" json:"ip_address"`
    UserAgent    string    `gorm:"size:500" json:"user_agent"`
    CreatedAt    time.Time `json:"created_at"`

    Operator     *User     `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
}

type CaseCollection struct {
    ID             string    `gorm:"primaryKey;size:36" json:"id"`
    ConsultationID string    `gorm:"size:36;not null;uniqueIndex" json:"consultation_id"`
    CollectorID    string    `gorm:"size:36;not null" json:"collector_id"`
    Tags           JSONType  `gorm:"type:jsonb" json:"tags"`  // ["劳动纠纷", "合同审核"]
    CreatedAt      time.Time `json:"created_at"`

    Consultation  *Consultation `gorm:"foreignKey:ConsultationID" json:"consultation,omitempty"`
    Collector     *User         `gorm:"foreignKey:CollectorID" json:"collector,omitempty"`
}
```

### 8. 咨询类型配置（引导配置）

```go
// backend/internal/models/consultation_config.go

package models

type ConsultationTypeConfig struct {
    ID          string    `gorm:"primaryKey;size:36" json:"id"`
    Type        string    `gorm:"uniqueIndex;size:50;not null" json:"type"`  // complaint/contract/labor/ip/dispute
    Name        string    `gorm:"size:100;not null" json:"name"`
    Keywords    JSONType  `gorm:"type:jsonb" json:"keywords"`  // 触发关键词列表
    Fields      JSONType  `gorm:"type:jsonb" json:"fields"`  // 引导字段配置
    SortOrder   int       `gorm:"default:0" json:"sort_order"`
    Enabled     bool      `gorm:"default:true" json:"enabled"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Fields JSON结构示例
// {
//   "fields": [
//     {
//       "name": "客诉来源",
//       "type": "select",
//       "options": ["12315平台", "第三方投诉平台", "内部客服", "举报", "其他"],
//       "required": false
//     },
//     {
//       "name": "消费者诉求",
//       "type": "multiselect",
//       "options": ["仅退款", "退一赔三", "道歉", "要求下架产品", "经济补偿", "其他"],
//       "required": false
//     }
//   ]
// }
```

---

## API 定义

### 1. 认证相关

#### POST /api/v1/auth/login

用户登录

**Request:**
```typescript
interface LoginRequest {
    employee_id: string   // 工号
    password: string      // 密码
}
```

**Response:**
```typescript
interface LoginResponse {
    success: boolean
    data?: {
        token: string
        user: {
            id: string
            employee_id: string
            name: string
            role: string
            department: { id: string, name: string } | null
        }
    }
    error?: { code: string, message: string }
}
```

#### POST /api/v1/auth/logout

用户登出

**Response:**
```typescript
interface LogoutResponse {
    success: boolean
}
```

---

### 2. 法律咨询

#### POST /api/v1/consultations

创建咨询（员工）

**Request:**
```typescript
interface CreateConsultationRequest {
    title: string                    // 问题标题（≤100字）
    description: string              // 问题描述（≤5000字）
    urgency: "normal" | "urgent" | "very_urgent"
    consultation_type?: "complaint" | "contract" | "labor" | "ip" | "dispute" | "other"
    extension_data?: {
        // 客诉类
        complaint_source?: string           // 客诉来源
        complaint_ticket_no?: string       // 投诉编号
        consumer_claims?: string[]         // 消费者诉求
        is_professional_claimer?: boolean  // 是否职业打假人
        current_claim?: string             // 当前诉求
        // 合同类
        contract_type?: string
        contract_party?: string
        contract_amount?: number
        contract_signed?: boolean
        main_disputes?: string
        // 劳动类
        employee_type?: string
        issue_type?: string
        people_count?: number
        current_claim?: string
        // 知识产权类
        ip_type?: string
        is_registered?: boolean
        registration_no?: string
        current_claim?: string
        // 纠纷类
        current_claim?: string
        opposing_party?: string
        dispute_amount?: number
        case_summary?: string
    }
    attachments?: File[]             // 附件（FormData）
}
```

**Response:**
```typescript
interface CreateConsultationResponse {
    success: boolean
    data?: {
        id: string
        ticket_no: string  // CONS-20260410-001
    }
}
```

#### GET /api/v1/consultations

获取咨询列表（员工查看自己的/法务查看待处理的）

**Query Parameters:**
```typescript
interface ListConsultationsQuery {
    role: "staff" | "legal_staff"
    status?: string
    start_date?: string
    end_date?: string
    keyword?: string
    page?: number
    page_size?: number
}
```

**Response:**
```typescript
interface ListConsultationsResponse {
    success: boolean
    data?: {
        items: Array<{
            id: string
            ticket_no: string
            title: string
            urgency: string
            status: string
            consultation_type: string
            submitted_at: string
            handler?: { id: string, name: string }
        }>
        total: number
        page: number
        page_size: number
    }
}
```

#### GET /api/v1/consultations/:id

获取咨询详情

**Response:**
```typescript
interface ConsultationDetailResponse {
    success: boolean
    data?: {
        id: string
        ticket_no: string
        title: string
        description: string
        urgency: string
        status: string
        consultation_type: string
        extension_data: object
        internal_category?: string
        complex_sub_category?: string
        submitter: { id: string, name: string, department: string }
        handler?: { id: string, name: string }
        rating?: number
        submitted_at: string
        accepted_at?: string
        first_replied_at?: string
        closed_at?: string
        attachments: Array<{
            id: string
            filename: string
            file_size: number
            created_at: string
        }>
        replies: Array<{
            id: string
            reply_type: "staff" | "legal" | "system"
            content: string
            author?: { id: string, name: string }
            attachments: Array<{...}>
            created_at: string
        }>
        operation_logs: Array<{
            id: string
            action: string
            operator: { id: string, name: string }
            created_at: string
            details?: object
        }>
    }
}
```

#### POST /api/v1/consultations/:id/accept

接单（法务）

**Request:**
```typescript
interface AcceptConsultationRequest {
    internal_category: "simple" | "complex"
    complex_sub_category?: "dispute" | "contract" | "labor" | "complaint" | "ip" | "other"
}
```

**Response:**
```typescript
interface AcceptConsultationResponse {
    success: boolean
}
```

#### POST /api/v1/consultations/:id/reply

回复咨询（法务）

**Request:**
```typescript
interface ReplyConsultationRequest {
    content: string
    attachments?: File[]
}
```

**Response:**
```typescript
interface ReplyConsultationResponse {
    success: boolean
}
```

#### POST /api/v1/consultations/:id/request-supplement

要求补充资料

**Request:**
```typescript
interface RequestSupplementRequest {
    message: string  // 说明需要补充什么资料
}
```

#### POST /api/v1/consultations/:id/close

标记结案

**Request:**
```typescript
interface CloseConsultationRequest {
    internal_category?: "simple" | "complex"
    complex_sub_category?: string
}
```

#### POST /api/v1/consultations/:id/transfer

变更处理人

**Request:**
```typescript
interface TransferConsultationRequest {
    new_handler_id: string
    reason?: string
}
```

#### POST /api/v1/consultations/:id/rate

评价

**Request:**
```typescript
interface RateConsultationRequest {
    rating: number  // 1-5
}
```

#### GET /api/v1/consultations/:id/similar

相似问题推荐

**Response:**
```typescript
interface SimilarConsultationsResponse {
    success: boolean
    data?: Array<{
        id: string
        ticket_no: string
        title: string
        similarity: number
    }>
}
```

---

### 3. 合同模板申请

#### POST /api/v1/template-requests

创建模板申请

**Request:**
```typescript
interface CreateTemplateRequestRequest {
    request_type: "new" | "update"
    existing_template_id?: string
    contract_type: string
    title: string
    business_scenario?: string
    business_flow?: string
    key_clauses?: string
    diff_from_existing?: string
    expected_date?: string
    attachments?: File[]
}
```

#### GET /api/v1/template-requests

获取申请列表

#### GET /api/v1/template-requests/:id

获取申请详情

#### POST /api/v1/template-requests/:id/approve

L1审批（业务主管）

**Request:**
```typescript
interface ApproveTemplateRequest {
    action: "approve" | "reject" | "return_for_supplement"
    comment?: string
}
```

#### POST /api/v1/template-requests/:id/draft

拟写模板（法务专员）

**Request:**
```typescript
interface DraftTemplateRequest {
    name: string
    description: string
    file: File  // docx
    editable_clauses?: string
    expected_date?: string
}
```

#### POST /api/v1/template-requests/:id/save-draft

保存草稿

#### POST /api/v1/template-requests/:id/review

审核（法务负责人）

**Request:**
```typescript
interface ReviewTemplateRequest {
    action: "approve" | "return_for_modification"
    comment?: string
}
```

---

### 4. 合同模板库

#### GET /api/v1/templates

获取模板列表

#### GET /api/v1/templates/:id

获取模板详情

#### GET /api/v1/templates/:id/download

下载模板

#### GET /api/v1/templates/:id/versions

获取版本历史

#### GET /api/v1/templates/compare

版本对比

**Query Parameters:**
```typescript
interface CompareVersionsQuery {
    version_a_id: string
    version_b_id: string
}
```

---

### 5. 法务工作台

#### GET /api/v1/legal/dashboard

获取工作台统计

**Response:**
```typescript
interface LegalDashboardResponse {
    success: boolean
    data?: {
        consultation_stats: {
            today_new: number
            today_replied: number
            pending: number
            week_closed: number
        }
        template_stats: {
            pending_draft: number
            pending_review: number
            today_new_requests: number
        }
        recent_records: Array<{
            id: string
            ticket_no: string
            title: string
            action_type: string
            created_at: string
        }>
    }
}
```

#### GET /api/v1/legal/consultation-pool

咨询池

#### GET /api/v1/legal/my-tasks

我的待办

#### GET /api/v1/legal/staff-list

获取法务专员列表（用于变更处理人）

---

### 6. 统计

#### GET /api/v1/statistics/overview

统计概览

**Query Parameters:**
```typescript
interface StatisticsQuery {
    start_date: string
    end_date: string
}
```

#### GET /api/v1/statistics/export

导出统计报表

---

### 7. 搜索与案例

#### GET /api/v1/consultations/search

全文搜索

**Query Parameters:**
```typescript
interface SearchQuery {
    keyword: string
    category?: string
    status?: string
    start_date?: string
    end_date?: string
    page?: number
    page_size?: number
}
```

#### POST /api/v1/cases/:consultation_id/collect

收藏案例

#### DELETE /api/v1/cases/:consultation_id/collect

取消收藏

#### GET /api/v1/cases/my

我的案例库

---

### 8. 系统管理

#### GET /api/v1/admin/users

用户列表

#### POST /api/v1/admin/users

创建用户

#### PUT /api/v1/admin/users/:id

编辑用户

#### POST /api/v1/admin/users/:id/reset-password

重置密码

#### PUT /api/v1/admin/users/:id/toggle-status

启用/禁用用户

#### GET /api/v1/admin/departments

部门列表

#### POST /api/v1/admin/departments

创建部门

#### PUT /api/v1/admin/departments/:id

编辑部门

#### DELETE /api/v1/admin/departments/:id

删除部门

#### GET /api/v1/admin/contract-types

合同类型列表

#### PUT /api/v1/admin/contract-types

更新合同类型

#### GET /api/v1/admin/consultation-types

咨询类型配置

#### PUT /api/v1/admin/consultation-types/:type

更新咨询类型配置

#### GET /api/v1/admin/system-config

系统配置

#### PUT /api/v1/admin/system-config

更新系统配置

#### POST /api/v1/admin/system-config/test-notification

测试通知

---

## 错误码

| 错误码 | 说明 |
|--------|------|
| ERR_AUTH_INVALID_CREDENTIALS | 工号或密码错误 |
| ERR_AUTH_TOKEN_EXPIRED | Token已过期 |
| ERR_AUTH_FORBIDDEN | 无权限访问 |
| ERR_CONSULTATION_NOT_FOUND | 咨询不存在 |
| ERR_CONSULTATION_ALREADY_ACCEPTED | 咨询已被接单 |
| ERR_CONSULTATION_ALREADY_CLOSED | 咨询已结案 |
| ERR_CONSULTATION_CANNOT_TRANSFER | 无法转让给同一人 |
| ERR_TEMPLATE_NOT_FOUND | 模板不存在 |
| ERR_TEMPLATE_ALREADY_PUBLISHED | 模板已发布 |
| ERR_TEMPLATE_REQUEST_NOT_FOUND | 申请不存在 |
| ERR_USER_NOT_FOUND | 用户不存在 |
| ERR_DEPARTMENT_NOT_EMPTY | 部门下有用户 |
| ERR_FILE_TYPE_NOT_ALLOWED | 文件类型不允许 |
| ERR_FILE_SIZE_EXCEEDED | 文件大小超出限制 |
| ERR_NOTIFICATION_FAILED | 通知发送失败 |

---

## 文件变更清单

### 后端 (Go)

| 文件 | 操作 | 内容 |
|------|------|------|
| backend/main.go | 新增 | 应用入口，路由注册 |
| backend/go.mod | 新增 | Go模块定义 |
| backend/go.sum | 新增 | 依赖锁定 |
| backend/internal/config/config.go | 新增 | 配置加载（YAML） |
| backend/internal/models/user.go | 新增 | User, Department模型 |
| backend/internal/models/consultation.go | 新增 | Consultation, ConsultationReply模型 |
| backend/internal/models/attachment.go | 新增 | Attachment模型 |
| backend/internal/models/template.go | 新增 | TemplateRequest, Template模型 |
| backend/internal/models/config.go | 新增 | NotificationConfig, OperationLog, CaseCollection模型 |
| backend/internal/models/consultation_config.go | 新增 | ConsultationTypeConfig模型 |
| backend/internal/database/database.go | 新增 | 数据库连接初始化 |
| backend/internal/repository/user_repo.go | 新增 | UserRepository用户仓储 |
| backend/internal/repository/consultation_repo.go | 新增 | ConsultationRepository咨询仓储 |
| backend/internal/repository/template_repo.go | 新增 | TemplateRepository模板仓储 |
| backend/internal/repository/attachment_repo.go | 新增 | AttachmentRepository附件仓储 |
| backend/internal/service/auth_service.go | 新增 | AuthService认证服务 |
| backend/internal/service/consultation_service.go | 新增 | ConsultationService咨询服务 |
| backend/internal/service/template_service.go | 新增 | TemplateService模板服务 |
| backend/internal/service/notification_service.go | 新增 | NotificationService通知服务 |
| backend/internal/service/statistics_service.go | 新增 | StatisticsService统计服务 |
| backend/internal/service/case_service.go | 新增 | CaseService案例服务 |
| backend/internal/handler/auth_handler.go | 新增 | AuthHandler认证处理器 |
| backend/internal/handler/consultation_handler.go | 新增 | ConsultationHandler咨询处理器 |
| backend/internal/handler/template_handler.go | 新增 | TemplateHandler模板处理器 |
| backend/internal/handler/legal_handler.go | 新增 | LegalHandler法务工作台处理器 |
| backend/internal/handler/statistics_handler.go | 新增 | StatisticsHandler统计处理器 |
| backend/internal/handler/admin_handler.go | 新增 | AdminHandler系统管理处理器 |
| backend/internal/middleware/auth.go | 新增 | 认证中间件 |
| backend/internal/middleware/permission.go | 新增 | 权限中间件 |
| backend/internal/middleware/logging.go | 新增 | 日志中间件 |
| backend/internal/utils/response.go | 新增 | 统一响应工具 |
| backend/internal/utils/password.go | 新增 | 密码加密工具 |
| backend/internal/utils/ticket_no.go | 新增 | 工单号生成工具 |
| backend/internal/utils/file.go | 新增 | 文件处理工具 |
| backend/config.yaml | 新增 | 配置文件 |

### 前端 (Next.js)

| 文件 | 操作 | 内容 |
|------|------|------|
| frontend/package.json | 新增 | 项目依赖 |
| frontend/next.config.js | 新增 | Next.js配置 |
| frontend/tailwind.config.ts | 新增 | Tailwind配置 |
| frontend/tsconfig.json | 新增 | TypeScript配置 |
| frontend/src/app/layout.tsx | 新增 | 根布局 |
| frontend/src/app/page.tsx | 新增 | 首页（分流页） |
| frontend/src/app/login/page.tsx | 新增 | 登录页 |
| frontend/src/app/(dashboard)/layout.tsx | 新增 | 仪表盘布局 |
| frontend/src/app/(dashboard)/consultations/page.tsx | 新增 | 我的咨询列表页 |
| frontend/src/app/(dashboard)/consultations/new/page.tsx | 新增 | 发起咨询页 |
| frontend/src/app/(dashboard)/consultations/[id]/page.tsx | 新增 | 咨询详情页 |
| frontend/src/app/(dashboard)/template-requests/page.tsx | 新增 | 我的申请列表页 |
| frontend/src/app/(dashboard)/template-requests/new/page.tsx | 新增 | 发起申请页 |
| frontend/src/app/(dashboard)/legal/layout.tsx | 新增 | 法务布局 |
| frontend/src/app/(dashboard)/legal/dashboard/page.tsx | 新增 | 法务工作台 |
| frontend/src/app/(dashboard)/legal/pool/page.tsx | 新增 | 咨询池 |
| frontend/src/app/(dashboard)/legal/tasks/page.tsx | 新增 | 我的待办 |
| frontend/src/app/(dashboard)/legal/templates/page.tsx | 新增 | 待拟写模板 |
| frontend/src/app/(dashboard)/legal/templates/[id]/draft/page.tsx | 新增 | 拟写模板页 |
| frontend/src/app/(dashboard)/legal/review/page.tsx | 新增 | 待审核列表 |
| frontend/src/app/(dashboard)/legal/review/[id]/page.tsx | 新增 | 审核页 |
| frontend/src/app/(dashboard)/legal/statistics/page.tsx | 新增 | 统计页 |
| frontend/src/app/(dashboard)/legal/cases/page.tsx | 新增 | 案例库 |
| frontend/src/app/(dashboard)/legal/search/page.tsx | 新增 | 搜索页 |
| frontend/src/app/(dashboard)/template-library/page.tsx | 新增 | 模板库 |
| frontend/src/app/(dashboard)/template-library/[id]/page.tsx | 新增 | 模板详情 |
| frontend/src/app/(dashboard)/admin/users/page.tsx | 新增 | 用户管理 |
| frontend/src/app/(dashboard)/admin/departments/page.tsx | 新增 | 部门管理 |
| frontend/src/app/(dashboard)/admin/contract-types/page.tsx | 新增 | 合同类型管理 |
| frontend/src/app/(dashboard)/admin/consultation-types/page.tsx | 新增 | 咨询类型配置 |
| frontend/src/app/(dashboard)/admin/system/page.tsx | 新增 | 系统设置 |
| frontend/src/components/ui/* | 新增 | Shadcn UI组件 |
| frontend/src/components/shared/* | 新增 | 共享组件 |
| frontend/src/lib/api.ts | 新增 | API客户端 |
| frontend/src/lib/auth.ts | 新增 | 认证工具 |
| frontend/src/lib/constants.ts | 新增 | 常量定义 |
| frontend/src/hooks/* | 新增 | React Hooks |
| frontend/src/types/* | 新增 | TypeScript类型定义 |

---

## 引用的已有代码

本项目为全新系统，无已有代码引用。

---

## 影响分析

| 已有功能 | 影响 | 风险等级 |
|---------|------|---------|
| 无 | - | - |

---

## 技术决策

| 决策 | 选择 | 理由 |
|------|------|------|
| Web框架 | Gin | Go生态成熟，性能优秀 |
| ORM | GORM | Go最流行的ORM，社区活跃 |
| 数据库 | PostgreSQL | 需求指定，支持JSONB存储扩展信息 |
| 附件存储 | 本地文件系统 | 需求指定，/uploads目录 |
| 前端框架 | Next.js App Router | 最新Next.js架构 |
| UI组件库 | Shadcn UI | Tailwind + Radix，定制灵活 |
| 状态管理 | React Context + SWR | 简单场景够用，无需引入Redux |
| 认证方式 | JWT Token | 无状态，适合API服务 |

---

## 风险点

| 风险 | 影响 | 应对 |
|------|------|------|
| 文件上传安全 | 可能上传恶意文件 | 后端严格校验Content-Type，限制文件扩展名 |
| 大文件上传 | 占用服务器带宽和存储 | 限制单文件50MB，多文件10个 |
| 通知失败 | 用户无法及时收到通知 | 通知发送失败不影响主流程，记录日志 |
| 并发接单 | 同一咨询被多人接单 | 使用数据库事务+乐观锁 |
| 文本搜索性能 | 大量数据全文搜索慢 | 使用PostgreSQL全文索引，考虑后期引入Elasticsearch |

---

## 需要人决策

- [ ] 附件存储路径：使用 `/uploads/` 还是云存储（如阿里云OSS）？
- [ ] 是否需要引入缓存（Redis）？
- [ ] 版本对比功能：前端实现还是后端实现？

