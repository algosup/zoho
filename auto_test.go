package zoho

import (
	"log"
	"testing"
)

func TestAuto(t *testing.T) {
	err := AutoUpdateContact("477339000004039009")
	if err != nil {
		panic(err)
	}
}

func TestAutoAll(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := AutoUpdateAllContacts()
	if err != nil {
		panic(err)
	}
}
