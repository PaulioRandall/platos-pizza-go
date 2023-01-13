package database

import (
	"fmt"
)

type PizzaType struct {
	// Unique identifier for each pizza type
	Id string

	// Name of the pizza as shown in the menu
	Name string

	// Category that the pizza fall under in the menu
	// (Classic, Chicken, Supreme, or Veggie)
	Category string

	// Comma-delimited ingredients used in the pizza as shown in the menu
	// (they all include Mozzarella Cheese, even if not specified; and they all
	// include Tomato Sauce, unless another sauce is specified)
	Ingredients string
}

func PrintPizzaTypes(pizzaTypes []PizzaType) {
	fmt.Println("[Pizza types]")
	fmt.Println(`"ID", "Name", "Category", "Ingredients"`)
	for _, v := range pizzaTypes {
		fmt.Printf("%q, %q, %q, %q\n", v.Id, v.Name, v.Category, v.Ingredients)
	}
}

func QueryPrintPizzaTypes(db PlatosPizzaDatabase) error {
	records, e := db.HeadPizzaTypes()

	if e != nil {
		return ErrDatabase.CausedBy(e, "database.QueryPrintPizzaTypes")
	}

	PrintPizzaTypes(records)
	fmt.Println("...")

	return nil
}

func InsertPizzaTypesFromCSV(db PlatosPizzaDatabase, filename string) error {
	records, e := readCSV(filename)
	if e != nil {
		return ErrDatabase.CausedBy(e, "Failed to read pizza types %q", filename)
	}

	for i, record := range records {
		pizzaType := PizzaType{
			Id:          record[0],
			Name:        record[1],
			Category:    record[2],
			Ingredients: record[3],
		}

		if e = db.InsertPizzaTypes(pizzaType); e != nil {
			return ErrDatabase.CausedBy(e,
				"Failed to insert pizza type at line %d", lineNumber(i),
			)
		}
	}

	return nil
}
