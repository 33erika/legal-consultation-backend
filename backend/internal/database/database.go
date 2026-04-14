package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"legal-consultation/internal/config"
	"legal-consultation/internal/models"
)

var db *gorm.DB

func Initialize(cfg *config.DatabaseConfig) error {
	var err error

	// 支持 SQLite 作为开发数据库
	var dsn string
	if cfg.Driver == "sqlite" || cfg.Host == "sqlite" {
		dsn = cfg.DBName + ".db"
		log.Printf("Using SQLite database: %s", dsn)
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		// PostgreSQL
		dsn = cfg.DSN()
		log.Printf("Connecting to database: %s:%d/%s", cfg.Host, cfg.Port, cfg.DBName)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	// 初始化基础数据
	if err := seedData(); err != nil {
		log.Printf("Warning: seed data failed: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func autoMigrate() error {
	return db.AutoMigrate(
		// 基础表
		&models.Department{},
		&models.User{},

		// 咨询相关
		&models.Consultation{},
		&models.ConsultationReply{},
		&models.Attachment{},
		&models.ConsultationAttachment{},

		// 系统配置
		&models.NotificationConfig{},
		&models.OperationLog{},
		&models.CaseCollection{},
		&models.ConsultationTypeConfig{},
	)
}

func seedData() error {
	// 检查是否已有数据
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return nil // 已有数据，跳过
	}

	log.Println("Seeding initial data...")

	// 创建部门
	departments := []models.Department{
		{ID: "dept-001", Name: "法务部", ParentID: nil},
		{ID: "dept-002", Name: "人力资源部", ParentID: nil},
		{ID: "dept-003", Name: "销售部", ParentID: nil},
	}
	for _, dept := range departments {
		db.Create(&dept)
	}

	// 创建测试用户 (密码统一为: 密码, bcrypt hash)
	users := []models.User{
		{
			ID:           "user-admin",
			EmployeeID:   "admin",
			Password:     "$2a$10$zzsnejcVLottn6SpZy6ImOs8dxygqMHW4v1CHEVeewl61h64ITwmK", // 密码
			Name:         "系统管理员",
			Email:        "admin@example.com",
			Role:         models.RoleAdmin,
			DepartmentID: stringPtr("dept-001"),
			Status:       models.UserStatusActive,
		},
		{
			ID:           "user-legal-1",
			EmployeeID:   "legal001",
			Password:     "$2a$10$zzsnejcVLottn6SpZy6ImOs8dxygqMHW4v1CHEVeewl61h64ITwmK", // 密码
			Name:         "法务专员张三",
			Email:        "legal@example.com",
			Role:         models.RoleLegalStaff,
			DepartmentID: stringPtr("dept-001"),
			Status:       models.UserStatusActive,
		},
		{
			ID:           "user-legal-head",
			EmployeeID:   "legal-head",
			Password:     "$2a$10$zzsnejcVLottn6SpZy6ImOs8dxygqMHW4v1CHEVeewl61h64ITwmK", // 密码
			Name:         "法务负责人李四",
			Email:        "legal-head@example.com",
			Role:         models.RoleLegalHead,
			DepartmentID: stringPtr("dept-001"),
			Status:       models.UserStatusActive,
		},
		{
			ID:           "user-employee-1",
			EmployeeID:   "emp001",
			Password:     "$2a$10$zzsnejcVLottn6SpZy6ImOs8dxygqMHW4v1CHEVeewl61h64ITwmK", // 密码
			Name:         "员工王五",
			Email:        "employee@example.com",
			Role:         models.RoleEmployee,
			DepartmentID: stringPtr("dept-003"),
			Status:       models.UserStatusActive,
		},
	}
	for _, user := range users {
		db.Create(&user)
	}

	// 创建咨询类型配置
	consultationTypes := []models.ConsultationTypeConfig{
		{
			ID:        "ct-complaint",
			Type:      "complaint",
			Name:      "客诉类",
			Keywords:  models.JSONType{"keywords": []string{"投诉", "消费者", "12315", "打假", "退款", "赔偿"}},
			Fields:    models.JSONType{"fields": []map[string]interface{}{
				{"name": "客诉来源", "type": "select", "options": []string{"12315平台", "第三方投诉平台", "内部客服", "举报", "其他"}, "required": false},
				{"name": "投诉编号/工单号", "type": "text", "required": false},
				{"name": "消费者诉求", "type": "multiselect", "options": []string{"仅退款", "退一赔三", "道歉", "要求下架产品", "经济补偿", "其他"}, "required": false},
				{"name": "是否为职业打假人", "type": "radio", "options": []string{"是", "否"}, "required": false},
				{"name": "当前诉求", "type": "textarea", "max_length": 500, "required": false},
			}},
			SortOrder: 1,
			Enabled:   true,
		},
		{
			ID:        "ct-contract",
			Type:      "contract",
			Name:      "合同类",
			Keywords:  models.JSONType{"keywords": []string{"合同", "协议", "签署", "盖章", "条款"}},
			Fields:    models.JSONType{"fields": []map[string]interface{}{
				{"name": "合同类型", "type": "select", "options": []string{"采购合同", "销售合同", "租赁合同", "服务合同", "劳动合同", "其他"}, "required": false},
				{"name": "合同相对方", "type": "text", "required": false},
				{"name": "合同金额", "type": "number", "required": false},
				{"name": "合同是否已签署", "type": "radio", "options": []string{"是", "否"}, "required": false},
				{"name": "主要争议点", "type": "textarea", "max_length": 500, "required": false},
			}},
			SortOrder: 2,
			Enabled:   true,
		},
		{
			ID:        "ct-labor",
			Type:      "labor",
			Name:      "劳动类",
			Keywords:  models.JSONType{"keywords": []string{"劳动", "工资", "赔偿", "加班", "社保", "解除", "试用期", "离职"}},
			Fields:    models.JSONType{"fields": []map[string]interface{}{
				{"name": "员工身份", "type": "select", "options": []string{"在职员工", "离职员工", "试用期员工", "实习生"}, "required": false},
				{"name": "问题类型", "type": "select", "options": []string{"工资", "赔偿金", "加班费", "社保", "合同解除", "其他"}, "required": false},
				{"name": "涉及人数", "type": "number", "required": false},
				{"name": "当前诉求", "type": "textarea", "max_length": 500, "required": false},
			}},
			SortOrder: 3,
			Enabled:   true,
		},
		{
			ID:        "ct-ip",
			Type:      "ip",
			Name:      "知识产权类",
			Keywords:  models.JSONType{"keywords": []string{"商标", "专利", "著作权", "版权", "侵权", "抄袭"}},
			Fields:    models.JSONType{"fields": []map[string]interface{}{
				{"name": "知识产权类型", "type": "select", "options": []string{"商标", "专利", "著作权", "商业秘密", "其他"}, "required": false},
				{"name": "是否已进行登记", "type": "radio", "options": []string{"是", "否"}, "required": false},
				{"name": "登记号/证书号", "type": "text", "required": false},
				{"name": "当前诉求", "type": "textarea", "max_length": 500, "required": false},
			}},
			SortOrder: 4,
			Enabled:   true,
		},
		{
			ID:        "ct-dispute",
			Type:      "dispute",
			Name:      "纠纷咨询类",
			Keywords:  models.JSONType{"keywords": []string{"纠纷", "起诉", "诉讼", "判决", "执行", "对方当事人"}},
			Fields:    models.JSONType{"fields": []map[string]interface{}{
				{"name": "当前诉求", "type": "textarea", "max_length": 500, "required": false},
				{"name": "对方当事人", "type": "text", "required": false},
				{"name": "纠纷金额", "type": "number", "required": false},
				{"name": "案件概要", "type": "textarea", "max_length": 500, "required": false},
			}},
			SortOrder: 5,
			Enabled:   true,
		},
	}
	for _, ct := range consultationTypes {
		db.Create(&ct)
	}

	log.Println("Seed data created successfully")
	return nil
}

func GetDB() *gorm.DB {
	return db
}

func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// 辅助函数
func stringPtr(s string) *string {
	return &s
}
