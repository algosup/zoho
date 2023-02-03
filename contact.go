package zoho

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Contact struct {
	ID                string `json:"id,omitempty"`
	Type              string `json:"Type,omitempty"`
	FirstName         string `json:"First_Name,omitempty"`
	LastName          string `json:"Last_Name,omitempty"`
	Email             string `json:"Email,omitempty"`
	Phone             string `json:"Phone,omitempty"`
	StudyLevel        string `json:"Study_level,omitempty"`
	LeadSource        string `json:"Lead_Source,omitempty"`
	GameStart         *Time  `json:"Game_Start,omitempty"` // Pointer to support omitempty
	GameScore30       int    `json:"Game_Score_30,omitempty"`
	GameMaxScore      int    `json:"Game_Max_Score,omitempty"`
	GameDurationMin   int    `json:"Game_Duration_Min,omitempty"`
	OpenHouse         *Date  `json:"Open_House,omitempty"`  // Pointer to support omitempty
	CodingCamp        *Date  `json:"Coding_Camp,omitempty"` // Pointer to support omitempty
	Immersion         *Date  `json:"Immersion,omitempty"`   // Pointer to support omitempty
	EmailsReceived    *int   `json:"Emails_Received,omitempty"`
	EmailsSent        *int   `json:"Emails_Sent,omitempty"`
	NotesCount        *int   `json:"Notes_Count,omitempty"`
	LastEmailSent     *Time  `json:"Last_Email_Sent,omitempty"`     // Pointer to support omitempty
	LastEmailReceived *Time  `json:"Last_Email_Received,omitempty"` // Pointer to support omitempty
	LastNote          *Time  `json:"Last_Note,omitempty"`           // Pointer to support omitempty
	LastUpdate        *Time  `json:"Last_Update,omitempty"`         // Pointer to support omitempty
}

type contact struct {
	Data                 []Contact `json:"data"`
	DuplicateCheckFields []string  `json:"duplicate_check_fields"`
	Trigger              []string  `json:"trigger"`
}

type ContactEmail struct {
	Cc                  interface{} `json:"cc"`
	HasThreadAttachment bool        `json:"has_thread_attachment"`
	Summary             interface{} `json:"summary"`
	Owner               struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"owner"`
	Read          bool        `json:"read"`
	Subject       string      `json:"subject"`
	MessageID     string      `json:"message_id"`
	HasAttachment bool        `json:"has_attachment"`
	Source        string      `json:"source"`
	Sent          bool        `json:"sent"`
	Intent        interface{} `json:"intent"`
	SentimentInfo interface{} `json:"sentiment_info"`
	LinkedRecord  interface{} `json:"linked_record"`
	Emotion       interface{} `json:"emotion"`
	From          struct {
		UserName string `json:"user_name"`
		Email    string `json:"email"`
	} `json:"from"`
	To []struct {
		UserName string `json:"user_name"`
		Email    string `json:"email"`
	} `json:"to"`
	Time   time.Time `json:"time"`
	Status []struct {
		Type string `json:"type"`
	} `json:"status"`
}

type contactEmails struct {
	Emails []ContactEmail `json:"Emails"`
	Info   struct {
		PerPage     int  `json:"per_page"`
		NextIndex   int  `json:"next_index"`
		Count       int  `json:"count"`
		PrevIndex   int  `json:"prev_index"`
		MoreRecords bool `json:"more_records"`
	} `json:"info"`
}

type contactNotes struct {
	Data []struct {
		CreatedTime time.Time `json:"Created_Time"`
		ID          string    `json:"id"`
	} `json:"data"`
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

func GetContact(id string) (*GetContactItem, error) {
	req, err := http.NewRequest("GET", "https://www.zohoapis.eu/crm/v3/Contacts/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		log.Println(string(b))
		return nil, errors.New(r.Status)
	}

	var res getContactResponse
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res.Data[0], nil
}

func DeleteContact(id string) error {
	req, err := http.NewRequest("DELETE", "https://www.zohoapis.eu/crm/v3/Contacts/"+id, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		log.Println(string(b))
		return errors.New(r.Status)
	}

	return nil
}

func GetContactEmails(id string) ([]ContactEmail, error) {
	req, err := http.NewRequest("GET", "https://www.zohoapis.eu/crm/v3/contacts/"+id+"/emails", nil)
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

	if r.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		log.Println(string(b))
		return nil, errors.New(r.Status)
	}

	var res contactEmails
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return res.Emails, nil
}

func GetContactNotesCount(id string) (int, *time.Time, error) {
	req, err := http.NewRequest("GET", "https://www.zohoapis.eu/crm/v3/contacts/"+id+"/notes?sort_order=desc&fields=Created_Time", nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNoContent {
		return 0, nil, nil
	}

	if r.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return 0, nil, err
		}
		log.Println(string(b))
		return 0, nil, errors.New(r.Status)
	}

	var res contactNotes
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return 0, nil, err
	}

	var last *time.Time
	if len(res.Data) > 0 {
		last = &res.Data[0].CreatedTime
	}
	return res.Info.Count, last, nil
}

func CreateContact(item Contact) (string, error) {
	if item.LastName == "" {
		item.LastName = "MISSING"
	}

	b, err := json.Marshal(&contact{Data: []Contact{item}})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://www.zohoapis.eu/crm/v3/Contacts", bytes.NewReader(b))
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

func UpdateContact(item Contact) error {
	b, err := json.Marshal(&contact{Data: []Contact{item}})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", "https://www.zohoapis.eu/crm/v3/Contacts", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		b, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		log.Println(string(b))
		return errors.New(r.Status)
	}

	var res updateContactResult
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return err
	}
	return nil
}

// scope=ZohoCRM.coql.READ (and) scope=ZohoCRM.modules.all

func FindContact(email string) (string, error) {

	b, err := json.Marshal(&selectQuery{Query: fmt.Sprintf("SELECT id FROM Contacts WHERE Email='%s'", email)})
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
