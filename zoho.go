package zoho

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// https://api-console.zoho.eu/
// ZohoCRM.modules.ALL,ZohoSearch.securesearch.READ

// curl -X POST -F grant_type=authorization_code -F client_id=1000.XXXXXXXXXX -F client_secret=XXXXX -F code=1000.XXXXXXXXXXXXXXXXXX https://accounts.zoho.eu/oauth/v2/token

// curl -X POST "https://accounts.zoho.eu/oauth/v2/token?refresh_token=1000.XXXXXXXXXXXXXX&client_id=1000.XXXXXXXXXXXXX&client_secret=XXXXXXXXXX&grant_type=refresh_token"

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

type Contact struct {
	ID              string `json:"id,omitempty"`
	Type            string `json:"Type,omitempty"`
	FirstName       string `json:"First_Name,omitempty"`
	LastName        string `json:"Last_Name,omitempty"`
	Email           string `json:"Email,omitempty"`
	Phone           string `json:"Phone,omitempty"`
	StudyLevel      string `json:"Study_level,omitempty"`
	LeadSource      string `json:"Lead_Source,omitempty"`
	GameStart       *Time  `json:"Game_Start,omitempty"` // Pointer to support omitempty
	GameScore30     int    `json:"Game_Score_30,omitempty"`
	GameMaxScore    int    `json:"Game_Max_Score,omitempty"`
	GameDurationMin int    `json:"Game_Duration_Min,omitempty"`
	OpenHouse       *Date  `json:"Open_House,omitempty"`  // Pointer to support omitempty
	CodingCamp      *Date  `json:"Coding_Camp,omitempty"` // Pointer to support omitempty
	Immersion       *Date  `json:"Immersion,omitempty"`   // Pointer to support omitempty
}

type contact struct {
	Data                 []Contact `json:"data"`
	DuplicateCheckFields []string  `json:"duplicate_check_fields"`
	Trigger              []string  `json:"trigger"`
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

type GetContactItem struct {
	Origin interface{} `json:"Origin"`
	Owner  struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"Owner"`
	Email             string      `json:"Email"`
	CurrencySymbol    string      `json:"$currency_symbol"`
	FieldStates       interface{} `json:"$field_states"`
	OtherPhone        interface{} `json:"Other_Phone"`
	MailingState      interface{} `json:"Mailing_State"`
	SharingPermission string      `json:"$sharing_permission"`
	Immersion         interface{} `json:"Immersion"`
	LastActivityTime  time.Time   `json:"Last_Activity_Time"`
	State             string      `json:"$state"`
	UnsubscribedMode  interface{} `json:"Unsubscribed_Mode"`
	ProcessFlow       bool        `json:"$process_flow"`
	GameScore30       int         `json:"Game_Score_30"`
	MailingCountry    interface{} `json:"Mailing_Country"`
	ID                string      `json:"id"`
	Approval          struct {
		Delegate bool `json:"delegate"`
		Approve  bool `json:"approve"`
		Reject   bool `json:"reject"`
		Resubmit bool `json:"resubmit"`
	} `json:"$approval"`
	EnrichStatusS interface{} `json:"Enrich_Status__s"`
	CreatedTime   time.Time   `json:"Created_Time"`
	PotentialLoan interface{} `json:"Potential_loan"`
	Editable      bool        `json:"$editable"`
	GameStart     interface{} `json:"Game_Start"`
	TestINE       interface{} `json:"test_INE"`
	CreatedBy     struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"Created_By"`
	SecondaryEmail interface{} `json:"Secondary_Email"`
	GameMaxScore   int         `json:"Game_Max_Score"`
	Description    interface{} `json:"Description"`
	MailingZip     interface{} `json:"Mailing_Zip"`
	VendorName     interface{} `json:"Vendor_Name"`
	ReviewProcess  struct {
		Approve  bool `json:"approve"`
		Reject   bool `json:"reject"`
		Resubmit bool `json:"resubmit"`
	} `json:"$review_process"`
	MailingStreet        interface{}   `json:"Mailing_Street"`
	CanvasID             interface{}   `json:"$canvas_id"`
	Salutation           interface{}   `json:"Salutation"`
	OpenHouse            interface{}   `json:"Open_House"`
	FirstName            string        `json:"First_Name"`
	FullName             string        `json:"Full_Name"`
	School               interface{}   `json:"School"`
	Review               interface{}   `json:"$review"`
	GameDurationMin      interface{}   `json:"Game_Duration_Min"`
	Phone                string        `json:"Phone"`
	StudyLevel           interface{}   `json:"Study_level"`
	TestScoringCampaigns interface{}   `json:"test_Scoring_Campaigns"`
	AccountName          interface{}   `json:"Account_Name"`
	AdmissionLevel       interface{}   `json:"Admission_Level"`
	EmailOptOut          bool          `json:"Email_Opt_Out"`
	ZiaVisions           interface{}   `json:"$zia_visions"`
	CodingCamp           interface{}   `json:"Coding_Camp"`
	DateOfBirth          interface{}   `json:"Date_of_Birth"`
	MailingCity          interface{}   `json:"Mailing_City"`
	UnsubscribedTime     interface{}   `json:"Unsubscribed_Time"`
	PlaceOfBirth         interface{}   `json:"Place_of_birth"`
	JobTitle             interface{}   `json:"Job_Title"`
	Orchestration        interface{}   `json:"$orchestration"`
	Pipeline             interface{}   `json:"Pipeline"`
	ProgramingExperience interface{}   `json:"Programing_Experience"`
	Type                 interface{}   `json:"Type"`
	S                    interface{}   `json:"s"`
	LastName             string        `json:"Last_Name"`
	InMerge              bool          `json:"$in_merge"`
	LeadSource           string        `json:"Lead_Source"`
	Tag                  []interface{} `json:"Tag"`
	ApprovalState        string        `json:"$approval_state"`
	Pathfinder           interface{}   `json:"$pathfinder"`
	LastEnrichedTimeS    interface{}   `json:"Last_Enriched_Time__s"`
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
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			time.Sleep(time.Duration(auth.ExpiresIn-60) * time.Second)
			err := authenticate()
			if err != nil {
				log.Fatalln(err)
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
	r, err := http.Post("https://accounts.zoho.eu/oauth/v2/token?refresh_token="+os.Getenv("ZOHO_REFRESH_TOKEN")+"&client_id="+os.Getenv("ZOHO_CLIENT_ID")+"&client_secret="+os.Getenv("ZOHO_CLIENT_SECRET")+"&grant_type=refresh_token", "", nil)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		log.Fatalln(string(b))
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

// Scope:  ZohoCRM.modules.all and ZohoSearch.securesearch.READ

func FindContact(email string) (string, error) {
	req, err := http.NewRequest("GET", "https://www.zohoapis.eu/crm/v3/Contacts/search?email="+email, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNoContent {
		return "", nil
	}
	if r.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		log.Fatalln(string(b))
	}
	var res getContactResponse
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return "", err
	}

	return res.Data[0].ID, nil
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
		b, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		log.Fatalln(string(b))
	}
	var res createNote
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return "", err
	}
	return res.Data[0].Details.ID, nil
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
		log.Fatalln(string(b))
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
		log.Fatalln(string(b))
	}

	var res updateContactResult
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return err
	}
	return nil
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
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		log.Fatalln(string(b))
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
		log.Fatalln(string(b))
	}

	var res getContactResponse
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res.Data[0], nil
}
