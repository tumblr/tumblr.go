package tumblr

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestGetFollowingFail(t *testing.T) {
	mockError := errors.New("Mock error")
	client := newTestClient("{}", mockError)
	if _, err := GetFollowing(client, 0, 0); err != mockError {
		t.Fatal("Failed request does not return correct error")
	}
}

func TestGetFollowingJsonError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := GetFollowing(client, 0, 0); err == nil {
		t.Fatal("Bad JSON response does not generate an error")
	}
}

func TestGetFollowingSuccess(t *testing.T) {
	client := newTestClient("{}", nil)
	expectedParams := url.Values{}
	expectedParams.Set("offset", "0")
	expectedParams.Set("limit", "0")
	client.confirmExpectedSet = expectClientCallParams(t, "GetFollowing", http.MethodGet, "/user/following", expectedParams)
	if result, err := GetFollowing(client, 0, 0); err != nil {
		t.Fatal("Failed to get from following", err)
	} else {
		if result.offset != 0 {
			t.Fatal("Offset should have been 0")
		}
		if result.limit != 0 {
			t.Fatal("Limit should have been 0")
		}
	}

	// test correct params go in the proper place
	expectedParams.Set("offset", "2")
	expectedParams.Set("limit", "7")
	if result, err := GetFollowing(client, 2, 7); err != nil {
		t.Fatal("Failed to get from following", err)
	} else {
		if result.offset != 2 {
			t.Fatal("Offset should have been 0")
		}
		if result.limit != 7 {
			t.Fatal("Limit should have been 0")
		}
	}
}

func TestFollowingNext(t *testing.T) {
	response := getFollowerString(5, Blog{}, Blog{}, Blog{})
	client := newTestClient(response, nil)
	result, _ := GetFollowing(client, 0, 0)
	// test getting next
	expectedParams := url.Values{}
	expectedParams.Set("limit", "3")
	expectedParams.Set("offset", "3")
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Following.Next",
		http.MethodGet,
		"/user/following",
		expectedParams,
	)

	nextResult, err := result.Next()
	if err != nil {
		t.Fatal("Failed to get next page")
	}
	if nextResult.offset != 3 {
		t.Fatal("Failed to set correct offset")
	}
	if _, err := nextResult.Next(); err != NoNextPageError {
		t.Fatal("Offset exceeding total should mean no next page")
	}
}

func TestFollowingPrevWithoutLimit(t *testing.T) {
	client := newTestClient("{}", nil)
	result, _ := GetFollowing(client, 4, 0)
	result.Blogs = []Blog{{}, {}}
	expectedParams := url.Values{}
	expectedParams.Set("offset", "2")
	expectedParams.Set("limit", "2")
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Following.Next",
		http.MethodGet,
		"/user/following",
		expectedParams,
	)
	if result, err := result.Prev(); err != nil {
		t.Fatal("Prev should succeed while next offset is positive")
	} else if result.limit != 2 {
		t.Fatal("Unspecified limit should be set from size of result set")
	}

}

func TestFollowingPrev(t *testing.T) {
	// total of 5 blogs, 3 in current result set
	response := getFollowerString(5, Blog{}, Blog{}, Blog{})
	client := newTestClient(response, nil)
	// offset 2, limit 5
	result, _ := GetFollowing(client, 2, 5)
	// test getting next
	expectedParams := url.Values{}
	expectedParams.Set("offset", "0")
	expectedParams.Set("limit", "5")
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Following.Next",
		http.MethodGet,
		"/user/following",
		expectedParams,
	)

	prevResult, err := result.Prev()
	if err != nil {
		t.Fatal("Failed to get prev page")
	}
	if prevResult.offset != 0 {
		t.Fatal("Failed to set correct offset after previous page")
	}
	if prevResult.limit != 5 {
		t.Fatal("Failed to set correct limit after previous page")
	}
	_, err = prevResult.Prev()
	if err != NoPrevPageError {
		t.Fatal("Previous page from 0 offset should generate an error")
	}
}

func TestGetFollowersWithError(t *testing.T) {
	mockError := errors.New("Mock error")
	client := newTestClient("{}", mockError)
	if _, err := GetFollowers(client, "david", 0, 0); err != mockError {
		t.Fatal("Failed request does not return correct error")
	}
}

