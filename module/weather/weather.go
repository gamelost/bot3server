package weather 

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"regexp"
	"log"
)

type WeatherLocation struct {
	Location string
}

type WUAPIResponse struct {
	
	Response APIResponse
	Current_observation APICurrentObservation
	Forecast APIForecast
}

type APIResponse struct {
	Version string
	Error APIError
	Results []map[string]interface{}
}

type APIError struct {
	Type string 
	Description string
}

type APIDisplayLocation struct {
	Full string
}

type APIForecast struct {
	Txt_forecast APITextForecast
}

type APITextForecast struct {
	Date string
	Forecastday []APIForecastDay
}

type APIForecastDay struct {
	Period int
	Title string
	Fcttext string
}

type APICurrentObservation struct {
	Display_location APIDisplayLocation
	Weather string
	Temperature_string string 
	Wind_string string 
}

func (r *WUAPIResponse) isResponseError() bool {
	
	if r.Response.Error.Type != "" {
		return true
	} else {
		return false
	}
}

func DoWeatherAPICall(call string) (weatherData *WUAPIResponse, err error) {
	
	resp, err := http.Get(call)
	if err != nil {
		return
	}	
	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return
	}
	
	// check for api error
	if weatherData.isResponseError() {
		log.Println("WUG api reported error.")
		err = errors.New("Weather Underground API reported an error.")
		return
	}
	
	log.Println("WUG call succeeeded")
	return
}

func (wl *WeatherLocation) IsZipCode() bool {
	
	matched, err := regexp.MatchString("^\\d{5}([\\-]?\\d{4})?$", wl.Location)
	
	if err != nil {
		return false
	} else {
		return matched
	}
}