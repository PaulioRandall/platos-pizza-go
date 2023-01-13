package sqlite

import (
	"database/sql"
	"time"

	"github.com/PaulioRandall/platos-pizzas-go/pkg/database"
)

func (db *sqliteDB) InsertOrders(orders ...database.Order) error {
	buildOrdersInsertSQL := func(batch []database.Order) (sql string, params []any) {
		rowCount := len(batch)
		paramCount := 2

		valuesSQL := buildValuesSQL(rowCount, paramCount)
		sql = joinLines(
			`INSERT INTO orders (`,
			`	id,`,
			`	datetime`,
			`) VALUES `+valuesSQL+";",
		)

		for _, v := range batch {
			params = append(params, v.Id, v.Datetime)
		}

		return sql, params
	}

	return sqlitePartitionedInsert(db, orders, buildOrdersInsertSQL)
}

func (db *sqliteDB) HeadOrders() ([]database.Order, error) {
	sql := joinLines(
		`SELECT`,
		`	id,`,
		`	strftime('%Y-%m-%d %H:%M:%S', orders.datetime) AS datetime`,
		`FROM`,
		`	orders`,
		// Could have used the value directly as SQL injection is not possible here
		// But it does mean the SQL driver will handle the type conversion for me
		`LIMIT ?;`,
	)

	rows, e := db.conn.Query(sql, queryHeadMax)
	if e != nil {
		e = database.ErrQuerying.CausedBy(e, "Querying orders")
		return nil, ErrSQLite.Wrap(e)
	}
	defer rows.Close()

	return scanOrderRows(rows)
}

func scanOrderRows(rows *sql.Rows) ([]database.Order, error) {
	var results []database.Order

	for rows.Next() {
		var order database.Order
		var datetimeStr string

		e := rows.Scan(&order.Id, &datetimeStr)
		if e != nil {
			return nil, database.ErrParsing.Wrap(e)
		}

		order.Datetime, e = time.Parse(database.DatetimeFormat, datetimeStr)
		if e != nil {
			e = database.ErrParsing.Wrap(e)
			return nil, ErrSQLite.Wrap(e)
		}

		results = append(results, order)
	}

	return results, nil
}
