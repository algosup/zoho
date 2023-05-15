package zoho

import (
	"testing"
)

// go test -run TestGetContactCalls
func TestGetContactCalls(t *testing.T) {
	_, err := GetContactCalls("477339000006512947")
	if err != nil {
		panic(err)
	}
}
