package main

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SQL_STMT = `
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
	foreign key (jot_id) references jots (id) on delete cascade,
	unique (jot_id, title)
);
`
	ERR_TRACKED = "UNIQUE constraint failed: jots.path"
)

func getDBPath() string {
	jotHome := os.Getenv("JOTHOME")
	if jotHome == "" {
		jotHome = path.Join(os.Getenv("HOME"), ".jot")
	}
	return jotHome
}

func setupDB(dbDir string) (*DB, error) {
	// create the directory if it doesn't exist
	var isNew bool = false
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if mkerr := os.MkdirAll(dbDir, 0775); mkerr != nil {
			return nil, fmt.Errorf("couldn't set up db:", mkerr)
		}
		isNew = true
	}

	// check if db exists, return connection if it does
	dbPath := filepath.Join(dbDir, "jot.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		isNew = true
	}
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_foreign_keys=on", dbPath))
	if err != nil {
		return nil, fmt.Errorf("couldn't open database: %s", err)
	}
	if !isNew {
		return &DB{db}, nil
	}

	// create tables if it's a new database
	_, err = db.Exec(SQL_STMT)
	if err != nil {
		return nil, fmt.Errorf("Error executing statement:", err)
	}
	return &DB{db}, nil
}

type DB struct {
	*sql.DB
}

func (d *DB) uninitialize(curPath string) error {
	stmt, err := d.Prepare("delete from jots where path = ?")
	if err != nil {
		return fmt.Errorf("init: couldn't setup prepared statement: %s", err)
	}
	res, err := stmt.Exec(curPath)
	if err != nil {
		return fmt.Errorf("init: couldn't insert: %s", err)
	}
	no, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if no == 0 {
		return fmt.Errorf("didn't untrack jot dir '%s', be sure to CWD to the tracked dir before untracking", curPath)
	}
	return nil
}

func (d *DB) initialize(curPath string) error {
	stmt, err := d.Prepare("insert into jots (path) values (?)")
	if err != nil {
		return fmt.Errorf("init: couldn't setup prepared statement: %s", err)
	}
	if _, err := stmt.Exec(curPath); err != nil {
		if err.Error() == ERR_TRACKED {
			return fmt.Errorf("directory already tracked")
		}
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

	//@TODO: add a limit here, say of only 10 upper dirs
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
	if !id.Valid {
		return 0, fmt.Errorf("jot dir not initialized: %s", jotPath)
	}
	return id.Int64, nil
}

func (d *DB) listAllDirs() (*sql.Rows, error) {
	return d.Query("select path from jots")
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

type Jot struct {
	jotId       int64
	title       string
	contents    *string
	lastUpdated time.Time
}

func (d *DB) get(jotId int64, title string) (*Jot, error) {
	var j Jot
	var contents sql.NullString
	q := `select content, last_update from entries
		where jot_id = ? and title = ?`
	err := d.QueryRow(q, jotId, title).Scan(&contents, &(j.lastUpdated))
	if err != nil {
		return nil, err
	}
	if !contents.Valid {
		return nil, fmt.Errorf("no jot named:", title)
	}
	j.contents = &contents.String
	j.title = title
	j.jotId = jotId
	return &j, nil
}

func (d *DB) delete(jotId int64, title string) error {
	q := "delete from entries where jot_id = ? and title = ?"
	stmt, err := d.Prepare(q)
	if err != nil {
		return fmt.Errorf("couldn't delete:", err)
	}
	_, err = stmt.Exec(jotId, title)
	return err
}
