package api

import (
	"testing"
)

func TestReturn(t *testing.T) {
	if 1 != api_return(1) {
		t.Error("api return should be 1")
	}
}

func TestIncrease(t *testing.T) {
	if 1 != api_increase(1) {
		t.Error("api return should be 1")
	}
}
