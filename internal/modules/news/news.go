package news

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Rocksus/pogo/configs"
)

var def Repository

type Repository interface {
	GetNewsByKeyword(parameter NewsSearchRequestParam) (Data, error)
	GetTopNews(parameter TopNewsRequestParam) (Data, error)
	GetNewsByTopic(parameter NewsTopicRequestParam) (Data, error)
}

func Init(config configs.NewsConfig) error {
	def = &newsRepo{
		APIKey: config.APIKey,
	}
	return nil
}

func (n *newsRepo) GetNewsByKeyword(parameter NewsSearchRequestParam) (Data, error) {
	var data Data
	requestURL := fmt.Sprintf("%s/search", apiURL)
	u, err := url.Parse(requestURL)
	if err != nil {
		return data, fmt.Errorf("[News][GetNewsByKeyword] Error on parsing query: %s", err.Error())
	}
	q := u.Query()
	q.Set("token", n.APIKey)
	if parameter.Query == "" {
		return data, fmt.Errorf("[News][GetNewsByKeyword] No search string provided")
	}
	q.Set("q", parameter.Query)
	if parameter.Country != "" {
		q.Set("country", parameter.Country)
	}
	if parameter.Language != "" {
		q.Set("lang", parameter.Language)
	}
	if parameter.Max != 0 {
		q.Set("max", fmt.Sprintf("%d", parameter.Max))
	}
	if parameter.Image {
		q.Set("image", "required")
	}
	if parameter.TitleOnly {
		q.Set("in", "title")
	}
	if !parameter.MinDate.IsZero() {
		q.Set("mindate", parameter.MinDate.Format(queryDateTimeFormat))
	}
	if !parameter.MaxDate.IsZero() {
		q.Set("maxdate", parameter.MaxDate.Format(queryDateTimeFormat))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return data, fmt.Errorf("[News][GetNewsByKeyword] Module Internal Error, %s", err.Error())
	}
	client := &http.Client{}
	rctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req = req.WithContext(rctx)

	resp, err := client.Do(req)
	if err != nil {
		return data, fmt.Errorf("[News][GetNewsByKeyword] Module Internal Error, %s", err.Error())
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, fmt.Errorf("[News][GetNewsByKeyword] Module Internal Error, %s", err.Error())
	}

	return data, nil
}

func (n *newsRepo) GetTopNews(parameter TopNewsRequestParam) (Data, error) {
	var data Data
	requestURL := fmt.Sprintf("%s/top-news", apiURL)
	u, err := url.Parse(requestURL)
	if err != nil {
		return data, fmt.Errorf("[News][GetTopNews] Error on parsing query: %s", err.Error())
	}
	q := u.Query()
	q.Set("token", n.APIKey)
	if parameter.Country != "" {
		q.Set("country", parameter.Country)
	}
	if parameter.Language != "" {
		q.Set("lang", parameter.Language)
	}
	if parameter.Max != 0 {
		q.Set("max", fmt.Sprintf("%d", parameter.Max))
	}
	if parameter.Image {
		q.Set("image", "required")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return data, fmt.Errorf("[News][GetTopNews] Module Internal Error, %s", err.Error())
	}
	client := &http.Client{}
	rctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req = req.WithContext(rctx)

	resp, err := client.Do(req)
	if err != nil {
		return data, fmt.Errorf("[News][GetTopNews] Module Internal Error, %s", err.Error())
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, fmt.Errorf("[News][GetTopNews] Module Internal Error, %s", err.Error())
	}

	return data, nil
}

func (n *newsRepo) GetNewsByTopic(parameter NewsTopicRequestParam) (Data, error) {
	var data Data
	if parameter.Topic == "" {
		return data, fmt.Errorf("[News][GetNewsByTopic] No topic found")
	}
	requestURL := fmt.Sprintf("%s/topics/%s", apiURL, parameter.Topic)
	u, err := url.Parse(requestURL)
	if err != nil {
		return data, fmt.Errorf("[News][GetNewsByTopic] Error on parsing query: %s", err.Error())
	}
	q := u.Query()
	q.Set("token", n.APIKey)
	if parameter.Country != "" {
		q.Set("lang", parameter.Country)
	}
	if parameter.Language != "" {
		q.Set("lang", parameter.Language)
	}
	if parameter.Max != 0 {
		q.Set("max", fmt.Sprintf("%d", parameter.Max))
	}
	if parameter.Image {
		q.Set("image", "required")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return data, fmt.Errorf("[News][GetNewsByTopic] Module Internal Error, %s", err.Error())
	}
	client := &http.Client{}
	rctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req = req.WithContext(rctx)

	resp, err := client.Do(req)
	if err != nil {
		return data, fmt.Errorf("[News][GetNewsByTopic] Module Internal Error, %s", err.Error())
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, fmt.Errorf("[News][GetNewsByTopic] Module Internal Error, %s", err.Error())
	}

	return data, nil
}
