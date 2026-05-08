package subscription

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	subscrdomain "github.com/IvanKuchsh-600/subtracker/internal/domain/subscription"
)

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*subscrdomain.Subscription, error) {

	inp, err := validateCreateInput(input)
	if err != nil {
		s.logger.Warn("Validation failed", "error", err)
		return nil, err
	}

	sub := &subscrdomain.Subscription{
		ServiceName: inp.ServiceName,
		Price:       inp.Price,
		UserID:      inp.UserID,
	}
	startDate, _ := time.Parse("01-2006", inp.StartDate)
	endDate, _ := time.Parse("01-2006", inp.StartDate)
	sub.StartDate = startDate
	sub.EndDate = &endDate

	now := time.Now().UTC()
	sub.CreatedAt = now
	sub.UpdatedAt = now

	createdSubscr, err := s.repo.Create(ctx, sub)
	if err != nil {
		s.logger.Error("Failed to create subscription", "error", err)
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	s.logger.Info("Subscription created", "id", sub.ID, "user_id", sub.UserID)
	return createdSubscr, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*subscrdomain.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get subscription", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	if sub == nil {
		return nil, nil
	}
	return sub, nil
}

func (s *Service) Update(ctx context.Context, id string, input UpdateInput) (*subscrdomain.Subscription, error) {
	updates, err := validateUpdateInput(input)
	if err != nil {
		s.logger.Warn("Validation failed", "error", err)
		return nil, err
	}

	sub := &subscrdomain.Subscription{}

	if updates.ServiceName != nil {
		sub.ServiceName = *updates.ServiceName
	}
	if updates.Price != nil {
		sub.Price = *updates.Price
	}
	if updates.StartDate != nil {
		parsed, _ := time.Parse("01-2006", *updates.StartDate)
		sub.StartDate = parsed
	}
	if updates.EndDate != nil {
		if *updates.EndDate == "" {
			sub.EndDate = nil
		} else {
			parsed, _ := time.Parse("01-2006", *updates.EndDate)
			sub.EndDate = &parsed
		}
	}

	now := time.Now().UTC()
	sub.UpdatedAt = now

	updated, err := s.repo.Update(ctx, id, sub)
	if err != nil {
		s.logger.Error("Failed to update subscription", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	s.logger.Info("Subscription updated", "id", id)
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete subscription", "id", id, "error", err)
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	s.logger.Info("Subscription deleted", "id", id)

	return nil
}

func (s *Service) List(ctx context.Context, page, pageSize int) (*SubscriptionsListResponse, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	subscriptions, total, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		s.logger.Error("Failed to list subscriptions", "error", err)
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return &SubscriptionsListResponse{
		Subscriptions: subscriptions,
		Total:         total,
	}, nil
}

func (s *Service) GetTotalCost(ctx context.Context, req TotalCostRequest) (int, error) {
	total, err := s.repo.GetTotalCost(ctx, req.FromDate, req.ToDate, req.UserID, req.ServiceName)
	if err != nil {
		s.logger.Error("Failed to calculate total cost", "error", err)
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}
	return total, nil
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	if input.ServiceName == "" {
		return CreateInput{}, fmt.Errorf("%w: service_name is required", ErrInvalidInput)
	}
	if input.Price <= 0 {
		return CreateInput{}, fmt.Errorf("%w: price must be greater than 0", ErrInvalidInput)
	}
	if input.UserID == "" {
		return CreateInput{}, fmt.Errorf("%w: user_id is required", ErrInvalidInput)
	}
	if input.StartDate == "" {
		return CreateInput{}, fmt.Errorf("%w: start_date is required", ErrInvalidInput)
	}

	_, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		return CreateInput{}, fmt.Errorf("%w: start_date must be in MM-YYYY format", ErrInvalidInput)
	}

	if input.EndDate != nil && *input.EndDate != "" {
		if _, err := time.Parse("01-2006", *input.EndDate); err != nil {
			return CreateInput{}, fmt.Errorf("%w: end_date must be in MM-YYYY format", ErrInvalidInput)
		}
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	if input.ServiceName != nil && *input.ServiceName == "" {
		return UpdateInput{}, fmt.Errorf("%w: service_name cannot be empty", ErrInvalidInput)
	}

	if input.Price != nil && *input.Price <= 0 {
		return UpdateInput{}, fmt.Errorf("%w: price must be greater than 0", ErrInvalidInput)
	}

	if input.StartDate != nil {
		if *input.StartDate == "" {
			return UpdateInput{}, fmt.Errorf("%w: start_date cannot be empty if provided", ErrInvalidInput)
		}
		if _, err := time.Parse("01-2006", *input.StartDate); err != nil {
			return UpdateInput{}, fmt.Errorf("%w: start_date must be in MM-YYYY format", ErrInvalidInput)
		}
	}

	if input.EndDate != nil {
		if *input.EndDate != "" {
			if _, err := time.Parse("01-2006", *input.EndDate); err != nil {
				return UpdateInput{}, fmt.Errorf("%w: end_date must be in MM-YYYY format", ErrInvalidInput)
			}
		}
	}

	return input, nil
}
