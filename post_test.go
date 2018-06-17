package tumblr

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestPostRefLike(t *testing.T) {
	client := newTestClient("{}", nil)
	var postId uint64 = 1986
	reblogKey := "reblog-key"
	ref := PostRef{
		client: client,
		MiniPost: MiniPost{
			Id:        postId,
			ReblogKey: reblogKey,
		},
	}
	params := setPostId(postId, url.Values{})
	params.Set("reblog_key", reblogKey)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"PostRef.Like",
		http.MethodPost,
		"/user/like",
		params,
	)
	ref.Like()
}

func TestPostRefUnlike(t *testing.T) {
	client := newTestClient("{}", nil)
	var postId uint64 = 1986
	reblogKey := "reblog-key"
	ref := PostRef{
		client: client,
		MiniPost: MiniPost{
			Id:        postId,
			ReblogKey: reblogKey,
		},
	}
	params := setPostId(postId, url.Values{})
	params.Set("reblog_key", reblogKey)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"PostRef.Like",
		http.MethodPost,
		"/user/unlike",
		params,
	)
	ref.Unlike()
}

func TestMakePostFromType(t *testing.T) {
	testCases := map[string]string{
		"quote":  "QuotePost",
		"chat":   "ChatPost",
		"photo":  "PhotoPost",
		"text":   "TextPost",
		"link":   "LinkPost",
		"answer": "AnswerPost",
		"audio":  "AudioPost",
		"video":  "VideoPost",
	}
	classPrefix := "*tumblr."
	for postType, postClass := range testCases {
		post, err := makePostFromType(postType)
		if err != nil {
			t.Errorf("Unexpected error creating post of type `%s`", postType)
		}
		postClass = classPrefix + postClass
		if actual := reflect.TypeOf(post).String(); actual != postClass {
			t.Errorf("Expected `%s` type to generate struct type `%s`, got `%s` instead", postType, postClass, actual)
		}
	}
	// test default case
	_, err := makePostFromType("")
	if err == nil {
		t.Fatal("Unexpected type should generate an error")
	}
}

func TestStringifyPost(t *testing.T) {
	post := Post{}
	if post.String() != jsonStringify(post) {
		t.Fatal("Post stringify does not conform to expected JSON output")
	}
}

func TestPostDynamicAccessor(t *testing.T) {
	post := Post{}
	post.Id = 1986
	if _, err := post.GetProperty("DoesNotExistProperty"); err == nil {
		t.Fatal("Dynamic accessor should error on property that does not exist")
	}
	actual, err := post.GetProperty("Id")
	if err != nil {
		t.Fatal("Dynamic accessor incorrectly errored")
	}
	if fmt.Sprintf("%v", actual) != "1986" {
		t.Fatalf("Dynamic accessor does not return the proper value; expected %d got %v", post.Id, actual)
	}
}

func TestQueryPostsReturnsClientError(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("", clientErr)
	if _, err := queryPosts(client, "", "", url.Values{}); err == nil {
		t.Fatal("Client error should be returned")
	}
}

func TestQueryPostsReturnsJsonError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := queryPosts(client, "", "", url.Values{}); err == nil {
		t.Fatal("JSON Unmarshal error should be returned")
	}
}

func TestQueryPostsSuccess(t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	path := "/blog/%s/something"
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"queryPosts",
		http.MethodGet,
		blogPath(path, blogName),
		params,
	)
	response, err := queryPosts(
		client,
		path,
		blogName,
		params,
	)
	if err != nil {
		t.Fatal("Posts should have been returned")
	}
	if string(response.response.body) != string(client.response.body) {
		t.Fatal("Response should match client's response")
	}
}

func TestGetPosts(t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"",
		http.MethodGet,
		blogPath("/blog/%s/posts", blogName),
		params,
	)
	if _, err := GetPosts(client, blogName, params); err != nil {
		t.Fatal("Posts should have been returned")
	}
}

func TestGetQueue(t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"",
		http.MethodGet,
		blogPath("/blog/%s/posts/queue", blogName),
		params,
	)
	if _, err := GetQueue(client, blogName, params); err != nil {
		t.Fatal("Posts should have been returned")
	}
}

func TestGetDrafts(t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"",
		http.MethodGet,
		blogPath("/blog/%s/posts/draft", blogName),
		params,
	)
	if _, err := GetDrafts(client, blogName, params); err != nil {
		t.Fatal("Posts should have been returned")
	}
}

