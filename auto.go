package zoho

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

/*
func AutoUpdateContact(id string) error {

		c, err := GetContact(id)
		if err != nil {
			return err
		}

		lang := c.Language
		if lang == "" {
			lang = "fr-FR"
		}

		phone, otherPhone := normalizePhone(c.Phone, c.OtherPhone)

		did, err := FindDealByContactID(id)
		if err != nil {
			return err
		}
		if did != "" {

			err = UpdateDealLeadSource(did, c.LeadSource)
			if err != nil {
				return err
			}
		}

		emails, err := GetContactEmails(id)
		if err != nil {
			return err
		}

		var lastSent *time.Time
		var lastReceived *time.Time
		var sent int
		var received int

		for _, m := range emails {
			if m.Sent {
				sent++
				if lastSent == nil || lastSent.Before(m.Time) {
					lastSent = &m.Time
				}
			} else {
				received++
				if lastReceived == nil || lastReceived.Before(m.Time) {
					lastSent = &m.Time
				}
			}
		}

		notes, lastNote, err := GetContactNotesCount(id)
		if err != nil {
			return err
		}

		p := &phone
		if phone == "" {
			p = nil
		}
		po := &otherPhone
		if otherPhone == "" {
			po = nil
		}
		return updateAutoContact(autoContact{
			ID:         id,
			Language:   lang,
			Phone:      p,
			OtherPhone: po,

			LastEmailSent:     AsTime(lastSent),
			LastEmailReceived: AsTime(lastReceived),
			LastNote:          AsTime(lastNote),
			EmailsSent:        &sent,
			EmailsReceived:    &received,
			NotesCount:        &notes,
			LastUpdate:        &Time{time.Now()},
		})
	}
*/
func AutoUpdateContact(id string) error {
	c, err := GetContact(id)
	if err != nil {
		return err
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
				c.Pipeline = "2026-2027"
			case "Première":
				c.Pipeline = "2025-2026"
			case "Terminale", "BAC", "BAC+1", "BAC+2", "BAC+3":
				c.Pipeline = "2024-2025"
			}
		case "Relative":
			switch c.StudyLevel {
			case "Seconde":
				c.Pipeline = "Relative 2026-2027"
			case "Première":
				c.Pipeline = "Relative 2025-2026"
			case "Terminale", "BAC", "BAC+1", "BAC+2", "BAC+3":
				c.Pipeline = "Relative 2024-2025"
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
	}

	phone, otherPhone := normalizePhone(c.Phone, c.OtherPhone)

	p := &phone
	if phone == "" {
		p = nil
	}
	po := &otherPhone
	if otherPhone == "" {
		po = nil
	}
	a := autoContact{
		ID:         id,
		Language:   lang,
		Phone:      p,
		OtherPhone: po,
		Pipeline:   c.Pipeline,
		LastUpdate: &Time{time.Now()},
	}

	template, short, final := templateToSend(c)
	a.GameTooShortEmailSent = short
	a.GameFinalEmailSent = final
	err = updateAutoContact(a)
	if err != nil {
		log.Println(a)
		return err
	}

	if template != "" {
		log.Println("Sending", template, "to", c.Email)
		err = SendMail("Natacha BOEZ", "natacha.boez@algosup.com", c.FirstName+" "+c.LastName, c.Email, id, template)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func templateToSend(contact *GetContactItem) (template string, short *bool, final *bool) {
	t := true
	f := false

	if contact.GameStart.IsZero() {
		final = &f
		return
	}

	if contact.GameStart.Before(time.Date(2023, 5, 19, 15, 0, 0, 0, time.UTC)) {
		// Game started before new messages where in place

		final = &t
		return
	}

	if contact.GameFinalEmailSent {
		return
	}

	if contact.GameStart.After(time.Now().Add(-time.Hour)) {
		// Game started less than an hour ago
		return
	}

	if contact.GameScore30 >= 18 {
		// Good score
		if contact.Type == "Relative" {
			if contact.Language == "fr-FR" {
				return "477339000003757016", short, &t
			} else {
				return "477339000006666221", short, &t
			}
		}
		if contact.Type == "Prospect" {
			if contact.Pipeline == "2024-2025" {
				if contact.Language == "fr-FR" {
					return "477339000002799418", short, &t
				} else {
					return "477339000006666202", short, &t
				}
			}
			if contact.Pipeline == "2025-2026" {
				if contact.Language == "fr-FR" {
					return "477339000002907148", short, &t
				} else {
					return "477339000006666174", short, &t
				}
			}
			if contact.Pipeline == "2026-2027" {
				if contact.Language == "fr-FR" {
					return "477339000006666309", short, &t
				} else {
					return "477339000006666298", short, &t
				}
			}
		}
		if contact.Pipeline == "" {
			if contact.Language == "fr-FR" {
				return "477339000002907121", short, &t
			} else {
				return "477339000006666263", short, &t
			}
		}

	} else {
		// Bad score
		if contact.Type == "Relative" {
			if contact.Language == "fr-FR" {
				return "477339000003757051", short, &t
			} else {
				return "477339000006666248", short, &t
			}
		}
		if contact.GameDurationMin >= 30 {
			// Long

			if contact.Type == "Prospect" {
				if contact.Language == "fr-FR" {
					return "477339000003757081", short, &t
				} else {
					return "477339000006666320", short, &t
				}
			}
		} else {
			// Short
			if contact.GameTooShortEmailSent {
				// No more than one email for duration
				return
			}
			if contact.Type == "Prospect" {
				if contact.Language == "fr-FR" {
					return "477339000002799429", &t, final
				} else {
					return "477339000006666151", &t, final
				}
			}
		}
	}

	return
}

func AsTime(t *time.Time) *Time {
	if t == nil {
		return nil
	}
	return &Time{*t}
}

func getContactsFromQuery(query string) (*findContactResponse, error) {
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
		return &findContactResponse{}, nil
	}
	if r.StatusCode != http.StatusOK {
		b, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		log.Println(string(b))
		return nil, errors.New(r.Status)
	}
	var res findContactResponse
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

/*
	func AutoUpdateAllContacts() error {
		c, err := getContactsFromQuery("SELECT id FROM Contacts WHERE Last_Update is null ORDER BY Last_Update ASC")
		if err != nil {
			return err
		}
		for _, d := range c.Data {
			err = AutoUpdateContact(d.ID)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		c, err = getContactsFromQuery(fmt.Sprintf("SELECT id FROM Contacts WHERE Last_Update <= '%s' ORDER BY Last_Update ASC", time.Now().Add(-24*time.Hour).Format("2006-01-02T15:04:05-07:00")))
		if err != nil {
			return err
		}
		for _, d := range c.Data {
			err = AutoUpdateContact(d.ID)
			if err != nil {
				log.Println(err)
				continue
			}
		}
		return nil
	}
*/
func AutoUpdateAllContacts() error {
	c, err := getContactsFromQuery("SELECT id FROM Contacts WHERE Last_Update is null ORDER BY Last_Update ASC")
	if err != nil {
		return err
	}
	for _, d := range c.Data {
		err = AutoUpdateContact(d.ID)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	c, err = getContactsFromQuery(fmt.Sprintf("SELECT id FROM Contacts WHERE Last_Update <= '%s' ORDER BY Last_Update ASC", time.Now().Add(-24*time.Hour).Format("2006-01-02T15:04:05-07:00")))
	if err != nil {
		return err
	}
	for _, d := range c.Data {
		err = AutoUpdateContact(d.ID)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
