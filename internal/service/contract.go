package service

import (
	"context"

	"github.com/vnchk1/subscription-aggregator/internal/models"
	"github.com/vnchk1/subscription-aggregator/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.SubscriptionResponse, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (*models.SubscriptionResponse, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, req *models.UpdateSubscriptionRequest) (*models.SubscriptionResponse, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	ListSubscriptions(ctx context.Context, userID *uuid.UUID, page, limit int) (*models.ListResponse, error)
	CalculateTotalCost(ctx context.Context, req *models.TotalCostRequest) (*models.TotalCostResponse, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		repo: repo,
	}
}
