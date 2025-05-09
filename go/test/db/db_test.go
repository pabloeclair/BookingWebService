package main

import (
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"database/sql"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {

	t.Run("CreateFirstUser", func(t *testing.T) {
		user := db.User{
			Email:      "test1@mail.ru",
			FirstName:  "Гленфорд",
			SecondName: "Дж.",
			Patronymic: "Майерс",
			Password:   "ArtSoftwareTesting",
		}

		res, err := db.CreateUser(user)
		if err != nil {
			t.Fatalf("Expected error: err = nil; actual error: err = %v", err)
		}
		if res.ID != 1 || res.Role != pb.Role_ADMIN {
			t.Fatalf("Expected result: id = 1, role = ADMIN; actual result: id = %d, role = %s", res.ID, res.Role.String())
		}
	})

	t.Run("CreateWithoutPatronymicUser", func(t *testing.T) {
		user := db.User{
			Email:      "test2@mail.ru",
			FirstName:  "Майкл",
			SecondName: "Болтон",
			Password:   "RapidSoftwareTesting",
		}

		res, err := db.CreateUser(user)
		if err != nil {
			t.Fatalf("Expected error: err = nil; actual error: err = %v", err)
		}
		if res.ID != 2 || res.Role != pb.Role_USER {
			t.Fatalf("Expected result: id = 2, role = USER; actual result: id = %d, role = %s", res.ID, res.Role.String())
		}
	})

	t.Run("CreateSameUser", func(t *testing.T) {
		user := db.User{
			Email:      "test2@mail.ru",
			FirstName:  "SameМайкл",
			SecondName: "Same Болтон",
			Password:   "RapidSoftwareTesting",
		}

		_, err := db.CreateUser(user)
		if err == nil {
			t.Fatal(`Expected error: err = creating user: insert error: duplicate key value violates unique constraint "users_email_key"; actual error: err = nil"`)
		}
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code != "23505" {
				t.Fatalf(`Expected error: code = 23505, msg = duplicate key value violates unique constraint "users_email_key"; actual error: code = %s, msg = %s`, pgErr.Code, pgErr.Message)
			}
		}
		t.Fatalf(`Expected error: err = creating user: insert error: duplicate key value violates unique constraint "users_email_key"; actual error: err = %v`, err)
	})

}

func TestGetUserByEmail(t *testing.T) {

	t.Run("GetFirstUser", func(t *testing.T) {
		email := "test1@mail.ru"

		user, err := db.GetUserByEmail(email)
		if err != nil {
			t.Fatalf("Expected error: err = nil; actual error: err = %v", err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("ArtSoftwareTesting"))
		if user.Email != "test1@mail.ru" || err != nil {
			t.Fatalf("Expected result: email = %s, errComparePassword = nil; actual result: email = %s, errComparePassword = %v", email, user.Email, err)
		}
	})

	t.Run("GetSecondUser", func(t *testing.T) {
		email := "test2@mail.ru"

		user, err := db.GetUserByEmail(email)
		if err != nil {
			t.Fatalf("Expected error: err = nil; actual error: err = %v", err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("RapidSoftwareTesting"))
		if user.Email != "test2@mail.ru" || err != nil {
			t.Fatalf("Expected result: email = %s, errComparePassword = nil; actual result: email = %s, errComparePassword = %v", email, user.Email, err)
		}
	})

	t.Run("GetUserByInvalidEmail", func(t *testing.T) {
		email := "notexistsemail@gmail.com"

		_, err := db.GetUserByEmail(email)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("Expected error: err = ErrNoRows; actual error: err = %v", err)
		}
	})

}
