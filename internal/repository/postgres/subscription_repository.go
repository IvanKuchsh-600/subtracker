package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	subscrerrors "github.com/IvanKuchsh-600/subtracker/internal/domain/errors"
	subscrdomain "github.com/IvanKuchsh-600/subtracker/internal/domain/subscription"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewSubscriptionRepository(connStr string, logger *slog.Logger) (*SubscriptionRepository, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Error("Failed to parse database config", "error", err)
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Error("Failed to create connection pool", "error", err)
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("Database ping failed", "error", err)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connected successfully")
	return &SubscriptionRepository{pool: pool, logger: logger}, nil
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *subscrdomain.Subscription) (*subscrdomain.Subscription, error) {
	query := `
        INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
    `

	sub.ID = uuid.New().String()
	row := r.pool.QueryRow(ctx, query, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)

	var created subscrdomain.Subscription
	var endDate *string

	err := row.Scan(&created.ID, &created.ServiceName, &created.Price, &created.UserID,
		&created.StartDate, &endDate, &created.CreatedAt, &created.UpdatedAt)

	if err != nil {
		r.logger.Error("Failed to create subscription", "error", err)
		return nil, err
	}

	if endDate != nil {
		created.EndDate = endDate
	}
	return &created, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*subscrdomain.Subscription, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        WHERE id = $1
    `

	var sub subscrdomain.Subscription
	var endDate *string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get subscription", "id", id, "error", err)
		return nil, err
	}

	sub.EndDate = endDate
	return &sub, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, id string, updates *subscrdomain.Subscription) (*subscrdomain.Subscription, error) {
	setClauses := []string{}
	args := []interface{}{}
	argCounter := 1

	if updates.ServiceName != "" {
		setClauses = append(setClauses, fmt.Sprintf("service_name = $%d", argCounter))
		args = append(args, updates.ServiceName)
		argCounter++
	}
	if updates.Price != 0 {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", argCounter))
		args = append(args, updates.Price)
		argCounter++
	}
	if updates.UserID != "" {
		setClauses = append(setClauses, fmt.Sprintf("user_id = $%d", argCounter))
		args = append(args, updates.UserID)
		argCounter++
	}
	if updates.StartDate != "" {
		setClauses = append(setClauses, fmt.Sprintf("start_date = $%d", argCounter))
		args = append(args, updates.StartDate)
		argCounter++
	}
	if updates.EndDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("end_date = $%d", argCounter))
		if *updates.EndDate == "" {
			args = append(args, nil)
		} else {
			args = append(args, updates.EndDate)
		}
		argCounter++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	query := fmt.Sprintf("UPDATE subscriptions SET %s WHERE id = $%d RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at",
		strings.Join(setClauses, ", "), argCounter)
	args = append(args, id)

	updated := subscrdomain.Subscription{}
	var endDate *string

	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&updated.ID, &updated.ServiceName, &updated.Price, &updated.UserID,
		&updated.StartDate, &endDate, &updated.CreatedAt, &updated.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, subscrerrors.ErrNotFound
		}
		r.logger.Error("Failed to update subscription", "id", id, "error", err)
		return nil, err
	}

	updated.EndDate = endDate
	return &updated, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		r.logger.Error("Failed to delete subscription", "id", id, "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return subscrerrors.ErrNotFound
	}
	return nil
}

func (r *SubscriptionRepository) List(ctx context.Context) ([]subscrdomain.Subscription, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        ORDER BY created_at DESC
    `)
	if err != nil {
		r.logger.Error("Failed to list subscriptions", "error", err)
		return nil, err
	}
	defer rows.Close()

	var subscriptions []subscrdomain.Subscription
	for rows.Next() {
		var sub subscrdomain.Subscription
		var endDate *string
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt); err != nil {
			r.logger.Error("Failed to scan subscription", "error", err)
			return nil, err
		}
		sub.EndDate = endDate
		subscriptions = append(subscriptions, sub)
	}
	return subscriptions, nil
}

func (r *SubscriptionRepository) GetTotalCost(ctx context.Context, fromDate, toDate, userID, serviceName string) (int, error) {
	query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE 
            (end_date IS NULL OR TO_DATE(end_date, 'MM-YYYY') >= TO_DATE($1, 'MM-YYYY'))
            AND TO_DATE(start_date, 'MM-YYYY') <= TO_DATE($2, 'MM-YYYY')
    `
	args := []interface{}{fromDate, toDate}
	argIndex := 3

	if userID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}
	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, serviceName)
	}

	var total int
	err := r.pool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		r.logger.Error("Failed to calculate total cost", "error", err)
		return 0, err
	}
	return total, nil
}
