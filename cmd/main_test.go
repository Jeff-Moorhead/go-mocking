package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jeff-moorhead/go-mocking/friends"
)

// Remember: It's all just data.

var server *Server
var backend *MockFriendStore

func setUp() {
	backend = NewMockFriendStore()
	server = NewServer(backend)
}
func TestHandlePost(t *testing.T) {

	// Create a new request
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"Emily","age":25,"occupation":"Teacher"}`))

	rr := httptest.NewRecorder() // httptest.ResponseRecorder is an implementation of http.ResponseWriter
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Incorrect status code returned: got %v, expected %v", status, http.StatusCreated)
	}

	expectedreturn := `[{"name":"Emily","age":25,"occupation":"Teacher"}]`

	// rr.Body is the buffer where the ResponseWriter.Write() calls send their data to
	// Body.String() is a bytes.Buffer method that returns the buffer's contents as a string
	if body := rr.Body.String(); body != expectedreturn {
		t.Errorf("Unexpected body returned: got %#q, expected %#q", body, expectedreturn)
	}

	expectedcontents := []friends.Friend{
		{
			Name:       "Emily",
			Age:        25,
			Occupation: "Teacher",
		},
	}

	if !backend.CheckContentsEquality(expectedcontents) {
		t.Errorf("Unexpected backend contents: got %#q, expected %#q", backend.friends, expectedcontents)
	}
}

func TestHandleAlreadyExists(t *testing.T) {
	contents := []friends.Friend{
		{
			Name:       "Emily",
			Age:        25,
			Occupation: "Teacher",
		},
	}
	backend.SetContents(contents)

	// Create a new request
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"Emily","age":25,"occupation":"Teacher"}`))

	rr := httptest.NewRecorder() // httptest.ResponseRecorder is an implementation of http.ResponseWriter
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Incorrect status code returned: got %v, expected %v", status, http.StatusCreated)
	}

	expectedcontents := []friends.Friend{
		{
			Name:       "Emily",
			Age:        25,
			Occupation: "Teacher",
		},
	}

	if !backend.CheckContentsEquality(expectedcontents) {
		t.Errorf("Unexpected backend contents: got %q, expected %q", backend.friends, expectedcontents)
	}
}

func TestHandleGet(t *testing.T) {
	contents := []friends.Friend{
		{
			Name:       "Jeff",
			Age:        25,
			Occupation: "Programmer",
		},
	}
	backend.SetContents(contents)

	req := httptest.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Incorrect status code returned: got %v, expected %v", status, http.StatusOK)
	}

	expected := `[{"name":"Jeff","age":25,"occupation":"Programmer"}]`

	if body := rr.Body.String(); body != expected {
		t.Errorf("Unexpected body returned: got %#q, expected %#q", body, expected)
	}
}

func TestHandleDelete(t *testing.T) {
	contents := []friends.Friend{
		{
			Name:       "Jeff",
			Age:        25,
			Occupation: "Programmer",
		},
		{
			Name:       "Emily",
			Age:        25,
			Occupation: "Teacher",
		},
		{
			Name:       "Paul",
			Age:        59,
			Occupation: "Cabinet Maker",
		},
	}
	backend.SetContents(contents)

	req := httptest.NewRequest("DELETE", "/", nil)
	q := req.URL.Query()
	q.Add("name", "Emily")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("Incorrect status code returned: got %v, expected %v", status, http.StatusOK)
	}

	expected := []friends.Friend{
		{
			Name:       "Jeff",
			Age:        25,
			Occupation: "Programmer",
		},
		{
			Name:       "Paul",
			Age:        59,
			Occupation: "Cabinet Maker",
		},
	}
	if !backend.CheckContentsEquality(expected) {
		t.Errorf("Incorrect backend contents: got %v, expected %v", backend.friends, expected)
	}
}

func TestHandleInvalidDelete(t *testing.T) {
	contents := []friends.Friend{
		{
			Name:       "Jeff",
			Age:        25,
			Occupation: "Programmer",
		},
		{
			Name:       "Emily",
			Age:        25,
			Occupation: "Teacher",
		},
		{
			Name:       "Paul",
			Age:        59,
			Occupation: "Cabinet Maker",
		},
	}
	backend.SetContents(contents)

	req := httptest.NewRequest("DELETE", "/", nil)
	q := req.URL.Query()
	q.Add("name", "NONEXISTANT-NAME")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Incorrect status code returned: got %v, expected %v", status, http.StatusOK)
	}

	expected := []friends.Friend{
		{
			Name:       "Jeff",
			Age:        25,
			Occupation: "Programmer",
		},
		{
			Name:       "Emily",
			Age:        25,
			Occupation: "Teacher",
		},
		{
			Name:       "Paul",
			Age:        59,
			Occupation: "Cabinet Maker",
		},
	}
	if !backend.CheckContentsEquality(expected) {
		t.Errorf("Incorrect backend contents: got %v, expected %v", backend.friends, expected)
	}
}

func TestMain(m *testing.M) {
	setUp()
	m.Run()
}
