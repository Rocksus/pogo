package gnews

const (
	PublishedAtTimeFormat = "2006-01-02 15:04:05 -0700"
	apiURL                = "https://gnews.io/api/v3/search?"
)

type NewsParameter 

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
