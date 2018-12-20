package totable

import (
	"fmt"
	"strings"
	"time"

	"github.com/Benjar12/knock_challenge/model/schemaelement"
	"github.com/Benjar12/knock_challenge/service/mysqldb"
	"github.com/Benjar12/knock_challenge/util/csvreader"
)

// So first off I could not think of a great name for the package.
// Second this prob should not be in model as it's more impl logic.

// ToTable is created once per file. Originally I was going to hold all
// of the lines in memory, but eventually that would crash the process. I
// was trying reduce algorithmic complexity and IOPS by random sampling
// or holding it all in memory. The issue would have been if all but
// one value parses as int things will blow up.
type ToTable struct {
	tableName string
	schema    []schemaelement.SchemaElement
	headers   []string // This is used only on load rows. It was a bit of a short cut I took.
}

func (t *ToTable) findSchema(lineIndex int, rowData []string) error {
	for i, val := range rowData {
		// If it's the first line we create the intial schema records
		if lineIndex == 0 {
			sv, err := schemaelement.NewSchemaElement(val)
			if err != nil {
				return err
			}
			t.schema = append(t.schema, sv)
			continue
		}
		// Otherwise we evaluate the type of the field to determain what type it should
		// be in mysql.
		t.schema[i].UpdateType(val)
	}

	return nil
}

func (t *ToTable) createTable() error {
	sql := fmt.Sprintf("CREATE TABLE `%s` ( `id` int(11) unsigned NOT NULL AUTO_INCREMENT, ", t.tableName)
	for _, sc := range t.schema {
		sql += fmt.Sprintf("`%s` ", sc.Name)

		if sc.ElmType == "int" {
			sql += "int(11) DEFAULT NULL, "
		} else if sc.ElmType == "float" {
			sql += "float DEFAULT NULL, "
		} else if sc.ElmType == "bool" {
			sql += "tinyint(1) DEFAULT NULL, "
		} else {
			sql += "varchar(255) DEFAULT NULL, "
		}
	}

	sql += "PRIMARY KEY (`id`) )"

	err := mysqldb.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func (t *ToTable) loadRow(lineIndex int, rowData []string) error {
	// If the line is just the headers then return
	if lineIndex == 0 {
		t.headers = rowData
		return nil
	}

	sql := fmt.Sprintf("INSERT INTO `%s` (`%s`) ", t.tableName, strings.Join(t.headers, "`,`"))
	sql += fmt.Sprintf("VALUES (")

	for i, colDat := range rowData {
		colType := t.schema[i].ElmType
		if colDat == "" {
			sql += "NULL,"
		} else if colType == "int" || colType == "float" {
			sql += fmt.Sprintf("%s,", colDat)
		} else if colType == "bool" {
			if colDat == "true" {
				sql += "1,"
			} else {
				sql += "0,"
			}
		} else {
			sql += fmt.Sprintf("'%s',", colDat)
		}
	}
	sql = strings.TrimRight(sql, ",")

	sql += ");"

	err := mysqldb.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// ProcessFile is called after the api finishes saving the file to a tmp dir. It first
// reads through the file once to find the schema and any discrepancies. After that it
// creates the table. From there it loads the rows.
func ProcessFile(tableName string, filePath string) error {
	tn := time.Now().Format("2006_01_02_15_04_05")
	tt := ToTable{tableName: fmt.Sprintf("%s_%s", tableName, tn)}

	err := csvreader.ReadAndIterate(filePath, tt.findSchema)
	if err != nil {
		return err
	}

	err = tt.createTable()
	if err != nil {
		return err
	}

	err = csvreader.ReadAndIterate(filePath, tt.loadRow)
	if err != nil {
		return err
	}

	return nil
}
