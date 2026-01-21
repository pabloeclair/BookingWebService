package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	ErrConDB    error = errors.New("ошибка подключения к базе данных")
	ErrConflict error = errors.New("произошел конфликт данных")
	ErrNotFound error = errors.New("по запросу ничего не было найдено")
	ErrBadKey   error = errors.New("передан некорректный ключ")
	DSN         string
)

// Структура пользователя для сохранения в базе данных.
type User struct {
	Id         uint32 `db:"id"`
	Email      string `db:"email"`
	FirstName  string `db:"first_name"`
	SecondName string `db:"second_name"`
	Patronymic string `db:"patronymic"`
	Password   string `db:"password"`
	Role       string `db:"role"`
}

// Преобразование ФИО в Title Case, т.к. изначально эти данные передаются в нижнем регистре
func (u *User) finallyFieldsProcessing() {
	u.FirstName = strings.ToTitle(string(u.FirstName[0])) + u.FirstName[1:]
	u.SecondName = strings.ToTitle(string(u.SecondName[0])) + u.SecondName[1:]
	if u.Patronymic != "" {
		u.Patronymic = strings.ToTitle(string(u.Patronymic[0])) + u.Patronymic[1:]
	}
}

// Подключение к базе данных.
func connectToDb() (*sqlx.DB, error) {
	if DSN == "" {
		return nil, fmt.Errorf("connecting to db: %w: %s", ErrConDB, "отсутствует DSN")
	}
	db, err := sqlx.Connect("pgx", DSN)
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w: %s", ErrConDB, err.Error())
	}

	return db, nil
}

// Создание таблицы пользователей, если ее не существует.
func CreateTable() error {
	db, err := connectToDb()
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
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

	return nil
}

