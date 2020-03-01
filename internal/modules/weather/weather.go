package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Rocksus/pogo/configs"
)

var def Repository

// Repository holds the package module's available functions
type Repository interface {
	QueryLocation(location string) (Data, error)
}

// QueryLocation gets the weather data of a city based on locationID
func (w *weatherRepo) QueryLocation(location string) (Data, error) {
	var data Data
	locationID, ok := locationIDs[location]
	if !ok {
		return data, fmt.Errorf("Invalid location string")
	}
	requestURL := fmt.Sprintf("%sid=%d&APPID=%s", weatherURL, locationID, w.APIKey)
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

	data.Details.TemperatureCelcius = kelvinToCelcius(data.Details.TemperatureKelvin)

	return data, nil
}

// QueryLocation is the global level of the function
func QueryLocation(location string) (Data, error) {
	return def.QueryLocation(location)
}

// Init sets up the default variable to handle package functions
func Init(config configs.WeatherConfig) error {
	newRepo := &weatherRepo{
		APIKey: config.APIKey,
	}
	_, err := newRepo.QueryLocation("jakarta")
	if err != nil {
		return fmt.Errorf("[Weather][Init] Error initalizing the weather module, err: %s", err.Error())
	}
	log.Print("[Weather][Init] Weather module initialized successfully.")
	def = newRepo
	return nil
}
