package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vnchk1/subscription-aggregator/internal/models"
	"github.com/vnchk1/subscription-aggregator/internal/service/mocks"
)

func TestSubscriptionService_CreateSubscription_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	req := &models.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       799,
		UserID:      userID,
		StartDate:   "01-2024",
	}

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, sub *models.Subscription) error {
			assert.Equal(t, "Netflix", sub.ServiceName)
			assert.Equal(t, 799, sub.Price)
			assert.Equal(t, userID, sub.UserID)
			assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), sub.StartDate)
			assert.Nil(t, sub.EndDate)
			return nil
		})

	result, err := service.CreateSubscription(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "Netflix", result.ServiceName)
	assert.Equal(t, 799, result.Price)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "01-2024", result.StartDate)
}

func TestSubscriptionService_CreateSubscription_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()

	testCases := []struct {
		name    string
		request *models.CreateSubscriptionRequest
		wantErr string
	}{
		{
			name: "empty service name",
			request: &models.CreateSubscriptionRequest{
				ServiceName: "",
				Price:       100,
				UserID:      uuid.New(),
				StartDate:   "01-2024",
			},
			wantErr: "service name is required",
		},
		{
			name: "zero price",
			request: &models.CreateSubscriptionRequest{
				ServiceName: "Netflix",
				Price:       0,
				UserID:      uuid.New(),
				StartDate:   "01-2024",
			},
			wantErr: "price must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

			result, err := service.CreateSubscription(ctx, tc.request)

			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestSubscriptionService_CreateSubscription_InvalidDateFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	req := &models.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       799,
		UserID:      uuid.New(),
		StartDate:   "invalid-date",
	}

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

	result, err := service.CreateSubscription(ctx, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start date format")
}

func TestSubscriptionService_CreateSubscription_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	req := &models.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       799,
		UserID:      uuid.New(),
		StartDate:   "01-2024",
	}

	expectedErr := errors.New("database error")
	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(expectedErr)

	result, err := service.CreateSubscription(ctx, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create subscription")
}

func TestSubscriptionService_CreateSubscription_ServiceNameTooLong(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()


	longName := string(make([]byte, 256))
	req := &models.CreateSubscriptionRequest{
		ServiceName: longName,
		Price:       799,
		UserID:      uuid.New(),
		StartDate:   "01-2024",
	}

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

	result, err := service.CreateSubscription(ctx, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service name too long")
}

func TestSubscriptionService_GetSubscription_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()
	expectedSub := &models.Subscription{
		ID:          subscriptionID,
		ServiceName: "Yandex Plus",
		Price:       399,
		UserID:      uuid.New(),
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.EXPECT().
		GetByID(ctx, subscriptionID).
		Return(expectedSub, nil)

	result, err := service.GetSubscription(ctx, subscriptionID)

	require.NoError(t, err)
	assert.Equal(t, expectedSub.ID, result.ID)
	assert.Equal(t, "Yandex Plus", result.ServiceName)
	assert.Equal(t, 399, result.Price)
}

func TestSubscriptionService_GetSubscription_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()

	mockRepo.EXPECT().
		GetByID(ctx, subscriptionID).
		Return(nil, models.ErrNotFound)

	result, err := service.GetSubscription(ctx, subscriptionID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, models.ErrNotFound)
}

func TestSubscriptionService_GetSubscription_EmptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()

	mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(0)

	result, err := service.GetSubscription(ctx, uuid.Nil)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "subscription ID is required")
}

func TestSubscriptionService_GetSubscription_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()

	expectedErr := errors.New("connection failed")
	mockRepo.EXPECT().
		GetByID(ctx, subscriptionID).
		Return(nil, expectedErr)

	result, err := service.GetSubscription(ctx, subscriptionID)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get subscription")
}

