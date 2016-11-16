package tumblr

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"net/url"
)

// Representation of a list of Posts
type Posts struct {
	client ClientInterface
	response Response
	parsedPosts []PostInterface
	Posts []MiniPost `json:"posts"`
	TotalPosts int64 `json:"total_posts"`
}

// Method to retrieve fully fleshed post data from stubs and cache result
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

// Method to retrieve a single Post entity at a given index; returns nil if index is out of bounds
func (p *Posts) Get(index uint) (PostInterface) {
	if posts,err := p.All(); err == nil {
		if index >= uint(len(posts)) {
			return nil
		}
		return posts[index]
	}
	return nil
}

// The basics for what is needed in a Post
type MiniPost struct {
	Id uint64 `json:"id"`
	Type string `json:"type"`
	BlogName string `json:"blog_name"`
	ReblogKey string `json:"reblog_key"`
}

// Starting point for performing operations on a post
type PostRef struct {
	MiniPost
	client ClientInterface
}

// The common fields on any post, no matter what type
type Post struct {
	PostRef
	Body string `json:"body"`
	CanLike bool `json:"can_like"`
	CanReblog bool `json:"can_reblog"`
	CanReply bool `json:"can_reply"`
	CanSendInMessage bool `json:"can_send_in_message"`
	Caption string `json:"caption"`
	Date string `json:"date"`
	DisplayAvatar bool `json:"display_avatar"`
	Followed bool `json:"followed"`
	Format string `json:"format"`
	Highlighted []interface{} `json:"highlighted"`
	Liked bool `json:"liked"`
	NoteCount uint64 `json:"note_count"`
	PermalinkUrl string `json:"permalink_url"`
	PostUrl string `json:"post_url"`
	Reblog struct {
		 Comment string `json:"comment"`
		 TreeHTML string `json:"tree_html"`
	       } `json:"reblog"`
	RecommendedColor string `json:"recommended_color"`
	RecommendedSource bool `json:"recommended_source"`
	ShortUrl string `json:"short_url"`
	Slug string `json:"slug"`
	SourceTitle string `json:"source_title"`
	SourceUrl string `json:"source_url"`
	State string `json:"state"`
	Summary string `json:"summary"`
	Tags []string `json:"tags"`
	Timestamp uint64 `json:"timestamp"`
	FeaturedTimestamp uint64 `json:"featured_timestamp,omitempty"`
	TrackName string `json:"track_name,omitempty"`
	Trail []ReblogTrailItem `json:"trail"`
}

// Post substructure
type ReblogTrailItem struct {
	Blog Blog `json:"blog"`
	Content string `json:"content"`
	ContentRaw string `json:"content_raw"`
	IsCurrentItem bool `json:"is_current_item"`
	Post struct {
		     // sometimes an actual int, sometimes a numeric string, always a headache
		     Id interface{} `json:"id"`
	     } `json:"post"`
}

// PostInterface for use in typed structures which could contain any of the below subtypes
type PostInterface interface {
	GetProperty(key string) (interface{}, error)
	GetSelf() (*Post)
}

// Post subtype
type QuotePost struct {
	Post
	Source string `json:"source,omitempty"`
	Text string `json:"text"`
}

// Post subtype
type ChatPost struct {
	Post
	Dialog []struct{
		Label string `json:"label"`
		Name string `json:"name"`
		Phrase string `json:"phrase"`
	} `json:"dialog"`
}

// Post subtype
type TextPost struct {
	Post
	Title string `json:"title"`
}

// Post subtype
type LinkPost struct {
	Post
	Description string `json:"description"`
	Excerpt string `json:"excerpt"`
	LinkAuthor string `json:"link_author"`
	Title string `json:"title"`
	Url string `json:"url"`
}

// Post subtype
type AnswerPost struct {
	Post
	Answer string `json:"answer"`
	AskingName string `json:"asking_name"`
	AskingUrl string `json:"asking_url"`
	Publisher string `json:"publisher"`
	Question string `json:"question"`
}

// Post subtype
type AudioPost struct {
	Post
	AlbumArt string `json:"album_art"`
	Artist string `json:"artist"`
	AudioSourceUrl string `json:"audio_source_url"`
	AudioType string `json:"audio_type"`
	AudioUrl string `json:"audio_url"`
	Embed string `json:"embed"`
	Player string `json:"player"`
	Plays uint64 `json:"plays"`
}

