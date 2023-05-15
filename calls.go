package zoho

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

func GetContactCalls(id string) ([]ContactEmail, error) {
	req, err := http.NewRequest("GET", "https://www.zohoapis.eu/crm/v3/contacts/"+id+"/calls?sort_order=desc&fields=Created_Time", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNoContent {
		return []ContactEmail{}, nil
	}

	//	if r.StatusCode != http.StatusOK {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(b))
	return nil, errors.New(r.Status)
	//	}

	var res contactEmails
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return res.Emails, nil
}
