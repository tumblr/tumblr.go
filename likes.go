package tumblr

import (
	"net/url"
	"encoding/json"
)

type Likes struct {
	client ClientInterface
	response *Response
	parsedPosts []PostInterface
	Posts []MiniPost `json:"liked_posts"`
	TotalLikes uint64 `json:"liked_count"`
}

// Retrieves a Users's list of Posts they have liked
// URL values can include:
// 	limit (int)
//	offset (int)
//	before (timestamp)
//	after (timestamp)
func GetLikes(client ClientInterface, params url.Values) (*Likes, error) {
	response, err := client.GetWithParams("/user/likes", params)
	if err != nil {
		return nil, err
	}

	result := struct {
		Response Likes `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &result); err != nil {
		return nil, err
	}
	result.Response.client = client
	result.Response.response = &response
	return &result.Response, nil
}

// Convenience method for performing a like/unlike operation
func doLike(client ClientInterface, path string, postId uint64, reblogKey string) error {
	params := url.Values{}
	params.Set("reblog_key", reblogKey)
	_, err := client.PostWithParams(path, setPostId(postId, params))
	return err
}

// Like a post on behalf of a user
func LikePost(client ClientInterface, postId uint64, reblogKey string) error {
	return doLike(client, "/user/like", postId, reblogKey)
}

// Unlike a post on behalf of a user
func UnlikePost(client ClientInterface, postId uint64, reblogKey string) error {
	return doLike(client, "/user/unlike", postId, reblogKey)
}

// Return an array of full post objects (instead of the default array of MiniPosts initially created)
func (l *Likes)Full() ([]PostInterface, error) {
	var err error = nil
	if l.parsedPosts == nil {
		r := struct {
			Response struct{
					 Posts []PostInterface `json:"liked_posts"`
				 } `json:"response"`
		}{}
		r.Response.Posts = makePostsFromMinis(l.Posts, l.client)
		if err = json.Unmarshal(l.response.body, &r); err != nil {
			l.parsedPosts = []PostInterface{}
		} else {
			l.parsedPosts = r.Response.Posts
		}
	}
	return l.parsedPosts, err
}
