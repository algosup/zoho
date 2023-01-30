package zoho

import (
	"testing"
)

/*
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
			Stage:     "Prospect",
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
*/
func TestDoubleDeal(t *testing.T) {
	cid, err := FindContact("fordeal@test.com")
	if err != nil {
		panic(err)
	}

	did, err := FindDealByContactID(cid)
	if err != nil {
		panic(err)
	}

	if did == "" {
		panic("deal not found")
	}

}
