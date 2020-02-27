package gnews

import (
	"fmt"

	"github.com/Rocksus/pogo/configs"
)

var def Repository

type Repository interface {
	GetNewsByKeyword(parameter map[string]interface{}) (Data, error)
	GetTopNews(parameter map[string]interface{}) (Data, error)
	GetNewsByTopic(parameter map[string]interface{}) (Data, error)
}

func Init(config configs.NewsConfig) error {
	def = &newsRepo{
		APIKey: config.APIKey,
	}
	return nil
}

func (n *newsRepo) GetNewsByKeyword(parameter map[string]interface{}) (Data, error) {
	requestURL := fmt.Sprintf("%s/search?", apiURL)
}

func (n *newsRepo) GetTopNews(parameter map[string]interface{}) (Data, error) {
	requestURL := fmt.Sprintf("%s/top-news?token=%s", n.APIKey)
}

func (n *newsRepo) GetNewsByTopic(parameter map[string]interface{}) (Data, error) {
	requestURL := fmt.Sprintf("%s/topics/%s?token=%s", apiURL, "topic", n.APIKey)
}