// Post subtype
type VideoPost struct {
	Post
	Html5Capable bool `json:"html5_capable"`
	PermalinkUrl string `json:"permalink_url"`
	Players []struct {
		EmbedCode string `json:"embed_code"`
		Width interface{} `json:"width"`
	} `json:"player"`
	ThumbnailHeight uint32 `json:"thumbnail_height"`
	ThumbnailUrl string `json:"thumbnail_url"`
	ThumbnailWidth uint32 `json:"thumbnail_width"`
	Video map[string]struct {
		Height uint32 `json:"height"`
		Width uint32 `json:"width"`
		VideoId string `json:"video_id"`
	} `json:"video"`
	VideoType string `json:"video_type"`
}

// Post subtype
type PhotoPost struct {
	Post
	ImagePermalink string `json:"image_permalink"`
	Photos []Photo `json:"photos"`
}

// Photo post substructure
type Photo struct {
	AltSizes []PhotoSize `json:"alt_sizes"`
	Caption string `json:"caption"`
	OriginalSize PhotoSize `json:"original_size"`
}

// Photo substructure
type PhotoSize struct {
	Height uint32 `json:"height"`
	Width uint32 `json:"width"`
	Url string `json:"url"`
}

// Convenience method for ease of use- renders a Post as a JSON string
func (p *Post) String() string {
	return jsonStringify(*p)
}

// Convenience method for easy retrieval of one-off values
func (p *Post) GetProperty(key string) (interface{},error) {
	if field,exists := reflect.TypeOf(p).Elem().FieldByName(key); exists {
		return reflect.ValueOf(p).Elem().FieldByIndex(field.Index), nil
	}
	return nil, errors.New(fmt.Sprintf("Property %s does not exist", key))
}

// Useful for converting a PostInterface into a Post
func (p *Post) GetSelf() (*Post) {
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

// Retrieve a blog's posts, in the API docs you can find how to filter by ID, type, etc
func GetPosts(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts", name, params)
}

// Retrieve a blog's Queue
func GetQueue(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts/queue", name, params)
}

// Retrieve a blog's drafts
func GetDrafts(client ClientInterface, name string, params url.Values) (*Posts, error) {
	return queryPosts(client, "/blog/%s/posts/draft", name, params)
}

// Retrieve a blog's submsisions
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
		Response struct{
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

// Creates a PostRef with the given properties set
func NewPostRefById(client ClientInterface, id uint64) (*PostRef) {
	return &PostRef{
		client: client,
		MiniPost: MiniPost{
			Id: id,
		},
	}
}

// Creates a PostRef with the given properties set
func NewPostRef(client ClientInterface, post *MiniPost) (*PostRef) {
	return &PostRef{
		client: client,
		MiniPost: *post,
	}
}

// Sets PostRef's client
func (r *PostRef)SetClient(c ClientInterface) {
	r.client = c
}

// Create a post, return the ID on success, error on failure
func CreatePost(client ClientInterface, name string, params url.Values) (*PostRef, error) {
	return doPost(client, "/blog/%s/post", name, params)
}

// Edit a given post, returns nil if successful, error on failure
func EditPost(client ClientInterface, blogName string, postId uint64, params url.Values) error {
	_, err := client.PostWithParams(blogPath("/blog/%s/post/edit", blogName), setPostId(postId, params))
	return err
}

// Convenience method to allow calling post.Edit(params)
func (p *PostRef) Edit(params url.Values) error {
	return EditPost(p.client, p.BlogName, p.Id, params)
}

// Reblog a given post to the given blog, returns the reblog's post id if successful, else the error
func ReblogPost(client ClientInterface, blogName string, postId uint64, reblogKey string, params url.Values) (*PostRef, error) {
	if reblogKey == "" {
		return nil, errors.New("No reblog key provided")
	}
	params.Set("reblog_key", reblogKey)
	return doPost(client, "/blog/%s/post/reblog", blogName, setPostId(postId, params))
}

// Convenience method to allow calling post.Reblog(params)
func (p *PostRef) ReblogOnBlog(name string, params url.Values) (*PostRef, error) {
	return ReblogPost(p.client, name, p.Id, p.ReblogKey, params)
}

// Delete a given blog's post by ID, nil if successful, error on failure
func DeletePost(client ClientInterface, name string, postId uint64) error {
	_, err := client.PostWithParams(blogPath("/blog/%s/post/delete", name), setPostId(postId, url.Values{}))
	return err
}

// Convenience method to allow calling post.Delete()
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

// Likes a Post on behalf of the current user
func (p *PostRef) Like() error {
	return LikePost(p.client, p.Id, p.ReblogKey)
}

// Unlikes a Post on behalf of the current user
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