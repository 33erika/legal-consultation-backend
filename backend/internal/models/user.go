package models

import (
	"time"
)

// 用户角色
const (
	RoleEmployee     = "employee"      // 普通员工
	RoleSupervisor   = "supervisor"    // 业务主管
	RoleLegalStaff   = "legal_staff"   // 法务专员
	RoleLegalHead    = "legal_head"    // 法务负责人
	RoleAdmin        = "admin"         // 系统管理员
)

// 用户状态
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
)

type User struct {
	ID           string      `gorm:"primaryKey;size:36" json:"id"`
	EmployeeID   string      `gorm:"uniqueIndex;size:32;not null" json:"employee_id"` // 工号
	Password     string      `gorm:"size:255;not null" json:"-"`                     // bcrypt加密
	Name         string      `gorm:"size:100;not null" json:"name"`
	Email        string      `gorm:"size:255" json:"email"`
	Role         string      `gorm:"size:50;not null" json:"role"` // employee/supervisor/legal_staff/legal_head/admin
	DepartmentID *string     `gorm:"size:36" json:"department_id"`
	Status       string      `gorm:"size:20;default:active" json:"status"` // active/inactive
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`

	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type Department struct {
	ID        string       `gorm:"primaryKey;size:36" json:"id"`
	Name      string       `gorm:"size:100;not null" json:"name"`
	ParentID  *string      `gorm:"size:36" json:"parent_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	Parent   *Department  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Users    []User       `gorm:"foreignKey:DepartmentID" json:"users,omitempty"`
}

func (Department) TableName() string {
	return "departments"
}
