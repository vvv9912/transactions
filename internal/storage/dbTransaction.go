package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type PostgresTransaction struct {
	db *sqlx.DB
}

func NewTransStorage(db *sqlx.DB) *PostgresTransaction {
	return &PostgresTransaction{db: db}
}

// Добавить транзакцию
func (db *PostgresTransaction) AddTransaction(ctx context.Context, numTransaction string, userid int64, status int, timestamp time.Time) (int64, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	//
	var id int64

	row := conn.QueryRowxContext(
		ctx,
		"INSERT INTO transaction (num_transaction, user_id, status,CREATED_AT) VALUES ($1, $2,$3, $4) RETURNING id",
		numTransaction,
		userid,
		status,
		timestamp,
	)
	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
func (db *PostgresTransaction) SetTransactionById(ctx context.Context, id int64, status int) error {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, `UPDATE transaction SET status = $1 WHERE id = $2`, status, id)
	return err
}
func (db *PostgresTransaction) GetStatus(ctx context.Context, numTransaction string) (int, error) {

	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var status int
	if err := conn.GetContext(ctx, &status, `SELECT status FROM transaction WHERE num_transaction= $1`, numTransaction); err != nil {
		return 0, err
	}
	return status, err
}
