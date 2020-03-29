package lib

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Table interface {
	Name() string
	Fields() (fields []string, dst []interface{})
	PrimaryKey() (fiedls []string, dst []interface{})
	Structur() Table
}

type Params struct {
	Field string
	Op    string
	Value interface{}
}

type RequestParams struct {
	Limit   int
	Offset  int
	Param   []Params
	OrderBy string
	Sort    string
}

func CreateDatabase(db *sql.DB, nama string) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v", nama)
	_, err := db.Exec(query)
	return err
}

func Connect(user, password, host, port, dbname string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True", user, password, host, port, dbname))
	fmt.Println(db)
	return db, err

}

func DropDB(db *sql.DB, name string) error {
	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", name)
	_, err := db.Exec(query)
	return err
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

func Delete(db *sql.DB, tb Table) error {
	pk, dst := tb.PrimaryKey()
	setWheres := fmt.Sprintf("%s = ?", pk[0])

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tb.Name(), setWheres)
	_, err := db.Exec(query, dst...)
	if err != nil {
		return err
	}
	return err
}

func Get(db *sql.DB, tb Table) error {
	pk, dstPk := tb.PrimaryKey()
	_, dst := tb.Fields()
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tb.Name(), pk[0])
	if err := db.QueryRow(query, dstPk...).Scan(dst...); err != nil {
		return err
	}

	return nil
}

func Fetch(db *sql.DB, tb Table, p RequestParams) ([]Table, error) {
	var param []interface{}
	var where []string
	query := fmt.Sprintf("SELECT * FROM %s", tb.Name())

	if len(p.Param) != 0 {
		for _, item := range p.Param {
			where = append(where, fmt.Sprintf("%s = ?", item.Field))

			param = append(param, item.Value)
		}

		whereKondisi := strings.Join(where, " AND ")
		query = query + " WHERE " + whereKondisi
	}

	if p.OrderBy != "" && p.Sort != "" {
		query = query + " ORDER BY " + p.OrderBy + p.Sort
	}

	rows, err := db.Query(query, param...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Table
	for rows.Next() {
		each := tb.Structur()
		_, dst := each.Fields()
		var err = rows.Scan(dst...)
		if err != nil {
			return nil, err
		}
		res = append(res, each)
	}
	return res, nil
}

func PlaceHolder(jml int) string {
	jumlah := make([]string, jml)
	for i, _ := range jumlah {
		jumlah[i] = "?"
	}
	placeholder := strings.Join(jumlah, ",")
	return placeholder

}

func QueryLimitOffset(values url.Values) (int, int, error) {
	lim := values.Get("limit")
	off := values.Get("offset")
	if lim != "" && off != "" {
		limit, err := strconv.Atoi(lim)
		if err != nil {
			return 0, 0, err
		}
		offset, err := strconv.Atoi(off)
		if err != nil {
			return 0, 0, err
		}
		return limit, offset, err
	}
	return 0, 0, nil
}

func QueryParams(values url.Values) ([]Params, error) {
	val := values.Get("params")
	var params []Params
	if val != "" {
		param := strings.Split(val, ":")
		for _, item := range param {
			temp := Params{}
			obj := strings.Split(item, ",")
			if len(obj) != 3 {
				return nil, fmt.Errorf("Parameter Tidak Sesuai")
			}
			temp.Field = obj[0]
			temp.Op = obj[1]
			temp.Value = obj[2]
			params = append(params, temp)
		}
	}
	return params, nil
}

func QueryOrder(value url.Values) (string, error) {
	val := value.Get("orderby")
	return val, nil
}

func QuerySort(values url.Values) (string, error) {
	return values.Get("sort"), nil
}

func countColumn(db *sql.DB, tb Table) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tb.Name())
	fmt.Println(query)
	return 0, nil
}
