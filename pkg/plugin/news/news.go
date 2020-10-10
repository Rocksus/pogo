package news

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nickylogan/go-log"
	"github.com/pkg/errors"
)

type Plugin struct {
	client *http.Client
	apiKey string
}

func InitPlugin(apiKey string) *Plugin {
	return &Plugin{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (p *Plugin) GetDaily(ctx context.Context, recipientID string) (linebot.SendingMessage, error) {
	news, err := p.getTopNews(ctx, TopNewsRequestParam{Max: 3})
	if err != nil {
		log.WithError(err).Errorln("failed to get top news")
		return nil, err
	}

	var sb strings.Builder
	sb.WriteString("Top news for today:\n")
	for _, v := range news.Articles {
		sb.WriteString(fmt.Sprintf(
			"%s: %s\n"+
				"%s",
			v.Source.Name, v.Title,
			v.URL,
		))
	}
	return linebot.NewTextMessage(sb.String()), nil
}

func (p *Plugin) getNewsByKeyword(ctx context.Context, parameter NewsSearchRequestParam) (data Data, err error) {
	// validate params
	if parameter.Query == "" {
		return data, errors.New("no search string provided")
	}

	// apiURL is guaranteed to be valid
	u, _ := url.Parse(fmt.Sprintf("%s/search", apiURL))

	q := u.Query()
	q.Set("token", p.apiKey)
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

	data, err = p.doRequest(ctx, u.String())
	return
}

func (p *Plugin) getTopNews(ctx context.Context, parameter TopNewsRequestParam) (data Data, err error) {
	// apiURL is guaranteed to be valid
	u, _ := url.Parse(fmt.Sprintf("%s/top-news", apiURL))

	q := u.Query()
	q.Set("token", p.apiKey)
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

	data, err = p.doRequest(ctx, u.String())
	return
}

func (p *Plugin) getNewsByTopic(ctx context.Context, parameter NewsTopicRequestParam) (data Data, err error) {
	if parameter.Topic == "" {
		err = errors.New("no topic provided")
		return
	}

	u, err := url.Parse(fmt.Sprintf("%s/topics/%s", apiURL, parameter.Topic))
	if err != nil {
		err = errors.WithMessage(err, "failed to parse url")
		return
	}

	q := u.Query()
	q.Set("token", p.apiKey)
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

	data, err = p.doRequest(ctx, u.String())
	return
}

func (p *Plugin) doRequest(ctx context.Context, url string) (data Data, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err = errors.WithMessage(err, "failed to create request")
		return
	}

	resp, err := p.client.Do(req.WithContext(ctx))
	if err != nil {
		err = errors.WithMessage(err, "failed to do request")
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		err = errors.WithMessage(err, "failed to decode json")
		return
	}
	return
}
