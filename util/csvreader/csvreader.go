package csvreader

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

// ieadFileToCSVReader is here to decouple filesystem operations from IterateLines.
// This will allow us to test IterateLines for functionality not just integrations.
// We could use a moke or fake, but for this code test I think it's good enough.
func readFileToCSVReader(filePath string) (*csv.Reader, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return csv.NewReader(bufio.NewReader(csvFile)), nil
}

// iterateLines is a helper function so we don't have to repeat boilerplate code.
// This function takes in a pointer to a reader and a function to handle each line.
func iterateLines(reader *csv.Reader, handler func(index int, nextRow []string) error) error {
	index := 0
	// Rather than use ReadAll i'm ++ index manually to keep our memory footprint low.
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		err = handler(index, line)
		if err != nil {
			return err
		}
		index++
	}

	return nil
}

// ReadAndIterate just combines the functions above.
func ReadAndIterate(filePath string, handler func(index int, rowData []string) error) error {
	r, err := readFileToCSVReader(filePath)
	if err != nil {
		return err
	}

	err = iterateLines(r, handler)
	if err != nil {
		return err
	}

	return nil
}
