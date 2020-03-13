package sqlmocktest_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	smt "github.com/makiuchi-d/sqlmocktest"
)

func TestInit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()
	repo := smt.NewRepo(sqlx.NewDb(db, "sqlite3"))

	mock.ExpectExec("CREATE TABLE user ").WillReturnResult(sqlmock.NewResult(0, 0))
	if err := repo.Init(); err != nil {
		t.Fatalf("Init error: %+v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInitFaill(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()
	repo := smt.NewRepo(sqlx.NewDb(db, "sqlite3"))

	expErr := fmt.Errorf("mock error")
	mock.ExpectExec("CREATE TABLE user ").WillReturnError(expErr)
	if err := repo.Init(); !errors.Is(err, expErr) {
		t.Fatalf("Init must expected error: %+v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()
	repo := smt.NewRepo(sqlx.NewDb(db, "sqlite3"))

	user := smt.User{
		Id:   10,
		Name: "Ten",
	}

	mock.ExpectExec("insert into user ").WithArgs(10, "Ten").WillReturnResult(sqlmock.NewResult(10, 1))
	res, err := repo.Insert(&user)
	if err != nil {
		t.Fatalf("Insert err: %v", err)
	}
	i, _ := res.LastInsertId()
	a, _ := res.RowsAffected()
	if i != 10 || a != 1 {
		t.Fatalf("result: {%v, %v} wants {10, 1}", i, a)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}
	defer db.Close()
	repo := smt.NewRepo(sqlx.NewDb(db, "sqlite3"))

	mock.ExpectQuery("SELECT \\* FROM user WHERE id = \\?").WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(5, "Five"))

	u, err := repo.GetUser(5)
	if err != nil {
		t.Fatalf("GetUser(5) error: %v", err)
	}
	if u.Id != 5 || u.Name != "Five" {
		t.Fatalf("GetUser(5) = %v, expect &{5 \"Five\"}", u)
	}

	mock.ExpectQuery("SELECT \\* FROM user WHERE id = \\?").WithArgs(10).WillReturnError(sql.ErrNoRows)

	u, err = repo.GetUser(10)
	if err != nil {
		t.Fatalf("GetUser(10) error: %v", err)
	}
	if u != nil {
		t.Fatalf("GetUser(10) = %v, expect nil", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
