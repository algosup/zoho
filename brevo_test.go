package zoho

import (
	"testing"
)

// go test -v -run TestCreateBrevoContact
func TestCreateBrevoContact(t *testing.T) {
	c := BrevoContact{
		Email:         "testbrevo@test.com",
		UpdateEnabled: true,
		Attributes: Attributes{
			LastName:   "Doe",
			FirstName:  "John",
			Pipeline:   "Sales",
			Type:       "Lead",
			LeadSource: "Referral",
			MaxScore:   100,
			SMS:        "+33630901690",
			WhatsApp:   "+33630901690",
		},
		ListIds: []int{17},
	}
	err := CreateBrevoContact(c)
	if err != nil {
		panic(err)
	}
}
