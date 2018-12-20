package csvreader

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	lastValueFirstColumn string
}

func (ts *testStruct) setFirstColumnValue(index int, rowData []string) error {
	ts.lastValueFirstColumn = rowData[0]
	return nil
}

// TestStructMethodCanBePassed is a pretty lame test. I was not sure if this would work
// so before I went super deep into development I'd rather write a throwaway test.
func TestStructMethodCanBePassed(t *testing.T) {
	in := `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"`
	r := csv.NewReader(strings.NewReader(in))
	ts := testStruct{}
	iterateLines(r, ts.setFirstColumnValue)
	assert.Equal(t, "Robert", ts.lastValueFirstColumn)
}

func TestFailsOnBadPath(t *testing.T) {
	_, err := readFileToCSVReader("/badpath/does/not/exist.csv")

	assert.Error(t, err)
}
