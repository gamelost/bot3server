package catfacts

import "testing"

var svc *CatFactsService

func init() {
	svc = svc.NewService().(*CatFactsService)
}

func TestCatFacts(t *testing.T) {
}