// Удаление таблицы. Применяется для интеграционных тестов.
func DeleteTable() error {

	db, err := connectToDb()
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
	defer cancel()

	_, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS users;`)
	if err != nil {
		return fmt.Errorf("deleting table: drop error: %w", err)
	}
	return nil
}

// Сохранение нового пользователя в базу данных. Может вернуть ErrConflict и ErrNotFound.
//
// Принимает контекст запроса и указатель на объект user с описанием всех
// полей, кроме id. Также необязательным полем является Patronymic. При успехе
// возвращает id пользователя.
func CreateUser(ctx context.Context, user *User) (uint32, error) {

	// в случае, если пользователь с указанной почтой уже существует, возвращает ErrConflict.
	if _, err := GetUserByEmail(ctx, user.Email); !errors.Is(err, ErrNotFound) {
		if err != nil {
			return 0, fmt.Errorf("creating user: ошибка проверки почты: %w", err)
		}
		return 0, fmt.Errorf("%w: пользователь с почтой %s уже существует", ErrConflict, user.Email)
	}

	db, err := connectToDb()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var id uint32
	// если пользователь первый, он обретает права MAIN_ADMIN
	if err = db.GetContext(ctx, &id, `SELECT id FROM users LIMIT 1`); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user.Role = "MAIN_ADMIN"
		} else {
			return 0, fmt.Errorf("creating user: select error: %w", err)
		}
	}

	// сохранение пользователя в бд
	queryInsert := `INSERT INTO users (email, first_name, second_name, patronymic, password, role) 
		VALUES (:email, :first_name, :second_name, :patronymic, :password, :role);`
	if _, err = db.NamedExecContext(ctx, queryInsert, &user); err != nil {
		return 0, fmt.Errorf("creating user: insert error: %w", err)
	}

	var idRes uint32
	// получение id пользователя и проверка, что все успешно сохранилось
	querySelect := `SELECT id FROM users WHERE email = $1;`
	if err := db.GetContext(ctx, &idRes, querySelect, user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, idRes)
		}
		return 0, fmt.Errorf("getting user by id: select error: %w", err)
	}

	return idRes, nil
}

// Получение информации о пользователе по email. Может вернуть ErrNotFound.
//
// Принимает контекст запроса и строку почты. При успехе возвращает указатель
// на объект пользователя.
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	db, err := connectToDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var res User
	query := `SELECT * FROM users WHERE LOWER(email) = $1;`
	if err := db.GetContext(ctx, &res, query, strings.ToLower(email)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: пользователь с почтой %s не существует", ErrNotFound, email)
		}
		return nil, fmt.Errorf("getting user by email: select error: %w", err)
	}

	res.finallyFieldsProcessing()
	return &res, nil
}

// Получение пользователя по id. Может вернуть ErrNotFound.
//
// Принимает контекст запроса и id пользователя. При успехе возвращает указатель
// на объект пользователя.
func GetUserById(ctx context.Context, id uint32) (*User, error) {
	db, err := connectToDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var res User
	query := `SELECT * FROM users WHERE id = $1;`
	if err := db.GetContext(ctx, &res, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, id)
		}
		return nil, fmt.Errorf("getting user by id: select error: %w", err)
	}

	res.finallyFieldsProcessing()
	return &res, nil
}

// // Получение списка пользователей по ключу. Может вернуть ErrBadKey.
// //
// // Принимает контекст запроса, категорию сортировки и искомое значение.
// // При успехе возвращает указатель на список объектов пользователей или пустой список,
// // если ничего не найдено.
// func GetUserByKey(ctx context.Context, sortBy *pb.By, sortValue string) ([]*User, error) {
// 	db, err := connectToDb()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	// todo: избавиться от sql инъекции
// 	var res []*User
// 	query := `SELECT * FROM users`
// 	if sortBy != pb.By_NONE.Enum() {
// 		query += " WHERE " + strings.ToLower(sortBy.String()) + " LIKE $1"
// 	}
// 	if err := db.SelectContext(ctx, &res, query, "%"+sortValue+"%"); err != nil {
// 		return nil, fmt.Errorf("getting user by id: select error: %w", err)
// 	}

// 	for _, u := range res {
// 		u.finallyFieldsProcessing()
// 	}
// 	return res, nil
// }

// Обновление информации о пользователе в базе данных. Может вернуть ErrConflict и ErrNotFound.
//
// Принимает контекст запроса и указатель на объект пользователя с полной информацией, кроме роли.
// При успехе не вернет ошибку.
func UpdateUser(ctx context.Context, user *User) error {

	// получение старой почты
	old, err := GetUserById(ctx, user.Id)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}
	if _, err := GetUserByEmail(ctx, user.Email); !errors.Is(err, sql.ErrNoRows) && old.Email != user.Email {
		if err == nil {
			return fmt.Errorf("%w: пользователь с почтой %s уже существует", ErrConflict, user.Email)
		}
		return fmt.Errorf("updating user: %w", err)
	}

	db, err := connectToDb()
	if err != nil {
		return err
	}
	defer db.Close()

	// сохранение данных
	query := `UPDATE users 
		SET email = :email, first_name = :first_name, second_name = :second_name, 
		patronymic = :patronymic WHERE id = :id`

	if _, err := db.NamedExecContext(ctx, query, &user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, user.Id)
		}
		return fmt.Errorf("updating user: update error: %w", err)
	}

	return nil
}

// Обновление пароля пользователя. Может вернуть ErrNotFound.
//
// Принимает контекст запроса, id и зашифрованный пароль пользователя.
// При успехе не вернет ошибку.
func UpdatePassword(ctx context.Context, id uint32, password string) error {
	db, err := connectToDb()
	if err != nil {
		return fmt.Errorf("updating password: %w", err)
	}
	defer db.Close()

	// сохранение пароля
	query := `UPDATE users SET password = $1 WHERE id = $2`
	if _, err := db.ExecContext(ctx, query, password, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, id)
		}
		return fmt.Errorf("updating password: update error: %w", err)
	}

	return nil
}

// Обновление роли пользователя. Может вернуть NotFound.
//
// Принимает контекст запроса, id и роль обновляемого пользователя.
// При успехе вернет обновленную информацию о пользователе.
func UpdateRole(ctx context.Context, id uint32, role string) (*User, error) {
	db, err := connectToDb()
	if err != nil {
		return nil, fmt.Errorf("updating password: %w", err)
	}
	defer db.Close()

	// сохранение роли
	query := `UPDATE users SET role = $1 WHERE id = $2`
	if _, err := db.ExecContext(ctx, query, role, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, id)
		}
		return nil, fmt.Errorf("updating password: update error: %w", err)
	}

	// получение итогового пользователя
	res, err := GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Удаляет пользователя по id. Может вернуть ErrNotFound.
//
// Принимает контекст запроса и id пользователя. При успехе не вернет ошибку.
func DeleteUser(ctx context.Context, id uint32) error {
	db, err := connectToDb()
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: пользователь с id = %d не существует", ErrNotFound, id)
		}
		return fmt.Errorf("deleting user: db error: %w", err)
	}
	return nil
}
