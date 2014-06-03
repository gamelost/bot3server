package forecast

import (
	"errors"
	"fmt"
	"github.com/gamelost/bot3server/module/weather"
	"github.com/gamelost/bot3server/server"
	"log"
	"net/url"
	"strings"
)

type WeatherForecastService struct{}

var stateCityAPICallUrl string
var cityAPICallUrl string
var airportAPICallUrl string

func (svc *WeatherForecastService) NewService() server.BotHandler {
	apiurl, _ := server.ServerConfig.GetString("weather", "wundergroundapiurl")
	apikey, _ := server.ServerConfig.GetString("weather", "wundergroundapikey")
	apipath := apiurl + apikey
	stateCityAPICallUrl = apipath + "/forecast/q/%s/%s.json"
	cityAPICallUrl = apipath + "/forecast/q/%s.json"
	airportAPICallUrl = apipath + "/forecast/q/%s.json"
	return &WeatherForecastService{}
}

func (svc *WeatherForecastService) Handle(botRequest *server.BotRequest, botResponse *server.BotResponse) {

	var err error
	var resp []string

	wStr := strings.TrimSpace(botRequest.LineTextWithoutCommand())
	weatherCmd := &weather.WeatherLocation{Location: wStr}

	resp, err = getWeatherForecastForLocation(weatherCmd.Location)

	if err != nil {
		log.Printf("error is: %v", err)
		botResponse.SetSingleLineResponse(err.Error())
	} else {
		botResponse.SetMultipleLineResponse(resp)
	}
}

func getWeatherForecastForLocation(command string) (weatherResponse []string, err error) {

	// split if there is a state
	args := strings.SplitN(command, ",", 2)

	// do we have two args?
	if len(args) == 2 {

		state := url.QueryEscape(strings.Replace(strings.TrimSpace(args[1]), " ", "_", -1))
		city := url.QueryEscape(strings.Replace(strings.TrimSpace(args[0]), " ", "_", -1))
		call := fmt.Sprintf(stateCityAPICallUrl, state, city)
		log.Printf(call)
		weatherData, callErr := weather.DoWeatherAPICall(call)

		//log.Printf("data: %v", weatherData)
		//log.Printf("error: %s", callErr.Error())
		if callErr != nil {
			err = callErr
			return
		}

		if len(weatherData.Response.Results) > 0 {
			err = errors.New("WUG API could not find exact match on this?  Your father was a hamster.")
			return
		}

		weatherResponse = parseWeatherDataIntoResponseString(weatherData)

	} else if len(args) == 1 {

		city := url.QueryEscape(strings.Replace(strings.TrimSpace(args[0]), " ", "_", -1))
		callUrl := fmt.Sprintf(cityAPICallUrl, city)
		weatherData, callErr := weather.DoWeatherAPICall(callUrl)

		// if an error in response
		if callErr != nil {
			err = callErr
			return
		}

		if len(weatherData.Response.Results) > 0 {
			err = errors.New("Multiple results returned.  Try to specify a state along with your city?")
			return
		}

		weatherResponse = parseWeatherDataIntoResponseString(weatherData)

	} else {
		err = errors.New("Unable to parse location string.  Try it in 'City,State' format.")
	}

	return
}

func parseWeatherDataIntoResponseString(weatherData *weather.WUAPIResponse) []string {

	response := make([]string, 3)

	response[0] = fmt.Sprintf("Weather forecast for %s is %s", weatherData.Forecast.Txt_forecast.Forecastday[0].Title, weatherData.Forecast.Txt_forecast.Forecastday[0].Fcttext)
	response[1] = fmt.Sprintf("%s is %s", weatherData.Forecast.Txt_forecast.Forecastday[1].Title, weatherData.Forecast.Txt_forecast.Forecastday[1].Fcttext)
	response[2] = fmt.Sprintf("%s is %s", weatherData.Forecast.Txt_forecast.Forecastday[2].Title, weatherData.Forecast.Txt_forecast.Forecastday[2].Fcttext)
	return response
}
