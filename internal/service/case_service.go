package service

import (
	"errors"

	"github.com/google/uuid"

	"legal-consultation/internal/models"
	"legal-consultation/internal/repository"
)

type CaseService struct {
	caseRepo         *repository.CaseRepository
	consultationRepo *repository.ConsultationRepository
}

func NewCaseService(
	caseRepo *repository.CaseRepository,
	consultationRepo *repository.ConsultationRepository,
) *CaseService {
	return &CaseService{
		caseRepo:         caseRepo,
		consultationRepo: consultationRepo,
	}
}

// 收藏案例
func (s *CaseService) CollectCase(consultationID, collectorID string, tags []string) error {
	// 检查是否已收藏
	existing, err := s.caseRepo.GetByConsultationID(consultationID)
	if err == nil && existing != nil {
		return errors.New("该咨询已被收藏")
	}

	// 检查咨询是否存在且已结案
	consultation, err := s.consultationRepo.GetByID(consultationID)
	if err != nil {
		return errors.New("咨询不存在")
	}
	if consultation.Status != models.ConsultationStatusClosed {
		return errors.New("只能收藏已结案的咨询")
	}

	caseCollection := &models.CaseCollection{
		ID:             uuid.New().String(),
		ConsultationID: consultationID,
		CollectorID:    collectorID,
		Tags:           models.JSONType{"tags": tags},
	}

	return s.caseRepo.Create(caseCollection)
}

// 取消收藏
func (s *CaseService) UncollectCase(consultationID, collectorID string) error {
	return s.caseRepo.Delete(consultationID, collectorID)
}

// 检查是否已收藏
func (s *CaseService) IsCollected(consultationID, collectorID string) bool {
	caseCollection, err := s.caseRepo.GetByConsultationID(consultationID)
	if err != nil || caseCollection == nil {
		return false
	}
	return caseCollection.CollectorID == collectorID
}

// 我的案例库
func (s *CaseService) ListMyCases(collectorID string, page, pageSize int) ([]models.CaseCollection, int64, error) {
	return s.caseRepo.ListByCollector(collectorID, page, pageSize)
}

// 搜索案例
func (s *CaseService) SearchCases(keyword, tag string, page, pageSize int) ([]models.CaseCollection, int64, error) {
	return s.caseRepo.Search(keyword, tag, page, pageSize)
}
