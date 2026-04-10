package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"legal-consultation/internal/config"
	"legal-consultation/internal/models"
	"legal-consultation/internal/repository"
	"legal-consultation/internal/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.JWTConfig
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.JWTConfig) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type LoginRequest struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  *UserInfo   `json:"user"`
}

type UserInfo struct {
	ID         string            `json:"id"`
	EmployeeID string            `json:"employee_id"`
	Name       string            `json:"name"`
	Role       string            `json:"role"`
	Department *DepartmentInfo  `json:"department"`
}

type DepartmentInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByEmployeeID(req.EmployeeID)
	if err != nil {
		return nil, errors.New("工号或密码错误")
	}

	// 检查状态
	if user.Status != models.UserStatusActive {
		return nil, errors.New("账号已被禁用")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("工号或密码错误")
	}

	// 生成 Token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	// 构建响应
	deptInfo := &DepartmentInfo{}
	if user.Department != nil {
		deptInfo.ID = user.Department.ID
		deptInfo.Name = user.Department.Name
	}

	return &LoginResponse{
		Token: token,
		User: &UserInfo{
			ID:         user.ID,
			EmployeeID: user.EmployeeID,
			Name:       user.Name,
			Role:       user.Role,
			Department: deptInfo,
		},
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*UserInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		user, err := s.userRepo.GetByID(userID)
		if err != nil {
			return nil, err
		}

		deptInfo := &DepartmentInfo{}
		if user.Department != nil {
			deptInfo.ID = user.Department.ID
			deptInfo.Name = user.Department.Name
		}

		return &UserInfo{
			ID:         user.ID,
			EmployeeID: user.EmployeeID,
			Name:       user.Name,
			Role:       user.Role,
			Department: deptInfo,
		}, nil
	}

	return nil, errors.New("无效的Token")
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"employee_id": user.EmployeeID,
		"role":       user.Role,
		"exp":        time.Now().Add(time.Hour * time.Duration(s.cfg.ExpireHours)).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

// 创建用户
func (s *AuthService) CreateUser(user *models.User) error {
	user.ID = uuid.New().String()
	// 密码已经在 Handler 层加密
	return s.userRepo.Create(user)
}

// 获取用户
func (s *AuthService) GetUser(id string) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// 更新用户
func (s *AuthService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

// 重置密码
func (s *AuthService) ResetPassword(id string) (string, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return "", err
	}

	// 生成随机密码
	newPassword := uuid.New().String()[:12]
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return "", err
	}

	user.Password = hash
	if err := s.userRepo.Update(user); err != nil {
		return "", err
	}

	return newPassword, nil
}
