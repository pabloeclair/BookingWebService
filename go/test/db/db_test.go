package main

import (
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"database/sql"
	"errors"
	"testing"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {

	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal(err)
	}
	t.Run("CreateFirstUser", func(t *testing.T) {
		user := db.User{
			Email:      "test1@mail.ru",
			FirstName:  "Гленфорд",
			SecondName: "Дж.",
			Patronymic: "Майерс",
			Password:   "ArtSoftwareTesting",
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			t.Fatalf("Internal error: bcrypt: %v", err)
		}
		user.Password = string(hashedPassword)

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

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			t.Fatalf("Internal error: bcrypt: %v", err)
		}
		user.Password = string(hashedPassword)

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

		expectedError := `creating user: insert error: ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`

		_, err := db.CreateUser(user)
		if err.Error() != expectedError {
			t.Fatalf(`Expected error: err = %s; actual error: err = %v`, expectedError, err)
		}
	})

	t.Run("CreateInvalidUser", func(t *testing.T) {
		user := db.User{
			Email:    "test3@mail.ru",
			Password: "invalidpassword",
		}

		expectedError := `creating user: insert error: ERROR: null value in column "first_name" violates not-null constraint (SQLSTATE 23502)`

		_, err := db.CreateUser(user)
		if err.Error() != expectedError {
			t.Fatalf("Expected error: err = %s; actual error: err = %v", expectedError, err)
		}
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
		if user.ID != 1 || user.Email != "test1@mail.ru" || err != nil {
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
		if user.ID != 2 || user.Email != "test2@mail.ru" || err != nil {
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

func TestUpdateUser(t *testing.T) {

	t.Run("UpdateFirstUser", func(t *testing.T) {
		user := db.User{
			Email:      "test1@mail.ru",
			FirstName:  "Гленфорд",
			SecondName: "Дж.",
			Password:   "NewArtSoftwareTesting",
		}

		err := db.UpdateUser(user.Email, user)
		if err != nil {
			t.Fatalf("Expected error: err = nil; actual error: err = %v", err)
		}

		actualUser, err := db.GetUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Expected error when getting actual user: err = nil; actual error: err = %v", err)
		}

		if actualUser.Password != user.Password {
			t.Fatalf("Expected result: password = NewArtSoftwareTesting; actual result: password = %s", actualUser.Password)
		}

	})

	t.Run("FullUpdateSecondUser", func(t *testing.T) {
		user := db.User{
			Email:      "updatetest2@gmail.com",
			FirstName:  "Святослав",
			SecondName: "Куликов",
			Patronymic: "Святославович",
			Password:   "RelationalDatabasesInTheExamples",
		}

		err := db.UpdateUser("test2@mail.ru", user)
		if err != nil {
			t.Fatalf("Expected error: err = nil; actual error: err = %v", err)
		}

		actualUser, err := db.GetUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Expected error when getting actual user: err = nil; actual error: err = %v", err)
		}

		if actualUser.Password != user.Password {
			t.Fatalf("Expected result: password = RelationalDatabasesInTheExamples; actual result: password = %s", actualUser.Password)
		}
	})

	t.Run("InvalidUpdateSecondUser", func(t *testing.T) {
		user := db.User{
			Email:    "updatetest2@gmail.com",
			Password: "RelationalDatabasesInTheExamples",
		}

		expectedError := `creating user: insert error: ERROR: null value in column "first_name" violates not-null constraint (SQLSTATE 23502)`
		err := db.UpdateUser("updatetest2@gmail.com", user)
		if err.Error() != expectedError {
			t.Fatalf("Expected error: err = %s; actual error: err = %v", expectedError, err)
		}
	})
}
