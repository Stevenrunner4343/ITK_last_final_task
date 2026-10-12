package repository

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"gw-currency-wallet/internal/storages/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepo struct {
	db *pgxpool.Pool
}

var ErrNotFound = errors.New("Юзер не найден")
var ErrWrongPassword = errors.New("Не верный Пароль")

func NewWalletRepo(db *pgxpool.Pool) *WalletRepo {
	return &WalletRepo{db: db}
}

func (r *WalletRepo) SendMoney(ctx context.Context, id int, operation string, amount int) error {
	_, err := r.db.Exec(ctx, "INSERT INTO operations (user_id, operation, amount) VALUES ($1, $2, $3)", id, operation, amount)
	if err != nil {
		return err
	}
	return nil

}

func (r *WalletRepo) FindUser(ctx context.Context, UserId int) (int, error) {
	var count int
	err := r.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM operations WHERE user_id = $1",
		UserId).Scan(&count)
	if err != nil {

		return 0, err

	}
	if count == 0 {
		return 0, ErrNotFound
	}

	var balance int

	err = r.db.QueryRow(context.Background(),
		"SELECT COALESCE(SUM(CASE WHEN operation = 'deposit' THEN amount WHEN operation = 'withdraw' THEN -amount END), 0) FROM operations WHERE user_id = $1",
		UserId).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil

}
func (r *WalletRepo) MakeNewUser(ctx context.Context, a model.Account) (int, error) {
	var username string
	var userID int

	hash := sha256.Sum256([]byte(a.Password))
	password := fmt.Sprintf("%x", hash)

	err := r.db.QueryRow(ctx, "select username FROM accounts where username = $1 AND email = $2",
		a.Username, a.Email).Scan(&username)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.db.QueryRow(ctx, "INSERT INTO accounts (username,password,email) VALUES ($1,$2,$3) RETURNING id",
				a.Username, password, a.Email).Scan(&userID)
			return userID, nil
		} else {
			return 0, ErrNotFound
		}
	} else {
		return 0, nil
	}

}

func (r *WalletRepo) CheckUser(ctx context.Context, acc model.Account) error {
	var storedHash string
	err := r.db.QueryRow(ctx,
		"SELECT password FROM accounts WHERE username = $1",
		acc.Username).Scan(&storedHash)

	if err == pgx.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	inputHash := fmt.Sprintf("%x", sha256.Sum256([]byte(acc.Password)))
	if storedHash != inputHash {
		return ErrWrongPassword
	}

	return nil
}