func TestGetFollowersWithJsonParseError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := GetFollowers(client, "david", 0, 0); err == nil {
		t.Fatal("GetFollowers should return an error if failed to parse JSON")
	}

}
func TestGetFollowers(t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	var offset uint = 1
	var limit uint = 2
	params := url.Values{}
	params.Set("offset", "1")
	params.Set("limit", "2")
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"GetFollowers",
		http.MethodGet,
		blogPath("/blog/%s/followers", blogName),
		params,
	)
	if result, err := GetFollowers(client, blogName, offset, limit); err != nil {
		t.Fatal("Failed request does not return correct error")
	} else {
		if result.client != client {
			t.Fatal("GetFollowers result should have a client")
		}
		if result.name != blogName {
			t.Fatal("GetFollowers result should have the blog name")
		}
		if result.offset != offset {
			t.Fatal("GetFollowers result should have correct offset")
		}
		if result.limit != limit {
			t.Fatal("GetFollowers result should have correct limit")
		}
	}
}

func TestFollowerNextFailsOnEmptyResult(t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	result, _ := GetFollowers(client, blogName, 1, 2)
	_, err := result.Next()
	if err != NoNextPageError {
		t.Fatal("Expected no next page error on empty result")
	}
	result.Total = 3
	result.Followers = []Follower{{}}
	if err != NoNextPageError {
		t.Fatal("Expected no next page on offset exceeding total")
	}
	result.Total = 5
	params := url.Values{}
	params.Set("offset", "3")
	params.Set("limit", "2")
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Followers.Next",
		http.MethodGet,
		blogPath("/blog/%s/followers", blogName),
		params,
	)
	_, err = result.Next()
	if err != nil {
		t.Fatal("Next should succeed when offset + limit is less than total")
	}
}

func TestFollowerNextUsesResultSizeIfNoLimit(t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	result, _ := GetFollowers(client, blogName, 0, 0)
	result.Total = 10
	result.Followers = []Follower{{}, {}}
	params := url.Values{}
	params.Set("offset", "2")
	params.Set("limit", "2")
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Followers.Next",
		http.MethodGet,
		blogPath("/blog/%s/followers", blogName),
		params,
	)
	result.Next()
}

func TestFollowerPrev(t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	result, _ := GetFollowers(client, blogName, 3, 2)
	params := url.Values{}
	params.Set("offset", "1")
	params.Set("limit", "2")
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Followers.Next",
		http.MethodGet,
		blogPath("/blog/%s/followers", blogName),
		params,
	)
	nextResult, err := result.Prev()
	if err != nil {
		t.Fatal("Should be able to paginate back while offset > 0")
	}
	params.Set("offset", "0")
	lastResult, err := nextResult.Prev()
	if err != nil {
		t.Fatal("Should be able to paginate back while offset > 0")
	}
	_, err = lastResult.Prev()
	if err != NoPrevPageError {
		t.Fatal("Should get error when attempting previous page on first result set")
	}
}

func TestFollowerPrevUsesResultSizeIfNoLimit(t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	result, _ := GetFollowers(client, blogName, 4, 0)
	result.Total = 10
	result.Followers = []Follower{{}, {}}
	params := url.Values{}
	params.Set("offset", "2")
	params.Set("limit", "2")
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Followers.Prev",
		http.MethodGet,
		blogPath("/blog/%s/followers", blogName),
		params,
	)
	result.Prev()
}

func TestFollow(t *testing.T) {
	blogName := "david"
	params := url.Values{}
	params.Set("url", normalizeBlogName(blogName))
	client := newTestClient("{}", nil)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Follow",
		http.MethodPost,
		"/user/follow",
		params,
	)
	err := Follow(client, blogName)
	if err != nil {
		t.Fatal("Follow fails")
	}
}

func TestUnfollow(t *testing.T) {
	blogName := "david"
	params := url.Values{}
	params.Set("url", normalizeBlogName(blogName))
	client := newTestClient("{}", nil)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Unfollow",
		http.MethodPost,
		"/user/unfollow",
		params,
	)
	err := Unfollow(client, blogName)
	if err != nil {
		t.Fatal("Unfollow fails")
	}
}

func getFollowerString(total uint, blogs ...Blog) string {
	return jsonStringify(map[string]interface{}{
		"response": map[string]interface{}{
			"total_blogs": total,
			"blogs":       blogs,
		},
	})
}
