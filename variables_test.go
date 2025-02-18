package zoho

import (
	"log"
	"testing"
	"time"
)

// go test -v -run TestZohoVariable

func TestZohoVariable(t *testing.T) {
	time := time.Now()
	err := SetLastAutoZoho(time)
	if err != nil {
		log.Fatal(err)
	}
	nt, err := GetLastAutoZoho()
	log.Println(nt)
	if err != nil {
		log.Fatal(err)
	}

	if time.UTC().Year() != nt.Year() {
		log.Fatal("bad month")
	}
	if time.UTC().Month() != nt.Month() {
		log.Fatal("bad month")
	}
	if time.UTC().Day() != nt.Day() {
		log.Fatal("bad month")
	}
	if time.UTC().Hour() != nt.Hour() {
		log.Fatal("bad hour")
	}
	if time.UTC().Minute() != nt.Minute() {
		log.Fatal("bad hour")
	}
}
