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

	if phone == "" {
		return phone, otherPhone
	}
	if phone[0:3] == "+33" {
		return phone, otherPhone
	}
	if phone[0] == '+' {
		return "", phone // Move foreign country to other phone
	}

	if len(phone) == 10 && (phone[0:2] == "06" || phone[0:2] == "07") {
		return "+33" + phone[1:], otherPhone
	}

	if len(phone) == 9 && (phone[0] == '6' || phone[0] == '7') {
		return "+33" + phone, otherPhone
	}

	return "", phone
}

func AutoUpdateContact(id string) error {

	c, err := GetContact(id)
	if err != nil {
		return err
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
	return updateAutoContact(autoContact{
		ID:         id,
		Phone:      phone,
		OtherPhone: otherPhone,

		LastEmailSent:     AsTime(lastSent),
		LastEmailReceived: AsTime(lastReceived),
		LastNote:          AsTime(lastNote),
		EmailsSent:        &sent,
		EmailsReceived:    &received,
		NotesCount:        &notes,
		LastUpdate:        &Time{time.Now()},
	})
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
