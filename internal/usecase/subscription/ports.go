package subscription

import (
	"context"

	subscrdomain "github.com/IvanKuchsh-600/subtracker/internal/domain/subscription"
)

type Repository interface {
	Create(ctx context.Context, sub *subscrdomain.Subscription) (*subscrdomain.Subscription, error)
	Update(ctx context.Context, id string, updates *subscrdomain.Subscription) (*subscrdomain.Subscription, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]subscrdomain.Subscription, error)
	GetByID(ctx context.Context, id string) (*subscrdomain.Subscription, error)
	GetTotalCost(ctx context.Context, fromDate, toDate, userID, serviceName string) (int, error)
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*subscrdomain.Subscription, error)
	GetByID(ctx context.Context, id string) (*subscrdomain.Subscription, error)
	Update(ctx context.Context, id string, input UpdateInput) (*subscrdomain.Subscription, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]subscrdomain.Subscription, error)
	GetTotalCost(ctx context.Context, req TotalCostRequest) (int, error)
}

type CreateInput struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     *string
}

type UpdateInput struct {
	ServiceName *string
	Price       *int
	StartDate   *string
	EndDate     *string
}

type TotalCostRequest struct {
	FromDate    string
	ToDate      string
	UserID      string
	ServiceName string
}
