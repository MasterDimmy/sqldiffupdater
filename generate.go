package sqldiffupdater

/*
 */
import (
	"fmt"
	"reflect"
	"strings"
)

/*
Parameters:

	tableName - database table to update
	newvar	- new object
	oldvar	- old object of same type
	key - field name of the primary key in tableName

Result:

	sql string of "update tableName set ..... where tableName.Id = oldvar.Id"
	map of key:value properties of newvar changes from oldvar
	error if happens

Notes:

	If no changes found, map[string]interface{} has len of 1 (just Id)!
	newvar should have key property as primary key for update!
*/
func Generate(tableName string, key string, newvar_ interface{}, oldvar_ interface{}) (string, map[string]interface{}, error) {
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
	_, ok := newType.FieldByName(key)
	if !ok {
		return "", nil, fmt.Errorf("newvar and oldvar must have an 'Id' field")
	}

	// Check each field of the struct for differences
	var setValues []string
	values := make(map[string]interface{})

	compareFields(newValue, oldValue, reflect.TypeOf(newvar), &setValues, values)

	// Construct the SQL query
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s=:%s", tableName, strings.Join(setValues, ", "), key, key)

	// Add the Id value to the values map
	idFieldValue := newValue.FieldByName(key)
	values[key] = idFieldValue.Interface()

	// Return the SQL query and the values map
	return sql, values, nil
}

// `update` is a Go function that generates SQL code for updating specific fields in
// a database table based on a given object and a list of fields to update.
// The function takes the table name, the object containing the new values,
// and a list of fields to update. The object must have an key field as a primary key.
// The function returns a tuple containing the SQL code for updating,
// a map of updated fields and their values, and an error if one occurs.
func Update(tableName string, key string, newvar_ interface{}, fields []string) (string, map[string]interface{}, error) {
	// Check if newobj is a pointer
	newvar := newvar_
	if reflect.TypeOf(newvar_).Kind() == reflect.Ptr {
		newvar = reflect.ValueOf(newvar_).Elem().Interface()
	}

	// Check that the object has an "Id" field
	_, ok := reflect.TypeOf(newvar).FieldByName(key)
	if !ok {
		return "", nil, fmt.Errorf("newvar must have an " + key + " field")
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
	sql = strings.TrimSuffix(sql, ",") + " WHERE " + key + "=:" + key

	// Add the Id value to the values map
	values[key] = newVal.FieldByName(key).Interface()

	return sql, values, nil
}

func compareFields(newVal, oldVal reflect.Value, newType reflect.Type, setValues *[]string, values map[string]interface{}) {
	for i := 0; i < newVal.NumField(); i++ {
		field := newType.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Handle nested structures
		if field.Type.Kind() == reflect.Struct {
			compareFields(newVal.Field(i), oldVal.Field(i), field.Type, setValues, values)
			continue
		}

		// Compare field values
		newFieldValue := newVal.Field(i)
		oldFieldValue := oldVal.Field(i)

		if !reflect.DeepEqual(newFieldValue.Interface(), oldFieldValue.Interface()) {
			fieldName := field.Name
			*setValues = append(*setValues, fmt.Sprintf("%s=:%s", fieldName, fieldName))
			values[fieldName] = newFieldValue.Interface()
		}
	}
}
