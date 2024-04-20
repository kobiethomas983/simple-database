package database_manager

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/r3labs/diff/v3"

	"github.com/google/uuid"
)


type Address struct {
	Street string
	ZipCode int
	State string
}

type Users struct {
	ID string
	FirstName string
	LastName string
	Email string
	Address *Address
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func exists (path string) bool {
	_, err := os.Stat(path)
	if err == nil { return true}
	if os.IsNotExist(err) {return false}
	return false
}

const FILEPATH = "datastore/users"

type Driver struct {
	mu sync.Mutex
}


func (d *Driver) Create(user *Users) (*Users, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !exists(FILEPATH) {
		err := os.MkdirAll(FILEPATH, 0755)
		check(err)
	}

	user.ID = uuid.New().String()
	dataBytes, err := json.Marshal(user)
	check(err)

	if len(user.FirstName) < 1 {
		return nil, errors.New("user must have a first name")
	}

	fileName := fmt.Sprintf("%s/%s.json", FILEPATH, user.FirstName)
	if exists(fileName) {
		return nil, errors.New("this user already exists")
	}
	err = os.WriteFile(fileName, dataBytes, 0644)
	check(err)

	return user, nil
}

func (d *Driver) Update(user *Users) (*Users, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	filePath := fmt.Sprintf("%s/%s.json", FILEPATH, user.FirstName)
	if !exists(filePath) {
		return nil, errors.Wrapf(errors.New("does not exists"), "user %v", user.FirstName)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading file")
	}

	var userObject *Users
	err = json.Unmarshal(data, &userObject)
	if err != nil {
		return nil, errors.Wrap(err, "unserializing data error")
	}

	_, err = diff.Merge(userObject, user, &userObject)
	if err != nil {
		return nil, errors.Wrap(err, "error merging change difference")
	}

	
	dataBytes, err := json.Marshal(userObject)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(filePath, dataBytes, 0644)
	if err != nil {
		return nil, err
	}

	// err = file.Truncate(0)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "error while trying to empty file")
	// }

	// _, err = file.Seek(0,0)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "error while trying to seek")
	// }
	// _, err = fmt.Fprintf(file, "%v", dataBytes)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "error writing to file")
	// }

	return userObject, nil
}
