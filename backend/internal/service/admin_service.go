package service

import (
	"legal-consultation/internal/repository"
)

type AdminService struct {
	userRepo         *repository.UserRepository
	notificationSvc *NotificationService
}

func NewAdminService(
	userRepo *repository.UserRepository,
	notificationSvc *NotificationService,
) *AdminService {
	return &AdminService{
		userRepo:         userRepo,
		notificationSvc: notificationSvc,
	}
}

func (s *AdminService) GetUserRepo() *repository.UserRepository {
	return s.userRepo
}

func (s *AdminService) GetNotificationSvc() *NotificationService {
	return s.notificationSvc
}
