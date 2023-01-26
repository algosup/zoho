package zoho

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func SendMail(fromName, fromEmail, toName, toEmail, contactID, templateID string) error {
	b, err := json.Marshal(&mail{
		Data: []mailData{{From: userEmail{UserName: fromName, Email: fromEmail},
			To:       []userEmail{{UserName: toName, Email: toEmail}},
			Template: &template{ID: templateID},
			OrgEmail: true,
		}}})
	if err != nil {
		return err
	}

	//log.Println((string)(b))
	req, err := http.NewRequest("POST", "https://www.zohoapis.eu/crm/v3/contacts/"+contactID+"/actions/send_mail", bytes.NewReader(b))
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

	return nil
}

func CheckValidFrom() {
	req, err := http.NewRequest("GET", "https://www.zohoapis.eu/crm/v3/settings/emails/actions/from_addresses", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Zoho-oauthtoken "+auth.AccessToken)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(b))
}

type userEmail struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type template struct {
	ID string `json:"id"`
}
type mailData struct {
	From          userEmail    `json:"from"`
	To            []userEmail  `json:"to"`
	Cc            []*userEmail `json:"cc,omitempty"`
	Bcc           []*userEmail `json:"bcc,omitempty"`
	ReplyTo       *userEmail   `json:"reply_to,omitempty"`
	OrgEmail      bool         `json:"org_email"`
	InReplyTo     string       `json:"in_reply_to,omitempty"`
	ScheduledTime *time.Time   `json:"scheduled_time,omitempty"`
	Subject       string       `json:"subject,omitempty"`
	Content       string       `json:"content,omitempty"`
	MailFormat    string       `json:"mail_format,omitempty"`
	Attachments   []struct {
		ID string `json:"id"`
	} `json:"attachments,omitempty"`
	Template *template `json:"template,omitempty"`
}
type mail struct {
	Data []mailData `json:"data"`
}
