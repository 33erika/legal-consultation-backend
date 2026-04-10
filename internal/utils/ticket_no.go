package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	mu         sync.Mutex
	seqCounter int
	lastDate   string
)

// GenerateTicketNo 生成工单号
// 格式: CONS-YYYYMMDD-序号 (如 CONS-20260410-001)
func GenerateTicketNo(prefix string) string {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	today := now.Format("20060102")

	// 如果日期变了，重置计数器
	if today != lastDate {
		lastDate = today
		seqCounter = 0
	}

	seqCounter++
	return fmt.Sprintf("%s-%s-%03d", prefix, today, seqCounter)
}

// GenerateConsultationTicketNo 生成咨询工单号
func GenerateConsultationTicketNo() string {
	return GenerateTicketNo("CONS")
}

// GenerateTemplateRequestNo 生成模板申请编号
func GenerateTemplateRequestNo() string {
	return GenerateTicketNo("TMPL")
}
