package main

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const SQL_STMT = `
create table jots (
	id integer not null primary key,
	path text unique
);

create table entries (
	id integer not null primary key,
	jot_id integer,
	title text,
	content text,
	last_update timestamp default (datetime('now','localtime')),
	foreign key (jot_id) references jots (id),
	unique (jot_id, title)
);
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

	info, err := os.Stat(dbPath)
	if !os.IsNotExist(err) && !info.IsDir() {
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

func (d *DB) initialize(curPath string) error {
	stmt, err := d.Prepare("insert into jots (path) values (?)")
	if err != nil {
		return fmt.Errorf("init: couldn't setup prepared statement: %s", err)
	}
	if _, err := stmt.Exec(curPath); err != nil {
		return fmt.Errorf("init: couldn't insert: %s", err)
	}
	return nil
}

func (d *DB) createJot(jotId int64, title, content string) (int64, error) {
	ins := "insert into entries (jot_id, title, content) values (?, ?, ?)"
	stmt, err := d.Prepare(ins)
	if err != nil {
		return 0, fmt.Errorf("couldn't setup prepared statement: %s", err)
	}
	res, err := stmt.Exec(jotId, title, content)
	if err != nil {
		return 0, fmt.Errorf("couldn't execute prep statment:", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("couldn't get last id:", err)
	}
	return id, err
}

func (d *DB) getJotDir(jotPath string) (int64, error) {
	q := "select id, max(length(path)) from jots where path in (%s)"
	pathParts := strings.Split(jotPath, "/")
	partLen := len(pathParts)
	paths := make([]interface{}, partLen, partLen)
	qstns := make([]string, partLen, partLen)
	for i := 0; i < partLen; i++ {
		x := partLen - i
		paths[i] = "/" + path.Join(pathParts[:x]...)
		qstns[i] = "?"
	}
	var id, l sql.NullInt64
	var query string = fmt.Sprintf(q, strings.Join(qstns, ","))
	err := d.QueryRow(query, paths...).Scan(&id, &l)
	if err != nil {
		return 0, err
	}
	if id.Valid {
		return id.Int64, nil
	}
	return 0, fmt.Errorf("couldn't find any jots for dir: %s", jotPath)
}

func (d *DB) listByPath(jotPath string) (*sql.Rows, error) {
	id, err := d.getJotDir(jotPath)
	if err != nil {
		return nil, err
	}
	q := `select title, last_update
		from  entries where jot_id = ?`
	return d.Query(q, id)
}

func (d *DB) get(path, title string) *sql.Row {
	q := `select b.content from jots a
		inner join entries b on a.id = b.jot_id
		where a.path like ?`
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

func bain() {
	//db, err := setupDB()
	//if err != nil {
	//fmt.Println("error:", err)
	//return
	//}
	//defer func() { fmt.Println("closed the db"); db.Close() }()

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