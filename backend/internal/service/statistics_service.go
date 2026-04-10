package service

import (
	"time"

	"legal-consultation/internal/database"
	"legal-consultation/internal/models"
	"legal-consultation/internal/repository"
)

type StatisticsService struct {
	consultationRepo *repository.ConsultationRepository
	templateRepo     *repository.TemplateRepository
}

func NewStatisticsService(
	consultationRepo *repository.ConsultationRepository,
	templateRepo *repository.TemplateRepository,
) *StatisticsService {
	return &StatisticsService{
		consultationRepo: consultationRepo,
		templateRepo:     templateRepo,
	}
}

// 统计概览
func (s *StatisticsService) GetOverview(startDate, endDate string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 咨询统计
	consultationStats, err := s.consultationRepo.GetStats()
	if err != nil {
		return nil, err
	}
	result["consultation_stats"] = consultationStats

	// 模板统计
	templateStats, err := s.templateRepo.GetRequestStats()
	if err != nil {
		return nil, err
	}
	result["template_stats"] = templateStats

	return result, nil
}

// 咨询分类分布
func (s *StatisticsService) GetCategoryDistribution(startDate, endDate string) (map[string]interface{}, error) {
	// 简化实现，实际应按日期范围查询
	result := make(map[string]interface{})

	db := database.GetDB()
	var simpleCount, complexCount int64
	db.Model(&models.Consultation{}).Where("internal_category = ?", models.InternalCategorySimple).Count(&simpleCount)
	db.Model(&models.Consultation{}).Where("internal_category = ?", models.InternalCategoryComplex).Count(&complexCount)

	result["simple"] = simpleCount
	result["complex"] = complexCount

	return result, nil
}

// 处理效率统计
func (s *StatisticsService) GetProcessingEfficiency(startDate, endDate string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	db := database.GetDB()

	// 计算平均处理时长
	var avgDuration float64
	db.Model(&models.Consultation{}).
		Where("closed_at IS NOT NULL").
		Select("AVG(EXTRACT(EPOCH FROM (closed_at - submitted_at)) / 3600)").
		Scan(&avgDuration)

	result["avg_processing_hours"] = avgDuration

	// 计算平均响应时长
	var avgResponse float64
	db.Model(&models.Consultation{}).
		Where("first_replied_at IS NOT NULL").
		Select("AVG(EXTRACT(EPOCH FROM (first_replied_at - submitted_at)) / 3600)").
		Scan(&avgResponse)

	result["avg_response_hours"] = avgResponse

	return result, nil
}

// 导出统计报表
func (s *StatisticsService) ExportReport(startDate, endDate string) ([]map[string]interface{}, error) {
	db := database.GetDB()
	// 获取所有已结案的咨询
	var consultations []models.Consultation
	db.
		Preload("Submitter").
		Preload("Handler").
		Where("status = ?", models.ConsultationStatusClosed).
		Where("submitted_at >= ? AND submitted_at <= ?", startDate, endDate).
		Find(&consultations)

	var result []map[string]interface{}
	for _, c := range consultations {
		processingHours := 0.0
		if c.ClosedAt != nil {
			processingHours = c.ClosedAt.Sub(c.SubmittedAt).Hours()
		}

		rating := 0
		if c.Rating != nil {
			rating = *c.Rating
		}

		submitterName := ""
		if c.Submitter != nil {
			submitterName = c.Submitter.Name
		}
		handlerName := ""
		if c.Handler != nil {
			handlerName = c.Handler.Name
		}

		closedAtStr := ""
		if c.ClosedAt != nil {
			closedAtStr = c.ClosedAt.Format(time.RFC3339)
		}

		result = append(result, map[string]interface{}{
			"ticket_no":        c.TicketNo,
			"category":        c.InternalCategory,
			"sub_category":    c.ComplexSubCategory,
			"submitter":       submitterName,
			"submitted_at":    c.SubmittedAt.Format(time.RFC3339),
			"closed_at":       closedAtStr,
			"processing_hours": processingHours,
			"handler":         handlerName,
			"rating":          rating,
		})
	}

	return result, nil
}

// 法务专员工作量统计
func (s *StatisticsService) GetStaffWorkload(startDate, endDate string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	db := database.GetDB()

	var handlers []struct {
		HandlerID string
		Count     int64
	}

	db.Model(&models.Consultation{}).
		Select("handler_id, COUNT(*) as count").
		Where("closed_at >= ? AND closed_at <= ?", startDate, endDate).
		Group("handler_id").
		Scan(&handlers)

	for _, h := range handlers {
		results = append(results, map[string]interface{}{
			"handler_id": h.HandlerID,
			"count":      h.Count,
		})
	}

	return results, nil
}
