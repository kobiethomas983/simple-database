package main

import (
	"encoding/json"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Address struct {
	Street string
	ZipCode int
	State string
}

type Users struct {
	FirstName string
	LastName string
	Email string
	Address *Address
}

func main() {
	err := os.MkdirAll("datastore/users", 0755 )
	check(err)

	data := &Users{
		FirstName: "lamar",
		LastName: "jackson",
		Email: "lamar@gmail.com",
		Address: &Address{
			Street: "123 road",
			ZipCode: 1123,
			State: "New Jersey",
		},
	}

	marshalledData, err := json.Marshal(data)
	check(err)

	createEmptyFile := func(fileName string) {
		check(os.WriteFile(fileName, marshalledData, 0644))
	}

	createEmptyFile("datastore/users/lamar")
}