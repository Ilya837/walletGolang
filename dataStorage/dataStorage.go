package datastorage

import (
	"database/sql"
	"errors"
	"fmt"
	"walletGolang/server"

	_ "github.com/lib/pq"
)

type WalletStorage interface {
	Get(uuid string) (float32, error)
	Check(uuid string) (bool, error)
	ChangeBalance(sum float32, uuid string) error
	CreateWallet(uuid string) error
}

type Postgres struct {
	data *sql.DB
}

func NewPostgres(host, port, user, password, dbname string) (Postgres, error) {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		return Postgres{}, err
	}

	if err := db.Ping(); err != nil {
		return Postgres{}, err
	}

	return Postgres{data: db}, nil
}

func (postgres Postgres) Get(uuid string) (float32, error) {
	err := postgres.data.Ping()

	if err != nil {
		return 0, err
	}

	rows, err := postgres.data.Query("select balance from public.wallets where id=$1", uuid)

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if rows.Next() {
		var balance float32
		err = rows.Scan(&balance)
		if err != nil {
			return 0, err
		}

		return balance, nil

	} else {
		return 0, server.UUIDUndefined{}
	}
}

func (postgres Postgres) Check(uuid string) (bool, error) {

	err := postgres.data.Ping()

	if err != nil {
		return false, err
	}

	rows, err := postgres.data.Query("select * from public.wallets where id=$1", uuid)

	if err != nil {
		return false, err
	}

	defer rows.Close()

	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}

}

func (postgres Postgres) ChangeBalance(sum float32, uuid string) error {
	err := postgres.data.Ping()

	if err != nil {
		return err
	}

	tx, err := postgres.data.Begin()

	defer tx.Rollback()

	exsist, err := postgres.Check(uuid)

	if !exsist {
		return server.UUIDUndefined{}
	}

	result, err := postgres.data.Exec("UPDATE wallets SET balance = balance + $1 WHERE id = $2", sum, uuid)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows != 1 {
		return errors.New("not enough balance")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (postgres Postgres) CreateWallet(uuid string) error {
	err := postgres.data.Ping()

	if err != nil {
		return err
	}

	exsist, err := postgres.Check(uuid)

	if exsist {
		return server.UUIDExists{}
	}

	result, err := postgres.data.Exec("INSERT INTO wallets (id, balance) VALUES ($1,0)", uuid)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 1 {
		return nil
	} else {
		return errors.New("strange insert behavior")
	}
}
