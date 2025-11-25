package models

import (
	"time"

	"github.com/google/uuid"
)

type TotalCostRequest struct {
	UserID      *uuid.UUID `query:"user_id"`
	ServiceName *string    `query:"service_name"`
	StartPeriod string     `query:"start_period"`
	EndPeriod   string     `query:"end_period"`
}

type TotalCostResponse struct {
	TotalCost int    `json:"total_cost"`
	Currency  string `json:"currency"`
	Period    string `json:"period"`
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartDate   time.Time
	EndDate     time.Time
}

func (r *TotalCostRequest) ParseDates() (time.Time, time.Time, error) {
	startDate, err := time.Parse("01-2006", r.StartPeriod)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err := time.Parse("01-2006", r.EndPeriod)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startDate, endDate, nil
}
