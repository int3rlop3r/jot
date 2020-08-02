package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const SQL_STMT = `
create table jots (
	id integer not null primary key,
	title text,
	path text,
	content text,
	last_update timestamp default current_timestamp
)
`

func getDBPath() string {
	jothome := os.Getenv("JOTHOME")
	if jothome == "" {
		jothome = filepath.Join(os.Getenv("HOME"), ".jot")
	}
	return filepath.Join(jothome, "jot.db")
}

func setupDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open database: %s", err)
	}
	fmt.Println("opened the db")

	info, err := os.Stat(dbPath)
	if !os.IsNotExist(err) && !info.IsDir() {
		fmt.Println("DB exists already")
		return &DB{db}, nil
	}

	_, err = db.Exec(SQL_STMT)
	if err != nil {
		return nil, fmt.Errorf("Error executing statement:", err)
	}
	return &DB{db}, nil
}

type DB struct {
	*sql.DB
}

func (d *DB) insert(path, title, content string) (int64, error) {
	ins := "insert into jots (title, path, content) values (?, ?, ?)"
	stmt, err := d.Prepare(ins)
	if err != nil {
		return 0, fmt.Errorf("couldn't setup prepared statement: %s", err)
	}
	res, err := stmt.Exec(title, path, content)
	if err != nil {
		return 0, fmt.Errorf("couldn't execute prep statment:", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("couldn't get last id:", err)
	}
	return id, err
}

func (d *DB) listByPath(path string) (*sql.Rows, error) {
	q := "select * from jots where path = ?"
	return d.Query(q, path)
}

func (d *DB) get(path, title string) *sql.Row {
	q := "select * from jots where path = ? and title = ?"
	return d.QueryRow(q, path, title)
}

func (d *DB) delete(path, title string) error {
	q := "delete from jots where path = ? and title = ?"
	stmt, err := d.Prepare(q)
	if err != nil {
		return fmt.Errorf("couldn't delete:", err)
	}
	_, err = stmt.Exec(path, title)
	return err
}

func main() {
	db, err := setupDB()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer func() { fmt.Println("closed the db"); db.Close() }()

	//// test insert
	//id, err := db.insert("some path3", "some title", "another thing I had to say")
	//if err != nil {
	//fmt.Println("error while inserting:", err)
	//return
	//}
	//fmt.Println("Last insert id:", id)

	//// test select
	//rows, err := db.listByPath("some path3")
	//if err != nil {
	//fmt.Println("error while getting rows:", err)
	//return
	//}
	//var id int
	//var title, path, content string
	//var last_updated time.Time
	//for rows.Next() {
	//e := rows.Scan(&id, &title, &path, &content, &last_updated)
	//if e != nil {
	//fmt.Println("scan error:", e)
	//}
	//fmt.Println(id, title, path, content, last_updated)
	//}

	//// test delete
	//err = db.delete("some path3", "some title")
	//if err != nil {
	//fmt.Println(err)
	//}
}
