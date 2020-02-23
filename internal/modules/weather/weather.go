package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Rocksus/pogo/configs"
)

var def Repository

type Repository interface {
}

func (w *weatherRepo) testAPI() error {
	requestURL := fmt.Sprintf("%s%s%s", weatherURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return fmt.Errorf("Module Internal Error, %s", err.Error())
	}
	client := &http.Client{}
	rctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req = req.WithContext(rctx)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Module Internal Error, %s", err.Error())
	}
	defer resp.Body.Close()

	var data Data

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("Module Internal Error, %s", err.Error())
	}
	if data.Response != 200 {
		return fmt.Errorf("API Error, %s", data.Message)
	}

	return nil
}

func Init(config configs.WeatherConfig) error {
	newRepo := &weatherRepo{
		APIKey: config.APIKey,
	}
	err := newRepo.testAPI()
	if err != nil {
		return errors.New(fmt.Sprintf("[Weather][Init] Error initalizing the weather module, err: %s"), err.Error())
	}
	def = newRepo
	return nil
}
