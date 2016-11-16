package tumblr

import (
	"encoding/json"
	"errors"
	"net/http"
)

// API Response structure which we'll use to pass back to behaviors
type Response struct {
	body []byte
	Headers http.Header
	Meta map[string]interface{} `json:"meta"`
	Result map[string]interface{} `json:"response"`
	Errors map[string]interface{} `json:"errors"`
}

// Create a response object from the body bytestream and the headers structure
func NewResponse(body []byte, headers http.Header) *Response {
	return &Response{body: body, Headers: headers}
}

// Get the raw response body
func (r *Response) GetBody() []byte {
	return r.body
}

// Utility function for populating the Response's fields
func (r *Response) PopulateFromBody() error {
	if len(r.body) < 1 {
		return errors.New("Unable to populate from empty body")
	}
	// already populated, don't do again
	if r.Meta != nil || r.Result != nil || r.Errors != nil {
		return nil
	}
	if e := json.Unmarshal(r.body, r); e != nil {
		return e
	}
	return nil
}