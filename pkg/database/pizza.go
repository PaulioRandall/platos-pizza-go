package database

import (
	"fmt"
	"strconv"
)

type Pizza struct {
	// Unique identifier for each pizza (constituted by its type and size)
	Id string

	// Foreign key that ties each pizza to its broader pizza type
	TypeId string

	// Size of the pizza (Small, Medium, Large, X Large, or XX Large)
	Size string

	// Price of the pizza in USD
	Price float64
}

func PrintPizzas(pizzas []Pizza) {
	fmt.Println("[Pizzas]")
	fmt.Println(`"ID", "Type ID", "Size", "Price"`)
	for _, v := range pizzas {
		fmt.Printf("%q, %q, %q, %.2f\n", v.Id, v.TypeId, v.Size, v.Price)
	}
}

func QueryPrintPizzas(db PlatosPizzaDatabase) error {
	records, e := db.HeadPizzas()

	if e != nil {
		return ErrDatabase.CausedBy(e, "database.QueryPrintPizzas")
	}

	PrintPizzas(records)
	fmt.Println("...")

	return nil
}

func InsertPizzasFromCSV(db PlatosPizzaDatabase, filename string) error {
	records, e := readCSV(filename)
	if e != nil {
		return ErrDatabase.CausedBy(e, "Failed to read pizzas %q", filename)
	}

	for i, record := range records {
		price, e := strconv.ParseFloat(record[3], 64)
		if e != nil {
			return ErrDatabase.CausedBy(e, "Bad price value discovered")
		}

		pizza := Pizza{
			Id:     record[0],
			TypeId: record[1],
			Size:   record[2],
			Price:  price,
		}

		if e = db.InsertPizzas(pizza); e != nil {
			return ErrDatabase.CausedBy(e, "Failed to insert pizza at line %d", lineNumber(i))
		}
	}

	return nil
}
