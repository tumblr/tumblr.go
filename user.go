package tumblr

import (
	"encoding/json"
)

type User struct {
	Following uint32 `json:"following"`
	DefaultPostFormat string `json:"default_post_format"`
	Name string `json:"name"`
	Likes uint64 `json:"likes"`
	Blogs []ShortBlog `json:"blogs"`
}

// Retrieves the current user's info (based on the client's token/secret values)
func GetUserInfo(client ClientInterface) (*User, error) {
	response, err := client.Get("/user/info")
	if err != nil {
		return nil, err
	}
	result := struct{
		Response struct {
			User User `json:"user"`
		} `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &result); err != nil {
		return nil, err
	}
	return &result.Response.User, nil
}
