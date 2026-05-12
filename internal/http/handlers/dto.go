package handlers

import (
	"time"

	subscrdomain "github.com/IvanKuchsh-600/subtracker/internal/domain/subscription"
	subscrusecase "github.com/IvanKuchsh-600/subtracker/internal/usecase/subscription"
)

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price       int     `json:"price" binding:"required,min=1" example:"400"`
	UserID      string  `json:"user_id" binding:"required" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string  `json:"start_date" binding:"required" example:"07-2025"`
	EndDate     *string `json:"end_date" example:"12-2026"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty" example:"Netflix"`
	Price       *int    `json:"price,omitempty" example:"799"`
	StartDate   *string `json:"start_date,omitempty" example:"01-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2026"`
}

type TotalCostRequest struct {
	FromDate    string `form:"from_date"`
	ToDate      string `form:"to_date"`
	UserID      string `form:"user_id"`
	ServiceName string `form:"service_name"`
}

type SubscriptionResponse struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceName string    `json:"service_name" example:"Yandex Plus"`
	Price       int       `json:"price" example:"400"`
	UserID      string    `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date" example:"07-2025"`
	EndDate     *string   `json:"end_date" example:"12-2026"`
	CreatedAt   time.Time `json:"created_at" example:"2026-05-11T10:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2026-05-11T10:00:00Z"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost" example:"1200"`
}

func toSubscriptionResponse(sub *subscrdomain.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}
}

// для ошибок валидации (500)
type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

// для ошибок валидации (400)
type ValidationErrorResponse struct {
	Error string `json:"error" example:"service_name is required"`
}

// для ошибок 404
type NotFoundErrorResponse struct {
	Error string `json:"error" example:"subscription not found"`
}

func toCreateInput(req CreateSubscriptionRequest) subscrusecase.CreateInput {
	return subscrusecase.CreateInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
}

func toUpdateInput(req UpdateSubscriptionRequest) subscrusecase.UpdateInput {
	return subscrusecase.UpdateInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
}

func toTotalCostRequest(req TotalCostRequest) subscrusecase.TotalCostRequest {
	return subscrusecase.TotalCostRequest{
		FromDate:    req.FromDate,
		ToDate:      req.ToDate,
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
	}
}

func toSubscriptionResponseList(subs []subscrdomain.Subscription) []SubscriptionResponse {
	result := make([]SubscriptionResponse, len(subs))
	for i, sub := range subs {
		result[i] = toSubscriptionResponse(&sub)
	}
	return result
}
