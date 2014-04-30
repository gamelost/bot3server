package nextwedding

import (
	"testing"
	"log"
)

func TestDuration(t *testing.T) {
	output := durationToWeddingDate()
	log.Println(output)
}