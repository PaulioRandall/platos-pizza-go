package database

import (
	"fmt"
	"strconv"
)

// OrderDetail represents a specific pizza type order, one or more pizzas,
// within an order
type OrderDetail struct {
	// Unique identifier for each pizza placed within each order
	// (pizzas of the same type and size are kept in the same row, and the quantity increases)
	Id int

	// Foreign key that ties the details in each order to the order itself
	OrderId int

	// Foreign key that ties the pizza ordered to its details, like size and price
	PizzaId string

	// Quantity ordered for each pizza of the same type and size
	Quantity int
}

func PrintOrderDetails(orderDetails []OrderDetail) {
	fmt.Println("[Order details]")
	fmt.Println(`"ID", "Order ID", "Pizza ID", "Quantity"`)
	for _, v := range orderDetails {
		fmt.Printf("%d, %d, %q, %d\n", v.Id, v.OrderId, v.PizzaId, v.Quantity)
	}
}

func QueryPrintOrderDetails(db PlatosPizzaDatabase) error {
	records, e := db.HeadOrderDetails()

	if e != nil {
		return ErrDatabase.CausedBy(e, "database.QueryPrintOrderDetails")
	}

	PrintOrderDetails(records)
	fmt.Println("...")

	return nil
}

func InsertOrderDetailsFromCSV(db PlatosPizzaDatabase, filename string) error {
	records, e := readCSV(filename)
	if e != nil {
		return ErrDatabase.CausedBy(e, "Failed to read order details %q", filename)
	}

	for i, record := range records {
		id, e := strconv.Atoi(record[0])
		if e != nil {
			return ErrDatabase.CausedBy(e, "Bad order details ID discovered")
		}

		orderId, e := strconv.Atoi(record[1])
		if e != nil {
			return ErrDatabase.CausedBy(e, "Bad order ID discovered")
		}

		quantity, e := strconv.Atoi(record[3])
		if e != nil {
			return ErrDatabase.CausedBy(e, "Bad quantity value discovered")
		}

		orderDetail := OrderDetail{
			Id:       id,
			OrderId:  orderId,
			PizzaId:  record[2],
			Quantity: quantity,
		}

		if e = db.InsertOrderDetails(orderDetail); e != nil {
			return ErrDatabase.CausedBy(e,
				"Failed to insert order detail at line %d", lineNumber(i),
			)
		}
	}

	return nil
}
