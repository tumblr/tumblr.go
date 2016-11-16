package tumblr

import (
	"strings"
	"errors"
	"encoding/json"
	"fmt"
	"net/url"
)

// Reference to a blog which can be used to perform further blog actions
type BlogRef struct {
	client ClientInterface
	Name string `json:"name"`
}

// Subset of blog information, returned in the list of blogs belonging to a user (see: GetUserInfo)
type ShortBlog struct {
	BlogRef
	Url string `json:"url"`
	Title string `json:"title"`
	IsPrimary bool `json:"primary"`
	FollowerCount uint32 `json:"followers"`
	PostToTwitter string `json:"tweet"`
	PostToFacebook string `json:"facebook"`
	Visibility string `json:"type"`
}

// Tumblelog struct
type Blog struct {
	BlogRef
	Url string `json:"url"`
	Title string `json:"title"`
	Posts int64 `json:"posts"`
	Ask bool `json:"ask"`
	AskAnon bool `json:"ask_anon"`
	AskAnonPageTitle string `json:"ask_page_title"`
	CanSendFanMail bool `json:"can_send_fan_mail"`
	CanSubmit bool `json:"can_submit"`
	CanSubscribe bool `json:"can_subscribe"`
	Description string `json:"description"`
	Followed bool `json:"followed"`
	IsBlockedFromPrimary bool `json:"is_blocked_from_primary"`
	IsNSFW bool `json:"is_nsfw"`
	ShareLikes bool `json:"share_likes"`
	SubmissionPageTitle string `json:"submission_page_title"`
	Subscribed bool `json:"subscribed"`
	TotalPosts int64 `json:"total_posts"`
	Updated int64 `json:"updated"`
}

// Convenience method converting a Blog into a JSON representation
func (b *Blog) String() string {
	return jsonStringify(*b)
}

// Retrieve information about a blog
func GetBlogInfo(client ClientInterface, name string) (*Blog, error) {
	response, err := client.Get(blogPath("/blog/%s/info", name))
	if err != nil {
		return nil, err
	}
	blog := struct {
		Response struct {
			Blog Blog `json:"blog"`
		} `json:"response"`
	}{}
	//blog := blogResponse{}
	err = json.Unmarshal(response.body, &blog)
	if err != nil {
		return nil, err
	}
	blog.Response.Blog.client = client
	return &blog.Response.Blog, nil
}

// Retrieve Blog's Avatar URI
func GetAvatar(client ClientInterface, name string) (string, error) {
	response, err := client.Get(blogPath("/blog/%s/avatar", name))
	if err != nil {
		return "", err
	}
	if location := response.Headers.Get("Location"); len(location) > 0 {
		return location, nil
	}
	if err = response.PopulateFromBody(); err != nil {
		return "", err
	}
	if l, ok := response.Result["location"]; ok {
		if location, ok := l.(string); ok {
			return location, nil
		}
	}
	return "", errors.New("Unable to detect avatar location")
}

// Create a BlogRef
func NewBlogRef(client ClientInterface, name string) (*BlogRef) {
	return &BlogRef{
		Name: name,
		client: client,
	}
}

// Retrieves blog info for the given blog reference
func (b *BlogRef) GetInfo() (*Blog, error) {
	return GetBlogInfo(b.client, b.Name)
}

// Retrieves blog avatar for the given blog reference
func (b *BlogRef) GetAvatar() (string, error) {
	return GetAvatar(b.client, b.Name)
}

// Retrieves blog's followers for the given blog reference
func (b *BlogRef) GetFollowers() (*FollowerList, error) {
	return GetFollowers(b.client, b.Name, 0, 0)
}

// Retrieves blog's posts for the given blog reference
func (b *BlogRef) GetPosts(params url.Values) (*Posts, error) {
	return GetPosts(b.client, b.Name, params)
}

// Retrieves blog's queue for the given blog reference
func (b *BlogRef) GetQueue(params url.Values) (*Posts, error) {
	return GetQueue(b.client, b.Name, params)
}

// Retrieves blog's drafts for the given blog reference
func (b *BlogRef) GetDrafts(params url.Values) (*Posts, error) {
	return GetDrafts(b.client, b.Name, params)
}

// Creates a post on the blog represented by BlogRef
func (b *BlogRef) CreatePost(params url.Values) (*PostRef, error) {
	return CreatePost(b.client, b.Name, params)
}

// Reblogs a post to the blog represented by BlogRef
func (b *BlogRef) ReblogPost(p *PostRef, params url.Values) (*PostRef, error) {
	return p.ReblogOnBlog(b.Name, params)
}

// Retrieves name property
func (b *BlogRef) getName() string {
	return b.Name
}

// Follows this blog for the current user (based on OAuth user token/secret)
func (b *BlogRef) Follow() error {
	return Follow(b.client, b.getName())
}

// Unfollows this blog for the current user (based on OAuth user token/secret)
func (b *BlogRef) Unfollow() error {
	return Unfollow(b.client, b.getName())
}

// Helper function to allow for less verbose code
func normalizeBlogName(name string) string {
	if !strings.Contains(name, ".") {
		name = fmt.Sprintf("%s.tumblr.com", name)
	}
	return name
}

// Expects path to contain a single %s placeholder to be substituted with the result of normalizeBlogName
func blogPath(path, name string) string {
	return fmt.Sprintf(path, normalizeBlogName(name))
}