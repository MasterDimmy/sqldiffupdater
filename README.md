## sqldiffupdater
The `sqldiffupdater.Generate()` is a Golang function that generates an SQL update statement based on the comparison of the properties of two objects, `newvar` and `oldvar`. The function takes three arguments: `tableName`, `newvar`, and `oldvar`. 

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