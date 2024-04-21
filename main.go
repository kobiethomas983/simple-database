package main

import (
	"github.com/kobie/simple-database/database_manager"
	"fmt"
)


func main() {
	dbDriver := database_manager.Driver{}

	user := &database_manager.Users{
		FirstName: "lamar",
		LastName: "jackson",
		Email: "lamar@gmail.com",
		Address: &database_manager.Address{
			Street: "123 road",
			ZipCode: 1123,
			State: "New Jersey",
		},
	}

	user, err := dbDriver.Create(user)
	if err != nil {
		fmt.Printf("error creating user: %v\n", err)
	}

	updates := map[string]interface{}{
		"FirstName": "Jamal",
		"LastName": "Brown",
		"Email": "lt@gmail.com",
	}

	_, err = dbDriver.Update(updates, user.ID, "users")
	if err != nil {
		fmt.Println("error happened", err)
	}

	user, err = dbDriver.GetByID(user.ID)
	if err != nil {
		fmt.Printf("error updating: %v\n", err)
	} else {
		fmt.Printf("user: %+v\n", user)
	}
}