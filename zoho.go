package zoho

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Auth struct {
	AccessToken string `json:"access_token"`
	APIDomain   string `json:"api_domain"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`

	Code    string `json:"code"`
	Details struct {
	} `json:"details"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type notes struct {
	Data []note `json:"data"`
}

type note struct {
	Title   string `json:"Note_Title,omitempty"`
	Content string `json:"Note_Content,omitempty"`
}

type updateContact struct {
	Code    string `json:"code"`
	Details struct {
		ModifiedTime time.Time `json:"Modified_Time"`
		ModifiedBy   struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"Modified_By"`
		CreatedTime time.Time `json:"Created_Time"`
		ID          string    `json:"id"`
		CreatedBy   struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"Created_By"`
		ApprovalState string `json:"$approval_state"`
	} `json:"details"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
type updateContactResult struct {
	Data []updateContact `json:"data"`
}

type getContactResponse struct {
	Data []GetContactItem `json:"data"`
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

type findContactResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

type modifiedContactResponse struct {
	Data []struct {
		ID       string `json:"id"`
		Modified string `json:"Modified_Time"`
	} `json:"data"`
}

type createNote struct {
	Data []struct {
		Message string `json:"message"`
		Details struct {
			CreatedBy struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"created_by"`
			ID         string `json:"id"`
			ModifiedBy struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"modified_by"`
			ModifiedTime time.Time `json:"modified_time"`
			CreatedTime  time.Time `json:"created_time"`
		} `json:"details"`
		Status string `json:"status"`
		Code   string `json:"code"`
	} `json:"data"`
}

var auth Auth

func init() {
	var err error
	for i := 0; i < 10; i++ {
		err := authenticate()
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	go func() {
		for {
			time.Sleep(time.Duration(auth.ExpiresIn-60) * time.Second)
			for i := 0; i < 10; i++ {
				err := authenticate()
				if err == nil {
					break
				}
				log.Println(err)
				time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			}
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

type Date struct {
	time.Time
}

type Time struct {
	time.Time
}

func (t Date) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))
	return []byte(stamp), nil
}

func (t *Date) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.Format("2006-01-02T15:04:05Z07:00"))
	return []byte(stamp), nil
}
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05Z07:00"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func authenticate() error {
	url := "https://accounts.zoho.eu/oauth/v2/token?refresh_token=" + os.Getenv("ZOHO_REFRESH_TOKEN") + "&client_id=" + os.Getenv("ZOHO_CLIENT_ID") + "&client_secret=" + os.Getenv("ZOHO_CLIENT_SECRET") + "&grant_type=refresh_token"
	r, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		log.Println(string(b))
		return errors.New(r.Status)
	}
	err = json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		return err
	}

	if auth.Status == "error" {
		return errors.New(auth.Message)
	}
	return nil
}

type selectQuery struct {
	Query string `json:"select_query"`
}

func AddContactNote(id string, title string, content string) (string, error) {

	b, err := json.Marshal(notes{Data: []note{{Title: title, Content: content}}})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://www.zohoapis.eu/crm/v3/Contacts/"+id+"/Notes", bytes.NewReader(b))
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
		b, err = io.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		log.Println(string(b))
		return "", errors.New(r.Status)
	}
	var res createNote
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return "", err
	}
	return res.Data[0].Details.ID, nil
}

func clearContactField(id string, field string) error {
	json := fmt.Sprintf(`{"data":[{"id":"%s","%s":null}]}`, id, field)
	req, err := http.NewRequest("PUT", "https://www.zohoapis.eu/crm/v3/Contacts", strings.NewReader(json))
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
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		log.Println(string(b))
		return errors.New(r.Status)
	}

	return nil
}

func CancelOpenHouse(contactId string) error {
	return clearContactField(contactId, "Open_House")
}

func CancelImmersion(contactId string) error {
	return clearContactField(contactId, "Immersion")
}

func CancelCodingCamp(contactId string) error {
	return clearContactField(contactId, "Coding_Camp")
}
