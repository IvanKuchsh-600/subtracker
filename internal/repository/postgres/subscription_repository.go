package postgres

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

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
		logger.Error("Failed to connect to database", "error", err)
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("Database ping failed", "error", err)
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")

	return &SubscriptionRepository{pool: pool}, nil
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *subscrdomain.Subscription) (*subscrdomain.Subscription, error) {
	query := `
        INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
    `

	sub.ID = uuid.New().String()

	// var endDate time.Time
	// if sub.EndDate != nil {
	// 	endDate = *sub.EndDate
	// }

	row := r.pool.QueryRow(ctx, query, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)

	var (
		created *subscrdomain.Subscription
		endDate *time.Time
	)

	err := row.Scan(&created.ID, &created.ServiceName, &created.Price, &created.UserID,
		&created.StartDate, &endDate, &created.CreatedAt, &created.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if endDate != nil {
		created.EndDate = endDate
	}

	return created, err
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*subscrdomain.Subscription, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        WHERE id = $1
    `

	var sub subscrdomain.Subscription
	var endDate *time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
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

	if !updates.StartDate.IsZero() {
		setClauses = append(setClauses, fmt.Sprintf("start_date = $%d", argCounter))
		args = append(args, updates.StartDate.Format("01-2006"))
		argCounter++
	}

	if updates.EndDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("end_date = $%d", argCounter))
		if updates.EndDate.IsZero() {
			args = append(args, nil)
		} else {
			args = append(args, updates.EndDate.Format("01-2006"))
		}
		argCounter++
	}

	if len(setClauses) == 0 {
		return nil, nil
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf("UPDATE subscriptions SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), argCounter)
	args = append(args, id)

	row := r.pool.QueryRow(ctx, query, args...)

	var (
		updated *subscrdomain.Subscription
		endDate *time.Time
	)

	err := row.Scan(&updated.ID, &updated.ServiceName, &updated.Price, &updated.UserID,
		&updated.StartDate, &updated.EndDate, &updated.EndDate, &endDate)

	if err != nil {
		return nil, err
	}

	if endDate != nil {
		updated.EndDate = endDate
	}

	return updated, err
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM subscriptions WHERE id = $1"
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *SubscriptionRepository) List(ctx context.Context, limit, offset int) ([]subscrdomain.Subscription, int, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subscriptions []subscrdomain.Subscription
	for rows.Next() {
		var sub subscrdomain.Subscription
		var endDate *time.Time
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &endDate, &sub.CreatedAt, &sub.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		sub.EndDate = endDate
		subscriptions = append(subscriptions, sub)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM subscriptions"
	err = r.pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func (r *SubscriptionRepository) GetTotalCost(ctx context.Context, fromDate, toDate, userID, serviceName string) (int, error) {
	query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE (end_date IS NULL OR end_date >= $1)
        AND start_date <= $2
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
	return total, err
}
