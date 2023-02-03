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

func TestCreateDelete(t *testing.T) {
	_, err := CreateContact(Contact{
		Email:      "test132@test.com",
		FirstName:  "Test",
		LastName:   "TEST TEST",
		LeadSource: "Test",
		OpenHouse:  &Date{time.Now()},
	})
	if err != nil {
		panic(err)
	}

	id, err := FindContact("test132@test.com")
	if err != nil {
		panic(err)
	}
	if id == "" {
		panic("not found")
	}

	_, err = GetContact(id)
	if err != nil {
		panic(err)
	}
	err = DeleteContact(id)
	if err != nil {
		panic(err)
	}
}

func TestFindContact2(t *testing.T) {
	id, err := FindContact("eb232235@gmail.com")
	if err != nil {
		panic(err)
	}
	_, err = GetContact(id)
	if err != nil {
		panic(err)
	}
}

func TestContactEmails(t *testing.T) {
	_, err := GetContactEmails("477339000004039009")
	if err != nil {
		panic(err)
	}
}

func TestContactNotes(t *testing.T) {
	c, d, err := GetContactNotesCount("477339000004039009")
	if err != nil {
		panic(err)
	}
	t.Log(c, d)
}
