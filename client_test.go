package tumblr

import (
	"net/http"
	"net/url"
	"testing"
)

var testJsonStringifyCases = []jsonStringifyTestCase{
	{
		output: "{\n  \"id\": 0,\n  \"type\": \"\",\n  \"blog_name\": \"\",\n  \"reblog_key\": \"\",\n  \"body\": \"\",\n  \"can_like\": false,\n  \"can_reblog\": false,\n  \"can_reply\": false,\n  \"can_send_in_message\": false,\n  \"caption\": \"\",\n  \"date\": \"\",\n  \"display_avatar\": false,\n  \"followed\": false,\n  \"format\": \"\",\n  \"highlighted\": null,\n  \"liked\": false,\n  \"note_count\": 0,\n  \"permalink_url\": \"\",\n  \"post_url\": \"\",\n  \"reblog\": {\n    \"comment\": \"\",\n    \"tree_html\": \"\"\n  },\n  \"recommended_color\": \"\",\n  \"recommended_source\": false,\n  \"short_url\": \"\",\n  \"slug\": \"\",\n  \"source_title\": \"\",\n  \"source_url\": \"\",\n  \"state\": \"\",\n  \"summary\": \"\",\n  \"tags\": null,\n  \"timestamp\": 0,\n  \"trail\": null\n}",
		input:  Post{},
	},
	{
		output: "{\n  \"name\": \"\",\n  \"url\": \"\",\n  \"title\": \"\",\n  \"posts\": 0,\n  \"ask\": false,\n  \"ask_anon\": false,\n  \"ask_page_title\": \"\",\n  \"can_send_fan_mail\": false,\n  \"can_submit\": false,\n  \"can_subscribe\": false,\n  \"description\": \"\",\n  \"followed\": false,\n  \"is_blocked_from_primary\": false,\n  \"is_nsfw\": false,\n  \"share_likes\": false,\n  \"submission_page_title\": \"\",\n  \"subscribed\": false,\n  \"total_posts\": 0,\n  \"updated\": 0\n}",
		input:  Blog{},
	},
}

type jsonStringifyTestCase struct {
	input  interface{}
	output string
}

// Basic test of stringified objects
func TestJsonStringify(t *testing.T) {
	for _, testCase := range testJsonStringifyCases {
		if out := jsonStringify(testCase.input); out != testCase.output {
			t.Errorf("Expected %s for json stringify of %v. Got %s", testCase.output, testCase.input, out)
		}
	}
}

func TestCopyParams(t *testing.T) {
	orig := url.Values{}
	copied := copyParams(orig)
	copied.Set("key", "value")
	if orig.Get("key") == copied.Get("key") {
		t.Fatal("Copy params allows for modification of source map")
	}
}

// shortcut for the most common case
func TestSetPostId(t *testing.T) {
	params := setPostId(1986, url.Values{})
	if params.Get("id") != "1986" || len(params) != 1 {
		t.Fatal("Did not correctly set id on params")
	}
}

// convenience function for setting an int
func TestSetParamsUint(t *testing.T) {
	params := setParamsUint(1986, url.Values{}, "key")
	if params.Get("key") != "1986" || len(params) != 1 {
		t.Fatal("Did not correctly set key on params")
	}
}

type testClient struct {
	response           Response
	err                error
	confirmExpectedSet func(method, path string, params url.Values)
}

func newTestClient(response string, err error) *testClient {
	return &testClient{
		response: Response{
			body: []byte(response),
		},
		err: err,
	}
}

func expectClientCallParams(t *testing.T, methodName, expectedMethod, expectedPath string, expectedParams url.Values) func(method, path string, params url.Values) {
	return func(method, path string, params url.Values) {
		if method != expectedMethod {
			t.Fatalf("%s expected a %s request method, attempted %s", methodName, expectedMethod, method)
		}
		if path != expectedPath {
			t.Fatalf("%s expected a request to path `%s`, attempted request to `%s`", methodName, expectedPath, path)
		}
		// check params
		if len(expectedParams) != len(params) {
			t.Fatalf("%s expected %d params, saw %d", methodName, len(expectedParams), len(params))
		}
		// check all params equal
		for key, val := range expectedParams {
			if len(val) < 1 {
				t.Fatalf("%s specified a param key `%s` with no values", methodName, key)
			}
			if v := params.Get(key); v != val[0] {
				t.Fatalf("%s expected param %s => %s, saw value %s instead", methodName, key, val[0], v)
			}
		}
	}
}

func (c *testClient) checkCallParams(method, path string, params url.Values) {
	if c.confirmExpectedSet != nil {
		c.confirmExpectedSet(method, path, params)
	}
}

func (c *testClient) Get(endpoint string) (Response, error) {
	c.checkCallParams(http.MethodGet, endpoint, url.Values{})
	return c.response, c.err
}

func (c *testClient) GetWithParams(endpoint string, params url.Values) (Response, error) {
	c.checkCallParams(http.MethodGet, endpoint, params)
	return c.response, c.err
}

func (c *testClient) Post(endpoint string) (Response, error) {
	c.checkCallParams(http.MethodPost, endpoint, url.Values{})
	return c.response, c.err
}

func (c *testClient) PostWithParams(endpoint string, params url.Values) (Response, error) {
	c.checkCallParams(http.MethodPost, endpoint, params)
	return c.response, c.err
}

func (c *testClient) Put(endpoint string) (Response, error) {
	c.checkCallParams(http.MethodPut, endpoint, url.Values{})
	return c.response, c.err
}

func (c *testClient) PutWithParams(endpoint string, params url.Values) (Response, error) {
	c.checkCallParams(http.MethodPut, endpoint, params)
	return c.response, c.err
}

func (c *testClient) Delete(endpoint string) (Response, error) {
	c.checkCallParams(http.MethodDelete, endpoint, url.Values{})
	return c.response, c.err
}

func (c *testClient) DeleteWithParams(endpoint string, params url.Values) (Response, error) {
	c.checkCallParams(http.MethodDelete, endpoint, params)
	return c.response, c.err
}
