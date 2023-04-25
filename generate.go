package sqldiffupdater

/*
 */
import (
	"fmt"
	"reflect"
	"strings"
)

func Generate(tableName string, newvar_ interface{}, oldvar_ interface{}) (string, map[string]interface{}, error) {
	// Check if newvar is a pointer
	newvar := newvar_
	if reflect.TypeOf(newvar_).Kind() == reflect.Ptr {
		newvar = reflect.ValueOf(newvar_).Elem().Interface()
	}

	// Check if oldvar is a pointer
	oldvar := oldvar_
	if reflect.TypeOf(oldvar_).Kind() == reflect.Ptr {
		oldvar = reflect.ValueOf(oldvar_).Elem().Interface()
	}

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

//`update` is a Go function that generates SQL code for updating specific fields in
// a database table based on a given object and a list of fields to update.
// The function takes the table name, the object containing the new values,
// and a list of fields to update. The object must have an `Id` field as a primary key.
// The function returns a tuple containing the SQL code for updating,
// a map of updated fields and their values, and an error if one occurs.
func Update(tableName string, newvar_ interface{}, fields []string) (string, map[string]interface{}, error) {
	// Check if newobj is a pointer
	newvar := newvar_
	if reflect.TypeOf(newvar_).Kind() == reflect.Ptr {
		newvar = reflect.ValueOf(newvar_).Elem().Interface()
	}

	// Check that the object has an "Id" field
	_, ok := reflect.TypeOf(newvar).FieldByName("Id")
	if !ok {
		return "", nil, fmt.Errorf("newvar must have an 'Id' field")
	}

	// Get the values of the fields that have changed
	values := make(map[string]interface{})
	newVal := reflect.ValueOf(newvar)
	for _, field := range fields {
		value := newVal.FieldByName(field).Interface()
		values[field] = value
	}

	// Build the update SQL statement
	sql := fmt.Sprintf("UPDATE %s SET", tableName)
	for _, field := range fields {
		sql += fmt.Sprintf(" %s=:%s,", field, field)
	}
	sql = strings.TrimSuffix(sql, ",") + " WHERE Id=:Id"

	// Add the Id value to the values map
	values["Id"] = newVal.FieldByName("Id").Interface()

	return sql, values, nil
}
