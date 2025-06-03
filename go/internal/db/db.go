package db

import (
	"context"
	"cu_coworking_book/go/internal/pb"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	ErrConDB      = errors.New("connection to the database failed")
	ErrBadRequest = errors.New("409 error")
	ErrNotFound   = errors.New("404 error")
)

type User struct {
	ID         uint32 `db:"id"`
	Email      string `db:"email"`
	FirstName  string `db:"first_name"`
	SecondName string `db:"second_name"`
	Patronymic string `db:"patronymic"`
	Password   string `db:"password"`
	Role       string `db:"role"`
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

func connectToDb() (context.Context, func(), *sqlx.DB, error) {

	db, err := sqlx.Connect("pgx", os.Getenv("DSN"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating user: connection to db: %w", err)
	}

	timeout, err := getSqlTimeout()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating user: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return ctx, cancel, db, nil
}

func CreateTable() error {

	ctx, cancel, db, err := connectToDb()
	if err != nil {
		return err
	}
	defer cancel()
	defer db.Close()

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

	return nil
}

func DeleteTable() error {

	ctx, cancel, db, err := connectToDb()
	if err != nil {
		return err
	}
	defer cancel()
	defer db.Close()

	_, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS users;`)
	if err != nil {
		return fmt.Errorf("deleting table: drop error: %w", err)
	}
	return nil
}

func CreateUser(user User) (uint32, error) {

	if _, err := GetUserByEmail(user.Email); !errors.Is(err, ErrNotFound) {
		if err != nil {
			return 0, fmt.Errorf("creating user: ошибка проверки почты: %w", err)
		}
		return 0, fmt.Errorf("%w: пользователь с почтой %s уже существует", ErrBadRequest, user.Email)
	}

	ctx, cancel, db, err := connectToDb()
	if err != nil {
		return 0, err
	}
	defer cancel()
	defer db.Close()

	var id uint32
	if err = db.GetContext(ctx, &id, `SELECT id FROM users LIMIT 1`); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user.Role = pb.Role_MAIN_ADMIN.String()
		} else {
			return 0, fmt.Errorf("creating user: select error: %w", err)
		}
	}

	queryInsert := `INSERT INTO users (email, first_name, second_name, patronymic, password, role) 
		VALUES (:email, :first_name, :second_name, :patronymic, :password, :role);`
	if _, err = db.NamedExecContext(ctx, queryInsert, &user); err != nil {
		return 0, fmt.Errorf("creating user: insert error: %w", err)
	}

	var idRes uint32
	querySelect := `SELECT id FROM users WHERE email = $1;`
	if err := db.GetContext(ctx, &idRes, querySelect, user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, idRes)
		}
		return 0, fmt.Errorf("getting user by id: select error: %w", err)
	}

	return idRes, nil
}

func GetUserByEmail(email string) (User, error) {

	var res User
	ctx, cancel, db, err := connectToDb()
	if err != nil {
		return res, err
	}
	defer cancel()
	defer db.Close()

	query := `SELECT * FROM users WHERE email = $1;`
	if err := db.GetContext(ctx, &res, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, fmt.Errorf("%w: пользователь с почтой %s не существует", ErrNotFound, email)
		}
		return res, fmt.Errorf("getting user by email: select error: %w", err)
	}
	return res, nil
}

func GetUserById(id uint32) (User, error) {

	var res User
	ctx, cancel, db, err := connectToDb()
	if err != nil {
		return res, err
	}
	defer cancel()
	defer db.Close()

	query := `SELECT * FROM users WHERE id = $1;`
	if err := db.GetContext(ctx, &res, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, id)
		}
		return res, fmt.Errorf("getting user by id: select error: %w", err)
	}
	return res, nil
}

func GetUserByKey(sortBy *pb.By, sortValue string) ([]User, error) {

	var res []User
	ctx, cancel, db, err := connectToDb()
	if err != nil {
		return res, err
	}
	defer cancel()
	defer db.Close()

	query := `SELECT * FROM users`
	if sortBy != pb.By_NONE.Enum() {
		query += " WHERE " + strings.ToLower(sortBy.String()) + " LIKE $1"
	}

	if err := db.SelectContext(ctx, &res, query, "%"+sortValue+"%"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, fmt.Errorf("%w: пользователь с полем %s = %s не существует", ErrNotFound, sortBy.String(), sortValue)
		}
		return res, fmt.Errorf("getting user by id: select error: %w", err)
	}
	return res, nil
}

func UpdateUser(user User) error {

	test, err := GetUserById(user.ID)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}
	if _, err := GetUserByEmail(user.Email); !errors.Is(err, sql.ErrNoRows) && test.Email != user.Email {
		if err == nil {
			return fmt.Errorf("%w: пользователь с почтой %s уже существует", ErrBadRequest, user.Email)
		}
		return fmt.Errorf("updating user: %w", err)
	}

	ctx, cancel, db, err := connectToDb()
	if err != nil {
		return err
	}
	defer cancel()
	defer db.Close()

	query := `UPDATE users 
		SET email = :email, first_name = :first_name, second_name = :second_name, 
		patronymic = :patronymic WHERE id = :id`

	if _, err := db.NamedExecContext(ctx, query, &user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, user.ID)
		}
		return fmt.Errorf("updating user: update error: %w", err)
	}

	return nil
}

func UpdatePassword(id uint32, password string) error {

	ctx, cancel, db, err := connectToDb()
	if err != nil {
		return fmt.Errorf("updating password: %w", err)
	}
	defer cancel()
	defer db.Close()

	query := `UPDATE users SET password = $1 WHERE id = $2`

	if _, err := db.ExecContext(ctx, query, password, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, id)
		}
		return fmt.Errorf("updating password: update error: %w", err)
	}

	return nil
}

func DeleteUser(id uint32) error {

	ctx, cancel, db, err := connectToDb()
	if err != nil {
		return err
	}
	defer cancel()
	defer db.Close()

	if _, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, id)
		}
		return fmt.Errorf("deleting user: db error: %w", err)
	}
	return nil
}
