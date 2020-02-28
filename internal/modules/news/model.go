package news

import "time"

const (
	PublishedAtTimeFormat = "2006-01-02 15:04:05 -0700"
	queryDateTimeFormat   = "2006-01-02"
	apiURL                = "https://gnews.io/api/v3"
)

type newsRepo struct {
	APIKey string
}

type NewsSearchRequestParam struct {
	Query     string
	Language  string
	Country   string
	Max       int  // maximum is 100
	Image     bool // required or optional
	MinDate   time.Time
	MaxDate   time.Time
	TitleOnly bool // all or title
}

type TopNewsRequestParam struct {
	Language string
	Country  string
	Max      int  // maximum is 100
	Image    bool // required or optional
}

type NewsTopicRequestParam struct {
	Topic    string
	Language string
	Country  string
	Max      int  // maximum is 100
	Image    bool // required or optional
}

type Data struct {
	Timestamp    int64     `json:"timestamp"`
	ArticleCount int64     `json:"articleCount"`
	Articles     []Article `json:"articles"`
}

type Article struct {
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	URL               string     `json:"url"`
	ImageURL          string     `json:"image"`
	PublishedAtString string     `json:"publishedAt"`
	Source            SourceData `json:"source"`
}

type SourceData struct {
	Name string `json:"name"`
	URL  string `json:""`
}
