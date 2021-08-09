package main

import (
	"encoding/json"
	"fmt"

	"github.com/jeff-moorhead/go-mocking/friends"
)

// A stubbed implementation of the DataStore interface for testing
type MockFriendStore struct {
	friends []friends.Friend
}

func NewMockFriendStore() *MockFriendStore {
	return &MockFriendStore{
		make([]friends.Friend, 0),
	}
}

func (self *MockFriendStore) Marshal() ([]byte, error) {
	return json.Marshal(&self.friends)
}

func (self *MockFriendStore) Save() error {
	return nil
}

func (self *MockFriendStore) Refresh() error {
	return nil
}

func (self *MockFriendStore) Add(newfriend friends.Friend) error {
	for _, v := range self.friends {
		if v.Equals(newfriend) {
			return fmt.Errorf("Entry %q already exists", v)
		}
	}

	self.friends = append(self.friends, newfriend)
	return nil
}

func (self *MockFriendStore) Delete(name string) bool {
	for i, v := range self.friends {
		if v.Name == name {
			self.friends = append(self.friends[:i], self.friends[i+1:]...)
			return true
		}
	}
	return false
}

func (self *MockFriendStore) CheckContentsEquality(expected []friends.Friend) bool {
	if len(self.friends) != len(expected) {
		return false
	}

	for i, v := range self.friends {
		if v != expected[i] {
			return false
		}
	}

	return true
}

func (self *MockFriendStore) SetContents(values []friends.Friend) {
	self.friends = values
}
