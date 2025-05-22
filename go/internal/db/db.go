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
	ID         uint32  `db:"id"`
	Email      string  `db:"email"`
	FirstName  string  `db:"first_name"`
	SecondName string  `db:"second_name"`
	Patronymic string  `db:"patronymic"`
	Password   string  `db:"password"`
	Role       pb.Role `db:"role"`
}

type SignUpResponse struct {
	ID   uint32
	Role pb.Role
}

type LoginUserData struct {
	ID       uint32 `db:"email"`
	Email    string `db:"email"`
	Password string `db:"password"`
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
		return fmt.Errorf("creating table: %w: %v", ErrConDB, err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return fmt.Errorf("creating table: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(50) UNIQUE NOT NULL,
		first_name VARCHAR(20) NOT NULL,
		second_name VARCHAR(20) NOT NULL,
		patronymic VARCHAR(20),
		password TEXT NOT NULL,
		role VARCHAR(10) NOT NULL
	);`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("creating table: create error: %w", err)
	}

	log.Println("Created users table")
	return nil
}

func DeleteTable() error {

	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return fmt.Errorf("deleting table: connection to db: %w: %v", ErrConDB, err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return fmt.Errorf("deleting table: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS users;`)
	if err != nil {
		return fmt.Errorf("deleting table: drop error: %w", err)
	}
	return nil
}

func CreateUser(user User) (SignUpResponse, error) {

	var res SignUpResponse
	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return res, fmt.Errorf("creating user: connection to db: %w", err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return res, fmt.Errorf("creating user: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var id uint32
	role := pb.Role_USER
	if err = db.GetContext(ctx, &id, `SELECT id FROM users LIMIT 1`); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			role = pb.Role_ADMIN
		} else {
			return res, fmt.Errorf("creating user: select error: %w", err)
		}
	}
	user.Role = role
	res.Role = role

	queryInsert := `INSERT INTO users (email, first_name, second_name, patronymic, password, role) 
		VALUES (:email, :first_name, :second_name, :patronymic, :password, :role);`
	if _, err = db.NamedExecContext(ctx, queryInsert, &user); err != nil {
		return res, fmt.Errorf("creating user: insert error: %w", err)
	}

	if err := db.GetContext(ctx, &id, `SELECT id FROM users WHERE email = $1`, user.Email); err != nil {
		return res, fmt.Errorf("creating user: get id error: %w", err)
	}
	res.ID = id

	return res, nil
}

func GetUserByEmail(email string) (LoginUserData, error) {

	var out LoginUserData
	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return out, fmt.Errorf("getting user by email: connection to db: %w", err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return out, fmt.Errorf("getting user by email: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `SELECT id, email, password FROM users WHERE email = $1;`
	if err := db.GetContext(ctx, &out, query, email); err != nil {
		return out, fmt.Errorf("getting user by email: select error: %w", err)
	}
	return out, nil
}

func UpdateUser(oldEmail string, user User) error {

	res, err := GetUserByEmail(oldEmail)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}
	user.ID = res.ID

	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return fmt.Errorf("getting user by email: connection to db: %w", err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return fmt.Errorf("getting user by email: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	query := `UPDATE users 
		SET email = :email, first_name := first_name, second_name := second_name, 
		patronymic := patronymic, password := password WHERE id := id`

	if _, err := db.NamedExecContext(ctx, query, &user); err != nil {
		return fmt.Errorf("updating user: update error: %w", err)
	}

	return nil
}

func DeleteUser(email string) error {

	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return fmt.Errorf("deleting user: connection to db: %w", err)
	}
	defer db.Close()

	timeout, err := getSqlTimeout()
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err := db.ExecContext(ctx, `DELETE FROM users WHERE email = $1`, email); err != nil {
		return fmt.Errorf("deleting user: db error: %w", err)
	}
	return nil
}
