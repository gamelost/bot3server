package conditions

import (
	"testing"
	"log"
)

/*
func TestGetWeatherReportForAirportCode(t *testing.T) {
	str, _ := getWeatherReportForAirportCode("LAX")
	log.Printf("%s", str)
}
*/

func TestGetWeatherConditionForLocation1(t *testing.T) {
	str, _ := getWeatherConditionsForLocation("lottery tickets")
	log.Printf("%s", str)
}