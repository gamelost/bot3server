package weather 

import (
	"testing"
)

func TestIsZipCode1(t *testing.T) {
	
	loc := &WeatherLocation{ Location: "55405" }
	if !loc.IsZipCode() {
		t.Fail()
	}
}

func TestIsZipCode3(t *testing.T) {
	
	loc := &WeatherLocation{ Location: "aug" }
	if loc.IsZipCode() {
		t.Fail()
	}
}

func TestIsZipCode5(t *testing.T) {
	
	loc := &WeatherLocation{ Location: "11" }
	if loc.IsZipCode() {
		t.Fail()	
	}
}