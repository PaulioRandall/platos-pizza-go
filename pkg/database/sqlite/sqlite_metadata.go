package sqlite

import (
	"database/sql"

	"github.com/PaulioRandall/platos-pizza-go/pkg/database"
)

func (db *sqliteDB) InsertMetadata(entries ...database.MetadataEntry) error {
	rowCount := len(entries)
	paramCount := 3

	valuesSQL := buildValuesSQL(rowCount, paramCount)
	sql := joinLines(
		`INSERT INTO metadata (`,
		`	table_name,`,
		`	field_name,`,
		`	description`,
		`) VALUES `+valuesSQL+";",
	)

	var params []any
	for _, v := range entries {
		params = append(params, v.Table, v.Field, v.Description)
	}

	if e := db.insert(sql, params); e != nil {
		return ErrSQLite.Wrap(e)
	}

	return nil
}

func (db *sqliteDB) AllMetadata() ([]database.MetadataEntry, error) {
	sql := joinLines(
		`SELECT`,
		`	table_name,`,
		`	field_name,`,
		`	description`,
		`FROM`,
		`	metadata;`,
	)

	rows, e := db.conn.Query(sql)
	if e != nil {
		e = database.ErrQuerying.Wrap(e)
		return nil, ErrSQLite.Wrap(e)
	}
	defer rows.Close()

	return scanMetadataRows(rows)
}

func scanMetadataRows(rows *sql.Rows) ([]database.MetadataEntry, error) {
	var results []database.MetadataEntry

	for rows.Next() {
		var entry database.MetadataEntry

		e := rows.Scan(&entry.Table, &entry.Field, &entry.Description)
		if e != nil {
			e = database.ErrParsing.Wrap(e)
			return nil, ErrSQLite.Wrap(e)
		}

		results = append(results, entry)
	}

	return results, nil
}
