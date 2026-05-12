package subscription

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	subscrerrors "github.com/IvanKuchsh-600/subtracker/internal/domain/errors"
	subscrdomain "github.com/IvanKuchsh-600/subtracker/internal/domain/subscription"
	"github.com/google/uuid"
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
		StartDate:   inp.StartDate,
		EndDate:     inp.EndDate,
	}
	if inp.EndDate != nil && *inp.EndDate == "" {
		sub.EndDate = nil
	}

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
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid uuid format", subscrerrors.ErrInvalidInput)
	}

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get subscription", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	if sub == nil {
		return nil, subscrerrors.ErrNotFound
	}

	return sub, nil
}

func (s *Service) Update(ctx context.Context, id string, input UpdateInput) (*subscrdomain.Subscription, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid uuid format", subscrerrors.ErrInvalidInput)
	}

	inp, err := validateUpdateInput(input)
	if err != nil {
		s.logger.Warn("Validation failed", "error", err)
		return nil, err
	}

	sub := &subscrdomain.Subscription{}

	if inp.ServiceName != nil {
		sub.ServiceName = *inp.ServiceName
	}
	if inp.Price != nil {
		sub.Price = *inp.Price
	}
	if inp.StartDate != nil {
		sub.StartDate = *inp.StartDate
	}
	if inp.EndDate != nil {
		sub.EndDate = inp.EndDate
	}

	updated, err := s.repo.Update(ctx, id, sub)
	if err != nil {
		if errors.Is(err, subscrerrors.ErrNotFound) {
			s.logger.Warn("Subscription not found", "id", id)
			return nil, subscrerrors.ErrNotFound
		}
		s.logger.Error("Failed to update subscription", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	s.logger.Info("Subscription updated", "id", id)
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: invalid uuid format", subscrerrors.ErrInvalidInput)
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, subscrerrors.ErrNotFound) {
			s.logger.Warn("Subscription not found for deletion", "id", id)
			return subscrerrors.ErrNotFound
		}
		s.logger.Error("Failed to delete subscription", "id", id, "error", err)
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	s.logger.Info("Subscription deleted", "id", id)
	return nil
}

func (s *Service) List(ctx context.Context) ([]subscrdomain.Subscription, error) {
	subscriptions, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	return subscriptions, nil
}

func (s *Service) GetTotalCost(ctx context.Context, req TotalCostRequest) (int, error) {
	if req.FromDate == "" {
		return 0, fmt.Errorf("%w: from_date is required", subscrerrors.ErrInvalidInput)
	}
	if req.ToDate == "" {
		return 0, fmt.Errorf("%w: to_date is required", subscrerrors.ErrInvalidInput)
	}

	_, err := time.Parse("01-2006", req.FromDate)
	if err != nil {
		return 0, fmt.Errorf("%w: from_date must be in MM-YYYY format", subscrerrors.ErrInvalidInput)
	}

	_, err = time.Parse("01-2006", req.ToDate)
	if err != nil {
		return 0, fmt.Errorf("%w: to_date must be in MM-YYYY format", subscrerrors.ErrInvalidInput)
	}

	from, _ := time.Parse("01-2006", req.FromDate)
	to, _ := time.Parse("01-2006", req.ToDate)
	if from.After(to) {
		return 0, fmt.Errorf("%w: from_date must be before or equal to to_date", subscrerrors.ErrInvalidInput)
	}

	total, err := s.repo.GetTotalCost(ctx, req.FromDate, req.ToDate, req.UserID, req.ServiceName)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	s.logger.Info("Total cost calculated",
		"from_date", req.FromDate,
		"to_date", req.ToDate,
		"user_id", req.UserID,
		"service_name", req.ServiceName,
		"total", total,
	)

	return total, nil
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	if input.ServiceName == "" {
		return CreateInput{}, fmt.Errorf("%w: service_name is required", subscrerrors.ErrInvalidInput)
	}
	if input.Price <= 0 {
		return CreateInput{}, fmt.Errorf("%w: price must be greater than 0", subscrerrors.ErrInvalidInput)
	}
	if input.UserID == "" {
		return CreateInput{}, fmt.Errorf("%w: user_id is required", subscrerrors.ErrInvalidInput)
	}
	if input.StartDate == "" {
		return CreateInput{}, fmt.Errorf("%w: start_date is required", subscrerrors.ErrInvalidInput)
	}

	startDate, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		return CreateInput{}, fmt.Errorf("%w: start_date must be in MM-YYYY format", subscrerrors.ErrInvalidInput)
	}

	if input.EndDate != nil && *input.EndDate != "" {
		endDate, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			return CreateInput{}, fmt.Errorf("%w: end_date must be in MM-YYYY format", subscrerrors.ErrInvalidInput)
		}
		if endDate.Before(startDate) {
			return CreateInput{}, fmt.Errorf("%w: end_date cannot be before start_date", subscrerrors.ErrInvalidInput)
		}
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	if input.ServiceName != nil && *input.ServiceName == "" {
		return UpdateInput{}, fmt.Errorf("%w: service_name cannot be empty", subscrerrors.ErrInvalidInput)
	}

	if input.Price != nil && *input.Price <= 0 {
		return UpdateInput{}, fmt.Errorf("%w: price must be greater than 0", subscrerrors.ErrInvalidInput)
	}

	if input.StartDate != nil {
		if *input.StartDate == "" {
			return UpdateInput{}, fmt.Errorf("%w: start_date cannot be empty if provided", subscrerrors.ErrInvalidInput)
		}

		_, err := time.Parse("01-2006", *input.StartDate)
		if err != nil {
			return UpdateInput{}, fmt.Errorf("%w: start_date must be in MM-YYYY format", subscrerrors.ErrInvalidInput)
		}
	}

	if input.EndDate != nil && *input.EndDate != "" {
		if _, err := time.Parse("01-2006", *input.EndDate); err != nil {
			return UpdateInput{}, fmt.Errorf("%w: end_date must be in MM-YYYY format", subscrerrors.ErrInvalidInput)
		}
	}

	return input, nil
}
