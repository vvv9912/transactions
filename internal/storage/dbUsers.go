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
		"INSERT INTO users (id_users, account) VALUES ($1, $2) RETURNING id",
		Users.ID,
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

//// Получение баланса
//func (db *PostgresUsers) Account(ctx context.Context) (model.Users, error) {
//
//}
//
//// Добавить баланс по акку
//func (db *PostgresUsers) AddAccountById(ctx context.Context, Users model.Users) error {
//
//}
//func (db *PostgresUsers) SubAccountById(ctx context.Context, Users model.Users) error {
//
//}
