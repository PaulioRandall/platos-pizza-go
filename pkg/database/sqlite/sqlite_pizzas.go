package sqlite

import (
	"database/sql"

	"github.com/PaulioRandall/platos-pizza-go/pkg/database"
)

func (db *sqliteDB) InsertPizzas(pizzas ...database.Pizza) error {
	buildPizzasInsertSQL := func(batch []database.Pizza) (sql string, params []any) {
		rowCount := len(batch)
		paramCount := 4

		valuesSQL := buildValuesSQL(rowCount, paramCount)
		sql = joinLines(
			`INSERT INTO pizzas (`,
			`	id,`,
			`	type_id,`,
			`	size,`,
			`	price`,
			`) VALUES `+valuesSQL+";",
		)

		for _, v := range batch {
			params = append(params, v.Id, v.TypeId, v.Size, v.Price)
		}

		return sql, params
	}

	return sqlitePartitionedInsert(db, pizzas, buildPizzasInsertSQL)
}

func (db *sqliteDB) HeadPizzas() ([]database.Pizza, error) {
	sql := joinLines(
		`SELECT`,
		`	id,`,
		`	type_id,`,
		`	size,`,
		`	price`,
		`FROM`,
		`	pizzas`,
		`LIMIT ?;`,
	)

	rows, e := db.conn.Query(sql, queryHeadMax)
	if e != nil {
		e = database.ErrQuerying.CausedBy(e, "Querying pizzas")
		return nil, ErrSQLite.Wrap(e)
	}
	defer rows.Close()

	return scanPizzaRows(rows)
}

func scanPizzaRows(rows *sql.Rows) ([]database.Pizza, error) {
	var results []database.Pizza

	for rows.Next() {
		var pizza database.Pizza

		e := rows.Scan(
			&pizza.Id,
			&pizza.TypeId,
			&pizza.Size,
			&pizza.Price,
		)

		if e != nil {
			e = database.ErrParsing.CausedBy(e, "Row scanning failed")
			return nil, ErrSQLite.Wrap(e)
		}

		results = append(results, pizza)
	}

	return results, nil
}
