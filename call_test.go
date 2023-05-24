package zoho

import (
	"log"
	"testing"
	"time"
)

// go test -run TestGetContactCalls
func TestGetContactCalls(t *testing.T) {
	calls, err := GetOpenCalls("477339000003902012")
	if err != nil {
		panic(err)
	}
	log.Println(calls)
}

// go test -run TestGetCalls
func TestGetCalls(t *testing.T) {
	calls, err := GetCalls("477339000003902012")
	if err != nil {
		panic(err)
	}
	log.Println(calls)
}

// go test -run TestCreateCall
func TestCreateCall(t *testing.T) {
	_, err := CreateCall(NewCall{
		Subject:       "Hello",
		CallType:      "Sortant",
		CallStartTime: Time{time.Now()},
		CallDuration:  "0",
	}, "477339000003902012")
	if err != nil {
		panic(err)
	}
}