func TestGetSubmissions(t *testing.T) {
	client := newTestClient("{}", nil)
	blogName := "david"
	params := url.Values{}
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"",
		http.MethodGet,
		blogPath("/blog/%s/posts/submission", blogName),
		params,
	)
	if _, err := GetSubmissions(client, blogName, params); err != nil {
		t.Fatal("Posts should have been returned")
	}
}

func TestDoPostMissingBlogError(t *testing.T) {
	client := newTestClient("{}", nil)
	if _, err := doPost(client, "", "", url.Values{}); err == nil {
		t.Fatal("Blog name should be required")
	}
}

func TestDoPostClientError(t *testing.T) {
	clientErr := errors.New("Client error")
	client := newTestClient("{}", clientErr)
	if _, err := doPost(client, "", "blog", url.Values{}); err != clientErr {
		t.Fatal("Client error should be returned")
	}
}

func TestDoPostJsonError(t *testing.T) {
	client := newTestClient("{", nil)
	if _, err := doPost(client, "", "blog", url.Values{}); err == nil {
		t.Fatal("Json error should be returned")
	}
}

func TestDoPostSuccess(t *testing.T) {
	var postId uint64 = 1986
	client := newTestClient(fmt.Sprintf("{\"response\": {\"id\": %d}}", postId), nil)
	params := url.Values{}
	path := "/blog/%s/test"
	blog := "david"
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"doPost",
		http.MethodPost,
		blogPath(path, blog),
		params,
	)
	if result, err := doPost(client, path, blog, params); err != nil {
		t.Fatal("Do post should succeed")
	} else {
		if result.Id != postId {
			t.Fatal("Incorrectly parsed the post id")
		} else if result.BlogName != blog {
			t.Fatal("Incorrectly assigned blog name")
		}
	}
}

func TestNewPostRefById(t *testing.T) {
	testClient := newTestClient("", nil)
	var postId uint64 = 1986
	ref := NewPostRefById(testClient, postId)
	if ref.client != testClient {
		t.Fatal("Client was not assigned")
	}
	if ref.Id != postId {
		t.Fatal("Post ID was not assigned")
	}
}

func TestNewPostRef(t *testing.T) {
	testClient := newTestClient("", nil)
	mini := MiniPost{}
	ref := NewPostRef(testClient, &mini)
	if ref.client != testClient {
		t.Fatal("Client was not assigned")
	}
	if ref.MiniPost != mini {
		t.Fatal("Minipost was not assigned")
	}
}

func TestPostRefSetClient(t *testing.T) {
	ref := &PostRef{client: newTestClient("", nil)}
	client2 := newTestClient("{}", nil)
	if ref.SetClient(client2); ref.client != client2 {
		t.Fatal("Client setter failed")
	}
}
func TestCreatePost(t *testing.T) {
	client := newTestClient("{}", nil)
	params := url.Values{}
	blog := "david"
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"CreatePost",
		http.MethodPost,
		blogPath("/blog/%s/post", blog),
		params,
	)
	CreatePost(client, blog, params)
}

func TestEditPost(t *testing.T) {
	client := newTestClient("{}", nil)
	blog := "david"
	var postId uint64 = 1986
	params := setPostId(postId, url.Values{})
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"EditPost",
		http.MethodPost,
		blogPath("/blog/%s/post/edit", blog),
		params,
	)
	EditPost(client, blog, postId, params)
}

func TestPostRef_Edit(t *testing.T) {
	client := newTestClient("{}", nil)
	blog := "david"
	var postId uint64 = 1986
	params := setPostId(postId, url.Values{})
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"EditPost",
		http.MethodPost,
		blogPath("/blog/%s/post/edit", blog),
		params,
	)
	ref := NewPostRefById(client, postId)
	ref.BlogName = blog
	ref.Edit(params)
}

func TestDeletePost(t *testing.T) {
	client := newTestClient("{}", nil)
	blog := "david"
	var postId uint64 = 1986
	params := setPostId(postId, url.Values{})
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"DeletePost",
		http.MethodPost,
		blogPath("/blog/%s/post/delete", blog),
		params,
	)
	DeletePost(client, blog, postId)
}

