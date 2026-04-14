package com.legalsys.service;

import com.legalsys.entity.Consultation;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.Map;

@Service
@RequiredArgsConstructor
@Slf4j
public class NotificationService {

    @Value("${dingtalk.webhook-url:}")
    private String webhookUrl;

    @Value("${dingtalk.enabled:false}")
    private boolean enabled;

    private final HttpClient httpClient = HttpClient.newHttpClient();

    public void notifyNewConsultation(Consultation c) {
        if (!enabled || webhookUrl == null || webhookUrl.isEmpty()) {
            log.info("DingTalk disabled, skipping notification for new consultation: {}", c.getTicketNo());
            return;
        }

        String message = String.format(
            "📋 您有一个新的法律咨询待处理\n工单号：%s\n紧急度：%s\n标题：%s",
            c.getTicketNo(),
            getUrgencyText(c.getUrgency()),
            c.getTitle()
        );
        send(message);
    }

    public void notifyConsultationAccepted(Consultation c) {
        if (!enabled || webhookUrl == null || webhookUrl.isEmpty()) {
            log.info("DingTalk disabled, skipping notification for accepted consultation");
            return;
        }

        String handlerName = c.getHandler() != null ? c.getHandler().getName() : "未知";
        String message = String.format(
            "🔔 您的咨询已被接单\n工单号：%s\n处理人：%s",
            c.getTicketNo(),
            handlerName
        );
        send(message);
    }

    public void notifyConsultationReplied(Consultation c) {
        if (!enabled || webhookUrl == null || webhookUrl.isEmpty()) {
            log.info("DingTalk disabled, skipping notification for replied consultation");
            return;
        }

        String message = String.format(
            "💬 您的咨询有新回复\n工单号：%s\n点击查看详情",
            c.getTicketNo()
        );
        send(message);
    }

    public void notifyConsultationClosed(Consultation c) {
        if (!enabled || webhookUrl == null || webhookUrl.isEmpty()) {
            log.info("DingTalk disabled, skipping notification for closed consultation");
            return;
        }

        String message = String.format(
            "✅ 您的咨询已处理完毕\n工单号：%s\n请对本次服务进行评价",
            c.getTicketNo()
        );
        send(message);
    }

    public void notifyConsultationTransferred(Consultation c, String oldHandlerId) {
        if (!enabled || webhookUrl == null || webhookUrl.isEmpty()) {
            log.info("DingTalk disabled, skipping notification for transferred consultation");
            return;
        }

        String handlerName = c.getHandler() != null ? c.getHandler().getName() : "未知";
        String message = String.format(
            "📤 咨询已被转交\n工单号：%s\n新处理人：%s",
            c.getTicketNo(),
            handlerName
        );
        send(message);
    }

    public boolean testNotification() {
        if (!enabled || webhookUrl == null || webhookUrl.isEmpty()) {
            log.info("DingTalk disabled, cannot send test notification");
            return false;
        }
        send("这是一条测试消息，用于验证钉钉机器人配置是否正确。");
        return true;
    }

    private void send(String content) {
        try {
            String body = String.format("{\"msgtype\":\"text\",\"text\":{\"content\":\"%s\"}}", content);
            HttpRequest request = HttpRequest.newBuilder()
                    .uri(URI.create(webhookUrl))
                    .header("Content-Type", "application/json")
                    .POST(HttpRequest.BodyPublishers.ofString(body))
                    .build();

            HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
            if (response.statusCode() == 200) {
                log.info("Notification sent successfully: {}", content.substring(0, Math.min(50, content.length())));
            } else {
                log.error("Failed to send notification: {}", response.body());
            }
        } catch (Exception e) {
            log.error("Error sending notification: {}", e.getMessage());
        }
    }

    private String getUrgencyText(String urgency) {
        return switch (urgency) {
            case "very_urgent" -> "非常紧急";
            case "urgent" -> "紧急";
            default -> "一般";
        };
    }
}
