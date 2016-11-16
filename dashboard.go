package tumblr

import (
	"net/url"
	"encoding/json"
	"strconv"
	"errors"
)

type Dashboard struct {
	client ClientInterface
	params url.Values
	bySince bool
	byOffset bool
	Posts []PostInterface `json:"posts"`
}

// Retreive a User's dashboard
func GetDashboard(client ClientInterface, params url.Values) (*Dashboard, error) {
	if params.Get("offset") != "" && params.Get("since_id") != "" {
		return nil, errors.New("Cannot specify both offset and since_id")
	}

	response, err := client.GetWithParams("/user/dashboard", params)
	if err != nil {
		return nil, err
	}
	result := struct {
		Response struct {
				 Posts []MiniPost `json:"posts"`
			 } `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &result); err != nil {
		return nil, err
	}
	minis := result.Response.Posts
	full := struct {
		Response Dashboard `json:"response"`
	}{
		Response: Dashboard{
			client: client,
			params: params,
			byOffset: params.Get("offset") != "",
			bySince: params.Get("since_id") != "",
		},
	}
	full.Response.Posts = makePostsFromMinis(minis, client)
	if err = json.Unmarshal(response.body, &full); err != nil {
		return nil, err
	}
	return &full.Response, nil
}

// Error generated when a Dashboard result set it attempting to change pagination methods
var MixedPaginationMethodsError error = errors.New("Cannot mix pagination between SinceId and Offset")

// Returns the next page of a user's dashboard using the current page's last Post id
func (d *Dashboard)NextBySinceId() (*Dashboard, error) {
	if d.byOffset {
		return nil, MixedPaginationMethodsError
	}
	size := len(d.Posts)
	if size < 1 {
		return nil, NoNextPageError
	}
	lastId := d.Posts[size - 1].GetSelf().Id
	params := setParamsUint(lastId, copyParams(d.params), "since_id")
	return GetDashboard(d.client, params)
}

// Returns the next page of a user's dashboard using the current page's offset
func (d *Dashboard)NextByOffset() (*Dashboard, error) {
	if d.bySince {
		return nil, MixedPaginationMethodsError
	}
	if len(d.Posts) < 1 {
		return nil, NoNextPageError
	}
	params := copyParams(d.params)
	offset, err := strconv.Atoi(params.Get("offset"))
	if err != nil {
		offset = 0
	}
	offset += len(d.Posts)
	params.Set("offset", strconv.Itoa(offset))
	return GetDashboard(d.client, params)
}

