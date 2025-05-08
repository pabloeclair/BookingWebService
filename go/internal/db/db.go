package db

import (
	"context"
	"cu_coworking_book/go/internal/pb"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var ErrConDB = errors.New("connection to the database failed")

type User struct {
	Email      string `db:"email"`
	FirstName  string `db:"first_name"`
	SecondName string `db:"second_name"`
	MidName    string `db:"mid_name"`
	Password   string `db:"password"`
}

func getSqlTimeout() (time.Duration, error) {

	var timeout time.Duration
	if os.Getenv("SQL_TIMEOUT") == "" {
		timeout = time.Second * 7
	} else {
		t, err := strconv.Atoi(os.Getenv("SQL_TIMEOUT"))
		if err != nil {
			return 0, errors.New("invalid environment variable: SQL_TIMEOUT must be int")
		}
		timeout = time.Second * time.Duration(t)
	}
	return timeout, nil
}

func CreateTable() error {

	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return fmt.Errorf("create table error: %w: %v", ErrConDB, err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return fmt.Errorf("create table error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR NOT NULL,
		first_name VARCHAR NOT NULL,
		second_name VARCHAR NOT NULL,
		mid_name VARCHAR,
		password TEXT NOT NULL
	);`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("create table error: sql error: %w", err)
	}

	log.Println("Created users table")
	return nil
}

func CreateUser(user User) (uint32, error) {

	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return 0, fmt.Errorf("create user error: connect to db: %w", err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return 0, fmt.Errorf("create user error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	queryInsert := `INSERT INTO users (email, first_name, second_name, mid_name, password) 
	VALUES (:email, :first_name, :second_name, :mid_name, :password);`

	_, err = db.NamedExecContext(ctx, queryInsert, user)
	if err != nil {
		return 0, fmt.Errorf("create user error: insert data: %w", err)
	}

	queryGet := `SELECT id FROM USERS ORDER BY id DESC LIMIT 1`
	var id uint32
	if err = db.GetContext(ctx, &id, queryGet); err != nil {
		return 0, fmt.Errorf("create user error: select id: %w", err)
	}
	return id, nil
}

func GetUserByEmail(email string) (*pb.User, error) {

	var user *pb.User
	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return nil, fmt.Errorf("get user by email error: connect to db: %w", err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return nil, fmt.Errorf("get user by email error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `SELECT * FROM users WHERE email=$1;`
	if err := db.GetContext(ctx, user, query); err != nil {
		return nil, fmt.Errorf("get user by email error: %w", err)
	}
	return user, nil
}
