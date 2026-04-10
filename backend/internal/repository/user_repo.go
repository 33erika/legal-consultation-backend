package repository

import (
	"gorm.io/gorm"

	"legal-consultation/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Department").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmployeeID(employeeID string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Department").First(&user, "employee_id = ?", employeeID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepository) List(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Model(&models.User{}).Count(&total)
	err := r.db.Preload("Department").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) Search(keyword string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{}).
		Where("employee_id LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	query.Count(&total)
	err := query.Preload("Department").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) ListByRole(role string) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role = ? AND status = ?", role, models.UserStatusActive).Find(&users).Error
	return users, err
}

func (r *UserRepository) ToggleStatus(id string, status string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *UserRepository) GetDepartmentUsers(departmentID string) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("department_id = ?", departmentID).Find(&users).Error
	return users, err
}
