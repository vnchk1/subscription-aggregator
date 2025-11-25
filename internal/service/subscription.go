package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vnchk1/subscription-aggregator/internal/models"

	"github.com/google/uuid"
)

func (s *subscriptionService) CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (*models.SubscriptionResponse, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	subscription := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     nil,
	}

	if err = s.validateSubscription(subscription); err != nil {
		return nil, err
	}

	if err = s.repo.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return s.toResponse(subscription), nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id uuid.UUID) (*models.SubscriptionResponse, error) {
	if id == uuid.Nil {
		return nil, errors.New("subscription ID is required")
	}

	subscription, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return s.toResponse(subscription), nil
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id uuid.UUID, req *models.UpdateSubscriptionRequest) (*models.SubscriptionResponse, error) {
	if id == uuid.Nil {
		return nil, errors.New("subscription ID is required")
	}

	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing subscription: %w", err)
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	var endDate *time.Time

	if req.EndDate != nil {
		parsed, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}

		endDate = &parsed
	}

	existing.ServiceName = req.ServiceName
	existing.Price = req.Price
	existing.StartDate = startDate
	existing.EndDate = endDate

	if err := s.validateSubscription(existing); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return s.toResponse(existing), nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("subscription ID is required")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, userID *uuid.UUID, page, limit int) (*models.ListResponse, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	subscriptions, total, err := s.repo.List(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	responseData := make([]*models.SubscriptionResponse, len(subscriptions))
	for i, sub := range subscriptions {
		responseData[i] = s.toResponse(sub)
	}

	return &models.ListResponse{
		Total: total,
		Data:  responseData,
	}, nil
}

func (s *subscriptionService) CalculateTotalCost(ctx context.Context, req *models.TotalCostRequest) (*models.TotalCostResponse, error) {

	if err := s.validateTotalCostRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	startDate, endDate, err := req.ParseDates()
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	if endDate.Before(startDate) {
		return nil, errors.New("end period cannot be before start period")
	}

	filter := &models.SubscriptionFilter{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	totalCost, err := s.repo.GetTotalCost(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	return &models.TotalCostResponse{
		TotalCost: totalCost,
		Currency:  "RUB",
		Period:    fmt.Sprintf("%s - %s", req.StartPeriod, req.EndPeriod),
	}, nil
}

func (s *subscriptionService) validateCreateRequest(req *models.CreateSubscriptionRequest) error {
	if req.ServiceName == "" {
		return errors.New("service name is required")
	}

	if len(req.ServiceName) > 255 {
		return errors.New("service name too long")
	}

	if req.Price <= 0 {
		return errors.New("price must be positive")
	}

	if req.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}

	if req.StartDate == "" {
		return errors.New("start date is required")
	}

	return nil
}

func (s *subscriptionService) validateUpdateRequest(req *models.UpdateSubscriptionRequest) error {
	if req.ServiceName == "" {
		return errors.New("service name is required")
	}

	if len(req.ServiceName) > 255 {
		return errors.New("service name too long")
	}

	if req.Price <= 0 {
		return errors.New("price must be positive")
	}

	if req.StartDate == "" {
		return errors.New("start date is required")
	}

	return nil
}

func (s *subscriptionService) validateSubscription(sub *models.Subscription) error {
	if sub.ServiceName == "" {
		return errors.New("service name is required")
	}

	if len(sub.ServiceName) > 255 {
		return errors.New("service name too long")
	}

	if sub.Price <= 0 {
		return errors.New("price must be positive")
	}

	if sub.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}

	if sub.StartDate.IsZero() {
		return errors.New("start date is required")
	}

	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return errors.New("end date cannot be before start date")
	}

	return nil
}

func (s *subscriptionService) validateTotalCostRequest(req *models.TotalCostRequest) error {
	if req.StartPeriod == "" {
		return errors.New("start period is required")
	}

	if req.EndPeriod == "" {
		return errors.New("end period is required")
	}

	return nil
}

func (s *subscriptionService) toResponse(sub *models.Subscription) *models.SubscriptionResponse {
	response := &models.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format("01-2006"),
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}

	if sub.EndDate != nil {
		endDateStr := sub.EndDate.Format("01-2006")
		response.EndDate = &endDateStr
	}

	return response
}
