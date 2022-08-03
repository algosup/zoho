package zoho

import (
	"testing"
)

func TestCreateContact(t *testing.T) {
	id, err := CreateContact(Contact{FirstName: "Test", LastName: "TEST", LeadSource: "Test"})
	if err != nil {
		panic(err)
	}

	id, err = AddContactNote(id, "This is a note", "Bla bla bla")
	if err != nil {
		panic(err)
	}

}
