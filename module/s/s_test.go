package s

import "testing"

var svc *SService

func TestS(t *testing.T) {
	sub := "/guinea/package"
	written := "I had a guinea golden - I lost it in the sand"
	expectedSub := "I had a package golden - I lost it in the sand"

	subbedWritten := svc.SubStringToStatement(sub, written)

	if expectedSub != subbedWritten {
		t.Errorf("Failed to substitute - got:[%s] when expecting:[%s]", subbedWritten, expectedSub)
	}
}
