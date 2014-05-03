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

const stateCityAPICallUrl = "http://api.wunderground.com/api/84a6128036456b39/forecast/q/%s/%s.json"
const cityAPICallUrl = "http://api.wunderground.com/api/84a6128036456b39/forecast/q/%s.json"
const airportAPICallUrl = "http://api.wunderground.com/api/84a6128036456b39/forecast/q/%s.json"

func (svc *WeatherForecastService) NewService() server.BotHandler {
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
		botResponse.Response = resp
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
