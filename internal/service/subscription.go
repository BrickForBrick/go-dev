package service

import (
	"fmt"
	"go-dev/internal/models"
	"go-dev/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	sub := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	err := s.repo.Create(sub)
	return sub, err
}

func (s *SubscriptionService) GetByID(id int) (*models.Subscription, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, fmt.Errorf("subscription not found")
	}
	return sub, nil
}

func (s *SubscriptionService) List(userID *uuid.UUID, serviceName *string, limit, offset int) ([]*models.Subscription, error) {
	return s.repo.List(userID, serviceName, limit, offset)
}

func (s *SubscriptionService) Update(id int, req *models.UpdateSubscriptionRequest) error {
	updates := make(map[string]interface{})

	if req.ServiceName != nil {
		updates["service_name"] = *req.ServiceName
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.StartDate != nil {
		updates["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = *req.EndDate
	}

	return s.repo.Update(id, updates)
}

func (s *SubscriptionService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *SubscriptionService) GetTotalCost(userID *uuid.UUID, serviceName *string, startPeriod, endPeriod string) (*models.TotalCostResponse, error) {
	totalCost, err := s.repo.GetTotalCost(userID, serviceName, startPeriod, endPeriod)
	if err != nil {
		return nil, err
	}

	filters := make(map[string]string)
	if userID != nil {
		filters["user_id"] = userID.String()
	}
	if serviceName != nil {
		filters["service_name"] = *serviceName
	}

	return &models.TotalCostResponse{
		TotalCost: totalCost,
		Period:    fmt.Sprintf("%s to %s", startPeriod, endPeriod),
		Filters:   filters,
	}, nil
}
