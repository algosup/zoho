package zoho

import (
	"testing"
)

func TestSendMail(t *testing.T) {
	id, err := FindContact("email@franckjeannin.com")
	if id == "" {
		panic("not found")
	}
	if err != nil {
		panic(err)
	}

	//CheckValidFrom()

	err = SendMail("Natacha BOEZ", "natacha.boez@algosup.com", "Bob", "email@franckjeannin.com", id, "477339000002907121")
	if err != nil {
		panic(err)
	}

}
