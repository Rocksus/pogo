package gnews

import "github.com/Rocksus/pogo/configs"

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

}

func (n *newsRepo) GetTopNews(parameter map[string]interface{}) (Data, error) {

}

func (n *newsRepo) GetNewsByTopic(parameter map[string]interface{}) (Data, error) {

}
