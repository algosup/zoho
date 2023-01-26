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

func TestDeal(t *testing.T) {
	id, err := CreateContact(Contact{
		Email:      "fordeal@test.com",
		FirstName:  "ForDeal",
		LastName:   "FORDEAL",
		LeadSource: "Test",
	})
	if err != nil {
		panic(err)
	}
	_, err = CreateDeal(Deal{
		DealName:  "ForDeal FORDEAL",
		Stage:     "Candidature",
		ContactId: id,
		Pipeline:  "TEST",
	})
	if err != nil {
		panic(err)
	}
	err = DeleteContact(id)
	if err != nil {
		panic(err)
	}
}
