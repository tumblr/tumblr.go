package tumblr

import (
	"testing"
	"errors"
	"net/http"
	"net/url"
	"fmt"
	"strconv"
)

func TestBlog_String(t *testing.T) {
	blog := Blog{}
	if blog.String() != jsonStringify(blog) {
		t.Fatal("Blog JSON representation is incorrect")
	}
}

func TestGetBlogInfoClientError(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("{}", clientErr)
	if _, err := GetBlogInfo(client, "name"); err != clientErr {
		t.Fatal("Blog info must return client error")
	}
}
func TestGetBlogInfoJsonError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := GetBlogInfo(client, "name"); err == nil {
		t.Fatal("Blog info must return json unmarshal error")
	}
}
func TestGetBlogInfo(t *testing.T) {
	client := newTestClient("{}", nil)
	blog := "david"
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"BlogInfo",
		http.MethodGet,
		blogPath("/blog/%s/info", blog),
		url.Values{},
	)
	if result, err := GetBlogInfo(client, blog); err != nil {
		t.Fatal("Blog info failed")
	} else if result.client != client {
		t.Fatal("Result must have client set.")
	}
}

func TestGetAvatarClientError(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("{}", clientErr)
	if _, err := GetAvatar(client, "blog"); err != clientErr {
		t.Fatal("Avatar must return client error")
	}
}

func TestGetAvatarUsesLocationHeader(t *testing.T) {
	client := newTestClient("{", nil)
	uri := "http://placekitten.com"
	client.response.Headers = http.Header{}
	client.response.Headers.Set("Location", uri)
	if actual, err := GetAvatar(client, "blog"); err != nil {
		t.Fatal("Get Avatar should use location header before attempting decode")
	} else if actual != uri {
		t.Fatal("Get AVatar did not retrieve location URI")
	}
}

func TestGetAvatarJsonError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := GetAvatar(client, "blog"); err == nil {
		t.Fatal("Avatar must return json unmarshal error")
	}
}

func TestGetAvatarUsesJSONWithoutHeader(t *testing.T) {
	uri := "http://placekitten.com"
	client := newTestClient(fmt.Sprintf("{\"response\": {\"location\": \"%s\"}}", uri), nil)
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"GetAvatar",
		http.MethodGet,
		blogPath("/blog/%s/avatar", "david"),
		params,
	)
	if actual, err := GetAvatar(client, "david"); err != nil {
		t.Fatal("Get Avatar should use location header before attempting decode")
	} else if actual != uri {
		t.Fatal("Get Avatar did not retrieve location URI")
	}
}

func TestGetAvatarFailsOnBadResponse(t *testing.T) {
	client := newTestClient("{}", nil)
	if _, err := GetAvatar(client, "blog"); err == nil {
		t.Fatal("Get Avatar should fail on response without Location header or location JSON")
	}
}

func TestNewBlogRef(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	if ref := NewBlogRef(client, blog); ref.Name != blog || ref.client == nil {
		t.Fatal("NewBlogRef does not initialize correctly")
	}

}

func TestBlogRef_CreatePost(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"CreatePost",
		http.MethodPost,
		blogPath("/blog/%s/post", blog),
		params,
	)
	ref := NewBlogRef(client, blog)
	if _,err := ref.CreatePost(url.Values{}); err != nil {
		t.Fatal("BlogRef_")
	}
}

func TestBlogRef_Follow(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	params := url.Values{}
	params.Set("url", normalizeBlogName(blog))
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Follow",
		http.MethodPost,
		"/user/follow",
		params,
	)
	ref := NewBlogRef(client, blog)
	if err := ref.Follow(); err != nil {
		t.Fatal("BlogRef_Follow")
	}
}

func TestBlogRef_GetAvatar(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	client.response.Headers = http.Header{}
	uri := "http://placekitten.com"
	client.response.Headers.Set("Location", uri)
	ref := NewBlogRef(client, blog)
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"BlogRef_GetAvatar",
		http.MethodGet,
		blogPath("/blog/%s/avatar", blog),
		params,
	)
	if actual,err := ref.GetAvatar(); err != nil || actual != uri {
		t.Fatal("Ref failed to retrieve avatar")
	}
}

func TestBlogRef_GetDrafts(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"GetDrafts",
		http.MethodGet,
		blogPath("/blog/%s/posts/draft", blog),
		params,
	)
	ref := NewBlogRef(client, blog)
	if _,err := ref.GetDrafts(url.Values{}); err != nil {
		t.Fatal("BlogRef_GetDrafts failed")
	}
}

func TestBlogRef_GetFollowers(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	params := url.Values{}
	params.Set("offset", "0")
	params.Set("limit", "0")
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"GetFollowers",
		http.MethodGet,
		blogPath("/blog/%s/followers", blog),
		params,
	)
	ref := NewBlogRef(client, blog)
	if _,err := ref.GetFollowers(); err != nil {
		t.Fatal("BlogRef_GetFollowers failed")
	}

}

func TestBlogRef_GetInfo(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"GetInfo",
		http.MethodGet,
		blogPath("/blog/%s/info", blog),
		params,
	)
	ref := NewBlogRef(client, blog)
	if _,err := ref.GetInfo(); err != nil {
		t.Fatal("BlogRef_GetInfo")
	}
}

func TestBlogRef_GetPosts(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"GetPosts",
		http.MethodGet,
		blogPath("/blog/%s/posts", blog),
		params,
	)
	ref := NewBlogRef(client, blog)
	if _,err := ref.GetPosts(url.Values{}); err != nil {
		t.Fatal("BlogRef_GetPosts")
	}
}

func TestBlogRef_GetQueue(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"GetQueue",
		http.MethodGet,
		blogPath("/blog/%s/posts/queue", blog),
		params,
	)
	ref := NewBlogRef(client, blog)
	if _,err := ref.GetQueue(url.Values{}); err != nil {
		t.Fatal("BlogRef_GetQueue")
	}
}

func TestBlogRef_ReblogPost(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	params := url.Values{}
	reblogKey := "reblog-key"
	var postId uint64 = 1986
	params.Set("reblog_key", reblogKey)
	params.Set("id", strconv.FormatUint(postId, 10))
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"ReblogPost",
		http.MethodPost,
		blogPath("/blog/%s/post/reblog", blog),
		params,
	)
	ref := NewBlogRef(client, blog)
	post := NewPostRef(client, &MiniPost{})
	post.ReblogKey = reblogKey
	post.Id = postId
	if _,err := ref.ReblogPost(post, url.Values{}); err != nil {
		t.Fatal("BlogRef_ReblogPost")
	}
}

func TestBlogRef_Unfollow(t *testing.T) {
	blog := "david"
	client := newTestClient("{}", nil)
	params := url.Values{}
	params.Set("url", normalizeBlogName(blog))
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"Unfollow",
		http.MethodPost,
		"/user/unfollow",
		params,
	)
	ref := NewBlogRef(client, blog)
	if err := ref.Unfollow(); err != nil {
		t.Fatal("BlogRef_Unfollow")
	}
}