func TestSubscriptionService_UpdateSubscription_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()
	req := &models.UpdateSubscriptionRequest{
		ServiceName: "Yandex Plus Premium",
		Price:       599,
		StartDate:   "02-2024",
		EndDate:     stringPtr("12-2024"),
	}

	existingSub := &models.Subscription{
		ID:          subscriptionID,
		ServiceName: "Yandex Plus",
		Price:       399,
		UserID:      uuid.New(),
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.EXPECT().
		GetByID(ctx, subscriptionID).
		Return(existingSub, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, sub *models.Subscription) error {
			assert.Equal(t, subscriptionID, sub.ID)
			assert.Equal(t, "Yandex Plus Premium", sub.ServiceName)
			assert.Equal(t, 599, sub.Price)
			assert.Equal(t, time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), sub.StartDate)
			assert.Equal(t, time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC), *sub.EndDate)
			return nil
		})

	result, err := service.UpdateSubscription(ctx, subscriptionID, req)

	require.NoError(t, err)
	assert.Equal(t, "Yandex Plus Premium", result.ServiceName)
	assert.Equal(t, 599, result.Price)
	assert.Equal(t, "02-2024", result.StartDate)
	assert.Equal(t, "12-2024", *result.EndDate)
}

func TestSubscriptionService_UpdateSubscription_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()
	req := &models.UpdateSubscriptionRequest{
		ServiceName: "Updated Service",
		Price:       999,
		StartDate:   "01-2024",
	}

	mockRepo.EXPECT().
		GetByID(ctx, subscriptionID).
		Return(nil, models.ErrNotFound)

	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

	result, err := service.UpdateSubscription(ctx, subscriptionID, req)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, models.ErrNotFound)
}

func TestSubscriptionService_UpdateSubscription_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()

	testCases := []struct {
		name    string
		request *models.UpdateSubscriptionRequest
		wantErr string
	}{
		{
			name: "empty service name",
			request: &models.UpdateSubscriptionRequest{
				ServiceName: "",
				Price:       100,
				StartDate:   "01-2024",
			},
			wantErr: "service name is required",
		},
		{
			name: "zero price",
			request: &models.UpdateSubscriptionRequest{
				ServiceName: "Service",
				Price:       0,
				StartDate:   "01-2024",
			},
			wantErr: "price must be positive",
		},
		{
			name: "negative price",
			request: &models.UpdateSubscriptionRequest{
				ServiceName: "Service",
				Price:       -100,
				StartDate:   "01-2024",
			},
			wantErr: "price must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(0)
			mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

			result, err := service.UpdateSubscription(ctx, subscriptionID, tc.request)

			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestSubscriptionService_UpdateSubscription_EndDateBeforeStartDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()
	req := &models.UpdateSubscriptionRequest{
		ServiceName: "Service",
		Price:       100,
		StartDate:   "12-2024",
		EndDate:     stringPtr("01-2024"),
	}

	existingSub := &models.Subscription{
		ID:          subscriptionID,
		ServiceName: "Old Service",
		Price:       50,
		UserID:      uuid.New(),
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.EXPECT().
		GetByID(ctx, subscriptionID).
		Return(existingSub, nil)

	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

	result, err := service.UpdateSubscription(ctx, subscriptionID, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "end date cannot be before start date")
}

func TestSubscriptionService_UpdateSubscription_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()
	req := &models.UpdateSubscriptionRequest{
		ServiceName: "Updated Service",
		Price:       999,
		StartDate:   "01-2024",
	}

	existingSub := &models.Subscription{
		ID:          subscriptionID,
		ServiceName: "Old Service",
		Price:       50,
		UserID:      uuid.New(),
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.EXPECT().
		GetByID(ctx, subscriptionID).
		Return(existingSub, nil)

	expectedErr := errors.New("update failed")
	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(expectedErr)

	result, err := service.UpdateSubscription(ctx, subscriptionID, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update subscription")
}

func TestSubscriptionService_UpdateSubscription_EmptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	req := &models.UpdateSubscriptionRequest{
		ServiceName: "Service",
		Price:       100,
		StartDate:   "01-2024",
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(0)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

	result, err := service.UpdateSubscription(ctx, uuid.Nil, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "subscription ID is required")
}

func TestSubscriptionService_DeleteSubscription_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()

	mockRepo.EXPECT().
		Delete(ctx, subscriptionID).
		Return(nil)

	err := service.DeleteSubscription(ctx, subscriptionID)

	assert.NoError(t, err)
}

func TestSubscriptionService_DeleteSubscription_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()

	mockRepo.EXPECT().
		Delete(ctx, subscriptionID).
		Return(models.ErrNotFound)

	err := service.DeleteSubscription(ctx, subscriptionID)

	assert.ErrorIs(t, err, models.ErrNotFound)
}

func TestSubscriptionService_DeleteSubscription_EmptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()

	mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(0)

	err := service.DeleteSubscription(ctx, uuid.Nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "subscription ID is required")
}

