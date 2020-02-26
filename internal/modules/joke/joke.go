package joke

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var def Repository

type Repository interface {
	GetRandomJoke() (Data, error)
	GetRandomJokeByCategory(category string) (Data, error)
	GetJokeByID(id int64) (Data, error)
}

func Init() {
	newRepo := &jokeRepo{}
	def = newRepo
	log.Print("[Joke][Init] Joke module initialized successfully.")
}

func (j *jokeRepo) GetRandomJoke() (Data, error) {
	var data Data
	requestURL := fmt.Sprintf("%s/jokes", apiURL)

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
		return data, fmt.Errorf("API Error, %s", data.Message)
	}

	return data, nil

}

func (j *jokeRepo) GetRandomJokeByCategory(category string) (Data, error) {
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
		return data, fmt.Errorf("API Error, %s", data.Message)
	}

	return data, nil

}

func (j *jokeRepo) GetJokeByID(id int64) (Data, error) {
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
		return data, fmt.Errorf("API Error, %s", data.Message)
	}

	return data, nil

}

func GetRandomJoke() (Data, error) {
	return def.GetRandomJoke()
}

func GetRandomJokeByCategory(category string) (Data, error) {
	return def.GetRandomJokeByCategory(category)
}

func GetJokeByID(id int64) (Data, error) {
	return def.GetJokeByID(id)
}
