package zoho

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func CreateDeal(item Deal) (string, error) {
	b, err := json.Marshal(&deal{Data: []Deal{item}})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://www.zohoapis.eu/crm/v3/Deals", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusCreated {
		b, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		log.Println(string(b))
		return "", errors.New(r.Status)
	}
	var res updateContactResult
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return "", err
	}
	return res.Data[0].Details.ID, nil
}

func FindDealByContactID(id string) (string, error) {
	b, err := json.Marshal(&selectQuery{Query: fmt.Sprintf("SELECT id FROM Deals WHERE Contact_Name='%s'", id)})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://www.zohoapis.eu/crm/v3/coql", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	req.Header.Set("Content-Type", "text/plain")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNoContent {
		return "", nil
	}

	if r.StatusCode != http.StatusOK {
		b, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		log.Println(string(b))
		return "", errors.New(r.Status)
	}
	var res findContactResponse
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return "", err
	}
	return res.Data[0].ID, nil
}
