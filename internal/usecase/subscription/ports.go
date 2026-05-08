package subscription

import (
	"context"

	subscrdomain "github.com/IvanKuchsh-600/subtracker/internal/domain/subscription"
)

type Repository interface {
	Create(ctx context.Context, sub *subscrdomain.Subscription) (*subscrdomain.Subscription, error)
	Update(ctx context.Context, id string, updates *subscrdomain.Subscription) (*subscrdomain.Subscription, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]subscrdomain.Subscription, int, error)
	GetByID(ctx context.Context, id string) (*subscrdomain.Subscription, error)
	GetTotalCost(ctx context.Context, fromDate, toDate, userID, serviceName string) (int, error)
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*subscrdomain.Subscription, error)
	GetByID(ctx context.Context, id string) (*subscrdomain.Subscription, error)
	Update(ctx context.Context, id string, input UpdateInput) (*subscrdomain.Subscription, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int) (*SubscriptionsListResponse, error)
}

type CreateInput struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price" `
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type UpdateInput struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type SubscriptionsListResponse struct {
	Subscriptions []subscrdomain.Subscription `json:"subscriptions"`
	Total         int                         `json:"total"`
}

type TotalCostRequest struct {
	FromDate    string `form:"from_date" binding:"required"`
	ToDate      string `form:"to_date" binding:"required"`
	UserID      string `form:"user_id"`
	ServiceName string `form:"service_name"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}
