package database

import (
	"fmt"

	"github.com/PaulioRandall/go-trackerr"
)

const (
	QueryHeadMax   = 8
	DatetimeFormat = "2006-01-02 15:04:05"
)

var (
	ErrDatabase = trackerr.Checkpoint("Database error")

	ErrCreating  = trackerr.Track("Failed to create database or tables")
	ErrPreparing = trackerr.Track("Failed to prepare database query or statement")
	ErrInserting = trackerr.Track("Failed to execute data insert into database")
	ErrQuerying  = trackerr.Track("Failed to execute query on database")
	ErrParsing   = trackerr.Track("Failed to read or parse database results")
	ErrPrinting  = trackerr.Track("Failed to print rows from database")
	ErrClosed    = trackerr.Track("Can't execute requests on a closed database")

	ErrCSVFile = trackerr.Track("Error handling CSV file")
)

// PlatosPizzaDatabase represents an interface to a database of orders, pizzas,
// and information useful for analysing Plato's Pizzeria customer buying
// habits.
type PlatosPizzaDatabase interface {
	InsertMetadata(...MetadataEntry) error
	InsertOrders(...Order) error
	InsertOrderDetails(...OrderDetail) error
	InsertPizzas(...Pizza) error
	InsertPizzaTypes(...PizzaType) error

	AllMetadata() ([]MetadataEntry, error)
	HeadOrders() ([]Order, error)
	HeadOrderDetails() ([]OrderDetail, error)
	HeadPizzas() ([]Pizza, error)
	HeadPizzaTypes() ([]PizzaType, error)

	Close()
}

func Print(db PlatosPizzaDatabase) error {
	queryPrintFuncs := []func(PlatosPizzaDatabase) error{
		QueryPrintMetadata,
		QueryPrintOrders,
		QueryPrintOrderDetails,
		QueryPrintPizzas,
		QueryPrintPizzaTypes,
	}

	for i, f := range queryPrintFuncs {
		if i != 0 {
			fmt.Println()
		}

		if e := f(db); e != nil {
			return ErrPrinting.Wrap(e)
		}
	}

	return nil
}
