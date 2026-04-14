package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"legal-consultation/internal/config"
	"legal-consultation/internal/models"
)

type NotificationService struct {
	cfg *config.DingTalkConfig
}

func NewNotificationService(cfg *config.DingTalkConfig) *NotificationService {
	return &NotificationService{cfg: cfg}
}

// 钉钉消息格式
type DingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    *TextContent `json:"text,omitempty"`
}

type TextContent struct {
	Content string `json:"content"`
}

// 发送通知
func (s *NotificationService) send(message string) error {
	if !s.cfg.Enabled || s.cfg.WebhookURL == "" {
		log.Printf("DingTalk notification disabled or no webhook URL configured, skipping: %s", message)
		return nil
	}

	msg := DingTalkMessage{
		MsgType: "text",
		Text:    &TextContent{Content: message},
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.cfg.WebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to send DingTalk notification: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("DingTalk notification returned non-200 status: %d", resp.StatusCode)
		return err
	}

	log.Printf("DingTalk notification sent: %s", message)
	return nil
}

// 新咨询通知（通知法务）
func (s *NotificationService) NotifyNewConsultation(c *models.Consultation) {
	message := "📋 您有一个新的法律咨询待处理\n" +
		"工单号：" + c.TicketNo + "\n" +
		"紧急度：" + getUrgencyText(c.Urgency) + "\n" +
		"标题：" + c.Title
	s.send(message)
}

// 咨询被接单通知（通知员工）
func (s *NotificationService) NotifyConsultationAccepted(c *models.Consultation) {
	handlerName := ""
	if c.Handler != nil {
		handlerName = c.Handler.Name
	}
	message := "🔔 您的咨询已被接单\n" +
		"工单号：" + c.TicketNo + "\n" +
		"处理人：" + handlerName
	s.send(message)
}

// 咨询有新回复通知（通知员工）
func (s *NotificationService) NotifyConsultationReplied(c *models.Consultation) {
	message := "💬 您的咨询有新回复\n" +
		"工单号：" + c.TicketNo + "\n" +
		"点击查看详情"
	s.send(message)
}

// 咨询已结案通知（通知员工）
func (s *NotificationService) NotifyConsultationClosed(c *models.Consultation) {
	message := "✅ 您的咨询已处理完毕\n" +
		"工单号：" + c.TicketNo + "\n" +
		"请对本次服务进行评价"
	s.send(message)
}

// 咨询处理人变更通知
func (s *NotificationService) NotifyConsultationTransferred(c *models.Consultation, oldHandlerID *string) {
	newHandlerName := ""
	if c.Handler != nil {
		newHandlerName = c.Handler.Name
	}
	message := "📤 咨询已被转交\n" +
		"工单号：" + c.TicketNo + "\n" +
		"新处理人：" + newHandlerName
	s.send(message)
}

// 测试通知
func (s *NotificationService) TestNotification() error {
	return s.send("这是一条测试消息，用于验证钉钉机器人配置是否正确。")
}

func getUrgencyText(urgency string) string {
	switch urgency {
	case "very_urgent":
		return "非常紧急"
	case "urgent":
		return "紧急"
	default:
		return "一般"
	}
}

// 获取通知配置
func (s *NotificationService) GetConfig() *config.DingTalkConfig {
	return s.cfg
}

// 更新通知配置
func (s *NotificationService) UpdateConfig(cfg *config.DingTalkConfig) {
	s.cfg = cfg
}

// 辅助函数：检查配置是否有效
func (s *NotificationService) IsConfigured() bool {
	return s.cfg != nil && s.cfg.Enabled && s.cfg.WebhookURL != ""
}
