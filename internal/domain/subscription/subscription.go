package subscription

import "time"

type Subscription struct {
	ID          string
	ServiceName string
	Price       int
	UserID      string
	StartDate   string
	EndDate     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
