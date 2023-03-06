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

type Deal struct {
	ID         string `json:"id,omitempty"`
	ContactId  string `json:"Contact_Name,omitempty"`
	DealName   string `json:"Deal_Name,omitempty"`
	Stage      string `json:"Stage,omitempty"`
	Pipeline   string `json:"Pipeline,omitempty"`
	LeadSource string `json:"Lead_Source,omitempty"`
}

type deal struct {
	Data                 []Deal   `json:"data"`
	DuplicateCheckFields []string `json:"duplicate_check_fields"`
	Trigger              []string `json:"trigger"`
}

type DealGet struct {
	Owner struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"Owner"`
	Description    interface{} `json:"Description"`
	CurrencySymbol string      `json:"$currency_symbol"`
	FieldStates    interface{} `json:"$field_states"`
	ReviewProcess  struct {
		Approve  bool `json:"approve"`
		Reject   bool `json:"reject"`
		Resubmit bool `json:"resubmit"`
	} `json:"$review_process"`
	Followers            interface{} `json:"$followers"`
	SharingPermission    string      `json:"$sharing_permission"`
	CanvasID             interface{} `json:"$canvas_id"`
	ClosingDate          interface{} `json:"Closing_Date"`
	LastActivityTime     time.Time   `json:"Last_Activity_Time"`
	Review               interface{} `json:"$review"`
	LeadConversionTime   interface{} `json:"Lead_Conversion_Time"`
	State                string      `json:"$state"`
	ProcessFlow          bool        `json:"$process_flow"`
	DealName             string      `json:"Deal_Name"`
	OverallSalesDuration interface{} `json:"Overall_Sales_Duration"`
	Stage                string      `json:"Stage"`
	AccountName          interface{} `json:"Account_Name"`
	ID                   string      `json:"id"`
	AdmissionLevel       interface{} `json:"Admission_Level"`
	ZiaVisions           interface{} `json:"$zia_visions"`
	Approval             struct {
		Delegate bool `json:"delegate"`
		Approve  bool `json:"approve"`
		Reject   bool `json:"reject"`
		Resubmit bool `json:"resubmit"`
	} `json:"$approval"`
	LeadName      interface{} `json:"Lead_Name"`
	Amount        int         `json:"Amount"`
	Followed      bool        `json:"$followed"`
	Probability   int         `json:"Probability"`
	NextStep      interface{} `json:"Next_Step"`
	Editable      bool        `json:"$editable"`
	Orchestration interface{} `json:"$orchestration"`
	Pipeline      string      `json:"Pipeline"`
	ContactName   struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"Contact_Name"`
	SalesCycleDuration interface{}   `json:"Sales_Cycle_Duration"`
	InMerge            bool          `json:"$in_merge"`
	LeadSource         string        `json:"Lead_Source"`
	Tag                []interface{} `json:"Tag"`
	ZiaOwnerAssignment interface{}   `json:"$zia_owner_assignment"`
	ApprovalState      string        `json:"$approval_state"`
	Pathfinder         interface{}   `json:"$pathfinder"`
}
type dealGet struct {
	Data []DealGet `json:"data"`
}

func GetDeal(id string) (*DealGet, error) {
	req, err := http.NewRequest("GET", "https://www.zohoapis.eu/crm/v3/Deals/"+id, nil)
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

	var res dealGet
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res.Data[0], nil
}

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

func UpdateDealStage(dealID string, stage string) error {
	b, err := json.Marshal(&deal{Data: []Deal{{ID: dealID, Stage: stage}}})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", "https://www.zohoapis.eu/crm/v3/Deals", bytes.NewReader(b))
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

func UpdateDealLeadSource(dealID string, source string) error {
	b, err := json.Marshal(&deal{Data: []Deal{{ID: dealID, LeadSource: source}}})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", "https://www.zohoapis.eu/crm/v3/Deals", bytes.NewReader(b))
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