func TestSubscriptionService_DeleteSubscription_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()

	expectedErr := errors.New("delete failed")
	mockRepo.EXPECT().
		Delete(ctx, subscriptionID).
		Return(expectedErr)

	err := service.DeleteSubscription(ctx, subscriptionID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete subscription")
}

func TestSubscriptionService_ListSubscriptions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")

	expectedSubs := []*models.Subscription{
		{
			ID:          uuid.New(),
			ServiceName: "Netflix",
			Price:       799,
			UserID:      userID,
			StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.EXPECT().
		List(ctx, &userID, 10, 0).
		Return(expectedSubs, 1, nil)

	result, err := service.ListSubscriptions(ctx, &userID, 1, 10)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Data, 1)
}

func TestSubscriptionService_ListSubscriptions_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()

	mockRepo.EXPECT().
		List(ctx, nil, 20, 0).
		Return([]*models.Subscription{}, 0, nil)

	result, err := service.ListSubscriptions(ctx, nil, 1, 20)

	require.NoError(t, err)
	assert.Equal(t, 0, result.Total)
	assert.Len(t, result.Data, 0)
}

func TestSubscriptionService_ListSubscriptions_WithPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()


	mockRepo.EXPECT().
		List(ctx, nil, 5, 5).
		Return([]*models.Subscription{}, 0, nil)

	result, err := service.ListSubscriptions(ctx, nil, 2, 5)

	require.NoError(t, err)
	assert.Equal(t, 0, result.Total)
}

func TestSubscriptionService_ListSubscriptions_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()

	expectedErr := errors.New("list failed")
	mockRepo.EXPECT().
		List(ctx, nil, 20, 0).
		Return(nil, 0, expectedErr)

	result, err := service.ListSubscriptions(ctx, nil, 1, 20)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list subscriptions")
}

func TestSubscriptionService_CalculateTotalCost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	req := &models.TotalCostRequest{
		StartPeriod: "01-2024",
		EndPeriod:   "12-2024",
	}

	expectedTotal := 2500

	mockRepo.EXPECT().
		GetTotalCost(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, filter *models.SubscriptionFilter) (int, error) {
			assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), filter.StartDate)
			assert.Equal(t, time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC), filter.EndDate)
			return expectedTotal, nil
		})

	result, err := service.CalculateTotalCost(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, expectedTotal, result.TotalCost)
	assert.Equal(t, "RUB", result.Currency)
	assert.Equal(t, "01-2024 - 12-2024", result.Period)
}

func TestSubscriptionService_CalculateTotalCost_WithUserFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	req := &models.TotalCostRequest{
		UserID:      &userID,
		StartPeriod: "01-2024",
		EndPeriod:   "12-2024",
	}

	expectedTotal := 1500
	mockRepo.EXPECT().
		GetTotalCost(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, filter *models.SubscriptionFilter) (int, error) {
			assert.Equal(t, &userID, filter.UserID)
			return expectedTotal, nil
		})

	result, err := service.CalculateTotalCost(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, expectedTotal, result.TotalCost)
}

