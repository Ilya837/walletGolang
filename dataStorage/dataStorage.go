package poolstorage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UUIDUndefined struct {
}

func (_ UUIDUndefined) Error() string {
	return "UUID undifined"
}

type UUIDExists struct {
}

func (_ UUIDExists) Error() string {
	return "UUID already exists"
}

type DBError struct {
}

func (_ DBError) Error() string {
	return "DB error"
}

type WalletStorage interface {
	Get(uuid string) (float64, error)
	Check(uuid string) (bool, error)
	ChangeBalance(sum float64, uuid string) (bool, error)
	CreateWallet(uuid string) error
}

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(host, port, user, password, dbName string) (Postgres, error) {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	pool, err := pgxpool.New(context.Background(), psqlconn)

	if err != nil {
		return Postgres{}, errors.New("database not created")
	}

	return Postgres{pool: pool}, nil
}

func (postgres Postgres) Get(uuid string) (float64, error) {
	var balance float64
	err := postgres.pool.QueryRow(context.Background(),
		"select balance from public.wallets where id=$1;",
		uuid).Scan(&balance)

	if err == pgx.ErrNoRows {
		log.Println("error in Get method: ", err)
		return 0, UUIDUndefined{}
	} else {

		if err != nil {
			log.Println("error in Get method: ", err)
			return 0, DBError{}
		}

		return balance, nil
	}
}

func (postgres Postgres) Check(uuid string) (bool, error) {

	var balance float64
	err := postgres.pool.QueryRow(context.Background(),
		"select balance from public.wallets where id=$1;",
		uuid).Scan(&balance)

	if err == pgx.ErrNoRows {
		return false, nil
	} else {

		if err != nil {
			log.Println("error in Check method: ", err)
			return false, DBError{}
		}

		return true, nil
	}

}

func (postgres Postgres) ChangeBalance(sum float64, uuid string) (bool, error) {

	cmdTag, err := postgres.pool.Exec(context.Background(),
		"UPDATE wallets SET balance = TRUNC( (balance + $1)::NUMERIC , 2) WHERE id = $2",
		sum, uuid)

	if err != nil {

		if err.(*pgconn.PgError).Code == "23514" {
			log.Println("balance too small for operation ")
			return false, nil

		} else {
			log.Println("error in ChangeBalance method: ", err)
			log.Println("error code: ")
			return false, DBError{}
		}
	}

	if cmdTag.RowsAffected() == 0 {
		log.Println("error in ChangeBalance method: ", err)
		return false, UUIDUndefined{}
	}

	return true, nil
}

func (postgres Postgres) CreateWallet(uuid string) error {

	_, err := postgres.pool.Exec(context.Background(),
		"INSERT INTO wallets (id, balance) VALUES ($1,0)", uuid)

	if err != nil {
		log.Println("error in CreateWallet method: ", err)
		return DBError{}
	}

	return nil

}
