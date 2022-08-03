package zoho

import (
	"testing"
	"time"
)

func TestCreateContact(t *testing.T) {
	id, err := CreateContact(Contact{
		FirstName:  "Test",
		LastName:   "TEST",
		LeadSource: "Test",
		OpenHouse:  &Date{time.Now()},
	})
	if err != nil {
		panic(err)
	}

	err = CancelOpenHouse(id)
	if err != nil {
		panic(err)
	}

	id, err = AddContactNote(id, "This is a note", "Bla bla bla")
	if err != nil {
		panic(err)
	}

}
