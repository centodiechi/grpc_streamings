package unarystreaming

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger"
)

type Role int

const (
	Admin Role = iota
	Editor
	Viewer
)

var DataBase *StorageProvider

type StorageProvider struct {
	DB *badger.DB
}

type User struct {
	UID      string  `json:"uid"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Profile  Profile `json:"profile"`
	Role     Role    `json:"role"`
}

type Profile struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

func NewStorageProvider(dbPath string) (*StorageProvider, error) {
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return nil, err
	}
	return &StorageProvider{DB: db}, nil
}

func (sp *StorageProvider) CreateUser(user User) error {
	err := sp.DB.View(func(txn *badger.Txn) error {
		userKey := []byte("user/" + user.UID)
		_, err := txn.Get(userKey)
		return err
	})
	if err == nil {
		return fmt.Errorf("user with this email already exists")
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return sp.DB.Update(func(txn *badger.Txn) error {
		userKey := []byte("user/" + user.UID)
		return txn.Set(userKey, data)
	})
}

func (sp *StorageProvider) AuthenticateUser(creds string) error {
	err := sp.DB.View(func(txn *badger.Txn) error {
		userKey := []byte("user/" + creds)
		_, err := txn.Get(userKey)
		return err
	})
	if err != nil {
		return fmt.Errorf("invalid Username or Password")
	}
	return nil
}

func (sp *StorageProvider) Close() {
	sp.DB.Close()
}
