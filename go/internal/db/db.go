package db

import (
	"context"
	"cu_coworking_book/go/internal/pb"
	"database/sql"
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
	ID             uint32  `db:"id"`
	Email          string  `db:"email"`
	FirstName      string  `db:"first_name"`
	SecondName     string  `db:"second_name"`
	Patronymic     string  `db:"patronymic"`
	HashedPassword string  `db:"password"`
	Role           pb.Role `db:"role"`
}

type SignUpResponse struct {
	ID   uint32
	Role pb.Role
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
		email VARCHAR UNIQUE NOT NULL,
		first_name VARCHAR NOT NULL,
		second_name VARCHAR NOT NULL,
		patronymic VARCHAR,
		password TEXT NOT NULL,
		role VARCHAR NOT NULL
	);`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("create table error: sql error: %w", err)
	}

	log.Println("Created users table")
	return nil
}

func CreateUser(user User) (SignUpResponse, error) {

	var res SignUpResponse
	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return res, fmt.Errorf("create user error: connect to db: %w", err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return res, fmt.Errorf("create user error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	users := []User{}
	role := pb.Role_USER
	if err = db.SelectContext(ctx, &users, `SELECT * FROM users`); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			role = pb.Role_ADMIN
		} else {
			return res, fmt.Errorf("create user error: select error: %w", err)
		}
	}
	user.Role = role
	res.Role = role

	queryInsert := `INSERT INTO users (email, first_name, second_name, patronymic, password, role) 
		VALUES (:email, :first_name, :second_name, :patronymic, :password, :role);`
	if _, err = db.ExecContext(ctx, queryInsert, &user); err != nil {
		return res, fmt.Errorf("create user error: insert error: %w", err)
	}

	var id uint32
	if err := db.GetContext(ctx, &id, `SELECT id FROM users WHERE email = $1`, user.Email); err != nil {
		return res, fmt.Errorf("create user error: get id error: %w", err)
	}
	res.ID = id

	return res, nil
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

	query := `SELECT * FROM users WHERE email = $1;`
	if err := db.GetContext(ctx, user, query, email); err != nil {
		return nil, fmt.Errorf("get user by email error: %w", err)
	}
	return user, nil
}

func UpdateUser(*pb.User) error {

	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return fmt.Errorf("get user by email error: connect to db: %w", err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return fmt.Errorf("get user by email error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `UPDATE users SET email := email, first_name := first_name, 
		second_name := second_name, mid_name := mid_name WHERE id := id`

}
