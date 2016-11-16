package tumblr

import (
	"net/url"
	"encoding/json"
)

type FollowingList struct {
	client ClientInterface
	Total  uint32 `json:"total_blogs"`
	Blogs  []Blog `json:"blogs"`
	offset uint
	limit  uint
}

// Object from the lsit of followers response
type FollowerList struct {
	client    ClientInterface
	Total     uint32 `json:"total_users"`
	Followers []Follower `json:"users"`
	name      string
	offset    uint
	limit     uint
}

// FollowerList substructure
type Follower struct {
	Following bool `json:"following"`
	Name string `json:"name"`
	Updated int64 `json:"updated"`
	Url string `json:"url"`
}

// Retrieves the list of blogs this user follows
func GetFollowing(client ClientInterface, offset, limit uint) (*FollowingList, error) {
	params := setParamsUint(uint64(limit), url.Values{}, "limit")
	params = setParamsUint(uint64(offset), params, "offset")
	result, err := client.GetWithParams("/user/following", params)
	if err != nil {
		return nil, err
	}
	response := struct{
		Response FollowingList `json:"response"`
	}{
		Response: FollowingList{
			client: client,
			limit: limit,
			offset: offset,
		},
	}
	if err = json.Unmarshal(result.body, &response); err != nil {
		return nil, err
	}
	return &response.Response, nil
}

// Retrieves the next page of followers
func (f *FollowingList)Next() (*FollowingList, error) {
	limit := f.limit
	if limit < 1 {
		limit = uint(len(f.Blogs))
	}
	offset := f.offset + limit
	if offset >= uint(f.Total) {
		return nil, NoNextPageError
	}
	return GetFollowing(f.client, offset, limit)
}

// Retrieves the previous page of followers
func (f *FollowingList)Prev() (*FollowingList, error) {
	if f.offset <= 0 {
		return nil, NoPrevPageError
	}
	limit := f.limit
	if limit < 1 {
		limit = uint(len(f.Blogs))
	}
	var newOffset uint = f.offset - limit
	if limit >= f.offset {
		newOffset = 0
	}
	return GetFollowing(f.client, newOffset, limit)
}

// Retrieve User's followers
func GetFollowers(client ClientInterface, name string, offset, limit uint) (*FollowerList, error) {
	params := setParamsUint(uint64(offset), url.Values{}, "offset")
	params = setParamsUint(uint64(limit), params, "limit")
	response, err := client.GetWithParams(blogPath("/blog/%s/followers", name), params)
	if err != nil {
		return nil, err
	}
	followers := struct {
		Followers FollowerList `json:"response"`
	}{
		Followers: FollowerList{
			client: client,
			name: name,
			limit: limit,
			offset: offset,
		},
	}
	if err = json.Unmarshal(response.body, &followers); err == nil {
		return &followers.Followers, nil
	}
	return nil, err
}

// Get next page of a user's followers
func (f *FollowerList)Next() (*FollowerList, error){
	limit := f.limit
	if limit < 1 {
		limit = uint(len(f.Followers))
	}
	offset := f.offset + limit
	if uint32(offset) >= f.Total || len(f.Followers) < 1 {
		return nil, NoNextPageError
	}
	return GetFollowers(f.client, f.name, offset, limit)
}

// Get previous page of a user's followers
func (f *FollowerList)Prev() (*FollowerList, error){
	if f.offset <= 0 {
		return nil, NoPrevPageError
	}
	limit := f.limit
	if limit < 1 {
		limit = uint(len(f.Followers))
	}
	offset := f.offset - limit
	if limit >= f.offset {
		offset = 0
	}
	return GetFollowers(f.client, f.name, offset, limit)
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
