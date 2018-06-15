package tumblr

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestTaggedSearchFailsWithClientError(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("", clientErr)
	if _, err := TaggedSearch(client, "", url.Values{}); err != clientErr {
		t.Fatal("Client error should be returned")
	}
}

func TestTaggedSearchFailsWithJsonError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := TaggedSearch(client, "", url.Values{}); err == nil {
		t.Fatal("JSON error should be returned")
	}
}

func TestTaggedSearchSuccess(t *testing.T) {
	client := newTestClient("{}", nil)
	tag := "some-tag"
	params := url.Values{}
	params.Set("akey", "aval")
	paramsCopy := copyParams(params)
	params.Set("tag", tag)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"TaggedSearch",
		http.MethodGet,
		"/tagged",
		params,
	)
	if _, err := TaggedSearch(client, tag, paramsCopy); err != nil {
		t.Fatal("JSON error should be returned")
	}
}

func TestSearchResults_Next(t *testing.T) {
	client := newTestClient("{}", nil)
	tag := "some-tag"
	result, _ := TaggedSearch(client, tag, url.Values{})
	if _, err := result.Next(); err != NoNextPageError {
		t.Fatal("Next on empty results should generate error")
	}
	lastPost := Post{Timestamp: 123456}
	result.Posts = []PostInterface{&Post{}, &Post{}, &lastPost}
	params := url.Values{}
	params.Set("before", "123456")
	params.Set("tag", tag)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"SearchResults.Next",
		http.MethodGet,
		"/tagged",
		params,
	)
	if _, err := result.Next(); err != nil {
		t.Fatal("Failed to retrieve next page of results")
	}
	lastPost = Post{Timestamp: 123456, FeaturedTimestamp: 7891011}
	result.Posts = []PostInterface{&Post{}, &Post{}, &lastPost}
	params = url.Values{}
	params.Set("before", "7891011")
	params.Set("tag", tag)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"SearchResults.Next",
		http.MethodGet,
		"/tagged",
		params,
	)
	if _, err := result.Next(); err != nil {
		t.Fatal("Failed to use the correct timestamp")
	}

}
