package zoho

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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
