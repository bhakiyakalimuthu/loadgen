package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strings"
)

func main() {
	db := sqlx.MustConnect("postgres", "postgres://postgres:postgres@localhost:5432/test?sslmode=disable")
	c := []Contact{
		{FirstName: "first1", LastName: "last1", Email: "email1"},
		{FirstName: "first2", LastName: "last2", Email: "email2"},
	}
	BulkInsert(c, db)
}

type Contact struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

func BulkInsert(cs []Contact, db *sqlx.DB) error {
	var (
		placeholders []string
		vals         []interface{}
	)

	for index, contact := range cs {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d)",
			index*3+1,
			index*3+2,
			index*3+3,
		))

		vals = append(vals, contact.FirstName, contact.LastName, contact.Email)
	}

	txn, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not start a new transaction %v", err)
	}

	insertStatement := fmt.Sprintf("INSERT INTO contacts(first_name,last_name,email) VALUES %s", strings.Join(placeholders, ","))
	_, err = txn.Exec(insertStatement, vals...)
	if err != nil {
		txn.Rollback()
		return fmt.Errorf("failed to insert multiple records at once %v", err)
	}

	if err := txn.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction %v", err)
	}

	return nil
}
