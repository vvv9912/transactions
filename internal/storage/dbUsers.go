package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"transaction/internal/model"
)

type PostgresUsers struct {
	db *sqlx.DB
}

func NewUsersStorage(db *sqlx.DB) *PostgresUsers {
	return &PostgresUsers{db: db}
}

// Добавить пользователя
func (db *PostgresUsers) AddUsers(ctx context.Context, Users model.Users) (int64, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	//
	var id int64

	row := conn.QueryRowxContext(
		ctx,
		"INSERT INTO users (id_user, account) VALUES ($1, $2) RETURNING id",
		Users.UserID,
		Users.Account,
	)
	if err := row.Err(); err != nil {
		return 0, err
	}

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// // Получение баланса
// func (db *PostgresUsers) Account(ctx context.Context) (model.Users, error) {
//
// }
//
// func (db *PostgresUsers) AddAccountById(ctx context.Context, Users model.Users) error {
//
// }
// func (db *PostgresUsers) SubAccountById(ctx context.Context, Users model.Users) error {
//
// }
// Существует ли пользователь
func (db *PostgresUsers) CheckId(ctx context.Context, idUser int64) (int64, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var id int64
	if err := conn.GetContext(ctx, &id, `SELECT id FROM users WHERE ID_USER = $1`, idUser); err != nil {
		return 0, err
	}
	return id, err
}
func (db *PostgresUsers) GetAccount(ctx context.Context, idUser int64) (float64, error) {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var account float64
	if err := conn.GetContext(ctx, &account, `SELECT Account FROM users WHERE id= $1`, idUser); err != nil {
		return 0, err
	}
	return account, err
}

// Добавить баланс по акку
func (db *PostgresUsers) AddAccountById(ctx context.Context, id int64, account float64) error {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, `UPDATE users SET account = (account + $1) WHERE id = $2`, account, id)
	return err
}
func (db *PostgresUsers) SetAccountById(ctx context.Context, id int64, account float64) error {
	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, `UPDATE users SET account = $1 WHERE id = $2`, account, id)
	return err
}
