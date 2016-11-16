package tumblr

import (
	"net/url"
	"encoding/json"
	"fmt"
	"errors"
	"strconv"
)

// If you wish to use your own client, simply make sure it implements this interface
type ClientInterface interface {
	// Issue GET request to Tumblr API
	Get(endpoint string) (Response, error)
	// Issue GET request to Tumblr API with param values
	GetWithParams(endpoint string, params url.Values) (Response, error)
	// Issue POST request to Tumblr API
	Post(endpoint string) (Response, error)
	// Issue POST request to Tumblr API with param values
	PostWithParams(endpoint string, params url.Values) (Response, error)
	// Issue PUT request to Tumblr API
	Put(endpoint string) (Response, error)
	// Issue PUT request to Tumblr API with param values
	PutWithParams(endpoint string, params url.Values) (Response, error)
	// Issue DELETE request to Tumblr API
	Delete(endpoint string) (Response, error)
	// Issue DELETE request to Tumblr API with param values
	DeleteWithParams(endpoint string, params url.Values) (Response, error)
}

// shortcut for the most common case
func setPostId(id uint64, params url.Values) url.Values {
	return setParamsUint(id, params, "id")
}

// convenience function for setting an int
func setParamsUint(id uint64, params url.Values, key string) url.Values {
	params.Set(key, strconv.FormatUint(id, 10))
	return params
}

// Helper function to JSON stringify a given value
func jsonStringify(b interface{}) string {
	out, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Sprint("error:", err)
	}
	return string(out)
}

// Create a shallow copy of a params object
func copyParams(src url.Values) url.Values {
	dest := url.Values{}
	for k,v := range src {
		dest[k] = v
	}
	return dest
}

// Error returned for a collection's Next() invocation if no next page is possible
var NoNextPageError error = errors.New("No next page.")
// Error returned for a collection's Prev() invocation if no previous page is possible
var NoPrevPageError error = errors.New("No prev page.")
