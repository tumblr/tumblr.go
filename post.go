package tumblr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
)

// Posts represents a list of MiniPosts, which have a minimal set of information.
type Posts struct {
	client      ClientInterface
	response    Response
	parsedPosts []PostInterface
	Posts       []MiniPost `json:"posts"`
	TotalPosts  int64      `json:"total_posts"`
}

// All will retrieve fully fleshed post data from stubs and cache result.
func (p *Posts) All() ([]PostInterface, error) {
	var err error = nil
	if p.parsedPosts == nil {
		r := struct {
			Response struct {
				Posts []PostInterface `json:"posts"`
			} `json:"response"`
		}{}
		r.Response.Posts = makePostsFromMinis(p.Posts, p.client)
		if err = json.Unmarshal(p.response.body, &r); err != nil {
			p.parsedPosts = []PostInterface{}
		} else {
			p.parsedPosts = r.Response.Posts
		}
	}
	return p.parsedPosts, err
}

// Get retrieves a single Post entity at a given index or returns nil if index is out of bounds.
func (p *Posts) Get(index uint) PostInterface {
	if posts, err := p.All(); err == nil {
		if index >= uint(len(posts)) {
			return nil
		}
		return posts[index]
	}
	return nil
}

// MiniPost stores the basics for what is needed in a Post.
type MiniPost struct {
	Id        uint64 `json:"id"`
	Type      string `json:"type"`
	BlogName  string `json:"blog_name"`
	ReblogKey string `json:"reblog_key"`
}

// PostRef is a base struct used as a starting point for performing operations on a post.
type PostRef struct {
	MiniPost
	client ClientInterface
}

// Post holds the common fields of any post type.
type Post struct {
	PostRef
	Body             string        `json:"body"`
	CanLike          bool          `json:"can_like"`
	CanReblog        bool          `json:"can_reblog"`
	CanReply         bool          `json:"can_reply"`
	CanSendInMessage bool          `json:"can_send_in_message"`
	Caption          string        `json:"caption"`
	Date             string        `json:"date"`
	DisplayAvatar    bool          `json:"display_avatar"`
	Followed         bool          `json:"followed"`
	Format           string        `json:"format"`
	Highlighted      []interface{} `json:"highlighted"`
	Liked            bool          `json:"liked"`
	NoteCount        uint64        `json:"note_count"`
	PermalinkUrl     string        `json:"permalink_url"`
	PostUrl          string        `json:"post_url"`
	Reblog           struct {
		Comment  string `json:"comment"`
		TreeHTML string `json:"tree_html"`
	} `json:"reblog"`
	RecommendedColor  string            `json:"recommended_color"`
	RecommendedSource bool              `json:"recommended_source"`
	ShortUrl          string            `json:"short_url"`
	Slug              string            `json:"slug"`
	SourceTitle       string            `json:"source_title"`
	SourceUrl         string            `json:"source_url"`
	State             string            `json:"state"`
	Summary           string            `json:"summary"`
	Tags              []string          `json:"tags"`
	Timestamp         uint64            `json:"timestamp"`
	FeaturedTimestamp uint64            `json:"featured_timestamp,omitempty"`
	TrackName         string            `json:"track_name,omitempty"`
	Trail             []ReblogTrailItem `json:"trail"`
}

// ReblogTrailItem represents an item in the "trail" to the original, root Post.
type ReblogTrailItem struct {
	Blog          Blog   `json:"blog"`
	Content       string `json:"content"`
	ContentRaw    string `json:"content_raw"`
	IsCurrentItem bool   `json:"is_current_item"`
	Post          struct {
		// sometimes an actual int, sometimes a numeric string, always a headache
		Id interface{} `json:"id"`
	} `json:"post"`
}

// PostInterface is the interface for any concrete Post type to retrieve a property.
type PostInterface interface {
	GetProperty(key string) (interface{}, error)
	GetSelf() *Post
}

// QuotePost represents a Quote Post.
type QuotePost struct {
	Post
	Source string `json:"source,omitempty"`
	Text   string `json:"text"`
}

// ChatPost represents a Chat Post.
type ChatPost struct {
	Post
	Dialog []struct {
		Label  string `json:"label"`
		Name   string `json:"name"`
		Phrase string `json:"phrase"`
	} `json:"dialog"`
}

// TextPost represents a Text Post.
type TextPost struct {
	Post
	Title string `json:"title"`
}

