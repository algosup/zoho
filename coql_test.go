package zoho

import (
	"log"
	"testing"
)

type id struct {
	ID string `json:"id"`
}

// go test -run TestQuery01
func TestQuery01(t *testing.T) {
	//	calls, err := Select[id]("SELECT id FROM Contacts WHERE Phone is not null AND Pipeline='2023-2024'")
	calls, err := Select[id]("SELECT id FROM Contacts WHERE First_Name='Test'")
	if err != nil {
		panic(err)
	}
	log.Println("len:", len(calls))
}
