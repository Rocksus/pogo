package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Rocksus/pogo/pkg/plugin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nickylogan/go-log"
	"github.com/pkg/errors"
)

type Plugin struct {
	client *http.Client
	apiKey string
}

// InitPlugin sets up the default variable to handle package functions
func InitPlugin(apiKey string) *Plugin {
	p := &Plugin{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	return p
}

func (p *Plugin) Reply(ctx context.Context, message plugin.Message, replyCh chan<- linebot.SendingMessage) {
	greet, fail := p.generateGreetingPair()
	replyCh <- linebot.NewTextMessage(greet)

	// TODO: temporarily set to jakarta only. Extract from entities later on
	data, err := p.queryLocation(ctx, "jakarta")
	if err != nil {
		replyCh <- linebot.NewTextMessage(fail)
	}

	replyCh <- linebot.NewTextMessage(fmt.Sprintf("Got it! Here's the weather in %s, %s", data.Name, data.System.Country))
	replyCh <- linebot.NewTextMessage(fmt.Sprintf("%s: %s", data.Weather[0].Main, data.Weather[0].Description))
	replyCh <- linebot.NewTextMessage(fmt.Sprintf("Humidity: %d", data.Details.Humidity))
	replyCh <- linebot.NewTextMessage(fmt.Sprintf("Temperature: %.2f°C", data.Details.TemperatureCelcius))
}

// QueryLocation gets the weather data of a city based on locationID
func (p *Plugin) queryLocation(ctx context.Context, location string) (data Data, err error) {
	locationID, ok := locationIDs[location]
	if !ok {
		log.Errorln("unknown location:", location)
		return data, ErrUnknownLocation
	}

	requestURL := fmt.Sprintf("%sid=%d&APPID=%s", weatherURL, locationID, p.apiKey)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return data, errors.WithMessage(err, "failed to create request")
	}

	resp, err := p.client.Do(req.WithContext(ctx))
	if err != nil {
		return data, errors.WithMessage(err, "failed to do request")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return data, errors.WithMessage(err, "failed to decode json")
	}
	if data.Response != 200 {
		return data, errors.WithMessage(errors.New(data.Message), "api error")
	}

	data.Details.TemperatureCelcius = kelvinToCelcius(data.Details.TemperatureKelvin)

	return data, nil
}

func (p *Plugin) generateGreetingPair() (greeting, fail string) {
	pairs := [][2]string{
		{"Checking my weather machine...", "Aww my weather machine is broken. Check back again later?"},
		{"Asking the weather gods...", "The weather gods are not responding. Check back again later?"},
	}

	i := rand.Intn(len(pairs))
	return pairs[i][0], pairs[i][1]
}
