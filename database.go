package lib

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Table interface {
	Name() string
	Fields() (fields []string, dst []interface{})
	PrimaryKey() (fiedls []string, dst []interface{})
}

func CreateDatabase(db *sql.DB, nama string) error {
	query := fmt.Sprintf("CREATE DATABASE %v", nama)
	_, err := db.Exec(query)
	return err
}

func Connect(user, password, host, port, dbname string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", user, password, host, port, dbname))
	fmt.Println(db)
	return db, err

}

func Use(db *sql.DB, name string) error {
	query := fmt.Sprintf("Use %v", name)
	_, err := db.Exec(query)
	return err
}

func CreateTable(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	return err
}

func Insert(db *sql.DB, tb Table) error {
	return err
}

func PlaceHolder(jml int) string {
	jumlah := make([]string, jml)
	for i, _ := range jumlah {
		jumlah[i] = "?"
	}
	placeholder := strings.Join(jumlah, ",")
	return placeholder
}
