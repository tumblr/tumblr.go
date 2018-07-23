package tumblr

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestGetLikesClientErrorReturnsError(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("", clientErr)
	if _, err := GetLikes(client, url.Values{}); err != clientErr {
		t.Fatal("Client error should be returned")
	}
}

func TestGetLikesJSONErrorReturnsError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := GetLikes(client, url.Values{}); err == nil {
		t.Fatal("JSON unmarshal error should be returned")
	}
}

func TestGetLikesSuccess(t *testing.T) {
	client := newTestClient("{}", nil)
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"GetLikes",
		http.MethodGet,
		"/user/likes",
		params,
	)
	if response, err := GetLikes(client, params); err != nil || response == nil {
		t.Fatal("Request should succeed")
	} else if response.response == nil {
		t.Fatal("Response should be set")
	} else if response.client != client {
		t.Fatal("Response should set client")
	}
}

func TestDoLikeError(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("{}", clientErr)
	path := "/any/given/path"
	var postId uint64 = 1986
	reblogKey := "some-key"
	params := setPostId(postId, url.Values{})
	params.Set("reblog_key", reblogKey)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"doLike",
		http.MethodPost,
		path,
		params,
	)
	if err := doLike(client, path, postId, reblogKey); err != clientErr {
		t.Fatal("Client error should be returned")
	}
}

func TestDoLikeSuccess(t *testing.T) {
	client := newTestClient("{}", nil)
	path := "/any/given/path"
	var postId uint64 = 1986
	reblogKey := "some-key"
	params := setPostId(postId, url.Values{})
	params.Set("reblog_key", reblogKey)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"doLike",
		http.MethodPost,
		path,
		params,
	)
	if err := doLike(client, path, postId, reblogKey); err != nil {
		t.Fatal("No error should be generated")
	}
}

func TestLike(t *testing.T) {
	client := newTestClient("{}", nil)
	var postId uint64 = 1986
	reblogKey := "some-key"
	params := setPostId(postId, url.Values{})
	params.Set("reblog_key", reblogKey)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"doLike",
		http.MethodPost,
		"/user/like",
		params,
	)
	if err := LikePost(client, postId, reblogKey); err != nil {
		t.Fatal("No error should be generated")
	}
}

func TestUnlike(t *testing.T) {
	client := newTestClient("{}", nil)
	var postId uint64 = 1986
	reblogKey := "some-key"
	params := setPostId(postId, url.Values{})
	params.Set("reblog_key", reblogKey)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"doLike",
		http.MethodPost,
		"/user/unlike",
		params,
	)
	if err := UnlikePost(client, postId, reblogKey); err != nil {
		t.Fatal("No error should be generated")
	}
}

func TestFullLikesWithJsonError(t *testing.T) {
	client := newTestClient("{}", nil)
	response, _ := GetLikes(client, url.Values{})
	if response == nil {
		t.Fatal("Unable to get likes")
	}
	if response.parsedPosts != nil {
		t.Fatal("Parsed posts should start out uninitialized")
	}
	response.response.body = []byte("{")
	_, err := response.Full()
	if err == nil {
		t.Fatal("JSON Unmarshal failure should be returned")
	}

}

func TestFullLikesSuccess(t *testing.T) {
	client := newTestClient("{}", nil)
	response, _ := GetLikes(client, url.Values{})
	if response == nil {
		t.Fatal("Unable to get likes")
	}
	if response.parsedPosts != nil {
		t.Fatal("Parsed posts should start out uninitialized")
	}
	_, err := response.Full()
	if err != nil {
		t.Fatal("Full like posts should be returned")
	}

}
