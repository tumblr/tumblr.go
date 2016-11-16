package tumblr

import (
	"net/url"
	"encoding/json"
	"strconv"
)

type SearchResults struct {
	client ClientInterface
	Posts []PostInterface `json:"response"`
	params url.Values
}

// gets page of posts
func TaggedSearch(client ClientInterface, tag string, params url.Values) (*SearchResults, error) {
	params.Set("tag", tag)
	response, err := client.GetWithParams("/tagged", params)
	if err != nil {
		return nil, err
	}
	result := struct {
		Response []MiniPost `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &result); err != nil {
		return nil, err
	}
	minis := result.Response
	full := SearchResults{
		Posts: makePostsFromMinis(minis, client),
		client: client,
		params: params,
	}
	if err = json.Unmarshal(response.body, &full); err != nil {
		return nil, err
	}
	return &full, nil
}

// returns next page of results
func (s *SearchResults) Next() (*SearchResults, error) {
	// get last timestamp
	var size = len(s.Posts)
	if size < 1 {
		return nil, NoNextPageError
	}
	lastPost := s.Posts[size - 1].GetSelf()
	lastTs := lastPost.FeaturedTimestamp
	if lastTs < 1 {
		lastTs = lastPost.Timestamp
	}
	params := s.params
	params.Set("before", strconv.FormatUint(lastTs, 10))
	return TaggedSearch(s.client, params.Get("tag"), params)
}