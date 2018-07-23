package tumblr

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestGetUserInfoClientError(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("{}", clientErr)
	if _, err := GetUserInfo(client); err != clientErr {
		t.Fatal("User info must return client error")
	}
}
func TestGetUserInfoJsonError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := GetUserInfo(client); err == nil {
		t.Fatal("User info must return json unmarshal error")
	}
}
func TestGetUserInfo(t *testing.T) {
	client := newTestClient("{}", nil)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"UserInfo",
		http.MethodGet,
		"/user/info",
		url.Values{},
	)
	if _, err := GetUserInfo(client); err != nil {
		t.Fatal("User info failed")
	}
}
