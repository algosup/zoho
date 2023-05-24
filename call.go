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

type getCalls struct {
	Data []ContactCall `json:"data"`
	Info struct {
		PerPage           int         `json:"per_page"`
		NextPageToken     interface{} `json:"next_page_token"`
		Count             int         `json:"count"`
		Page              int         `json:"page"`
		PreviousPageToken interface{} `json:"previous_page_token"`
		PageTokenExpiry   interface{} `json:"page_token_expiry"`
		MoreRecords       bool        `json:"more_records"`
	} `json:"info"`
}

type ContactCall struct {
	CallType   string `json:"Call_Type"`
	CallResult string `json:"Call_Result"`
	Owner      struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"Owner"`
	CallStartTime         Time   `json:"Call_Start_Time"`
	CallDurationInSeconds int    `json:"Call_Duration_in_seconds"`
	ID                    string `json:"id"`
	Subject               string `json:"Subject"`
}

func GetOpenCalls(contactID string) ([]ContactCall, error) {
	req, err := http.NewRequest("GET", "https://www.zohoapis.eu/crm/v4/contacts/"+contactID+"/calls?sort_order=desc&fields=Subject,Call_Type,Call_Start_Time,Owner,Call_Result,Call_Status,Call_Duration_in_seconds", nil)
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
		return []ContactCall{}, nil
	}

	if r.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		log.Println(string(b))
		return nil, errors.New(r.Status)
	}

	var res getCalls
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

type calls struct {
	Data []NewCall `json:"data"`
}

type Owner struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Email string `json:"email"`
}
type NameID struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
type NewCall struct {
	Subject       string `json:"Subject"`
	CallType      string `json:"Call_Type"`
	CallStartTime Time   `json:"Call_Start_Time"`
	CallDuration  string `json:"Call_Duration"`

	Owner    Owner  `json:"Owner"`
	SeModule string `json:"$se_module"`
	WhoID    NameID `json:"Who_Id"` // Contact
	/*
		Description string      `json:"Description"`
		CallerID    interface{} `json:"Caller_ID"`
		CallAgenda  string      `json:"Call_Agenda"`
		CallPurpose string      `json:"Call_Purpose"`
		CallStatus  string      `json:"Call_Status"`

		Reminder    interface{} `json:"Reminder"`
		CallResult  string      `json:"Call_Result"`
		//WhatID                NameID        `json:"What_Id"`
		CallDurationInSeconds string        `json:"Call_Duration_in_seconds"`
		Tag                   []interface{} `json:"Tag"`
		DialledNumber         interface{}   `json:"Dialled_Number"`*/
}

func CreateCall(call NewCall, contactId string) (string, error) {
	call.Owner = Owner{ID: "477339000000309001", Email: "julie.dupoty@algosup.com"} // Julie
	call.WhoID.ID = contactId
	call.SeModule = "Activities"

	b, err := json.Marshal(&calls{Data: []NewCall{call}})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://www.zohoapis.eu/crm/v4/Calls", bytes.NewReader(b))
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

func GetCalls(contactID string) ([]ContactCall, error) {
	//"SELECT id,Who_Id FROM Calls WHERE Owner.id='477339000000309001'"
	b, err := json.Marshal(&selectQuery{Query: fmt.Sprintf("SELECT Subject,Call_Type,Call_Start_Time,Owner,Call_Result,Call_Duration_in_seconds FROM Calls WHERE Who_Id.id='%s'", contactID)})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://www.zohoapis.eu/crm/v4/coql", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	req.Header.Set("Content-Type", "text/plain")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	if r.StatusCode != http.StatusOK {
		b, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		log.Println(string(b))
		return nil, errors.New(r.Status)
	}
	var res getCalls
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}
