package friends

import (
	"encoding/json"
	"fmt"
	"os"
)

type DataStore interface {
	Refresh() error
	Save() error
	Marshal() ([]byte, error)
	Add(Friend) error
	Delete(name string) bool
}

// Implements DataStore
type FriendStore struct {
	friends []Friend
	file    string
}

func NewFriendStore(datafile string) *FriendStore {

	friends := FriendStore{
		make([]Friend, 0),
		datafile,
	}
	return &friends
}

func (self *FriendStore) Refresh() error {

	b, err := os.ReadFile(self.file)
	if err != nil && os.IsNotExist(err) {
		err := self.Save()
		if err != nil {
			return err
		}

		return nil

	} else if err != nil {
		return err
	}

	return json.Unmarshal(b, &self.friends)
}

func (self *FriendStore) Save() error {

	b, err := self.Marshal()
	if err != nil {
		return err
	}

	return os.WriteFile(self.file, b, 0644)
}

func (self *FriendStore) Marshal() ([]byte, error) {

	b, err := json.Marshal(self.friends)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (self *FriendStore) Add(newfriend Friend) error {

	for _, v := range self.friends {
		if v.Equals(newfriend) {
			return fmt.Errorf("Entity %#q already exists", newfriend.Name)
		}
	}

	// Need a pointer here to modify the receiver directly. Value just throws out the copy
	self.friends = append(self.friends, newfriend)
	return nil
}

func (self *FriendStore) Delete(name string) bool {
	for i, v := range self.friends {
		if v.Name == name {
			self.friends = append(self.friends[:i], self.friends[i+1:]...)
			return true
		}
	}

	return false
}
