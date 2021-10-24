package joke

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/pkg/errors"

	"github.com/Rocksus/pogo/pkg/plugin"
)

type Plugin struct {
	client *http.Client
}

func InitPlugin() *Plugin {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &Plugin{
		client: client,
	}
}

func (p *Plugin) Reply(ctx context.Context, message plugin.Message, replyCh chan<- linebot.SendingMessage) {
	joke, err := p.getRandomJoke(ctx)
	if err != nil {
		replyCh <- p.getDefaultMessage()
		return
	}

	replyCh <- linebot.NewTextMessage(joke.Setup)
	time.Sleep(2 * time.Second)
	replyCh <- linebot.NewTextMessage(joke.Punchline)
}

func (p *Plugin) getRandomJoke(ctx context.Context) (data Data, err error) {
	requestURL := fmt.Sprintf("%s/jokes", apiURL)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		err = errors.WithMessage(err, "failed to create request")
		return
	}
	req = req.WithContext(ctx)

	resp, err := p.client.Do(req)
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
	if data.Response != 200 {
		err = errors.WithMessage(errors.New(data.Error), "api error")
		return
	}

	return data, nil
}

func (p *Plugin) getDefaultMessage() linebot.SendingMessage {
	defaults := []string{
		"Boo hoo, I ran out of funny juice.",
		"I can't come up with a joke.",
		"I'm sorry, I'm not funny :(. Come back later when I can think of a good one",
		"Knock knock? My joke machine is broken.",
		"Jokes are not in my dictionary. I need to update my dictionary stats.",
	}

	i := rand.Intn(len(defaults))
	return linebot.NewTextMessage(defaults[i])
}

func (p *Plugin) getRandomJokeByCategory(ctx context.Context, category string) (Data, error) {
	var data Data
	requestURL := fmt.Sprintf("%s/jokes/%s", apiURL, category)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return data, fmt.Errorf("Module Internal Error, %s", err.Error())
	}
	client := &http.Client{}
	rctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req = req.WithContext(rctx)

	resp, err := client.Do(req)
	if err != nil {
		return data, fmt.Errorf("Module Internal Error, %s", err.Error())
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, fmt.Errorf("Module Internal Error, %s", err.Error())
	}
	if data.Response != 200 {
		return data, fmt.Errorf("API Error, %s", data.Error)
	}

	return data, nil

}

func (p *Plugin) getJokeByID(id int64) (Data, error) {
	var data Data
	requestURL := fmt.Sprintf("%s/joke/%d", apiURL, id)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return data, fmt.Errorf("Module Internal Error, %s", err.Error())
	}
	client := &http.Client{}
	rctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req = req.WithContext(rctx)

	resp, err := client.Do(req)
	if err != nil {
		return data, fmt.Errorf("Module Internal Error, %s", err.Error())
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, fmt.Errorf("Module Internal Error, %s", err.Error())
	}
	if data.Response != 200 {
		return data, fmt.Errorf("API Error, %s", data.Error)
	}

	return data, nil

}
