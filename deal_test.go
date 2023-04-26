package zoho

import (
	"fmt"
	"log"
	"os"
	"testing"
)

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

// go test -run TestGetStageHistory
func TestGetStageHistory(t *testing.T) {
	h, err := GetDealStageHistory("477339000002064028")
	if err != nil {
		panic(err)
	}
	log.Println(h)
}

// go test -run TestGetAllDealIds
func TestGetAllDealIds(t *testing.T) {
	ids, err := GetDealIdsFromPipeline("2023-2024")
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("History 2023-2024.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for _, id := range ids {
		h, err := GetDealStageHistory(id)
		if err != nil {
			panic(err)
		}
		fmt.Println(h)
		for _, ds := range h {
			if _, err := f.WriteString(fmt.Sprintf("%s;%s;%s\r\n", id, ds.Stage, ds.ModifiedTime)); err != nil {
				panic(err)
			}
		}
	}
}
