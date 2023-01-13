package sqlite

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/PaulioRandall/go-trackerr"

	"github.com/PaulioRandall/platos-pizza-go/pkg/database"
)

// TODO: Create SQLite helper library?
// TODO: Convert to using transactions for bulk inserts

const (
	insertBatchSize = 256
	queryHeadMax    = database.QueryHeadMax
)

var (
	ErrSQLite = trackerr.Checkpoint("SQLite database error")

	// See (SQLite) https://www.sqlite.org/pragma.html
	// See (go-sqlite3) https://github.com/mattn/go-sqlite3#connection-string
	sqlitePragma = []string{
		// 5x the default page & cache sizes to speed things up (I'm on desktop)
		"_page_size=5120",
		"_cache_size=10000",

		// Turn off stuff that slows inserts down, just while developing.
		"_journal=MEMORY",
		//"_foreign_keys=OFF",
		//"_ignore_check_constraints=ON",
		"_sync=OFF",
	}
)

type sqliteDB struct {
	conn *sql.DB
}

func OpenDatabase(file string) (*sqliteDB, error) {
	fileURL := file + "?" + strings.Join(sqlitePragma, "&")

	conn, e := sql.Open("sqlite3", fileURL)
	if e != nil {
		return nil, ErrSQLite.Wrap(e)
	}

	db := &sqliteDB{
		conn: conn,
	}

	if e = db.createTables(); e != nil {
		return nil, ErrSQLite.Wrap(e)
	}

	return db, nil
}

func (db *sqliteDB) createTables() error {
	sql := joinLines(
		`CREATE TABLE metadata (`,
		`	id          INTEGER NOT NULL PRIMARY KEY,`, // Alias for SQLite 'rowid'
		`	table_name  TEXT    NOT NULL,`,
		`	field_name  TEXT    NOT NULL,`,
		`	description TEXT    NOT NULL`,
		`);`,
		``,
		`CREATE TABLE pizza_types (`,
		`	id          TEXT NOT NULL PRIMARY KEY,`,
		`	name        TEXT NOT NULL,`,
		`	category    TEXT NOT NULL,`,
		`	ingredients TEXT NOT NULL`,
		`);`,
		``,
		`CREATE TABLE pizzas (`,
		`	id      TEXT NOT NULL PRIMARY KEY,`,
		`	type_id TEXT NOT NULL,`,
		`	size    TEXT NOT NULL,`,
		`	price   REAL NOT NULL,`,
		`	FOREIGN KEY(type_id) REFERENCES pizza_types(id)`,
		`);`,
		``,
		`CREATE TABLE orders (`,
		`	id       INTEGER NOT NULL PRIMARY KEY,`,
		`	datetime TEXT    NOT NULL`,
		`);`,
		``,
		`CREATE TABLE order_details (`,
		`	id       INTEGER NOT NULL PRIMARY KEY,`,
		`	order_id INTEGER NOT NULL,`,
		`	pizza_id TEXT    NOT NULL,`,
		`	quantity INTEGER NOT NULL,`,
		`	FOREIGN KEY(order_id) REFERENCES orders(id),`,
		`	FOREIGN KEY(pizza_id) REFERENCES pizzas(id)`,
		`);`,
	)

	if _, e := db.conn.Exec(sql); e != nil {
		return database.ErrCreating.Wrap(e)
	}

	return nil
}

func (db *sqliteDB) insert(sql string, params []any) error {
	stmt, e := db.conn.Prepare(sql)
	if e != nil {
		e = database.ErrPreparing.Wrap(e)
		return database.ErrInserting.Wrap(e)
	}
	defer stmt.Close()

	if _, e := stmt.Exec(params...); e != nil {
		return database.ErrInserting.Wrap(e)
	}

	return nil
}

func (db *sqliteDB) Close() {
	db.conn.Close()
}
