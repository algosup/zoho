package zoho

import (
	"log"
	"testing"
)

// go test -v -run TestZohoVariable

func TestZohoVariable(t *testing.T) {
	err := SetLastAutoZoho("2006-01-02T15:04:05-07:00")
	if err != nil {
		log.Fatal(err)
	}
	nt, err := GetLastAutoZoho()
	log.Println(nt)
	if err != nil {
		log.Fatal(err)
	}

	if nt != "2006-01-02T15:04:05-07:00" {
		log.Fatal()
	}
}
