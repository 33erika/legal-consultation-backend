package handler

import (
	"github.com/gin-gonic/gin"

	"legal-consultation/internal/service"
	"legal-consultation/internal/utils"
)

type StatisticsHandler struct {
	statisticsSvc *service.StatisticsService
}

func NewStatisticsHandler(statisticsSvc *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{statisticsSvc: statisticsSvc}
}

// Overview 统计概览
func (h *StatisticsHandler) Overview(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.statisticsSvc.GetOverview(startDate, endDate)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

// CategoryDistribution 分类分布
func (h *StatisticsHandler) CategoryDistribution(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.statisticsSvc.GetCategoryDistribution(startDate, endDate)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

// ProcessingEfficiency 处理效率
func (h *StatisticsHandler) ProcessingEfficiency(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.statisticsSvc.GetProcessingEfficiency(startDate, endDate)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

// Export 导出报表
func (h *StatisticsHandler) Export(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := h.statisticsSvc.ExportReport(startDate, endDate)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, data)
}

// StaffWorkload 法务专员工作量
func (h *StatisticsHandler) StaffWorkload(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := h.statisticsSvc.GetStaffWorkload(startDate, endDate)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, data)
}
