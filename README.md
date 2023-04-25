# sqldiffupdater

## Function: `Generate`

Generates an SQL update statement based on the comparison of the properties of two objects, `newvar` and `oldvar`. The function takes three arguments: `tableName`, `newvar`, and `oldvar`. 

The `tableName` argument is a string that represents the name of the database table that the update statement will be executed against. The `newvar` and `oldvar` arguments are the new and old versions of the object that will be compared.

The function uses reflection to get the names and values of the properties of the `newvar` and `oldvar` objects, and generates an SQL update statement based on the differences between the two objects. The function checks that the `newvar` and `oldvar` objects have the same type, and that they have an `Id` field as a primary key.

The function returns a string representing the generated SQL update statement, a map containing the values of the changed properties, and an error if one occurs during the execution of the function.

### Parameters

- `tableName` (string): The name of the database table that the update statement will be executed against.
- `newvar` (interface{}): The new version of the object that will be compared.
- `oldvar` (interface{}): The old version of the object that will be compared.

### Returns

- `string`: The generated SQL update statement.
- `map[string]interface{}`: A map containing the values of the changed properties.
- `error`: An error if one occurs during the execution of the function.

### Example

```go
type User struct {
    Id        int
    FirstName string
    LastName  string
    Email     string
}

newUser := User{Id: 1, FirstName: "John", LastName: "Doe", Email: "johndoe@example.com"}
oldUser := User{Id: 1, FirstName: "Jane", LastName: "Doe", Email: "janedoe@example.com"}

sql, values, err := Generate("users", newUser, oldUser)
if err != nil {
    // handle error
}

// sql: "UPDATE users SET FirstName=:FirstName, Email=:Email WHERE Id=:Id"
// here changed only FirstName and Email
// values: map[string]interface{}{"FirstName": "John", "Email": "johndoe@example.com", "Id": 1}
// db : your sql=db provider
err = db.Exec(sql, values)
if err != nil {
    // handle error
}

```


## Function: `Update`

Generates SQL code for updating specific fields in a database table based on a given object and a list of fields to update.

### Prototype

```go
func update(tableName string, newvar interface{}, fields []string) (string, map[string]interface{}, error)
```

### Parameters

- `tableName` (string): the name of the database table to update.
- `newvar` (interface{}): an object containing the new values to update in the database table.
- `fields` ([]string): a list of fields to update in the database table.

### Returns

- `string`: the SQL code for updating the specified fields in the database table.
- `map[string]interface{}`: a map of the updated fields and their new values.
- `error`: an error if one occurs.

### Notes

- The object provided as `newvar` must have an `Id` field as the primary key in the corresponding database table.
- The `fields` parameter should contain the names of the fields in the `newvar` object that need to be updated in the database table.
- This function can be used to efficiently update only the necessary fields in a database table, saving time and resources.

### Example

```go
import (
    "fmt"
    "github.com/MasterDimmy/sqldiffupdater
)

type User struct {
    Id        int
    FirstName string
    LastName  string
    Email     string
}

func main() {
    user := User{Id: 1, FirstName: "John", LastName: "Doe", Email: "johndoe@example.com"}
    fieldsToUpdate := []string{"FirstName", "Email"}

    sql, updatedFields, err := repo.Update("users", user, fieldsToUpdate)
    if err != nil {
        panic(err)
    }

    fmt.Println(sql)
    fmt.Println(updatedFields)
}
```

Output:
```
UPDATE users SET FirstName=:FirstName, Email=:Email WHERE Id=:Id
map[FirstName:John Email:johndoe@example.com]
```
