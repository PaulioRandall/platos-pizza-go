package database

import (
	"fmt"
	"strconv"
	"time"
)

// Order represents an order of pizzas, one or many pizzas per order
type Order struct {
	// Unique identifier for each order placed by a table
	Id int

	// Date & time the order was placed
	// (entered into the system prior to cooking & serving)
	Datetime time.Time
}

func PrintOrders(orders []Order) {
	fmt.Println("[Orders]")
	fmt.Println(`"ID", "Datetime"`)
	for _, entry := range orders {
		fmt.Printf("%d, %q\n", entry.Id, entry.Datetime)
	}
}

func QueryPrintOrders(db PlatosPizzaDatabase) error {
	records, e := db.HeadOrders()

	if e != nil {
		return ErrDatabase.CausedBy(e, "database.QueryPrintOrders")
	}

	PrintOrders(records)
	fmt.Println("...")

	return nil
}

func InsertOrdersFromCSV(db PlatosPizzaDatabase, filename string) error {
	records, e := readCSV(filename)
	if e != nil {
		return ErrDatabase.CausedBy(e, "Failure to read orders %q", filename)
	}

	orders := make([]Order, len(records))
	for i, record := range records {
		id, e := strconv.Atoi(record[0])
		if e != nil {
			return ErrDatabase.CausedBy(e, "Bad order ID discovered")
		}

		datetime, e := time.Parse(DatetimeFormat, record[1]+" "+record[2])
		if e != nil {
			return ErrDatabase.CausedBy(e, "Bad order date or time discovered")
		}

		orders[i] = Order{
			Id:       id,
			Datetime: datetime,
		}
	}

	if e = db.InsertOrders(orders...); e != nil {
		return ErrDatabase.Wrap(e)
	}

	return nil
}
