package register

import (
	"testing"

	"github.com/google/uuid"
)

func TestUUIDCreation(t *testing.T) {
	uuidAsString, err := uuid.FromBytes([]byte("S167357X6205283" + " "))
	if err != nil {
		t.Error(err)
	}
	t.Logf("got: %s", uuidAsString)

	uuidAsString2, err := uuid.FromBytes([]byte("S167357X6205283" + " "))
	if err != nil {
		t.Error(err)
	}
	if uuidAsString != uuidAsString2 {
		t.Errorf("expected same uuid, got different: %s vs: %s", uuidAsString, uuidAsString2)
	}
}
