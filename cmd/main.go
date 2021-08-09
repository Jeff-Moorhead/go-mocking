package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jeff-moorhead/go-mocking/friends"
)

// Takes any DataStore implementation
type Server struct {
	friendstore friends.DataStore
}

func NewServer(datastore friends.DataStore) *Server {
	server := &Server{
		datastore,
	}
	err := server.friendstore.Refresh()
	if err != nil {
		log.Fatalf("Unable to populate datastore: %q", err)
	}
	return server
}

func DecodeRequest(r *http.Request, f *friends.Friend) error {
	return json.NewDecoder(r.Body).Decode(f)
}

func (self *Server) handleGet(w http.ResponseWriter, r *http.Request) {

	err := self.friendstore.Refresh()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":%#q}`, err)
		return
	}

	b, err := self.friendstore.Marshal()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":%#q}`, err)
		return
	}

	fmt.Fprint(w, string(b))
}

func (self *Server) handlePost(w http.ResponseWriter, r *http.Request) {

	var newfriend friends.Friend

	err := DecodeRequest(r, &newfriend)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":%#q}`, err)
		return
	}

	err = self.friendstore.Add(newfriend)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, `{"error":%#q}`, err)
		return
	}

	err = self.friendstore.Save()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":%#q}`, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	b, err := self.friendstore.Marshal()
	if err != nil {
		fmt.Fprintf(w, `{"error":%#q}`, err)
		return
	}

	fmt.Fprint(w, string(b))
}

func (self *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := q.Get("name")

	deleted := self.friendstore.Delete(name)
	if !deleted {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"Unable to delete, no record exists with name %q"}`, name)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (self *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		self.handleGet(w, r)
	case http.MethodPost:
		self.handlePost(w, r)
	case http.MethodDelete:
		self.handleDelete(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func main() {
	fs := friends.NewFriendStore("friends.dat")
	server := NewServer(fs)
	http.Handle("/", server)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
