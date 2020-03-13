package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	smt "github.com/makiuchi-d/sqlmocktest"
)

func main() {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := smt.NewRepo(db)
	if err := repo.Init(); err != nil {
		panic(err)
	}

	r, err := repo.Insert(&smt.User{1, "First"})
	if err != nil {
		panic(err)
	}
	lastId, err := r.LastInsertId()
	if err != nil {
		panic(err)
	}
	fmt.Printf("last insert: %v\n", lastId)

	u, err := repo.GetUser(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("user[1] = %v\n", u)

	u, err = repo.GetUser(2)
	if err != nil {
		fmt.Printf("%T\n", err)
		panic(err)
	}
	fmt.Printf("user[2] = %v\n", u)
}
