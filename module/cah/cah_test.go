package cah

import "testing"

// import "log"

var svc *CahService

func init() {
	svc = NewService()
}

func TestCah1(t *testing.T) {

	question := "The _ went to the _."
	answers := []string{"cat", "food store."}
	expectedAnswer := "The cat went to the food store."

	testAnswer := svc.MessageFromQuestionAndAnswers(question, answers)

	if expectedAnswer != testAnswer {
		t.Errorf("Failed - got:[%s] when expecting:[%s]", testAnswer, expectedAnswer)
	}
}

func TestCah2(t *testing.T) {

	question := "Who ate me?"
	answers := []string{"The Cat"}
	expectedAnswer := "Who ate me? The Cat."

	testAnswer := svc.MessageFromQuestionAndAnswers(question, answers)

	if expectedAnswer != testAnswer {
		t.Errorf("Failed - got:[%s] when expecting:[%s]", testAnswer, expectedAnswer)
	}
}
