package handler

import (
	"net/http"
	"strconv"

	"github.com/vnchk1/subscription-aggregator/internal/models"
	"github.com/vnchk1/subscription-aggregator/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

// @Router /subscriptions [post].
func (h *SubscriptionHandler) CreateSubscription(c echo.Context) error {
	var req models.CreateSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	subscription, err := h.service.CreateSubscription(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Failed to create subscription",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, subscription)
}

// @Router /subscriptions/{id} [get].
func (h *SubscriptionHandler) GetSubscription(c echo.Context) error {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid subscription ID",
			Message: "Subscription ID must be a valid UUID",
		})
	}

	subscription, err := h.service.GetSubscription(c.Request().Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, subscription)
}

// @Router /subscriptions/{id} [put].
func (h *SubscriptionHandler) UpdateSubscription(c echo.Context) error {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid subscription ID",
			Message: "Subscription ID must be a valid UUID",
		})
	}

	var req models.UpdateSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	subscription, err := h.service.UpdateSubscription(c.Request().Context(), id, &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, subscription)
}

// @Router /subscriptions/{id} [delete].
func (h *SubscriptionHandler) DeleteSubscription(c echo.Context) error {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid subscription ID",
			Message: "Subscription ID must be a valid UUID",
		})
	}

	if err := h.service.DeleteSubscription(c.Request().Context(), id); err != nil {
		return h.handleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Router /subscriptions [get].
func (h *SubscriptionHandler) ListSubscriptions(c echo.Context) error {
	// Парсинг параметров запроса
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var userID *uuid.UUID

	if userIDStr := c.QueryParam("user_id"); userIDStr != "" {
		id, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid user ID",
				Message: "User ID must be a valid UUID",
			})
		}

		userID = &id
	}

	response, err := h.service.ListSubscriptions(c.Request().Context(), userID, page, limit)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

// @Router /subscriptions/total-cost [get].
func (h *SubscriptionHandler) CalculateTotalCost(c echo.Context) error {
	var req models.TotalCostRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid query parameters",
			Message: err.Error(),
		})
	}

	if userIDStr := c.QueryParam("user_id"); userIDStr != "" {
		id, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid user ID",
				Message: "User ID must be a valid UUID",
			})
		}

		req.UserID = &id
	}

	response, err := h.service.CalculateTotalCost(c.Request().Context(), &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *SubscriptionHandler) handleError(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	status := http.StatusInternalServerError
	errorMsg := err.Error()

	switch {
	case err.Error() == "subscription not found":
		status = http.StatusNotFound
	case err.Error() == "invalid subscription data":
		status = http.StatusBadRequest
	default:
		switch {
		case contains(err.Error(), "validation failed"):
			status = http.StatusBadRequest
		case contains(err.Error(), "not found"):
			status = http.StatusNotFound
		case contains(err.Error(), "invalid") || contains(err.Error(), "required"):
			status = http.StatusBadRequest
		}
	}

	return c.JSON(status, models.ErrorResponse{
		Error:   http.StatusText(status),
		Message: errorMsg,
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && len(substr) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