func TestSubscriptionService_CalculateTotalCost_WithServiceFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	serviceName := "Netflix"
	req := &models.TotalCostRequest{
		ServiceName: &serviceName,
		StartPeriod: "01-2024",
		EndPeriod:   "12-2024",
	}

	expectedTotal := 799
	mockRepo.EXPECT().
		GetTotalCost(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, filter *models.SubscriptionFilter) (int, error) {
			assert.Equal(t, &serviceName, filter.ServiceName)
			return expectedTotal, nil
		})

	result, err := service.CalculateTotalCost(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, expectedTotal, result.TotalCost)
}

func TestSubscriptionService_CalculateTotalCost_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	req := &models.TotalCostRequest{
		StartPeriod: "01-2024",
		EndPeriod:   "12-2024",
	}

	mockRepo.EXPECT().
		GetTotalCost(ctx, gomock.Any()).
		Return(0, nil)

	result, err := service.CalculateTotalCost(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 0, result.TotalCost)
}

func TestSubscriptionService_CalculateTotalCost_InvalidPeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	req := &models.TotalCostRequest{
		StartPeriod: "12-2024",
		EndPeriod:   "01-2024",
	}

	mockRepo.EXPECT().GetTotalCost(gomock.Any(), gomock.Any()).Times(0)

	result, err := service.CalculateTotalCost(ctx, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "end period cannot be before start period")
}

func TestSubscriptionService_CalculateTotalCost_MissingRequiredParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()

	testCases := []struct {
		name    string
		request *models.TotalCostRequest
		wantErr string
	}{
		{
			name: "missing start period",
			request: &models.TotalCostRequest{
				EndPeriod: "12-2024",
			},
			wantErr: "start period is required",
		},
		{
			name: "missing end period",
			request: &models.TotalCostRequest{
				StartPeriod: "01-2024",
			},
			wantErr: "end period is required",
		},
		{
			name:    "both periods missing",
			request: &models.TotalCostRequest{},
			wantErr: "start period is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().GetTotalCost(gomock.Any(), gomock.Any()).Times(0)

			result, err := service.CalculateTotalCost(ctx, tc.request)

			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestSubscriptionService_CalculateTotalCost_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	req := &models.TotalCostRequest{
		StartPeriod: "01-2024",
		EndPeriod:   "12-2024",
	}

	expectedErr := errors.New("calculation failed")
	mockRepo.EXPECT().
		GetTotalCost(ctx, gomock.Any()).
		Return(0, expectedErr)

	result, err := service.CalculateTotalCost(ctx, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to calculate total cost")
}

func TestSubscriptionService_CalculateTotalCost_InvalidDateFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	req := &models.TotalCostRequest{
		StartPeriod: "invalid-date",
		EndPeriod:   "12-2024",
	}

	mockRepo.EXPECT().GetTotalCost(gomock.Any(), gomock.Any()).Times(0)

	result, err := service.CalculateTotalCost(ctx, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid date format")
}

func TestSubscriptionService_DateValidation_EndDateBeforeStartDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubscriptionRepository(ctrl)
	service := NewSubscriptionService(mockRepo)

	ctx := context.Background()
	subscriptionID := uuid.New()

	existingSub := &models.Subscription{
		ID:          subscriptionID,
		ServiceName: "Old Service",
		Price:       50,
		UserID:      uuid.New(),
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	req := &models.UpdateSubscriptionRequest{
		ServiceName: "Service",
		Price:       100,
		StartDate:   "12-2024",
		EndDate:     stringPtr("01-2024"),
	}

	mockRepo.EXPECT().
		GetByID(ctx, subscriptionID).
		Return(existingSub, nil)

	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

	result, err := service.UpdateSubscription(ctx, subscriptionID, req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "end date cannot be before start date")
}

func stringPtr(s string) *string {
	return &s
}
