package weather

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const weatherProvider = "http://api.apixu.com/v1"
const currentContext = "current.json"

//const forecastContext = "forecast.json"

// NewWeather returns Weather struct
// func NewWeather(key string) (*Weather, error) {
// 	var w Weather
// 	w.Key = key

// 	return nil, nil
// }

// GetCurrentWeather returnes parsed JSON for current weather response
func GetCurrentWeather(conf *Weather) (*CurWeather, error) {
	url := weatherProvider + "/" + currentContext + "?key=" + conf.Key + "&q=Moscow&lang=ru"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error in GetCurrentWeather func while calling API: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var result CurWeather

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error in GetCurrentWeather func while Unmarshalling: %s", err)
		return nil, err
	}

	return &result, nil

}
