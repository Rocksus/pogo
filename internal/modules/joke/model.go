package joke

const (
	apiURL = "https://horrible-jokes.appspot.com"
)

// Data holds the data for the joke API
type Data struct {
	ID        int64  `json:"id"`
	Category  string `json:"category"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
	Language  string `json:"lang"`
	Response  int    `json:"status"`
	Error     string `json:"error"`
}

type jokeRepo struct {
}
