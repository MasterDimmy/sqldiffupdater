package sdu

import (
	"reflect"
	"testing"
)

func Test_Generate(t *testing.T) {
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
		//handle err
		t.Fatal()
	}

	// sql: "UPDATE users SET FirstName=:FirstName, Email=:Email WHERE Id=:Id"
	// here changed only FirstName and Email
	// values: map[string]interface{}{"FirstName": "John", "Email": "johndoe@example.com", "Id": 1}
	// db : your sql=db provider

	if sql != "UPDATE users SET FirstName=:FirstName, Email=:Email WHERE Id=:Id" {
		t.Fatal()
	}

	if !reflect.DeepEqual(values, map[string]interface{}{
		"Email":     "johndoe@example.com",
		"FirstName": "John",
		"Id":        1,
	}) {
		t.Fatal()
	}

	t.Logf("sql: %s\n", sql)
	t.Logf("values: %+v\n", values)
}

func Test_Update(t *testing.T) {
	type User struct {
		Id        int
		FirstName string
		LastName  string
		Email     string
	}

	newUser := User{Id: 1, FirstName: "John", LastName: "Doe", Email: "johndoe@example.com"}

	sql, values, err := Update("users", newUser, []string{"FirstName", "Email"})
	if err != nil {
		//handle err
		t.Fatal()
	}
	// sql: "UPDATE users SET FirstName=:FirstName, Email=:Email WHERE Id=:Id"
	// here changed only FirstName and Email
	// values: map[string]interface{}{"FirstName": "John", "Email": "johndoe@example.com", "Id": 1}
	// db : your sql=db provider

	if sql != "UPDATE users SET FirstName=:FirstName, Email=:Email WHERE Id=:Id" {
		t.Fatal()
	}

	if !reflect.DeepEqual(values, map[string]interface{}{
		"Email":     "johndoe@example.com",
		"FirstName": "John",
		"Id":        1,
	}) {
		t.Fatal()
	}

	t.Logf("sql: %s\n", sql)
	t.Logf("values: %+v\n", values)
}
