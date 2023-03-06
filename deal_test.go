package zoho

import "testing"

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

	/*err = DeleteContact(id)
	if err != nil {
		panic(err)
	}*/
}

// go test -run TestGetDeal
func TestGetDeal(t *testing.T) {
	d, err := GetDeal("477339000004063513")
	if err != nil {
		panic(err)
	}
	if d.Stage != "Prospect" {
		panic(d)
	}
}

/*
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
*/

// go test -run TestUpdateStage
func TestUpdateStage(t *testing.T) {
	id, err := FindContact("fordeal@test.com")
	if err != nil {
		panic(err)
	}
	did, err := FindDealByContactID(id)
	if did == "" {
		panic("deal not found")
	}
	if err != nil {
		panic(err)
	}
	err = UpdateDealStage(did, "Meeting")
	if err != nil {
		panic(err)
	}
}
