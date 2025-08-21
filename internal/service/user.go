package service

import (
	"fmt"
	"go-dev/internal/models"
	"go-dev/internal/repository"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error) {
	// Проверяем, что пользователь с таким email не существует
	existing, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	err = s.repo.Create(user)
	return user, err
}

func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserService) List(limit, offset int) ([]*models.User, error) {
	return s.repo.List(limit, offset)
}

func (s *UserService) Update(id uuid.UUID, req *models.UpdateUserRequest) error {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		// Проверяем, что email уникален
		existing, err := s.repo.GetByEmail(*req.Email)
		if err != nil {
			return err
		}
		if existing != nil && existing.ID != id {
			return fmt.Errorf("user with email %s already exists", *req.Email)
		}
		updates["email"] = *req.Email
	}

	return s.repo.Update(id, updates)
}

func (s *UserService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
