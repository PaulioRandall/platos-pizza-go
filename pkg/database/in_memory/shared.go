package in_memory

import (
	"github.com/PaulioRandall/platos-pizza-go/pkg/database"
)

type query[T any] func() ([]T, error)

func inMemoryHead[T any](db *inMemory, items []T) ([]T, error) {
	return inMemoryExecute(db, func() ([]T, error) {
		if len(items) < queryHeadMax {
			return items, nil
		}
		return items[0:queryHeadMax], nil
	})
}

func inMemoryExecute[T any](db *inMemory, q query[T]) ([]T, error) {
	if db.closed {
		return nil, ErrInMemory.Wrap(database.ErrClosed)
	}

	result, e := q()
	if e != nil {
		e = database.ErrQuerying.Wrap(e)
		e = ErrInMemory.Wrap(e)
	}

	return result, e
}

func inMemoryInsert(db *inMemory, f func()) error {
	if db.closed {
		return ErrInMemory.Wrap(database.ErrClosed)
	}

	f()
	return nil
}
