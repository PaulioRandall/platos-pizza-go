package sqlite

import (
	"strings"
)

type sqlBuilder[T any] func([]T) (sql string, params []any)

func sqlitePartitionedInsert[T any](
	db *sqliteDB,
	items []T,
	buildInsertSQL sqlBuilder[T],
) error {
	for _, batch := range partition(items, insertBatchSize) {
		sql, params := buildInsertSQL(batch)

		if e := db.insert(sql, params); e != nil {
			return ErrSQLite.Wrap(e)
		}
	}

	return nil
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func buildValuesSQL(rowCount, paramCount int) string {
	sb := strings.Builder{}
	params := buildParamsSQL(paramCount)

	for i := 0; i < rowCount; i++ {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(params)
	}

	return sb.String()
}

func buildParamsSQL(paramCount int) string {
	sb := strings.Builder{}
	sb.WriteRune('(')

	for i := 0; i < paramCount; i++ {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteRune('?')
	}

	sb.WriteRune(')')
	return sb.String()
}

func partition[T any](items []T, batchSize int) [][]T {
	var batches [][]T
	var batch []T

	for _, v := range items {
		if len(batch) == batchSize {
			batches = append(batches, batch)
			batch = nil
		}

		batch = append(batch, v)
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}
