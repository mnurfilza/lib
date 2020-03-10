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
	fields, dst := tb.Fields()
	query := fmt.Sprintf("INSERT INTO %s VALUES(%s)", tb.Name(), PlaceHolder(len(fields)))

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(dst...)
	if err != nil {
		return err
	}
	return err
}

func Update(db *sql.DB, tb Table, change map[string]interface{}) error {
	var w, set []string
	var setQuery, setWhere string
	var args []interface{}
	pk, dst := tb.PrimaryKey()

	// ini untuk set
	for val, v := range change {
		args = append(args, v)
		temp := fmt.Sprintf("%s = ?", val)
		set = append(set, temp)
		setQuery = strings.Join(set, ",")

	}

	for _, prim := range pk {
		setWheres := fmt.Sprintf("%s = ?", prim)
		w = append(w, setWheres)
		setWhere = strings.Join(w, ",")
	}

	args = append(args, dst...)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tb.Name(), setQuery, setWhere)
	fmt.Println(query)
	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
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
