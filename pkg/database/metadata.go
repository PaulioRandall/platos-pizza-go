package database

import (
	"fmt"
)

type MetadataEntry struct {
	Table       string
	Field       string
	Description string
}

func PrintMetadata(entries []MetadataEntry) {
	fmt.Println("[Metadata]")
	fmt.Println(`"Table", "Field", "Description"`)
	for _, entry := range entries {
		fmt.Printf("%q, %q, %q\n", entry.Table, entry.Field, entry.Description)
	}
}

func QueryPrintMetadata(db PlatosPizzaDatabase) error {
	records, e := db.AllMetadata()

	if e != nil {
		return ErrDatabase.CausedBy(e, "Quering all metadata")
	}

	PrintMetadata(records)
	return nil
}

func InsertMetadataFromCSV(db PlatosPizzaDatabase, filename string) error {
	records, e := readCSV(filename)
	if e != nil {
		return ErrDatabase.CausedBy(e, "Failed to read metadata %q", filename)
	}

	for i, record := range records {
		entry := MetadataEntry{
			Table:       record[0],
			Field:       record[1],
			Description: record[2],
		}

		e := db.InsertMetadata(entry)
		if e != nil {
			return ErrDatabase.CausedBy(e,
				"Failed to insert metadata record at line %d", lineNumber(i),
			)
		}
	}

	return nil
}
