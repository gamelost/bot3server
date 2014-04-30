package forecast

import (
	"testing"
	"log"
)

func TestGetWeatherForecastForLocation(t *testing.T) {
	str, err := getWeatherForecastForLocation("Monterey Bay, ca")
	
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	
	log.Printf("%s", str)
}