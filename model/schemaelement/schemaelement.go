package schemaelement

import "strconv"

// SchemaElement will first initalized when the headers are read
// from the file. As we loop through all the rows ElmType be set to
// int, bool, float, or string.
type SchemaElement struct {
	Name    string
	ElmType string
}

// UpdateType is called for each value within a given column. Our
// default type is string. I'm sure there is a much better way to do
// this but for now lets just get something working.
func (s *SchemaElement) UpdateType(columnValue string) {
	if s.ElmType == "string" {
		return
	}

	// Check if currect column value is int
	_, err := strconv.ParseInt(columnValue, 10, 64)
	if (err == nil) && (s.ElmType == "" || s.ElmType == "int") {
		s.ElmType = "int"
		return
	}

	// Check if currect column value is float
	_, err = strconv.ParseFloat(columnValue, 64)
	if (err == nil) && (s.ElmType == "" || s.ElmType == "float") {
		s.ElmType = "float"
		return
	}

	// Check if currect column value is float
	// TODO: this does not handle string case or if 0,1 are being
	// used as bool. Fix that.
	if (columnValue == "true" || columnValue == "false") && (s.ElmType == "" || s.ElmType == "bool") {
		s.ElmType = "bool"
		return
	}

	// Default to string
	s.ElmType = "string"
}

// NewSchemaElement is just a sudo constructor. normally would do
// more stuff.
func NewSchemaElement(name string) (SchemaElement, error) {
	return SchemaElement{
		Name: name,
	}, nil
}
