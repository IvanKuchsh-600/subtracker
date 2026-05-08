package handlers

import (
	"time"

	subscrdomain "github.com/IvanKuchsh-600/subtracker/internal/domain/subscription"
)

type subscriptionMutationDTO struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type subscriptionDTO struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newSubscriptionDTO(subscription *subscrdomain.Subscription) subscriptionDTO {
	startDateStr := subscription.StartDate.Format("01-2006")

	var endDateStr *string
	if subscription.EndDate != nil {
		formatted := subscription.EndDate.Format("01-2006")
		endDateStr = &formatted
	}
	return subscriptionDTO{
		ID:          subscription.ID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID,
		StartDate:   startDateStr,
		EndDate:     endDateStr,
		CreatedAt:   subscription.CreatedAt,
		UpdatedAt:   subscription.UpdatedAt,
	}
}
