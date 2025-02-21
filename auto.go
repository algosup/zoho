package zoho

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func normalizePhone(phone, otherPhone string) (string, string) {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, ".", "")
	phone = strings.ReplaceAll(phone, "-", "")
	otherPhone = strings.ReplaceAll(otherPhone, " ", "")
	otherPhone = strings.ReplaceAll(otherPhone, ".", "")
	otherPhone = strings.ReplaceAll(otherPhone, "-", "")

	if phone == "" && len(otherPhone) == 11 && otherPhone[0:2] == "33" {
		return "+" + otherPhone, ""
	}

	if phone == "" && len(otherPhone) == 13 && otherPhone[0:3] == "+33" && (otherPhone[3:5] == "06" || otherPhone[3:5] == "07") {
		return otherPhone[0:3] + otherPhone[4:], ""
	}

	if phone == "" {
		return phone, otherPhone
	}

	if phone == otherPhone {
		return phone, ""
	}

	if phone == "+"+otherPhone {
		return phone, ""
	}

	if len(phone) < 9 {
		return "", phone
	}

	if len(phone) == 13 && phone[0:3] == "+33" && (phone[3:5] == "06" || phone[3:5] == "07") {
		return phone[0:3] + phone[4:], otherPhone
	}

	if phone[0:3] == "+33" {
		return phone, otherPhone
	}
	if phone[0:2] == "33" && len(phone) == 11 {
		return "+" + phone, otherPhone
	}

	if len(phone) == 10 && (phone[0:2] == "06" || phone[0:2] == "07") {
		return "+33" + phone[1:], otherPhone
	}

	if len(phone) == 9 && (phone[0] == '6' || phone[0] == '7') {
		return "+33" + phone, otherPhone
	}

	if len(phone) > 12 {
		return "", phone
	}
	if phone[0] == '+' {
		return "", phone // Move foreign country to other phone
	}

	return "", phone
}

func AutoUpdateContact(id string) error {
	c, err := GetContact(id)
	if err != nil {
		return err
	}
	original := autoContact{
		ID:                 c.ID,
		Language:           c.Language,
		Phone:              c.Phone,
		OtherPhone:         c.OtherPhone,
		Pipeline:           c.Pipeline,
		Stage:              c.Stage,
		GameFinalEmailSent: c.GameFinalEmailSent,
	}

	lang := c.Language
	if lang == "" {
		lang = "fr-FR"
	}

	did, err := FindDealByContactID(id)
	if err != nil {
		return err
	}
	c.Pipeline = ""
	if did == "" {
		switch c.Type {
		case "Prospect":
			switch c.StudyLevel {
			case "Seconde":
				c.Pipeline = "2027-2028"
			case "Première":
				c.Pipeline = "2026-2027"
			case "Terminale", "BAC", "BAC+1", "BAC+2", "BAC+3":
				c.Pipeline = "2025-2026"
			}
		case "Relative", "Student Relative":
			switch c.StudyLevel {
			case "Seconde":
				c.Pipeline = "Relative 2027-2028"
			case "Première":
				c.Pipeline = "Relative 2026-2027"
			case "Terminale", "BAC", "BAC+1", "BAC+2", "BAC+3":
				c.Pipeline = "Relative 2025-2026"
			}
		}

		stage := "Prospect"
		if c.CodingCamp != nil || c.OpenHouse != nil {
			stage = "Meeting scheduled"
		}
		if c.Pipeline != "" {
			_, err := CreateDeal(Deal{
				ContactId:  id,
				DealName:   c.FirstName + " " + c.LastName,
				Stage:      stage,
				Pipeline:   c.Pipeline,
				LeadSource: c.LeadSource,
				Amount:     9500,
			})
			if err != nil {
				return err
			}
		}
	} else {
		deal, err := GetDeal(did)
		if err != nil {
			return err
		}
		c.Pipeline = deal.Pipeline
		c.Stage = deal.Stage
	}

	phone, otherPhone := normalizePhone(c.Phone, c.OtherPhone)

	a := original
	a.Phone = phone
	a.OtherPhone = otherPhone
	template := templateToSend(c)
	a.GameFinalEmailSent = c.GameFinalEmailSent || template != ""
	a.Pipeline = c.Pipeline
	a.Stage = c.Stage
	a.Language = lang

	if !reflect.DeepEqual(original, a) {
		//log.Println(a, original)
		err = updateAutoContact(a)
		if err != nil {
			log.Println(a)
			return err
		}
	}

	if template != "" {
		log.Println("Sending", template, "to", c.Email)
		err = SendMail("Natacha BOEZ", "natacha.boez@algosup.com", c.FirstName+" "+c.LastName, c.Email, id, template)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	brevo := BrevoContact{
		Email:         c.Email,
		UpdateEnabled: true,
		Attributes: Attributes{
			LastName:   c.LastName,
			FirstName:  c.FirstName,
			Pipeline:   c.Pipeline,
			Type:       c.Type,
			LeadSource: c.LeadSource,
			MaxScore:   c.GameMaxScore,
			SMS:        c.Phone,
			WhatsApp:   c.Phone,
		},
		ListIds: []int{17},
	}
	err = CreateBrevoContact(brevo)
	if err != nil {
		log.Println(err)
		log.Println(brevo)

		brevo.Attributes.SMS = ""
		brevo.Attributes.WhatsApp = ""
		return CreateBrevoContact(brevo) // Try again without phone number
	}
	return nil
}

func templateToSend(contact *GetContactItem) string {
	if contact.GameStart.IsZero() {
		return ""
	}

	if contact.GameStart.Before(time.Date(2023, 5, 19, 15, 0, 0, 0, time.UTC)) {
		// Game started before new messages where in place

		return ""
	}

	if contact.GameFinalEmailSent {
		return ""
	}

	if contact.GameStart.After(time.Now().Add(-time.Hour)) {
		// Game started less than an hour ago
		return ""
	}

	if contact.Language == "fr-FR" {
		if contact.GameMaxScore >= 16 {
			return "477339000034193061"
		} else {
			return "477339000034193078"
		}
	} else {
		if contact.GameMaxScore >= 16 {
			return "477339000034553001"
		} else {
			return "477339000034553020"
		}
	}
}

func AsTime(t *time.Time) *Time {
	if t == nil {
		return nil
	}
	return &Time{*t}
}

func getContactsFromQuery(query string) (*modifiedContactResponse, error) {
	b, err := json.Marshal(&selectQuery{Query: query})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://www.zohoapis.eu/crm/v3/coql", bytes.NewReader(b))
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
		return &modifiedContactResponse{}, nil
	}
	if r.StatusCode != http.StatusOK {
		b, err = io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		log.Println(string(b))
		return nil, errors.New(r.Status)
	}
	var res modifiedContactResponse
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func AutoUpdateAllContacts() error {
	last, err := GetLastAutoZoho()
	if err != nil {
		return err
	}
	c, err := getContactsFromQuery(fmt.Sprintf("SELECT id,Modified_Time FROM Contacts WHERE Created_Time > '2025-02-06T00:17:00-00:00' AND Modified_Time > '%s' ORDER BY Modified_Time ASC", last))
	if err != nil {
		return err
	}
	//log.Println(c.Data)
	l := last
	for _, d := range c.Data {
		err = AutoUpdateContact(d.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		l = d.Modified
	}
	if l != last {
		err = SetLastAutoZoho(l)
	}
	return err
}
