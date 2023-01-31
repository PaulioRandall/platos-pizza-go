package database

import (
	"encoding/csv"
	"os"
)

func lineNumber(i int) int {
	i++ // Convert from index to count
	i++ // Skip the header
	return i
}

func readCSV(filename string) ([][]string, error) {
	f, e := os.Open(filename)
	if e != nil {
		return nil, ErrCSVFile.CausedBy(e, "Could not open file %q", filename)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, e := r.ReadAll()

	if e != nil {
		return nil, ErrCSVFile.CausedBy(e, "Could not read file %q", filename)
	}

	records = records[1:] // Remove header
	return records, nil
}
