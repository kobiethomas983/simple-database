package database_manager

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"

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

	//Get first 8 digits of uuid
	user.ID = strings.Split(uuid.New().String(), "-")[0]
	dataBytes, err := json.Marshal(user)
	check(err)

	fileName := fmt.Sprintf("%s/%s.json", FILEPATH, user.ID)
	if exists(fileName) {
		return nil, errors.New("this user already exists")
	}
	err = os.WriteFile(fileName, dataBytes, 0644)
	check(err)

	return user, nil
}

func (d *Driver) Update(fieldUpdates map[string]interface{}, primaryID, collection string) (*Users, error) {
	switch collection {
	case "users":
		user, err := d.GetByID(primaryID)
		if err != nil {
			return nil, err
		}

		d.mu.Lock()
		defer d.mu.Unlock()

		sourceObject := reflect.ValueOf(&user).Elem()
		if sourceObject.Kind() == reflect.Ptr {
			sourceObject = sourceObject.Elem()
		}
		for field, value := range fieldUpdates {
			objectField := sourceObject.FieldByName(field)
			if !objectField.IsValid() {
				return nil, errors.New("field doesn't exist")
			}
			if objectField.CanSet() {
				fieldValue := reflect.ValueOf(value)
				objectField.Set(fieldValue)
			}
		}

		filePath := fmt.Sprintf("%s/%s.json", FILEPATH, user.ID)
		databytes, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(filePath, databytes, 0644)
		if err != nil {
			return nil, errors.Wrap(err, "error writing update")
		}
		return user, nil

	default:
		return nil, errors.New("unsupported collection")
	}
}

func (d *Driver) GetByID(ID string) (*Users, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	filePath := fmt.Sprintf("%s/%s.json", FILEPATH, ID)
	if !exists(filePath) {
		return nil, errors.Wrapf(errors.New("does not exists"), "user %s", ID)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading file")
	}

	var user *Users
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, errors.Wrap(err, "unserializing data error")
	}

	return user, nil
}

func (d *Driver) Delete(ID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	filePath := fmt.Sprintf("%s/%s.json", FILEPATH, ID)
	err := os.Remove(filePath)
	if err != nil {
		return errors.New("error remove row")
	}
	return nil
}