// LinkPost represents a Link Post.
type LinkPost struct {
	Post
	Description string `json:"description"`
	Excerpt     string `json:"excerpt"`
	LinkAuthor  string `json:"link_author"`
	Title       string `json:"title"`
	Url         string `json:"url"`
}

// AnswerPost represents an Answer Post.
type AnswerPost struct {
	Post
	Answer     string `json:"answer"`
	AskingName string `json:"asking_name"`
	AskingUrl  string `json:"asking_url"`
	Publisher  string `json:"publisher"`
	Question   string `json:"question"`
}

// AudioPost represents an Audio Post.
type AudioPost struct {
	Post
	AlbumArt       string `json:"album_art"`
	Artist         string `json:"artist"`
	AudioSourceUrl string `json:"audio_source_url"`
	AudioType      string `json:"audio_type"`
	AudioUrl       string `json:"audio_url"`
	Embed          string `json:"embed"`
	Player         string `json:"player"`
	Plays          uint64 `json:"plays"`
}

// VideoPost represents a Video Post.
type VideoPost struct {
	Post
	Html5Capable bool   `json:"html5_capable"`
	PermalinkUrl string `json:"permalink_url"`
	Players      []struct {
		EmbedCode StringOrBool `json:"embed_code"`
		Width     interface{}  `json:"width"`
	} `json:"player"`
	ThumbnailHeight uint32 `json:"thumbnail_height"`
	ThumbnailUrl    string `json:"thumbnail_url"`
	ThumbnailWidth  uint32 `json:"thumbnail_width"`
	Video           map[string]struct {
		Height  uint32 `json:"height"`
		Width   uint32 `json:"width"`
		VideoId string `json:"video_id"`
	} `json:"video"`
	VideoType string `json:"video_type"`
}

// PhotoPost represents a Photo Post.
type PhotoPost struct {
	Post
	ImagePermalink string  `json:"image_permalink"`
	Photos         []Photo `json:"photos"`
}

// Photo represents one photo in a PhotoPost.
type Photo struct {
	AltSizes     []PhotoSize `json:"alt_sizes"`
	Caption      string      `json:"caption"`
	OriginalSize PhotoSize   `json:"original_size"`
}

// PhotoSize represents a particular size for a Photo.
type PhotoSize struct {
	Height uint32 `json:"height"`
	Width  uint32 `json:"width"`
	Url    string `json:"url"`
}

// StringOrBool is a string that can be unmarshalled from a string or a boolean value.
type StringOrBool string

// UnmarshalJSON implements the json.Unmarshaler interface to ingest strings or boolean values.
func (sb *StringOrBool) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		return json.Unmarshal(b, (*string)(sb))
	}

	var bl bool
	if err := json.Unmarshal(b, &bl); err != nil {
		return err
	}

	if bl {
		*sb = StringOrBool("true")
	} else {
		*sb = StringOrBool("false")
	}
	return nil
}

// String returns the Post as a JSON string.
func (p *Post) String() string {
	return jsonStringify(*p)
}

// GetProperty uses reflection to retrieve one-off field values.
func (p *Post) GetProperty(key string) (interface{}, error) {
	if field, exists := reflect.TypeOf(p).Elem().FieldByName(key); exists {
		return reflect.ValueOf(p).Elem().FieldByIndex(field.Index), nil
	}
	return nil, errors.New(fmt.Sprintf("Property %s does not exist", key))
}

// GetSelf returns the Post from a PostInterface.
func (p *Post) GetSelf() *Post {
	return p
}

