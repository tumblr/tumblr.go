package tumblr

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type FollowingList struct {
	client ClientInterface
	Total  int    `json:"total_blogs"`
	Blogs  []Blog `json:"blogs"`
	offset int
	limit  int
}

// Object from the lsit of followers response
type FollowerList struct {
	client    ClientInterface
	Total     int        `json:"total_users"`
	Followers []Follower `json:"users"`
	name      string
	offset    int
	limit     int
}

// FollowerList substructure
type Follower struct {
	Following bool   `json:"following"`
	Name      string `json:"name"`
	Updated   int    `json:"updated"`
	Url       string `json:"url"`
}

// Retrieves the list of blogs this user follows
func GetFollowing(client ClientInterface, params url.Values) (*FollowingList, error) {
	result, err := client.GetWithParams("/user/following", params)
	if err != nil {
		return nil, err
	}
	limit, _ := strconv.Atoi(params.Get("limit"))
	offset, _ := strconv.Atoi(params.Get("offset"))
	response := struct {
		Response FollowingList `json:"response"`
	}{
		Response: FollowingList{
			client: client,
			limit:  limit,
			offset: offset,
		},
	}
	if err = json.Unmarshal(result.body, &response); err != nil {
		return nil, err
	}
	return &response.Response, nil
}

// GetFollowingOfBlog comment
func GetFollowingOfBlog(client ClientInterface, name string, params url.Values) (*FollowingList, error) {
	result, err := client.GetWithParams(blogPath("/blog/%s/following", name), params)
	if err != nil {
		return nil, err
	}
	limit, _ := strconv.Atoi(params.Get("limit"))
	offset, _ := strconv.Atoi(params.Get("offset"))
	response := struct {
		Response FollowingList `json:"response"`
	}{
		Response: FollowingList{
			client: client,
			limit:  limit,
			offset: offset,
		},
	}
	if err = json.Unmarshal(result.body, &response); err != nil {
		return nil, err
	}
	return &response.Response, nil
}

// Retrieves the next page of followers
func (f *FollowingList) Next() (*FollowingList, error) {
	limit := f.limit
	if limit < 1 {
		limit = int(len(f.Blogs))
	}
	offset := f.offset + limit
	if offset >= f.Total {
		return nil, NoNextPageError
	}
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("offset", fmt.Sprintf("%d", offset))
	return GetFollowing(f.client, params)
}

// Retrieves the previous page of followers
func (f *FollowingList) Prev() (*FollowingList, error) {
	if f.offset <= 0 {
		return nil, NoPrevPageError
	}
	limit := f.limit
	if limit < 1 {
		limit = int(len(f.Blogs))
	}
	var newOffset = f.offset - limit
	if limit >= f.offset {
		newOffset = 0
	}
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("offset", fmt.Sprintf("%d", newOffset))
	return GetFollowing(f.client, params)
}

// Retrieve User's followers
func GetFollowers(client ClientInterface, name string, params url.Values) (*FollowerList, error) {
	response, err := client.GetWithParams(blogPath("/blog/%s/followers", name), params)
	if err != nil {
		return nil, err
	}
	limit, _ := strconv.Atoi(params.Get("limit"))
	offset, _ := strconv.Atoi(params.Get("offset"))
	followers := struct {
		Followers FollowerList `json:"response"`
	}{
		Followers: FollowerList{
			client: client,
			name:   name,
			limit:  limit,
			offset: offset,
		},
	}
	if err = json.Unmarshal(response.body, &followers); err == nil {
		return &followers.Followers, nil
	}
	return nil, err
}

// Get next page of a user's followers
func (f *FollowerList) Next() (*FollowerList, error) {
	limit := f.limit
	if limit < 1 {
		limit = int(len(f.Followers))
	}
	offset := f.offset + limit
	if int(offset) >= f.Total || len(f.Followers) < 1 {
		return nil, NoNextPageError
	}
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("offset", fmt.Sprintf("%d", offset))
	return GetFollowers(f.client, f.name, params)
}

// Get previous page of a user's followers
func (f *FollowerList) Prev() (*FollowerList, error) {
	if f.offset <= 0 {
		return nil, NoPrevPageError
	}
	limit := f.limit
	if limit < 1 {
		limit = int(len(f.Followers))
	}
	offset := f.offset - limit
	if limit >= f.offset {
		offset = 0
	}
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("offset", fmt.Sprintf("%d", offset))
	return GetFollowers(f.client, f.name, params)
}

// Follow a blog
func Follow(client ClientInterface, blogName string) error {
	_, err := client.PostWithParams("/user/follow", url.Values{
		"url": []string{normalizeBlogName(blogName)},
	})
	return err
}

// Unfollow a blog
func Unfollow(client ClientInterface, blogName string) error {
	_, err := client.PostWithParams("/user/unfollow", url.Values{
		"url": []string{normalizeBlogName(blogName)},
	})
	return err
}
