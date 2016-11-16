package tumblr

import (
	"testing"
	"net/http"
)

func TestNewResponse(t *testing.T) {
	body := []byte("This is a byte arra")
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	r := NewResponse(body, headers)
	if string(r.body) != string(body) || string(r.GetBody()) != string(body) {
		t.Fatal("NewResponse should generate response with the correct body data", r.body, body)
	}
	if len(r.Headers) != len(headers) {
		t.Fatal("NewResponse should generate response with the correct header data")
	}
	for k,_ := range headers {
		if r.Headers.Get(k) != headers.Get(k) {
			t.Fatal("NewResponse should generate response with identical header data")
		}
	}
}

func TestResponse_PopulateFromBody(t *testing.T) {
	r := NewResponse([]byte{}, nil)
	err := r.PopulateFromBody()
	if err == nil {
		t.Fatal("Populate from body should return error if body is empty")
	}
	if r.Meta != nil || r.Errors != nil || r.Result != nil {
		t.Fatal("Populate from body error should not change Meta, Errors, or Result")
	}
	r.body = []byte("{\"meta\": {}, \"response\": {}}")
	err = r.PopulateFromBody()
	if err != nil {
		t.Fatal("Populate from body should succeed with appropriate body")
	}
	r.body = []byte("{")
	if err = r.PopulateFromBody(); err != nil {
		t.Fatal("Populate from body should succeed after successful parse")
	}
	r.Result = nil
	r.Meta = nil
	r.Errors = nil
	if err = r.PopulateFromBody(); err == nil {
		t.Fatal("Populate from body should return unmarshal error on invalid JSON")
	}
}