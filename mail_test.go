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

	CheckValidFrom()

	//https://crm.zoho.eu/crm/org20082710153/settings/templates?type=email&templateId=477339000002907121
	err = SendMail("Franck JEANNIN", "franck.jeannin@algosup.com", "Franck JEANNIN", "email@franckjeannin.com", id, "477339000002907121")
	if err != nil {
		panic(err)
	}

}