func TestPostRef_Delete(t *testing.T) {
	client := newTestClient("{}", nil)
	blog := "david"
	var postId uint64 = 1986
	params := setPostId(postId, url.Values{})
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"DeletePost",
		http.MethodPost,
		blogPath("/blog/%s/post/delete", blog),
		params,
	)
	ref := NewPostRefById(client, postId)
	ref.BlogName = blog
	ref.Delete()
}

func TestReblogPostWithoutKey(t *testing.T) {
	client := newTestClient("{}", nil)
	if _, err := ReblogPost(client, "blog", 1986, "", url.Values{}); err == nil {
		t.Fatal("Reblogging with empty key value should generate error")
	}
}

func TestReblogPost(t *testing.T) {
	client := newTestClient("{}", nil)
	blog := "david"
	reblogKey := "reblog-key"
	var postId uint64 = 1986
	params := setPostId(postId, url.Values{})
	params.Set("reblog_key", reblogKey)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"ReblogPost",
		http.MethodPost,
		blogPath("/blog/%s/post/reblog", blog),
		params,
	)
	ReblogPost(client, blog, postId, reblogKey, params)
}

func TestPostRef_ReblogOnBlog(t *testing.T) {
	client := newTestClient("{}", nil)
	blog := "david"
	reblogKey := "reblog-key"
	var postId uint64 = 1986
	params := setPostId(postId, url.Values{})
	params.Set("reblog_key", reblogKey)
	client.confirmExpectedSet = expectClientCallParams(
		t,
		"ReblogPost",
		http.MethodPost,
		blogPath("/blog/%s/post/reblog", blog),
		params,
	)
	ref := NewPostRef(client, &MiniPost{Id: postId, ReblogKey: reblogKey})
	ref.ReblogOnBlog(blog, params)

}

func TestPosts_All(t *testing.T) {
	client := newTestClient("{}", nil)
	posts, err := GetPosts(client, "blog", url.Values{})
	if err != nil {
		t.Fatal("Failed to get posts")
	}
	posts.Posts = []MiniPost{{Type: "quote"}}
	if posts.parsedPosts != nil {
		t.Fatal("Posts initialized with non-nil parsed posts")
	}
	all, err := posts.All()
	if err != nil {
		t.Fatal("Failed to parse Posts")
	}
	if posts.parsedPosts == nil {
		t.Fatal("Posts does not cache parsed posts after All()")
	}
	if len(all) != 1 {
		t.Fatal("Failed to correctly posts from mini posts array")
	}
}

func TestPosts_AllWithJsonError(t *testing.T) {
	client := newTestClient("{}", nil)
	posts, err := GetPosts(client, "blog", url.Values{})
	if err != nil {
		t.Fatal("Failed to get posts")
	}
	posts.response.body = []byte("{")
	_, err = posts.All()
	if err == nil {
		t.Fatal("Failed to return JSON parse error")
	}
}

func TestPosts_Get(t *testing.T) {
	client := newTestClient("{}", nil)
	posts, err := GetPosts(client, "blog", url.Values{})
	if err != nil {
		t.Fatal("Failed to get posts")
	}
	mini := MiniPost{Type: "quote"}
	posts.Posts = []MiniPost{mini}
	mockResponse := struct {
		Response struct {
			Posts []QuotePost `json:"posts"`
		} `json:"response"`
	}{}
	mockResponse.Response.Posts = []QuotePost{{Post: Post{PostRef: PostRef{MiniPost: mini}}}}
	posts.response.body = []byte(jsonStringify(mockResponse))
	if posts.parsedPosts != nil {
		t.Fatal("Posts initialized with non-nil parsed posts")
	}
	post := posts.Get(0)
	if posts.parsedPosts == nil {
		t.Fatal("Get() should cache full array of parsed posts")
	}
	if post.GetSelf().Type != mini.Type {
		t.Fatalf("Get() should return correct type of mini post; expected %s, got %s", post.GetSelf().Type, mini.Type)
	}
	if post := posts.Get(10); post != nil {
		t.Fatal("Getting out of bounds should generate error")
	}
}

func TestPosts_GetWithAllError(t *testing.T) {
	client := newTestClient("{}", nil)
	posts, err := GetPosts(client, "blog", url.Values{})
	if err != nil {
		t.Fatal("Failed to get posts")
	}
	posts.response.body = []byte("{")
	if post := posts.Get(0); post != nil {
		t.Fatal("Get() should return nil on error from All()")
	}
}
