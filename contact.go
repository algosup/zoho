package zoho

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type autoContact struct {
	ID                    string  `json:"id,omitempty"`
	Language              string  `json:"Language,omitempty"`
	Phone                 *string `json:"Phone"`
	OtherPhone            *string `json:"Other_Phone"`
	Pipeline              string  `json:"Pipeline"`
	Stage                 string  `json:"Stage"`
	GameTooShortEmailSent *bool   `json:"GameTooShortEmailSent,omitempty"` // Pointer to support omitempty
	GameFinalEmailSent    *bool   `json:"GameFinalEmailSent,omitempty"`    // Pointer to support omitempty
	LastUpdate            *Time   `json:"Last_Update,omitempty"`           // Pointer to support omitempty

	EmailsReceived    *int  `json:"Emails_Received,omitempty"`
	EmailsSent        *int  `json:"Emails_Sent,omitempty"`
	NotesCount        *int  `json:"Notes_Count,omitempty"`
	LastEmailSent     *Time `json:"Last_Email_Sent,omitempty"`     // Pointer to support omitempty
	LastEmailReceived *Time `json:"Last_Email_Received,omitempty"` // Pointer to support omitempty
	LastNote          *Time `json:"Last_Note,omitempty"`           // Pointer to support omitempty
}

type GetContactItem struct {
	Language              string `json:"Language"`
	GameTooShortEmailSent bool   `json:"GameTooShortEmailSent"`
	GameFinalEmailSent    bool   `json:"GameFinalEmailSent"`

	Origin interface{} `json:"Origin"`
	Owner  struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"Owner"`
	Email             string      `json:"Email"`
	CurrencySymbol    string      `json:"$currency_symbol"`
	FieldStates       interface{} `json:"$field_states"`
	OtherPhone        string      `json:"Other_Phone"`
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
	GameStart     time.Time   `json:"Game_Start"`
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
	GameDurationMin      int           `json:"Game_Duration_Min"`
	Phone                string        `json:"Phone"`
	StudyLevel           string        `json:"Study_level"`
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
	Pipeline             string        `json:"Pipeline"`
	Stage                string        `json:"Stage"`
	ProgramingExperience interface{}   `json:"Programing_Experience"`
	Type                 string        `json:"Type"`
	S                    interface{}   `json:"s"`
	LastName             string        `json:"Last_Name"`
	InMerge              bool          `json:"$in_merge"`
	LeadSource           string        `json:"Lead_Source"`
	Tag                  []interface{} `json:"Tag"`
	ApprovalState        string        `json:"$approval_state"`
	Pathfinder           interface{}   `json:"$pathfinder"`
	LastEnrichedTimeS    interface{}   `json:"Last_Enriched_Time__s"`
}

type Contact struct {
	ID              string `json:"id,omitempty"`
	Language        string `json:"Language,omitempty"`
	Type            string `json:"Type,omitempty"`
	Salutation      string `json:"Salutation,omitempty"`
	FirstName       string `json:"First_Name,omitempty"`
	LastName        string `json:"Last_Name,omitempty"`
	Email           string `json:"Email,omitempty"`
	Phone           string `json:"Phone,omitempty"`
	OtherPhone      string `json:"Other_Phone,omitempty"`
	StudyLevel      string `json:"Study_level,omitempty"`
	LeadSource      string `json:"Lead_Source,omitempty"`
	GameStart       *Time  `json:"Game_Start,omitempty"` // Pointer to support omitempty
	GameScore30     int    `json:"Game_Score_30,omitempty"`
	GameMaxScore    int    `json:"Game_Max_Score,omitempty"`
	GameDurationMin int    `json:"Game_Duration_Min,omitempty"`
	OpenHouse       *Date  `json:"Open_House,omitempty"`  // Pointer to support omitempty
	CodingCamp      *Date  `json:"Coding_Camp,omitempty"` // Pointer to support omitempty
	Immersion       *Date  `json:"Immersion,omitempty"`   // Pointer to support omitempty

	GameTooShortEmailSent *bool `json:"GameTooShortEmailSent,omitempty"` // Pointer to support omitempty
	GameFinalEmailSent    *bool `json:"GameFinalEmailSent,omitempty"`    // Pointer to support omitempty

	LastUpdate *Time `json:"Last_Update,omitempty"` // Pointer to support omitempty

	EmailsReceived    *int  `json:"Emails_Received,omitempty"`
	EmailsSent        *int  `json:"Emails_Sent,omitempty"`
	NotesCount        *int  `json:"Notes_Count,omitempty"`
	LastEmailSent     *Time `json:"Last_Email_Sent,omitempty"`     // Pointer to support omitempty
	LastEmailReceived *Time `json:"Last_Email_Received,omitempty"` // Pointer to support omitempty
	LastNote          *Time `json:"Last_Note,omitempty"`           // Pointer to support omitempty

	StudyingFor    string `json:"Studying_For,omitempty"`
	Pathway        string `json:"Pathway,omitempty"`
	MailingStreet  string `json:"Mailing_Street,omitempty"`
	MailingCountry string `json:"Mailing_Country,omitempty"`
	MailingZip     string `json:"Mailing_Zip,omitempty"`
	MailingCity    string `json:"Mailing_City,omitempty"`

	NSI                           bool `json:"NSI"`
	SciencesDeLIngNieur           bool `json:"Sciences_de_l_Ing_nieur"`
	MathMatiques                  bool `json:"Math_matiques"`
	PhysiqueChimie                bool `json:"Physique_Chimie"`
	Biologie                      bool `json:"Biologie"`
	SciencesEconomiquesEtSociales bool `json:"Sciences_Economiques_et_Sociales"`
	Anglais                       bool `json:"Anglais"`
	Espagnol                      bool `json:"Espagnol"`
	LittRaturePhilosophie         bool `json:"Litt_rature_Philosophie"`
	HistoireGOgraphie             bool `json:"Histoire_G_ographie"`
	SNT                           bool `json:"SNT"`
	Environnement                 bool `json:"Environnement"`
	SciencesPolitiques            bool `json:"Sciences_Politiques"`
	SIN                           bool `json:"SIN"`

	/*
		Origin                        interface{}   `json:"Origin"`
		CurrencySymbol                string        `json:"$currency_symbol"`
		FieldStates                   interface{}   `json:"$field_states"`

		MailingState                  interface{}   `json:"Mailing_State"`
		SharingPermission             string        `json:"$sharing_permission"`
		LastActivityTime              time.Time     `json:"Last_Activity_Time"`
		State                         string        `json:"$state"`
		UnsubscribedMode              interface{}   `json:"Unsubscribed_Mode"`
		ProcessFlow                   bool          `json:"$process_flow"`
		TestNewsletterOptOut          bool          `json:"test_Newsletter_Opt_Out"`
		EnrichStatusS                 interface{}   `json:"Enrich_Status__s"`
		PotentialLoan                 interface{}   `json:"Potential_loan"`
		Editable                      bool          `json:"$editable"`
		ITEC                          bool          `json:"ITEC"`
		EquipmentNumber               interface{}   `json:"Equipment_Number"`


		TestINE                       interface{}   `json:"test_INE"`
		MoonshotProjectRepository     interface{}   `json:"Moonshot_Project_Repository"`
		ZiaOwnerAssignment            interface{}   `json:"$zia_owner_assignment"`
		SecondaryEmail                interface{}   `json:"Secondary_Email"`
		VendorName                    interface{}   `json:"Vendor_Name"`
		CanvasID                      interface{}   `json:"$canvas_id"`
		School                        interface{}   `json:"School"`
		Review                        interface{}   `json:"$review"`

		AdmissionLevel                interface{}   `json:"Admission_Level"`
		EmailOptOut                   bool          `json:"Email_Opt_Out"`
		ZiaVisions                    interface{}   `json:"$zia_visions"`
		DateOfBirth                   interface{}   `json:"Date_of_Birth"`
		CommunicationMarketing        bool          `json:"Communication_Marketing"`
		UnsubscribedTime              interface{}   `json:"Unsubscribed_Time"`
		PlaceOfBirth                  interface{}   `json:"Place_of_birth"`
		JobTitle                      interface{}   `json:"Job_Title"`
		Orchestration                 interface{}   `json:"$orchestration"`
		Pipeline                      interface{}   `json:"Pipeline"`
		ProgramingExperience          interface{}   `json:"Programing_Experience"`
		SchoolFair1                   interface{}   `json:"School_Fair1"`
		GestionFinanceRH              bool          `json:"Gestion_Finance_RH"`
		InMerge                       bool          `json:"$in_merge"`
		Tag                           []interface{} `json:"Tag"`
		ApprovalState                 string        `json:"$approval_state"`
		Pathfinder                    interface{}   `json:"$pathfinder"`
		LastEnrichedTimeS             interface{}   `json:"Last_Enriched_Time__s"`
	*/
}

type contact struct {
	Data                 []Contact `json:"data"`
	DuplicateCheckFields []string  `json:"duplicate_check_fields"`
	Trigger              []string  `json:"trigger"`
}

type autoContactData struct {
	Data                 []autoContact `json:"data"`
	DuplicateCheckFields []string      `json:"duplicate_check_fields"`
	Trigger              []string      `json:"trigger"`
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
		PerPage     int    `json:"per_page"`
		NextIndex   string `json:"next_index"`
		Count       int    `json:"count"`
		PrevIndex   string `json:"prev_index"`
		MoreRecords bool   `json:"more_records"`
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
		b, err := io.ReadAll(r.Body)
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
		b, err := io.ReadAll(r.Body)
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
		b, err := io.ReadAll(r.Body)
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
		b, err := io.ReadAll(r.Body)
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
		b, err = io.ReadAll(r.Body)
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
		b, err = io.ReadAll(r.Body)
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

func updateAutoContact(item autoContact) error {
	//fmt.Println(item)
	b, err := json.Marshal(&autoContactData{Data: []autoContact{item}})
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

	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusNoContent {
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
