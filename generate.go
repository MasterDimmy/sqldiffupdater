package sdu

/*
 */
import (
	"fmt"
	"reflect"
	"strings"
)

func Generate(tableName string, newvar interface{}, oldvar interface{}) (string, map[string]interface{}, error) {
	// Get the type of the new and old variables
	newType := reflect.TypeOf(newvar)
	oldType := reflect.TypeOf(oldvar)

	// Make sure that the types are the same
	if newType != oldType {
		return "", nil, fmt.Errorf("newvar and oldvar must be of the same type")
	}

	// Get the values of the new and old variables
	newValue := reflect.ValueOf(newvar)
	oldValue := reflect.ValueOf(oldvar)

	// Make sure that the values are valid
	if !newValue.IsValid() || !oldValue.IsValid() {
		return "", nil, fmt.Errorf("newvar and oldvar must be valid")
	}

	// Get the primary key field
	_, ok := newType.FieldByName("Id")
	if !ok {
		return "", nil, fmt.Errorf("newvar and oldvar must have an 'Id' field")
	}

	// Check each field of the struct for differences
	var setValues []string
	values := make(map[string]interface{})
	for i := 0; i < newType.NumField(); i++ {
		field := newType.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Get the values of the field in the new and old variables
		newFieldValue := newValue.Field(i)
		oldFieldValue := oldValue.Field(i)

		// Skip fields that are equal
		if reflect.DeepEqual(newFieldValue.Interface(), oldFieldValue.Interface()) {
			continue
		}

		// Get the name of the field
		fieldName := field.Name

		// Skip the Id field
		if fieldName == "Id" {
			continue
		}

		// Add the field to the setValues slice
		setValues = append(setValues, fmt.Sprintf("%s=:%s", fieldName, fieldName))

		// Add the value to the values map
		values[fieldName] = newFieldValue.Interface()
	}

	// Construct the SQL query
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE Id=:Id", tableName, strings.Join(setValues, ", "))

	// Add the Id value to the values map
	idFieldValue := newValue.FieldByName("Id")
	values["Id"] = idFieldValue.Interface()

	// Return the SQL query and the values map
	return sql, values, nil
}