// helper method for querying a given path which should return a list of posts
func queryPosts(client ClientInterface, path, name string, params url.Values) (*Posts, error) {
	response, err := client.GetWithParams(blogPath(path, name), params)
	if err != nil {
		return nil, err
	}
	posts := struct {
		Response Posts `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &posts); err == nil {
		posts.Response.response = response
		posts.Response.client = client
		// store
		return &posts.Response, nil
	}
	return nil, err
}

// GetPosts retrieves a blog's posts, in the API docs you can find how to filter by ID, type, etc.
func GetPosts(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts", name, params)
}

// GetQueue retrieves a blog's queued posts.
func GetQueue(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts/queue", name, params)
}

// GetDrafts retrieves a blog's draft posts.
func GetDrafts(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts/draft", name, params)
}

// GetSubmissions retrieves a blog's submission posts.
func GetSubmissions(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts/submission", name, params)
}

// Util method for decoding the response and converting the resulting ID into a PostRef
func doPost(client ClientInterface, path, blogName string, params url.Values) (*PostRef, error) {
	if blogName == "" {
		return nil, errors.New("No blog name provided")
	}
	response, err := client.PostWithParams(blogPath(path, blogName), params)
	if err != nil {
		return nil, err
	}
	post := struct {
		Response struct {
			Id uint64 `json:"id"`
		} `json:"response"`
	}{}
	if err = json.Unmarshal(response.body, &post); err == nil {
		ref := NewPostRefById(client, post.Response.Id)
		ref.BlogName = blogName
		return ref, nil
	}
	return nil, err
}

// NewPostRefById creates a PostRef for the id.
func NewPostRefById(client ClientInterface, id uint64) *PostRef {
	return &PostRef{
		client: client,
		MiniPost: MiniPost{
			Id: id,
		},
	}
}

// NewPostRef creates a PostRef for the MiniPost.
func NewPostRef(client ClientInterface, post *MiniPost) *PostRef {
	return &PostRef{
		client:   client,
		MiniPost: *post,
	}
}

// SetClient sets the client member of the PostRef.
func (r *PostRef) SetClient(c ClientInterface) {
	r.client = c
}

// CreatePost will create a Post on tumblr for the blog in name.
func CreatePost(client ClientInterface, name string, params url.Values) (*PostRef, error) {
	return doPost(client, "/blog/%s/post", name, params)
}

// EditPost will update a Post on tumblr for the blog in name and Post in postId.
func EditPost(client ClientInterface, blogName string, postId uint64, params url.Values) error {
	_, err := client.PostWithParams(blogPath("/blog/%s/post/edit", blogName), setPostId(postId, params))
	return err
}

// Edit will update this Post on tumblr.
func (p *PostRef) Edit(params url.Values) error {
	return EditPost(p.client, p.BlogName, p.Id, params)
}

// ReblogPost will reblog the post in postId and reblogKey to the blog blogName.
func ReblogPost(client ClientInterface, blogName string, postId uint64, reblogKey string, params url.Values) (*PostRef, error) {
	if reblogKey == "" {
		return nil, errors.New("No reblog key provided")
	}
	params.Set("reblog_key", reblogKey)
	return doPost(client, "/blog/%s/post/reblog", blogName, setPostId(postId, params))
}

// ReblogOnBlog will reblog this Post to the blog in name.
func (p *PostRef) ReblogOnBlog(name string, params url.Values) (*PostRef, error) {
	return ReblogPost(p.client, name, p.Id, p.ReblogKey, params)
}

// DeletePost will delete a Post on tumblr for the blog in name and Post in postId.
func DeletePost(client ClientInterface, name string, postId uint64) error {
	_, err := client.PostWithParams(blogPath("/blog/%s/post/delete", name), setPostId(postId, url.Values{}))
	return err
}

// Delete will delete this Post on tumblr.
func (p *PostRef) Delete() error {
	return DeletePost(p.client, p.BlogName, p.Id)
}

// Utility function to create the proper instance of Post and return a reference to the generic interface
func makePostFromType(t string) (PostInterface, error) {
	switch t {
	case "quote":
		return &QuotePost{}, nil
	case "chat":
		return &ChatPost{}, nil
	case "photo":
		return &PhotoPost{}, nil
	case "text":
		return &TextPost{}, nil
	case "link":
		return &LinkPost{}, nil
	case "answer":
		return &AnswerPost{}, nil
	case "audio":
		return &AudioPost{}, nil
	case "video":
		return &VideoPost{}, nil
	}
	return &Post{}, errors.New(fmt.Sprintf("Unknown type %s", t))
}

// Like will like this Post on behalf of the current user.
func (p *PostRef) Like() error {
	return LikePost(p.client, p.Id, p.ReblogKey)
}

// Unlikes will unlike a Post on behalf of the current user.
func (p *PostRef) Unlike() error {
	return UnlikePost(p.client, p.Id, p.ReblogKey)
}

// Create an array of PostInterfaces based on the array of MiniPost objects provided
func makePostsFromMinis(minis []MiniPost, client ClientInterface) []PostInterface {
	posts := []PostInterface{}
	for _, mini := range minis {
		post, _ := makePostFromType(mini.Type)
		post.GetSelf().client = client
		posts = append(posts, post)
	}
	return posts
}
